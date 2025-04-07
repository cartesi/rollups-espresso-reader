// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package espressoreader

import (
	"context"
	"encoding/hex"
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
	"github.com/cartesi/rollups-espresso-reader/pkg/contracts/iinputbox"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	eth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type EpochInputMap = map[*model.Epoch][]*model.Input
type MockRepository struct {
	mock.Mock
	inputs     []model.Input
	epochs     []*model.Epoch
	nonce      int
	inputIndex int
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
		found := false
		for loopIndex, loopEpoch := range m.epochs {
			if loopEpoch.LastBlock == e.LastBlock && e.LastBlock != 0 {
				found = true
				m.epochs[loopIndex].Status = e.Status
			}
		}
		if found == false {
			m.epochs = append(m.epochs, e)
		}

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
	return uint64(m.nonce), args.Error(1)
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
	return uint64(m.inputIndex), args.Error(1)
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
func (m *MockRepository) ListApplications(ctx context.Context, f repository.ApplicationFilter, p repository.Pagination) ([]*model.Application, uint64, error) {
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

// StoreClaimAndProofs implements repository.Repository.
func (m *MockRepository) StoreClaimAndProofs(ctx context.Context, epoch *model.Epoch, outputs []*model.Output) error {
	panic("unimplemented")
}

// UpdateApplication implements repository.Repository.
func (m *MockRepository) UpdateApplication(ctx context.Context, app *model.Application) error {
	panic("unimplemented")
}

// UpdateApplicationState implements repository.Repository.
func (m *MockRepository) UpdateApplicationState(ctx context.Context, appID int64, state model.ApplicationState, reason *string) error {
	panic("unimplemented")
}

// UpdateEventLastCheckBlock implements repository.Repository.
func (m *MockRepository) UpdateEventLastCheckBlock(ctx context.Context, appIDs []int64, event model.MonitoredEvent, blockNumber uint64) error {
	return nil
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

// UpdateEspressoNonce implements repository.Repository.
func (m *MockRepository) UpdateEspressoNonce(ctx context.Context, senderAddress string, nameOrAddress string) error {
	args := m.Called(ctx, senderAddress, nameOrAddress)
	m.nonce++
	return args.Error(0)
}

// UpdateExecutionParameters implements repository.Repository.
func (m *MockRepository) UpdateExecutionParameters(ctx context.Context, ep *model.ExecutionParameters) error {
	panic("unimplemented")
}

// UpdateInputIndex implements repository.Repository.
func (m *MockRepository) UpdateInputIndex(ctx context.Context, nameOrAddress string) error {
	args := m.Called(ctx, nameOrAddress)
	m.inputIndex++
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

func (m *MockRepository) Close() {
	panic("unimplemented")
}

func (m *MockRepository) GetEspressoConfig(
	ctx context.Context,
	nameOrAddress string,
) (uint64, uint64, error) {
	panic("unimplemented")
}

func (m *MockRepository) UpdateEspressoConfig(
	ctx context.Context,
	nameOrAddress string,
	startingBlock uint64, namespace uint64,
) error {
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

type MockInputSource struct {
	mock.Mock
}

func (m *MockInputSource) RetrieveInputs(opts *bind.FilterOpts, appAddresses []common.Address, index []*big.Int) ([]iinputbox.IInputBoxInputAdded, error) {
	args := m.Called(opts, appAddresses, index)
	events, _ := args.Get(0).([]iinputbox.IInputBoxInputAdded)
	return events, args.Error(1)
}

type MockEspressoHelper struct {
	mock.Mock
}

func (m *MockEspressoHelper) getL1FinalizedHeight(ctx context.Context, espressoBlockHeight uint64, delay uint64, url string) (uint64, uint64) {
	args := m.Called(ctx, espressoBlockHeight, delay, url)
	return uint64(args.Int(0)), uint64(args.Int(1))
}

func (m *MockEspressoHelper) readEspressoHeadersByRange(ctx context.Context, from uint64, until uint64, delay uint64, url string) string {
	panic("unimplemented")
}

func (m *MockEspressoHelper) getNSTableByRange(ctx context.Context, from uint64, until uint64, delay uint64, url string) (string, error) {
	panic("unimplemented")
}

func (m *MockEspressoHelper) extractNS(nsTable []byte) []uint32 {
	panic("unimplemented")
}

type EspressoReaderUnitTestSuite struct {
	suite.Suite
	espressoReader     EspressoReader
	mockEspressoClient *MockEspressoClient
	mockDatabase       *MockRepository
	mockEthClient      *MockEthClient
	mockInputSource    *MockInputSource
	mockEspressoHelper *MockEspressoHelper
}

var transactions = []types.Bytes{
	[]byte(`{"typedData":{"domain":{"name":"Cartesi","version":"0.1.0","chainId":11155111,"verifyingContract":"0x0000000000000000000000000000000000000000"},"types":{"EIP712Domain":[{"name":"name","type":"string"},{"name":"version","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"address"}],"CartesiMessage":[{"name":"app","type":"address"},{"name":"nonce","type":"uint64"},{"name":"max_gas_price","type":"uint128"},{"name":"data","type":"bytes"}]},"primaryType":"CartesiMessage","message":{"app":"0x5a205fcb6947e200615b75c409ac0aa486d77649","nonce":0,"data":"0xdeadbeef","max_gas_price":"10"}},"account":"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266","signature":"0xebbaaf80346db193a06c82eb2aab587c59bc4626c84a51d552c2b4457371892c7361223b2a09d935f32af010081a27d42e60fa4e11a58e8dc056bafe6f848c9f1c"}`),
	[]byte(`{"typedData":{"domain":{"name":"Cartesi","version":"0.1.0","chainId":11155111,"verifyingContract":"0x0000000000000000000000000000000000000000"},"types":{"EIP712Domain":[{"name":"name","type":"string"},{"name":"version","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"address"}],"CartesiMessage":[{"name":"app","type":"address"},{"name":"nonce","type":"uint64"},{"name":"max_gas_price","type":"uint128"},{"name":"data","type":"bytes"}]},"primaryType":"CartesiMessage","message":{"app":"0x5a205fcb6947e200615b75c409ac0aa486d77649","nonce":1,"data":"0xdeadbeef","max_gas_price":"10"}},"account":"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266","signature":"0x8edcdc950107fd3a64192d2803da1624cadcac1e5a5a2b80301af25662e4107f44cb7841bd74ce6b571452494193115364b7f8f40adf4d5a2f6ec07a4073324d1b"}`),
}

func (s *EspressoReaderUnitTestSuite) SetupTest() {
	espressoApiURL := ""
	chainId := 13370
	maxRetries := 10
	maxDelay := 1
	mockDatabase := new(MockRepository)
	mockEthClient := new(MockEthClient)
	mockInputSource := new(MockInputSource)
	s.mockInputSource = mockInputSource
	evmReader := evmreader.NewEvmReader(
		mockEthClient, nil, mockDatabase, "0", nil, true,
	)
	blockchainHttpEndpoint := "http://localhost:8545"
	s.espressoReader = NewEspressoReader(
		espressoApiURL,
		mockDatabase,
		&evmReader,
		uint64(chainId),
		uint64(maxRetries),
		uint64(maxDelay),
		blockchainHttpEndpoint,
	)
	mockEspressoClient := new(MockEspressoClient)
	s.mockEthClient = mockEthClient
	s.mockEspressoClient = mockEspressoClient
	s.mockDatabase = mockDatabase
	s.espressoReader.client = mockEspressoClient
	mockEspressoHelper := new(MockEspressoHelper)
	s.mockEspressoHelper = mockEspressoHelper
	s.espressoReader.espressoHelper = mockEspressoHelper
}

func ignoreError(val interface{}, _ error) interface{} {
	return val
}

func (s *EspressoReaderUnitTestSuite) TestReadEspressoInput() {
	transaction := transactions[0]
	espressoInput, err := s.espressoReader.readEspressoInput([]byte(transaction))
	s.Require().NoError(err)
	s.Equal("0x5a205fcb6947e200615b75c409ac0aa486d77649", espressoInput.app)
	s.Equal(ignoreError(hex.DecodeString("deadbeef")), espressoInput.data)
	s.Equal(uint64(0), espressoInput.nonce)
	s.Equal("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", espressoInput.msgSender.Hex())
	s.Equal("0x0ba026185fb624135e746f7ecae1734085f2a2483b5a0b18b34332e47a88a02e", espressoInput.sigHash)
}

func (s *EspressoReaderUnitTestSuite) TestReadEspresso() {
	transaction := transactions[0]
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
		InputSource: s.mockInputSource,
	}

	s.espressoReader.readEspresso(ctx, appEvmType, uint64(currentBlockHeight), 55555, uint64(l1FinalizedLatestHeight), uint64(l1FinalizedTimestamp))

	s.Equal(1, len(s.mockDatabase.inputs))
	input := s.mockDatabase.inputs[0]
	s.Equal(currentInputIndex, int(input.Index))
	s.Equal(l1FinalizedLatestHeight, int(input.BlockNumber))
	s.Equal("415bf3630000000000000000000000000000000000000000000000000000000000007a690000000000000000000000005a205fcb6947e200615b75c409ac0aa486d77649000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000110000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000004deadbeef00000000000000000000000000000000000000000000000000000000", common.Bytes2Hex(input.RawData))
	s.Equal(model.InputCompletionStatus("NONE"), input.Status)
	s.Equal("0x0ba026185fb624135e746f7ecae1734085f2a2483b5a0b18b34332e47a88a02e", input.TransactionReference.Hex())

	s.mockEspressoClient.AssertExpectations(s.T())
	s.mockDatabase.AssertExpectations(s.T())
	s.mockEthClient.AssertExpectations(s.T())
}

func (s *EspressoReaderUnitTestSuite) TestReadEspressoWith2Transactions() {
	transactions := []types.Bytes{
		transactions[0],
		transactions[1],
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
		InputSource: s.mockInputSource,
	}

	s.espressoReader.readEspresso(ctx, appEvmType, uint64(currentBlockHeight), 55555, uint64(l1FinalizedLatestHeight), uint64(l1FinalizedTimestamp))

	s.Equal(2, len(s.mockDatabase.inputs))
	input := s.mockDatabase.inputs[0]
	s.Equal(0, int(input.Index))
	s.Equal(l1FinalizedLatestHeight, int(input.BlockNumber))
	s.Equal("415bf3630000000000000000000000000000000000000000000000000000000000007a690000000000000000000000005a205fcb6947e200615b75c409ac0aa486d77649000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000110000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000004deadbeef00000000000000000000000000000000000000000000000000000000", common.Bytes2Hex(input.RawData))
	s.Equal(model.InputCompletionStatus("NONE"), input.Status)
	s.Equal("0x0ba026185fb624135e746f7ecae1734085f2a2483b5a0b18b34332e47a88a02e", input.TransactionReference.Hex())

	input2 := s.mockDatabase.inputs[1]
	s.Equal(1, int(input2.Index))
	s.Equal(l1FinalizedLatestHeight, int(input2.BlockNumber))
	s.Equal("415bf3630000000000000000000000000000000000000000000000000000000000007a690000000000000000000000005a205fcb6947e200615b75c409ac0aa486d77649000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000110000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000004deadbeef00000000000000000000000000000000000000000000000000000000", common.Bytes2Hex(input.RawData))
	s.Equal(model.InputCompletionStatus("NONE"), input2.Status)
	s.Equal("0xfa89de5b3417b7abb60eaf2ca208459a46611eae5d81a4bf120e721dd1ad5e5c", input2.TransactionReference.Hex())

	s.mockEspressoClient.AssertExpectations(s.T())
	s.mockDatabase.AssertExpectations(s.T())
	s.mockEthClient.AssertExpectations(s.T())
}

func (s *EspressoReaderUnitTestSuite) TestEdgeCaseInputAtTheEndOfEpoch() {
	transaction := transactions[0]
	transactions := []types.Bytes{
		[]byte(transaction),
	}
	transactionsInBlock := client.TransactionsInBlock{
		Transactions: transactions,
	}

	ctx := context.Background()

	currentBlockHeight := 10
	lastProcessedBlock := 3
	l1FinalizedLatestHeight := 5
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
		uint64(0),     // index
	).Return(model.Epoch{
		Index:      0,
		FirstBlock: 0,
		LastBlock:  uint64(lastProcessedBlock),
		Status:     model.EpochStatus_Open,
	}, nil)

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

	s.mockInputSource.On("RetrieveInputs",
		mock.Anything, //  *bind.FilterOpts
		mock.Anything, // appAddresses
		mock.Anything, // index
	).Return(nil, nil)

	s.mockEspressoHelper.On("getL1FinalizedHeight",
		mock.Anything, // context.Context,
		mock.Anything, // espressoBlockHeight
		mock.Anything, // delay
		mock.Anything, // url
	).Return(l1FinalizedLatestHeight, l1FinalizedLatestHeight)

	application := model.Application{
		ID:                  33331,
		IApplicationAddress: common.HexToAddress("0x5a205fcb6947e200615b75c409ac0aa486d77649"),
		EpochLength:         uint64(lastProcessedBlock + 1), // test the edge case
	}
	appEvmType := evmreader.TypeExportApplication{
		Application: application,
		InputSource: s.mockInputSource,
	}

	s.espressoReader.readEspresso(ctx, appEvmType, uint64(currentBlockHeight), 55555, uint64(lastProcessedBlock), uint64(l1FinalizedTimestamp))
	var apps []evmreader.TypeExportApplication
	apps = append(apps, appEvmType)

	s.espressoReader.readL1(ctx, appEvmType, uint64(currentBlockHeight), uint64(lastProcessedBlock))

	s.Equal(model.EpochStatus_Closed, s.mockDatabase.epochs[0].Status)

	s.mockEspressoClient.AssertExpectations(s.T())
	s.mockDatabase.AssertExpectations(s.T())
	s.mockEthClient.AssertExpectations(s.T())
	s.mockEspressoHelper.AssertExpectations(s.T())
}

func (s *EspressoReaderUnitTestSuite) TestEdgeCaseSkippingL1Blocks() {
	transaction := transactions[0]
	transactions := []types.Bytes{
		[]byte(transaction),
	}
	transactionsInBlock := client.TransactionsInBlock{
		Transactions: transactions,
	}

	ctx := context.Background()

	currentBlockHeight := 10
	lastProcessedBlock := 3
	l1FinalizedLatestHeight := 5
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
		uint64(0),     // index
	).Return(model.Epoch{
		Index:      0,
		FirstBlock: 0,
		LastBlock:  uint64(lastProcessedBlock),
		Status:     model.EpochStatus_Closed,
	}, nil)

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

	appContractAddress := common.HexToAddress("0x5a205fcb6947e200615b75c409ac0aa486d77649")
	indexValue := new(big.Int).SetUint64(1)
	inputData := []byte{0x01, 0x02, 0x03, 0x4a, 0x5b, 0x6c}
	rawLog := eth_types.Log{
		Address:     common.HexToAddress("0x5a205fcb6947e200615b75c409ac0aa486d77649"),
		Topics:      []common.Hash{common.HexToHash("0xabcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789")},
		Data:        []byte{0xaa, 0xbb, 0xcc},
		BlockNumber: 5,
		TxHash:      common.HexToHash("0x9876543210abcdef9876543210abcdef9876543210abcdef9876543210abcdef"),
		TxIndex:     0,
		BlockHash:   common.HexToHash("0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"),
		Index:       1,
		Removed:     false,
	}
	mockEvent := iinputbox.IInputBoxInputAdded{
		AppContract: appContractAddress,
		Index:       indexValue,
		Input:       inputData,
		Raw:         rawLog,
	}
	var mockEvents []iinputbox.IInputBoxInputAdded
	mockEvents = append(mockEvents, mockEvent)
	s.mockInputSource.On("RetrieveInputs",
		mock.Anything, //  *bind.FilterOpts
		mock.Anything, // appAddresses
		mock.Anything, // index
	).Return(mockEvents, nil)

	s.mockEspressoHelper.On("getL1FinalizedHeight",
		mock.Anything, // context.Context,
		mock.Anything, // espressoBlockHeight
		mock.Anything, // delay
		mock.Anything, // url
	).Return(l1FinalizedLatestHeight, l1FinalizedLatestHeight)

	application := model.Application{
		ID:                  33331,
		IApplicationAddress: common.HexToAddress("0x5a205fcb6947e200615b75c409ac0aa486d77649"),
		EpochLength:         uint64(lastProcessedBlock + 1), // test the edge case
	}
	appEvmType := evmreader.TypeExportApplication{
		Application: application,
		InputSource: s.mockInputSource,
	}

	s.espressoReader.readEspresso(ctx, appEvmType, uint64(currentBlockHeight), 55555, uint64(lastProcessedBlock), uint64(l1FinalizedTimestamp))
	var apps []evmreader.TypeExportApplication
	apps = append(apps, appEvmType)

	s.espressoReader.readL1(ctx, appEvmType, uint64(currentBlockHeight), uint64(lastProcessedBlock))

	s.Equal(2, len(s.mockDatabase.inputs))

	s.mockEspressoClient.AssertExpectations(s.T())
	s.mockDatabase.AssertExpectations(s.T())
	s.mockEthClient.AssertExpectations(s.T())
	s.mockEspressoHelper.AssertExpectations(s.T())
}

func TestEspressoReaderUnitTestSuite(t *testing.T) {
	suite.Run(t, new(EspressoReaderUnitTestSuite))
}
