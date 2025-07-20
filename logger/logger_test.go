package logger

import (
	"os"
	"sync"
	"testing"

	"go.uber.org/zap/zapcore"
)

// TestDefaultConfig 测试默认配置
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Level != "info" {
		t.Errorf("Expected default level 'info', got '%s'", cfg.Level)
	}

	if cfg.ToFile != false {
		t.Errorf("Expected default ToFile false, got %v", cfg.ToFile)
	}

	if cfg.FilePath != "logs/vhagar.log" {
		t.Errorf("Expected default FilePath 'logs/vhagar.log', got '%s'", cfg.FilePath)
	}

	if cfg.MaxSize != 100 {
		t.Errorf("Expected default MaxSize 100, got %d", cfg.MaxSize)
	}

	if cfg.MaxBackups != 10 {
		t.Errorf("Expected default MaxBackups 10, got %d", cfg.MaxBackups)
	}

	if cfg.MaxAge != 7 {
		t.Errorf("Expected default MaxAge 7, got %d", cfg.MaxAge)
	}

	if cfg.Compress != true {
		t.Errorf("Expected default Compress true, got %v", cfg.Compress)
	}

	if cfg.Format != "console" {
		t.Errorf("Expected default Format 'console', got '%s'", cfg.Format)
	}
}

// TestInitLogger 测试日志器初始化
func TestInitLogger(t *testing.T) {
	// 清理测试环境
	defer func() {
		Logger = nil
		once = sync.Once{}
		// 清理测试日志文件
		_ = os.RemoveAll("test_logs")
	}()

	tests := []struct {
		name   string
		config Config
		valid  bool
	}{
		{
			name: "基本控制台日志配置",
			config: Config{
				Level:  "info",
				ToFile: false,
				Format: "console",
			},
			valid: true,
		},
		{
			name: "文件日志配置",
			config: Config{
				Level:      "debug",
				ToFile:     true,
				FilePath:   "test_logs/test.log",
				MaxSize:    50,
				MaxBackups: 5,
				MaxAge:     3,
				Compress:   true,
				Format:     "json",
			},
			valid: true,
		},
		{
			name: "无效日志级别",
			config: Config{
				Level:  "invalid",
				ToFile: false,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重置全局变量
			Logger = nil
			once = sync.Once{}

			err := InitLogger(tt.config)

			if tt.valid {
				if err != nil {
					t.Errorf("InitLogger() error = %v, want nil", err)
				}
				if Logger == nil {
					t.Error("InitLogger() did not set global Logger")
				}
			} else {
				if err == nil {
					t.Error("InitLogger() expected error, got nil")
				}
			}
		})
	}
}

// TestInitLoggerWithConfig 测试向后兼容的日志初始化函数
func TestInitLoggerWithConfig(t *testing.T) {
	// 清理测试环境
	defer func() {
		Logger = nil
		once = sync.Once{}
		_ = os.RemoveAll("logs")
	}()

	tests := []struct {
		name   string
		level  string
		toFile bool
	}{
		{
			name:   "控制台输出_info级别",
			level:  "info",
			toFile: false,
		},
		{
			name:   "文件输出_debug级别",
			level:  "debug",
			toFile: true,
		},
		{
			name:   "控制台输出_error级别",
			level:  "error",
			toFile: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重置全局变量
			Logger = nil
			once = sync.Once{}

			InitLoggerWithConfig(tt.level, tt.toFile)

			if Logger == nil {
				t.Error("InitLoggerWithConfig() did not set global Logger")
			}
		})
	}
}

// TestGetLogger 测试获取日志器实例功能
func TestGetLogger(t *testing.T) {
	// 清理测试环境
	defer func() {
		Logger = nil
		once = sync.Once{}
	}()

	// 测试未初始化时的默认行为
	Logger = nil
	once = sync.Once{}

	logger := GetLogger()
	if logger == nil {
		t.Error("GetLogger() returned nil, expected default logger")
	}

	// 测试已初始化时的行为
	if Logger == nil {
		t.Error("GetLogger() did not initialize global Logger")
	}
}

// TestGetEncoderConfig 测试编码器配置
func TestGetEncoderConfig(t *testing.T) {
	tests := []struct {
		name   string
		format string
	}{
		{
			name:   "控制台格式",
			format: "console",
		},
		{
			name:   "JSON格式",
			format: "json",
		},
		{
			name:   "默认格式",
			format: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := getEncoderConfig(tt.format)

			// 验证基本字段存在
			if config.TimeKey == "" {
				t.Error("Expected TimeKey to be set")
			}
			if config.LevelKey == "" {
				t.Error("Expected LevelKey to be set")
			}
			if config.MessageKey == "" {
				t.Error("Expected MessageKey to be set")
			}
		})
	}
}

// TestCreateConsoleCore 测试控制台核心创建
func TestCreateConsoleCore(t *testing.T) {
	cfg := Config{
		Format: "console",
	}
	level := zapcore.InfoLevel

	core := createConsoleCore(cfg, level)
	if core == nil {
		t.Error("createConsoleCore() returned nil")
	}
}

// TestNoOpLogger 测试空操作日志器
func TestNoOpLogger(t *testing.T) {
	logger := &noOpLogger{}

	// 测试所有方法都不会 panic
	logger.Debug("test")
	logger.Debugf("test %s", "arg")
	logger.Debugw("test", "key", "value")
	logger.Info("test")
	logger.Infof("test %s", "arg")
	logger.Infow("test", "key", "value")
	logger.Warn("test")
	logger.Warnf("test %s", "arg")
	logger.Warnw("test", "key", "value")
	logger.Error("test")
	logger.Errorf("test %s", "arg")
	logger.Errorw("test", "key", "value")
	logger.Fatal("test")
	logger.Fatalf("test %s", "arg")
	logger.Fatalw("test", "key", "value")

	err := logger.Sync()
	if err != nil {
		t.Errorf("noOpLogger.Sync() returned error: %v", err)
	}
}

// TestLoggerInterface 测试日志器接口实现
func TestLoggerInterface(t *testing.T) {
	// 清理测试环境
	defer func() {
		Logger = nil
		once = sync.Once{}
	}()

	// 初始化日志器
	cfg := Config{
		Level:  "debug",
		ToFile: false,
		Format: "console",
	}

	Logger = nil
	once = sync.Once{}

	err := InitLogger(cfg)
	if err != nil {
		t.Fatalf("InitLogger() failed: %v", err)
	}

	// 测试日志器接口方法
	logger := GetLogger()

	// 这些调用不应该 panic
	logger.Debug("debug message")
	logger.Debugf("debug message: %s", "formatted")
	logger.Debugw("debug message", "key", "value")

	logger.Info("info message")
	logger.Infof("info message: %s", "formatted")
	logger.Infow("info message", "key", "value")

	logger.Warn("warn message")
	logger.Warnf("warn message: %s", "formatted")
	logger.Warnw("warn message", "key", "value")

	logger.Error("error message")
	logger.Errorf("error message: %s", "formatted")
	logger.Errorw("error message", "key", "value")

	// 测试 Sync 方法
	err = logger.Sync()
	if err != nil {
		t.Logf("Logger.Sync() returned error (this may be expected): %v", err)
	}
}
