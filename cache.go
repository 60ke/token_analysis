package main

import (
	"math/big"
	"sort"
	"strings"

	lru "github.com/hashicorp/golang-lru"
)

var (
	// bsc数据缓存
	BscCache *lru.Cache
	// eth数据缓存
	EthCache *lru.Cache
	// bsc持有者合约地址缓存
	BscHolders = make(Holders)
	// eth持有者合约地址缓存
	EthHolders = make(Holders)
	// bsc当前追块起始数据区块号缓存
	BscCurrentBlock uint64
	// eth当前追块起始数据区块号缓存
	EthCurrentBlock uint64
	// TODO 全局checkChainType
	ChainTypes = []string{"BSC", "ETH"}
)

type Holders map[string]*big.Int

func CacheInit(len int) {

	// getHolders()
	BscCache, _ = lru.New(len)
	EthCache, _ = lru.New(len)

	readHolders2Cache("BSC")
	readHolders2Cache("ETH")

}

func getHoldersCache(chainType string) Holders {
	if chainType == "BSC" {
		return BscHolders
	} else {
		return EthHolders
	}
}

// 将验证者从数据库加载至缓存
//
//	@param chainType
func readHolders2Cache(chainType string) {
	Logger.Infof("加载%s验证者...", chainType)
	cache := getHoldersCache(chainType)
	holders := GetDbHolders(chainType)
	for _, holder := range holders {
		balance := getBalanceNumber(holder.Balance)
		cache[holder.Addr] = balance
	}
	Logger.Infof("%s验证者加载完成", chainType)
}

// 更新新的持有者
func UpdateHoldersCache(chainType, addr, balance string) {
	if chainType == "BSC" {
		UpdateBscHoldersCache(addr, balance)
	} else {
		UpdateEthHoldersCache(addr, balance)
	}
}

// 更新BSC持有者
func UpdateBscHoldersCache(addr, balance string) {
	num := new(big.Int)
	balance = strings.TrimPrefix(balance, "0x")
	num.SetString(balance, 16)
	BscHolders[addr] = num

}

// 更新ETH持有者
func UpdateEthHoldersCache(addr, balance string) {
	num := new(big.Int)
	balance = strings.TrimPrefix(balance, "0x")
	num.SetString(balance, 16)
	EthHolders[addr] = num

}

func HolderSort(chainType string) []string {
	var Holders map[string]*big.Int
	if chainType == "BSC" {
		Holders = BscHolders
	} else {
		Holders = EthHolders
	}
	keys := make([]string, 0, len(Holders))
	for key := range Holders {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return Holders[keys[i]].Cmp(Holders[keys[j]]) > 0
	})
	return keys
}

// func getHolders(chainType string) []string {
// 	var holders []string
// 	if chainType == "BSC" {
// 		for holder := range BscHolders {
// 			holders = append(holders, holder)
// 		}
// 	} else {
// 		for holder := range EthHolders {
// 			holders = append(holders, holder)
// 		}
// 	}
// 	return holders
// }

// 返回持有者合约地址
func GetHolders(chainType string) Holders {

	if chainType == "BSC" {
		return BscHolders
	} else {
		return EthHolders
	}

}
func UpdateFromBlockCache(fromBlock uint64, chainType string) {
	if chainType == "BSC" {
		BscCurrentBlock = fromBlock
	} else {
		EthCurrentBlock = fromBlock
	}
}

// 从缓存中获取当前追块起始数据区块号
//
//	@param chainType
//	@return uint64
func GetCurrentBlock(chainType string) uint64 {
	if chainType == "BSC" {
		return BscCurrentBlock
	} else {
		return EthCurrentBlock
	}
}
