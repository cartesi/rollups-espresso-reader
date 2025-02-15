// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// The config package manages the node configuration, which comes from environment variables.
// The sub-package generate specifies these environment variables.
package config

import (
	"fmt"
	"os"
)

// NodeConfig contains all the Node variables.
// See the corresponding environment variable for the variable documentation.
type NodeConfig struct {
	LogLevel                       LogLevel
	LogPrettyEnabled               bool
	BlockchainHttpEndpoint         Redacted[string]
	BlockchainWsEndpoint           Redacted[string]
	LegacyBlockchainEnabled        bool
	EvmReaderDefaultBlock          DefaultBlock
	BlockchainBlockTimeout         int
	SnapshotDir                    string
	PostgresEndpoint               Redacted[string]
	HttpAddress                    string
	HttpPort                       int
	FeatureClaimSubmissionEnabled  bool
	FeatureMachineHashCheckEnabled bool
	Auth                           Auth
	AdvancerPollingInterval        Duration
	ValidatorPollingInterval       Duration
	ClaimerPollingInterval         Duration
	EspressoBaseUrl                string
	EspressoStartingBlock          uint64
	EspressoNamespace              uint64
	EspressoServiceEndpoint        string
	MaxRetries                     uint64
	MaxDelay                       Duration
}

// Auth is used to sign transactions.
type Auth any

// AuthPrivateKey allows signing through private keys.
type AuthPrivateKey struct {
	PrivateKey Redacted[string]
}

// AuthMnemonic allows signing through mnemonics.
type AuthMnemonic struct {
	Mnemonic     Redacted[string]
	AccountIndex Redacted[int]
}

// AuthAWS allows signing through AWS services.
type AuthAWS struct {
	KeyID  Redacted[string]
	Region Redacted[string]
}

// Redacted is a wrapper that redacts a given field from the logs.
type Redacted[T any] struct {
	Value T
}

func (r Redacted[T]) String() string {
	return "[REDACTED]"
}

// FromEnv loads the config from environment variables.
func FromEnv() NodeConfig {
	var config NodeConfig
	config.LogLevel = GetLogLevel()
	config.LogPrettyEnabled = GetLogPrettyEnabled()
	config.BlockchainHttpEndpoint = Redacted[string]{GetBlockchainHttpEndpoint()}
	config.BlockchainWsEndpoint = Redacted[string]{GetBlockchainWsEndpoint()}
	config.LegacyBlockchainEnabled = GetLegacyBlockchainEnabled()
	config.EvmReaderDefaultBlock = GetEvmReaderDefaultBlock()
	config.BlockchainBlockTimeout = GetBlockchainBlockTimeout()
	config.SnapshotDir = GetSnapshotDir()
	config.PostgresEndpoint = Redacted[string]{GetPostgresEndpoint()}
	config.HttpAddress = GetHttpAddress()
	config.HttpPort = GetHttpPort()
	config.FeatureClaimSubmissionEnabled = GetFeatureClaimSubmissionEnabled()
	config.FeatureMachineHashCheckEnabled = GetFeatureMachineHashCheckEnabled()
	if config.FeatureClaimSubmissionEnabled {
		config.Auth = AuthFromEnv()
	}
	config.AdvancerPollingInterval = GetAdvancerPollingInterval()
	config.ValidatorPollingInterval = GetValidatorPollingInterval()
	config.ClaimerPollingInterval = GetClaimerPollingInterval()
	config.EspressoBaseUrl = GetBaseUrl() + "/v0"
	config.EspressoStartingBlock = GetStartingBlock()
	config.EspressoNamespace = GetNamespace()
	config.EspressoServiceEndpoint = GetServiceEndpoint()
	config.MaxRetries = GetPolicyMaxRetries()
	config.MaxDelay = GetPolicyMaxDelay()
	return config
}

func AuthFromEnv() Auth {
	switch GetAuthKind() {
	case AuthKindPrivateKeyVar:
		return AuthPrivateKey{
			PrivateKey: Redacted[string]{GetAuthPrivateKey()},
		}
	case AuthKindPrivateKeyFile:
		path := GetAuthPrivateKeyFile()
		privateKey, err := os.ReadFile(path)
		if err != nil {
			panic(fmt.Sprintf("failed to read private-key file: %v", err))
		}
		return AuthPrivateKey{
			PrivateKey: Redacted[string]{string(privateKey)},
		}
	case AuthKindMnemonicVar:
		return AuthMnemonic{
			Mnemonic:     Redacted[string]{GetAuthMnemonic()},
			AccountIndex: Redacted[int]{GetAuthMnemonicAccountIndex()},
		}
	case AuthKindMnemonicFile:
		path := GetAuthMnemonicFile()
		mnemonic, err := os.ReadFile(path)
		if err != nil {
			panic(fmt.Sprintf("failed to read mnemonic file: %v", err))
		}
		return AuthMnemonic{
			Mnemonic:     Redacted[string]{string(mnemonic)},
			AccountIndex: Redacted[int]{GetAuthMnemonicAccountIndex()},
		}
	case AuthKindAWS:
		return AuthAWS{
			KeyID:  Redacted[string]{GetAuthAwsKmsKeyId()},
			Region: Redacted[string]{GetAuthAwsKmsRegion()},
		}
	default:
		panic("invalid auth kind")
	}
}
