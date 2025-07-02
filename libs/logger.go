package libs

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.SugaredLogger
	once   sync.Once
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
			fileWriter, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}
			fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
			core = zapcore.NewTee(
				zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapLevel),
				zapcore.NewCore(fileEncoder, zapcore.AddSync(fileWriter), zapLevel),
			)
		} else {
			core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapLevel)
		}
		logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		Logger = logger.Sugar()
	})
}
