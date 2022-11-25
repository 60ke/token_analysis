package main

import (
	"time"

	"golang.org/x/exp/slices"
)

// GetTopHolders 获取
//
//	@param top
//	@param chainType

// top统计
type TopStat struct {
	Top10  string `json:"top10"`
	Top20  string `json:"top20"`
	Top50  string `json:"top50"`
	Top100 string `json:"top100"`
	// example : 2022-11-23
	Time string `json:"time"`
}

// 返回前{num}的持有者列表
//
//	@param top
//	@param chainType
//	@return []string
func GetTopHolders(num int, chainType string) []string {
	keys := HolderSort(chainType)
	return keys[:num]

}

// 返回近{num}天的日期
//
//	@param num
//	@return []string
func GetDate(num int) []string {
	var dates []string
	now := time.Now().Format(time.RFC3339)[:10]
	dates = append(dates, now)
	for i := 1; i < num; i++ {
		day := -i
		date := time.Now().AddDate(0, 0, day).Format(time.RFC3339)[:10]
		dates = append(dates, date)
	}
	return dates

}

// 返回近{num}天的top 统计表
//
//	@param chainType
//	@return []TopStat
func GetTopChart(num int, chainType string) []TopStat {
	var topStats []TopStat
	m := GetDate(num)
	topCharts := getDbTopChart(chainType)

	for _, topChart := range topCharts {
		if slices.Contains(m, topChart.Time) {
			topStats = append(topStats, topChart.TopStat)
		}
	}
	return topStats
}
