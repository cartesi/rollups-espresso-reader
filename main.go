// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/cartesi/rollups-espresso-reader/internal/config"
	"github.com/cartesi/rollups-espresso-reader/internal/espressoreader"
	"github.com/cartesi/rollups-espresso-reader/internal/repository/factory"
	"github.com/cartesi/rollups-espresso-reader/internal/services/startup"

	"github.com/spf13/cobra"
)

var (
	// Should be overridden during the final release build with ldflags
	// to contain the actual version number
	buildVersion = "devel"
)

const (
	CMD_NAME = "espresso-sequencer"
)

var Cmd = &cobra.Command{
	Use:   CMD_NAME,
	Short: "Runs Espresso Reader",
	Long:  `Runs Espresso Reader`,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	startTime := time.Now()

	ctx := cmd.Context()

	c := config.FromEnv()

	// setup log
	startup.ConfigLogs(c.LogLevel, c.LogPrettyEnabled)

	slog.Info("Starting the Cartesi Rollups Node Espresso Reader", "config", c)

	database, err := factory.NewRepositoryFromConnectionString(ctx, c.PostgresEndpoint.Value)
	if err != nil {
		slog.Error("EVM Reader couldn't connect to the database", "error", err)
		os.Exit(1)
	}

	_, err = startup.SetupNodeConfig(ctx, database, c)
	if err != nil {
		slog.Error("EVM Reader couldn't connect to the database", "error", err)
		os.Exit(1)
	}

	// create Espresso Reader Service
	service := espressoreader.NewEspressoReaderService(
		c.BlockchainHttpEndpoint.Value,
		c.BlockchainHttpEndpoint.Value,
		database,
		c.EspressoBaseUrl,
		c.EspressoStartingBlock,
		c.EspressoNamespace,
		c.EvmReaderRetryPolicyMaxRetries,
		c.EvmReaderRetryPolicyMaxDelay,
		c.BlockchainID,
		uint64(c.ContractsInputBoxDeploymentBlockNumber),
		c.EspressoServiceEndpoint,
	)

	// logs startup time
	ready := make(chan struct{}, 1)
	go func() {
		select {
		case <-ready:
			duration := time.Since(startTime)
			slog.Info("EVM Reader is ready", "after", duration)
		case <-ctx.Done():
		}
	}()

	// start service
	if err := service.Start(ctx, ready); err != nil {
		slog.Error("Espresso Reader exited with an error", "error", err)
		os.Exit(1)
	}
}

func main() {
	err := Cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
