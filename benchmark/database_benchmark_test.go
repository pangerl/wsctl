package benchmark

import (
	"testing"
	"vhagar/database"
)

// 测试数据库配置验证性能
func BenchmarkDatabaseConfigValidation(b *testing.B) {
	// 创建有效的数据库配置
	cfg := database.Config{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "password",
		Database: "vhagar",
	}

	// 重置计时器
	b.ResetTimer()

	// 执行基准测试
	for i := 0; i < b.N; i++ {
		_ = cfg.HasValue()
	}
}

// 测试Redis配置验证性能
func BenchmarkRedisConfigValidation(b *testing.B) {
	// 创建有效的Redis配置
	cfg := database.RedisConfig{
		Addr:     "localhost:6379",
		Password: "password",
		DB:       0,
	}

	// 重置计时器
	b.ResetTimer()

	// 执行基准测试
	for i := 0; i < b.N; i++ {
		_ = cfg.HasValue()
	}
}
