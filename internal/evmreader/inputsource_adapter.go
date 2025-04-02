// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package evmreader

import (
	"math/big"

	. "github.com/cartesi/rollups-espresso-reader/internal/model"
	"github.com/cartesi/rollups-espresso-reader/pkg/contracts/iinputbox"
	"github.com/cartesi/rollups-espresso-reader/pkg/ethutil"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// InputBox Wrapper
type InputSourceAdapter struct {
	inputbox        *iinputbox.IInputBox
	client          *ethclient.Client
	inputBoxAddress common.Address
}

func NewInputSourceAdapter(
	inputBoxAddress common.Address,
	client *ethclient.Client,
) (*InputSourceAdapter, error) {
	inputbox, err := iinputbox.NewIInputBox(inputBoxAddress, client)
	if err != nil {
		return nil, err
	}
	return &InputSourceAdapter{
		inputbox:        inputbox,
		client:          client,
		inputBoxAddress: inputBoxAddress,
	}, nil
}

func buildInputAddedFilterQuery(
	opts *bind.FilterOpts,
	inputBoxAddress common.Address,
	appContract []common.Address,
	index []*big.Int,
) (q ethereum.FilterQuery, err error) {
	c, err := iinputbox.IInputBoxMetaData.GetAbi()
	if err != nil {
		return q, err
	}

	var appContractRule []any
	for _, appContractItem := range appContract {
		appContractRule = append(appContractRule, appContractItem)
	}
	var indexRule []any
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}

	topics, err := abi.MakeTopics(
		[]any{c.Events[MonitoredEvent_InputAdded.String()].ID},
		appContractRule,
		indexRule,
	)
	if err != nil {
		return q, err
	}

	q = ethereum.FilterQuery{
		Addresses: []common.Address{inputBoxAddress},
		FromBlock: new(big.Int).SetUint64(opts.Start),
		Topics:    topics,
	}
	if opts.End != nil {
		q.ToBlock = new(big.Int).SetUint64(*opts.End)
	}
	return q, err
}

func (i *InputSourceAdapter) RetrieveInputs(
	opts *bind.FilterOpts,
	appContract []common.Address,
	index []*big.Int,
) ([]iinputbox.IInputBoxInputAdded, error) {
	q, err := buildInputAddedFilterQuery(opts, i.inputBoxAddress, appContract, index)
	if err != nil {
		return nil, err
	}

	itr, err := ethutil.ChunkedFilterLogs(opts.Context, i.client, q)
	if err != nil {
		return nil, err
	}

	var events []iinputbox.IInputBoxInputAdded
	for log, err := range itr {
		if err != nil {
			return nil, err
		}
		ev, err := i.inputbox.ParseInputAdded(*log)
		if err != nil {
			return nil, err
		}
		events = append(events, *ev)
	}
	return events, nil
}
