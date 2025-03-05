// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	. "github.com/cartesi/rollups-espresso-reader/internal/model"
	"github.com/ethereum/go-ethereum/common"
)

type Pagination struct {
	Limit  int64
	Offset int64
}

type ApplicationFilter struct {
	State   *ApplicationState
	Name    *string
	Address *string
}

type EpochFilter struct {
	Status      *EpochStatus
	BeforeBlock *uint64
}

type InputFilter struct {
	InputIndex           *uint64
	Status               *InputCompletionStatus
	NotStatus            *InputCompletionStatus
	TransactionReference *[]byte
}

type Range struct {
	Start uint64
	End   uint64
}

type OutputFilter struct {
	InputIndex *uint64
	BlockRange *Range
}

type ReportFilter struct {
	InputIndex *uint64
}

type ApplicationRepository interface {
	CreateApplication(ctx context.Context, app *Application) (int64, error)
	GetApplication(ctx context.Context, nameOrAddress string) (*Application, error)
	UpdateApplication(ctx context.Context, app *Application) error
	UpdateApplicationState(ctx context.Context, app *Application) error
	DeleteApplication(ctx context.Context, id int64) error
	ListApplications(ctx context.Context, f ApplicationFilter, p Pagination) ([]*Application, error)
}

type EpochRepository interface {
	CreateEpoch(ctx context.Context, nameOrAddress string, e *Epoch) error
	CreateEpochsAndInputs(ctx context.Context, nameOrAddress string, epochInputMap map[*Epoch][]*Input, blockNumber uint64) error

	GetEpoch(ctx context.Context, nameOrAddress string, index uint64) (*Epoch, error)
	GetEpochByVirtualIndex(ctx context.Context, nameOrAddress string, index uint64) (*Epoch, error)

	UpdateEpoch(ctx context.Context, nameOrAddress string, e *Epoch) error
	UpdateEpochsInputsProcessed(ctx context.Context, nameOrAddress string) error

	ListEpochs(ctx context.Context, nameOrAddress string, f EpochFilter, p Pagination) ([]*Epoch, error)
}

type InputRepository interface {
	GetInput(ctx context.Context, nameOrAddress string, inputIndex uint64) (*Input, error)
	GetInputByTxReference(ctx context.Context, nameOrAddress string, ref *common.Hash) (*Input, error)
	GetLastInput(ctx context.Context, appAddress string, epochIndex uint64) (*Input, error) // FIXME remove me (list, filter and order)
	ListInputs(ctx context.Context, nameOrAddress string, f InputFilter, p Pagination) ([]*Input, error)
}

type OutputRepository interface {
	GetOutput(ctx context.Context, nameOrAddress string, outputIndex uint64) (*Output, error)
	ListOutputs(ctx context.Context, nameOrAddress string, f OutputFilter, p Pagination) ([]*Output, error)
}

type ReportRepository interface {
	GetReport(ctx context.Context, nameOrAddress string, reportIndex uint64) (*Report, error)
	ListReports(ctx context.Context, nameOrAddress string, f ReportFilter, p Pagination) ([]*Report, error)
}

type NodeConfigRepository interface {
	SaveNodeConfigRaw(ctx context.Context, key string, rawJSON []byte) error
	LoadNodeConfigRaw(ctx context.Context, key string) (rawJSON []byte, createdAt, updatedAt time.Time, err error)
}

type EspressoRepository interface {
	GetEspressoNonce(
		ctx context.Context,
		senderAddress string,
		nameOrAddress string,
	) (uint64, error)
	UpdateEspressoNonce(
		ctx context.Context,
		senderAddress string,
		nameOrAddress string,
	) error
	GetInputIndex(
		ctx context.Context,
		nameOrAddress string,
	) (uint64, error)
	UpdateInputIndex(
		ctx context.Context,
		nameOrAddress string,
	) error
	GetLastProcessedEspressoBlock(
		ctx context.Context,
		nameOrAddress string,
	) (uint64, error)
	UpdateLastProcessedEspressoBlock(
		ctx context.Context,
		nameOrAddress string,
		lastProcessedEspressoBlock uint64,
	) error
}

type Repository interface {
	ApplicationRepository
	EpochRepository
	InputRepository
	OutputRepository
	ReportRepository
	NodeConfigRepository
	EspressoRepository
}

func SaveNodeConfig[T any](
	ctx context.Context,
	repo NodeConfigRepository,
	nc *NodeConfig[T],
) error {
	data, err := json.Marshal(nc.Value)
	if err != nil {
		return fmt.Errorf("marshal node_config value failed: %w", err)
	}
	err = repo.SaveNodeConfigRaw(ctx, nc.Key, data)
	if err != nil {
		return err
	}
	return nil
}

func LoadNodeConfig[T any](
	ctx context.Context,
	repo NodeConfigRepository,
	key string,
) (*NodeConfig[T], error) {
	raw, createdAt, updatedAt, err := repo.LoadNodeConfigRaw(ctx, key)
	if err != nil {
		return nil, err
	}
	var val T
	if err := json.Unmarshal(raw, &val); err != nil {
		return nil, fmt.Errorf("unmarshal node_config value failed: %w", err)
	}
	return &NodeConfig[T]{
		Key:       key,
		Value:     val,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
