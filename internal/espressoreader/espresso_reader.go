// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package espressoreader

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/cartesi/rollups-espresso-reader/internal/evmreader"
	"github.com/cartesi/rollups-espresso-reader/internal/model"
	"github.com/cartesi/rollups-espresso-reader/internal/repository"
	"github.com/cartesi/rollups-espresso-reader/internal/services/retry"

	"github.com/EspressoSystems/espresso-sequencer-go/client"
	"github.com/EspressoSystems/espresso-sequencer-go/types"
	espresso "github.com/EspressoSystems/espresso-sequencer-go/types/common"
	"github.com/ethereum/go-ethereum/common"
)

type EspressoClient interface {
	FetchVidCommonByHeight(ctx context.Context, blockHeight uint64) (espresso.VidCommon, error)
	FetchLatestBlockHeight(ctx context.Context) (uint64, error)
	FetchHeaderByHeight(ctx context.Context, blockHeight uint64) (types.HeaderImpl, error)
	FetchHeadersByRange(ctx context.Context, from uint64, until uint64) ([]types.HeaderImpl, error)
	FetchTransactionByHash(ctx context.Context, hash *types.TaggedBase64) (types.TransactionQueryData, error)
	FetchBlockMerkleProof(ctx context.Context, rootHeight uint64, hotshotHeight uint64) (types.HotShotBlockMerkleProof, error)
	FetchTransactionsInBlock(ctx context.Context, blockHeight uint64, namespace uint64) (client.TransactionsInBlock, error)
	SubmitTransaction(ctx context.Context, tx types.Transaction) (*types.TaggedBase64, error)
}

type EspressoHelperInterface interface {
	getL1FinalizedHeight(ctx context.Context, espressoBlockHeight uint64, delay uint64, url string) (uint64, uint64)
	readEspressoHeadersByRange(ctx context.Context, from uint64, until uint64, delay uint64, url string) string
	getNSTableByRange(ctx context.Context, from uint64, until uint64, delay uint64, url string) (string, error)
	extractNS(nsTable []byte) []uint32
}

type EspressoReader struct {
	url                     string
	client                  EspressoClient
	espressoHelper          EspressoHelperInterface
	startingBlock           uint64
	namespace               uint64
	repository              repository.Repository
	evmReader               *evmreader.EvmReader
	chainId                 uint64
	inputBoxDeploymentBlock uint64
	maxRetries              uint64
	maxDelay                uint64
}

func NewEspressoReader(url string, startingBlock uint64, namespace uint64, repository repository.Repository, evmReader *evmreader.EvmReader, chainId uint64, inputBoxDeploymentBlock uint64, maxRetries uint64, maxDelay uint64) EspressoReader {
	client := client.NewClient(url)
	espressoHelper := &EspressoHelper{}
	return EspressoReader{url: url, client: client, espressoHelper: espressoHelper, startingBlock: startingBlock, namespace: namespace, repository: repository, evmReader: evmReader, chainId: chainId, inputBoxDeploymentBlock: inputBoxDeploymentBlock, maxRetries: maxRetries, maxDelay: maxDelay}
}

func (e *EspressoReader) Run(ctx context.Context, ready chan<- struct{}) error {
	ready <- struct{}{}

	for {
		select {
		case <-ctx.Done():
			slog.Info("exiting espresso reader")
			return ctx.Err()
		default:
			// fetch latest espresso block height
			latestBlockHeight, err := retry.CallFunctionWithRetryPolicy(
				e.client.FetchLatestBlockHeight,
				ctx,
				e.maxRetries,
				time.Duration(e.maxDelay),
				"EspressoReader::FetchLatestBlockHeight",
			)
			if err != nil {
				slog.Error("failed fetching latest espresso block height", "error", err)
				return err
			}
			slog.Debug("Espresso:", "latestBlockHeight", latestBlockHeight)

			apps := e.getAppsForEvmReader(ctx)
			if len(apps) > 0 {
				for _, app := range apps {
					lastProcessedEspressoBlock, err := e.repository.GetLastProcessedEspressoBlock(ctx, app.Application.IApplicationAddress.Hex())
					if err != nil {
						slog.Error("failed reading lastProcessedEspressoBlock", "error", err)
						continue
					}
					lastProcessedL1Block := app.Application.LastProcessedBlock
					appAddress := app.Application.IApplicationAddress
					if lastProcessedL1Block < e.inputBoxDeploymentBlock {
						lastProcessedL1Block = e.inputBoxDeploymentBlock - 1
					}
					if lastProcessedEspressoBlock == 0 {
						if e.startingBlock != 0 {
							lastProcessedEspressoBlock = e.startingBlock - 1
						} else {
							lastProcessedEspressoBlock = latestBlockHeight - 1
						}
					}
					// bootstrap if there are more than 100 blocks to catch up
					if latestBlockHeight-lastProcessedEspressoBlock > 100 {
						// bootstrap
						slog.Debug("bootstrapping:", "app", appAddress, "from-block", lastProcessedEspressoBlock+1, "to-block", latestBlockHeight)
						err = e.bootstrap(ctx, app, lastProcessedEspressoBlock, latestBlockHeight, lastProcessedL1Block)
						if err != nil {
							slog.Error("failed reading inputs", "error", err)
							continue
						}

						// update lastProcessedEspressoBlock in db
						err = e.repository.UpdateLastProcessedEspressoBlock(ctx, app.Application.IApplicationAddress.Hex(), latestBlockHeight)
						if err != nil {
							slog.Error("failed updating last processed espresso block", "error", err)
						}
					} else {
						// in sync. Process espresso blocks one-by-one
						currentBlockHeight := lastProcessedEspressoBlock + 1
						for ; currentBlockHeight <= latestBlockHeight; currentBlockHeight++ {
							slog.Debug("Espresso:", "app", appAddress, "currentBlockHeight", currentBlockHeight)
							//** read base layer **//
							var l1FinalizedTimestamp uint64
							lastProcessedL1Block, l1FinalizedTimestamp = e.readL1(ctx, app, currentBlockHeight, lastProcessedL1Block)
							//** read espresso **//
							e.readEspresso(ctx, app, currentBlockHeight, lastProcessedL1Block, l1FinalizedTimestamp)

							// update lastProcessedEspressoBlock in db
							err = e.repository.UpdateLastProcessedEspressoBlock(ctx, app.Application.IApplicationAddress.Hex(), latestBlockHeight)
							if err != nil {
								slog.Error("failed updating last processed espresso block", "error", err)
							}
						}
					}

				}
			}

			// take a break :)
			var delay time.Duration = 1000
			time.Sleep(delay * time.Millisecond)
		}
	}
}

func (e *EspressoReader) bootstrap(ctx context.Context, app evmreader.TypeExportApplication, lastProcessedEspressoBlock uint64, latestBlockHeight uint64, l1FinalizedHeight uint64) error {
	var l1FinalizedTimestamp uint64
	var nsTables []string
	batchStartingBlock := lastProcessedEspressoBlock + 1
	batchLimit := uint64(100)
	for latestBlockHeight >= batchStartingBlock {
		select {
		case <-ctx.Done():
			slog.Info("exiting espresso reader")
			return ctx.Err()
		default:
			var batchEndingBlock uint64
			if batchStartingBlock+batchLimit > latestBlockHeight+1 {
				batchEndingBlock = latestBlockHeight + 1
			} else {
				batchEndingBlock = batchStartingBlock + batchLimit
			}
			nsTable, err := e.espressoHelper.getNSTableByRange(ctx, batchStartingBlock, batchEndingBlock, e.maxDelay, e.url)
			if err != nil {
				return err
			}
			nsTableBytes := []byte(nsTable)
			err = json.Unmarshal(nsTableBytes, &nsTables)
			if err != nil {
				slog.Error("failed fetching ns tables", "error", err, "ns table", nsTables)
			} else {
				for index, nsTable := range nsTables {
					nsTableBytes, _ := base64.StdEncoding.DecodeString(nsTable)
					ns := e.espressoHelper.extractNS(nsTableBytes)
					if slices.Contains(ns, uint32(e.namespace)) {
						currentEspressoBlock := batchStartingBlock + uint64(index)
						slog.Debug("found namespace contained in", "block", currentEspressoBlock)
						l1FinalizedHeight, l1FinalizedTimestamp = e.readL1(ctx, app, currentEspressoBlock, l1FinalizedHeight)
						e.readEspresso(ctx, app, currentEspressoBlock, l1FinalizedHeight, l1FinalizedTimestamp)
					}
				}
			}
			// update loop var
			batchStartingBlock += batchLimit
		}
	}
	return nil
}

func (e *EspressoReader) readL1(ctx context.Context, app evmreader.TypeExportApplication, currentBlockHeight uint64, lastProcessedL1Block uint64) (uint64, uint64) {
	l1FinalizedLatestHeight, l1FinalizedTimestamp := e.espressoHelper.getL1FinalizedHeight(ctx, currentBlockHeight, e.maxDelay, e.url)
	// read L1 if there might be update
	if l1FinalizedLatestHeight > lastProcessedL1Block {
		slog.Debug("L1 finalized", "app", app.Application.IApplicationAddress, "from", lastProcessedL1Block, "to", l1FinalizedLatestHeight)

		var apps []evmreader.TypeExportApplication
		apps = append(apps, app) // make app into 1-element array

		// start reading from the block after the prev height
		e.evmReader.ReadAndStoreInputs(ctx, lastProcessedL1Block, l1FinalizedLatestHeight, apps)
		// check for claim status and output execution
		// e.evmReader.CheckForClaimStatus(ctx, apps, l1FinalizedLatestHeight) // checked by the node
		// e.evmReader.CheckForOutputExecution(ctx, apps, l1FinalizedLatestHeight) // checked by the node
	}
	return l1FinalizedLatestHeight, l1FinalizedTimestamp
}

type EspressoInput struct {
	app       string
	data      []byte
	nonce     uint64
	msgSender common.Address
	sigHash   string
}

func (e *EspressoReader) readEspressoInput(transaction espresso.Bytes) (*EspressoInput, error) {
	msgSender, typedData, sigHash, err := ExtractSigAndData(string(transaction))
	if err != nil {
		slog.Error("failed to extract espresso tx", "error", err)
		return nil, err
	}

	var nonce uint64
	nonceFloat64, ok := typedData.Message["nonce"].(float64)
	if !ok {
		nonceStr, ok := typedData.Message["nonce"].(string)
		if !ok {
			nonceU64, err := strconv.ParseUint(nonceStr, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to cast nonce to float: %s", err)
			}
			nonce = nonceU64
		}
	} else {
		nonce = uint64(nonceFloat64)
	}
	payload, ok := typedData.Message["data"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to cast data to string")
	}
	appAddressStr, ok := typedData.Message["app"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to cast app address to string")
	}

	slog.Info("Espresso input", "msgSender", msgSender, "nonce", nonce, "payload", payload, "appAddrss", appAddressStr, "tx-id", sigHash)

	payloadBytes := []byte(payload)
	if strings.HasPrefix(payload, "0x") {
		payload = payload[2:] // remove 0x
		payloadBytes, err = hex.DecodeString(payload)
		if err != nil {
			slog.Error("failed to decode hex string", "error", err)
			return nil, err
		}
	}
	espressoInput := EspressoInput{
		app:       strings.ToLower(appAddressStr),
		data:      payloadBytes,
		nonce:     nonce,
		msgSender: msgSender,
		sigHash:   sigHash,
	}

	return &espressoInput, nil
}

func (e *EspressoReader) encodeEvmAdvance(espressoInput *EspressoInput, chainId *big.Int, app common.Address,
	indexUint64 uint64, l1FinalizedLatestHeightBig *big.Int, l1FinalizedTimestampBig *big.Int,
	prevRandao *big.Int,
) ([]byte, error) {
	payloadBytes := espressoInput.data
	index := &big.Int{}
	index.SetUint64(indexUint64)
	payloadAbi, err := e.evmReader.IOAbi.Pack("EvmAdvance", chainId, app, espressoInput.msgSender, l1FinalizedLatestHeightBig, l1FinalizedTimestampBig, prevRandao, index, payloadBytes)
	if err != nil {
		slog.Error("failed to abi encode", "error", err)
		return nil, err
	}
	return payloadAbi, nil
}

func toBigInt(value uint64) *big.Int {
	bigInt := &big.Int{}
	bigInt.SetUint64(value)
	return bigInt
}

func (e *EspressoReader) buildInput(ctx context.Context,
	espressoInput *EspressoInput, chainId *big.Int,
	app common.Address,
	l1FinalizedLatestHeightBig *big.Int,
	l1FinalizedTimestampBig *big.Int,
	prevRandao *big.Int,
	l1FinalizedLatestHeight uint64,
) (*model.Input, error) {
	indexUint64, err := e.repository.GetInputIndex(ctx, espressoInput.app)
	if err != nil {
		slog.Error("failed to read index", "app", espressoInput.app, "error", err)
	}

	// abi encode payload
	payloadAbi, err := e.encodeEvmAdvance(espressoInput, chainId, app, indexUint64, l1FinalizedLatestHeightBig, l1FinalizedTimestampBig, prevRandao)
	if err != nil {
		slog.Error("failed to abi encode", "error", err)
		return nil, err
	}
	// build input
	sigHashHexBytes, err := hex.DecodeString(espressoInput.sigHash[2:])
	if err != nil {
		slog.Error("could not obtain bytes for tx-id", "err", err)
		return nil, err
	}
	input := model.Input{
		Index:                indexUint64,
		Status:               model.InputCompletionStatus_None,
		RawData:              payloadAbi,
		BlockNumber:          l1FinalizedLatestHeight,
		TransactionReference: common.BytesToHash(sigHashHexBytes),
	}
	return &input, nil
}

func (e *EspressoReader) isNonceValid(ctx context.Context, espressoInput *EspressoInput) (bool, error) {
	nonceInDb, err := e.repository.GetEspressoNonce(ctx, espressoInput.msgSender.Hex(), espressoInput.app)
	if err != nil {
		// Here, I maintain the same behavior, but it's better to handle the error by retrying or halting.
		slog.Error("failed to get espresso nonce from db", "error", err)
		return false, err
	}
	if espressoInput.nonce != nonceInDb {
		slog.Error("Espresso nonce is incorrect. May be a duplicate tx", "nonce from espresso", espressoInput.nonce, "nonce in db", nonceInDb)
		return false, nil
	}
	return true, nil
}

func (e *EspressoReader) findOrBuildNewEpoch(ctx context.Context, appEvmType evmreader.TypeExportApplication, readingAppAddress string, l1FinalizedLatestHeight uint64) (*model.Epoch, error) {
	// if current epoch is not nil, the epoch is open
	// espresso inputs do not close epoch
	epochIndex := evmreader.CalculateEpochIndex(appEvmType.EpochLength, l1FinalizedLatestHeight)
	currentEpoch, err := e.repository.GetEpoch(ctx, readingAppAddress, epochIndex)
	if err != nil {
		slog.Error("could not obtain current epoch", "err", err)
		return nil, err
	}
	if currentEpoch == nil {
		currentEpoch = &model.Epoch{
			Index:      epochIndex,
			FirstBlock: epochIndex * appEvmType.EpochLength,
			LastBlock:  (epochIndex * appEvmType.EpochLength) + appEvmType.EpochLength - 1,
			Status:     model.EpochStatus_Open,
		}
	}
	return currentEpoch, nil
}

func (e *EspressoReader) readEspresso(ctx context.Context, appEvmType evmreader.TypeExportApplication, currentBlockHeight uint64, l1FinalizedLatestHeight uint64, l1FinalizedTimestamp uint64) {
	app := appEvmType.Application.IApplicationAddress
	transactions, err := e.client.FetchTransactionsInBlock(ctx, currentBlockHeight, e.namespace)
	if err != nil {
		slog.Error("failed fetching espresso tx", "error", err)
		return
	}
	readingAppAddress := strings.ToLower(app.Hex())
	chainId := toBigInt(e.chainId)
	l1FinalizedLatestHeightBig := toBigInt(l1FinalizedLatestHeight)
	l1FinalizedTimestampBig := toBigInt(l1FinalizedTimestamp)
	prevRandao, err := readPrevRandao(ctx, l1FinalizedLatestHeight, e.evmReader.GetEthClient())
	if err != nil {
		slog.Error("failed to read prevrandao", "error", err)
	}

	numTx := len(transactions.Transactions)
	for i := 0; i < numTx; i++ {
		transaction := transactions.Transactions[i]
		espressoInput, err := e.readEspressoInput(transaction)
		if err != nil {
			slog.Error("failed to extract espresso tx", "error", err)
			continue
		}
		if espressoInput.app != readingAppAddress {
			slog.Debug("skipping... Looking for txs for", "app", app, "found tx for app", espressoInput.app)
			continue
		}

		nonceValid, err := e.isNonceValid(ctx, espressoInput)
		if !nonceValid || err != nil {
			continue
		}

		input, err := e.buildInput(ctx, espressoInput, chainId, app, l1FinalizedLatestHeightBig, l1FinalizedTimestampBig, prevRandao, l1FinalizedLatestHeight)
		if err != nil {
			slog.Error("failed to build input", "error", err)
			continue
		}

		currentEpoch, err := e.findOrBuildNewEpoch(ctx, appEvmType, readingAppAddress, l1FinalizedLatestHeight)
		if err != nil {
			slog.Error("could not obtain current epoch", "err", err)
			continue
		}

		// build epochInputMap
		// Initialize epochs inputs map
		// -> It always has only one input
		var epochInputMap = make(map[*model.Epoch][]*model.Input)
		currentInputs, ok := epochInputMap[currentEpoch]
		if !ok {
			currentInputs = []*model.Input{}
		}
		epochInputMap[currentEpoch] = append(currentInputs, input)

		// Store everything
		// future optimization: bundle tx by address to fully utilize `epochInputMap`
		// -> begin db transaction
		err = e.repository.CreateEpochsAndInputs(
			ctx,
			readingAppAddress,
			epochInputMap,
			l1FinalizedLatestHeight,
		)
		if err != nil {
			slog.Error("could not store Espresso input", "err", err)
			continue
		}

		// update nonce
		err = e.repository.UpdateEspressoNonce(ctx, espressoInput.msgSender.Hex(), readingAppAddress)
		if err != nil {
			slog.Error("!!!could not update Espresso nonce!!!", "err", err)
			continue
		}
		// update input index
		err = e.repository.UpdateInputIndex(ctx, readingAppAddress)
		if err != nil {
			slog.Error("failed to update index", "app", readingAppAddress, "error", err)
		}
		// -> end db transaction
	}

}

//////// evm reader related ////////

func (e *EspressoReader) getAppsForEvmReader(ctx context.Context) []evmreader.TypeExportApplication {
	// Get All Applications
	runningApps, err := e.evmReader.GetAllRunningApplications(ctx)
	if err != nil {
		slog.Error("Error retrieving running applications",
			"error",
			err,
		)
	}

	// Build Contracts
	var apps []evmreader.TypeExportApplication
	for _, app := range runningApps {
		applicationContract, consensusContract, err := e.evmReader.GetAppContracts(*app)
		if err != nil {
			slog.Error("Error retrieving application contracts", "app", app, "error", err)
			continue
		}
		apps = append(apps, evmreader.TypeExportApplication{Application: *app,
			ApplicationContract: applicationContract,
			ConsensusContract:   consensusContract})
	}

	if len(apps) == 0 {
		slog.Info("No correctly configured applications running")
	}

	return apps
}

func readPrevRandao(ctx context.Context, l1FinalizedLatestHeight uint64, client *evmreader.EthClient) (*big.Int, error) {
	header, err := (*client).HeaderByNumber(ctx, big.NewInt(int64(l1FinalizedLatestHeight)))
	if err != nil {
		return &big.Int{}, fmt.Errorf("espresso read block header error: %w", err)
	}
	prevRandao := header.MixDigest.Big()
	slog.Debug("readPrevRandao", "prevRandao", prevRandao, "blockNumber", l1FinalizedLatestHeight)
	return prevRandao, nil
}
