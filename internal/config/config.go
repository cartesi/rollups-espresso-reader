// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

// The config package manages the node configuration, which comes from environment variables.
// The sub-package generate specifies these environment variables.
package config

// NodeConfig contains all the Node variables.
// See the corresponding environment variable for the variable documentation.
type NodeConfig struct {
	LogLevel                LogLevel
	LogPrettyEnabled        bool
	BlockchainHttpEndpoint  Redacted[string]
	BlockchainWsEndpoint    Redacted[string]
	EvmReaderDefaultBlock   DefaultBlock
	BlockchainBlockTimeout  int
	DatabaseConnection      Redacted[string]
	EspressoBaseUrl         string
	EspressoStartingBlock   uint64
	EspressoNamespace       uint64
	EspressoServiceEndpoint string
	MaxRetries              uint64
	MaxDelay                Duration
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
	config.EvmReaderDefaultBlock = GetEvmReaderDefaultBlock()
	config.DatabaseConnection = Redacted[string]{GetDatabaseConnection()}
	config.EspressoBaseUrl = GetBaseUrl() + "/v0"
	config.EspressoStartingBlock = GetStartingBlock()
	config.EspressoNamespace = GetNamespace()
	config.EspressoServiceEndpoint = GetServiceEndpoint()
	config.MaxRetries = GetPolicyMaxRetries()
	config.MaxDelay = GetPolicyMaxDelay()
	return config
}
