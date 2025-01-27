// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)
package startup

import (
	"log/slog"
	"os"

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
