package qosapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// QOSClient QOS行情API客户端
type QOSClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// NewClient 创建新的QOS客户端
func NewClient(apiKey string) *QOSClient {
	return &QOSClient{
		apiKey:     apiKey,
		httpClient: &http.Client{},
		baseURL:    HTTPBaseURL,
	}
}

// SetHTTPClient 设置自定义HTTP客户端
func (c *QOSClient) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

// SetBaseURL 设置基础URL
func (c *QOSClient) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// doRequest 执行HTTP请求
func (c *QOSClient) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	// 添加API Key到请求头
	req.Header.Add("key", c.apiKey)
	req.Header.Add("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

// GetInstrumentInfo 获取交易品种的基础信息
func (c *QOSClient) GetInstrumentInfo(codes []string) ([]InstrumentInfo, error) {
	req := struct {
		Codes []string `json:"codes"`
	}{
		Codes: codes,
	}

	resp, err := c.doRequest("POST", "/instrument-info", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Msg  string           `json:"msg"`
		Data []InstrumentInfo `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Msg != "OK" {
		return nil, fmt.Errorf("API error: %s", result.Msg)
	}

	return result.Data, nil
}

// GetSnapshot 获取交易品种的实时行情快照
func (c *QOSClient) GetSnapshot(codes []string) ([]Snapshot, error) {
	req := struct {
		Codes []string `json:"codes"`
	}{
		Codes: codes,
	}

	resp, err := c.doRequest("POST", "/snapshot", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Msg  string     `json:"msg"`
		Data []Snapshot `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Msg != "OK" {
		return nil, fmt.Errorf("API error: %s", result.Msg)
	}

	return result.Data, nil
}

// GetDepth 获取交易品种的实时最新盘口深度
func (c *QOSClient) GetDepth(codes []string) ([]Depth, error) {
	req := struct {
		Codes []string `json:"codes"`
	}{
		Codes: codes,
	}

	resp, err := c.doRequest("POST", "/depth", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Msg  string  `json:"msg"`
		Data []Depth `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Msg != "OK" {
		return nil, fmt.Errorf("API error: %s", result.Msg)
	}

	return result.Data, nil
}

// GetTrade 获取交易品种的实时最新逐笔成交明细
func (c *QOSClient) GetTrade(codes []string, count int) ([]Trade, error) {
	req := struct {
		Codes []string `json:"codes"`
		Count int      `json:"count"`
	}{
		Codes: codes,
		Count: count,
	}

	resp, err := c.doRequest("POST", "/trade", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Msg  string  `json:"msg"`
		Data []Trade `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Msg != "OK" {
		return nil, fmt.Errorf("API error: %s", result.Msg)
	}

	return result.Data, nil
}

// GetKLine 获取交易品种的K线
func (c *QOSClient) GetKLine(requests []KLineRequest) ([][]KLine, error) {
	req := struct {
		KLineReqs []KLineRequest `json:"kline_reqs"`
	}{
		KLineReqs: requests,
	}

	resp, err := c.doRequest("POST", "/kline", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Msg  string `json:"msg"`
		Data []struct {
			Code string  `json:"c"`
			K    []KLine `json:"k"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Msg != "OK" {
		return nil, fmt.Errorf("API error: %s", result.Msg)
	}

	klineData := make([][]KLine, len(result.Data))
	for i, item := range result.Data {
		klineData[i] = item.K
	}

	return klineData, nil
}

// GetHistoryKLine 获取交易品种的历史K线
func (c *QOSClient) GetHistoryKLine(requests []KLineRequest) ([][]KLine, error) {
	req := struct {
		KLineReqs []KLineRequest `json:"kline_reqs"`
	}{
		KLineReqs: requests,
	}

	resp, err := c.doRequest("POST", "/history", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Msg  string `json:"msg"`
		Data []struct {
			Code string  `json:"c"`
			K    []KLine `json:"k"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Msg != "OK" {
		return nil, fmt.Errorf("API error: %s", result.Msg)
	}

	klineData := make([][]KLine, len(result.Data))
	for i, item := range result.Data {
		klineData[i] = item.K
	}

	return klineData, nil
}
