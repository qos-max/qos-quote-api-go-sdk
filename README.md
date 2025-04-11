# QOSAPI Go客户端

[![Go参考文档](https://pkg.go.dev/badge/github.com/qos-max/qos-quote-api-go-sdk.svg)](https://pkg.go.dev/github.com/qos-max/qos-quote-api-go-sdk)
[![许可证](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

QOS行情数据API(Quote Ocean System)的Go客户端库。提供HTTP和WebSocket接口访问股票和加密货币的实时行情数据。

- **官网**：[https://qos.hk](https://qos.hk)
- **API文档**：[https://qos.hk/api.html](https://qos.hk/api.html)
- **注册API KEY**：[https://qos.hk](https://qos.hk)

## 功能特性

- 完整实现所有QOS API接口
- 同时支持HTTP和WebSocket协议
- 类型安全的数据结构
- WebSocket自动重连
- 内置心跳机制
- 完善的错误处理

## 安装

```bash
go get github.com/qos-max/qos-quote-api-go-sdk/qosapi
# 或者
go mod tidy
```

## 使用示例

### HTTP客户端示例

```go
package main

import (
	"fmt"
	"github.com/qos-max/qos-quote-api-go-sdk/qosapi"
	"log"
	"time"
)
//官网:https://qos.hk
//API文档:https://qos.hk/api.html
//注册API KEY:https://qos.hk
func main() {
	// 替换为你的API Key
	apiKey := "your-api-key"
	client := qosapi.NewClient(apiKey)

	// 示例1: 获取基础信息
	codes := []string{
		"US:AAPL,TSLA",
		"HK:00700,09988",
		"SH:600519,600518",
		"SZ:000001,002594",
	}

	info, err := client.GetInstrumentInfo(codes)
	if err != nil {
		log.Fatalf("Failed to get instrument info: %v", err)
	}
	fmt.Println("Instrument Info:")
	for _, item := range info {
		fmt.Printf("%s: %s (%s)\n", item.Code, item.NameCN, item.NameEN)
	}

	// 示例2: 获取行情快照
	snapshots, err := client.GetSnapshot(codes)
	if err != nil {
		log.Fatalf("Failed to get snapshots: %v", err)
	}
	fmt.Println("\nSnapshots:")
	for _, s := range snapshots {
		fmt.Printf("%s: %s (High: %s, Low: %s)\n", s.Code, s.LastPrice, s.High, s.Low)
	}

	// 示例3: 获取盘口深度
	depths, err := client.GetDepth(codes)
	if err != nil {
		log.Fatalf("Failed to get depths: %v", err)
	}
	fmt.Println("\nDepths:")
	for _, d := range depths {
		fmt.Printf("%s: Bids %d, Asks %d\n", d.Code, len(d.Bids), len(d.Asks))
	}

	// 示例4: 获取逐笔成交
	trades, err := client.GetTrade(codes, 5)
	if err != nil {
		log.Fatalf("Failed to get trades: %v", err)
	}
	fmt.Println("\nTrades:")
	for _, t := range trades {
		fmt.Printf("%s: Price %s, Volume %s, Direction %d\n", t.Code, t.Price, t.Volume, t.Direction)
	}

	// 示例5: 获取K线
	klineReqs := []qosapi.KLineRequest{
		{
			Codes:     "US:AAPL,TSLA",
			Count:     2,
			Adjust:    0,
			KLineType: qosapi.KLineTypeDay,
		},
		{
			Codes:     "CF:BTCUSDT,ETHUSDT",
			Count:     2,
			Adjust:    0,
			KLineType: qosapi.KLineTypeDay,
		},
	}

	klineData, err := client.GetKLine(klineReqs)
	if err != nil {
		log.Fatalf("Failed to get KLine: %v", err)
	}
	fmt.Println("\nKLine:")
	for _, klines := range klineData {
		for _, k := range klines {
			fmt.Printf("%s: O:%s C:%s H:%s L:%s V:%s\n",
				k.Code, k.Open, k.Close, k.High, k.Low, k.Volume)
		}
	}

	// 示例6: 获取历史K线
	historyReqs := []qosapi.KLineRequest{
		{
			Codes:     "US:AAPL",
			Count:     2,
			Adjust:    0,
			KLineType: qosapi.KLineTypeDay,
			EndTime:   time.Now().Unix(),
		},
	}

	historyData, err := client.GetHistoryKLine(historyReqs)
	if err != nil {
		log.Fatalf("Failed to get history KLine: %v", err)
	}
	fmt.Println("\nHistory KLine:")
	for _, klines := range historyData {
		for _, k := range klines {
			fmt.Printf("%s: O:%s C:%s H:%s L:%s V:%s\n",
				k.Code, k.Open, k.Close, k.High, k.Low, k.Volume)
		}
	}
}

```

### WebSocket客户端示例

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/qos-max/qos-quote-api-go-sdk/qosapi"
)

// 官网:https://qos.hk
// API文档:https://qos.hk/api.html
// 注册API KEY:https://qos.hk
func main() {
	// 替换为你的API Key
	apiKey := "your-api-key"
	client := qosapi.NewWSClient(apiKey)

	// 连接到WebSocket服务器
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	// 启动心跳
	client.StartHeartbeat(20 * time.Second)

	// 示例1: 订阅实时快照
	codes := []string{"US:AAPL,TSLA", "HK:700"}
	err := client.SubscribeSnapshot(codes, func(s qosapi.WSSnapshot) {
		fmt.Printf("Snapshot - %s: %s (Time: %v)\n", s.Code, s.LastPrice, time.Unix(s.Timestamp, 0))
	})
	if err != nil {
		log.Printf("Failed to subscribe snapshot: %v", err)
	}

	// 示例2: 订阅实时逐笔成交
	err = client.SubscribeTrade(codes, func(t qosapi.WSTrade) {
		fmt.Printf("Trade - %s: %s %s (Dir: %d, Time: %v)\n",
			t.Code, t.Price, t.Volume, t.Direction, time.Unix(t.Timestamp, 0))
	})
	if err != nil {
		log.Printf("Failed to subscribe trade: %v", err)
	}

	// 示例3: 订阅实时盘口
	err = client.SubscribeDepth(codes, func(d qosapi.WSDepth) {
		fmt.Printf("Depth - %s: Bids %d, Asks %d (Time: %v)\n",
			d.Code, len(d.Bids), len(d.Asks), time.Unix(d.Timestamp, 0))
	})
	if err != nil {
		log.Printf("Failed to subscribe depth: %v", err)
	}

	// 示例4: 订阅实时K线
	err = client.SubscribeKLine([]string{"CF:BTCUSDT"}, qosapi.KLineTypeDay, func(k qosapi.WSKLine) {
		fmt.Printf("KLine - %s: O:%s C:%s H:%s L:%s V:%s\n",
			k.Code, k.Open, k.Close, k.High, k.Low, k.Volume)
	})
	if err != nil {
		log.Printf("Failed to subscribe KLine: %v", err)
	}

	// 示例5: 请求实时快照
	go func() {
		time.Sleep(5 * time.Second)
		snapshots, err := client.RequestSnapshot(codes)
		if err != nil {
			log.Printf("Failed to request snapshots: %v", err)
			return
		}
		fmt.Println("\nRequested Snapshots:")
		for _, s := range snapshots {
			fmt.Printf("%s: %s (High: %s, Low: %s)\n", s.Code, s.LastPrice, s.High, s.Low)
		}
	}()

	// 示例6: 请求实时逐笔成交
	go func() {
		time.Sleep(6 * time.Second)
		trades, err := client.RequestTrade(codes, 5)
		if err != nil {
			log.Printf("Failed to request trades: %v", err)
			return
		}
		fmt.Println("\nRequested Trades:")
		for _, t := range trades {
			fmt.Printf("%s: Price %s, Volume %s, Direction %d\n", t.Code, t.Price, t.Volume, t.Direction)
		}
	}()

	// 示例7: 请求实时盘口
	go func() {
		time.Sleep(7 * time.Second)
		depths, err := client.RequestDepth(codes)
		if err != nil {
			log.Printf("Failed to request depths: %v", err)
			return
		}
		fmt.Println("\nRequested Depths:")
		for _, d := range depths {
			fmt.Printf("%s: Bids %d, Asks %d\n", d.Code, len(d.Bids), len(d.Asks))
		}
	}()

	// 示例8: 请求实时K线
	go func() {
		time.Sleep(8 * time.Second)
		klineReqs := []qosapi.KLineRequest{
			{
				Codes:     "US:AAPL,TSLA",
				Count:     2,
				Adjust:    0,
				KLineType: qosapi.KLineTypeDay,
			},
		}

		klineData, err := client.RequestKLine(klineReqs)
		if err != nil {
			log.Printf("Failed to request KLine: %v", err)
			return
		}
		fmt.Println("\nRequested KLine:")
		for _, klines := range klineData {
			for _, k := range klines {
				fmt.Printf("%s: O:%s C:%s H:%s L:%s V:%s\n",
					k.Code, k.Open, k.Close, k.High, k.Low, k.Volume)
			}
		}
	}()

	// 示例9: 请求历史K线
	go func() {
		time.Sleep(9 * time.Second)
		historyReqs := []qosapi.KLineRequest{
			{
				Codes:     "US:AAPL",
				Count:     2,
				Adjust:    0,
				KLineType: qosapi.KLineTypeDay,
				EndTime:   time.Now().Unix(),
			},
		}

		historyData, err := client.RequestHistoryKLine(historyReqs)
		if err != nil {
			log.Printf("Failed to request history KLine: %v", err)
			return
		}
		fmt.Println("\nRequested History KLine:")
		for _, klines := range historyData {
			for _, k := range klines {
				fmt.Printf("%s: O:%s C:%s H:%s L:%s V:%s\n",
					k.Code, k.Open, k.Close, k.High, k.Low, k.Volume)
			}
		}
	}()

	// 示例10: 请求基础信息
	go func() {
		time.Sleep(10 * time.Second)
		info, err := client.RequestInstrumentInfo(codes)
		if err != nil {
			log.Printf("Failed to request instrument info: %v", err)
			return
		}
		fmt.Println("\nRequested Instrument Info:")
		for _, item := range info {
			fmt.Printf("%s: %s (%s)\n", item.Code, item.NameCN, item.NameEN)
		}
	}()

	// 保持连接运行
	select {}
}

```

## API文档

### HTTP接口

- `NewClient(apiKey string) *QOSClient` - 创建HTTP客户端
- `GetInstrumentInfo(codes []string) ([]InstrumentInfo, error)` - 获取品种基础信息
- `GetSnapshot(codes []string) ([]Snapshot, error)` - 获取实时行情快照
- `GetDepth(codes []string) ([]Depth, error)` - 获取盘口深度数据
- `GetTrade(codes []string, count int) ([]Trade, error)` - 获取逐笔成交数据
- `GetKLine(requests []KLineRequest) ([][]KLine, error)` - 获取K线数据
- `GetHistoryKLine(requests []KLineRequest) ([][]KLine, error)` - 获取历史K线数据

### WebSocket接口

- `NewWSClient(apiKey string) *WSClient` - 创建WebSocket客户端
- `Connect() error` - 连接服务器
- `Close() error` - 关闭连接
- `SubscribeSnapshot(codes []string, callback func(WSSnapshot)) error` - 订阅行情快照
- `SubscribeTrade(codes []string, callback func(WSTrade)) error` - 订阅逐笔成交
- `SubscribeDepth(codes []string, callback func(WSDepth)) error` - 订阅盘口深度
- `SubscribeKLine(codes []string, klineType int, callback func(WSKLine)) error` - 订阅K线数据
- `SendHeartbeat() error` - 发送心跳
- `StartHeartbeat(interval time.Duration)` - 启动定时心跳

## 许可证

本项目采用MIT许可证 - 详情见LICENSE文件
License: [MIT](./LICENSE)
