// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cartesi/rollups-espresso-reader/internal/repository/postgres/db/rollupsdb/espresso/table"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-jet/jet/v2/postgres"
)

func (r *postgresRepository) GetEspressoConfig(
	ctx context.Context,
	nameOrAddress string,
) (uint64, uint64, error) {
	app := common.HexToAddress(nameOrAddress)
	sel := table.AppInfo.
		SELECT(
			table.AppInfo.StartingBlock,
			table.AppInfo.Namespace,
		).
		FROM(table.AppInfo).
		WHERE(
			table.AppInfo.ApplicationAddress.EQ(postgres.Bytea(app.Bytes())),
		)

	sqlStr, args := sel.Sql()
	row := r.db.QueryRow(ctx, sqlStr, args...)

	var startingBlock, namespace uint64
	err := row.Scan(
		&startingBlock,
		&namespace,
	)
	if err != nil {
		return 0, 0, err
	}
	return startingBlock, namespace, nil
}

func (r *postgresRepository) UpdateEspressoConfig(
	ctx context.Context,
	nameOrAddress string,
	startingBlock uint64,
	namespace uint64,
) error {
	app := common.HexToAddress(nameOrAddress)
	insertStmt := table.AppInfo.INSERT(
		table.AppInfo.ApplicationAddress,
		table.AppInfo.StartingBlock,
		table.AppInfo.Namespace,
	).VALUES(
		postgres.Bytea(app.Bytes()),
		startingBlock,
		namespace,
	)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	sqlStr, args := insertStmt.Sql()
	_, err = tx.Exec(ctx, sqlStr, args...)

	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}
	return nil
}

func (r *postgresRepository) GetEspressoNonce(
	ctx context.Context,
	senderAddress string,
	nameOrAddress string,
) (uint64, error) {
	// assume all are hex address string
	sender := common.HexToAddress(senderAddress)
	app := common.HexToAddress(nameOrAddress)
	sel := table.EspressoNonce.
		SELECT(table.EspressoNonce.Nonce).
		FROM(table.EspressoNonce).
		WHERE(
			table.EspressoNonce.SenderAddress.EQ(postgres.Bytea(sender.Bytes())).
				AND(table.EspressoNonce.ApplicationAddress.EQ(postgres.Bytea(app.Bytes()))),
		)

	sqlStr, args := sel.Sql()
	row := r.db.QueryRow(ctx, sqlStr, args...)

	var nonce uint64
	err := row.Scan(
		&nonce,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

func (r *postgresRepository) UpdateEspressoNonce(
	ctx context.Context,
	senderAddress string,
	nameOrAddress string,
) error {
	// assume all are hex address string
	sender := common.HexToAddress(senderAddress)
	app := common.HexToAddress(nameOrAddress)
	nonce, err := r.GetEspressoNonce(ctx, senderAddress, nameOrAddress)
	if err != nil {
		return err
	}
	nextNonce := nonce + 1

	nonceInsertStmt := table.EspressoNonce.INSERT(
		table.EspressoNonce.SenderAddress,
		table.EspressoNonce.ApplicationAddress,
		table.EspressoNonce.Nonce,
	).VALUES(
		postgres.Bytea(sender.Bytes()),
		postgres.Bytea(app.Bytes()),
		nextNonce,
	)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	sqlStr, args := nonceInsertStmt.
		ON_CONFLICT(table.EspressoNonce.SenderAddress, table.EspressoNonce.ApplicationAddress).
		DO_UPDATE(postgres.SET(
			table.EspressoNonce.Nonce.SET(postgres.RawInt(fmt.Sprintf("%d", nextNonce))),
		)).Sql()
	_, err = tx.Exec(ctx, sqlStr, args...)

	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return errors.Join(err, tx.Rollback(ctx))
	}
	return nil
}

func (r *postgresRepository) GetInputIndex(
	ctx context.Context,
	nameOrAddress string,
) (uint64, error) {
	app := common.HexToAddress(nameOrAddress)
	sel := table.AppInfo.
		SELECT(table.AppInfo.Index).
		FROM(table.AppInfo).
		WHERE(
			table.AppInfo.ApplicationAddress.EQ(postgres.Bytea(app.Bytes())),
		)

	sqlStr, args := sel.Sql()
	row := r.db.QueryRow(ctx, sqlStr, args...)

	var index *uint64
	err := row.Scan(
		&index,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if index == nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return *index, nil
}

func (r *postgresRepository) UpdateInputIndex(
	ctx context.Context,
	nameOrAddress string,
) error {
	app := common.HexToAddress(nameOrAddress)
	index, err := r.GetInputIndex(ctx, nameOrAddress)
	if err != nil {
		return err
	}
	nextIndex := index + 1

	updateStmt := table.AppInfo.
		UPDATE(table.AppInfo.Index).
		SET(table.AppInfo.Index.SET(postgres.RawInt(fmt.Sprintf("%d", nextIndex)))).
		WHERE(table.AppInfo.ApplicationAddress.EQ(postgres.Bytea(app.Bytes())))

	sqlStr, args := updateStmt.Sql()
	_, err = r.db.Exec(ctx, sqlStr, args...)
	return err
}

func (r *postgresRepository) GetLastProcessedEspressoBlock(
	ctx context.Context,
	nameOrAddress string,
) (uint64, error) {
	app := common.HexToAddress(nameOrAddress)
	sel := table.AppInfo.
		SELECT(table.AppInfo.LastProcessedEspressoBlock).
		FROM(table.AppInfo).
		WHERE(
			table.AppInfo.ApplicationAddress.EQ(postgres.Bytea(app.Bytes())),
		)

	sqlStr, args := sel.Sql()
	row := r.db.QueryRow(ctx, sqlStr, args...)
	var lastProcessedEspressoBlock *uint64
	err := row.Scan(
		&lastProcessedEspressoBlock,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if lastProcessedEspressoBlock == nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return *lastProcessedEspressoBlock, nil
}

func (r *postgresRepository) UpdateLastProcessedEspressoBlock(
	ctx context.Context,
	nameOrAddress string,
	lastProcessedEspressoBlock uint64,
) error {
	app := common.HexToAddress(nameOrAddress)

	updateStmt := table.AppInfo.
		UPDATE(table.AppInfo.LastProcessedEspressoBlock).
		SET(table.AppInfo.LastProcessedEspressoBlock.SET(postgres.RawFloat(fmt.Sprintf("%d", lastProcessedEspressoBlock)))).
		WHERE(table.AppInfo.ApplicationAddress.EQ(postgres.Bytea(app.Bytes())))

	sqlStr, args := updateStmt.Sql()
	_, err := r.db.Exec(ctx, sqlStr, args...)
	return err
}
