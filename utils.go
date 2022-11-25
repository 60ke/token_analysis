package main

import (
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

// getPercent 获取账户余额在token中的占比
// 灵敏度为十万分之一
//
//	@param chainType
//	@param balance
//	@return string
func getPercent(chainType string, balance *big.Int) string {

	sensitivity, _ := new(big.Float).SetString("0.00001")
	if chainType == "BSC" {
		x, _ := new(big.Float).SetString(balance.String())
		y, _ := new(big.Float).SetString(BscTotal.String())
		z := new(big.Float).Quo(x, y)

		if z.Cmp(sensitivity) > 0 {
			return z.String()
		}
		return "0"
	} else {
		x, _ := new(big.Float).SetString(balance.String())
		y, _ := new(big.Float).SetString(EthTotal.String())
		z := new(big.Float).Quo(x, y)
		if z.Cmp(sensitivity) > 0 {
			return z.String()
		}
		return "0"
	}
}

// getBalanceNumber
// 将10进制的balance字符串转换为big int
func getBalanceNumber(s string) *big.Int {
	balance := new(big.Int)
	_, ok := balance.SetString(s, 10)
	if !ok {
		Logger.Error("getBalanceNumber Err: ", s)
	}
	return balance
}

// http get
//
//	@param url
//	@return []byte
//	@return error
//
// Returns true if the request was successful.
func get(url string) ([]byte, error) {
	// url := "http://106.3.133.179:46657/tri_block_info?height=104360"

	client := &http.Client{}
	client.Timeout = time.Second * 5
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// http post
//
//	@param url
//	@param payload
//	@return []byte
//	@return error
func post(url string, payload *strings.Reader) ([]byte, error) {

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func CheckChainType(chainType string) bool {
	return slices.Contains(ChainTypes, chainType)
}
