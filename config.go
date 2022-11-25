package main

import (
	"math/big"
	"time"

	"github.com/go-ini/ini"
)

// APP
type APP struct {
	TokenName string
	// 以太坊持有者获取API
	EthHolderAPI string
	// BSC持有者获取API
	BscHolderAPI string
	BscRpc       string
	// BSC token合约地址
	BscAddr string
	EthRpc  string
	// ETH token合约地址
	EthAddr  string
	LogLevel string
	LogName  string
	// 本地账户 用于从链上获取合约相关数据(账户余额,历史交易等)
	// 账户地址
	LocalAccount string
	// 账户私钥
	LocalPK string
	// 缓存大小
	CacheSize int
	// 程序内置账户私钥
	PrivateKey string
	// token发行量固定不再设置为可配置
	// EthTotal  string
	// BscTotal  string
	// bsc 追块步进
	BscStepper uint64
	// eth 追块步进
	EthStepper uint64
	// 固定为0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
	// 用于获取指定的log
	TransferTopic string
	// BSC同步区块间隔:秒
	BscInterval int64
	// ETH同步区块间隔:秒
	EthInterval int64
}

/*
bsc token 发行量为1500000000000000007368000000000000000000000,
当前(2022年11月21日16:31:08)最大账户余额为 3711298000000000000000000,地址百分占比均为0%,
为了节省计算时间,将故持币数量小于3711298000000000000000000 * 10000的占比均直接设置为0
*/
var (
	APPSetting        = &APP{}
	EthTotal          = big.NewInt(10000000)
	BscTotal, _       = new(big.Int).SetString("1500000000000000007368000000000000000000000", 10)
	PercentMeasure, _ = new(big.Int).SetString("37112980000000000000000000000", 10)
)

type Server struct {
	// 后端服务运行模式release or debug
	RunMode string
	// 后端服务端口
	HttpPort int

	// server超时设置
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

// 数据库相关配置
type Mysql struct {
	User                string
	Password            string
	Host                string
	SQLName             string
	BscTransferTable    string
	EthTransferTable    string
	BscHolderTable      string
	EthHolderTable      string
	BscTopTable         string
	EthTopTable         string
	BscHolderChartTable string
	EthHolderChartTable string
	BscTopChartTable    string
	EthTopChartTable    string
}

var DatabaseSetting = &Mysql{}

// 任务定时配置
type Task struct {
	BscHolderChart string
	EthHolderChart string
	BscTop100      string
	EthTop100      string
	FxhCrawler     string
}

var TaskSetting = &Task{}

var cfg *ini.File

func ConfInit() {
	var err error
	cfg, err = ini.Load("conf.ini")
	if err != nil {
		Logger.Fatalf("setting.Setup, fail to parse 'conf.ini': %v", err)
	}

	mapTo("app", APPSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("task", TaskSetting)
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second

}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		Logger.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
