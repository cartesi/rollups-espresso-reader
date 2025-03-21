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
	"github.com/cartesi/rollups-espresso-reader/internal/repository"
	"github.com/cartesi/rollups-espresso-reader/internal/repository/postgres/db/rollupsdb/public/table"
)

// ------------------------ ApplicationRepository Methods ------------------------ //

func (r *postgresRepository) CreateApplication(
	ctx context.Context,
	app *model.Application,
) (int64, error) {

	insertStmt := table.Application.
		INSERT(
			table.Application.Name,
			table.Application.IapplicationAddress,
			table.Application.IconsensusAddress,
			table.Application.TemplateHash,
			table.Application.TemplateURI,
			table.Application.EpochLength,
			table.Application.State,
			table.Application.LastProcessedBlock,
			table.Application.LastClaimCheckBlock,
			table.Application.LastOutputCheckBlock,
			table.Application.ProcessedInputs,
		).
		VALUES(
			app.Name,
			app.IApplicationAddress,
			app.IConsensusAddress,
			app.TemplateHash,
			app.TemplateURI,
			app.EpochLength,
			app.State,
			app.LastProcessedBlock,
			app.LastClaimCheckBlock,
			app.LastOutputCheckBlock,
			app.ProcessedInputs,
		).
		RETURNING(table.Application.ID)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}

	sqlStr, args := insertStmt.Sql()
	var newID int64
	err = tx.QueryRow(ctx, sqlStr, args...).Scan(&newID)
	if err != nil {
		return 0, errors.Join(err, tx.Rollback(ctx))
	}

	sqlStr, args = table.ExecutionParameters.
		INSERT(
			table.ExecutionParameters.ApplicationID,
		).
		VALUES(
			newID,
		).Sql()

	_, err = tx.Exec(ctx, sqlStr, args...)
	if err != nil {
		return 0, errors.Join(err, tx.Rollback(ctx))
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, errors.Join(err, tx.Rollback(ctx))
	}
	return newID, nil
}

// GetApplication retrieves one application by ID, optionally loading status & execution parameters.
func (r *postgresRepository) GetApplication(
	ctx context.Context,
	nameOrAddress string,
) (*model.Application, error) {

	whereClause, err := getWhereClauseFromNameOrAddress(nameOrAddress)
	if err != nil {
		return nil, err
	}

	stmt := table.Application.
		SELECT(
			table.Application.ID,
			table.Application.Name,
			table.Application.IapplicationAddress,
			table.Application.IconsensusAddress,
			table.Application.TemplateHash,
			table.Application.TemplateURI,
			table.Application.EpochLength,
			table.Application.State,
			table.Application.Reason,
			table.Application.LastProcessedBlock,
			table.Application.LastClaimCheckBlock,
			table.Application.LastOutputCheckBlock,
			table.Application.ProcessedInputs,
			table.Application.CreatedAt,
			table.Application.UpdatedAt,
			table.ExecutionParameters.ApplicationID,
			table.ExecutionParameters.SnapshotRetention,
			table.ExecutionParameters.AdvanceIncCycles,
			table.ExecutionParameters.AdvanceMaxCycles,
			table.ExecutionParameters.InspectIncCycles,
			table.ExecutionParameters.InspectMaxCycles,
			table.ExecutionParameters.AdvanceIncDeadline,
			table.ExecutionParameters.AdvanceMaxDeadline,
			table.ExecutionParameters.InspectIncDeadline,
			table.ExecutionParameters.InspectMaxDeadline,
			table.ExecutionParameters.LoadDeadline,
			table.ExecutionParameters.StoreDeadline,
			table.ExecutionParameters.FastDeadline,
			table.ExecutionParameters.MaxConcurrentInspects,
			table.ExecutionParameters.CreatedAt,
			table.ExecutionParameters.UpdatedAt,
		).
		FROM(
			table.Application.INNER_JOIN(
				table.ExecutionParameters,
				table.ExecutionParameters.ApplicationID.EQ(table.Application.ID),
			),
		).
		WHERE(whereClause)

	sqlStr, args := stmt.Sql()
	row := r.db.QueryRow(ctx, sqlStr, args...)

	var app model.Application
	err = row.Scan(
		&app.ID,
		&app.Name,
		&app.IApplicationAddress,
		&app.IConsensusAddress,
		&app.TemplateHash,
		&app.TemplateURI,
		&app.EpochLength,
		&app.State,
		&app.Reason,
		&app.LastProcessedBlock,
		&app.LastClaimCheckBlock,
		&app.LastOutputCheckBlock,
		&app.ProcessedInputs,
		&app.CreatedAt,
		&app.UpdatedAt,
		&app.ExecutionParameters.ApplicationID,
		&app.ExecutionParameters.SnapshotRetention,
		&app.ExecutionParameters.AdvanceIncCycles,
		&app.ExecutionParameters.AdvanceMaxCycles,
		&app.ExecutionParameters.InspectIncCycles,
		&app.ExecutionParameters.InspectMaxCycles,
		&app.ExecutionParameters.AdvanceIncDeadline,
		&app.ExecutionParameters.AdvanceMaxDeadline,
		&app.ExecutionParameters.InspectIncDeadline,
		&app.ExecutionParameters.InspectMaxDeadline,
		&app.ExecutionParameters.LoadDeadline,
		&app.ExecutionParameters.StoreDeadline,
		&app.ExecutionParameters.FastDeadline,
		&app.ExecutionParameters.MaxConcurrentInspects,
		&app.ExecutionParameters.CreatedAt,
		&app.ExecutionParameters.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // not found
	}
	if err != nil {
		return nil, err
	}

	return &app, nil
}

// UpdateApplication updates an existing application row.
func (r *postgresRepository) UpdateApplication(
	ctx context.Context,
	app *model.Application,
) error {

	updateStmt := table.Application.
		UPDATE(
			table.Application.Name,
			table.Application.IapplicationAddress,
			table.Application.IconsensusAddress,
			table.Application.TemplateHash,
			table.Application.TemplateURI,
			table.Application.EpochLength,
			table.Application.State,
			table.Application.Reason,
			table.Application.LastProcessedBlock,
			table.Application.LastClaimCheckBlock,
			table.Application.LastOutputCheckBlock,
			table.Application.ProcessedInputs,
		).
		SET(
			app.Name,
			app.IApplicationAddress,
			app.IConsensusAddress,
			app.TemplateHash,
			app.TemplateURI,
			app.EpochLength,
			app.State,
			app.Reason,
			app.LastProcessedBlock,
			app.LastClaimCheckBlock,
			app.LastOutputCheckBlock,
			app.ProcessedInputs,
		).
		WHERE(table.Application.ID.EQ(postgres.Int(app.ID)))

	sqlStr, args := updateStmt.Sql()
	_, err := r.db.Exec(ctx, sqlStr, args...)
	return err
}

func (r *postgresRepository) UpdateApplicationState(
	ctx context.Context,
	appID int64,
	state model.ApplicationState,
	reason *string,
) error {

	updateStmt := table.Application.
		UPDATE(
			table.Application.State,
			table.Application.Reason,
		).
		SET(
			state,
			reason,
		).
		WHERE(table.Application.ID.EQ(postgres.Int(appID)))

	sqlStr, args := updateStmt.Sql()
	_, err := r.db.Exec(ctx, sqlStr, args...)
	return err
}

// DeleteApplication removes the row from "application" by ID.
func (r *postgresRepository) DeleteApplication(
	ctx context.Context,
	id int64,
) error {

	delStmt := table.Application.
		DELETE().
		WHERE(table.Application.ID.EQ(postgres.Int(id)))

	sqlStr, args := delStmt.Sql()
	cmd, err := r.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() != 1 {
		return fmt.Errorf("application with ID %d not found", id)
	}
	return nil
}

// ListApplications queries multiple apps with optional filters & pagination.
func (r *postgresRepository) ListApplications(
	ctx context.Context,
	f repository.ApplicationFilter,
	p repository.Pagination,
) ([]*model.Application, uint64, error) {

	sel := table.Application.
		SELECT(
			table.Application.ID,
			table.Application.Name,
			table.Application.IapplicationAddress,
			table.Application.IconsensusAddress,
			table.Application.TemplateHash,
			table.Application.TemplateURI,
			table.Application.EpochLength,
			table.Application.State,
			table.Application.Reason,
			table.Application.LastProcessedBlock,
			table.Application.LastClaimCheckBlock,
			table.Application.LastOutputCheckBlock,
			table.Application.ProcessedInputs,
			table.Application.CreatedAt,
			table.Application.UpdatedAt,
			table.ExecutionParameters.ApplicationID,
			table.ExecutionParameters.SnapshotRetention,
			table.ExecutionParameters.AdvanceIncCycles,
			table.ExecutionParameters.AdvanceMaxCycles,
			table.ExecutionParameters.InspectIncCycles,
			table.ExecutionParameters.InspectMaxCycles,
			table.ExecutionParameters.AdvanceIncDeadline,
			table.ExecutionParameters.AdvanceMaxDeadline,
			table.ExecutionParameters.InspectIncDeadline,
			table.ExecutionParameters.InspectMaxDeadline,
			table.ExecutionParameters.LoadDeadline,
			table.ExecutionParameters.StoreDeadline,
			table.ExecutionParameters.FastDeadline,
			table.ExecutionParameters.MaxConcurrentInspects,
			table.ExecutionParameters.CreatedAt,
			table.ExecutionParameters.UpdatedAt,
			postgres.COUNT(postgres.STAR).OVER().AS("total_count"),
		).
		FROM(
			table.Application.INNER_JOIN(
				table.ExecutionParameters,
				table.ExecutionParameters.ApplicationID.EQ(table.Application.ID),
			),
		)

	conditions := []postgres.BoolExpression{}
	if f.State != nil {
		conditions = append(conditions, table.Application.State.EQ(postgres.NewEnumValue(f.State.String())))
	}

	if len(conditions) > 0 {
		sel = sel.WHERE(postgres.AND(conditions...))
	}

	sel.ORDER_BY(table.Application.Name.ASC())

	// Apply pagination
	if p.Limit > 0 {
		sel = sel.LIMIT(int64(p.Limit))
	}
	if p.Offset > 0 {
		sel = sel.OFFSET(int64(p.Offset))
	}

	sqlStr, args := sel.Sql()
	rows, err := r.db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var apps []*model.Application
	var total uint64
	for rows.Next() {
		var app model.Application
		err := rows.Scan(
			&app.ID,
			&app.Name,
			&app.IApplicationAddress,
			&app.IConsensusAddress,
			&app.TemplateHash,
			&app.TemplateURI,
			&app.EpochLength,
			&app.State,
			&app.Reason,
			&app.LastProcessedBlock,
			&app.LastClaimCheckBlock,
			&app.LastOutputCheckBlock,
			&app.ProcessedInputs,
			&app.CreatedAt,
			&app.UpdatedAt,
			&app.ExecutionParameters.ApplicationID,
			&app.ExecutionParameters.SnapshotRetention,
			&app.ExecutionParameters.AdvanceIncCycles,
			&app.ExecutionParameters.AdvanceMaxCycles,
			&app.ExecutionParameters.InspectIncCycles,
			&app.ExecutionParameters.InspectMaxCycles,
			&app.ExecutionParameters.AdvanceIncDeadline,
			&app.ExecutionParameters.AdvanceMaxDeadline,
			&app.ExecutionParameters.InspectIncDeadline,
			&app.ExecutionParameters.InspectMaxDeadline,
			&app.ExecutionParameters.LoadDeadline,
			&app.ExecutionParameters.StoreDeadline,
			&app.ExecutionParameters.FastDeadline,
			&app.ExecutionParameters.MaxConcurrentInspects,
			&app.ExecutionParameters.CreatedAt,
			&app.ExecutionParameters.UpdatedAt,
			&total,
		)
		if err != nil {
			return nil, 0, err
		}
		apps = append(apps, &app)
	}

	return apps, total, nil
}
