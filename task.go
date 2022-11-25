package main

import (
	"math/big"

	"github.com/robfig/cron/v3"
)

func StartTask() {
	c := cron.New()

	bscHolderChart := TaskSetting.BscHolderChart
	ethHolderChart := TaskSetting.EthHolderChart
	bscTop100 := TaskSetting.BscTop100
	ethTop100 := TaskSetting.EthTop100

	fxhCrawler := TaskSetting.FxhCrawler

	c.AddFunc(bscHolderChart, func() {
		insertHolder("BSC")
	})

	c.AddFunc(ethHolderChart, func() {
		insertHolder("ETH")
	})

	c.AddFunc(bscTop100, func() {
		insertTop100("BSC")
	})

	c.AddFunc(ethTop100, func() {
		insertTop100("ETH")
	})

	c.AddFunc(fxhCrawler, func() {
		CrawlerEthHolder()
	})

	c.Start()
}

// 插入HolderChart数据
func insertHolder(chainType string) {
	table := getHolderChartTable(chainType)
	holders := getHoldersCache(chainType)
	updateHolderChart(table, len(holders))

}

// 更新Top100 及 对应的表格数据
// chainType BSC/ETH
func insertTop100(chainType string) {
	Logger.Infof("start: %s insertTop100")
	table := getTopTable(chainType)
	chartTable := getTopChartTable(chainType)
	holders := getHoldersCache(chainType)
	top := HolderSort(chainType)[:100]

	total10 := new(big.Int)
	total20 := new(big.Int)
	total50 := new(big.Int)
	total100 := new(big.Int)

	for _, addr := range top[:10] {
		balance := holders[addr]
		total10.Add(total10, balance)
		updateTopTable(table, addr, balance.String())
	}
	total20.Add(total20, total10)
	for _, addr := range top[10:20] {
		balance := holders[addr]
		total20.Add(total20, balance)
		updateTopTable(table, addr, balance.String())
	}
	total50.Add(total20, total50)
	for _, addr := range top[20:50] {
		balance := holders[addr]
		total50.Add(total50, balance)
		updateTopTable(table, addr, balance.String())
	}
	total100.Add(total100, total50)
	for _, addr := range top[50:100] {
		balance := holders[addr]
		total100.Add(total100, balance)
		updateTopTable(table, addr, balance.String())
	}

	top10 := getPercent(chainType, total10)
	top20 := getPercent(chainType, total20)
	top50 := getPercent(chainType, total50)
	top100 := getPercent(chainType, total100)

	updateTopChart(chartTable, top10, top20, top50, top100)

}
