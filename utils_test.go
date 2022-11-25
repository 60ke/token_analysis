package main

import (
	"fmt"
	"math/big"
	"testing"
)

func Test_getPercent(t *testing.T) {
	t1, _ := new(big.Int).SetString("402086", 10)
	fmt.Println(getPercent("ETH", t1))
	fmt.Println(getPercent("BSC", t1))
}
