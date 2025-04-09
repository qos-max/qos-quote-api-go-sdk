package qosapi

// 基础信息
type InstrumentInfo struct {
	Code              string `json:"c"`  // 股票代码
	Exchange          string `json:"e"`  // 交易所
	TradeCurrency     string `json:"tc"` // 交易币种
	NameCN            string `json:"nc"` // 中文名称
	NameEN            string `json:"ne"` // 英文名称
	LotSize           int64  `json:"ls"` // 最小交易单位
	TotalShares       int64  `json:"ts"` // 总股本
	OutstandingShares int64  `json:"os"` // 流通股本
	EPS               string `json:"ep"` // 每股盈利
	NAV               string `json:"na"` // 每股净资产
	DividendYield     string `json:"dy"` // 股息率
}

// 行情快照
type Snapshot struct {
	Code             string        `json:"c"`  // 股票代码
	LastPrice        string        `json:"lp"` // 当前价格
	PrevClose        string        `json:"yp"` // 昨日收盘价
	Open             string        `json:"o"`  // 开盘价
	High             string        `json:"h"`  // 最高价
	Low              string        `json:"l"`  // 最低价
	Timestamp        int64         `json:"ts"` // 时间戳
	Volume           string        `json:"v"`  // 成交量
	Turnover         string        `json:"t"`  // 成交金额
	Suspended        int           `json:"s"`  // 是否停牌
	PreMarket        *SessionQuote `json:"pq"` // 盘前数据
	AfterMarket      *SessionQuote `json:"aq"` // 盘后数据
	NightMarket      *SessionQuote `json:"nq"` // 夜盘数据
	TradeSessionType int           `json:"tt"` // 交易时段类型
}

// 交易时段行情
type SessionQuote struct {
	LastPrice string `json:"lp"` // 当前价格
	PrevClose string `json:"yp"` // 上次收盘价
	High      string `json:"h"`  // 最高价
	Low       string `json:"l"`  // 最低价
	Timestamp int64  `json:"ts"` // 时间戳
	Volume    string `json:"v"`  // 成交量
	Turnover  string `json:"t"`  // 成交金额
}

// 盘口深度
type Depth struct {
	Code      string      `json:"c"`  // 股票代码
	Bids      []DepthItem `json:"b"`  // 买单数组
	Asks      []DepthItem `json:"a"`  // 卖单数组
	Timestamp int64       `json:"ts"` // 时间戳
}

// 盘口项
type DepthItem struct {
	Price  string `json:"p"` // 价格
	Volume string `json:"v"` // 数量
}

// 逐笔成交
type Trade struct {
	Code      string `json:"c"`  // 股票代码
	Price     string `json:"p"`  // 当前价格
	Volume    string `json:"v"`  // 当前成交量
	Timestamp int64  `json:"ts"` // 时间戳
	Direction int    `json:"d"`  // 交易方向
}

// K线数据
type KLine struct {
	Code      string `json:"c"`  // 股票代码
	Open      string `json:"o"`  // 开盘价
	Close     string `json:"cl"` // 收盘价
	High      string `json:"h"`  // 最高价
	Low       string `json:"l"`  // 最低价
	Volume    string `json:"v"`  // 成交量
	Timestamp int64  `json:"ts"` // 时间戳
	KLineType int    `json:"kt"` // K线类型
}

// K线请求
type KLineRequest struct {
	Codes     string `json:"c"`           // 股票代码，多个用逗号分隔
	Count     int    `json:"co"`          // 请求数量
	Adjust    int    `json:"a"`           // 复权类型 0:不复权 1:前复权
	KLineType int    `json:"kt"`          // K线类型
	EndTime   int64  `json:"e,omitempty"` // 结束时间戳(仅历史K线需要)
}

// 基础响应
type BaseResponse struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// WebSocket请求
type WSRequest struct {
	Type      string         `json:"type"`
	Codes     []string       `json:"codes,omitempty"`
	Count     int            `json:"count,omitempty"`
	KLineType int            `json:"kt,omitempty"`
	ReqID     int            `json:"reqid,omitempty"`
	KLineReqs []KLineRequest `json:"kline_reqs,omitempty"`
}

// WebSocket响应
type WSResponse struct {
	Type  string      `json:"type"`
	Msg   string      `json:"msg"`
	Time  int64       `json:"time,omitempty"`
	ReqID int         `json:"reqid,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

// WebSocket行情快照
type WSSnapshot struct {
	Type             string `json:"tp"` // 数据类型 S
	Code             string `json:"c"`  // 股票代码
	LastPrice        string `json:"lp"` // 当前价格
	PrevClose        string `json:"yp"` // 昨日收盘价
	Open             string `json:"o"`  // 开盘价
	High             string `json:"h"`  // 最高价
	Low              string `json:"l"`  // 最低价
	Timestamp        int64  `json:"ts"` // 时间戳
	Volume           string `json:"v"`  // 成交量
	Turnover         string `json:"t"`  // 成交金额
	Suspended        int    `json:"s"`  // 是否停牌
	TradeSessionType int    `json:"tt"` // 交易时段类型
}

// WebSocket逐笔成交
type WSTrade struct {
	Type      string `json:"tp"` // 数据类型 T
	Code      string `json:"c"`  // 股票代码
	Price     string `json:"p"`  // 当前价格
	Volume    string `json:"v"`  // 当前成交量
	Timestamp int64  `json:"ts"` // 时间戳
	Direction int    `json:"d"`  // 交易方向
}

// WebSocket盘口深度
type WSDepth struct {
	Type      string      `json:"tp"` // 数据类型 D
	Code      string      `json:"c"`  // 股票代码
	Bids      []DepthItem `json:"b"`  // 买单数组
	Asks      []DepthItem `json:"a"`  // 卖单数组
	Timestamp int64       `json:"ts"` // 时间戳
}

// WebSocket K线
type WSKLine struct {
	Type      string `json:"tp"` // 数据类型 K
	Code      string `json:"c"`  // 股票代码
	Open      string `json:"o"`  // 开盘价
	Close     string `json:"cl"` // 收盘价
	High      string `json:"h"`  // 最高价
	Low       string `json:"l"`  // 最低价
	Volume    string `json:"v"`  // 成交量
	Timestamp int64  `json:"ts"` // 时间戳
	KLineType int    `json:"kt"` // K线类型
}
