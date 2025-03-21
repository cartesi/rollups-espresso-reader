// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/cartesi/rollups-espresso-reader/internal/config"
	"github.com/cartesi/rollups-espresso-reader/internal/espressoreader"
	"github.com/cartesi/rollups-espresso-reader/internal/evmreader"
	"github.com/cartesi/rollups-espresso-reader/internal/repository"
	"github.com/cartesi/rollups-espresso-reader/internal/repository/factory"
	"github.com/cartesi/rollups-espresso-reader/internal/services/retry"
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

type loadNodeConfigArgs struct {
	ctx       context.Context
	database  repository.Repository
	configKey string
}

func run(cmd *cobra.Command, args []string) {
	startTime := time.Now()

	ctx := cmd.Context()

	c := config.FromEnv()

	// setup log
	startup.ConfigLogs(c.LogLevel, c.LogPrettyEnabled)

	slog.Info("Starting the Cartesi Rollups Node Espresso Reader", "config", c)

	// connect to database
	database, err := factory.NewRepositoryFromConnectionString(ctx, c.PostgresEndpoint.Value)
	if err != nil {
		slog.Error("Espresso Reader couldn't connect to the database", "error", err)
		os.Exit(1)
	}

	// load node configuration
	var config evmreader.PersistentConfig
	config, err = retry.CallFunctionWithRetryPolicy(
		loadNodeConfig,
		loadNodeConfigArgs{
			ctx:       ctx,
			database:  database,
			configKey: evmreader.EvmReaderConfigKey,
		},
		c.MaxRetries,
		c.MaxDelay,
		"Main::LoadNodeConfig",
	)
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// create Espresso Reader Service
	service := espressoreader.NewEspressoReaderService(
		c.BlockchainHttpEndpoint.Value,
		c.BlockchainWsEndpoint.Value,
		database,
		c.EspressoBaseUrl,
		c.EspressoStartingBlock,
		c.EspressoNamespace,
		config.ChainID,
		config.InputBoxDeploymentBlock,
		c.EspressoServiceEndpoint,
		c.MaxRetries,
		c.MaxDelay,
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

func loadNodeConfig(args loadNodeConfigArgs) (evmreader.PersistentConfig, error) {
	config, _ := repository.LoadNodeConfig[evmreader.PersistentConfig](args.ctx, args.database, args.configKey)
	return config.Value, nil
}

func main() {
	err := Cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
