package utils

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestDoPostRequestWithContext 测试带上下文的 POST 请求
func TestDoPostRequestWithContext(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// 读取并验证请求体
		var data map[string]string
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if data["key"] != "value" {
			t.Errorf("Expected request data key=value, got %s", data["key"])
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("context post response"))
	}))
	defer server.Close()

	// 测试正常情况
	ctx := context.Background()
	requestData := map[string]string{"key": "value"}
	response, err := DoPostRequestWithContext(ctx, server.URL, requestData)
	if err != nil {
		t.Errorf("DoPostRequestWithContext() error = %v", err)
		return
	}

	expected := "context post response"
	if string(response) != expected {
		t.Errorf("DoPostRequestWithContext() = %s, want %s", string(response), expected)
	}

	// 测试上下文取消
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消上下文
	_, err = DoPostRequestWithContext(cancelCtx, server.URL, requestData)
	if err == nil {
		t.Error("DoPostRequestWithContext() expected error for canceled context, got nil")
	}

	// 测试上下文超时
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(1 * time.Millisecond) // 确保上下文超时
	_, err = DoPostRequestWithContext(timeoutCtx, server.URL, requestData)
	if err == nil {
		t.Error("DoPostRequestWithContext() expected error for timeout context, got nil")
	}
}

// TestHTTPClient_RetryBehavior 测试 HTTP 客户端的重试行为
func TestHTTPClient_RetryBehavior(t *testing.T) {
	// 计数器，记录请求次数
	requestCount := 0

	// 创建测试服务器，前两次请求返回 500 错误，第三次返回成功
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server error"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success after retry"))
		}
	}))
	defer server.Close()

	// 创建客户端，设置 2 次重试（总共 3 次请求）
	client := NewHTTPClient(1*time.Second, 2)
	ctx := context.Background()

	// 发送请求
	response, err := client.Get(ctx, server.URL, nil)
	if err != nil {
		t.Errorf("HTTPClient.Get() with retries error = %v", err)
		return
	}

	// 验证结果
	expected := "success after retry"
	if string(response) != expected {
		t.Errorf("HTTPClient.Get() with retries = %s, want %s", string(response), expected)
	}

	// 验证请求次数
	if requestCount != 3 {
		t.Errorf("HTTPClient retry mechanism made %d requests, want 3", requestCount)
	}
}

// TestHTTPClient_ClientErrorNoRetry 测试客户端错误不重试
func TestHTTPClient_ClientErrorNoRetry(t *testing.T) {
	// 计数器，记录请求次数
	requestCount := 0

	// 创建测试服务器，返回 400 错误
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
	}))
	defer server.Close()

	// 创建客户端，设置 2 次重试
	client := NewHTTPClient(1*time.Second, 2)
	ctx := context.Background()

	// 发送请求
	_, err := client.Get(ctx, server.URL, nil)
	if err == nil {
		t.Error("HTTPClient.Get() with 400 error expected error, got nil")
		return
	}

	// 验证请求次数（应该只有 1 次，不重试）
	if requestCount != 1 {
		t.Errorf("HTTPClient made %d requests for client error, want 1 (no retry)", requestCount)
	}
}
