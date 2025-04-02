// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package espressoreader

import (
	"context"
	"encoding/hex"
	"log/slog"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/cartesi/rollups-espresso-reader/internal/config"
	"github.com/cartesi/rollups-espresso-reader/internal/evmreader"
	"github.com/cartesi/rollups-espresso-reader/internal/model"
	"github.com/cartesi/rollups-espresso-reader/internal/repository"
	"github.com/cartesi/rollups-espresso-reader/internal/repository/factory"
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
	appAddress := "0xA0a3baf0ABB307C7c944c2C1B60A50bBE34BD637"
	consensusAddress := "0xF957365e6a65687a9dcE4fD0244667FD5C09E783"

	templatePath := "applications/echo-dapp/"
	templateHash := "0x1a456963bec81fdb3eddeddb2881281737b1136e744f856d6aa43a8dcf920cd5"
	suite.application = model.Application{
		Name:                 "test-dapp",
		IApplicationAddress:  common.HexToAddress(appAddress),
		IConsensusAddress:    common.HexToAddress(consensusAddress),
		TemplateURI:          templatePath,
		TemplateHash:         common.HexToHash(templateHash),
		State:                model.ApplicationState_Enabled,
		LastInputCheckBlock:  0,
		LastOutputCheckBlock: 0,
		EpochLength:          10,
	}
	suite.senderAddress = common.HexToAddress("0x590F92fEa8df163fFF2d7Df266364De7CE8F9E16").String()

	suite.ctx = context.Background()
	suite.c = config.FromEnv()
	suite.database, _ = factory.NewRepositoryFromConnectionString(suite.ctx, suite.c.PostgresEndpoint.Value)
	_, err := suite.database.CreateApplication(suite.ctx, &suite.application)
	if err != nil {
		slog.Error("create application", "error", err)
	}
	suite.Nil(err)

	config, err := repository.LoadNodeConfig[evmreader.PersistentConfig](suite.ctx, suite.database, evmreader.EvmReaderConfigKey)
	if err != nil {
		slog.Error("db config", "error", err)
	}

	service := NewEspressoReaderService(
		suite.c.BlockchainHttpEndpoint.Value,
		suite.c.BlockchainWsEndpoint.Value,
		suite.database,
		suite.c.EspressoBaseUrl,
		suite.c.EspressoStartingBlock,
		suite.c.EspressoNamespace,
		config.Value.ChainID,
		suite.c.EspressoServiceEndpoint,
		suite.c.MaxRetries,
		suite.c.MaxDelay,
	)
	go service.Start(suite.ctx, make(chan struct{}, 1))
	// let reader run for some time
	time.Sleep(10 * time.Second)
}

func (suite *EspressoReaderTestSuite) TestInputs() {
	// input0
	input, err := suite.database.GetInput(suite.ctx, suite.application.IApplicationAddress.Hex(), 0)
	suite.Nil(err)
	suite.Equal(uint64(756613), input.EpochIndex)
	suite.Equal(uint64(0), input.Index)
	suite.Equal(uint64(7566135), input.BlockNumber)
	suite.Equal(EvmAdvanceEncode(suite.application.IApplicationAddress.Hex(), 7566135, 1737786720, "0xf4b4be2ff1373371ff705336b0eff0f70f117bcfe61fb04470e613a1b6e710f6", 0, "bb01"), input.RawData)
	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0x04cb2532275733f68678223314a231d44ad407832f7fce0902ec95328c870fbd"), input.TransactionReference)

	// input1
	input, err = suite.database.GetInput(suite.ctx, suite.application.IApplicationAddress.Hex(), 1)
	suite.Nil(err)
	suite.Equal(uint64(756613), input.EpochIndex)
	suite.Equal(uint64(1), input.Index)
	suite.Equal(uint64(7566135), input.BlockNumber)
	suite.Equal(EvmAdvanceEncode(suite.application.IApplicationAddress.Hex(), 7566135, 1737786720, "0xf4b4be2ff1373371ff705336b0eff0f70f117bcfe61fb04470e613a1b6e710f6", 1, "bb02"), input.RawData)
	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0x02b0d850e04f4008109d7f65c0080304e01c97defc1998084d4c319997d2d9fb"), input.TransactionReference)

	// input2
	input, err = suite.database.GetInput(suite.ctx, suite.application.IApplicationAddress.Hex(), 2)
	suite.Nil(err)
	suite.Equal(uint64(756621), input.EpochIndex)
	suite.Equal(uint64(2), input.Index)
	suite.Equal(uint64(7566214), input.BlockNumber)
	suite.Equal(EvmAdvanceEncode(suite.application.IApplicationAddress.Hex(), 7566214, 1737787692, "0x9171073480935a9d1bafa8c6380716a352bc1471c18dc091bcdfdec6d0c24072", 2, "aa03"), input.RawData)
	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0x00"), input.TransactionReference)

	// input3
	input, err = suite.database.GetInput(suite.ctx, suite.application.IApplicationAddress.Hex(), 3)
	suite.Nil(err)
	suite.Equal(uint64(756621), input.EpochIndex)
	suite.Equal(uint64(3), input.Index)
	suite.Equal(uint64(7566216), input.BlockNumber)
	suite.Equal(EvmAdvanceEncode(suite.application.IApplicationAddress.Hex(), 7566216, 1737787728, "0x63ed554373a61df9fbcbc91ea12e1d0b6a4e7e108837c88204194c244c8ff8fe", 3, "aa04"), input.RawData)
	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0x01"), input.TransactionReference)

	// input4
	input, err = suite.database.GetInput(suite.ctx, suite.application.IApplicationAddress.Hex(), 4)
	suite.Nil(err)
	suite.Equal(uint64(756622), input.EpochIndex)
	suite.Equal(uint64(4), input.Index)
	suite.Equal(uint64(7566228), input.BlockNumber)
	suite.Equal(EvmAdvanceEncode(suite.application.IApplicationAddress.Hex(), 7566228, 1737787872, "0x2bbd0c2b7dbbced136cfef527cde4c0fdd6d954d31341c5ee485e64cc15e6d83", 4, "bb05"), input.RawData)
	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0x5ded89e031600ae61b35ac3b383758651e707442f2acb26606e42ae8c88bcf3e"), input.TransactionReference)

	// input5
	input, err = suite.database.GetInput(suite.ctx, suite.application.IApplicationAddress.Hex(), 5)
	suite.Nil(err)
	suite.Equal(uint64(756622), input.EpochIndex)
	suite.Equal(uint64(5), input.Index)
	suite.Equal(uint64(7566228), input.BlockNumber)
	suite.Equal(EvmAdvanceEncode(suite.application.IApplicationAddress.Hex(), 7566228, 1737787872, "0x2bbd0c2b7dbbced136cfef527cde4c0fdd6d954d31341c5ee485e64cc15e6d83", 5, "bb06"), input.RawData)
	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0x73a160dde917feab97cacdc1e99d08c455c1f97621daaf75cabda53fabe31ee0"), input.TransactionReference)
}

func (suite *EspressoReaderTestSuite) TestEspressoNonce() {
	nonce, err := suite.database.GetEspressoNonce(suite.ctx, suite.senderAddress, suite.application.IApplicationAddress.Hex())
	suite.Nil(err)
	suite.Equal(uint64(4), nonce)
}

func (suite *EspressoReaderTestSuite) TestApplication() {
	appInDB, err := suite.database.GetApplication(suite.ctx, suite.application.IApplicationAddress.Hex())
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
	epoch0, err := suite.database.GetEpoch(suite.ctx, suite.application.IApplicationAddress.Hex(), 756613)
	suite.Nil(err)
	suite.Equal(int64(1), epoch0.ApplicationID)
	suite.Equal(uint64(756613), epoch0.Index)
	suite.Equal(uint64(7566130), epoch0.FirstBlock)
	suite.Equal(uint64(7566139), epoch0.LastBlock)
	// suite.Equal(, epoch0.ClaimHash)
	// suite.Equal(, epoch0.ClaimTransactionHash)
	// suite.Equal(, epoch0.Status)
	suite.Equal(uint64(0), epoch0.VirtualIndex)

	epoch1, err := suite.database.GetEpoch(suite.ctx, suite.application.IApplicationAddress.Hex(), 756621)
	suite.Nil(err)
	suite.Equal(int64(1), epoch1.ApplicationID)
	suite.Equal(uint64(756621), epoch1.Index)
	suite.Equal(uint64(7566210), epoch1.FirstBlock)
	suite.Equal(uint64(7566219), epoch1.LastBlock)
	// suite.Equal(, epoch1.ClaimHash)
	// suite.Equal(, epoch1.ClaimTransactionHash)
	// suite.Equal(, epoch1.Status)
	suite.Equal(uint64(1), epoch1.VirtualIndex)

	epoch2, err := suite.database.GetEpoch(suite.ctx, suite.application.IApplicationAddress.Hex(), 756622)
	suite.Nil(err)
	suite.Equal(int64(1), epoch2.ApplicationID)
	suite.Equal(uint64(756622), epoch2.Index)
	suite.Equal(uint64(7566220), epoch2.FirstBlock)
	suite.Equal(uint64(7566229), epoch2.LastBlock)
	// suite.Equal(, epoch2.ClaimHash)
	// suite.Equal(, epoch2.ClaimTransactionHash)
	// suite.Equal(, epoch2.Status)
	suite.Equal(uint64(2), epoch2.VirtualIndex)
}

func (suite *EspressoReaderTestSuite) TestNodeConfig() {
	nodeConfig, err := repository.LoadNodeConfig[evmreader.PersistentConfig](suite.ctx, suite.database, evmreader.EvmReaderConfigKey)
	suite.Nil(err)
	suite.Equal(model.DefaultBlock_Finalized, nodeConfig.Value.DefaultBlock)
	suite.Equal(uint64(11155111), nodeConfig.Value.ChainID)
}

func (suite *EspressoReaderTestSuite) TestOutput() {
	output0, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress.Hex(), 0)
	suite.Nil(err)
	suite.Equal(int64(1), output0.InputEpochApplicationID)
	suite.Equal(uint64(0), output0.InputIndex)
	suite.Equal(uint64(0), output0.Index)
	suite.Equal(VoucherEncode(suite.senderAddress, "bb01"), output0.RawData)
	// suite.Equal(, output0.Hash)
	// suite.Equal(, output0.OutputHashesSiblings)
	// suite.Equal(, output0.ExecutionTransactionHash)
	output1, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress.Hex(), 1)
	suite.Nil(err)
	suite.Equal(int64(1), output1.InputEpochApplicationID)
	suite.Equal(uint64(0), output1.InputIndex)
	suite.Equal(uint64(1), output1.Index)
	suite.Equal(NoticeEncode("bb01"), output1.RawData)
	// suite.Equal(, output1.Hash)
	// suite.Equal(, output1.OutputHashesSiblings)
	// suite.Equal(, output1.ExecutionTransactionHash)

	output2, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress.Hex(), 2)
	suite.Nil(err)
	suite.Equal(int64(1), output2.InputEpochApplicationID)
	suite.Equal(uint64(1), output2.InputIndex)
	suite.Equal(uint64(2), output2.Index)
	suite.Equal(VoucherEncode(suite.senderAddress, "bb02"), output2.RawData)
	// suite.Equal(, output2.Hash)
	// suite.Equal(, output2.OutputHashesSiblings)
	// suite.Equal(, output2.ExecutionTransactionHash)
	output3, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress.Hex(), 3)
	suite.Nil(err)
	suite.Equal(int64(1), output3.InputEpochApplicationID)
	suite.Equal(uint64(1), output3.InputIndex)
	suite.Equal(uint64(3), output3.Index)
	suite.Equal(NoticeEncode("bb02"), output3.RawData)
	// suite.Equal(, output3.Hash)
	// suite.Equal(, output3.OutputHashesSiblings)
	// suite.Equal(, output3.ExecutionTransactionHash)

	output4, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress.Hex(), 4)
	suite.Nil(err)
	suite.Equal(int64(1), output4.InputEpochApplicationID)
	suite.Equal(uint64(2), output4.InputIndex)
	suite.Equal(uint64(4), output4.Index)
	suite.Equal(VoucherEncode(suite.senderAddress, "aa03"), output4.RawData)
	// suite.Equal(, output4.Hash)
	// suite.Equal(, output4.OutputHashesSiblings)
	// suite.Equal(, output4.ExecutionTransactionHash)
	output5, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress.Hex(), 5)
	suite.Nil(err)
	suite.Equal(int64(1), output5.InputEpochApplicationID)
	suite.Equal(uint64(2), output5.InputIndex)
	suite.Equal(uint64(5), output5.Index)
	suite.Equal(NoticeEncode("aa03"), output5.RawData)
	// suite.Equal(, output5.Hash)
	// suite.Equal(, output5.OutputHashesSiblings)
	// suite.Equal(, output5.ExecutionTransactionHash)

	output6, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress.Hex(), 6)
	suite.Nil(err)
	suite.Equal(int64(1), output6.InputEpochApplicationID)
	suite.Equal(uint64(3), output6.InputIndex)
	suite.Equal(uint64(6), output6.Index)
	suite.Equal(VoucherEncode(suite.senderAddress, "aa04"), output6.RawData)
	// suite.Equal(, output6.Hash)
	// suite.Equal(, output6.OutputHashesSiblings)
	// suite.Equal(, output6.ExecutionTransactionHash)
	output7, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress.Hex(), 7)
	suite.Nil(err)
	suite.Equal(int64(1), output7.InputEpochApplicationID)
	suite.Equal(uint64(3), output7.InputIndex)
	suite.Equal(uint64(7), output7.Index)
	suite.Equal(NoticeEncode("aa04"), output7.RawData)
	// suite.Equal(, output7.Hash)
	// suite.Equal(, output7.OutputHashesSiblings)
	// suite.Equal(, output7.ExecutionTransactionHash)

	output8, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress.Hex(), 8)
	suite.Nil(err)
	suite.Equal(int64(1), output8.InputEpochApplicationID)
	suite.Equal(uint64(4), output8.InputIndex)
	suite.Equal(uint64(8), output8.Index)
	suite.Equal(VoucherEncode(suite.senderAddress, "bb05"), output8.RawData)
	// suite.Equal(, output8.Hash)
	// suite.Equal(, output8.OutputHashesSiblings)
	// suite.Equal(, output8.ExecutionTransactionHash)
	output9, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress.Hex(), 9)
	suite.Nil(err)
	suite.Equal(int64(1), output9.InputEpochApplicationID)
	suite.Equal(uint64(4), output9.InputIndex)
	suite.Equal(uint64(9), output9.Index)
	suite.Equal(NoticeEncode("bb05"), output9.RawData)
	// suite.Equal(, output9.Hash)
	// suite.Equal(, output9.OutputHashesSiblings)
	// suite.Equal(, output9.ExecutionTransactionHash)

	output10, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress.Hex(), 10)
	suite.Nil(err)
	suite.Equal(int64(1), output10.InputEpochApplicationID)
	suite.Equal(uint64(5), output10.InputIndex)
	suite.Equal(uint64(10), output10.Index)
	suite.Equal(VoucherEncode(suite.senderAddress, "bb06"), output10.RawData)
	// suite.Equal(, output10.Hash)
	// suite.Equal(, output10.OutputHashesSiblings)
	// suite.Equal(, output10.ExecutionTransactionHash)
	output11, err := suite.database.GetOutput(suite.ctx, suite.application.IApplicationAddress.Hex(), 11)
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
	report0, err := suite.database.GetReport(suite.ctx, suite.application.IApplicationAddress.Hex(), 0)
	suite.Nil(err)
	suite.Equal(int64(1), report0.InputEpochApplicationID)
	suite.Equal(uint64(0), report0.InputIndex)
	suite.Equal(uint64(0), report0.Index)
	expectedRawData, _ := hex.DecodeString("bb01")
	suite.Equal(expectedRawData, report0.RawData)

	report1, err := suite.database.GetReport(suite.ctx, suite.application.IApplicationAddress.Hex(), 1)
	suite.Nil(err)
	suite.Equal(int64(1), report1.InputEpochApplicationID)
	suite.Equal(uint64(1), report1.InputIndex)
	suite.Equal(uint64(1), report1.Index)
	expectedRawData, _ = hex.DecodeString("bb02")
	suite.Equal(expectedRawData, report1.RawData)

	report2, err := suite.database.GetReport(suite.ctx, suite.application.IApplicationAddress.Hex(), 2)
	suite.Nil(err)
	suite.Equal(int64(1), report2.InputEpochApplicationID)
	suite.Equal(uint64(2), report2.InputIndex)
	suite.Equal(uint64(2), report2.Index)
	expectedRawData, _ = hex.DecodeString("aa03")
	suite.Equal(expectedRawData, report2.RawData)

	report3, err := suite.database.GetReport(suite.ctx, suite.application.IApplicationAddress.Hex(), 3)
	suite.Nil(err)
	suite.Equal(int64(1), report3.InputEpochApplicationID)
	suite.Equal(uint64(3), report3.InputIndex)
	suite.Equal(uint64(3), report3.Index)
	expectedRawData, _ = hex.DecodeString("aa04")
	suite.Equal(expectedRawData, report3.RawData)

	report4, err := suite.database.GetReport(suite.ctx, suite.application.IApplicationAddress.Hex(), 4)
	suite.Nil(err)
	suite.Equal(int64(1), report4.InputEpochApplicationID)
	suite.Equal(uint64(4), report4.InputIndex)
	suite.Equal(uint64(4), report4.Index)
	expectedRawData, _ = hex.DecodeString("bb05")
	suite.Equal(expectedRawData, report4.RawData)

	report5, err := suite.database.GetReport(suite.ctx, suite.application.IApplicationAddress.Hex(), 5)
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
