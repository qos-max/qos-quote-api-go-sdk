# QOSAPI Go客户端

[![Go参考文档](https://pkg.go.dev/badge/github.com/qos-max/qos-quote-api-go-sdk.svg)](https://pkg.go.dev/github.com/qos-max/qos-quote-api-go-sdk)
[![许可证](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

QOS行情数据API(Quote Ocean System)的Go客户端库。提供HTTP和WebSocket接口访问股票和加密货币的实时行情数据。

**官网**：[https://qos.hk](https://qos.hk)
**API文档**：[https://qos.hk/api.html](https://qos.hk/api.html)
**注册API KEY**：[https://qos.hk](https://qos.hk)

## 功能特性

- 完整实现所有QOS API接口
- 同时支持HTTP和WebSocket协议
- 类型安全的数据结构
- WebSocket自动重连
- 内置心跳机制
- 完善的错误处理

## 安装

```bash
go get github.com/qos-max/qos-quote-api-go-sdk/qos-quote-api-go-sdk/qosapi
```

## 使用示例

### HTTP客户端示例

```go
package main

import (
"fmt"
"log"

    "github.com/qos-max/qos-quote-api-go-sdk/qosapi"
)

func main() {
// 创建客户端
client := qosapi.NewClient("您的API密钥")

    // 获取基础信息
    info, err := client.GetInstrumentInfo([]string{"US:AAPL", "HK:00700"})
    if err != nil {
        log.Fatal(err)
    }

    // 打印结果
    for _, item := range info {
        fmt.Printf("%s: %s(总股本:%d)\n", item.Code, item.NameCN, item.TotalShares)
    }
}
```

### WebSocket客户端示例

```go
package main

import (
"log"
"time"

    "github.com/qos-max/qos-quote-api-go-sdk/qosapi"
)

func main() {
// 创建WebSocket客户端
client := qosapi.NewWSClient("您的API密钥")

    // 连接服务器
    if err := client.Connect(); err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // 启动心跳(每20秒一次)
    client.StartHeartbeat(20 * time.Second)

    // 订阅实时行情快照
    if err := client.SubscribeSnapshot([]string{"US:AAPL"}, func(s qosapi.WSSnapshot) {
        log.Printf("行情快照: %s 最新价: %s", s.Code, s.LastPrice)
    }); err != nil {
        log.Fatal(err)
    }

    // 保持连接
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
