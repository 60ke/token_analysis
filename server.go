package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type HolderInfo struct {
	Addr    string `json:"addr"`
	Balance string `json:"balance"`
	Percent string `json:"percent"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode int, msg string, data interface{}) {
	g.C.JSON(httpCode, Response{
		Msg:  msg,
		Data: data,
	})
}

func getFlows(c *gin.Context) {
	// var data []FXHEthHolder
	app := Gin{C: c}
	if v, ok := EthCache.Get("flows"); ok {
		data := v.([]FXHEthHolder)
		Logger.Info(data)
		app.Response(http.StatusOK, "", data)
	}

}

func GetTopMonthChart(c *gin.Context) {
	app := Gin{C: c}
	chainType := c.Query("chain_type")
	if !CheckChainType(chainType) {
		msg := fmt.Sprintf("invaild chain_type: %s", chainType)
		app.Response(http.StatusInternalServerError, msg, nil)
		return
	}
	topStats := GetTopChart(30, chainType)
	Logger.Info(topStats)
	app.Response(http.StatusOK, "", topStats)
}

func GetTop100Holders(c *gin.Context) {
	app := Gin{C: c}
	var holderInfos []HolderInfo
	chainType := c.Query("chain_type")
	if !CheckChainType(chainType) {
		msg := fmt.Sprintf("invaild chain_type: %s", chainType)
		app.Response(http.StatusInternalServerError, msg, nil)
		return
	}
	keys := GetTopHolders(100, chainType)
	cache := getHoldersCache(chainType)
	Logger.Info(len(keys))
	for _, key := range keys {
		var holderInfo HolderInfo
		addr := key
		balance := cache[key].String()
		percent := getPercent(chainType, cache[key])
		holderInfo.Addr = addr
		holderInfo.Balance = balance
		holderInfo.Percent = percent
		holderInfos = append(holderInfos, holderInfo)
	}
	app.Response(http.StatusOK, "", holderInfos)

}

func StartServer(port int) {
	router := gin.Default()

	router.GET("/flows", getFlows)
	router.GET("/top100", GetTop100Holders)
	router.GET("/top_chart", GetTopMonthChart)
	// router.POST("/time_filter", handle)
	router.Run(fmt.Sprintf(":%d", port))
}
