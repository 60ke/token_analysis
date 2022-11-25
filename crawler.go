package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type FXHolder struct {
	Address string `json:"address"`
	// 持仓数量
	Quantity float64 `json:"quantity"`
	// 持仓占比
	Percentage float64 `json:"percentage"`
	// 账户身份(交易所名)
	Platform     string  `json:"platform"`
	PlatformName string  `json:"platform_name"`
	Logo         string  `json:"logo"`
	Change       float64 `json:"change"`
	Blockurl     string  `json:"blockurl"`
	ChangeAbs    float64 `json:"change_abs"`
	// 数据更新时间
	Updatetime  string `json:"updatetime"`
	Hidden      int    `json:"hidden"`
	Destroy     int    `json:"destroy"`
	Iscontract  int    `json:"iscontract"`
	Addressflag string `json:"addressflag"`
}

// 非小号前100持币者json接口
type FXHolders struct {
	Data struct {
		Top struct {
			Updatedate string  `json:"updatedate"`
			Addrcount  int     `json:"addrcount"`
			Top10Rate  float64 `json:"top10rate"`
			Top20Rate  float64 `json:"top20rate"`
			Top50Rate  float64 `json:"top50rate"`
			Top100Rate float64 `json:"top100rate"`
		} `json:"top"`
		Maincoins []struct {
			Maincoin        string `json:"maincoin"`
			Contractaddress string `json:"contractaddress"`
			Description     string `json:"description"`
			Urlprefix       string `json:"urlprefix"`
		} `json:"maincoins"`
		Toplist []FXHEthHolder `json:"toplist"`
		Holders []FXHEthHolder `json:"holders"`
		// 10大流动地址
		Flows    []FXHEthHolder `json:"flows"`
		Holdcoin struct {
			Summary struct {
				Rise       int    `json:"rise"`
				Riserate   int    `json:"riserate"`
				Addrcount  int    `json:"addrcount"`
				Updatedate string `json:"updatedate"`
			} `json:"summary"`
			List []struct {
				Updatedate int     `json:"updatedate"`
				Addrcount  int     `json:"addrcount"`
				Top10Rate  float64 `json:"top10rate"`
				Top20Rate  float64 `json:"top20rate"`
				Top50Rate  float64 `json:"top50rate"`
				Top100Rate float64 `json:"top100rate"`
			} `json:"list"`
		} `json:"holdcoin"`
	} `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Crawler Returns true if the request was successful.
//
//	当前只在非小号(www.feixiaohaozh.info)中测试
//	@param url
//	@return []byte
//	@return error
func crawler(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Logger.Error(err)
		return nil, err

	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"104\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"104\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		Logger.Error(err, url)
		return nil, err

	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)

}

// getVersion 获取当前持有者数据版本
func getVersion() int {
	var eth FXHEthHolder
	Sql.Table("flows").Order("version desc").Find(&eth)
	return eth.Version
}

func patchVersion(v int, holders []FXHEthHolder) []FXHEthHolder {
	for i := range holders {
		holders[i].Version = v
	}
	return holders
}

func CrawlerEthHolder() (FXHolders, error) {
	Logger.Info("start crawler ethHolder")
	var err error
	var data FXHolders
	url := APPSetting.EthHolderAPI
	ret, err := crawler(url)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(ret, &data)
	// 创建新的版本号
	v := getVersion() + 1
	EthCache.Add("version", v)
	flows := patchVersion(v, data.Data.Flows)
	EthCache.Add("flows", flows)
	Sql.Table("flows").Create(flows)
	top100 := patchVersion(v, data.Data.Toplist)
	EthCache.Add("top100", top100)
	Sql.Table("top100").Create(top100)
	return data, err
}
