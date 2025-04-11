// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"

	"github.com/cartesi/rollups-espresso-reader/internal/model"
	"github.com/cartesi/rollups-espresso-reader/internal/repository/postgres/db/rollupsdb/public/table"
)

func (r *PostgresRepository) GetOutput(
	ctx context.Context,
	nameOrAddress string,
	outputIndex uint64,
) (*model.Output, error) {

	whereClause, err := getWhereClauseFromNameOrAddress(nameOrAddress)
	if err != nil {
		return nil, err
	}

	sel := table.Output.
		SELECT(
			table.Output.InputEpochApplicationID,
			table.Output.InputIndex,
			table.Output.Index,
			table.Output.RawData,
			table.Output.Hash,
			table.Output.OutputHashesSiblings,
			table.Output.ExecutionTransactionHash,
			table.Output.CreatedAt,
			table.Output.UpdatedAt,
		).
		FROM(
			table.Output.
				INNER_JOIN(table.Application,
					table.Output.InputEpochApplicationID.EQ(table.Application.ID),
				),
		).
		WHERE(
			whereClause.
				AND(table.Output.Index.EQ(postgres.RawFloat(fmt.Sprintf("%d", outputIndex)))),
		)

	sqlStr, args := sel.Sql()
	row := r.db.QueryRow(ctx, sqlStr, args...)

	var o model.Output
	err = row.Scan(
		&o.InputEpochApplicationID,
		&o.InputIndex,
		&o.Index,
		&o.RawData,
		&o.Hash,
		&o.OutputHashesSiblings,
		&o.ExecutionTransactionHash,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}
