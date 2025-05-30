package espressoreader

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type SigAndData struct {
	TypedData apitypes.TypedData `json:"typedData"`
	Account   string             `json:"account"`
	Signature string             `json:"signature"`
}

type SigAndDataV0 struct {
	TypedData string `json:"typedData"`
	Account   string `json:"account"`
	Signature string `json:"signature"`
}

func ExtractSigAndData(raw string) (common.Address, apitypes.TypedData, string, error) {
	var sigAndData SigAndData
	if err := json.Unmarshal([]byte(raw), &sigAndData); err != nil {
		slog.Warn("Unmarshal error", "raw", raw)
		var sigAndDataV0 SigAndDataV0
		if err0 := json.Unmarshal([]byte(raw), &sigAndDataV0); err0 != nil {
			return common.HexToAddress("0x"), apitypes.TypedData{}, "", fmt.Errorf("unmarshal sigAndData: %w", err)
		}
		sigAndData.Account = sigAndDataV0.Account
		sigAndData.Signature = sigAndDataV0.Signature
		if err1 := json.Unmarshal([]byte(sigAndDataV0.TypedData), &sigAndData.TypedData); err1 != nil {
			return common.HexToAddress("0x"), apitypes.TypedData{}, "", fmt.Errorf("unmarshal sigAndData: %w", err)
		}
	}

	signature, err := hexutil.Decode(sigAndData.Signature)
	if err != nil {
		return common.HexToAddress("0x"), apitypes.TypedData{}, "", fmt.Errorf("decode signature: %w", err)
	}
	sigHash := crypto.Keccak256Hash(signature).String()

	typedData := sigAndData.TypedData
	dataHash, _, err := apitypes.TypedDataAndHash(typedData)
	if err != nil {
		return common.HexToAddress("0x"), apitypes.TypedData{}, "", fmt.Errorf("typed data hash: %w", err)
	}

	// update the recovery id
	// https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L442
	signature[64] -= 27

	// get the pubkey used to sign this signature
	sigPubkey, err := crypto.Ecrecover(dataHash, signature)
	if err != nil {
		return common.HexToAddress("0x"), apitypes.TypedData{}, "", fmt.Errorf("ecrecover: %w", err)
	}
	pubkey, err := crypto.UnmarshalPubkey(sigPubkey)
	if err != nil {
		return common.HexToAddress("0x"), apitypes.TypedData{}, "", fmt.Errorf("unmarshal: %w", err)
	}
	address := crypto.PubkeyToAddress(*pubkey)

	return address, typedData, sigHash, nil
}
