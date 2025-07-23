/*
Package logger 提供日志记录功能。

本包基于 zap 日志库，支持不同级别的日志记录、结构化日志和文件输出。
通过统一的日志接口，可以方便地记录应用程序的运行状态和错误信息。

主要功能：

  - 日志初始化与配置
  - 多级别日志记录
  - 结构化日志支持
  - 文件输出与轮转
  - 自定义日志格式

日志级别：

  - Debug: 调试信息，仅在开发环境使用
  - Info: 一般信息，记录正常操作
  - Warn: 警告信息，可能的问题但不影响正常运行
  - Error: 错误信息，影响功能但不导致程序崩溃
  - Fatal: 致命错误，导致程序无法继续运行

基本用法：

	// 初始化日志记录器
	logger.InitLogger(logger.Config{
		Level:  "info",
		ToFile: true,
		Format: "json",
	})

	// 记录不同级别的日志
	logger.Logger.Debug("调试信息")
	logger.Logger.Info("一般信息")
	logger.Logger.Warn("警告信息")
	logger.Logger.Error("错误信息")

	// 记录带字段的结构化日志
	logger.Logger.Infow("用户登录",
		"user_id", 123,
		"ip", "192.168.1.1",
		"time", time.Now(),
	)

	// 记录带格式的日志
	logger.Logger.Infof("处理了 %d 条记录，耗时 %.2f 秒", count, duration)

配置选项：

  - Level: 日志级别 (debug, info, warn, error)
  - ToFile: 是否输出到文件
  - Format: 日志格式 (json, console)
  - FilePath: 日志文件路径 (可选)

日志文件：

当启用文件输出时，日志文件将按日期命名，存放在 logs 目录下：

	logs/2025-07-21.log

注意事项：

  - 避免在性能敏感的代码中使用 Debug 级别的日志
  - 结构化日志的字段名应使用下划线命名法 (user_id 而非 userId)
  - 错误日志应包含足够的上下文信息，便于问题定位
*/
package logger
