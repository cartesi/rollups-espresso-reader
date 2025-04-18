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
type InputSource interface {
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

// EthClient mimics part of ethclient.Client functions to narrow down the
// interface needed by the EvmReader. It must be bound to an HTTP endpoint
type EthClient interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

// EthWsClient mimics part of ethclient.Client functions to narrow down the
// interface needed by the EvmReader. It must be bound to a WS endpoint
type EthWsClient interface {
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
}

type ApplicationContract interface {
	RetrieveOutputExecutionEvents(
		opts *bind.FilterOpts,
	) ([]*appcontract.IApplicationOutputExecuted, error)
}

type ContractFactory interface {
	NewApplication(address common.Address) (ApplicationContract, error)
	NewInputSource(address common.Address) (InputSource, error)
}

// Internal struct to hold application and it's contracts together
type application struct {
	Application
	ApplicationContract
	InputSource
}

// EvmReader reads Input Added, Claim Submitted and
// Output Executed events from the blockchain
type EvmReader struct {
	client            EthClient
	wsClient          EthWsClient
	repository        EvmReaderRepository
	contractFactory   ContractFactory
	defaultBlock      DefaultBlock
	hasEnabledApps    bool
	shouldModifyIndex bool // modify index in raw data if the main sequencer is espresso
	IOAbi             abi.ABI
}

func (r *EvmReader) String() string {
	return "evmreader"
}

// Creates a new EvmReader
func NewEvmReader(
	client EthClient,
	wsClient EthWsClient,
	repository EvmReaderRepository,
	defaultBlock DefaultBlock,
	contractFactory ContractFactory,
	shouldModifyIndex bool,
) EvmReader {
	ioABI, err := abi.JSON(strings.NewReader(abiData))
	if err != nil {
		panic(err)
	}
	evmReader := EvmReader{
		client:            client,
		wsClient:          wsClient,
		repository:        repository,
		defaultBlock:      defaultBlock,
		contractFactory:   contractFactory,
		hasEnabledApps:    true,
		shouldModifyIndex: shouldModifyIndex,
		IOAbi:             ioABI,
	}
	return evmReader
}

func (r *EvmReader) GetAllRunningApplications(ctx context.Context, er EvmReaderRepository) ([]*Application, uint64, error) {
	f := repository.ApplicationFilter{
		State:            Pointer(ApplicationState_Enabled),
		DataAvailability: &model.DataAvailability_InputBoxAndEspresso,
	}
	return er.ListApplications(ctx, f, repository.Pagination{})
}

// GetAppContracts retrieves the ApplicationContract and ConsensusContract for a given Application.
// Also validates if IConsensus configuration matches the blockchain registered one
func (r *EvmReader) GetAppContracts(app *Application,
) (ApplicationContract, InputSource, error) {
	if app == nil {
		return nil, nil, fmt.Errorf("Application reference is nil. Should never happen")
	}

	applicationContract, err := r.contractFactory.NewApplication(app.IApplicationAddress)
	if err != nil {
		return nil, nil, errors.Join(
			fmt.Errorf("error building application contract"),
			err,
		)

	}

	inputSource, err := r.contractFactory.NewInputSource(app.IInputBoxAddress)
	if err != nil {
		return nil, nil, errors.Join(
			fmt.Errorf("error building inputbox contract"),
			err,
		)

	}

	return applicationContract, inputSource, nil
}

func (r *EvmReader) GetEthClient() *EthClient {
	return &r.client
}
