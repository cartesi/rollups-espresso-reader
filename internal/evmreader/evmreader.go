// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package evmreader

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/cartesi/rollups-espresso-reader/internal/model"
	. "github.com/cartesi/rollups-espresso-reader/internal/model"
	"github.com/cartesi/rollups-espresso-reader/internal/repository"
	appcontract "github.com/cartesi/rollups-espresso-reader/pkg/contracts/iapplication"
	"github.com/cartesi/rollups-espresso-reader/pkg/contracts/iinputbox"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

//go:embed abi.json
var abiData string

const EvmReaderConfigKey = "evm-reader"

type PersistentConfig struct {
	DefaultBlock       DefaultBlock
	InputReaderEnabled bool
	ChainID            uint64
}

// Interface for Input reading
type InputSourceAdapter interface {
	// Wrapper for FilterInputAdded(), which is automatically generated
	// by go-ethereum and cannot be used for testing
	RetrieveInputs(opts *bind.FilterOpts, appAddresses []common.Address, index []*big.Int,
	) ([]iinputbox.IInputBoxInputAdded, error)
}

// Interface for the node repository
type EvmReaderRepository interface {
	ListApplications(ctx context.Context, f repository.ApplicationFilter, p repository.Pagination) ([]*Application, uint64, error)
	UpdateApplicationState(ctx context.Context, appID int64, state ApplicationState, reason *string) error
	UpdateEventLastCheckBlock(ctx context.Context, appIDs []int64, event MonitoredEvent, blockNumber uint64) error

	CreateEpochsAndInputs(
		ctx context.Context, nameOrAddress string,
		epochInputMap map[*Epoch][]*Input, blockNumber uint64,
		espressoUpdateInfo *model.EspressoUpdateInfo,
	) error
	GetEpoch(ctx context.Context, nameOrAddress string, index uint64) (*Epoch, error)

	GetInputIndex(
		ctx context.Context,
		nameOrAddress string,
	) (uint64, error)
}

// EthClientInterface defines the methods we need from ethclient.Client
type EthClientInterface interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
	ChainID(ctx context.Context) (*big.Int, error)
}

type ApplicationContractAdapter interface {
	RetrieveOutputExecutionEvents(
		opts *bind.FilterOpts,
	) ([]*appcontract.IApplicationOutputExecuted, error)
}

// Internal struct to hold application and it's contracts together
type application struct {
	Application
	ApplicationContractAdapter
	InputSourceAdapter
}

// EvmReader reads Input Added, Claim Submitted and
// Output Executed events from the blockchain
type EvmReader struct {
	client                 EthClientInterface
	wsClient               EthClientInterface
	ExportedAdapterFactory AdapterFactory
	repository             EvmReaderRepository
	chainId                uint64
	defaultBlock           DefaultBlock
	hasEnabledApps         bool
	inputReaderEnabled     bool
	shouldModifyIndex      bool // modify index in raw data if the main sequencer is espresso
	IOAbi                  abi.ABI
}

func (r *EvmReader) String() string {
	return "evmreader"
}

// Creates a new EvmReader
func NewEvmReader(
	client EthClientInterface,
	wsClient EthClientInterface,
	repository EvmReaderRepository,
	chainId uint64,
	defaultBlock DefaultBlock,
	inputReaderEnabled bool,
	shouldModifyIndex bool,
) EvmReader {
	ioABI, err := abi.JSON(strings.NewReader(abiData))
	if err != nil {
		panic(err)
	}
	evmReader := EvmReader{
		client:                 client,
		wsClient:               wsClient,
		ExportedAdapterFactory: defaultAdapterFactory,
		repository:             repository,
		chainId:                chainId,
		defaultBlock:           defaultBlock,
		hasEnabledApps:         true,
		inputReaderEnabled:     inputReaderEnabled,
		shouldModifyIndex:      shouldModifyIndex,
		IOAbi:                  ioABI,
	}
	return evmReader
}

func (r *EvmReader) GetAllRunningApplications(ctx context.Context, er EvmReaderRepository) ([]*Application, uint64, error) {
	f := repository.ApplicationFilter{
		State: Pointer(ApplicationState_Enabled),
	}
	return er.ListApplications(ctx, f, repository.Pagination{})
}

type AdapterFactory interface {
	CreateAdapters(app *Application, client EthClientInterface) (ApplicationContractAdapter, InputSourceAdapter, error)
}

type DefaultAdapterFactory struct{}

func (f *DefaultAdapterFactory) CreateAdapters(app *Application, client EthClientInterface) (ApplicationContractAdapter, InputSourceAdapter, error) {
	if app == nil {
		return nil, nil, fmt.Errorf("Application reference is nil. Should never happen")
	}

	// Type assertion to get the concrete client if possible
	ethClient, ok := client.(*ethclient.Client)
	if !ok {
		return nil, nil, fmt.Errorf("client is not an *ethclient.Client, cannot create adapters")
	}

	applicationContract, err := NewApplicationContractAdapter(app.IApplicationAddress, ethClient)
	if err != nil {
		return nil, nil, errors.Join(
			fmt.Errorf("error building application contract"),
			err,
		)
	}

	inputSource, err := NewInputSourceAdapter(app.IInputBoxAddress, ethClient)
	if err != nil {
		return nil, nil, errors.Join(
			fmt.Errorf("error building inputbox contract"),
			err,
		)
	}

	return applicationContract, inputSource, nil
}

var defaultAdapterFactory = &DefaultAdapterFactory{}

func (r *EvmReader) GetEthClient() *EthClientInterface {
	return &r.client
}
