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

		dbTx, err := r.repository.GetTx(ctx)
		if err != nil {
			slog.Error("Error beginning db tx",
				"application", app.Name,
				"address", address,
				"error", err,
			)
			continue
		}
		defer dbTx.Rollback(ctx)

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
						reason := "Received inputs for an epoch that was not open. Should never happen"

						slog.Error(reason,
							"application", app.Name,
							"address", address,
							"epoch_index", currentEpoch.Index,
							"status", currentEpoch.Status,
						)
						err := r.repository.UpdateApplicationState(ctx, app.ID, ApplicationState_Inoperable, &reason)

						if err != nil {
							slog.Error("failed to update application state to inoperable", "application", app.Name, "err", err)
						}
						return errors.New(reason)
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
			combinedIndex, err := r.repository.GetInputIndexWithTx(ctx, dbTx, address.Hex())
			if err != nil {
				slog.Error("evmreader: failed to read index", "app", address, "error", err)
				return errors.New("evmreader: failed to read index")
			}
			if combinedIndex != input.Index && r.shouldModifyIndex == true {
				slog.Info("evmreader: Overriding input index", "onchain-index", input.Index, "new-index", combinedIndex)
				input.Index = combinedIndex
				modifiedRawData, err := r.modifyIndexInRaw(input.RawData, combinedIndex)
				if err != nil {
					slog.Error("evmreader: failed to modify index", "app", address, "error", err)
					return errors.New("evmreader: failed to modify index")
				}
				input.RawData = modifiedRawData
			}
			// update input index
			err = r.repository.UpdateInputIndexWithTx(ctx, dbTx, address.Hex())
			if err != nil {
				slog.Error("failed to update index", "app", address, "error", err)
				return errors.New("evmreader: failed to update index")
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
			dbTx,
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
			return errors.New("evmreader: Error storing inputs and epochs")
		}
		// Commit transaction
		err = dbTx.Commit(ctx)
		if err != nil {
			slog.Error("could not commit db tx", "err", err)
			return errors.New("evmreader: could not commit db tx")
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

	// Update LastInputCheckBlock for applications that didn't have any inputs
	// (for apps with inputs, LastInputCheckBlock is already updated in CreateEpochsAndInputs)
	appsToUpdate := []int64{}
	// Find applications that didn't have any inputs in appInputsMap
	for _, app := range apps {
		appAddress := app.IApplicationAddress
		// If the app doesn't have any inputs in the map or has an empty slice
		if inputs, exists := appInputsMap[appAddress]; !exists || len(inputs) == 0 {
			appsToUpdate = append(appsToUpdate, app.ID)
		}
	}
	// Update LastInputCheckBlock for applications without inputs
	if len(appsToUpdate) > 0 {
		err := r.repository.UpdateEventLastCheckBlock(ctx, appsToUpdate, MonitoredEvent_InputAdded, mostRecentBlockNumber)
		if err != nil {
			slog.Error("Failed to update LastInputCheckBlock for applications without inputs",
				"app_ids", appsToUpdate,
				"block_number", mostRecentBlockNumber,
				"error", err,
			)
			// We don't return an error here as we've already processed the inputs
			// and this is just an update to the last check block
		} else {
			slog.Debug("Updated LastInputCheckBlock for applications without inputs",
				"app_ids", appsToUpdate,
				"block_number", mostRecentBlockNumber,
			)
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

	inputSource := apps[0].InputSourceAdapter
	opts := bind.FilterOpts{
		Context: ctx,
		Start:   startBlock,
		End:     &endBlock,
	}
	inputsEvents, err := inputSource.RetrieveInputs(&opts, appsAddresses, nil)
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
	return app.LastInputCheckBlock
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
