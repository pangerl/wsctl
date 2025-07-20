// Package utils provides common utility functions
// @Author lanpang
// @Date 2024/12/19
// @Desc HTTP-related utility functions
package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"vhagar/errors"
	"vhagar/logger"
)

// HTTPClient HTTP 客户端配置
type HTTPClient struct {
	Client  *http.Client
	Timeout time.Duration
	Retries int
}

// NewHTTPClient 创建新的 HTTP 客户端
func NewHTTPClient(timeout time.Duration, retries int) *HTTPClient {
	return &HTTPClient{
		Client: &http.Client{
			Timeout: timeout,
		},
		Timeout: timeout,
		Retries: retries,
	}
}

// DefaultHTTPClient 创建默认配置的 HTTP 客户端
func DefaultHTTPClient() *HTTPClient {
	return NewHTTPClient(30*time.Second, 3)
}

// Get 发送 GET 请求
func (c *HTTPClient) Get(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeNetworkFailed, "创建 GET 请求失败", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return c.doRequest(req)
}

// Post 发送 POST 请求
func (c *HTTPClient) Post(ctx context.Context, url string, data interface{}, headers map[string]string) ([]byte, error) {
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, errors.Wrap(errors.ErrCodeInvalidParam, "序列化请求数据失败", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeNetworkFailed, "创建 POST 请求失败", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return c.doRequest(req)
}

// Put 发送 PUT 请求
func (c *HTTPClient) Put(ctx context.Context, url string, data interface{}, headers map[string]string) ([]byte, error) {
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, errors.Wrap(errors.ErrCodeInvalidParam, "序列化请求数据失败", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, body)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeNetworkFailed, "创建 PUT 请求失败", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return c.doRequest(req)
}

// Delete 发送 DELETE 请求
func (c *HTTPClient) Delete(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeNetworkFailed, "创建 DELETE 请求失败", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return c.doRequest(req)
}

// doRequest 执行请求（带重试机制）
func (c *HTTPClient) doRequest(req *http.Request) ([]byte, error) {
	var lastErr error
	log := logger.GetLogger()

	for i := 0; i <= c.Retries; i++ {
		if i > 0 {
			log.Warnw("HTTP 请求重试", "attempt", i, "url", req.URL.String())
			// 简单的退避策略
			time.Sleep(time.Duration(i) * time.Second)
		}

		resp, err := c.Client.Do(req)
		if err != nil {
			lastErr = err
			log.Errorw("HTTP 请求失败", "error", err, "url", req.URL.String(), "attempt", i+1)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			lastErr = err
			log.Errorw("读取响应体失败", "error", err, "url", req.URL.String())
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			log.Debugw("HTTP 请求成功", "status", resp.StatusCode, "url", req.URL.String())
			return body, nil
		}

		lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		log.Errorw("HTTP 请求返回错误状态", "status", resp.StatusCode, "url", req.URL.String(), "response", string(body))

		// 对于客户端错误（4xx），不进行重试
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			break
		}
	}

	if lastErr != nil {
		return nil, errors.Wrap(errors.ErrCodeNetworkFailed, "HTTP 请求失败", lastErr)
	}

	return nil, errors.New(errors.ErrCodeNetworkFailed, "HTTP 请求失败，未知错误")
}

// DoRequest 简单的 GET 请求函数（向后兼容）
// 迁移自 task/utils.go 中的 DoRequest 函数
func DoRequest(url string) []byte {
	client := DefaultHTTPClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	data, err := client.Get(ctx, url, nil)
	if err != nil {
		log := logger.GetLogger()
		log.Errorw("DoRequest 失败", "url", url, "error", err)
		return nil
	}

	return data
}

// DoRequestWithContext 带上下文的简单 GET 请求
func DoRequestWithContext(ctx context.Context, url string) ([]byte, error) {
	client := DefaultHTTPClient()
	return client.Get(ctx, url, nil)
}

// DoPostRequest 简单的 POST 请求函数
func DoPostRequest(url string, data interface{}) ([]byte, error) {
	client := DefaultHTTPClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return client.Post(ctx, url, data, nil)
}

// DoPostRequestWithContext 带上下文的 POST 请求
func DoPostRequestWithContext(ctx context.Context, url string, data interface{}) ([]byte, error) {
	client := DefaultHTTPClient()
	return client.Post(ctx, url, data, nil)
}
