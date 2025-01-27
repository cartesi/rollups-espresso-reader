// Code generated by internal/config/generate.
// DO NOT EDIT.

// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/cartesi/rollups-espresso-reader/internal/model"
)

type (
	Duration     = time.Duration
	LogLevel     = slog.Level
	DefaultBlock = model.DefaultBlock
)

// ------------------------------------------------------------------------------------------------
// Auth Kind
// ------------------------------------------------------------------------------------------------

type AuthKind uint8

const (
	AuthKindPrivateKeyVar AuthKind = iota
	AuthKindPrivateKeyFile
	AuthKindMnemonicVar
	AuthKindMnemonicFile
	AuthKindAWS
)

// ------------------------------------------------------------------------------------------------
// Parsing functions
// ------------------------------------------------------------------------------------------------

func ToInt64FromString(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func ToUint64FromString(s string) (uint64, error) {
	value, err := strconv.ParseUint(s, 10, 64)
	return value, err
}

func ToStringFromString(s string) (string, error) {
	return s, nil
}

func ToDurationFromSeconds(s string) (time.Duration, error) {
	return time.ParseDuration(s + "s")
}

func ToLogLevelFromString(s string) (LogLevel, error) {
	var m = map[string]LogLevel{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
	if v, ok := m[s]; ok {
		return v, nil
	} else {
		var zeroValue LogLevel
		return zeroValue, fmt.Errorf("invalid log level '%s'", s)
	}
}

func ToDefaultBlockFromString(s string) (DefaultBlock, error) {
	var m = map[string]DefaultBlock{
		"latest":    model.DefaultBlock_Latest,
		"pending":   model.DefaultBlock_Pending,
		"safe":      model.DefaultBlock_Safe,
		"finalized": model.DefaultBlock_Finalized,
	}
	if v, ok := m[s]; ok {
		return v, nil
	} else {
		var zeroValue DefaultBlock
		return zeroValue, fmt.Errorf("invalid default block '%s'", s)
	}
}

func ToAuthKindFromString(s string) (AuthKind, error) {
	var m = map[string]AuthKind{
		"private_key":      AuthKindPrivateKeyVar,
		"private_key_file": AuthKindPrivateKeyFile,
		"mnemonic":         AuthKindMnemonicVar,
		"mnemonic_file":    AuthKindMnemonicFile,
		"aws":              AuthKindAWS,
	}
	if v, ok := m[s]; ok {
		return v, nil
	} else {
		var zeroValue AuthKind
		return zeroValue, fmt.Errorf("invalid auth kind '%s'", s)
	}
}

// Aliases to be used by the generated functions.
var (
	toBool         = strconv.ParseBool
	toInt          = strconv.Atoi
	toInt64        = ToInt64FromString
	toUint64       = ToUint64FromString
	toString       = ToStringFromString
	toDuration     = ToDurationFromSeconds
	toLogLevel     = ToLogLevelFromString
	toAuthKind     = ToAuthKindFromString
	toDefaultBlock = ToDefaultBlockFromString
)

// ------------------------------------------------------------------------------------------------
// Getters
// ------------------------------------------------------------------------------------------------

func GetAuthAwsKmsKeyId() string {
	s, ok := os.LookupEnv("CARTESI_AUTH_AWS_KMS_KEY_ID")
	if !ok {
		panic("missing env var CARTESI_AUTH_AWS_KMS_KEY_ID")
	}
	val, err := toString(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_AUTH_AWS_KMS_KEY_ID: %v", err))
	}
	return val
}

func GetAuthAwsKmsRegion() string {
	s, ok := os.LookupEnv("CARTESI_AUTH_AWS_KMS_REGION")
	if !ok {
		panic("missing env var CARTESI_AUTH_AWS_KMS_REGION")
	}
	val, err := toString(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_AUTH_AWS_KMS_REGION: %v", err))
	}
	return val
}

func GetAuthKind() AuthKind {
	s, ok := os.LookupEnv("CARTESI_AUTH_KIND")
	if !ok {
		s = "mnemonic"
	}
	val, err := toAuthKind(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_AUTH_KIND: %v", err))
	}
	return val
}

func GetAuthMnemonic() string {
	s, ok := os.LookupEnv("CARTESI_AUTH_MNEMONIC")
	if !ok {
		panic("missing env var CARTESI_AUTH_MNEMONIC")
	}
	val, err := toString(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_AUTH_MNEMONIC: %v", err))
	}
	return val
}

func GetAuthMnemonicAccountIndex() int {
	s, ok := os.LookupEnv("CARTESI_AUTH_MNEMONIC_ACCOUNT_INDEX")
	if !ok {
		s = "0"
	}
	val, err := toInt(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_AUTH_MNEMONIC_ACCOUNT_INDEX: %v", err))
	}
	return val
}

func GetAuthMnemonicFile() string {
	s, ok := os.LookupEnv("CARTESI_AUTH_MNEMONIC_FILE")
	if !ok {
		panic("missing env var CARTESI_AUTH_MNEMONIC_FILE")
	}
	val, err := toString(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_AUTH_MNEMONIC_FILE: %v", err))
	}
	return val
}

func GetAuthPrivateKey() string {
	s, ok := os.LookupEnv("CARTESI_AUTH_PRIVATE_KEY")
	if !ok {
		panic("missing env var CARTESI_AUTH_PRIVATE_KEY")
	}
	val, err := toString(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_AUTH_PRIVATE_KEY: %v", err))
	}
	return val
}

func GetAuthPrivateKeyFile() string {
	s, ok := os.LookupEnv("CARTESI_AUTH_PRIVATE_KEY_FILE")
	if !ok {
		panic("missing env var CARTESI_AUTH_PRIVATE_KEY_FILE")
	}
	val, err := toString(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_AUTH_PRIVATE_KEY_FILE: %v", err))
	}
	return val
}

func GetBlockchainBlockTimeout() int {
	s, ok := os.LookupEnv("CARTESI_BLOCKCHAIN_BLOCK_TIMEOUT")
	if !ok {
		s = "60"
	}
	val, err := toInt(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_BLOCKCHAIN_BLOCK_TIMEOUT: %v", err))
	}
	return val
}

func GetBlockchainHttpEndpoint() string {
	s, ok := os.LookupEnv("CARTESI_BLOCKCHAIN_HTTP_ENDPOINT")
	if !ok {
		panic("missing env var CARTESI_BLOCKCHAIN_HTTP_ENDPOINT")
	}
	val, err := toString(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_BLOCKCHAIN_HTTP_ENDPOINT: %v", err))
	}
	return val
}

func GetBlockchainWsEndpoint() string {
	s, ok := os.LookupEnv("CARTESI_BLOCKCHAIN_WS_ENDPOINT")
	if !ok {
		panic("missing env var CARTESI_BLOCKCHAIN_WS_ENDPOINT")
	}
	val, err := toString(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_BLOCKCHAIN_WS_ENDPOINT: %v", err))
	}
	return val
}

func GetEvmReaderDefaultBlock() DefaultBlock {
	s, ok := os.LookupEnv("CARTESI_EVM_READER_DEFAULT_BLOCK")
	if !ok {
		s = "finalized"
	}
	val, err := toDefaultBlock(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_EVM_READER_DEFAULT_BLOCK: %v", err))
	}
	return val
}

func GetLegacyBlockchainEnabled() bool {
	s, ok := os.LookupEnv("CARTESI_LEGACY_BLOCKCHAIN_ENABLED")
	if !ok {
		s = "false"
	}
	val, err := toBool(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_LEGACY_BLOCKCHAIN_ENABLED: %v", err))
	}
	return val
}

func GetBaseUrl() string {
	s, ok := os.LookupEnv("ESPRESSO_BASE_URL")
	if !ok {
		s = ""
	}
	val, err := toString(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse ESPRESSO_BASE_URL: %v", err))
	}
	return val
}

func GetNamespace() uint64 {
	s, ok := os.LookupEnv("ESPRESSO_NAMESPACE")
	if !ok {
		s = "0"
	}
	val, err := toUint64(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse ESPRESSO_NAMESPACE: %v", err))
	}
	return val
}

func GetServiceEndpoint() string {
	s, ok := os.LookupEnv("ESPRESSO_SERVICE_ENDPOINT")
	if !ok {
		s = "localhost:8080"
	}
	val, err := toString(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse ESPRESSO_SERVICE_ENDPOINT: %v", err))
	}
	return val
}

func GetStartingBlock() uint64 {
	s, ok := os.LookupEnv("ESPRESSO_STARTING_BLOCK")
	if !ok {
		s = "0"
	}
	val, err := toUint64(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse ESPRESSO_STARTING_BLOCK: %v", err))
	}
	return val
}

func GetFeatureClaimSubmissionEnabled() bool {
	s, ok := os.LookupEnv("CARTESI_FEATURE_CLAIM_SUBMISSION_ENABLED")
	if !ok {
		s = "true"
	}
	val, err := toBool(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_FEATURE_CLAIM_SUBMISSION_ENABLED: %v", err))
	}
	return val
}

func GetFeatureMachineHashCheckEnabled() bool {
	s, ok := os.LookupEnv("CARTESI_FEATURE_MACHINE_HASH_CHECK_ENABLED")
	if !ok {
		s = "true"
	}
	val, err := toBool(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_FEATURE_MACHINE_HASH_CHECK_ENABLED: %v", err))
	}
	return val
}

func GetHttpAddress() string {
	s, ok := os.LookupEnv("CARTESI_HTTP_ADDRESS")
	if !ok {
		s = "127.0.0.1"
	}
	val, err := toString(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_HTTP_ADDRESS: %v", err))
	}
	return val
}

func GetHttpPort() int {
	s, ok := os.LookupEnv("CARTESI_HTTP_PORT")
	if !ok {
		s = "10000"
	}
	val, err := toInt(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_HTTP_PORT: %v", err))
	}
	return val
}

func GetLogLevel() LogLevel {
	s, ok := os.LookupEnv("CARTESI_LOG_LEVEL")
	if !ok {
		s = "info"
	}
	val, err := toLogLevel(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_LOG_LEVEL: %v", err))
	}
	return val
}

func GetLogPrettyEnabled() bool {
	s, ok := os.LookupEnv("CARTESI_LOG_PRETTY_ENABLED")
	if !ok {
		s = "false"
	}
	val, err := toBool(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_LOG_PRETTY_ENABLED: %v", err))
	}
	return val
}

func GetPostgresEndpoint() string {
	s, ok := os.LookupEnv("CARTESI_POSTGRES_ENDPOINT")
	if !ok {
		s = ""
	}
	val, err := toString(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_POSTGRES_ENDPOINT: %v", err))
	}
	return val
}

func GetAdvancerPollingInterval() Duration {
	s, ok := os.LookupEnv("CARTESI_ADVANCER_POLLING_INTERVAL")
	if !ok {
		s = "7"
	}
	val, err := toDuration(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_ADVANCER_POLLING_INTERVAL: %v", err))
	}
	return val
}

func GetClaimerPollingInterval() Duration {
	s, ok := os.LookupEnv("CARTESI_CLAIMER_POLLING_INTERVAL")
	if !ok {
		s = "7"
	}
	val, err := toDuration(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_CLAIMER_POLLING_INTERVAL: %v", err))
	}
	return val
}

func GetEvmReaderRetryPolicyMaxDelay() Duration {
	s, ok := os.LookupEnv("CARTESI_EVM_READER_RETRY_POLICY_MAX_DELAY")
	if !ok {
		s = "3"
	}
	val, err := toDuration(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_EVM_READER_RETRY_POLICY_MAX_DELAY: %v", err))
	}
	return val
}

func GetEvmReaderRetryPolicyMaxRetries() uint64 {
	s, ok := os.LookupEnv("CARTESI_EVM_READER_RETRY_POLICY_MAX_RETRIES")
	if !ok {
		s = "3"
	}
	val, err := toUint64(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_EVM_READER_RETRY_POLICY_MAX_RETRIES: %v", err))
	}
	return val
}

func GetValidatorPollingInterval() Duration {
	s, ok := os.LookupEnv("CARTESI_VALIDATOR_POLLING_INTERVAL")
	if !ok {
		s = "7"
	}
	val, err := toDuration(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_VALIDATOR_POLLING_INTERVAL: %v", err))
	}
	return val
}

func GetSnapshotDir() string {
	s, ok := os.LookupEnv("CARTESI_SNAPSHOT_DIR")
	if !ok {
		s = "/var/lib/cartesi-rollups-node/snapshots"
	}
	val, err := toString(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_SNAPSHOT_DIR: %v", err))
	}
	return val
}

func GetMachineServerVerbosity() string {
	s, ok := os.LookupEnv("CARTESI_MACHINE_SERVER_VERBOSITY")
	if !ok {
		s = "info"
	}
	val, err := toString(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse CARTESI_MACHINE_SERVER_VERBOSITY: %v", err))
	}
	return val
}
