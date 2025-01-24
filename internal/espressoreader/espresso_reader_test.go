// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package espressoreader

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
	ctx           context.Context
	c             config.NodeConfig
	database      repository.Repository
	application   model.Application
	senderAddress string
}

func (suite *EspressoReaderTestSuite) SetupSuite() {
	appAddress := "0xdDedD051B0EeeF07A5099eE572Af85Cb816B3fdE"
	consensusAddress := "0xb66b8260646df96E20E3745489BAF0108E842B8e"

	templatePath := "applications/echo-dapp/"
	templateHash := "0x1a456963bec81fdb3eddeddb2881281737b1136e744f856d6aa43a8dcf920cd5"
	suite.application = model.Application{
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
	suite.senderAddress = strings.ToLower(common.HexToAddress("0x590F92fEa8df163fFF2d7Df266364De7CE8F9E16").String())

	suite.ctx = context.Background()
	suite.c = config.FromEnv()
	suite.database, _ = factory.NewRepositoryFromConnectionString(suite.ctx, suite.c.PostgresEndpoint.Value)
	_, err := suite.database.CreateApplication(suite.ctx, &suite.application)
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
	input, err := suite.database.GetInput(suite.ctx, suite.application.IApplicationAddress, 0)
	suite.Nil(err)
	suite.Equal(uint64(755175), input.EpochIndex)
	suite.Equal(uint64(0), input.Index)
	suite.Equal(uint64(7551758), input.BlockNumber)
	suite.Equal(EvmAdvanceEncode(suite.application.IApplicationAddress, 7551758, 1737607776, "0x639a5df764d365b4c75dc038b068d9aee72bed6465cafee577b29e3dc4b0a8f8", 0, "bb01"), input.RawData)
	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0xcbd6ab30a829bd312776f9411f8e11b5997b757fd64f202336200a65ebab2ff8"), input.TransactionReference)

	// input1
	input, err = suite.database.GetInput(suite.ctx, suite.application.IApplicationAddress, 1)
	suite.Nil(err)
	suite.Equal(uint64(755175), input.EpochIndex)
	suite.Equal(uint64(1), input.Index)
	suite.Equal(uint64(7551758), input.BlockNumber)
	suite.Equal(EvmAdvanceEncode(suite.application.IApplicationAddress, 7551758, 1737607776, "0x639a5df764d365b4c75dc038b068d9aee72bed6465cafee577b29e3dc4b0a8f8", 1, "bb02"), input.RawData)
	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0xe836569fc44b6ac012748f69dc13e40ab9f65020a453ee2d27605cbb7aa44f97"), input.TransactionReference)

	// input2
	input, err = suite.database.GetInput(suite.ctx, suite.application.IApplicationAddress, 2)
	suite.Nil(err)
	suite.Equal(uint64(755182), input.EpochIndex)
	suite.Equal(uint64(2), input.Index)
	suite.Equal(uint64(7551823), input.BlockNumber)
	suite.Equal(EvmAdvanceEncode(suite.application.IApplicationAddress, 7551823, 1737608604, "0x6e60b652c4567d42692930f668ac87ea5ace1c52ec9d36a24705c67f50748acc", 2, "aa03"), input.RawData)
	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0x00"), input.TransactionReference)

	// input3
	input, err = suite.database.GetInput(suite.ctx, suite.application.IApplicationAddress, 3)
	suite.Nil(err)
	suite.Equal(uint64(755182), input.EpochIndex)
	suite.Equal(uint64(3), input.Index)
	suite.Equal(uint64(7551827), input.BlockNumber)
	suite.Equal(EvmAdvanceEncode(suite.application.IApplicationAddress, 7551827, 1737608664, "0x3e5fd5b712082c3f1a1482f88ada563c61867b407ad00498dfe0c482d17c255a", 3, "aa04"), input.RawData)
	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0x01"), input.TransactionReference)

	// input4
	input, err = suite.database.GetInput(suite.ctx, suite.application.IApplicationAddress, 4)
	suite.Nil(err)
	suite.Equal(uint64(755184), input.EpochIndex)
	suite.Equal(uint64(4), input.Index)
	suite.Equal(uint64(7551847), input.BlockNumber)
	suite.Equal(EvmAdvanceEncode(suite.application.IApplicationAddress, 7551847, 1737608916, "0x93e406ac7cbc987af170daf5a8c37f8c410dd06119a7efcbb74e17387f301ed2", 4, "bb05"), input.RawData)
	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0xa5d16b4799a749218846823304cf59d4bd4f1c6f609e05b8332d0082871b3e37"), input.TransactionReference)

	// input5
	input, err = suite.database.GetInput(suite.ctx, suite.application.IApplicationAddress, 5)
	suite.Nil(err)
	suite.Equal(uint64(755184), input.EpochIndex)
	suite.Equal(uint64(5), input.Index)
	suite.Equal(uint64(7551847), input.BlockNumber)
	suite.Equal(EvmAdvanceEncode(suite.application.IApplicationAddress, 7551847, 1737608916, "0x93e406ac7cbc987af170daf5a8c37f8c410dd06119a7efcbb74e17387f301ed2", 5, "bb06"), input.RawData)
	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0x184d3fb1de731efd9e5643527b946f69b8ee1d6dd0873270ef82dd18ef2f31b5"), input.TransactionReference)
}

func (suite *EspressoReaderTestSuite) TestEspressoNonce() {
	nonce, err := suite.database.GetEspressoNonce(suite.ctx, suite.senderAddress, suite.application.IApplicationAddress)
	suite.Nil(err)
	suite.Equal(uint64(4), nonce)
}

func (suite *EspressoReaderTestSuite) TestApplication() {
	appInDB, err := suite.database.GetApplication(suite.ctx, suite.application.IApplicationAddress)
	suite.Nil(err)
	suite.Equal(int64(1), appInDB.ID)
	suite.Equal("test-dapp", appInDB.Name)
	suite.Equal(suite.application.IApplicationAddress, appInDB.IApplicationAddress)
	suite.Equal(suite.application.IConsensusAddress, appInDB.IConsensusAddress)
	suite.Equal(suite.application.TemplateHash, appInDB.TemplateHash)
	suite.Equal(suite.application.TemplateURI, appInDB.TemplateURI)
	suite.Equal(model.ApplicationState_Enabled, appInDB.State)
	suite.Equal(uint64(6), appInDB.ProcessedInputs)
}

func (suite *EspressoReaderTestSuite) TestEpoch() {
	epoch0, err := suite.database.GetEpoch(suite.ctx, suite.application.IApplicationAddress, 755175)
	suite.Nil(err)
	suite.Equal(int64(1), epoch0.ApplicationID)
	suite.Equal(uint64(755175), epoch0.Index)
	suite.Equal(uint64(7551750), epoch0.FirstBlock)
	suite.Equal(uint64(7551759), epoch0.LastBlock)
	// suite.Equal(, epoch0.ClaimHash)
	// suite.Equal(, epoch0.ClaimTransactionHash)
	// suite.Equal(, epoch0.Status)
	suite.Equal(uint64(0), epoch0.VirtualIndex)

	epoch1, err := suite.database.GetEpoch(suite.ctx, suite.application.IApplicationAddress, 755182)
	suite.Nil(err)
	suite.Equal(int64(1), epoch1.ApplicationID)
	suite.Equal(uint64(755182), epoch1.Index)
	suite.Equal(uint64(7551820), epoch1.FirstBlock)
	suite.Equal(uint64(7551829), epoch1.LastBlock)
	// suite.Equal(, epoch1.ClaimHash)
	// suite.Equal(, epoch1.ClaimTransactionHash)
	// suite.Equal(, epoch1.Status)
	suite.Equal(uint64(1), epoch1.VirtualIndex)

	epoch2, err := suite.database.GetEpoch(suite.ctx, suite.application.IApplicationAddress, 755184)
	suite.Nil(err)
	suite.Equal(int64(1), epoch2.ApplicationID)
	suite.Equal(uint64(755184), epoch2.Index)
	suite.Equal(uint64(7551840), epoch2.FirstBlock)
	suite.Equal(uint64(7551849), epoch2.LastBlock)
	// suite.Equal(, epoch2.ClaimHash)
	// suite.Equal(, epoch2.ClaimTransactionHash)
	// suite.Equal(, epoch2.Status)
	suite.Equal(uint64(2), epoch2.VirtualIndex)
}

func (suite *EspressoReaderTestSuite) TestNodeConfig() {
	nodeConfig, err := suite.database.GetNodeConfig(suite.ctx)
	suite.Nil(err)
	suite.Equal(model.DefaultBlock_Finalized, nodeConfig.DefaultBlock)
	suite.Equal(strings.ToLower(common.HexToAddress("0x593E5BCf894D6829Dd26D0810DA7F064406aebB6").String()), nodeConfig.InputBoxAddress)
	suite.Equal(uint64(11155111), nodeConfig.ChainID)
}

func (suite *EspressoReaderTestSuite) TestOutput() {
	output0, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress, 0)
	suite.Nil(err)
	suite.Equal(int64(1), output0.InputEpochApplicationID)
	suite.Equal(uint64(0), output0.InputIndex)
	suite.Equal(uint64(0), output0.Index)
	suite.Equal(VoucherEncode(suite.senderAddress, "bb01"), output0.RawData)
	// suite.Equal(, output0.Hash)
	// suite.Equal(, output0.OutputHashesSiblings)
	// suite.Equal(, output0.ExecutionTransactionHash)
	output1, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress, 1)
	suite.Nil(err)
	suite.Equal(int64(1), output1.InputEpochApplicationID)
	suite.Equal(uint64(0), output1.InputIndex)
	suite.Equal(uint64(1), output1.Index)
	suite.Equal(NoticeEncode("bb01"), output1.RawData)
	// suite.Equal(, output1.Hash)
	// suite.Equal(, output1.OutputHashesSiblings)
	// suite.Equal(, output1.ExecutionTransactionHash)

	output2, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress, 2)
	suite.Nil(err)
	suite.Equal(int64(1), output2.InputEpochApplicationID)
	suite.Equal(uint64(1), output2.InputIndex)
	suite.Equal(uint64(2), output2.Index)
	suite.Equal(VoucherEncode(suite.senderAddress, "bb02"), output2.RawData)
	// suite.Equal(, output2.Hash)
	// suite.Equal(, output2.OutputHashesSiblings)
	// suite.Equal(, output2.ExecutionTransactionHash)
	output3, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress, 3)
	suite.Nil(err)
	suite.Equal(int64(1), output3.InputEpochApplicationID)
	suite.Equal(uint64(1), output3.InputIndex)
	suite.Equal(uint64(3), output3.Index)
	suite.Equal(NoticeEncode("bb02"), output3.RawData)
	// suite.Equal(, output3.Hash)
	// suite.Equal(, output3.OutputHashesSiblings)
	// suite.Equal(, output3.ExecutionTransactionHash)

	output4, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress, 4)
	suite.Nil(err)
	suite.Equal(int64(1), output4.InputEpochApplicationID)
	suite.Equal(uint64(2), output4.InputIndex)
	suite.Equal(uint64(4), output4.Index)
	suite.Equal(VoucherEncode(suite.senderAddress, "aa03"), output4.RawData)
	// suite.Equal(, output4.Hash)
	// suite.Equal(, output4.OutputHashesSiblings)
	// suite.Equal(, output4.ExecutionTransactionHash)
	output5, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress, 5)
	suite.Nil(err)
	suite.Equal(int64(1), output5.InputEpochApplicationID)
	suite.Equal(uint64(2), output5.InputIndex)
	suite.Equal(uint64(5), output5.Index)
	suite.Equal(NoticeEncode("aa03"), output5.RawData)
	// suite.Equal(, output5.Hash)
	// suite.Equal(, output5.OutputHashesSiblings)
	// suite.Equal(, output5.ExecutionTransactionHash)

	output6, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress, 6)
	suite.Nil(err)
	suite.Equal(int64(1), output6.InputEpochApplicationID)
	suite.Equal(uint64(3), output6.InputIndex)
	suite.Equal(uint64(6), output6.Index)
	suite.Equal(VoucherEncode(suite.senderAddress, "aa04"), output6.RawData)
	// suite.Equal(, output6.Hash)
	// suite.Equal(, output6.OutputHashesSiblings)
	// suite.Equal(, output6.ExecutionTransactionHash)
	output7, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress, 7)
	suite.Nil(err)
	suite.Equal(int64(1), output7.InputEpochApplicationID)
	suite.Equal(uint64(3), output7.InputIndex)
	suite.Equal(uint64(7), output7.Index)
	suite.Equal(NoticeEncode("aa04"), output7.RawData)
	// suite.Equal(, output7.Hash)
	// suite.Equal(, output7.OutputHashesSiblings)
	// suite.Equal(, output7.ExecutionTransactionHash)

	output8, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress, 8)
	suite.Nil(err)
	suite.Equal(int64(1), output8.InputEpochApplicationID)
	suite.Equal(uint64(4), output8.InputIndex)
	suite.Equal(uint64(8), output8.Index)
	suite.Equal(VoucherEncode(suite.senderAddress, "bb05"), output8.RawData)
	// suite.Equal(, output8.Hash)
	// suite.Equal(, output8.OutputHashesSiblings)
	// suite.Equal(, output8.ExecutionTransactionHash)
	output9, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress, 9)
	suite.Nil(err)
	suite.Equal(int64(1), output9.InputEpochApplicationID)
	suite.Equal(uint64(4), output9.InputIndex)
	suite.Equal(uint64(9), output9.Index)
	suite.Equal(NoticeEncode("bb05"), output9.RawData)
	// suite.Equal(, output9.Hash)
	// suite.Equal(, output9.OutputHashesSiblings)
	// suite.Equal(, output9.ExecutionTransactionHash)

	output10, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress, 10)
	suite.Nil(err)
	suite.Equal(int64(1), output10.InputEpochApplicationID)
	suite.Equal(uint64(5), output10.InputIndex)
	suite.Equal(uint64(10), output10.Index)
	suite.Equal(VoucherEncode(suite.senderAddress, "bb06"), output10.RawData)
	// suite.Equal(, output10.Hash)
	// suite.Equal(, output10.OutputHashesSiblings)
	// suite.Equal(, output10.ExecutionTransactionHash)
	output11, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress, 11)
	suite.Nil(err)
	suite.Equal(int64(1), output11.InputEpochApplicationID)
	suite.Equal(uint64(5), output11.InputIndex)
	suite.Equal(uint64(11), output11.Index)
	suite.Equal(NoticeEncode("bb06"), output11.RawData)
	// suite.Equal(, output11.Hash)
	// suite.Equal(, output11.OutputHashesSiblings)
	// suite.Equal(, output11.ExecutionTransactionHash)
}

func (suite *EspressoReaderTestSuite) TestReport() {
	report0, err := suite.database.GetReport(suite.ctx, suite.application.IApplicationAddress, 0)
	suite.Nil(err)
	suite.Equal(int64(1), report0.InputEpochApplicationID)
	suite.Equal(uint64(0), report0.InputIndex)
	suite.Equal(uint64(0), report0.Index)
	expectedRawData, _ := hex.DecodeString("bb01")
	suite.Equal(expectedRawData, report0.RawData)

	report1, err := suite.database.GetReport(suite.ctx, suite.application.IApplicationAddress, 1)
	suite.Nil(err)
	suite.Equal(int64(1), report1.InputEpochApplicationID)
	suite.Equal(uint64(1), report1.InputIndex)
	suite.Equal(uint64(1), report1.Index)
	expectedRawData, _ = hex.DecodeString("bb02")
	suite.Equal(expectedRawData, report1.RawData)

	report2, err := suite.database.GetReport(suite.ctx, suite.application.IApplicationAddress, 2)
	suite.Nil(err)
	suite.Equal(int64(1), report2.InputEpochApplicationID)
	suite.Equal(uint64(2), report2.InputIndex)
	suite.Equal(uint64(2), report2.Index)
	expectedRawData, _ = hex.DecodeString("aa03")
	suite.Equal(expectedRawData, report2.RawData)

	report3, err := suite.database.GetReport(suite.ctx, suite.application.IApplicationAddress, 3)
	suite.Nil(err)
	suite.Equal(int64(1), report3.InputEpochApplicationID)
	suite.Equal(uint64(3), report3.InputIndex)
	suite.Equal(uint64(3), report3.Index)
	expectedRawData, _ = hex.DecodeString("aa04")
	suite.Equal(expectedRawData, report3.RawData)

	report4, err := suite.database.GetReport(suite.ctx, suite.application.IApplicationAddress, 4)
	suite.Nil(err)
	suite.Equal(int64(1), report4.InputEpochApplicationID)
	suite.Equal(uint64(4), report4.InputIndex)
	suite.Equal(uint64(4), report4.Index)
	expectedRawData, _ = hex.DecodeString("bb05")
	suite.Equal(expectedRawData, report4.RawData)

	report5, err := suite.database.GetReport(suite.ctx, suite.application.IApplicationAddress, 5)
	suite.Nil(err)
	suite.Equal(int64(1), report5.InputEpochApplicationID)
	suite.Equal(uint64(5), report5.InputIndex)
	suite.Equal(uint64(5), report5.Index)
	expectedRawData, _ = hex.DecodeString("bb06")
	suite.Equal(expectedRawData, report5.RawData)
}

func EvmAdvanceEncode(applicationAddress string, l1FinalizedLatestHeight uint64, l1FinalizedTimestamp uint64, prevRandaoHexStringWith0x string, index uint64, payloadHexStringNo0x string) []byte {
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
		}]`,
	))
	prevRandao, _ := hexutil.DecodeBig(prevRandaoHexStringWith0x)
	payloadBytes, _ := hex.DecodeString(payloadHexStringNo0x)
	abiBytes, _ := abiObject.Pack("EvmAdvance",
		(&big.Int{}).SetInt64(int64(11155111)),
		common.HexToAddress(applicationAddress),
		common.HexToAddress("0x590F92fEa8df163fFF2d7Df266364De7CE8F9E16"),
		(&big.Int{}).SetUint64(l1FinalizedLatestHeight),
		(&big.Int{}).SetUint64(l1FinalizedTimestamp),
		prevRandao,
		(&big.Int{}).SetUint64(index),
		payloadBytes)
	return abiBytes
}

func VoucherEncode(senderAddress string, payloadHexStringNo0x string) []byte {
	abiObject, _ := abi.JSON(strings.NewReader(
		`[{
			"type" : "function",
			"name" : "Voucher",
			"inputs" : [
				{ "type" : "address" },
				{ "type" : "uint256" },
				{ "type" : "bytes"   }
			]
		}]`,
	))
	voucherValue, _ := hexutil.DecodeBig("0xdeadbeef")
	payloadBytes, _ := hex.DecodeString(payloadHexStringNo0x)
	abiBytes, _ := abiObject.Pack("Voucher",
		common.HexToAddress(senderAddress),
		voucherValue,
		payloadBytes)
	return abiBytes
}

func NoticeEncode(payloadHexStringNo0x string) []byte {
	abiObject, _ := abi.JSON(strings.NewReader(
		`[{
			"type" : "function",
			"name" : "Notice",
			"inputs" : [
				{ "type" : "bytes"   }
			]
		}]`,
	))
	payloadBytes, _ := hex.DecodeString(payloadHexStringNo0x)
	abiBytes, _ := abiObject.Pack("Notice", payloadBytes)
	return abiBytes
}

func TestEspressoReaderTestSuite(t *testing.T) {
	suite.Run(t, new(EspressoReaderTestSuite))
}
