// Package errors 错误日志记录功能
// 提供错误日志记录和格式化功能
package errors

import (
	"fmt"
	"vhagar/logger"
)

// LogError 记录错误日志
// 参数:
//   - err: 要记录的错误
//   - context: 错误上下文信息
func LogError(err error, context string) {
	if err == nil {
		return
	}

	if logger.Logger == nil {
		return
	}

	if appErr, ok := err.(*AppError); ok {
		if appErr.Cause != nil {
			logger.Logger.Errorw("应用错误",
				"context", context,
				"code", appErr.Code,
				"message", appErr.Message,
				"detail", appErr.Detail,
				"cause", appErr.Cause.Error(),
			)
		} else {
			logger.Logger.Errorw("应用错误",
				"context", context,
				"code", appErr.Code,
				"message", appErr.Message,
				"detail", appErr.Detail,
				"cause", nil,
			)
		}
	} else {
		logger.Logger.Errorw("系统错误",
			"context", context,
			"error", err.Error(),
		)
	}
}

// LogErrorWithFields 记录带字段的错误日志
// 参数:
//   - err: 要记录的错误
//   - context: 错误上下文信息
//   - fields: 额外的字段信息
func LogErrorWithFields(err error, context string, fields map[string]any) {
	if err == nil {
		return
	}

	if logger.Logger == nil {
		return
	}

	// 合并字段
	allFields := make(map[string]any)
	if fields != nil {
		for k, v := range fields {
			allFields[k] = v
		}
	}

	if appErr, ok := err.(*AppError); ok {
		allFields["context"] = context
		allFields["code"] = appErr.Code
		allFields["message"] = appErr.Message
		allFields["detail"] = appErr.Detail
		if appErr.Cause != nil {
			allFields["cause"] = appErr.Cause.Error()
		} else {
			allFields["cause"] = nil
		}
		logger.Logger.Errorw("应用错误", fieldsToArgs(allFields)...)
	} else {
		allFields["context"] = context
		allFields["error"] = err.Error()
		logger.Logger.Errorw("系统错误", fieldsToArgs(allFields)...)
	}
}

// LogErrorf 记录格式化的错误日志
// 参数:
//   - err: 要记录的错误
//   - format: 格式化字符串
//   - args: 格式化参数
func LogErrorf(err error, format string, args ...any) {
	if err == nil {
		return
	}

	if logger.Logger == nil {
		return
	}

	message := fmt.Sprintf(format, args...)
	LogError(err, message)
}

// LogDebug 记录调试日志
// 参数:
//   - msg: 日志消息
//   - keysAndValues: 键值对参数
func LogDebug(msg string, keysAndValues ...any) {
	if logger.Logger == nil {
		return
	}
	logger.Logger.Debugw(msg, keysAndValues...)
}

// LogInfo 记录信息日志
// 参数:
//   - msg: 日志消息
//   - keysAndValues: 键值对参数
func LogInfo(msg string, keysAndValues ...any) {
	if logger.Logger == nil {
		return
	}
	logger.Logger.Infow(msg, keysAndValues...)
}

// LogWarn 记录警告日志
// 参数:
//   - msg: 日志消息
//   - keysAndValues: 键值对参数
func LogWarn(msg string, keysAndValues ...any) {
	if logger.Logger == nil {
		return
	}
	logger.Logger.Warnw(msg, keysAndValues...)
}

// fieldsToArgs 将字段映射转换为参数列表
// 参数:
//   - fields: 字段映射
//
// 返回:
//   - []any: 参数列表
func fieldsToArgs(fields map[string]any) []any {
	args := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return args
}
