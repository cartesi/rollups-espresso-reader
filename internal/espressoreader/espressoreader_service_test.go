package espressoreader

import (
	"context"
	"testing"
	"time"

	"github.com/cartesi/rollups-espresso-reader/internal/model"
	"github.com/cartesi/rollups-espresso-reader/internal/repository"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5"
)

type fakeRepo struct{}

// Close implements repository.Repository.
func (f *fakeRepo) Close() {
	panic("unimplemented")
}

// CreateApplication implements repository.Repository.
func (f *fakeRepo) CreateApplication(ctx context.Context, app *model.Application) (int64, error) {
	panic("unimplemented")
}

// CreateEpoch implements repository.Repository.
func (f *fakeRepo) CreateEpoch(ctx context.Context, nameOrAddress string, e *model.Epoch) error {
	panic("unimplemented")
}

// CreateEpochsAndInputs implements repository.Repository.
func (f *fakeRepo) CreateEpochsAndInputs(ctx context.Context, nameOrAddress string, epochInputMap map[*model.Epoch][]*model.Input, blockNumber uint64, espressoUpdateInfo *model.EspressoUpdateInfo) error {
	panic("unimplemented")
}

// DeleteApplication implements repository.Repository.
func (f *fakeRepo) DeleteApplication(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// GetApplication implements repository.Repository.
func (f *fakeRepo) GetApplication(ctx context.Context, nameOrAddress string) (*model.Application, error) {
	panic("unimplemented")
}

// GetEpoch implements repository.Repository.
func (f *fakeRepo) GetEpoch(ctx context.Context, nameOrAddress string, index uint64) (*model.Epoch, error) {
	panic("unimplemented")
}

// GetEpochByVirtualIndex implements repository.Repository.
func (f *fakeRepo) GetEpochByVirtualIndex(ctx context.Context, nameOrAddress string, index uint64) (*model.Epoch, error) {
	panic("unimplemented")
}

// GetEspressoConfig implements repository.Repository.
func (f *fakeRepo) GetEspressoConfig(ctx context.Context, nameOrAddress string) (uint64, uint64, error) {
	panic("unimplemented")
}

// GetEspressoNonce implements repository.Repository.
func (f *fakeRepo) GetEspressoNonce(ctx context.Context, senderAddress string, nameOrAddress string) (uint64, error) {
	panic("unimplemented")
}

// GetEspressoNonceWithTx implements repository.Repository.
func (f *fakeRepo) GetEspressoNonceWithTx(ctx context.Context, tx pgx.Tx, senderAddress string, nameOrAddress string) (uint64, error) {
	panic("unimplemented")
}

// GetInput implements repository.Repository.
func (f *fakeRepo) GetInput(ctx context.Context, nameOrAddress string, inputIndex uint64) (*model.Input, error) {
	panic("unimplemented")
}

// GetInputByTxReference implements repository.Repository.
func (f *fakeRepo) GetInputByTxReference(ctx context.Context, nameOrAddress string, ref *common.Hash) (*model.Input, error) {
	panic("unimplemented")
}

// GetInputIndex implements repository.Repository.
func (f *fakeRepo) GetInputIndex(ctx context.Context, nameOrAddress string) (uint64, error) {
	panic("unimplemented")
}

// GetInputIndexWithTx implements repository.Repository.
func (f *fakeRepo) GetInputIndexWithTx(ctx context.Context, tx pgx.Tx, nameOrAddress string) (uint64, error) {
	panic("unimplemented")
}

// GetLastInput implements repository.Repository.
func (f *fakeRepo) GetLastInput(ctx context.Context, appAddress string, epochIndex uint64) (*model.Input, error) {
	panic("unimplemented")
}

// GetLastProcessedEspressoBlock implements repository.Repository.
func (f *fakeRepo) GetLastProcessedEspressoBlock(ctx context.Context, nameOrAddress string) (uint64, error) {
	panic("unimplemented")
}

// GetOutput implements repository.Repository.
func (f *fakeRepo) GetOutput(ctx context.Context, nameOrAddress string, outputIndex uint64) (*model.Output, error) {
	panic("unimplemented")
}

// GetReport implements repository.Repository.
func (f *fakeRepo) GetReport(ctx context.Context, nameOrAddress string, reportIndex uint64) (*model.Report, error) {
	panic("unimplemented")
}

// InsertEspressoConfig implements repository.Repository.
func (f *fakeRepo) InsertEspressoConfig(ctx context.Context, nameOrAddress string, startingBlock uint64, namespace uint64) error {
	panic("unimplemented")
}

// ListApplications implements repository.Repository.
func (*fakeRepo) ListApplications(ctx context.Context, f repository.ApplicationFilter, p repository.Pagination) ([]*model.Application, uint64, error) {
	panic("unimplemented")
}

// LoadNodeConfigRaw implements repository.Repository.
func (f *fakeRepo) LoadNodeConfigRaw(ctx context.Context, key string) (rawJSON []byte, createdAt time.Time, updatedAt time.Time, err error) {
	panic("unimplemented")
}

// SaveNodeConfigRaw implements repository.Repository.
func (f *fakeRepo) SaveNodeConfigRaw(ctx context.Context, key string, rawJSON []byte) error {
	panic("unimplemented")
}

// UpdateApplication implements repository.Repository.
func (f *fakeRepo) UpdateApplication(ctx context.Context, app *model.Application) error {
	panic("unimplemented")
}

// UpdateApplicationState implements repository.Repository.
func (f *fakeRepo) UpdateApplicationState(ctx context.Context, appID int64, state model.ApplicationState, reason *string) error {
	panic("unimplemented")
}

// UpdateEpoch implements repository.Repository.
func (f *fakeRepo) UpdateEpoch(ctx context.Context, nameOrAddress string, e *model.Epoch) error {
	panic("unimplemented")
}

// UpdateEspressoBlock implements repository.Repository.
func (f *fakeRepo) UpdateEspressoBlock(ctx context.Context, appAddress common.Address, lastProcessedEspressoBlock uint64) error {
	panic("unimplemented")
}

// UpdateEventLastCheckBlock implements repository.Repository.
func (f *fakeRepo) UpdateEventLastCheckBlock(ctx context.Context, appIDs []int64, event model.MonitoredEvent, blockNumber uint64) error {
	panic("unimplemented")
}

func TestTrySetupEvmReader_InvalidEndpointReturnsError(t *testing.T) {
	ctx := context.Background()

	service := &EspressoReaderService{
		blockchainHttpEndpoint: "invalid-endpoint",
		blockchainWsEndpoint:   "invalid-endpoint",
	}

	r := &fakeRepo{}
	_, err := service.trySetupEvmReader(ctx, r)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
