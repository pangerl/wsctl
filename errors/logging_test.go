package errors

import (
	"fmt"
	"testing"
	"vhagar/logger"
)

// TestLogError 测试错误日志记录功能
func TestLogError(t *testing.T) {
	// 确保日志器已初始化
	if logger.Logger == nil {
		logger.InitLogger(logger.Config{
			Level: "info",
		})
	}

	// 测试应用错误
	appErr := New(ErrCodeInvalidParam, "应用错误")
	LogError(appErr, "测试上下文")

	// 测试普通错误
	normalErr := fmt.Errorf("普通错误")
	LogError(normalErr, "测试上下文")

	// 测试 nil 错误
	// 这应该不会 panic
	LogError(nil, "测试上下文")
}

// TestLogErrorWithFields 测试带字段的错误日志记录功能
func TestLogErrorWithFields(t *testing.T) {
	// 确保日志器已初始化
	if logger.Logger == nil {
		logger.InitLogger(logger.Config{
			Level: "info",
		})
	}

	// 测试应用错误
	appErr := New(ErrCodeInvalidParam, "应用错误")
	LogErrorWithFields(appErr, "测试上下文", map[string]any{
		"field1": "value1",
		"field2": 123,
	})

	// 测试普通错误
	normalErr := fmt.Errorf("普通错误")
	LogErrorWithFields(normalErr, "测试上下文", map[string]any{
		"field1": "value1",
		"field2": 123,
	})

	// 测试 nil 错误
	// 这应该不会 panic
	LogErrorWithFields(nil, "测试上下文", map[string]any{
		"field1": "value1",
	})

	// 测试空字段
	LogErrorWithFields(appErr, "测试上下文", nil)
	LogErrorWithFields(appErr, "测试上下文", map[string]any{})
}

// TestLogErrorf 测试格式化的错误日志记录功能
func TestLogErrorf(t *testing.T) {
	// 确保日志器已初始化
	if logger.Logger == nil {
		logger.InitLogger(logger.Config{
			Level: "info",
		})
	}

	// 测试应用错误
	appErr := New(ErrCodeInvalidParam, "应用错误")
	LogErrorf(appErr, "格式化消息: %s, %d", "参数1", 123)

	// 测试普通错误
	normalErr := fmt.Errorf("普通错误")
	LogErrorf(normalErr, "格式化消息: %s, %d", "参数1", 123)

	// 测试 nil 错误
	// 这应该不会 panic
	LogErrorf(nil, "格式化消息: %s", "参数1")
}

// TestLogLevels 测试不同级别的日志记录功能
func TestLogLevels(t *testing.T) {
	// 确保日志器已初始化
	if logger.Logger == nil {
		logger.InitLogger(logger.Config{
			Level: "debug",
		})
	}

	// 测试不同级别的日志
	LogDebug("调试消息", "key1", "value1")
	LogInfo("信息消息", "key1", "value1", "key2", 123)
	LogWarn("警告消息", "key1", "value1", "key2", 123, "key3", true)

	// 测试空字段
	LogDebug("调试消息")
	LogInfo("信息消息")
	LogWarn("警告消息")
}

// TestWrapWithDetail 测试包装错误并添加详细信息
func TestWrapWithDetail(t *testing.T) {
	originalErr := fmt.Errorf("原始错误")
	wrappedErr := WrapWithDetail(ErrCodeInternalErr, "包装错误", "详细信息", originalErr)

	if wrappedErr.Code != ErrCodeInternalErr {
		t.Errorf("Expected code %d, got %d", ErrCodeInternalErr, wrappedErr.Code)
	}

	if wrappedErr.Message != "包装错误" {
		t.Errorf("Expected message '包装错误', got '%s'", wrappedErr.Message)
	}

	if wrappedErr.Detail != "详细信息" {
		t.Errorf("Expected detail '详细信息', got '%s'", wrappedErr.Detail)
	}

	if wrappedErr.Cause != originalErr {
		t.Errorf("Expected cause to be original error, got %v", wrappedErr.Cause)
	}
}

// TestErrorCodeString 测试错误码字符串表示
func TestErrorCodeString(t *testing.T) {
	tests := []struct {
		code ErrorCode
		want string
	}{
		{ErrCodeSuccess, "0"},
		{ErrCodeInternalErr, "10001"},
		{ErrCodeInvalidParam, "10002"},
		{ErrCodeNotFound, "10003"},
		{ErrCodeUnauthorized, "10004"},
		{ErrCodeForbidden, "10005"},
		{ErrCodeAIProviderNotFound, "20001"},
		{ErrCodeAIRequestFailed, "20002"},
		{ErrCodeAIResponseInvalid, "20003"},
		{ErrCodeToolNotFound, "30001"},
		{ErrCodeToolCallFailed, "30002"},
		{ErrCodeToolRegFailed, "30003"},
		{ErrCodeConfigNotFound, "40001"},
		{ErrCodeConfigInvalid, "40002"},
		{ErrCodeNetworkTimeout, "50001"},
		{ErrCodeNetworkFailed, "50002"},
		{ErrCodeDBConnFailed, "60001"},
		{ErrCodeDBQueryFailed, "60002"},
		{ErrCodeDBUpdateFailed, "60003"},
		{ErrCodeDBDeleteFailed, "60004"},
		{ErrCodeDBInsertFailed, "60005"},
		{ErrCodeDBTxFailed, "60006"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("ErrorCode_%d", tt.code), func(t *testing.T) {
			got := fmt.Sprintf("%d", tt.code)
			if got != tt.want {
				t.Errorf("ErrorCode.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNilErrorHandling 测试处理 nil 错误的情况
func TestNilErrorHandling(t *testing.T) {
	// 测试 Wrap 函数处理 nil 错误
	wrappedErr := Wrap(ErrCodeInternalErr, "包装错误", nil)
	if wrappedErr.Detail != "" {
		t.Errorf("Expected empty detail for nil error, got '%s'", wrappedErr.Detail)
	}
	if wrappedErr.Cause != nil {
		t.Errorf("Expected nil cause, got %v", wrappedErr.Cause)
	}

	// 测试 IsAppError 函数处理 nil 错误
	_, ok := IsAppError(nil)
	if ok {
		t.Error("Expected IsAppError to return false for nil error")
	}

	// 测试 GetErrorCode 函数处理 nil 错误
	code := GetErrorCode(nil)
	if code != ErrCodeInternalErr {
		t.Errorf("Expected error code %d for nil error, got %d", ErrCodeInternalErr, code)
	}

	// 测试 GetHTTPStatusFromError 函数处理 nil 错误
	status := GetHTTPStatusFromError(nil)
	if status != 500 {
		t.Errorf("Expected HTTP status 500 for nil error, got %d", status)
	}
}
