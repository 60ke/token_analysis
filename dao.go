package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Sql *gorm.DB

// 非小号的持有者数据,暂时没用到
type FXHEthHolder struct {
	ID int `gorm:"index:"`
	// 数据更新版本
	Version int    `gorm:"column:version" `
	Address string `gorm:"column:address" `
	// 持仓数量
	Quantity float64 `gorm:"column:quantity"`
	// 持仓占比
	Percentage float64 `gorm:"column:percentage"`
	// 账户身份(交易所名)
	Platform     string `gorm:"column:platform"`
	PlatformName string `gorm:"column:platform_name"`
	// 交易所Logo
	Logo      string  `gorm:"column:logo"`
	Change    float64 `gorm:"column:change"`
	Blockurl  string  `gorm:"column:blockurl"`
	ChangeAbs float64 `gorm:"column:change_abs"`
	// 数据更新时间
	Updatetime  string `gorm:"column:updatetime"`
	Hidden      int    `gorm:"column:hidden"`
	Destroy     int    `gorm:"column:destroy"`
	Iscontract  int    `gorm:"column:iscontract"`
	Addressflag string `gorm:"column:addressflag"`
}

type Status struct {
	ID int `gorm:"index:"`
	// 当前块高,为了方便计算,块高类型设置为int64
	BscCurrentBlock string `gorm:"column:bsc_current_block,default:0x1628ca7"`
	EthCurrentBlock string `gorm:"column:eth_current_block,default:0xb44035"`
}

// 记录Transfer交易数据
type Transfer struct {
	ID    int    `gorm:"index:,"`
	Hash  string `gorm:"index:,unique,comment:交易哈希"`
	From  string `gorm:"comment:发送者地址"`
	To    string `gorm:"comment:接受者地址"`
	Value string `gorm:"comment:转账金额"`
	Time  int64  `gorm:"comment:数据更新时间戳"`
	Block string `gorm:"comment:交易所在区块号"`
}

// 持有者实时余额
type Holder struct {
	ID   int    `gorm:"index:,"`
	Addr string `gorm:"index:,"`
	// 余额为10进制表示的big int string
	Balance string
}

// 每天的top100账户余额数据
type Top struct {
	ID   int `gorm:"index:,"`
	Addr string
	// 余额为10进制表示的big int string
	Balance string
	// example : 2022-11-23
	Time string
}

// 大账户持币占比表格
type TopChart struct {
	ID int `gorm:"index:,"`
	TopStat
}

// 存储holders地址变化
type HolderChart struct {
	ID           int `gorm:"index:,"`
	Time         string
	HolderNumber int
}

func DBInit(user, pass, host, dbName string) *gorm.DB {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, dbName)
	// dsn := "k:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	Sql, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("db connect error")
	}
	fmt.Println("db connect success")
	return Sql
}

func DBCreate() {
	// eth_holder 存储top100持有者
	// Sql.Table("top100").AutoMigrate(&FXHEthHolder{})
	// flows 存储10大流动客户
	// Sql.Table("flows").AutoMigrate(&FXHEthHolder{})
	// 存储当前数据状态
	Sql.Table("status").AutoMigrate(&Status{})

	// 存储bsc transfer交易数据
	Sql.Table("bsc_transfer").AutoMigrate(&Transfer{})

	// 存储eth transfer交易数据
	Sql.Table("eth_transfer").AutoMigrate(&Transfer{})

	// 存储eth持有者
	Sql.Table("eth_holder").AutoMigrate(&Holder{})
	// 存储bsc持有者
	Sql.Table("bsc_holder").AutoMigrate(&Holder{})

	// 定时存储 eth top 100
	Sql.Table("eth_top100").AutoMigrate(&Top{})
	// 定时存储 eth top 100
	Sql.Table("bsc_top100").AutoMigrate(&Top{})

	// bsc top chart表格
	Sql.Table("bsc_top_chart").AutoMigrate(&TopChart{})
	// eth top chart表格
	Sql.Table("eth_top_chart").AutoMigrate(&TopChart{})

	// bsc持有者数量统计
	Sql.Table("bsc_holder_chart").AutoMigrate(&HolderChart{})
	Sql.Table("eth_holder_chart").AutoMigrate(&HolderChart{})

}

func DBClose() {
	db, _ := Sql.DB()
	db.Close()
}

// updateDbFromBlock 更新数据库追块游标
//
//	@param fromBlock
//	@param chainType
func updateDbFromBlock(fromBlock, chainType string) {
	var status Status
	Sql.Table("status").First(&status)
	if chainType == "BSC" {
		status.BscCurrentBlock = fromBlock
	} else {
		status.EthCurrentBlock = fromBlock
	}
	Sql.Table("status").Save(&status)
}

func getTransferTable(chainType string) string {
	if chainType == "BSC" {
		return DatabaseSetting.BscTransferTable
	} else {
		return DatabaseSetting.EthTransferTable
	}
}

func getHolderTable(chainType string) string {
	if chainType == "BSC" {
		return DatabaseSetting.BscHolderTable
	} else {
		return DatabaseSetting.EthHolderTable
	}
}

func getTopTable(chainType string) string {
	if chainType == "BSC" {
		return DatabaseSetting.BscTopTable
	} else {
		return DatabaseSetting.EthTopTable
	}
}

// 获取topchar表类型
//
//	@param chainType
//	@return string
func getTopChartTable(chainType string) string {
	if chainType == "BSC" {
		return DatabaseSetting.BscTopChartTable
	} else {
		return DatabaseSetting.EthTopChartTable
	}
}

func getHolderChartTable(chainType string) string {
	if chainType == "BSC" {
		return DatabaseSetting.BscHolderChartTable
	} else {
		return DatabaseSetting.EthHolderChartTable
	}
}

func updateTopTable(table, addr, balance string) {
	var top Top
	// 2022-11-23
	top.Time = time.Now().Format(time.RFC3339)[:10]
	top.Addr = addr
	top.Balance = balance
	Sql.Table(table).Create(&top)
}

func updateTopChart(table, top10, top20, top50, top100 string) {
	var topChart TopChart
	topChart.Time = time.Now().Format(time.RFC3339)[:10]
	topChart.Top10 = top10
	topChart.Top20 = top20
	topChart.Top50 = top50
	topChart.Top100 = top100
	Sql.Table(table).Create(&topChart)

}

func updateHolderChart(table string, num int) {
	var holderChart HolderChart
	holderChart.Time = time.Now().Format(time.RFC3339)[:10]
	holderChart.HolderNumber = num
	Sql.Table(table).Create(&holderChart)
}

func UpdateHoldersDB(chainType, addr, balance string) {
	table := getHolderTable(chainType)
	var holder Holder
	holder.Addr = addr
	holder.Balance = balance
	Sql.Table(table).Where(Holder{Addr: addr}).Attrs(Holder{Balance: balance}).FirstOrCreate(&holder)
	// Sql.Table(table).Create(&Holder{Addr: addr})

}

// GetDbHolders 从数据获取持有者
//
//	@param chainType : BSC/ETH
func GetDbHolders(chainType string) []Holder {
	table := getHolderTable(chainType)
	var holders []Holder
	Sql.Table(table).Find(&holders)
	return holders

}

// getDbTopChart
// TODO 返回近30天的数据而不是全部,以节省内存(不过每天一条的数据内存占用并不大)
//
//	@param chainType
//	@return []TopStat
func getDbTopChart(chainType string) []TopChart {
	var topCharts []TopChart
	table := getTopChartTable(chainType)
	Sql.Table(table).Find(&topCharts)
	return topCharts

}
