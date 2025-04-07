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

	"github.com/EspressoSystems/espresso-sequencer-go/client"
	"github.com/EspressoSystems/espresso-sequencer-go/types"
	"github.com/cartesi/rollups-espresso-reader/internal/config"
	"github.com/cartesi/rollups-espresso-reader/internal/evmreader"
	"github.com/cartesi/rollups-espresso-reader/internal/model"
	"github.com/cartesi/rollups-espresso-reader/internal/repository"
	"github.com/cartesi/rollups-espresso-reader/internal/repository/factory"
	"github.com/cartesi/rollups-espresso-reader/pkg/contracts/inputs"
	"github.com/cartesi/rollups-espresso-reader/pkg/ethutil"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type EspressoReaderTestSuite struct {
	suite.Suite
	ctx           context.Context
	c             config.NodeConfig
	database      repository.Repository
	application   model.Application
	senderAddress string
	chainId       uint64
}

func (suite *EspressoReaderTestSuite) prepareTxs(ctx context.Context, EspressoBaseUrl string, namespace uint64, blockchainHttpEndpoint string) error {
	client := client.NewClient(EspressoBaseUrl)
	var tx types.Transaction
	tx.Namespace = namespace

	tx.Payload = []byte(`{"typedData":{"domain":{"name":"Cartesi","version":"0.1.0","chainId":13370,"verifyingContract":"0x0000000000000000000000000000000000000000"},"types":{"EIP712Domain":[{"name":"name","type":"string"},{"name":"version","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"address"}],"CartesiMessage":[{"name":"app","type":"address"},{"name":"nonce","type":"uint64"},{"name":"max_gas_price","type":"uint128"},{"name":"data","type":"bytes"}]},"primaryType":"CartesiMessage","message":{"app":"0xcbefc92066a262b82840515a132a931bc79f9c1c","nonce":0,"data":"0xbb01","max_gas_price":"10"}},"account":"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266","signature":"0x55894fa93d9120120c875acf73d710d6504d92a80c41254beea3834fc15e486960c85ab0e0631d93e6e8d7d51f5a797d85a9887e15c9fe1e4cd9bd6ae3d83bb81c"}`)
	_, err := client.SubmitTransaction(ctx, tx)
	suite.Require().NoError(err)
	time.Sleep(1 * time.Second)

	tx.Payload = []byte(`{"typedData":{"domain":{"name":"Cartesi","version":"0.1.0","chainId":13370,"verifyingContract":"0x0000000000000000000000000000000000000000"},"types":{"EIP712Domain":[{"name":"name","type":"string"},{"name":"version","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"address"}],"CartesiMessage":[{"name":"app","type":"address"},{"name":"nonce","type":"uint64"},{"name":"max_gas_price","type":"uint128"},{"name":"data","type":"bytes"}]},"primaryType":"CartesiMessage","message":{"app":"0xcbefc92066a262b82840515a132a931bc79f9c1c","nonce":1,"data":"0xbb02","max_gas_price":"10"}},"account":"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266","signature":"0x36e608512962f157b218de25a0a6efb9eecb77acd79b2a95fc6b9f7255f958c0238c148b9dad5fa700179466a6b619a964f679a63aa387c0dd0df9724a1d31891c"}`)
	_, err = client.SubmitTransaction(ctx, tx)
	suite.Require().NoError(err)

	ethClient, err := ethclient.Dial(blockchainHttpEndpoint)
	suite.Require().NoError(err)
	defer ethClient.Close()
	privateKey, err := ethutil.MnemonicToPrivateKey("test test test test test test test test test test test junk", 0)
	suite.Require().NoError(err)
	txOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(int64(suite.chainId)))
	suite.Require().NoError(err)
	L1Payload, _ := hex.DecodeString("aa03")
	_, _, err = ethutil.AddInput(ctx, ethClient, txOpts, suite.application.IInputBoxAddress, suite.application.IApplicationAddress, L1Payload)
	suite.Require().NoError(err)

	L1Payload, _ = hex.DecodeString("aa04")
	_, _, err = ethutil.AddInput(ctx, ethClient, txOpts, suite.application.IInputBoxAddress, suite.application.IApplicationAddress, L1Payload)
	cobra.CheckErr(err)
	time.Sleep(65 * time.Second)

	tx.Payload = []byte(`{"typedData":{"domain":{"name":"Cartesi","version":"0.1.0","chainId":13370,"verifyingContract":"0x0000000000000000000000000000000000000000"},"types":{"EIP712Domain":[{"name":"name","type":"string"},{"name":"version","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"address"}],"CartesiMessage":[{"name":"app","type":"address"},{"name":"nonce","type":"uint64"},{"name":"max_gas_price","type":"uint128"},{"name":"data","type":"bytes"}]},"primaryType":"CartesiMessage","message":{"app":"0xcbefc92066a262b82840515a132a931bc79f9c1c","nonce":2,"data":"0xbb05","max_gas_price":"10"}},"account":"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266","signature":"0xd75d8f9706ef9468724adbddcce48f95a45fef2a9a726132a4fcfb8278d795250fd1a70a9e44e940d578264639db44604af5900cf442ec51b5db815942b23d261b"}`)
	_, err = client.SubmitTransaction(ctx, tx)
	suite.Require().NoError(err)
	time.Sleep(1 * time.Second)

	tx.Payload = []byte(`{"typedData":{"domain":{"name":"Cartesi","version":"0.1.0","chainId":13370,"verifyingContract":"0x0000000000000000000000000000000000000000"},"types":{"EIP712Domain":[{"name":"name","type":"string"},{"name":"version","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"address"}],"CartesiMessage":[{"name":"app","type":"address"},{"name":"nonce","type":"uint64"},{"name":"max_gas_price","type":"uint128"},{"name":"data","type":"bytes"}]},"primaryType":"CartesiMessage","message":{"app":"0xcbefc92066a262b82840515a132a931bc79f9c1c","nonce":3,"data":"0xbb06","max_gas_price":"10"}},"account":"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266","signature":"0x2ed0cbf17aaa41c1b0e0acefb100b3fbbdd52df7046b503d10a09ae8d26cfa36654227639050ac1a041f41ecdbaa24b4c6bddea7a7c2df30c699a2bbe555f15f1b"}`)
	_, err = client.SubmitTransaction(ctx, tx)
	suite.Require().NoError(err)
	time.Sleep(10 * time.Second)

	return nil
}

func (suite *EspressoReaderTestSuite) SetupSuite() {
	appAddress := "0x01e800bbE852aeb27cE65604709134Ea63782c6B" // appAddress := "0xcbefc92066a262b82840515a132a931bc79f9c1c"
	consensusAddress := "0x1d76bdb32803ae72fc5aed528779b3f581f93fed"
	inputboxAddress := "0xB6b39Fb3dD926A9e3FBc7A129540eEbeA3016a6c"
	templatePath := "applications/echo-dapp/"
	templateHash := "0x2fa07a837075faedd5be6215cef05e90848d01fd752e2f41b6039f3317bee84d" // templateHash := "0x1611c2f376328c21520ff4d521eb43d6ca581a5d51eab19e3c354f71ff4bdeae"
	suite.application = model.Application{
		Name:                 "test-dapp",
		IApplicationAddress:  common.HexToAddress(appAddress),
		IConsensusAddress:    common.HexToAddress(consensusAddress),
		IInputBoxAddress:     common.HexToAddress(inputboxAddress),
		TemplateURI:          templatePath,
		TemplateHash:         common.HexToHash(templateHash),
		State:                model.ApplicationState_Enabled,
		LastInputCheckBlock:  0,
		LastOutputCheckBlock: 0,
		EpochLength:          10,
		DataAvailability:     model.DataAvailability_InputBoxAndEspresso,
	}
	suite.senderAddress = common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266").String()

	suite.ctx = context.Background()
	suite.c = config.FromEnv()
	suite.database, _ = factory.NewRepositoryFromConnectionString(suite.ctx, suite.c.PostgresEndpoint.Value)
	// _, err := suite.database.CreateApplication(suite.ctx, &suite.application)
	// if err != nil {
	// 	slog.Error("create application", "error", err)
	// }
	// suite.Nil(err)

	config, err := repository.LoadNodeConfig[evmreader.PersistentConfig](suite.ctx, suite.database, evmreader.EvmReaderConfigKey)
	if err != nil {
		slog.Error("db config", "error", err)
	}
	slog.Debug("SetupSuite", "chainID", config.Value.ChainID)
	suite.chainId = config.Value.ChainID

	_, namespace, err := getEspressoConfig(suite.ctx, common.HexToAddress(appAddress), suite.database, suite.c.BlockchainHttpEndpoint.Value)
	suite.Nil(err)

	err = suite.prepareTxs(suite.ctx, suite.c.EspressoBaseUrl, namespace, suite.c.BlockchainHttpEndpoint.Value)
	suite.Nil(err)

	service := NewEspressoReaderService(
		suite.c.BlockchainHttpEndpoint.Value,
		suite.c.BlockchainWsEndpoint.Value,
		suite.database,
		suite.c.EspressoBaseUrl,
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
	// suite.Equal(uint64(756613), input.EpochIndex)
	suite.Equal(uint64(0), input.Index)
	// suite.Equal(uint64(7566135), input.BlockNumber)
	// suite.Equal(EvmAdvanceEncode(suite.application.IApplicationAddress.Hex(), input.BlockNumber, 1737786720, "0xf4b4be2ff1373371ff705336b0eff0f70f117bcfe61fb04470e613a1b6e710f6", 0, "bb01"), input.RawData)
	parsedAbi, err := inputs.InputsMetaData.GetAbi()
	suite.Nil(err)
	params, err := decodeInputData(input, parsedAbi)
	suite.Nil(err)
	suite.Equal(int64(suite.chainId), params.ChainId)
	suite.Equal(suite.application.IApplicationAddress.Hex(), params.AppContract)
	suite.Equal(suite.senderAddress, params.MsgSender)
	suite.Equal("0", params.Index)
	suite.Equal("bb01", params.Payload)

	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0x4d3c276e1b15b416979e2080022f24f5e8dfe0767591929e0c3595719cca0b34"), input.TransactionReference)

	// input1
	input, err = suite.database.GetInput(suite.ctx, suite.application.IApplicationAddress.Hex(), 1)
	suite.Nil(err)
	// suite.Equal(uint64(756613), input.EpochIndex)
	suite.Equal(uint64(1), input.Index)
	// suite.Equal(uint64(7566135), input.BlockNumber)
	// suite.Equal(EvmAdvanceEncode(suite.application.IApplicationAddress.Hex(), 7566135, 1737786720, "0xf4b4be2ff1373371ff705336b0eff0f70f117bcfe61fb04470e613a1b6e710f6", 1, "bb02"), input.RawData)
	parsedAbi, err = inputs.InputsMetaData.GetAbi()
	suite.Nil(err)
	params, err = decodeInputData(input, parsedAbi)
	suite.Nil(err)
	suite.Equal(int64(suite.chainId), params.ChainId)
	suite.Equal(suite.application.IApplicationAddress.Hex(), params.AppContract)
	suite.Equal(suite.senderAddress, params.MsgSender)
	suite.Equal("1", params.Index)
	suite.Equal("bb02", params.Payload)

	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0xb0ebbaccd97d9e04fa0c5b1d7df2db59843db74e40813ca56d7af19b587ce300"), input.TransactionReference)

	// input2
	input, err = suite.database.GetInput(suite.ctx, suite.application.IApplicationAddress.Hex(), 2)
	suite.Nil(err)
	// suite.Equal(uint64(756621), input.EpochIndex)
	suite.Equal(uint64(2), input.Index)
	// suite.Equal(uint64(7566214), input.BlockNumber)
	// suite.Equal(EvmAdvanceEncode(suite.application.IApplicationAddress.Hex(), 7566214, 1737787692, "0x9171073480935a9d1bafa8c6380716a352bc1471c18dc091bcdfdec6d0c24072", 2, "aa03"), input.RawData)
	parsedAbi, err = inputs.InputsMetaData.GetAbi()
	suite.Nil(err)
	params, err = decodeInputData(input, parsedAbi)
	suite.Nil(err)
	suite.Equal(int64(suite.chainId), params.ChainId)
	suite.Equal(suite.application.IApplicationAddress.Hex(), params.AppContract)
	suite.Equal(suite.senderAddress, params.MsgSender)
	suite.Equal("2", params.Index)
	suite.Equal("aa03", params.Payload)

	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0x00"), input.TransactionReference)

	// input3
	input, err = suite.database.GetInput(suite.ctx, suite.application.IApplicationAddress.Hex(), 3)
	suite.Nil(err)
	// suite.Equal(uint64(756621), input.EpochIndex)
	suite.Equal(uint64(3), input.Index)
	// suite.Equal(uint64(7566216), input.BlockNumber)
	// suite.Equal(EvmAdvanceEncode(suite.application.IApplicationAddress.Hex(), 7566216, 1737787728, "0x63ed554373a61df9fbcbc91ea12e1d0b6a4e7e108837c88204194c244c8ff8fe", 3, "aa04"), input.RawData)
	parsedAbi, err = inputs.InputsMetaData.GetAbi()
	suite.Nil(err)
	params, err = decodeInputData(input, parsedAbi)
	suite.Nil(err)
	suite.Equal(int64(suite.chainId), params.ChainId)
	suite.Equal(suite.application.IApplicationAddress.Hex(), params.AppContract)
	suite.Equal(suite.senderAddress, params.MsgSender)
	suite.Equal("3", params.Index)
	suite.Equal("aa04", params.Payload)

	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0x01"), input.TransactionReference)

	// input4
	input, err = suite.database.GetInput(suite.ctx, suite.application.IApplicationAddress.Hex(), 4)
	suite.Nil(err)
	// suite.Equal(uint64(756622), input.EpochIndex)
	suite.Equal(uint64(4), input.Index)
	// suite.Equal(uint64(7566228), input.BlockNumber)
	// suite.Equal(EvmAdvanceEncode(suite.application.IApplicationAddress.Hex(), 7566228, 1737787872, "0x2bbd0c2b7dbbced136cfef527cde4c0fdd6d954d31341c5ee485e64cc15e6d83", 4, "bb05"), input.RawData)
	parsedAbi, err = inputs.InputsMetaData.GetAbi()
	suite.Nil(err)
	params, err = decodeInputData(input, parsedAbi)
	suite.Nil(err)
	suite.Equal(int64(suite.chainId), params.ChainId)
	suite.Equal(suite.application.IApplicationAddress.Hex(), params.AppContract)
	suite.Equal(suite.senderAddress, params.MsgSender)
	suite.Equal("4", params.Index)
	suite.Equal("bb05", params.Payload)

	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0x4e63b1d5768dee61ede900a3bb95998079484dc412bc099a35bc902be7544102"), input.TransactionReference)

	// input5
	input, err = suite.database.GetInput(suite.ctx, suite.application.IApplicationAddress.Hex(), 5)
	suite.Nil(err)
	// suite.Equal(uint64(756622), input.EpochIndex)
	suite.Equal(uint64(5), input.Index)
	// suite.Equal(uint64(7566228), input.BlockNumber)
	// suite.Equal(EvmAdvanceEncode(suite.application.IApplicationAddress.Hex(), 7566228, 1737787872, "0x2bbd0c2b7dbbced136cfef527cde4c0fdd6d954d31341c5ee485e64cc15e6d83", 5, "bb06"), input.RawData)
	parsedAbi, err = inputs.InputsMetaData.GetAbi()
	suite.Nil(err)
	params, err = decodeInputData(input, parsedAbi)
	suite.Nil(err)
	suite.Equal(int64(suite.chainId), params.ChainId)
	suite.Equal(suite.application.IApplicationAddress.Hex(), params.AppContract)
	suite.Equal(suite.senderAddress, params.MsgSender)
	suite.Equal("5", params.Index)
	suite.Equal("bb06", params.Payload)

	suite.Equal(model.InputCompletionStatus_Accepted, input.Status)
	suite.Equal(common.HexToHash("0x248ffa7972cd567d963734e12813be3020245d2617701021e7fab47579de565f"), input.TransactionReference)
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
	suite.Equal("echo-dapp", appInDB.Name)
	suite.Equal(suite.application.IApplicationAddress, appInDB.IApplicationAddress)
	suite.Equal(suite.application.IConsensusAddress, appInDB.IConsensusAddress)
	suite.Equal(suite.application.TemplateHash, appInDB.TemplateHash)
	suite.Equal(suite.application.TemplateURI, appInDB.TemplateURI)
	suite.Equal(model.ApplicationState_Enabled, appInDB.State)
	suite.Equal(uint64(6), appInDB.ProcessedInputs)
}

// func (suite *EspressoReaderTestSuite) TestEpoch() {
// 	epoch0, err := suite.database.GetEpoch(suite.ctx, suite.application.IApplicationAddress.Hex(), 756613)
// 	suite.Nil(err)
// 	suite.Equal(int64(1), epoch0.ApplicationID)
// 	suite.Equal(uint64(756613), epoch0.Index)
// 	suite.Equal(uint64(7566130), epoch0.FirstBlock)
// 	suite.Equal(uint64(7566139), epoch0.LastBlock)
// 	// suite.Equal(, epoch0.ClaimHash)
// 	// suite.Equal(, epoch0.ClaimTransactionHash)
// 	// suite.Equal(, epoch0.Status)
// 	suite.Equal(uint64(0), epoch0.VirtualIndex)

// 	epoch1, err := suite.database.GetEpoch(suite.ctx, suite.application.IApplicationAddress.Hex(), 756621)
// 	suite.Nil(err)
// 	suite.Equal(int64(1), epoch1.ApplicationID)
// 	suite.Equal(uint64(756621), epoch1.Index)
// 	suite.Equal(uint64(7566210), epoch1.FirstBlock)
// 	suite.Equal(uint64(7566219), epoch1.LastBlock)
// 	// suite.Equal(, epoch1.ClaimHash)
// 	// suite.Equal(, epoch1.ClaimTransactionHash)
// 	// suite.Equal(, epoch1.Status)
// 	suite.Equal(uint64(1), epoch1.VirtualIndex)

// 	epoch2, err := suite.database.GetEpoch(suite.ctx, suite.application.IApplicationAddress.Hex(), 756622)
// 	suite.Nil(err)
// 	suite.Equal(int64(1), epoch2.ApplicationID)
// 	suite.Equal(uint64(756622), epoch2.Index)
// 	suite.Equal(uint64(7566220), epoch2.FirstBlock)
// 	suite.Equal(uint64(7566229), epoch2.LastBlock)
// 	// suite.Equal(, epoch2.ClaimHash)
// 	// suite.Equal(, epoch2.ClaimTransactionHash)
// 	// suite.Equal(, epoch2.Status)
// 	suite.Equal(uint64(2), epoch2.VirtualIndex)
// }

func (suite *EspressoReaderTestSuite) TestNodeConfig() {
	nodeConfig, err := repository.LoadNodeConfig[evmreader.PersistentConfig](suite.ctx, suite.database, evmreader.EvmReaderConfigKey)
	suite.Nil(err)
	// suite.Equal(model.DefaultBlock_Finalized, nodeConfig.Value.DefaultBlock)
	suite.Equal(uint64(13370), nodeConfig.Value.ChainID)
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

type EvmAdvance struct {
	ChainId        int64  `json:"chainId"`
	AppContract    string `json:"appContract"`
	MsgSender      string `json:"msgSender"`
	BlockNumber    string `json:"blockNumber"`
	BlockTimestamp string `json:"blockTimestamp"`
	PrevRandao     string `json:"prevRandao"`
	Index          string `json:"index"`
	Payload        string `json:"payload"`
}

func decodeInputData(input *model.Input, parsedAbi *abi.ABI) (EvmAdvance, error) {
	decoded := make(map[string]any)
	err := parsedAbi.Methods["EvmAdvance"].Inputs.UnpackIntoMap(decoded, input.RawData[4:])
	if err != nil {
		return EvmAdvance{}, err
	}

	params := EvmAdvance{
		ChainId:        decoded["chainId"].(*big.Int).Int64(),
		AppContract:    decoded["appContract"].(common.Address).Hex(),
		MsgSender:      decoded["msgSender"].(common.Address).Hex(),
		BlockNumber:    decoded["blockNumber"].(*big.Int).String(),
		BlockTimestamp: decoded["blockTimestamp"].(*big.Int).String(),
		PrevRandao:     decoded["prevRandao"].(*big.Int).String(),
		Index:          decoded["index"].(*big.Int).String(),
		Payload:        hex.EncodeToString(decoded["payload"].([]byte)),
	}

	return params, nil
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
