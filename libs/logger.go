package libs

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger *zap.SugaredLogger
	once   sync.Once

	// 占位使用，避免未使用变量警告
	_ = Logger
	_ = once
)

// InitLoggerWithConfig 支持自定义日志级别和是否写文件
func InitLoggerWithConfig(level string, toFile bool) {
	once.Do(func() {
		var zapLevel zapcore.Level
		switch level {
		case "debug":
			zapLevel = zapcore.DebugLevel
		case "info":
			zapLevel = zapcore.InfoLevel
		case "warn":
			zapLevel = zapcore.WarnLevel
		case "error":
			zapLevel = zapcore.ErrorLevel
		default:
			zapLevel = zapcore.InfoLevel
		}

		encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		var core zapcore.Core
		if toFile {
			logFile := "logs/vhagar.log"
			_ = os.MkdirAll("logs", 0755)
			fileWriter := &lumberjack.Logger{
				Filename:   logFile,
				MaxSize:    100,  // 单个日志文件最大100MB
				MaxBackups: 10,   // 最多保留10个备份文件
				MaxAge:     7,    // 只保留7天
				Compress:   true, // 启用压缩
			}
			fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
			consoleCore := zapcore.NewCore(
				encoder,
				zapcore.AddSync(os.Stdout),
				zapcore.WarnLevel, // 控制台只输出 warn 及以上
			)
			fileCore := zapcore.NewCore(
				fileEncoder,
				zapcore.AddSync(fileWriter),
				zapLevel, // 文件输出 info 及以上
			)
			core = zapcore.NewTee(consoleCore, fileCore)
		} else {
			core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapLevel)
		}
		logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		Logger = logger.Sugar()
	})
}
