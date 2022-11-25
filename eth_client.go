package main

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/chenzhijie/go-web3"
	"github.com/chenzhijie/go-web3/types"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ABI string = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`
)

// getBlance
//
//	获取账户余额
//	@param web3
//	@param chainType BSC or ETH
//	@param addr
//	@return *big.Int
func getBlance(web3 *web3.Web3, chainType, addr string) (*big.Int, error) {
	token := getToken(chainType)
	contract, err := web3.Eth.NewContract(ABI, token)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	// totalSupply, err := contract.Call("totalSupply")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Total supply %v\n", totalSupply)

	// data, _ := contract.EncodeABI("balanceOf", web3.Eth.Address())
	// fmt.Printf("%x\n", data)
	balance, err := contract.Call("balanceOf", common.HexToAddress(addr))
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	b := balance.(*big.Int)
	Logger.Infof("%s balance is %s", addr, b)
	return balance.(*big.Int), nil
}

func getWeb3(chainType string) (*web3.Web3, error) {
	var rpc string
	if chainType == "BSC" {
		rpc = APPSetting.BscRpc
	} else {
		rpc = APPSetting.EthRpc
	}
	Logger.Info(rpc)
	web3, err := web3.NewWeb3(rpc)
	if err != nil {
		return nil, err
	}
	web3.Eth.SetAccount(APPSetting.PrivateKey)
	blockNumber, err := web3.Eth.GetBlockNumber()
	if err != nil {
		Logger.Error(err)
		return nil, nil
		// panic(err)
	}
	Logger.Info("Current block number: ", blockNumber)
	return web3, nil
}

func getToken(chainType string) string {
	if chainType == "BSC" {
		return APPSetting.BscAddr
	} else {
		return APPSetting.EthAddr
	}
}

// 从数据库获取当前追块起始数据区块号
// 仅在程序初次运行时使用
//
//	@param chainType
//	@return string
func getDbCurrentBlock(chainType string) string {
	var status Status
	Sql.Table("status").Find(&status)
	if chainType == "BSC" {
		return status.BscCurrentBlock
	} else {
		return status.EthCurrentBlock
	}
}

func saveCurrentBlock(chainType, currentBlock string) {
	var status Status
	Sql.First(&status)
	if chainType == "BSC" {
		status.BscCurrentBlock = currentBlock
		Sql.Save(&status)
	} else {
		status.EthCurrentBlock = currentBlock
		Sql.Save(&status)
	}
}

func getInterval(chainType string) int64 {
	if chainType == "BSC" {
		return APPSetting.BscInterval
	} else {
		return APPSetting.EthInterval
	}
}

func getStepper(chainType string) uint64 {
	if chainType == "BSC" {
		return APPSetting.BscStepper
	} else {
		return APPSetting.EthStepper
	}
}

func Suscribe(chainType string) {
	Logger.Infof("订阅%s出块\n", chainType)

start:
	web3, err := getWeb3(chainType)
	if err != nil {
		Logger.Error(err)
		time.Sleep(10 * time.Second)
		goto start
	}
	var fromBlock string
	stepper := getStepper(chainType)
	token := getToken(chainType)
	fromBlock = getDbCurrentBlock(chainType)
	interval := getInterval(chainType)
	table := getTransferTable(chainType)
	for {
		fromBlockUint, _ := strconv.ParseUint(fromBlock[2:], 16, 64)
		// 获取当前链的块高
		chainblockNumber, err := web3.Eth.GetBlockNumber()
		if err != nil {
			Logger.Error(err)
			chainblockNumber = fromBlockUint
		}
		Logger.Debug("Current block number: ", chainblockNumber)
		Logger.Debug(fromBlock)

		toBlockUint := fromBlockUint
		var toBlock string
		// 追块到chainblockNumber
		for toBlockUint < chainblockNumber {

			toBlockUint += stepper
			if toBlockUint < chainblockNumber {
				toBlock = "0x" + fmt.Sprintf("%x", toBlockUint)

			} else {
				toBlock = "0x" + fmt.Sprintf("%x", chainblockNumber)
			}
			err := QueryBlock(fromBlock, toBlock, token, table, chainType, web3)
			if err != nil {
				// updateFromBlockCache(fromBlock, chainType)
				fromBlock = toBlock
				updateDbFromBlock(fromBlock, chainType)
				UpdateFromBlockCache(toBlockUint, chainType)

			}

		}
		time.Sleep(time.Duration(interval) * time.Second)

	}
}

func UpdateHolders(web3 *web3.Web3, chainType string, addrs ...string) {
	for _, addr := range addrs {
		balance, err := getBlance(web3, chainType, addr)
		// TODO 获取失败的账户应加入失败列表重新获取
		if err != nil {
			return
		}
		UpdateHoldersCache(chainType, addr, balance.String())
		UpdateHoldersDB(chainType, addr, balance.String())
	}

}

// 全量更新持有者余额
func UpdateAllBalance(chainType string) {
	// holders := getHolders(chainType)
	// table :=
	// for _,holder := range holders{

	// }
}

func QueryBlock(fromBlock, toBlock, token, table, chainType string, web3 *web3.Web3) error {
	Logger.Debug(fromBlock, toBlock, token)
	transferTopic := APPSetting.TransferTopic
	fliter := types.Fliter{
		Address:   common.HexToAddress(token),
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Topics:    []string{transferTopic},
	}
	Logger.Debug(chainType, fliter)
	logs, err := web3.Eth.GetLogs(&fliter)
	if err != nil {
		Logger.Error(err)
		return err
	}
	Logger.Debug(chainType, len(logs))
	for _, log := range logs {

		/*
			log data example:
			Data: 0x0000000000000000000000000000000000000000000000037ee5c1d792d1295e
			Topics: [0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef \
			0x00000000000000000000000091b43c14050a35af42eac1e5c9a3153127aa072e \
			0x000000000000000000000000118ed46a6ea1ad0938cd637536d4351734ce16ee]
			BlockNumber: 0x1628ca7
			TransactionHash: 0x0490f5e2a3b6c832530f6a798119eaaa65077b5f5ce2d371613980d32e562043

		*/

		var transfer Transfer
		fmt.Println(log.Data)
		fmt.Println(log.BlockNumber)
		fmt.Println(log.TransactionHash)
		transfer.Hash = log.TransactionHash.Hex()
		transfer.Block = log.BlockNumber

		transfer.From = common.HexToAddress(log.Topics[1]).String()
		transfer.To = common.HexToAddress(log.Topics[2]).String()
		transfer.Value = log.Data
		blockNumber, _ := new(big.Int).SetString(log.BlockNumber, 16)
		block, err := web3.Eth.GetBlocByNumber(blockNumber, true)
		if err != nil {
			Logger.Error(err)
			return err
		}
		transfer.Time = int64(block.Time())

		// unixTimeUTC := time.Unix(int64(block.Time()), 0)
		// unixTimeUTC.Format(time.RFC3339)
		Logger.Debug(transfer)
		Sql.Table(table).Create(&transfer)
		// Sql.Table("transfer").Create(transfer)
		UpdateHolders(web3, chainType, transfer.From, transfer.To)
		// UpdateBalance(chainType, log.Topics[1], log.Topics[2])

	}
	return nil
}
