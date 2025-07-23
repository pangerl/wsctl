/*
Package config 提供应用程序的配置管理功能。

本包负责配置文件的加载、解析、验证和访问，支持 TOML 格式的配置文件。
配置结构采用分层设计，将不同功能的配置项组织到对应的子结构中，便于管理和访问。

主要功能：

  - 配置文件加载与解析
  - 配置项验证
  - 配置访问接口
  - 默认配置提供

基本用法：

	// 加载配置文件
	cfg, err := config.InitConfig("config.toml")
	if err != nil {
		// 处理错误
	}

	// 访问配置项
	logLevel := config.Config.Global.LogLevel
	dbConfig := config.Config.Database.PG

	// 检查功能是否启用
	if config.Config.IsAIEnabled() {
		// 使用AI功能
	}

配置结构：

AppConfig 是主配置结构体，包含以下主要部分：
  - Global: 全局应用设置（日志级别、项目名称等）
  - Database: 数据库配置（PostgreSQL、Redis、Elasticsearch等）
  - Services: 外部服务配置（AI、天气API等）
  - Metric: 监控配置
  - Tenant: 租户配置

详细的配置项定义请参考 config.go 文件中的结构体定义。
*/
package config
