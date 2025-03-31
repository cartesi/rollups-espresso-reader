// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// This package contains functions to help using the Go-ethereum library.
// It is not the objective of this package to replace or hide Go-ethereum.
package ethutil

import (
	"context"
	"fmt"

	"github.com/cartesi/rollups-espresso-reader/pkg/contracts/iapplication"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

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
