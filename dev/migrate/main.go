// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package main

import (
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/cartesi/rollups-espresso-reader/internal/config"
	"github.com/cartesi/rollups-espresso-reader/internal/repository/postgres/schema"
)

func main() {
	var s *schema.Schema
	var err error

	postgresEndpoint := config.GetDatabaseConnection() + "&x-migrations-table=espresso_schema_migrations"

	uri, err := url.Parse(postgresEndpoint)
	if err == nil {
		uri.User = nil
	} else {
		slog.Error("Failed to parse PostgresEndpoint.", "error", err)
		os.Exit(1)
	}

	for i := 0; i < 5; i++ {
		s, err = schema.New(postgresEndpoint)
		if err == nil {
			break
		}
		slog.Warn("Connection to database failed. Trying again.", "PostgresEndpoint", uri.String())
		if i == 4 {
			slog.Error("Failed to connect to database.", "error", err)
			os.Exit(1)
		}
		time.Sleep(5 * time.Second) // wait before retrying
	}
	defer s.Close()

	err = s.Upgrade()
	if err != nil {
		slog.Error("Error while upgrading database schema", "error", err)
		os.Exit(1)
	}

	version, err := s.ValidateVersion()
	if err != nil {
		slog.Error("Error while validating database schema version", "error", err)
		os.Exit(1)
	}

	slog.Info("Database Schema successfully Updated.", "version", version)
}
