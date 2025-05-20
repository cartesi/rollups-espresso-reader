// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)
package ethutil

import (
	"context"
	"iter"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	minChunk = new(big.Int).SetInt64(64)
)

// read chunkedFilterLogs comment for additional information.
//
// NOTE: There is no standard reply among providers, add as needed. To handle a
// new provider add it to the table below and make queryBlockRangeTooLarge
// return true when encountering its RPC error code.
// ┌───────────────────────────────────────────────────┬───────┬────────┬────────────┐
// │          provider                                 │ limit │  code  │ checked at │
// ├───────────────────────────────────────────────────┼───────┼────────┼────────────┤
// │ https://cloudflare-eth.com/                       │   800 │ -32047 │ 2025-01-24 │
// │ https://eth-mainnet.g.alchemy.com/v2/{key} (free) │   500 │ -32600 │ 2025-05-13 │
// │ https://mainnet.infura.io/v3/{key} (free)         │ 10000 │ -32005 │ 2025-05-15 │
// │ https://site1.moralis-nodes.com/eth/{key} (free)  │   100 │    400 │ 2025-05-15 │
// └───────────────────────────────────────────────────┴───────┴────────┴────────────┘
func queryBlockRangeTooLarge(err error) bool {
	if err != nil {
		switch e := err.(type) {
		case rpc.Error:
			return (e.ErrorCode() == -32047) || // cloudflare (free)
				(e.ErrorCode() == -32600) || // alchemy (free)
				(e.ErrorCode() == -32005) || // infura (free)
				(e.ErrorCode() == 400) // moralis (free)
		}
	}
	return false
}

// chunkedFilterLogs is very similar to LogFilterer FilterLogs. Both functions
// query blockchain events (logs) and return the ones matching the filter
// criteria. In addition to the basic functionality, this version splits large
// (From, To) block ranges into multiple smaller calls when it detects the
// provider rejected the query for this specific reason. Detection is a
// heuristic and implemented in the function queryBlockRangeTooLarge. It
// potentially has to be adjusted to accomodate each provider.
func ChunkedFilterLogs(
	ctx context.Context,
	client *ethclient.Client,
	q ethereum.FilterQuery,
) (
	iter.Seq2[*types.Log, error],
	error,
) {
	if q.FromBlock == nil {
		q.FromBlock = big.NewInt(0)
	}
	if q.ToBlock == nil {
		end, err := client.BlockNumber(ctx)
		if err != nil {
			return nil, err
		}
		q.ToBlock = big.NewInt(0).SetUint64(end)
	}

	return func(yield func(log *types.Log, err error) bool) {
		one := big.NewInt(1)
		endBlock := new(big.Int).Set(q.ToBlock)
		for q.FromBlock.Cmp(endBlock) <= 0 {
			logs, err := client.FilterLogs(ctx, q)
			delta := new(big.Int).Sub(q.ToBlock, q.FromBlock)

			if queryBlockRangeTooLarge(err) {
				if delta.Cmp(minChunk) < 0 {
					yield(nil, err)
					return
				}
				// ToBlock -= ToBlock/2
				q.ToBlock.Sub(q.ToBlock, delta.Rsh(delta, 1))
				continue
			} else if err != nil {
				yield(nil, err)
				return
			}

			for _, log := range logs {
				if !yield(&log, nil) {
					return
				}
			}

			// ------------------------
			//  [  delta  |  delta  ]
			//  ^         ^
			// From       To
			// ------------------------
			//  [  delta  |  delta  ]
			//             ^        ^
			//            From      To
			// ------------------------
			q.FromBlock.Add(q.ToBlock, one)
			q.ToBlock.Add(q.FromBlock, delta)
			if q.ToBlock.Cmp(endBlock) > 0 {
				q.ToBlock.Set(endBlock)
			}
		}
	}, nil
}
