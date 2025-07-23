/*
Package database 提供数据库连接和操作的功能。

本包封装了常用数据库的连接配置和客户端创建逻辑，支持 MySQL、PostgreSQL、Redis 和 Elasticsearch。
每种数据库都有对应的配置结构体和客户端创建函数，提供统一的接口进行数据库操作。

主要功能：

  - 数据库配置管理
  - 数据库连接创建
  - 连接池管理
  - 错误处理和日志记录

支持的数据库：

  - MySQL
  - PostgreSQL
  - Redis
  - Elasticsearch

基本用法：

	// MySQL 连接示例
	cfg := database.Config{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "password",
		Database: "mydb",
	}

	client, err := database.NewMySQLClient(cfg, "mydb")
	if err != nil {
		// 处理错误
	}
	defer client.Close()

	// 使用客户端进行数据库操作
	rows, err := client.Query("SELECT * FROM users")

	// Redis 连接示例
	redisCfg := database.RedisConfig{
		Addr:     "localhost:6379",
		Password: "password",
		DB:       0,
	}

	redisClient, err := database.NewRedisClient(redisCfg)
	if err != nil {
		// 处理错误
	}
	defer redisClient.Close()

配置验证：

所有数据库配置结构体都实现了 HasValue 方法，用于验证配置的完整性：

	if !cfg.HasValue() {
		// 配置不完整，处理错误
	}

错误处理：

本包使用 vhagar/errors 包进行统一的错误处理，所有数据库操作错误都会被包装为应用错误，
包含错误码、错误消息和原始错误信息，便于问题定位和处理。
*/
package database
