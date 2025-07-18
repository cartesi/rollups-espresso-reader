// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package factory

import (
	"context"
	"fmt"
	"strings"
	"time"

	. "github.com/cartesi/rollups-espresso-reader/internal/repository"
	"github.com/cartesi/rollups-espresso-reader/internal/repository/postgres"
)

// NewRepositoryFromConnectionString chooses the backend based on the connection string.
// For instance:
//   - "postgres://user:pass@localhost/dbname" => Postgres
//   - "sqlite://some/path.db" => SQLite
//
// Then it initializes the repo, runs migrations, and returns it.
func NewRepositoryFromConnectionString(ctx context.Context, conn string) (Repository, error) {
	lowerConn := strings.ToLower(conn)
	switch {
	case strings.HasPrefix(lowerConn, "postgres://"), strings.HasPrefix(lowerConn, "postgresql://"):
		return newPostgresRepository(ctx, conn)
	default:
		return nil, fmt.Errorf("unrecognized connection string format: %s", conn)
	}
}

func newPostgresRepository(ctx context.Context, conn string) (Repository, error) {
	pgRepo, err := postgres.NewPostgresRepository(ctx, conn, 5, 3*time.Second) // FIXME: get from config
	if err != nil {
		return nil, err
	}

	return pgRepo, nil
}
