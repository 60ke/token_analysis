Token 数据展示: bsc与eth token数据分析,当前仅做了trias的数据,但可用于所有基于以太坊的公链(如bsc)的所有erc20

## 项目运行
`go run .`
获取编译运行

## Restful 接口说明

### /flows

近 7 天 10 大流动账户地址,接口源于非小号,仅支持以太坊;

> GET

> Param :

> Response:

- msg: 错误信息
- data: 近 7 天 10 大流动账户地址

example:
{
"msg :"",
"data":""
}

### /top100

返回前 100 的 Token 持有者信息

> GET

> Param :

- chain_type : string(链类型:BSC 或 ETH)

> Response:

- msg: 错误信息
- data: 前 100 的 Token 持有者信息

example:

```json

{
    "msg":"",
    "data":[
        {
            "addr":"0x118eD46a6Ea1aD0938cD637536D4351734cE16Ee",
            "balance":"84165811501827467878436",
            "percent":"0"
        },
        {
            "addr":"0x8fe471d0B6269a51D2ceFc4B926853b3375ABb40",
            "balance":"74120000000000000000000",
            "percent":"0"
        },
        {
            "addr":"0x997D5759C560b9c4B461d8C19cD5488A844e3Fcb",
            "balance":"66531706543600000000000",
            "percent":"0"
        }
    ]
}
```
**data长度为100**
### top_chart

近 1 个月的 top10,top20,top50,top100 持有者 Token 占比变化

> GET

> Param :

- chain_type : string(链类型:BSC 或 ETH)

> Response:

- msg: 错误信息
- data: 近 1 个月的 top10,top20,top50,top100 持有者 Token 占比变化

example:
```json
{
    "msg":"",
    "data":[
        {
            "top10":"0",
            "top20":"0",
            "top50":"0",
            "top100":"0",
            "time":"2022-11-25"
        }
    ]
}
```

## 数据抓取

### 链上交易

链上交易为实时数据,程序实时获取

### top100/top_chart

表格数据为统计数据,默认每天 0 点定时写入数据库;可通过 crontab 语法在`conf.ini`配置
