// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package evmreader

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"

	. "github.com/cartesi/rollups-espresso-reader/internal/model"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type TypeExportApplication = application

// ReadAndStoreInputs reads, inputs from the InputSource given specific filter options, indexes
// them into epochs and store the indexed inputs and epochs
func (r *EvmReader) ReadAndStoreInputs(
	ctx context.Context,
	lastProcessedBlock uint64,
	mostRecentBlockNumber uint64,
	apps []TypeExportApplication,
) error {

	if len(apps) == 0 {
		slog.Warn("No valid running applications")
		return nil
	}

	// Retrieve Inputs from blockchain
	nextSearchBlock := lastProcessedBlock + 1
	appInputsMap, err := r.readInputsFromBlockchain(ctx, apps, nextSearchBlock, mostRecentBlockNumber)
	if err != nil {
		return fmt.Errorf("failed to read inputs from block %v to block %v. %w",
			nextSearchBlock,
			mostRecentBlockNumber,
			err)
	}

	addrToApp := mapAddressToApp(apps)

	// Index Inputs into epochs and handle epoch finalization
	for address, inputs := range appInputsMap {

		app, exists := addrToApp[address]
		if !exists {
			slog.Error("Application address on input not found",
				"address", address)
			continue
		}
		epochLength := app.EpochLength

		// Retrieves last open epoch from DB
		currentEpoch, err := r.repository.GetEpoch(ctx, address.String(), CalculateEpochIndex(epochLength, lastProcessedBlock))
		if err != nil {
			slog.Error("Error retrieving existing current epoch",
				"application", app.Name,
				"address", address,
				"error", err,
			)
			continue
		}

		// Initialize epochs inputs map
		var epochInputMap = make(map[*Epoch][]*Input)
		// Index Inputs into epochs
		for _, input := range inputs {

			inputEpochIndex := CalculateEpochIndex(epochLength, input.BlockNumber)

			// If input belongs into a new epoch, close the previous known one
			if currentEpoch != nil {
				slog.Debug("Current epoch and new input",
					"application", app.Name,
					"address", address,
					"epoch_index", currentEpoch.Index,
					"epoch_status", currentEpoch.Status,
					"input_epoch_index", inputEpochIndex,
				)
				if currentEpoch.Index == inputEpochIndex {
					// Input can only be added to open epochs
					if currentEpoch.Status != EpochStatus_Open {
						slog.Error("Current epoch is not open",
							"application", app.Name,
							"address", address,
							"epoch_index", currentEpoch.Index,
							"status", currentEpoch.Status,
						)
						return fmt.Errorf("Current epoch is not open. Should never happen")
					}
				} else {
					if currentEpoch.Status == EpochStatus_Open {
						currentEpoch.Status = EpochStatus_Closed
						slog.Info("Closing epoch",
							"application", app.Name,
							"address", address,
							"epoch_index", currentEpoch.Index,
							"start", currentEpoch.FirstBlock,
							"end", currentEpoch.LastBlock)
						_, ok := epochInputMap[currentEpoch]
						if !ok {
							epochInputMap[currentEpoch] = []*Input{}
						}
					}
					currentEpoch = nil
				}
			}
			if currentEpoch == nil {
				currentEpoch = &Epoch{
					Index:      inputEpochIndex,
					FirstBlock: inputEpochIndex * epochLength,
					LastBlock:  (inputEpochIndex * epochLength) + epochLength - 1,
					Status:     EpochStatus_Open,
				}
				epochInputMap[currentEpoch] = []*Input{}
			}

			slog.Info("Found new Input",
				"application", app.Name,
				"address", address,
				"index", input.Index,
				"block", input.BlockNumber,
				"epoch_index", inputEpochIndex)

			currentInputs, ok := epochInputMap[currentEpoch]
			if !ok {
				currentInputs = []*Input{}
			}
			// overriding input index
			combinedIndex, err := r.repository.GetInputIndex(ctx, address.Hex())
			if err != nil {
				slog.Error("evmreader: failed to read index", "app", address, "error", err)
			}
			if combinedIndex != input.Index && r.shouldModifyIndex == true {
				slog.Info("evmreader: Overriding input index", "onchain-index", input.Index, "new-index", combinedIndex)
				input.Index = combinedIndex
				modifiedRawData, err := r.modifyIndexInRaw(input.RawData, combinedIndex)
				if err == nil {
					input.RawData = modifiedRawData
				}
			}
			err = r.repository.UpdateInputIndex(ctx, address.Hex())
			if err != nil {
				slog.Error("evmreader: failed to update index", "app", address, "error", err)
			}
			epochInputMap[currentEpoch] = append(currentInputs, input)

		}

		// Indexed all inputs. Check if it is time to close the last epoch
		if currentEpoch != nil && currentEpoch.Status == EpochStatus_Open &&
			mostRecentBlockNumber >= currentEpoch.LastBlock {
			currentEpoch.Status = EpochStatus_Closed
			slog.Info("Closing epoch",
				"application", app.Name,
				"address", address,
				"epoch_index", currentEpoch.Index,
				"start", currentEpoch.FirstBlock,
				"end", currentEpoch.LastBlock)
			// Add to inputMap so it is stored
			_, ok := epochInputMap[currentEpoch]
			if !ok {
				epochInputMap[currentEpoch] = []*Input{}
			}
		}

		err = r.repository.CreateEpochsAndInputs(
			ctx,
			address.String(),
			epochInputMap,
			mostRecentBlockNumber,
		)
		if err != nil {
			slog.Error("Error storing inputs and epochs",
				"application", app.Name,
				"address", address,
				"error", err,
			)
			continue
		}

		// Store everything
		if len(epochInputMap) > 0 {
			slog.Debug("Inputs and epochs stored successfully",
				"application", app.Name,
				"address", address,
				"start-block", nextSearchBlock,
				"end-block", mostRecentBlockNumber,
				"total epochs", len(epochInputMap),
				"total inputs", len(inputs),
			)
		} else {
			slog.Debug("No inputs or epochs to store")
		}

	}

	return nil
}

// readInputsFromBlockchain read the inputs from the blockchain ordered by Input index
func (r *EvmReader) readInputsFromBlockchain(
	ctx context.Context,
	apps []TypeExportApplication,
	startBlock, endBlock uint64,
) (map[common.Address][]*Input, error) {

	// Initialize app input map
	var appInputsMap = make(map[common.Address][]*Input)
	var appsAddresses = []common.Address{}
	for _, app := range apps {
		appInputsMap[app.IApplicationAddress] = []*Input{}
		appsAddresses = append(appsAddresses, app.IApplicationAddress)
	}

	opts := bind.FilterOpts{
		Context: ctx,
		Start:   startBlock,
		End:     &endBlock,
	}
	inputsEvents, err := r.inputSource.RetrieveInputs(&opts, appsAddresses, nil)
	if err != nil {
		return nil, err
	}

	// Order inputs as order is not enforced by RetrieveInputs method nor the APIs
	for _, event := range inputsEvents {
		slog.Debug("Received input",
			"address", event.AppContract,
			"index", event.Index,
			"block", event.Raw.BlockNumber)
		input := &Input{
			Index:                event.Index.Uint64(),
			Status:               InputCompletionStatus_None,
			RawData:              event.Input,
			BlockNumber:          event.Raw.BlockNumber,
			TransactionReference: common.BigToHash(event.Index),
		}

		// Insert Sorted
		appInputsMap[event.AppContract] = insertSorted(
			sortByInputIndex, appInputsMap[event.AppContract], input)
	}
	return appInputsMap, nil
}

// byLastProcessedBlock key extractor function intended to be used with `indexApps` function
func byLastProcessedBlock(app application) uint64 {
	return app.LastProcessedBlock
}

// getEpochLength reads the application epoch length given it's consensus contract
func getEpochLength(consensus ConsensusContract) (uint64, error) {

	epochLengthRaw, err := consensus.GetEpochLength(nil)
	if err != nil {
		return 0, errors.Join(
			fmt.Errorf("error retrieving application epoch length"),
			err,
		)
	}

	return epochLengthRaw.Uint64(), nil
}

func (r *EvmReader) modifyIndexInRaw(rawData []byte, currentIndex uint64) ([]byte, error) {
	// load contract ABI
	abiObject := r.IOAbi
	values, err := abiObject.Methods["EvmAdvance"].Inputs.Unpack(rawData[4:])
	if err != nil {
		slog.Error("Error unpacking abi", "err", err)
		return []byte{}, err
	}

	type EvmAdvance struct {
		ChainId        *big.Int
		AppContract    common.Address
		MsgSender      common.Address
		BlockNumber    *big.Int
		BlockTimestamp *big.Int
		PrevRandao     *big.Int
		Index          *big.Int
		Payload        []byte
	}

	data := EvmAdvance{
		ChainId:        values[0].(*big.Int),
		AppContract:    values[1].(common.Address),
		MsgSender:      values[2].(common.Address),
		BlockNumber:    values[3].(*big.Int),
		BlockTimestamp: values[4].(*big.Int),
		PrevRandao:     values[5].(*big.Int),
		Index:          values[6].(*big.Int),
		Payload:        values[7].([]byte),
	}

	// modify index
	data.Index = big.NewInt(int64(currentIndex))

	// abi encode again
	dataAbi, err := abiObject.Pack("EvmAdvance", data.ChainId, data.AppContract, data.MsgSender, data.BlockNumber, data.BlockTimestamp, data.PrevRandao, data.Index, data.Payload)
	if err != nil {
		slog.Error("failed to abi encode", "error", err)
		return []byte{}, err
	}

	return dataAbi, nil
}
