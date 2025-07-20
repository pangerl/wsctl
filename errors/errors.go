// Package errors 统一错误处理模块
// 提供统一的错误处理和错误码定义，支持错误包装和 HTTP 状态码映射
package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误码类型
// 用于标识不同类型的应用错误
type ErrorCode int

// 定义错误码常量
const (
	// 通用错误码 (0-19999)
	ErrCodeSuccess      ErrorCode = 0     // 成功
	ErrCodeInternalErr  ErrorCode = 10001 // 内部错误
	ErrCodeInvalidParam ErrorCode = 10002 // 参数错误
	ErrCodeNotFound     ErrorCode = 10003 // 资源不存在
	ErrCodeUnauthorized ErrorCode = 10004 // 未授权
	ErrCodeForbidden    ErrorCode = 10005 // 禁止访问

	// AI 相关错误码 (20000-29999)
	ErrCodeAIProviderNotFound ErrorCode = 20001 // AI 服务商未找到
	ErrCodeAIRequestFailed    ErrorCode = 20002 // AI 请求失败
	ErrCodeAIResponseInvalid  ErrorCode = 20003 // AI 响应无效

	// 工具相关错误码 (30000-39999)
	ErrCodeToolNotFound   ErrorCode = 30001 // 工具未找到
	ErrCodeToolCallFailed ErrorCode = 30002 // 工具调用失败
	ErrCodeToolRegFailed  ErrorCode = 30003 // 工具注册失败

	// 配置相关错误码 (40000-49999)
	ErrCodeConfigNotFound ErrorCode = 40001 // 配置文件未找到
	ErrCodeConfigInvalid  ErrorCode = 40002 // 配置文件格式错误

	// 网络相关错误码 (50000-59999)
	ErrCodeNetworkTimeout ErrorCode = 50001 // 网络超时
	ErrCodeNetworkFailed  ErrorCode = 50002 // 网络请求失败

	// 数据库相关错误码 (60000-69999)
	ErrCodeDBConnFailed   ErrorCode = 60001 // 数据库连接失败
	ErrCodeDBQueryFailed  ErrorCode = 60002 // 数据库查询失败
	ErrCodeDBUpdateFailed ErrorCode = 60003 // 数据库更新失败
	ErrCodeDBDeleteFailed ErrorCode = 60004 // 数据库删除失败
	ErrCodeDBInsertFailed ErrorCode = 60005 // 数据库插入失败
	ErrCodeDBTxFailed     ErrorCode = 60006 // 数据库事务失败
)

// AppError 应用错误结构体
// 包含错误码、消息、详细信息和原始错误
type AppError struct {
	Code    ErrorCode `json:"code"`             // 错误码
	Message string    `json:"message"`          // 错误消息
	Detail  string    `json:"detail,omitempty"` // 详细信息
	Cause   error     `json:"-"`                // 原始错误（不序列化）
}

// Error 实现 error 接口
// 返回格式化的错误字符串
func (e *AppError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 支持 Go 1.13+ 的错误链功能
// 返回包装的原始错误
func (e *AppError) Unwrap() error {
	return e.Cause
}

// GetHTTPStatus 根据错误码获取对应的 HTTP 状态码
// 返回:
//   - int: HTTP 状态码
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
	case ErrCodeDBConnFailed, ErrCodeDBQueryFailed, ErrCodeDBUpdateFailed,
		ErrCodeDBDeleteFailed, ErrCodeDBInsertFailed, ErrCodeDBTxFailed:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// New 创建新的应用错误
// 参数:
//   - code: 错误码
//   - message: 错误消息
//
// 返回:
//   - *AppError: 应用错误对象
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// NewError 创建新的应用错误（向后兼容）
// 参数:
//   - code: 错误码
//   - message: 错误消息
//
// 返回:
//   - *AppError: 应用错误对象
func NewError(code ErrorCode, message string) *AppError {
	return New(code, message)
}

// NewWithDetail 创建带详细信息的应用错误
// 参数:
//   - code: 错误码
//   - message: 错误消息
//   - detail: 详细信息
//
// 返回:
//   - *AppError: 应用错误对象
func NewWithDetail(code ErrorCode, message, detail string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

// NewErrorWithDetail 创建带详细信息的应用错误（向后兼容）
// 参数:
//   - code: 错误码
//   - message: 错误消息
//   - detail: 详细信息
//
// 返回:
//   - *AppError: 应用错误对象
func NewErrorWithDetail(code ErrorCode, message, detail string) *AppError {
	return NewWithDetail(code, message, detail)
}

// Wrap 包装已有错误
// 参数:
//   - code: 错误码
//   - message: 错误消息
//   - cause: 原始错误
//
// 返回:
//   - *AppError: 应用错误对象
func Wrap(code ErrorCode, message string, cause error) *AppError {
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

// WrapError 包装已有错误（向后兼容）
// 参数:
//   - code: 错误码
//   - message: 错误消息
//   - cause: 原始错误
//
// 返回:
//   - *AppError: 应用错误对象
func WrapError(code ErrorCode, message string, cause error) *AppError {
	return Wrap(code, message, cause)
}

// WrapWithDetail 包装已有错误并添加详细信息
// 参数:
//   - code: 错误码
//   - message: 错误消息
//   - detail: 详细信息
//   - cause: 原始错误
//
// 返回:
//   - *AppError: 应用错误对象
func WrapWithDetail(code ErrorCode, message, detail string, cause error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Detail:  detail,
		Cause:   cause,
	}
}

// 预定义常用错误实例
var (
	// 通用错误
	ErrInternalServer = New(ErrCodeInternalErr, "内部服务器错误")
	ErrInvalidParam   = New(ErrCodeInvalidParam, "参数错误")
	ErrNotFound       = New(ErrCodeNotFound, "资源不存在")
	ErrUnauthorized   = New(ErrCodeUnauthorized, "未授权访问")
	ErrForbidden      = New(ErrCodeForbidden, "禁止访问")

	// AI 相关错误
	ErrAIProviderNotFound = New(ErrCodeAIProviderNotFound, "AI 服务商配置未找到")
	ErrAIRequestFailed    = New(ErrCodeAIRequestFailed, "AI 请求失败")
	ErrAIResponseInvalid  = New(ErrCodeAIResponseInvalid, "AI 响应格式无效")

	// 工具相关错误
	ErrToolNotFound   = New(ErrCodeToolNotFound, "工具未找到")
	ErrToolCallFailed = New(ErrCodeToolCallFailed, "工具调用失败")
	ErrToolRegFailed  = New(ErrCodeToolRegFailed, "工具注册失败")

	// 配置相关错误
	ErrConfigNotFound = New(ErrCodeConfigNotFound, "配置文件未找到")
	ErrConfigInvalid  = New(ErrCodeConfigInvalid, "配置文件格式错误")

	// 网络相关错误
	ErrNetworkTimeout = New(ErrCodeNetworkTimeout, "网络请求超时")
	ErrNetworkFailed  = New(ErrCodeNetworkFailed, "网络请求失败")

	// 数据库相关错误
	ErrDBConnFailed   = New(ErrCodeDBConnFailed, "数据库连接失败")
	ErrDBQueryFailed  = New(ErrCodeDBQueryFailed, "数据库查询失败")
	ErrDBUpdateFailed = New(ErrCodeDBUpdateFailed, "数据库更新失败")
	ErrDBDeleteFailed = New(ErrCodeDBDeleteFailed, "数据库删除失败")
	ErrDBInsertFailed = New(ErrCodeDBInsertFailed, "数据库插入失败")
	ErrDBTxFailed     = New(ErrCodeDBTxFailed, "数据库事务失败")
)

// IsAppError 检查错误是否为应用错误
// 参数:
//   - err: 要检查的错误
//
// 返回:
//   - *AppError: 如果是应用错误则返回转换后的对象，否则返回 nil
//   - bool: 是否为应用错误
func IsAppError(err error) (*AppError, bool) {
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}

// GetErrorCode 获取错误的错误码
// 参数:
//   - err: 要检查的错误
//
// 返回:
//   - ErrorCode: 错误码，如果不是应用错误则返回 ErrCodeInternalErr
func GetErrorCode(err error) ErrorCode {
	if appErr, ok := IsAppError(err); ok {
		return appErr.Code
	}
	return ErrCodeInternalErr
}

// GetHTTPStatusFromError 从错误获取 HTTP 状态码
// 参数:
//   - err: 要检查的错误
//
// 返回:
//   - int: HTTP 状态码
func GetHTTPStatusFromError(err error) int {
	if appErr, ok := IsAppError(err); ok {
		return appErr.GetHTTPStatus()
	}
	return http.StatusInternalServerError
}
