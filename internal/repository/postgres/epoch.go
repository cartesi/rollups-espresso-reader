// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"sort"

	"github.com/cartesi/rollups-espresso-reader/internal/model"
	"github.com/cartesi/rollups-espresso-reader/internal/repository/postgres/db/rollupsdb/public/table"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5"
)

func (r *PostgresRepository) CreateEpoch(
	ctx context.Context,
	nameOrAddress string,
	e *model.Epoch,
) error {

	whereClause, err := getWhereClauseFromNameOrAddress(nameOrAddress)
	if err != nil {
		return err
	}

	selectQuery := postgres.SELECT(
		table.Application.ID,
		postgres.RawFloat(fmt.Sprintf("%d", e.Index)),
		postgres.RawFloat(fmt.Sprintf("%d", e.FirstBlock)),
		postgres.RawFloat(fmt.Sprintf("%d", e.LastBlock)),
		postgres.Bytea(e.ClaimHash),
		postgres.Bytea(e.ClaimTransactionHash),
		postgres.NewEnumValue(e.Status.String()),
		postgres.RawFloat(fmt.Sprintf("%d", e.VirtualIndex)),
	).FROM(
		table.Application,
	).WHERE(
		whereClause,
	)

	insertStmt := table.Epoch.INSERT(
		table.Epoch.ApplicationID,
		table.Epoch.Index,
		table.Epoch.FirstBlock,
		table.Epoch.LastBlock,
		table.Epoch.ClaimHash,
		table.Epoch.ClaimTransactionHash,
		table.Epoch.Status,
		table.Epoch.VirtualIndex,
	).QUERY(
		selectQuery,
	)

	sqlStr, args := insertStmt.Sql()
	_, err = r.db.Exec(ctx, sqlStr, args...)
	return err
}

func getEpochNextVirtualIndex(
	ctx context.Context,
	tx pgx.Tx,
	nameOrAddress string,
) (uint64, error) {

	whereClause, err := getWhereClauseFromNameOrAddress(nameOrAddress)
	if err != nil {
		return 0, err
	}

	query := table.Epoch.SELECT(
		postgres.COALESCE(
			postgres.Float(1).ADD(postgres.MAXf(table.Epoch.VirtualIndex)),
			postgres.Float(0),
		),
	).FROM(
		table.Epoch.INNER_JOIN(table.Application, table.Epoch.ApplicationID.EQ(table.Application.ID)),
	).WHERE(
		whereClause,
	)

	queryStr, args := query.Sql()
	var currentIndex uint64
	err = tx.QueryRow(ctx, queryStr, args...).Scan(&currentIndex)
	if err != nil {
		err = fmt.Errorf("failed to get the next epoch virtual index: %w", err)
		return 0, errors.Join(err, tx.Rollback(ctx))
	}
	return currentIndex, nil
}

func orderEpochs(epochInputsMap map[*model.Epoch][]*model.Input) []*model.Epoch {
	epochs := make([]*model.Epoch, 0, len(epochInputsMap))
	for e := range epochInputsMap {
		epochs = append(epochs, e)
	}

	sort.Slice(epochs, func(i, j int) bool {
		return epochs[i].FirstBlock < epochs[j].FirstBlock
	})

	return epochs
}

func (r *PostgresRepository) CreateEpochsAndInputs(
	ctx context.Context,
	tx pgx.Tx,
	nameOrAddress string,
	epochInputsMap map[*model.Epoch][]*model.Input,
	blockNumber uint64,
) error {

	whereClause, err := getWhereClauseFromNameOrAddress(nameOrAddress)
	if err != nil {
		return err
	}

	epochInsertStmt := table.Epoch.INSERT(
		table.Epoch.ApplicationID,
		table.Epoch.Index,
		table.Epoch.FirstBlock,
		table.Epoch.LastBlock,
		table.Epoch.Status,
		table.Epoch.VirtualIndex,
	)

	inputInsertStmt := table.Input.
		INSERT(
			table.Input.EpochApplicationID,
			table.Input.EpochIndex,
			table.Input.Index,
			table.Input.BlockNumber,
			table.Input.RawData,
			table.Input.Status,
			table.Input.TransactionReference,
		)

	epochs := orderEpochs(epochInputsMap)
	for _, epoch := range epochs {
		inputs := epochInputsMap[epoch]

		nextVirtualIndex, err := getEpochNextVirtualIndex(ctx, tx, nameOrAddress)
		if err != nil {
			return err
		}

		epochSelectQuery := table.Application.SELECT(
			table.Application.ID,
			postgres.RawFloat(fmt.Sprintf("%d", epoch.Index)),
			postgres.RawFloat(fmt.Sprintf("%d", epoch.FirstBlock)),
			postgres.RawFloat(fmt.Sprintf("%d", epoch.LastBlock)),
			postgres.NewEnumValue(epoch.Status.String()),
			postgres.RawFloat(fmt.Sprintf("%d", nextVirtualIndex)),
		).WHERE(
			whereClause,
		)

		sqlStr, args := epochInsertStmt.QUERY(epochSelectQuery).
			ON_CONFLICT(table.Epoch.ApplicationID, table.Epoch.Index).
			DO_UPDATE(postgres.SET(
				table.Epoch.Status.SET(postgres.NewEnumValue(epoch.Status.String())),
			)).Sql() // FIXME on conflict
		_, err = tx.Exec(ctx, sqlStr, args...)

		if err != nil {
			return err
		}

		for _, input := range inputs {
			slog.Debug("inserting input to db", "epoch", epoch.Index, "input block#", input.BlockNumber, "input index", input.Index)
			inputSelectQuery := table.Application.SELECT(
				table.Application.ID,
				postgres.RawFloat(fmt.Sprintf("%d", epoch.Index)),
				postgres.RawFloat(fmt.Sprintf("%d", input.Index)),
				postgres.RawFloat(fmt.Sprintf("%d", input.BlockNumber)),
				postgres.Bytea(input.RawData),
				postgres.NewEnumValue(input.Status.String()),
				postgres.Bytea(input.TransactionReference.Bytes()),
			).WHERE(
				whereClause,
			)

			sqlStr, args := inputInsertStmt.QUERY(inputSelectQuery).Sql()
			_, err := tx.Exec(ctx, sqlStr, args...)
			if err != nil {
				return err
			}
		}
	}

	// Update last processed L1 block
	appUpdateStmt := table.Application.
		UPDATE(
			table.Application.LastInputCheckBlock,
		).
		SET(
			postgres.RawFloat(fmt.Sprintf("%d", blockNumber)),
		).
		WHERE(whereClause)

	sqlStr, args := appUpdateStmt.Sql()
	_, err = tx.Exec(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) GetEpoch(
	ctx context.Context,
	nameOrAddress string,
	index uint64,
) (*model.Epoch, error) {

	whereClause, err := getWhereClauseFromNameOrAddress(nameOrAddress)
	if err != nil {
		return nil, err
	}

	stmt := table.Epoch.
		SELECT(
			table.Epoch.ApplicationID,
			table.Epoch.Index,
			table.Epoch.FirstBlock,
			table.Epoch.LastBlock,
			table.Epoch.ClaimHash,
			table.Epoch.ClaimTransactionHash,
			table.Epoch.Status,
			table.Epoch.VirtualIndex,
			table.Epoch.CreatedAt,
			table.Epoch.UpdatedAt,
		).
		FROM(
			table.Epoch.
				INNER_JOIN(table.Application,
					table.Epoch.ApplicationID.EQ(table.Application.ID),
				),
		).
		WHERE(
			whereClause.
				AND(table.Epoch.Index.EQ(postgres.RawFloat(fmt.Sprintf("%d", index)))),
		)

	sqlStr, args := stmt.Sql()
	row := r.db.QueryRow(ctx, sqlStr, args...)

	var ep model.Epoch
	err = row.Scan(
		&ep.ApplicationID,
		&ep.Index,
		&ep.FirstBlock,
		&ep.LastBlock,
		&ep.ClaimHash,
		&ep.ClaimTransactionHash,
		&ep.Status,
		&ep.VirtualIndex,
		&ep.CreatedAt,
		&ep.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ep, nil
}

func (r *PostgresRepository) GetEpochByVirtualIndex(
	ctx context.Context,
	nameOrAddress string,
	index uint64,
) (*model.Epoch, error) {

	whereClause, err := getWhereClauseFromNameOrAddress(nameOrAddress)
	if err != nil {
		return nil, err
	}

	stmt := table.Epoch.
		SELECT(
			table.Epoch.ApplicationID,
			table.Epoch.Index,
			table.Epoch.FirstBlock,
			table.Epoch.LastBlock,
			table.Epoch.ClaimHash,
			table.Epoch.ClaimTransactionHash,
			table.Epoch.Status,
			table.Epoch.VirtualIndex,
			table.Epoch.CreatedAt,
			table.Epoch.UpdatedAt,
		).
		FROM(
			table.Epoch.
				INNER_JOIN(table.Application,
					table.Epoch.ApplicationID.EQ(table.Application.ID),
				),
		).
		WHERE(
			whereClause.
				AND(table.Epoch.VirtualIndex.EQ(postgres.RawFloat(fmt.Sprintf("%d", index)))),
		)

	sqlStr, args := stmt.Sql()
	row := r.db.QueryRow(ctx, sqlStr, args...)

	var ep model.Epoch
	err = row.Scan(
		&ep.ApplicationID,
		&ep.Index,
		&ep.FirstBlock,
		&ep.LastBlock,
		&ep.ClaimHash,
		&ep.ClaimTransactionHash,
		&ep.Status,
		&ep.VirtualIndex,
		&ep.CreatedAt,
		&ep.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ep, nil
}

func (r *PostgresRepository) UpdateEpoch(
	ctx context.Context,
	nameOrAddress string,
	e *model.Epoch,
) error {

	whereClause, err := getWhereClauseFromNameOrAddress(nameOrAddress)
	if err != nil {
		return err
	}

	updStmt := table.Epoch.
		UPDATE(
			table.Epoch.ClaimHash,
			table.Epoch.ClaimTransactionHash,
			table.Epoch.Status,
		).
		SET(
			e.ClaimHash,
			e.ClaimTransactionHash,
			e.Status,
		).
		FROM(
			table.Application,
		).
		WHERE(
			whereClause.
				AND(table.Epoch.ApplicationID.EQ(table.Application.ID)).
				AND(table.Epoch.Index.EQ(postgres.RawFloat(fmt.Sprintf("%d", e.Index)))),
		)

	sqlStr, args := updStmt.Sql()
	cmd, err := r.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}
