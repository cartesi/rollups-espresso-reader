// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package espressoreader

// Basic imports
import (
	"context"
	"encoding/hex"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/cartesi/rollups-espresso-reader/internal/config"
	"github.com/cartesi/rollups-espresso-reader/internal/model"
	"github.com/cartesi/rollups-espresso-reader/internal/repository"
	"github.com/cartesi/rollups-espresso-reader/internal/repository/factory"
	"github.com/cartesi/rollups-espresso-reader/internal/services/startup"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/suite"
)

type EspressoReaderTestSuite struct {
	suite.Suite
	ctx      context.Context
	c        config.NodeConfig
	database repository.Repository
}

var application model.Application

func (suite *EspressoReaderTestSuite) SetupTest() {
	appAddress := "0xdDedD051B0EeeF07A5099eE572Af85Cb816B3fdE"
	consensusAddress := "0xb66b8260646df96E20E3745489BAF0108E842B8e"

	templatePath := "applications/echo-dapp/"
	templateHash := "0x1a456963bec81fdb3eddeddb2881281737b1136e744f856d6aa43a8dcf920cd5"
	application = model.Application{
		Name:                 "test-dapp",
		IApplicationAddress:  strings.ToLower(common.HexToAddress(appAddress).String()),
		IConsensusAddress:    strings.ToLower(common.HexToAddress(consensusAddress).String()),
		TemplateURI:          templatePath,
		TemplateHash:         common.HexToHash(templateHash),
		State:                model.ApplicationState_Enabled,
		LastProcessedBlock:   0,
		LastOutputCheckBlock: 0,
		LastClaimCheckBlock:  0,
	}

	suite.ctx = context.Background()
	suite.c = config.FromEnv()
	suite.database, _ = factory.NewRepositoryFromConnectionString(suite.ctx, suite.c.PostgresEndpoint.Value)
	_, err := suite.database.CreateApplication(suite.ctx, &application)
	suite.Nil(err)

	startup.SetupNodeConfig(suite.ctx, suite.database, suite.c)
	service := NewEspressoReaderService(
		suite.c.BlockchainHttpEndpoint.Value,
		suite.c.BlockchainHttpEndpoint.Value,
		suite.database,
		suite.c.EspressoBaseUrl,
		suite.c.EspressoStartingBlock,
		suite.c.EspressoNamespace,
		suite.c.EvmReaderRetryPolicyMaxRetries,
		suite.c.EvmReaderRetryPolicyMaxDelay,
		suite.c.BlockchainID,
		uint64(suite.c.ContractsInputBoxDeploymentBlockNumber),
		suite.c.EspressoServiceEndpoint,
	)
	go service.Start(suite.ctx, make(chan struct{}, 1))
	// let reader run for some time
	time.Sleep(10 * time.Second)
}

func (suite *EspressoReaderTestSuite) TestInputs() {
	// input0
	input, err := suite.database.GetInput(suite.ctx, application.IApplicationAddress, 0)
	suite.Nil(err)
	suite.Equal(uint64(755175), input.EpochIndex)
	suite.Equal(uint64(0), input.Index)
	suite.Equal(uint64(7551758), input.BlockNumber)
	suite.Equal(EvmAdvanceEncode(7551758, 1737607776, "0x639a5df764d365b4c75dc038b068d9aee72bed6465cafee577b29e3dc4b0a8f8", 0, "bb01"), input.RawData)
	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0xcbd6ab30a829bd312776f9411f8e11b5997b757fd64f202336200a65ebab2ff8"), input.TransactionReference)

	// input1
	input, err = suite.database.GetInput(suite.ctx, application.IApplicationAddress, 1)
	suite.Nil(err)
	suite.Equal(uint64(755175), input.EpochIndex)
	suite.Equal(uint64(1), input.Index)
	suite.Equal(uint64(7551758), input.BlockNumber)
	suite.Equal(EvmAdvanceEncode(7551758, 1737607776, "0x639a5df764d365b4c75dc038b068d9aee72bed6465cafee577b29e3dc4b0a8f8", 1, "bb02"), input.RawData)
	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0xe836569fc44b6ac012748f69dc13e40ab9f65020a453ee2d27605cbb7aa44f97"), input.TransactionReference)

	// input2
	input, err = suite.database.GetInput(suite.ctx, application.IApplicationAddress, 2)
	suite.Nil(err)
	suite.Equal(uint64(755182), input.EpochIndex)
	suite.Equal(uint64(2), input.Index)
	suite.Equal(uint64(7551823), input.BlockNumber)
	suite.Equal(EvmAdvanceEncode(7551823, 1737608604, "0x6e60b652c4567d42692930f668ac87ea5ace1c52ec9d36a24705c67f50748acc", 2, "aa03"), input.RawData)
	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0x00"), input.TransactionReference)

	// input3
	input, err = suite.database.GetInput(suite.ctx, application.IApplicationAddress, 3)
	suite.Nil(err)
	suite.Equal(uint64(755182), input.EpochIndex)
	suite.Equal(uint64(3), input.Index)
	suite.Equal(uint64(7551827), input.BlockNumber)
	suite.Equal(EvmAdvanceEncode(7551827, 1737608664, "0x3e5fd5b712082c3f1a1482f88ada563c61867b407ad00498dfe0c482d17c255a", 3, "aa04"), input.RawData)
	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0x01"), input.TransactionReference)

	// input4
	input, err = suite.database.GetInput(suite.ctx, application.IApplicationAddress, 4)
	suite.Nil(err)
	suite.Equal(uint64(755184), input.EpochIndex)
	suite.Equal(uint64(4), input.Index)
	suite.Equal(uint64(7551847), input.BlockNumber)
	suite.Equal(EvmAdvanceEncode(7551847, 1737608916, "0x93e406ac7cbc987af170daf5a8c37f8c410dd06119a7efcbb74e17387f301ed2", 4, "bb05"), input.RawData)
	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0xa5d16b4799a749218846823304cf59d4bd4f1c6f609e05b8332d0082871b3e37"), input.TransactionReference)

	// input5
	input, err = suite.database.GetInput(suite.ctx, application.IApplicationAddress, 5)
	suite.Nil(err)
	suite.Equal(uint64(755184), input.EpochIndex)
	suite.Equal(uint64(5), input.Index)
	suite.Equal(uint64(7551847), input.BlockNumber)
	suite.Equal(EvmAdvanceEncode(7551847, 1737608916, "0x93e406ac7cbc987af170daf5a8c37f8c410dd06119a7efcbb74e17387f301ed2", 5, "bb06"), input.RawData)
	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0x184d3fb1de731efd9e5643527b946f69b8ee1d6dd0873270ef82dd18ef2f31b5"), input.TransactionReference)
}

func EvmAdvanceEncode(l1FinalizedLatestHeight uint64, l1FinalizedTimestamp uint64, prevRandaoHexStringWith0x string, index uint64, payloadHexStringNo0x string) []byte {
	abiObject, _ := abi.JSON(strings.NewReader(
		`[{
			"type" : "function",
			"name" : "EvmAdvance",
			"inputs" : [
				{ "type" : "uint256" },
				{ "type" : "address" },
				{ "type" : "address" },
				{ "type" : "uint256" },
				{ "type" : "uint256" },
				{ "type" : "uint256" },
				{ "type" : "uint256" },
				{ "type" : "bytes"   }
			]
		}, {
			"type" : "function",
			"name" : "Voucher",
			"inputs" : [
				{ "type" : "address" },
				{ "type" : "uint256" },
				{ "type" : "bytes"   }
			]
		}, {
			"type" : "function",
			"name" : "Notice",
			"inputs" : [
				{ "type" : "bytes"   }
			]
		}]`,
	))
	prevRandao, _ := hexutil.DecodeBig(prevRandaoHexStringWith0x)
	payloadBytes, _ := hex.DecodeString(payloadHexStringNo0x)
	payloadAbi, _ := abiObject.Pack("EvmAdvance",
		(&big.Int{}).SetInt64(int64(11155111)),
		common.HexToAddress(application.IApplicationAddress),
		common.HexToAddress("0x590F92fEa8df163fFF2d7Df266364De7CE8F9E16"),
		(&big.Int{}).SetUint64(l1FinalizedLatestHeight),
		(&big.Int{}).SetUint64(l1FinalizedTimestamp),
		prevRandao,
		(&big.Int{}).SetUint64(index),
		payloadBytes)
	return payloadAbi
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(EspressoReaderTestSuite))
}
