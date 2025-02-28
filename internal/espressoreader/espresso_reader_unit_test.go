// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package espressoreader

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/EspressoSystems/espresso-sequencer-go/client"
	tagged_base64 "github.com/EspressoSystems/espresso-sequencer-go/tagged-base64"
	"github.com/EspressoSystems/espresso-sequencer-go/types"
	espresso "github.com/EspressoSystems/espresso-sequencer-go/types/common"
	"github.com/cartesi/rollups-espresso-reader/internal/evmreader"
	"github.com/cartesi/rollups-espresso-reader/internal/model"
	"github.com/cartesi/rollups-espresso-reader/internal/repository"
	"github.com/ethereum/go-ethereum/common"
	eth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type EpochInputMap = map[*model.Epoch][]*model.Input
type MockRepository struct {
	mock.Mock
	inputs []model.Input
}

// CreateApplication implements repository.Repository.
func (m *MockRepository) CreateApplication(ctx context.Context, app *model.Application) (int64, error) {
	panic("unimplemented")
}

// CreateEpoch implements repository.Repository.
func (m *MockRepository) CreateEpoch(ctx context.Context, nameOrAddress string, e *model.Epoch) error {
	panic("unimplemented")
}

// CreateEpochsAndInputs implements repository.Repository.
func (m *MockRepository) CreateEpochsAndInputs(ctx context.Context, nameOrAddress string, epochInputMap map[*model.Epoch][]*model.Input, blockNumber uint64) error {
	for e := range epochInputMap {
		for _, input := range epochInputMap[e] {
			m.inputs = append(m.inputs, *input)
		}
	}
	args := m.Called(ctx, nameOrAddress, epochInputMap, blockNumber)
	return args.Error(0)
}

// DeleteApplication implements repository.Repository.
func (m *MockRepository) DeleteApplication(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// GetApplication implements repository.Repository.
func (m *MockRepository) GetApplication(ctx context.Context, nameOrAddress string) (*model.Application, error) {
	panic("unimplemented")
}

// GetEpoch implements repository.Repository.
func (m *MockRepository) GetEpoch(ctx context.Context, nameOrAddress string, index uint64) (*model.Epoch, error) {
	args := m.Called(ctx, nameOrAddress, index)
	epoch := args.Get(0).(model.Epoch)
	return &epoch, args.Error(1)
}

// GetEpochByVirtualIndex implements repository.Repository.
func (m *MockRepository) GetEpochByVirtualIndex(ctx context.Context, nameOrAddress string, index uint64) (*model.Epoch, error) {
	panic("unimplemented")
}

// GetEspressoNonce implements repository.Repository.
func (m *MockRepository) GetEspressoNonce(ctx context.Context, senderAddress string, nameOrAddress string) (uint64, error) {
	args := m.Called(ctx, senderAddress, nameOrAddress)
	return uint64(args.Int(0)), args.Error(1)
}

// GetExecutionParameters implements repository.Repository.
func (m *MockRepository) GetExecutionParameters(ctx context.Context, applicationID int64) (*model.ExecutionParameters, error) {
	panic("unimplemented")
}

// GetInput implements repository.Repository.
func (m *MockRepository) GetInput(ctx context.Context, nameOrAddress string, inputIndex uint64) (*model.Input, error) {
	panic("unimplemented")
}

// GetInputByTxReference implements repository.Repository.
func (m *MockRepository) GetInputByTxReference(ctx context.Context, nameOrAddress string, ref *common.Hash) (*model.Input, error) {
	panic("unimplemented")
}

// GetInputIndex implements repository.Repository.
func (m *MockRepository) GetInputIndex(ctx context.Context, nameOrAddress string) (uint64, error) {
	args := m.Called(ctx, nameOrAddress)
	return uint64(args.Int(0)), args.Error(1)
}

// GetLastInput implements repository.Repository.
func (m *MockRepository) GetLastInput(ctx context.Context, appAddress string, epochIndex uint64) (*model.Input, error) {
	panic("unimplemented")
}

// GetLastProcessedEspressoBlock implements repository.Repository.
func (m *MockRepository) GetLastProcessedEspressoBlock(ctx context.Context, nameOrAddress string) (uint64, error) {
	panic("unimplemented")
}

// GetOutput implements repository.Repository.
func (m *MockRepository) GetOutput(ctx context.Context, nameOrAddress string, outputIndex uint64) (*model.Output, error) {
	panic("unimplemented")
}

// GetReport implements repository.Repository.
func (m *MockRepository) GetReport(ctx context.Context, nameOrAddress string, reportIndex uint64) (*model.Report, error) {
	panic("unimplemented")
}

// ListApplications implements repository.Repository.
func (m *MockRepository) ListApplications(ctx context.Context, f repository.ApplicationFilter, p repository.Pagination) ([]*model.Application, error) {
	panic("unimplemented")
}

// ListEpochs implements repository.Repository.
func (m *MockRepository) ListEpochs(ctx context.Context, nameOrAddress string, f repository.EpochFilter, p repository.Pagination) ([]*model.Epoch, error) {
	panic("unimplemented")
}

// ListInputs implements repository.Repository.
func (m *MockRepository) ListInputs(ctx context.Context, nameOrAddress string, f repository.InputFilter, p repository.Pagination) ([]*model.Input, error) {
	panic("unimplemented")
}

// ListOutputs implements repository.Repository.
func (m *MockRepository) ListOutputs(ctx context.Context, nameOrAddress string, f repository.OutputFilter, p repository.Pagination) ([]*model.Output, error) {
	panic("unimplemented")
}

// ListReports implements repository.Repository.
func (m *MockRepository) ListReports(ctx context.Context, nameOrAddress string, f repository.ReportFilter, p repository.Pagination) ([]*model.Report, error) {
	panic("unimplemented")
}

// LoadNodeConfigRaw implements repository.Repository.
func (m *MockRepository) LoadNodeConfigRaw(ctx context.Context, key string) (rawJSON []byte, createdAt time.Time, updatedAt time.Time, err error) {
	panic("unimplemented")
}

// SaveNodeConfigRaw implements repository.Repository.
func (m *MockRepository) SaveNodeConfigRaw(ctx context.Context, key string, rawJSON []byte) error {
	panic("unimplemented")
}

// SelectClaimPairsPerApp implements repository.Repository.
func (m *MockRepository) SelectClaimPairsPerApp(ctx context.Context) (map[common.Address]*model.ClaimRow, map[common.Address]*model.ClaimRow, error) {
	panic("unimplemented")
}

// SelectNewestAcceptedClaimPerApp implements repository.Repository.
func (m *MockRepository) SelectNewestAcceptedClaimPerApp(ctx context.Context) (map[common.Address]*model.ClaimRow, error) {
	panic("unimplemented")
}

// SelectOldestComputedClaimPerApp implements repository.Repository.
func (m *MockRepository) SelectOldestComputedClaimPerApp(ctx context.Context) (map[common.Address]*model.ClaimRow, error) {
	panic("unimplemented")
}

// StoreAdvanceResult implements repository.Repository.
func (m *MockRepository) StoreAdvanceResult(ctx context.Context, appId int64, ar *model.AdvanceResult) error {
	panic("unimplemented")
}

// StoreClaimAndProofs implements repository.Repository.
func (m *MockRepository) StoreClaimAndProofs(ctx context.Context, epoch *model.Epoch, outputs []*model.Output) error {
	panic("unimplemented")
}

// UpdateApplication implements repository.Repository.
func (m *MockRepository) UpdateApplication(ctx context.Context, app *model.Application) error {
	panic("unimplemented")
}

// UpdateApplicationState implements repository.Repository.
func (m *MockRepository) UpdateApplicationState(ctx context.Context, app *model.Application) error {
	panic("unimplemented")
}

// UpdateEpoch implements repository.Repository.
func (m *MockRepository) UpdateEpoch(ctx context.Context, nameOrAddress string, e *model.Epoch) error {
	panic("unimplemented")
}

// UpdateEpochWithSubmittedClaim implements repository.Repository.
func (m *MockRepository) UpdateEpochWithSubmittedClaim(ctx context.Context, application_id int64, index uint64, transaction_hash common.Hash) error {
	panic("unimplemented")
}

// UpdateEpochsClaimAccepted implements repository.Repository.
func (m *MockRepository) UpdateEpochsClaimAccepted(ctx context.Context, nameOrAddress string, epochs []*model.Epoch, mostRecentBlockNumber uint64) error {
	panic("unimplemented")
}

// UpdateEpochsInputsProcessed implements repository.Repository.
func (m *MockRepository) UpdateEpochsInputsProcessed(ctx context.Context, nameOrAddress string) error {
	panic("unimplemented")
}

// UpdateEspressoNonce implements repository.Repository.
func (m *MockRepository) UpdateEspressoNonce(ctx context.Context, senderAddress string, nameOrAddress string) error {
	args := m.Called(ctx, senderAddress, nameOrAddress)
	return args.Error(0)
}

// UpdateExecutionParameters implements repository.Repository.
func (m *MockRepository) UpdateExecutionParameters(ctx context.Context, ep *model.ExecutionParameters) error {
	panic("unimplemented")
}

// UpdateInputIndex implements repository.Repository.
func (m *MockRepository) UpdateInputIndex(ctx context.Context, nameOrAddress string) error {
	args := m.Called(ctx, nameOrAddress)
	return args.Error(0)
}

// UpdateLastProcessedEspressoBlock implements repository.Repository.
func (m *MockRepository) UpdateLastProcessedEspressoBlock(ctx context.Context, nameOrAddress string, lastProcessedEspressoBlock uint64) error {
	panic("unimplemented")
}

// UpdateOutputsExecution implements repository.Repository.
func (m *MockRepository) UpdateOutputsExecution(ctx context.Context, nameOrAddress string, executedOutputs []*model.Output, blockNumber uint64) error {
	panic("unimplemented")
}

type MockEspressoClient struct {
	mock.Mock
}

// FetchBlockMerkleProof implements EspressoClient.
func (m *MockEspressoClient) FetchBlockMerkleProof(ctx context.Context, rootHeight uint64, hotshotHeight uint64) (espresso.HotShotBlockMerkleProof, error) {
	panic("unimplemented")
}

// FetchHeaderByHeight implements EspressoClient.
func (m *MockEspressoClient) FetchHeaderByHeight(ctx context.Context, blockHeight uint64) (types.HeaderImpl, error) {
	panic("unimplemented")
}

// FetchHeadersByRange implements EspressoClient.
func (m *MockEspressoClient) FetchHeadersByRange(ctx context.Context, from uint64, until uint64) ([]types.HeaderImpl, error) {
	panic("unimplemented")
}

// FetchLatestBlockHeight implements EspressoClient.
func (m *MockEspressoClient) FetchLatestBlockHeight(ctx context.Context) (uint64, error) {
	return uint64(1), nil
}

// FetchTransactionByHash implements EspressoClient.
func (m *MockEspressoClient) FetchTransactionByHash(ctx context.Context, hash *tagged_base64.TaggedBase64) (espresso.TransactionQueryData, error) {
	panic("unimplemented")
}

// FetchVidCommonByHeight implements EspressoClient.
func (m *MockEspressoClient) FetchVidCommonByHeight(ctx context.Context, blockHeight uint64) (json.RawMessage, error) {
	panic("unimplemented")
}

// SubmitTransaction implements EspressoClient.
func (m *MockEspressoClient) SubmitTransaction(ctx context.Context, tx espresso.Transaction) (*tagged_base64.TaggedBase64, error) {
	panic("unimplemented")
}

func (m *MockEspressoClient) FetchTransactionsInBlock(ctx context.Context, blockHeight uint64, namespace uint64) (client.TransactionsInBlock, error) {
	args := m.Called(ctx, blockHeight, namespace)
	transactions, _ := args.Get(0).(client.TransactionsInBlock)
	return transactions, args.Error(1)
}

type MockEthClient struct {
	mock.Mock
}

func (m *MockEthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*eth_types.Header, error) {
	args := m.Called(ctx, number)
	header, _ := args.Get(0).(eth_types.Header)
	return &header, args.Error(1)
}

type EspressoReaderUnitTestSuite struct {
	suite.Suite
	espressoReader     EspressoReader
	mockEspressoClient *MockEspressoClient
	mockDatabase       *MockRepository
	mockEthClient      *MockEthClient
}

func (s *EspressoReaderUnitTestSuite) SetupSuite() {
	espressoApiURL := ""
	startingBlock := 1
	namespace := 55555
	chainId := 31337
	inputBoxDeploymentBlock := 10
	maxRetries := 10
	maxDelay := 1
	mockDatabase := new(MockRepository)
	mockEthClient := new(MockEthClient)
	evmReader := evmreader.NewEvmReader(
		mockEthClient, nil, nil, nil, uint64(inputBoxDeploymentBlock), "0", nil, false,
	)
	s.espressoReader = NewEspressoReader(
		espressoApiURL,
		uint64(startingBlock),
		uint64(namespace),
		mockDatabase,
		&evmReader,
		uint64(chainId),
		uint64(inputBoxDeploymentBlock),
		uint64(maxRetries),
		uint64(maxDelay),
	)
	mockEspressoClient := new(MockEspressoClient)
	s.mockEthClient = mockEthClient
	s.mockEspressoClient = mockEspressoClient
	s.mockDatabase = mockDatabase
	s.espressoReader.client = mockEspressoClient
}

func (s *EspressoReaderUnitTestSuite) TestReadEspresso() {
	transaction := `{"typedData":{"domain":{"name":"Cartesi","version":"0.1.0","chainId":11155111,"verifyingContract":"0x0000000000000000000000000000000000000000"},"types":{"EIP712Domain":[{"name":"name","type":"string"},{"name":"version","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"address"}],"CartesiMessage":[{"name":"app","type":"address"},{"name":"nonce","type":"uint64"},{"name":"max_gas_price","type":"uint128"},{"name":"data","type":"bytes"}]},"primaryType":"CartesiMessage","message":{"app":"0x5a205fcb6947e200615b75c409ac0aa486d77649","nonce":0,"data":"0xdeadbeef","max_gas_price":"10"}},"account":"0x04dd244dd1a881f29e59a14139691540c55799f070d1cb3ed133a5d0b6ed0004e237d974f007aad1d2fb333169f0fee69bd93e6f36b4531d41a0e23a963a65f28f","signature":"0x73f1c93f1a9bea6675d2271432878f2705119cec63a8d5d8c1275155c845a17a6699d21cb568910429a8eb53b11cc3cedc196f323e06bb72f83f09bf2ed113371c"}`
	transactions := []types.Bytes{
		[]byte(transaction),
	}
	transactionsInBlock := client.TransactionsInBlock{
		Transactions: transactions,
	}

	ctx := context.Background()

	currentBlockHeight := 10
	l1FinalizedLatestHeight := 3
	l1FinalizedTimestamp := 17
	currentInputIndex := 0

	s.mockEspressoClient.On(
		"FetchTransactionsInBlock",
		mock.Anything, // context.Context
		mock.Anything, // blockHeight
		mock.Anything, // namespace
	).Return(transactionsInBlock, nil)

	s.mockDatabase.On("GetEspressoNonce",
		mock.Anything, // context.Context
		mock.Anything, // msgSender
		mock.Anything, // appAddressStr
	).Return(0, nil)

	s.mockDatabase.On("GetInputIndex",
		mock.Anything, // context.Context
		mock.Anything, // nameOrAddress
	).Return(currentInputIndex, nil)

	s.mockDatabase.On("GetEpoch",
		mock.Anything, // context.Context
		mock.Anything, // nameOrAddress
		mock.Anything, // index
	).Return(model.Epoch{}, nil)

	s.mockDatabase.On("CreateEpochsAndInputs",
		mock.Anything, // context.Context
		mock.Anything, // nameOrAddress
		mock.Anything, // epochInputMap
		mock.Anything, // blockNumber
	).Return(nil)

	s.mockDatabase.On("UpdateEspressoNonce",
		mock.Anything, // context.Context
		mock.Anything, // senderAddress
		mock.Anything, // nameOrAddress
	).Return(nil)

	s.mockDatabase.On("UpdateInputIndex",
		mock.Anything, // context.Context
		mock.Anything, // nameOrAddress
	).Return(nil)

	s.mockEthClient.On("HeaderByNumber",
		mock.Anything, // context.Context
		mock.Anything, // number
	).Return(eth_types.Header{
		MixDigest: common.Hash{},
	}, nil)

	application := model.Application{
		ID:                  33331,
		IApplicationAddress: common.HexToAddress("0x5a205fcb6947e200615b75c409ac0aa486d77649"),
		EpochLength:         10, // cannot be zero
	}
	appEvmType := evmreader.TypeExportApplication{
		Application: application,
	}

	s.espressoReader.readEspresso(ctx, appEvmType, uint64(currentBlockHeight), uint64(l1FinalizedLatestHeight), uint64(l1FinalizedTimestamp))

	s.Equal(1, len(s.mockDatabase.inputs))
	input := s.mockDatabase.inputs[0]
	s.Equal(currentInputIndex, int(input.Index))
	s.Equal(l1FinalizedLatestHeight, int(input.BlockNumber))
	s.Equal("415bf3630000000000000000000000000000000000000000000000000000000000007a690000000000000000000000005a205fcb6947e200615b75c409ac0aa486d77649000000000000000000000000b5c1674c0527b6c31a5019fd04a6c1529396da37000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000110000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000004deadbeef00000000000000000000000000000000000000000000000000000000", common.Bytes2Hex(input.RawData))
	s.Equal(model.InputCompletionStatus("NONE"), input.Status)
	s.Equal("0xc355ddfeeedb4498f328c5bce91a1160af9a5d052b2b23ee2e339111cecf1aba", input.TransactionReference.Hex())

	s.mockEspressoClient.AssertExpectations(s.T())
	s.mockDatabase.AssertExpectations(s.T())
	s.mockEthClient.AssertExpectations(s.T())
}

func TestEspressoReaderUnitTestSuite(t *testing.T) {
	suite.Run(t, new(EspressoReaderUnitTestSuite))
}
