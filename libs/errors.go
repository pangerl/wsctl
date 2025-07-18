// Package libs 统一错误处理模块
// @Author lanpang
// @Date 2025/7/18
// @Desc 提供统一的错误处理和错误码定义
package libs

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误码类型
type ErrorCode int

// 定义错误码常量
const (
	// 通用错误码
	ErrCodeSuccess      ErrorCode = 0     // 成功
	ErrCodeInternalErr  ErrorCode = 10001 // 内部错误
	ErrCodeInvalidParam ErrorCode = 10002 // 参数错误
	ErrCodeNotFound     ErrorCode = 10003 // 资源不存在
	ErrCodeUnauthorized ErrorCode = 10004 // 未授权
	ErrCodeForbidden    ErrorCode = 10005 // 禁止访问

	// AI相关错误码
	ErrCodeAIProviderNotFound ErrorCode = 20001 // AI服务商未找到
	ErrCodeAIRequestFailed    ErrorCode = 20002 // AI请求失败
	ErrCodeAIResponseInvalid  ErrorCode = 20003 // AI响应无效

	// 工具相关错误码
	ErrCodeToolNotFound   ErrorCode = 30001 // 工具未找到
	ErrCodeToolCallFailed ErrorCode = 30002 // 工具调用失败
	ErrCodeToolRegFailed  ErrorCode = 30003 // 工具注册失败

	// 配置相关错误码
	ErrCodeConfigNotFound ErrorCode = 40001 // 配置文件未找到
	ErrCodeConfigInvalid  ErrorCode = 40002 // 配置文件格式错误

	// 网络相关错误码
	ErrCodeNetworkTimeout ErrorCode = 50001 // 网络超时
	ErrCodeNetworkFailed  ErrorCode = 50002 // 网络请求失败
)

// AppError 应用错误结构
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Detail  string    `json:"detail,omitempty"`
	Cause   error     `json:"-"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 支持errors.Unwrap
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewError 创建新的应用错误
func NewError(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// NewErrorWithDetail 创建带详细信息的应用错误
func NewErrorWithDetail(code ErrorCode, message, detail string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

// WrapError 包装已有错误
func WrapError(code ErrorCode, message string, cause error) *AppError {
	detail := ""
	if cause != nil {
		detail = cause.Error()
	}
	return &AppError{
		Code:    code,
		Message: message,
		Detail:  detail,
		Cause:   cause,
	}
}

// GetHTTPStatus 根据错误码获取HTTP状态码
func (e *AppError) GetHTTPStatus() int {
	switch e.Code {
	case ErrCodeSuccess:
		return http.StatusOK
	case ErrCodeInvalidParam:
		return http.StatusBadRequest
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeNotFound, ErrCodeToolNotFound, ErrCodeConfigNotFound:
		return http.StatusNotFound
	case ErrCodeNetworkTimeout:
		return http.StatusRequestTimeout
	default:
		return http.StatusInternalServerError
	}
}

// 预定义常用错误
var (
	ErrInternalServer = NewError(ErrCodeInternalErr, "内部服务器错误")
	ErrInvalidParam   = NewError(ErrCodeInvalidParam, "参数错误")
	ErrNotFound       = NewError(ErrCodeNotFound, "资源不存在")
	ErrUnauthorized   = NewError(ErrCodeUnauthorized, "未授权访问")
	ErrForbidden      = NewError(ErrCodeForbidden, "禁止访问")

	ErrAIProviderNotFound = NewError(ErrCodeAIProviderNotFound, "AI服务商配置未找到")
	ErrAIRequestFailed    = NewError(ErrCodeAIRequestFailed, "AI请求失败")
	ErrAIResponseInvalid  = NewError(ErrCodeAIResponseInvalid, "AI响应格式无效")

	ErrToolNotFound   = NewError(ErrCodeToolNotFound, "工具未找到")
	ErrToolCallFailed = NewError(ErrCodeToolCallFailed, "工具调用失败")
	ErrToolRegFailed  = NewError(ErrCodeToolRegFailed, "工具注册失败")

	ErrConfigNotFound = NewError(ErrCodeConfigNotFound, "配置文件未找到")
	ErrConfigInvalid  = NewError(ErrCodeConfigInvalid, "配置文件格式错误")

	ErrNetworkTimeout = NewError(ErrCodeNetworkTimeout, "网络请求超时")
	ErrNetworkFailed  = NewError(ErrCodeNetworkFailed, "网络请求失败")

	// 占位使用，避免未使用变量警告
	_ = ErrInternalServer
	_ = ErrInvalidParam
	_ = ErrNotFound
	_ = ErrUnauthorized
	_ = ErrForbidden
	_ = ErrAIProviderNotFound
	_ = ErrAIRequestFailed
	_ = ErrAIResponseInvalid
	_ = ErrToolNotFound
	_ = ErrToolCallFailed
	_ = ErrToolRegFailed
	_ = ErrConfigNotFound
	_ = ErrConfigInvalid
	_ = ErrNetworkTimeout
	_ = ErrNetworkFailed
)

// LogError 记录错误日志
func LogError(err error, context string) {
	if appErr, ok := err.(*AppError); ok {
		Logger.Errorw("应用错误",
			"context", context,
			"code", appErr.Code,
			"message", appErr.Message,
			"detail", appErr.Detail,
			"cause", appErr.Cause,
		)
	} else {
		Logger.Errorw("系统错误",
			"context", context,
			"error", err.Error(),
		)
	}
}

// LogErrorWithFields 记录带字段的错误日志
func LogErrorWithFields(err error, context string, fields map[string]interface{}) {
	logFields := []interface{}{"context", context}

	if appErr, ok := err.(*AppError); ok {
		logFields = append(logFields,
			"code", appErr.Code,
			"message", appErr.Message,
			"detail", appErr.Detail,
			"cause", appErr.Cause,
		)
	} else {
		logFields = append(logFields, "error", err.Error())
	}

	// 添加自定义字段
	for k, v := range fields {
		logFields = append(logFields, k, v)
	}

	Logger.Errorw("错误详情", logFields...)
}
