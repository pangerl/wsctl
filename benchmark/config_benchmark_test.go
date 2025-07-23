package benchmark

import (
	"testing"
	"vhagar/config"
)

// 测试配置访问性能（跳过加载）
func BenchmarkConfigAccess2(b *testing.B) {
	// 创建一个有效的配置对象
	cfg := config.DefaultConfig()
	cfg.Global.ProjectName = "vhagar"
	cfg.Global.LogLevel = "info"
	cfg.Services.AI.Enable = true
	cfg.Services.AI.Provider = "openai"
	cfg.Services.AI.Providers = map[string]config.ProviderConfig{
		"openai": {
			ApiKey: "sk-test",
			ApiUrl: "https://api.openai.com/v1",
			Model:  "gpt-3.5-turbo",
		},
	}
	cfg.Database.PG.Host = "localhost"
	cfg.Database.PG.Port = 5432
	cfg.Database.PG.Username = "postgres"
	cfg.Database.PG.Password = "password"
	cfg.Database.PG.Database = "vhagar"

	// 重置计时器
	b.ResetTimer()

	// 执行基准测试
	for i := 0; i < b.N; i++ {
		// 测试各种配置访问方法
		_ = cfg.IsAIEnabled()
		_, _ = cfg.GetAIProvider()
		_ = cfg.IsMetricEnabled()
		_ = cfg.GetDatabaseConfig("pg")
		_ = cfg.GetDatabaseConfig("redis")
	}
}

// 测试配置验证性能
func BenchmarkConfigValidate(b *testing.B) {
	// 创建一个有效的配置对象
	cfg := config.DefaultConfig()
	cfg.Global.ProjectName = "vhagar"
	cfg.Global.LogLevel = "info"
	cfg.Services.AI.Enable = true
	cfg.Services.AI.Provider = "openai"
	cfg.Services.AI.Providers = map[string]config.ProviderConfig{
		"openai": {
			ApiKey: "sk-test",
			ApiUrl: "https://api.openai.com/v1",
			Model:  "gpt-3.5-turbo",
		},
	}

	// 重置计时器
	b.ResetTimer()

	// 执行基准测试
	for i := 0; i < b.N; i++ {
		err := cfg.Validate()
		if err != nil {
			b.Fatalf("配置验证失败: %v", err)
		}
	}
}

// 测试配置访问性能
func BenchmarkConfigAccess(b *testing.B) {
	// 创建一个有效的配置对象
	cfg := config.DefaultConfig()
	cfg.Global.ProjectName = "vhagar"
	cfg.Global.LogLevel = "info"
	cfg.Services.AI.Enable = true
	cfg.Services.AI.Provider = "openai"
	cfg.Services.AI.Providers = map[string]config.ProviderConfig{
		"openai": {
			ApiKey: "sk-test",
			ApiUrl: "https://api.openai.com/v1",
			Model:  "gpt-3.5-turbo",
		},
	}
	cfg.Database.PG.Host = "localhost"
	cfg.Database.PG.Port = 5432
	cfg.Database.PG.Username = "postgres"
	cfg.Database.PG.Password = "password"
	cfg.Database.PG.Database = "vhagar"

	// 重置计时器
	b.ResetTimer()

	// 执行基准测试
	for i := 0; i < b.N; i++ {
		// 测试各种配置访问方法
		_ = cfg.IsAIEnabled()
		_, _ = cfg.GetAIProvider()
		_ = cfg.IsMetricEnabled()
		_ = cfg.GetDatabaseConfig("pg")
	}
}
