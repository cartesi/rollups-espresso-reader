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

	"github.com/EspressoSystems/espresso-network-go/client"
	"github.com/EspressoSystems/espresso-network-go/types"
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
	maxBlockRange           uint64
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
	maxBlockRange uint64,
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
		maxBlockRange:           maxBlockRange,
	}
}

func (s *EspressoReaderService) Start(
	ctx context.Context,
	ready chan<- struct{},
) error {

	evmReader := s.setupEvmReader(ctx, s.database)

	espressoReader := NewEspressoReader(s.EspressoBaseUrl, s.database, evmReader, s.chainId, s.maxRetries, uint64(s.maxDelay))

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
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Only POST method is allowed"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("could not read request body", "error", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	nonceRequest := &NonceRequest{}
	if err := json.Unmarshal(body, nonceRequest); err != nil {
		slog.Error("could not unmarshal request body", "error", err, "body", string(body))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	senderAddress := common.HexToAddress(nonceRequest.MsgSender)
	applicationAddress := common.HexToAddress(nonceRequest.AppContract)

	var nonce uint64
	if nonceCache[applicationAddress] == nil {
		nonceCache[applicationAddress] = make(map[common.Address]uint64)
	}
	_, exists := nonceCache[applicationAddress][senderAddress]
	if !exists {
		ctx := r.Context()
		nonce, err = s.queryNonceFromDb(ctx, senderAddress, applicationAddress)
		if err != nil {
			slog.Error("failed to query nonce from database", "error", err, "senderAddress", senderAddress, "applicationAddress", applicationAddress)
			http.Error(w, "Failed to retrieve nonce", http.StatusInternalServerError)
			return
		}
		nonceCache[applicationAddress][senderAddress] = nonce
	} else {
		nonce = nonceCache[applicationAddress][senderAddress]
	}

	slog.Debug("got nonce request", "senderAddress", senderAddress, "applicationAddress", applicationAddress)

	nonceResponse := NonceResponse{Nonce: nonce}

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
	applicationAddress common.Address) (uint64, error) {
	nonce, err := s.database.GetEspressoNonce(ctx, senderAddress.Hex(), applicationAddress.Hex())
	if err != nil {
		slog.Error("failed to get espresso nonce", "error", err)
		return 0, err
	}

	return nonce, nil
}

type SubmitResponse struct {
	Id string `json:"id,omitempty"`
}

func (s *EspressoReaderService) submit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Only POST method is allowed"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("could not read request body", "error", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	slog.Debug("got submit request", "request body", string(body))

	msgSender, typedData, sigHash, err := ExtractSigAndData(string(body))
	if err != nil {
		slog.Error("transaction not correctly formatted", "error", err)
		http.Error(w, "Invalid transaction format", http.StatusBadRequest)
		return
	}
	submitResponse := SubmitResponse{Id: sigHash}

	appAddressStr := typedData.Message["app"].(string)
	appAddress := common.HexToAddress(appAddressStr)
	client := client.NewClient(s.EspressoBaseUrl)
	ctx := r.Context()
	app, err := s.database.GetApplication(ctx, appAddressStr)
	if err != nil || app == nil {
		slog.Error("application not registered", "err", err)
		http.Error(w, "application not registered", http.StatusInternalServerError)
		return
	}
	_, namespace, err := getEspressoConfig(ctx, appAddress, s.database, app.DataAvailability)
	if err != nil {
		slog.Error("failed to get espresso config", "error", err, "appAddress", appAddress)
		http.Error(w, "Failed to get application configuration", http.StatusInternalServerError)
		return
	}
	var tx types.Transaction
	tx.Namespace = namespace
	tx.Payload = body
	_, err = client.SubmitTransaction(ctx, tx)
	if err != nil {
		slog.Error("espresso tx submit error", "error", err)
		http.Error(w, "Failed to submit transaction to Espresso", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(submitResponse); err != nil {
		slog.Error("failed to encode submit response", "error", err, "submitResponse", submitResponse)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// update nonce cache
	if nonceCache[appAddress] == nil {
		slog.Error("Should query nonce before submit", "appAddress", appAddress)
		http.Error(w, "Nonce not initialized for this application. Please request a nonce first.", http.StatusBadRequest)
		return
	}
	nonceInRequestFloat, ok := typedData.Message["nonce"].(float64)
	if !ok {
		slog.Error("Nonce not found or invalid type in request", "typedData", typedData)
		http.Error(w, "Invalid or missing nonce in request", http.StatusBadRequest)
		return
	}
	nonceInRequest := uint64(nonceInRequestFloat)

	cachedNonce, exists := nonceCache[appAddress][msgSender]
	if !exists {
		ctx := r.Context()
		nonceInDb, err := s.queryNonceFromDb(ctx, msgSender, appAddress)
		if err != nil {
			slog.Error("failed to query nonce from database during submit", "error", err, "senderAddress", msgSender, "applicationAddress", appAddress)
			http.Error(w, "Failed to retrieve nonce for validation", http.StatusInternalServerError)
			return
		}
		if nonceInRequest != nonceInDb {
			slog.Error("Nonce in request is incorrect", "nonceInRequest", nonceInRequest, "nonceInDb", nonceInDb, "senderAddress", msgSender, "applicationAddress", appAddress)
			http.Error(w, "Incorrect nonce", http.StatusBadRequest)
			return
		}
		nonceCache[appAddress][msgSender] = nonceInDb + 1
	} else {
		if nonceInRequest != cachedNonce {
			slog.Error("Nonce in request is incorrect", "requestNonce", nonceInRequest, "cachedNonce", cachedNonce, "senderAddress", msgSender, "applicationAddress", appAddress)
			http.Error(w, "Incorrect nonce", http.StatusBadRequest)
			return
		}
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
		s.maxBlockRange,
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
