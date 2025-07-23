// Package config 提供向后兼容的类型和函数
// 为了保持与旧代码的兼容性，重新导出主要的类型和函数
// @Author lanpang
// @Date 2024/12/19
// @Desc 向后兼容层，重新导出新的配置类型
package config

import (
	"time"
	"vhagar/database"
	"vhagar/utils"
)

// 重新导出数据库配置类型（向后兼容）
type DB = database.Config
type RedisConfig = database.RedisConfig

// 重新导出旧的配置类型（向后兼容）
type TenantCfg = TenantConfig
type CorpCfg = CorpConfig

// 重新导出任务相关配置类型（向后兼容）
type DorisCfg = DorisConfig
type RocketMQCfg = RocketMQConfig
type NacosCfg = NacosConfig
type MetricCfg = MetricConfig
type NotifierCfg = NotifierConfig

// GetRandomDuration 获取随机持续时间（向后兼容）
// 原位于 config/task.go 中，现已移至 utils/random.go
func GetRandomDuration() time.Duration {
	// 调用 utils 包中的实现
	return utils.GetRandomDuration()
}

// 向后兼容的配置访问函数
func GetConfig() *AppConfig {
	return Config
}

// 向后兼容的配置初始化函数
func LoadConfig(configPath string) (*AppConfig, error) {
	return loadConfig(configPath)
}
