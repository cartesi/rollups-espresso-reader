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

func (r *postgresRepository) GetReport(
	ctx context.Context,
	nameOrAddress string,
	reportIndex uint64,
) (*model.Report, error) {

	whereClause, err := getWhereClauseFromNameOrAddress(nameOrAddress)
	if err != nil {
		return nil, err
	}

	sel := table.Report.
		SELECT(
			table.Report.InputEpochApplicationID,
			table.Report.InputIndex,
			table.Report.Index,
			table.Report.RawData,
			table.Report.CreatedAt,
			table.Report.UpdatedAt,
		).
		FROM(
			table.Report.
				INNER_JOIN(table.Application,
					table.Report.InputEpochApplicationID.EQ(table.Application.ID),
				),
		).
		WHERE(
			whereClause.
				AND(table.Report.Index.EQ(postgres.RawFloat(fmt.Sprintf("%d", reportIndex)))),
		)

	sqlStr, args := sel.Sql()
	row := r.db.QueryRow(ctx, sqlStr, args...)

	var rp model.Report
	err = row.Scan(
		&rp.InputEpochApplicationID,
		&rp.InputIndex,
		&rp.Index,
		&rp.RawData,
		&rp.CreatedAt,
		&rp.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rp, nil
}
