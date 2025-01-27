// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)
package startup

import (
	"context"
	"log/slog"
	"os"

	"github.com/cartesi/rollups-espresso-reader/internal/model"
	"github.com/cartesi/rollups-espresso-reader/internal/repository"
	"github.com/jackc/pgx/v4"
	"github.com/lmittmann/tint"
)

// Configure the node logs
func ConfigLogs(logLevel slog.Level, logPrettyEnabled bool) {
	opts := &tint.Options{
		Level:     logLevel,
		AddSource: logLevel == slog.LevelDebug,
		// NoColor:    !logPrettyEnabled || !isatty.IsTerminal(os.Stdout.Fd()),
		NoColor:    false,
		TimeFormat: "2006-01-02T15:04:05.000", // RFC3339 with milliseconds and without timezone
	}
	handler := tint.NewHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func SetupNodeConfig(
	ctx context.Context,
	r repository.Repository,
	c *model.NodeConfig[model.NodeConfigValue],
) error {
	_, err := repository.LoadNodeConfig[model.NodeConfigValue](ctx, r, model.BaseConfigKey)
	if err == pgx.ErrNoRows {
		slog.Debug("Initializing node config", "config", c)
		err = repository.SaveNodeConfig(ctx, r, c)
		if err != nil {
			return err
		}
	} else if err == nil {
		slog.Warn("Node was already configured. Using previous persistent config", "config", c.Value)
	} else {
		slog.Error("Could not retrieve persistent config from Database. %w", "error", err)
	}
	return err
}
