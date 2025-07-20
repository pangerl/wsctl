// Package database 提供数据库连接工具和配置管理
// 包含 MySQL、PostgreSQL、Redis、Elasticsearch 等数据库的连接配置和工具函数
package database

// Config 数据库连接配置结构体
// 用于统一管理各种数据库的连接参数
type Config struct {
	Host     string `toml:"host"`     // 数据库主机地址（重命名自 Ip，符合 Go 命名规范）
	Port     int    `toml:"port"`     // 数据库端口号
	Username string `toml:"username"` // 数据库用户名
	Password string `toml:"password"` // 数据库密码
	Database string `toml:"database"` // 数据库名称
	SSLMode  bool   `toml:"ssl_mode"` // SSL 模式（重命名自 Sslmode，符合 Go 命名规范）
}

// HasValue 检查数据库配置是否包含所有必需的值
// 返回 true 表示配置完整，false 表示缺少必要参数
func (c Config) HasValue() bool {
	return c.Host != "" && c.Port != 0 && c.Username != "" && c.Password != ""
}

// RedisConfig Redis 连接配置结构体
// 专门用于 Redis 数据库的连接配置
type RedisConfig struct {
	Addr     string `toml:"addr"`     // Redis 服务器地址
	Password string `toml:"password"` // Redis 密码
	DB       int    `toml:"db"`       // Redis 数据库编号
}

// HasValue 检查 Redis 配置是否包含必需的值
// 返回 true 表示配置完整，false 表示缺少必要参数
func (r RedisConfig) HasValue() bool {
	return r.Addr != ""
}
