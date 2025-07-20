package utils

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewHTTPClient(t *testing.T) {
	timeout := 10 * time.Second
	retries := 5

	client := NewHTTPClient(timeout, retries)

	if client.Timeout != timeout {
		t.Errorf("NewHTTPClient() timeout = %v, want %v", client.Timeout, timeout)
	}

	if client.Retries != retries {
		t.Errorf("NewHTTPClient() retries = %v, want %v", client.Retries, retries)
	}

	if client.Client.Timeout != timeout {
		t.Errorf("NewHTTPClient() client timeout = %v, want %v", client.Client.Timeout, timeout)
	}
}

func TestDefaultHTTPClient(t *testing.T) {
	client := DefaultHTTPClient()

	expectedTimeout := 30 * time.Second
	expectedRetries := 3

	if client.Timeout != expectedTimeout {
		t.Errorf("DefaultHTTPClient() timeout = %v, want %v", client.Timeout, expectedTimeout)
	}

	if client.Retries != expectedRetries {
		t.Errorf("DefaultHTTPClient() retries = %v, want %v", client.Retries, expectedRetries)
	}
}

func TestHTTPClient_Get(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Check headers
		if r.Header.Get("X-Test-Header") != "test-value" {
			t.Errorf("Expected header X-Test-Header=test-value, got %s", r.Header.Get("X-Test-Header"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	client := DefaultHTTPClient()
	ctx := context.Background()
	headers := map[string]string{
		"X-Test-Header": "test-value",
	}

	response, err := client.Get(ctx, server.URL, headers)
	if err != nil {
		t.Errorf("HTTPClient.Get() error = %v", err)
		return
	}

	expected := "test response"
	if string(response) != expected {
		t.Errorf("HTTPClient.Get() = %s, want %s", string(response), expected)
	}
}

func TestHTTPClient_Post(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type=application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Read and verify request body
		var data map[string]string
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if data["key"] != "value" {
			t.Errorf("Expected request data key=value, got %s", data["key"])
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("post response"))
	}))
	defer server.Close()

	client := DefaultHTTPClient()
	ctx := context.Background()
	requestData := map[string]string{"key": "value"}

	response, err := client.Post(ctx, server.URL, requestData, nil)
	if err != nil {
		t.Errorf("HTTPClient.Post() error = %v", err)
		return
	}

	expected := "post response"
	if string(response) != expected {
		t.Errorf("HTTPClient.Post() = %s, want %s", string(response), expected)
	}
}

func TestHTTPClient_Put(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("put response"))
	}))
	defer server.Close()

	client := DefaultHTTPClient()
	ctx := context.Background()
	requestData := map[string]string{"key": "value"}

	response, err := client.Put(ctx, server.URL, requestData, nil)
	if err != nil {
		t.Errorf("HTTPClient.Put() error = %v", err)
		return
	}

	expected := "put response"
	if string(response) != expected {
		t.Errorf("HTTPClient.Put() = %s, want %s", string(response), expected)
	}
}

func TestHTTPClient_Delete(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("delete response"))
	}))
	defer server.Close()

	client := DefaultHTTPClient()
	ctx := context.Background()

	response, err := client.Delete(ctx, server.URL, nil)
	if err != nil {
		t.Errorf("HTTPClient.Delete() error = %v", err)
		return
	}

	expected := "delete response"
	if string(response) != expected {
		t.Errorf("HTTPClient.Delete() = %s, want %s", string(response), expected)
	}
}

func TestHTTPClient_ErrorHandling(t *testing.T) {
	// Test server that returns 500 error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}))
	defer server.Close()

	client := NewHTTPClient(1*time.Second, 1) // Short timeout and few retries for testing
	ctx := context.Background()

	_, err := client.Get(ctx, server.URL, nil)
	if err == nil {
		t.Error("HTTPClient.Get() expected error for 500 status, got nil")
	}
}

func TestHTTPClient_Timeout(t *testing.T) {
	// Test server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Delay longer than client timeout
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("delayed response"))
	}))
	defer server.Close()

	client := NewHTTPClient(500*time.Millisecond, 0) // Short timeout, no retries
	ctx := context.Background()

	_, err := client.Get(ctx, server.URL, nil)
	if err == nil {
		t.Error("HTTPClient.Get() expected timeout error, got nil")
	}
}

func TestDoRequest(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("simple response"))
	}))
	defer server.Close()

	response := DoRequest(server.URL)
	expected := "simple response"

	if string(response) != expected {
		t.Errorf("DoRequest() = %s, want %s", string(response), expected)
	}
}

func TestDoRequest_Error(t *testing.T) {
	// Test with invalid URL
	response := DoRequest("http://invalid-url-that-does-not-exist.local")
	if response != nil {
		t.Errorf("DoRequest() expected nil for invalid URL, got %s", string(response))
	}
}

func TestDoRequestWithContext(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("context response"))
	}))
	defer server.Close()

	ctx := context.Background()
	response, err := DoRequestWithContext(ctx, server.URL)
	if err != nil {
		t.Errorf("DoRequestWithContext() error = %v", err)
		return
	}

	expected := "context response"
	if string(response) != expected {
		t.Errorf("DoRequestWithContext() = %s, want %s", string(response), expected)
	}
}

func TestDoPostRequest(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("post response"))
	}))
	defer server.Close()

	requestData := map[string]string{"key": "value"}
	response, err := DoPostRequest(server.URL, requestData)
	if err != nil {
		t.Errorf("DoPostRequest() error = %v", err)
		return
	}

	expected := "post response"
	if string(response) != expected {
		t.Errorf("DoPostRequest() = %s, want %s", string(response), expected)
	}
}
