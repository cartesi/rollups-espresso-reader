// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package espressoreader

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

type EspressoHelper struct{}

func (eh *EspressoHelper) readEspressoHeader(espressoBlockHeight uint64, url string) string {
	requestURL := fmt.Sprintf("%s/availability/header/%d", url, espressoBlockHeight)
	res, err := http.Get(requestURL)
	if err != nil {
		slog.Error("error making http request", "err", err)
		return "" // it's not the best way to do it
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("could not read response body", "err", err)
		return "" // it's not the best way to do it
	}

	return string(resBody)
}

func (eh *EspressoHelper) getL1FinalizedHeight(ctx context.Context, espressoBlockHeight uint64, delay uint64, url string) (uint64, uint64) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("exiting espresso reader")
			return 0, 0
		default:
			espressoHeader := eh.readEspressoHeader(espressoBlockHeight, url)
			if len(espressoHeader) == 0 {
				slog.Error("error fetching espresso header", "at height", espressoBlockHeight, "header", espressoHeader)
				slog.Error("retrying fetching header")
				time.Sleep(time.Duration(delay))
				continue
			}

			l1FinalizedNumber := gjson.Get(espressoHeader, "fields.l1_finalized.number").Uint()
			l1FinalizedTimestampStr := gjson.Get(espressoHeader, "fields.l1_finalized.timestamp").Str
			if len(l1FinalizedTimestampStr) < 2 {
				slog.Debug("Espresso header not ready. Retry fetching", "height", espressoBlockHeight)
				time.Sleep(time.Duration(delay))
				continue
			}
			l1FinalizedTimestampInt, err := strconv.ParseInt(l1FinalizedTimestampStr[2:], 16, 64)
			if err != nil {
				slog.Error("hex to int conversion failed", "err", err)
				slog.Error("retrying")
				time.Sleep(time.Duration(delay))
				continue
			}
			l1FinalizedTimestamp := uint64(l1FinalizedTimestampInt)
			return l1FinalizedNumber, l1FinalizedTimestamp
		}
	}
}

func (eh *EspressoHelper) readEspressoHeadersByRange(ctx context.Context, from uint64, until uint64, delay uint64, url string) string {
	for {
		select {
		case <-ctx.Done():
			slog.Info("exiting espresso reader")
			return ""
		default:
			requestURL := fmt.Sprintf("%s/availability/header/%d/%d", url, from, until)
			res, err := http.Get(requestURL)
			if err != nil {
				slog.Error("error making http request", "err", err)
				slog.Error("retrying")
				time.Sleep(time.Duration(delay))
				continue
			}
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				slog.Error("could not read response body", "err", err)
				slog.Error("retrying")
				time.Sleep(time.Duration(delay))
				continue
			}

			return string(resBody)
		}
	}
}

func (eh *EspressoHelper) getNSTableByRange(ctx context.Context, from uint64, until uint64, delay uint64, url string) (string, error) {
	var nsTables string
	for len(nsTables) == 0 {
		select {
		case <-ctx.Done():
			slog.Info("exiting espresso reader")
			return "", ctx.Err()
		default:
			espressoHeaders := eh.readEspressoHeadersByRange(ctx, from, until, delay, url)
			nsTables = gjson.Get(espressoHeaders, "#.fields.ns_table.bytes").Raw
			if len(nsTables) == 0 {
				slog.Debug("ns table is empty in current block range. Retry fetching")
				var delay time.Duration = 2000
				time.Sleep(delay * time.Millisecond)
			}
		}
	}

	return nsTables, nil
}

func (eh *EspressoHelper) extractNS(nsTable []byte) []uint32 {
	var nsArray []uint32
	numNS := binary.LittleEndian.Uint32(nsTable[0:])
	for i := range numNS {
		nextNS := binary.LittleEndian.Uint32(nsTable[(4 + 8*i):])
		nsArray = append(nsArray, nextNS)
	}
	return nsArray
}
