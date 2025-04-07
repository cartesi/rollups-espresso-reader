// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains functions to help using the Go-ethereum library.
// It is not the objective of this package to replace or hide Go-ethereum.
package ethutil

import (
	"context"
	"fmt"
	"math/big"

	"github.com/cartesi/rollups-espresso-reader/pkg/contracts/iapplication"
	"github.com/cartesi/rollups-espresso-reader/pkg/contracts/iinputbox"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Gas limit when sending transactions.
const GasLimit = 30_000_000

// Add input to the input box for the given DApp address.
// This function waits until the transaction is added to a block and return the input index.
func AddInput(
	ctx context.Context,
	client *ethclient.Client,
	transactionOpts *bind.TransactOpts,
	inputBoxAddress common.Address,
	application common.Address,
	input []byte,
) (uint64, uint64, error) {
	if client == nil {
		return 0, 0, fmt.Errorf("AddInput: client is nil")
	}
	inputBox, err := iinputbox.NewIInputBox(inputBoxAddress, client)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to connect to InputBox contract: %v", err)
	}
	receipt, err := sendTransaction(
		ctx, client, transactionOpts, big.NewInt(0), GasLimit,
		func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
			return inputBox.AddInput(txOpts, application, input)
		},
	)
	if err != nil {
		return 0, 0, err
	}
	index, err := getInputIndex(inputBoxAddress, inputBox, receipt)
	return index, receipt.BlockNumber.Uint64(), nil
}

// Get input index in the transaction by looking at the event logs.
func getInputIndex(
	inputBoxAddress common.Address,
	inputBox *iinputbox.IInputBox,
	receipt *types.Receipt,
) (uint64, error) {
	for _, log := range receipt.Logs {
		if log.Address != inputBoxAddress {
			continue
		}
		inputAdded, err := inputBox.ParseInputAdded(*log)
		if err != nil {
			return 0, fmt.Errorf("failed to parse input added event: %v", err)
		}
		// We assume that uint64 will fit all dapp inputs for now
		return inputAdded.Index.Uint64(), nil
	}
	return 0, fmt.Errorf("input index not found")
}

func GetDataAvailability(
	ctx context.Context,
	client *ethclient.Client,
	appAddress common.Address,
) ([]byte, error) {
	if client == nil {
		return nil, fmt.Errorf("get dataAvailability: client is nil")
	}
	app, err := iapplication.NewIApplication(appAddress, client)
	if err != nil {
		return nil, fmt.Errorf("Failed to instantiate contract: %v", err)
	}
	dataAvailability, err := app.GetDataAvailability(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("error retrieving application epoch length: %v", err)
	}
	return dataAvailability, nil
}
