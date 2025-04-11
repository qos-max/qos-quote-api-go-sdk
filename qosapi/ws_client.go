package qosapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WSClient WebSocket客户端
type WSClient struct {
	apiKey      string
	conn        *websocket.Conn
	mu          sync.Mutex
	reqCounter  int
	callbacks   map[int]func(interface{}, error)
	subscribers map[string]func(interface{})
	closeChan   chan struct{}
}

// NewWSClient 创建新的WebSocket客户端
func NewWSClient(apiKey string) *WSClient {
	return &WSClient{
		apiKey:      apiKey,
		callbacks:   make(map[int]func(interface{}, error)),
		subscribers: make(map[string]func(interface{})),
		closeChan:   make(chan struct{}),
	}
}

// Connect 连接到WebSocket服务器
func (c *WSClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return nil
	}

	// 添加API Key到URL参数
	u, err := url.Parse(WSBaseURL)
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("key", c.apiKey)
	u.RawQuery = q.Encode()

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	c.conn = conn

	// 启动读取goroutine
	go c.readLoop()

	return nil
}

// Close 关闭WebSocket连接
func (c *WSClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return nil
	}

	close(c.closeChan)
	err := c.conn.Close()
	c.conn = nil

	return err
}

// readLoop 读取WebSocket消息的循环
func (c *WSClient) readLoop() {
	for {
		select {
		case <-c.closeChan:
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				return
			}

			var baseResp WSResponse
			if err := json.Unmarshal(message, &baseResp); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			if baseResp.Type == "" {
				baseResp.Type = baseResp.TP
			}

			// 处理订阅数据推送
			switch baseResp.Type {
			case "S":
				var snapshot WSSnapshot
				if err := json.Unmarshal(message, &snapshot); err != nil {
					log.Printf("Failed to unmarshal snapshot: %v", err)
					continue
				}
				if cb, ok := c.subscribers["S"]; ok {
					if len(snapshot.Code) > 0 {
						cb(snapshot)
					}
				}
			case "T":
				var trade WSTrade
				if err := json.Unmarshal(message, &trade); err != nil {
					log.Printf("Failed to unmarshal trade: %v", err)
					continue
				}
				if cb, ok := c.subscribers["T"]; ok {
					if len(trade.Code) > 0 {
						cb(trade)
					}
				}
			case "D":
				var depth WSDepth
				if err := json.Unmarshal(message, &depth); err != nil {
					log.Printf("Failed to unmarshal depth: %v", err)
					continue
				}
				if cb, ok := c.subscribers["D"]; ok {
					if len(depth.Code) > 0 {
						cb(depth)
					}
				}
			case "K":
				var kline WSKLine
				if err := json.Unmarshal(message, &kline); err != nil {
					log.Printf("Failed to unmarshal kline: %v", err)
					continue
				}
				if cb, ok := c.subscribers["K"]; ok {
					if len(kline.Code) > 0 {
						cb(kline)
					}
				}
			default:
				// 处理请求响应
				c.mu.Lock()
				if cb, ok := c.callbacks[baseResp.ReqID]; ok {
					delete(c.callbacks, baseResp.ReqID)
					c.mu.Unlock()

					if baseResp.Msg != "OK" {
						cb(nil, errors.New(baseResp.Msg))
					} else {
						cb(baseResp, nil)
					}
				} else {
					c.mu.Unlock()
				}
			}
		}
	}
}

// sendRequest 发送WebSocket请求
func (c *WSClient) sendRequest(req WSRequest, callback func(interface{}, error)) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return errors.New("WebSocket not connected")
	}

	c.reqCounter++
	req.ReqID = c.reqCounter

	if callback != nil {
		c.callbacks[req.ReqID] = callback
	}

	return c.conn.WriteJSON(req)
}

// SubscribeSnapshot 订阅实时快照
func (c *WSClient) SubscribeSnapshot(codes []string, callback func(WSSnapshot)) error {
	key := "S"
	c.subscribers[key] = func(data interface{}) {
		callback(data.(WSSnapshot))
	}

	return c.sendRequest(WSRequest{
		Type:  "S",
		Codes: codes,
	}, nil)
}

// UnsubscribeSnapshot 取消订阅实时快照
func (c *WSClient) UnsubscribeSnapshot(codes []string) error {
	return c.sendRequest(WSRequest{
		Type:  "SC",
		Codes: codes,
	}, nil)
}

// SubscribeTrade 订阅实时逐笔成交
func (c *WSClient) SubscribeTrade(codes []string, callback func(WSTrade)) error {
	key := "T"
	c.subscribers[key] = func(data interface{}) {
		callback(data.(WSTrade))
	}

	return c.sendRequest(WSRequest{
		Type:  "T",
		Codes: codes,
	}, nil)
}

// UnsubscribeTrade 取消订阅实时逐笔成交
func (c *WSClient) UnsubscribeTrade(codes []string) error {
	return c.sendRequest(WSRequest{
		Type:  "TC",
		Codes: codes,
	}, nil)
}

// SubscribeDepth 订阅实时盘口
func (c *WSClient) SubscribeDepth(codes []string, callback func(WSDepth)) error {
	key := "D"
	c.subscribers[key] = func(data interface{}) {
		callback(data.(WSDepth))
	}

	return c.sendRequest(WSRequest{
		Type:  "D",
		Codes: codes,
	}, nil)
}

// UnsubscribeDepth 取消订阅实时盘口
func (c *WSClient) UnsubscribeDepth(codes []string) error {
	return c.sendRequest(WSRequest{
		Type:  "DC",
		Codes: codes,
	}, nil)
}

// SubscribeKLine 订阅实时K线
func (c *WSClient) SubscribeKLine(codes []string, klineType int, callback func(WSKLine)) error {
	key := "K"
	c.subscribers[key] = func(data interface{}) {
		callback(data.(WSKLine))
	}

	return c.sendRequest(WSRequest{
		Type:      "K",
		Codes:     codes,
		KLineType: klineType,
	}, nil)
}

// UnsubscribeKLine 取消订阅实时K线
func (c *WSClient) UnsubscribeKLine(codes []string, klineType int) error {
	return c.sendRequest(WSRequest{
		Type:      "KC",
		Codes:     codes,
		KLineType: klineType,
	}, nil)
}

// RequestSnapshot 请求实时快照
func (c *WSClient) RequestSnapshot(codes []string) ([]Snapshot, error) {
	var result []Snapshot
	errChan := make(chan error, 1)

	err := c.sendRequest(WSRequest{
		Type:  "RS",
		Codes: codes,
	}, func(data interface{}, err error) {
		if err != nil {
			errChan <- err
			return
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			errChan <- err
			return
		}

		var resp struct {
			Data []Snapshot `json:"data"`
		}
		if err := json.Unmarshal(jsonData, &resp); err != nil {
			errChan <- err
			return
		}

		result = resp.Data
		errChan <- nil
	})

	if err != nil {
		return nil, err
	}

	return result, <-errChan
}

// RequestTrade 请求实时逐笔成交
func (c *WSClient) RequestTrade(codes []string, count int) ([]Trade, error) {
	var result []Trade
	errChan := make(chan error, 1)

	err := c.sendRequest(WSRequest{
		Type:  "RT",
		Codes: codes,
		Count: count,
	}, func(data interface{}, err error) {
		if err != nil {
			errChan <- err
			return
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			errChan <- err
			return
		}

		var resp struct {
			Data []Trade `json:"data"`
		}
		if err := json.Unmarshal(jsonData, &resp); err != nil {
			errChan <- err
			return
		}

		result = resp.Data
		errChan <- nil
	})

	if err != nil {
		return nil, err
	}

	return result, <-errChan
}

// RequestDepth 请求实时盘口
func (c *WSClient) RequestDepth(codes []string) ([]Depth, error) {
	var result []Depth
	errChan := make(chan error, 1)

	err := c.sendRequest(WSRequest{
		Type:  "RD",
		Codes: codes,
	}, func(data interface{}, err error) {
		if err != nil {
			errChan <- err
			return
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			errChan <- err
			return
		}

		var resp struct {
			Data []Depth `json:"data"`
		}
		if err := json.Unmarshal(jsonData, &resp); err != nil {
			errChan <- err
			return
		}

		result = resp.Data
		errChan <- nil
	})

	if err != nil {
		return nil, err
	}

	return result, <-errChan
}

// RequestKLine 请求实时K线
func (c *WSClient) RequestKLine(requests []KLineRequest) ([][]KLine, error) {
	var result [][]KLine
	errChan := make(chan error, 1)

	err := c.sendRequest(WSRequest{
		Type:      "RK",
		KLineReqs: requests,
	}, func(data interface{}, err error) {
		if err != nil {
			errChan <- err
			return
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			errChan <- err
			return
		}

		var resp struct {
			Data []struct {
				Code string  `json:"c"`
				K    []KLine `json:"k"`
			} `json:"data"`
		}
		if err := json.Unmarshal(jsonData, &resp); err != nil {
			errChan <- err
			return
		}

		result = make([][]KLine, len(resp.Data))
		for i, item := range resp.Data {
			result[i] = item.K
		}
		errChan <- nil
	})

	if err != nil {
		return nil, err
	}

	return result, <-errChan
}

// RequestHistoryKLine 请求历史K线
func (c *WSClient) RequestHistoryKLine(requests []KLineRequest) ([][]KLine, error) {
	var result [][]KLine
	errChan := make(chan error, 1)

	err := c.sendRequest(WSRequest{
		Type:      "RH",
		KLineReqs: requests,
	}, func(data interface{}, err error) {
		if err != nil {
			errChan <- err
			return
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			errChan <- err
			return
		}

		var resp struct {
			Data []struct {
				Code string  `json:"c"`
				K    []KLine `json:"k"`
			} `json:"data"`
		}
		if err := json.Unmarshal(jsonData, &resp); err != nil {
			errChan <- err
			return
		}

		result = make([][]KLine, len(resp.Data))
		for i, item := range resp.Data {
			result[i] = item.K
		}
		errChan <- nil
	})

	if err != nil {
		return nil, err
	}

	return result, <-errChan
}

// RequestInstrumentInfo 请求交易品种的基础信息
func (c *WSClient) RequestInstrumentInfo(codes []string) ([]InstrumentInfo, error) {
	var result []InstrumentInfo
	errChan := make(chan error, 1)

	err := c.sendRequest(WSRequest{
		Type:  "RI",
		Codes: codes,
	}, func(data interface{}, err error) {
		if err != nil {
			errChan <- err
			return
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			errChan <- err
			return
		}

		var resp struct {
			Data []InstrumentInfo `json:"data"`
		}
		if err := json.Unmarshal(jsonData, &resp); err != nil {
			errChan <- err
			return
		}

		result = resp.Data
		errChan <- nil
	})

	if err != nil {
		return nil, err
	}

	return result, <-errChan
}

// SendHeartbeat 发送心跳
func (c *WSClient) SendHeartbeat() error {
	return c.sendRequest(WSRequest{
		Type: "H",
	}, nil)
}

// StartHeartbeat 启动定时心跳
func (c *WSClient) StartHeartbeat(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := c.SendHeartbeat(); err != nil {
					log.Printf("Failed to send heartbeat: %v", err)
				}
			case <-c.closeChan:
				ticker.Stop()
				return
			}
		}
	}()
}
