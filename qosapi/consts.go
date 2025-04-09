package qosapi

// K线类型
const (
	KLineTypeMin1  = 1    // 1分钟
	KLineTypeMin5  = 5    // 5分钟
	KLineTypeMin15 = 15   // 15分钟
	KLineTypeMin30 = 30   // 30分钟
	KLineTypeHour1 = 60   // 1小时
	KLineTypeHour2 = 120  // 2小时
	KLineTypeHour4 = 240  // 4小时
	KLineTypeDay   = 1001 // 日线
	KLineTypeWeek  = 1007 // 周线
	KLineTypeMonth = 1030 // 月线
	KLineTypeYear  = 2001 // 年线
)

// 交易方向
const (
	TradeDirectionUnknown = 0 // 未知
	TradeDirectionBuy     = 1 // 买入
	TradeDirectionSell    = 2 // 卖出
)

// 美股交易时段类型
const (
	USTradeSessionUnknown    = 0 // 未知
	USTradeSessionNight      = 1 // 夜盘
	USTradeSessionPreMarket  = 2 // 盘前
	USTradeSessionIntraday   = 3 // 盘中
	USTradeSessionAfterHours = 4 // 盘后
)

// 市场代码
const (
	MarketUS = "US" // 美股
	MarketHK = "HK" // 港股
	MarketSH = "SH" // 沪市
	MarketSZ = "SZ" // 深市
	MarketCF = "CF" // 加密货币
)

const (
	HTTPBaseURL = "https://api.qos.hk"
	WSBaseURL   = "wss://api.qos.hk/ws"
)
