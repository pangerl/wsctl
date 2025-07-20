// Package errors 错误日志记录功能
package errors

import (
	"vhagar/logger"
)

// LogError 记录错误日志
// 参数:
//   - err: 要记录的错误
//   - context: 错误上下文信息
func LogError(err error, context string) {
	// 获取日志器实例
	log := logger.GetLogger()

	if appErr, ok := IsAppError(err); ok {
		// 记录应用错误
		log.Errorw("应用错误",
			"context", context,
			"code", appErr.Code,
			"message", appErr.Message,
			"detail", appErr.Detail,
			"cause", appErr.Cause,
		)
	} else {
		// 记录系统错误
		log.Errorw("系统错误",
			"context", context,
			"error", err.Error(),
		)
	}
}

// LogErrorWithFields 记录带字段的错误日志
// 参数:
//   - err: 要记录的错误
//   - context: 错误上下文信息
//   - fields: 额外的日志字段
func LogErrorWithFields(err error, context string, fields map[string]any) {
	// 获取日志器实例
	log := logger.GetLogger()

	// 构建日志字段
	logFields := []any{"context", context}

	if appErr, ok := IsAppError(err); ok {
		// 添加应用错误字段
		logFields = append(logFields,
			"code", appErr.Code,
			"message", appErr.Message,
			"detail", appErr.Detail,
			"cause", appErr.Cause,
		)
	} else {
		// 添加系统错误字段
		logFields = append(logFields, "error", err.Error())
	}

	// 添加自定义字段
	for k, v := range fields {
		logFields = append(logFields, k, v)
	}

	log.Errorw("错误详情", logFields...)
}

// LogErrorf 记录格式化的错误日志
// 参数:
//   - err: 要记录的错误
//   - format: 格式化字符串
//   - args: 格式化参数
func LogErrorf(err error, format string, args ...any) {
	// 获取日志器实例
	log := logger.GetLogger()

	if appErr, ok := IsAppError(err); ok {
		// 记录应用错误
		log.Errorw("应用错误",
			"message", format,
			"args", args,
			"code", appErr.Code,
			"error_message", appErr.Message,
			"detail", appErr.Detail,
			"cause", appErr.Cause,
		)
	} else {
		// 记录系统错误
		log.Errorw("系统错误",
			"message", format,
			"args", args,
			"error", err.Error(),
		)
	}
}

// LogWarn 记录警告日志
// 参数:
//   - message: 警告消息
//   - fields: 日志字段
func LogWarn(message string, fields ...any) {
	log := logger.GetLogger()
	log.Warnw(message, fields...)
}

// LogInfo 记录信息日志
// 参数:
//   - message: 信息消息
//   - fields: 日志字段
func LogInfo(message string, fields ...any) {
	log := logger.GetLogger()
	log.Infow(message, fields...)
}

// LogDebug 记录调试日志
// 参数:
//   - message: 调试消息
//   - fields: 日志字段
func LogDebug(message string, fields ...any) {
	log := logger.GetLogger()
	log.Debugw(message, fields...)
}
