// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package postgres

import (
	"context"
	"database/sql"
	"strings"

	"github.com/cartesi/rollups-espresso-reader/internal/model"
	"github.com/cartesi/rollups-espresso-reader/internal/repository/postgres/db/rollupsdb/public/table"
)

func (r *postgresRepository) CreateNodeConfig(
	ctx context.Context,
	nc *model.NodeConfig,
) error {

	insertStmt := table.NodeConfig.
		INSERT(
			table.NodeConfig.DefaultBlock,
			table.NodeConfig.InputBoxDeploymentBlock,
			table.NodeConfig.InputBoxAddress,
			table.NodeConfig.ChainID,
		).
		VALUES(
			nc.DefaultBlock,
			nc.InputBoxDeploymentBlock,
			strings.ToLower(nc.InputBoxAddress),
			nc.ChainID,
		)

	sqlStr, args := insertStmt.Sql()
	_, err := r.db.Exec(ctx, sqlStr, args...)
	return err
}

func (r *postgresRepository) GetNodeConfig(ctx context.Context) (*model.NodeConfig, error) {
	sel := table.NodeConfig.
		SELECT(
			table.NodeConfig.DefaultBlock,
			table.NodeConfig.InputBoxDeploymentBlock,
			table.NodeConfig.InputBoxAddress,
			table.NodeConfig.ChainID,
			table.NodeConfig.CreatedAt,
			table.NodeConfig.UpdatedAt,
		).
		LIMIT(1)

	sqlStr, args := sel.Sql()
	row := r.db.QueryRow(ctx, sqlStr, args...)

	var nc model.NodeConfig
	err := row.Scan(
		&nc.DefaultBlock,
		&nc.InputBoxDeploymentBlock,
		&nc.InputBoxAddress,
		&nc.ChainID,
		&nc.CreatedAt,
		&nc.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &nc, nil
}
