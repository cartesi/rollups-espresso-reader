// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package retrypolicy

import (
	"time"

	"github.com/cartesi/rollups-espresso-reader/internal/evmreader"
	"github.com/cartesi/rollups-espresso-reader/internal/services/retry"
	"github.com/cartesi/rollups-espresso-reader/pkg/contracts/iapplication"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type ApplicationRetryPolicyDelegator struct {
	delegate          evmreader.ApplicationContract
	maxRetries        uint64
	delayBetweenCalls time.Duration
}

func NewApplicationWithRetryPolicy(
	delegate evmreader.ApplicationContract,
	maxRetries uint64,
	delayBetweenCalls time.Duration,
) *ApplicationRetryPolicyDelegator {
	return &ApplicationRetryPolicyDelegator{
		delegate:          delegate,
		maxRetries:        maxRetries,
		delayBetweenCalls: delayBetweenCalls,
	}
}

func (d *ApplicationRetryPolicyDelegator) RetrieveOutputExecutionEvents(
	opts *bind.FilterOpts,
) ([]*iapplication.IApplicationOutputExecuted, error) {
	return retry.CallFunctionWithRetryPolicy(d.delegate.RetrieveOutputExecutionEvents,
		opts,
		d.maxRetries,
		d.delayBetweenCalls,
		"Application::RetrieveOutputExecutionEvents",
	)
}
