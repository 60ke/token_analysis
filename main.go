package main

func main() {

	ConfInit()
	LogInit(APPSetting.LogLevel, APPSetting.LogName)

	DBInit(DatabaseSetting.User, DatabaseSetting.Password, DatabaseSetting.Host, DatabaseSetting.SQLName)
	DBCreate()
	Logger.Info("加载数据")
	CacheInit(APPSetting.CacheSize)
	Logger.Info("数据加载完成")
	// StartServer()
	// go crontab()
	go Suscribe("BSC")
	go Suscribe("ETH")

	StartTask()
	StartServer(ServerSetting.HttpPort)
}
