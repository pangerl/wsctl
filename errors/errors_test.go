package errors

import (
	"fmt"
	"net/http"
	"testing"
)

// TestNew 测试创建新的应用错误
func TestNew(t *testing.T) {
	err := New(ErrCodeInvalidParam, "测试错误")

	if err.Code != ErrCodeInvalidParam {
		t.Errorf("Expected code %d, got %d", ErrCodeInvalidParam, err.Code)
	}

	if err.Message != "测试错误" {
		t.Errorf("Expected message '测试错误', got '%s'", err.Message)
	}

	if err.Detail != "" {
		t.Errorf("Expected empty detail, got '%s'", err.Detail)
	}

	if err.Cause != nil {
		t.Errorf("Expected nil cause, got %v", err.Cause)
	}
}

// TestNewWithDetail 测试创建带详细信息的应用错误
func TestNewWithDetail(t *testing.T) {
	err := NewWithDetail(ErrCodeInvalidParam, "测试错误", "详细信息")

	if err.Code != ErrCodeInvalidParam {
		t.Errorf("Expected code %d, got %d", ErrCodeInvalidParam, err.Code)
	}

	if err.Message != "测试错误" {
		t.Errorf("Expected message '测试错误', got '%s'", err.Message)
	}

	if err.Detail != "详细信息" {
		t.Errorf("Expected detail '详细信息', got '%s'", err.Detail)
	}
}

// TestWrap 测试包装已有错误
func TestWrap(t *testing.T) {
	originalErr := fmt.Errorf("原始错误")
	wrappedErr := Wrap(ErrCodeInternalErr, "包装错误", originalErr)

	if wrappedErr.Code != ErrCodeInternalErr {
		t.Errorf("Expected code %d, got %d", ErrCodeInternalErr, wrappedErr.Code)
	}

	if wrappedErr.Message != "包装错误" {
		t.Errorf("Expected message '包装错误', got '%s'", wrappedErr.Message)
	}

	if wrappedErr.Cause != originalErr {
		t.Errorf("Expected cause to be original error, got %v", wrappedErr.Cause)
	}

	if wrappedErr.Detail != "原始错误" {
		t.Errorf("Expected detail to be '原始错误', got '%s'", wrappedErr.Detail)
	}
}

// TestAppError_Error 测试错误字符串格式化
func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *AppError
		expected string
	}{
		{
			name: "只有消息",
			err: &AppError{
				Code:    ErrCodeInvalidParam,
				Message: "参数错误",
			},
			expected: "[10002] 参数错误",
		},
		{
			name: "有消息和详细信息",
			err: &AppError{
				Code:    ErrCodeInvalidParam,
				Message: "参数错误",
				Detail:  "用户名不能为空",
			},
			expected: "[10002] 参数错误: 用户名不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestAppError_Unwrap 测试错误解包功能
func TestAppError_Unwrap(t *testing.T) {
	originalErr := fmt.Errorf("原始错误")
	wrappedErr := Wrap(ErrCodeInternalErr, "包装错误", originalErr)

	unwrapped := wrappedErr.Unwrap()
	if unwrapped != originalErr {
		t.Errorf("Expected unwrapped error to be original error, got %v", unwrapped)
	}

	// 测试没有原始错误的情况
	simpleErr := New(ErrCodeInvalidParam, "简单错误")
	unwrapped = simpleErr.Unwrap()
	if unwrapped != nil {
		t.Errorf("Expected unwrapped error to be nil, got %v", unwrapped)
	}
}

// TestAppError_GetHTTPStatus 测试 HTTP 状态码映射
func TestAppError_GetHTTPStatus(t *testing.T) {
	tests := []struct {
		name     string
		code     ErrorCode
		expected int
	}{
		{
			name:     "成功",
			code:     ErrCodeSuccess,
			expected: http.StatusOK,
		},
		{
			name:     "参数错误",
			code:     ErrCodeInvalidParam,
			expected: http.StatusBadRequest,
		},
		{
			name:     "未授权",
			code:     ErrCodeUnauthorized,
			expected: http.StatusUnauthorized,
		},
		{
			name:     "禁止访问",
			code:     ErrCodeForbidden,
			expected: http.StatusForbidden,
		},
		{
			name:     "资源不存在",
			code:     ErrCodeNotFound,
			expected: http.StatusNotFound,
		},
		{
			name:     "网络超时",
			code:     ErrCodeNetworkTimeout,
			expected: http.StatusRequestTimeout,
		},
		{
			name:     "数据库连接失败",
			code:     ErrCodeDBConnFailed,
			expected: http.StatusInternalServerError,
		},
		{
			name:     "未知错误码",
			code:     ErrorCode(99999),
			expected: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.code, "测试错误")
			result := err.GetHTTPStatus()
			if result != tt.expected {
				t.Errorf("Expected HTTP status %d, got %d", tt.expected, result)
			}
		})
	}
}

// TestIsAppError 测试应用错误检查功能
func TestIsAppError(t *testing.T) {
	// 测试应用错误
	appErr := New(ErrCodeInvalidParam, "应用错误")
	result, ok := IsAppError(appErr)
	if !ok {
		t.Error("Expected IsAppError to return true for AppError")
	}
	if result != appErr {
		t.Error("Expected IsAppError to return the same AppError instance")
	}

	// 测试普通错误
	normalErr := fmt.Errorf("普通错误")
	result, ok = IsAppError(normalErr)
	if ok {
		t.Error("Expected IsAppError to return false for normal error")
	}
	if result != nil {
		t.Error("Expected IsAppError to return nil for normal error")
	}
}

// TestGetErrorCode 测试获取错误码功能
func TestGetErrorCode(t *testing.T) {
	// 测试应用错误
	appErr := New(ErrCodeInvalidParam, "应用错误")
	code := GetErrorCode(appErr)
	if code != ErrCodeInvalidParam {
		t.Errorf("Expected error code %d, got %d", ErrCodeInvalidParam, code)
	}

	// 测试普通错误
	normalErr := fmt.Errorf("普通错误")
	code = GetErrorCode(normalErr)
	if code != ErrCodeInternalErr {
		t.Errorf("Expected error code %d for normal error, got %d", ErrCodeInternalErr, code)
	}
}

// TestGetHTTPStatusFromError 测试从错误获取 HTTP 状态码
func TestGetHTTPStatusFromError(t *testing.T) {
	// 测试应用错误
	appErr := New(ErrCodeInvalidParam, "应用错误")
	status := GetHTTPStatusFromError(appErr)
	if status != http.StatusBadRequest {
		t.Errorf("Expected HTTP status %d, got %d", http.StatusBadRequest, status)
	}

	// 测试普通错误
	normalErr := fmt.Errorf("普通错误")
	status = GetHTTPStatusFromError(normalErr)
	if status != http.StatusInternalServerError {
		t.Errorf("Expected HTTP status %d for normal error, got %d", http.StatusInternalServerError, status)
	}
}

// TestPredefinedErrors 测试预定义错误实例
func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name string
		err  *AppError
		code ErrorCode
	}{
		{"ErrInternalServer", ErrInternalServer, ErrCodeInternalErr},
		{"ErrInvalidParam", ErrInvalidParam, ErrCodeInvalidParam},
		{"ErrNotFound", ErrNotFound, ErrCodeNotFound},
		{"ErrUnauthorized", ErrUnauthorized, ErrCodeUnauthorized},
		{"ErrForbidden", ErrForbidden, ErrCodeForbidden},
		{"ErrConfigNotFound", ErrConfigNotFound, ErrCodeConfigNotFound},
		{"ErrConfigInvalid", ErrConfigInvalid, ErrCodeConfigInvalid},
		{"ErrNetworkTimeout", ErrNetworkTimeout, ErrCodeNetworkTimeout},
		{"ErrNetworkFailed", ErrNetworkFailed, ErrCodeNetworkFailed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.code {
				t.Errorf("Expected code %d, got %d", tt.code, tt.err.Code)
			}
			if tt.err.Message == "" {
				t.Error("Expected non-empty message")
			}
		})
	}
}

// TestBackwardCompatibility 测试向后兼容性
func TestBackwardCompatibility(t *testing.T) {
	// 测试 NewError 函数（向后兼容）
	err1 := NewError(ErrCodeInvalidParam, "测试错误")
	err2 := New(ErrCodeInvalidParam, "测试错误")

	if err1.Code != err2.Code || err1.Message != err2.Message {
		t.Error("NewError should be compatible with New")
	}

	// 测试 NewErrorWithDetail 函数（向后兼容）
	err3 := NewErrorWithDetail(ErrCodeInvalidParam, "测试错误", "详细信息")
	err4 := NewWithDetail(ErrCodeInvalidParam, "测试错误", "详细信息")

	if err3.Code != err4.Code || err3.Message != err4.Message || err3.Detail != err4.Detail {
		t.Error("NewErrorWithDetail should be compatible with NewWithDetail")
	}

	// 测试 WrapError 函数（向后兼容）
	originalErr := fmt.Errorf("原始错误")
	err5 := WrapError(ErrCodeInternalErr, "包装错误", originalErr)
	err6 := Wrap(ErrCodeInternalErr, "包装错误", originalErr)

	if err5.Code != err6.Code || err5.Message != err6.Message {
		t.Error("WrapError should be compatible with Wrap")
	}
}
