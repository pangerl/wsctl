// Package logger 提供应用程序的日志工具和配置
// 支持结构化日志记录，具有不同级别和输出格式，包括控制台和文件输出，支持日志轮转功能
package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config 日志器配置结构体
type Config struct {
	Level      string `toml:"level" json:"level"`             // 日志级别 (debug, info, warn, error)
	ToFile     bool   `toml:"to_file" json:"to_file"`         // 是否写入文件
	FilePath   string `toml:"file_path" json:"file_path"`     // 日志文件路径
	MaxSize    int    `toml:"max_size" json:"max_size"`       // 最大文件大小（MB）
	MaxBackups int    `toml:"max_backups" json:"max_backups"` // 最大备份文件数
	MaxAge     int    `toml:"max_age" json:"max_age"`         // 保留旧日志文件的最大天数
	Compress   bool   `toml:"compress" json:"compress"`       // 是否压缩旧日志文件
	Format     string `toml:"format" json:"format"`           // 日志格式 (json, console)
}

// DefaultConfig 返回默认的日志器配置
func DefaultConfig() Config {
	return Config{
		Level:      "info",
		ToFile:     false,
		FilePath:   "logs/vhagar.log",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     7,
		Compress:   true,
		Format:     "console",
	}
}

var (
	// Logger 全局日志器实例
	Logger *zap.SugaredLogger
	// once 确保日志只初始化一次
	once sync.Once
)

// LoggerInterface 定义日志器实现的接口
type LoggerInterface interface {
	Debug(args ...any)
	Debugf(template string, args ...any)
	Debugw(msg string, keysAndValues ...any)
	Info(args ...any)
	Infof(template string, args ...any)
	Infow(msg string, keysAndValues ...any)
	Warn(args ...any)
	Warnf(template string, args ...any)
	Warnw(msg string, keysAndValues ...any)
	Error(args ...any)
	Errorf(template string, args ...any)
	Errorw(msg string, keysAndValues ...any)
	Fatal(args ...any)
	Fatalf(template string, args ...any)
	Fatalw(msg string, keysAndValues ...any)
	Sync() error
}

// InitLogger 使用配置结构体初始化日志器
// 参数:
//   - cfg: 日志配置信息
//
// 返回:
//   - error: 初始化错误信息
func InitLogger(cfg Config) error {
	var initErr error
	once.Do(func() {
		initErr = initLoggerInternal(cfg)
	})
	return initErr
}

// InitLoggerWithConfig 支持自定义日志级别和是否写文件（向后兼容）
// 参数:
//   - level: 日志级别字符串
//   - toFile: 是否输出到文件
func InitLoggerWithConfig(level string, toFile bool) {
	cfg := Config{
		Level:      level,
		ToFile:     toFile,
		FilePath:   "logs/vhagar.log",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     7,
		Compress:   true,
		Format:     "console",
	}
	_ = InitLogger(cfg)
}

// initLoggerInternal 内部日志初始化函数
func initLoggerInternal(cfg Config) error {
	// 解析日志级别
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("无效的日志级别 %s: %w", cfg.Level, err)
	}

	// 设置默认值
	if cfg.FilePath == "" {
		cfg.FilePath = "logs/vhagar.log"
	}
	if cfg.MaxSize == 0 {
		cfg.MaxSize = 100
	}
	if cfg.MaxBackups == 0 {
		cfg.MaxBackups = 10
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = 7
	}
	if cfg.Format == "" {
		cfg.Format = "console"
	}

	var core zapcore.Core

	if cfg.ToFile {
		core = createFileCore(cfg, level)
	} else {
		core = createConsoleCore(cfg, level)
	}

	// 创建日志器
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	Logger = logger.Sugar()

	return nil
}

// createConsoleCore 为日志器创建控制台输出核心
func createConsoleCore(cfg Config, level zapcore.Level) zapcore.Core {
	encoderConfig := getEncoderConfig(cfg.Format)
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
}

// createFileCore 为日志器创建带轮转的文件输出核心
func createFileCore(cfg Config, level zapcore.Level) zapcore.Core {
	// 确保日志目录存在
	if err := os.MkdirAll(filepath.Dir(cfg.FilePath), 0755); err != nil {
		// 如果文件创建失败，回退到控制台输出
		return createConsoleCore(cfg, level)
	}

	writer := &lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	encoderConfig := getEncoderConfig(cfg.Format)
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	if cfg.Format == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	return zapcore.NewCore(encoder, zapcore.AddSync(writer), level)
}

// getEncoderConfig 返回日志器的编码器配置
func getEncoderConfig(format string) zapcore.EncoderConfig {
	config := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if format == "console" {
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}
	}

	return config
}

// GetLogger 返回全局日志器实例
// 如果日志器未初始化，将使用默认配置初始化
func GetLogger() LoggerInterface {
	if Logger == nil {
		// 使用默认配置初始化
		defaultConfig := DefaultConfig()
		_ = InitLogger(defaultConfig)
	}

	if Logger == nil {
		// 如果仍然未初始化，返回空操作日志器
		return &noOpLogger{}
	}

	return Logger
}

// noOpLogger 空操作日志器实现
type noOpLogger struct{}

func (n *noOpLogger) Debug(args ...any)                       {}
func (n *noOpLogger) Debugf(template string, args ...any)     {}
func (n *noOpLogger) Debugw(msg string, keysAndValues ...any) {}
func (n *noOpLogger) Info(args ...any)                        {}
func (n *noOpLogger) Infof(template string, args ...any)      {}
func (n *noOpLogger) Infow(msg string, keysAndValues ...any)  {}
func (n *noOpLogger) Warn(args ...any)                        {}
func (n *noOpLogger) Warnf(template string, args ...any)      {}
func (n *noOpLogger) Warnw(msg string, keysAndValues ...any)  {}
func (n *noOpLogger) Error(args ...any)                       {}
func (n *noOpLogger) Errorf(template string, args ...any)     {}
func (n *noOpLogger) Errorw(msg string, keysAndValues ...any) {}
func (n *noOpLogger) Fatal(args ...any)                       {}
func (n *noOpLogger) Fatalf(template string, args ...any)     {}
func (n *noOpLogger) Fatalw(msg string, keysAndValues ...any) {}
func (n *noOpLogger) Sync() error                             { return nil }
