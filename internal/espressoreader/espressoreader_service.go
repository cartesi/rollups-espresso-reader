// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package espressoreader

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/cartesi/rollups-espresso-reader/internal/evmreader"
	"github.com/cartesi/rollups-espresso-reader/internal/repository"

	"github.com/EspressoSystems/espresso-sequencer-go/client"
	"github.com/EspressoSystems/espresso-sequencer-go/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// app address => sender address => nonce
var nonceCache map[common.Address]map[common.Address]uint64

// Service to manage InputReader lifecycle
type EspressoReaderService struct {
	blockchainHttpEndpoint  string
	blockchainWsEndpoint    string
	database                repository.Repository
	EspressoBaseUrl         string
	maxRetries              uint64
	maxDelay                time.Duration
	chainId                 uint64
	espressoServiceEndpoint string
}

func NewEspressoReaderService(
	blockchainHttpEndpoint string,
	blockchainWsEndpoint string,
	database repository.Repository,
	EspressoBaseUrl string,
	chainId uint64,
	espressoServiceEndpoint string,
	maxRetries uint64,
	maxDelay time.Duration,
) *EspressoReaderService {
	return &EspressoReaderService{
		blockchainHttpEndpoint:  blockchainHttpEndpoint,
		blockchainWsEndpoint:    blockchainWsEndpoint,
		database:                database,
		EspressoBaseUrl:         EspressoBaseUrl,
		chainId:                 chainId,
		espressoServiceEndpoint: espressoServiceEndpoint,
		maxRetries:              maxRetries,
		maxDelay:                maxDelay,
	}
}

func (s *EspressoReaderService) Start(
	ctx context.Context,
	ready chan<- struct{},
) error {

	evmReader := s.setupEvmReader(ctx, s.database)

	espressoReader := NewEspressoReader(s.EspressoBaseUrl, s.database, evmReader, s.chainId, s.maxRetries, uint64(s.maxDelay), s.blockchainHttpEndpoint)

	go s.setupNonceHttpServer()

	return espressoReader.Run(ctx, ready)
}

func (s *EspressoReaderService) String() string {
	return "espressoreader"
}

func (s *EspressoReaderService) setupNonceHttpServer() {
	nonceCache = make(map[common.Address]map[common.Address]uint64)

	http.HandleFunc("/nonce", s.requestNonce)
	http.HandleFunc("/submit", s.submit)

	err := http.ListenAndServe(s.espressoServiceEndpoint, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to start the transaction endpoints /nonce and /submit: %v", err))
	}
	slog.Debug("Transaction service started", "espressoServiceEndpoint", s.espressoServiceEndpoint)
}

type NonceRequest struct {
	// AppContract App contract address
	AppContract string `json:"app_contract"`

	// MsgSender Message sender address
	MsgSender string `json:"msg_sender"`
}

type NonceResponse struct {
	Nonce uint64 `json:"nonce"`
}

func (s *EspressoReaderService) requestNonce(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
	if r.Method != http.MethodPost {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("could not read body", "err", err)
	}
	nonceRequest := &NonceRequest{}
	if err := json.Unmarshal(body, nonceRequest); err != nil {
		slog.Error("could not unmarshal", "err", err)
	}

	senderAddress := common.HexToAddress(nonceRequest.MsgSender)
	applicationAddress := common.HexToAddress(nonceRequest.AppContract)

	var nonce uint64
	if nonceCache[applicationAddress] == nil {
		nonceCache[applicationAddress] = make(map[common.Address]uint64)
	}
	if nonceCache[applicationAddress][senderAddress] == 0 {
		ctx := r.Context()
		nonce = s.queryNonceFromDb(ctx, senderAddress, applicationAddress)
		nonceCache[applicationAddress][senderAddress] = nonce
	} else {
		nonce = nonceCache[applicationAddress][senderAddress]
	}

	slog.Debug("got nonce request", "senderAddress", senderAddress, "applicationAddress", applicationAddress)

	nonceResponse := NonceResponse{Nonce: nonce}
	if err != nil {
		slog.Error("error json marshal nonce response", "err", err)
	}

	err = json.NewEncoder(w).Encode(nonceResponse)
	if err != nil {
		slog.Info("Internal server error",
			"service", "espresso nonce querier",
			"err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *EspressoReaderService) queryNonceFromDb(
	ctx context.Context,
	senderAddress common.Address,
	applicationAddress common.Address) uint64 {
	nonce, err := s.database.GetEspressoNonce(ctx, senderAddress.Hex(), applicationAddress.Hex())
	if err != nil {
		slog.Error("failed to get espresso nonce", "error", err)
	}

	return nonce
}

type SubmitResponse struct {
	Id string `json:"id,omitempty"`
}

func (s *EspressoReaderService) submit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
	if r.Method != http.MethodPost {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("could not read body", "err", err)
	}
	slog.Debug("got submit request", "request body", string(body))

	msgSender, typedData, sigHash, err := ExtractSigAndData(string(body))
	if err != nil {
		slog.Error("transaction not correctly formatted", "error", err)
		return
	}
	submitResponse := SubmitResponse{Id: sigHash}

	appAddressStr := typedData.Message["app"].(string)
	appAddress := common.HexToAddress(appAddressStr)
	client := client.NewClient(s.EspressoBaseUrl)
	ctx := r.Context()
	_, namespace, err := getEspressoConfig(ctx, appAddress, s.database, s.blockchainHttpEndpoint)
	var tx types.Transaction
	tx.Namespace = namespace
	tx.Payload = body
	_, err = client.SubmitTransaction(ctx, tx)
	if err != nil {
		slog.Error("espresso tx submit error", "err", err)
		return
	}

	err = json.NewEncoder(w).Encode(submitResponse)
	if err != nil {
		slog.Info("Internal server error",
			"service", "espresso submit endpoint",
			"err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// update nonce cache
	if nonceCache[appAddress] == nil {
		slog.Error("Should query nonce before submit")
		return
	}
	nonceInRequest := uint64(typedData.Message["nonce"].(float64))
	if nonceCache[appAddress][msgSender] == 0 {
		ctx := r.Context()
		nonceInDb := s.queryNonceFromDb(ctx, msgSender, appAddress)
		if nonceInRequest != nonceInDb {
			slog.Error("Nonce in request is incorrect")
			return
		}
		nonceCache[appAddress][msgSender] = nonceInDb + 1
	} else {
		nonceCache[appAddress][msgSender]++
	}
}

func (s *EspressoReaderService) trySetupEvmReader(ctx context.Context, r repository.Repository) (*evmreader.EvmReader, error) {
	client, err := ethclient.DialContext(ctx, s.blockchainHttpEndpoint)
	if err != nil {
		return nil, fmt.Errorf("eth client http: %w", err)
	}
	defer client.Close()

	wsClient, err := ethclient.DialContext(ctx, s.blockchainWsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("eth client ws: %w", err)
	}
	defer wsClient.Close()

	config, err := repository.LoadNodeConfig[evmreader.PersistentConfig](ctx, r, evmreader.EvmReaderConfigKey)
	if err != nil {
		return nil, fmt.Errorf("db config: %w", err)
	}

	evmReader := evmreader.NewEvmReader(
		client,
		wsClient,
		r,
		config.Value.ChainID,
		config.Value.DefaultBlock,
		config.Value.InputReaderEnabled,
		true,
	)

	return &evmReader, nil
}

func (s *EspressoReaderService) setupEvmReader(ctx context.Context, r repository.Repository) *evmreader.EvmReader {
	evmReader, err := s.trySetupEvmReader(ctx, r)
	if err != nil {
		slog.Error("failed to setup evm reader", "error", err)
		os.Exit(1)
	}
	return evmReader
}
