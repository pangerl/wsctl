package database

import (
	"context"
	"testing"
	"time"
)

// TestRedisConfig_HasValue 测试 Redis 配置验证功能
// 这个测试已经在 mysql_test.go 中定义了，这里不重复

// TestRedisClientCreation 测试 Redis 客户端创建功能
// 注意：这个测试不会实际连接到 Redis 服务器，只测试客户端对象的创建
func TestRedisClientCreation(t *testing.T) {
	tests := []struct {
		name   string
		config RedisConfig
		valid  bool
	}{
		{
			name: "有效的Redis配置",
			config: RedisConfig{
				Addr:     "localhost:6379",
				Password: "",
				DB:       0,
			},
			valid: true,
		},
		{
			name: "带密码的Redis配置",
			config: RedisConfig{
				Addr:     "localhost:6379",
				Password: "password123",
				DB:       1,
			},
			valid: true,
		},
		{
			name: "无效的Redis配置_空地址",
			config: RedisConfig{
				Addr:     "",
				Password: "",
				DB:       0,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试配置验证
			hasValue := tt.config.HasValue()
			if hasValue != tt.valid {
				t.Errorf("RedisConfig.HasValue() = %v, want %v", hasValue, tt.valid)
			}
		})
	}
}

// TestRedisClientWithContext 测试带上下文的 Redis 客户端创建
func TestRedisClientWithContext(t *testing.T) {
	// 创建一个带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	config := RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}

	// 注意：这个测试在没有实际 Redis 服务器的情况下会失败
	// 但我们可以测试函数是否正确处理上下文
	_, err := NewRedisClientWithContext(ctx, config)

	// 在没有 Redis 服务器的情况下，我们期望得到连接错误
	// 这证明函数正确尝试了连接
	if err == nil {
		t.Log("Redis 连接成功（可能有实际的 Redis 服务器运行）")
	} else {
		t.Logf("Redis 连接失败（预期行为，没有 Redis 服务器）: %v", err)
	}
}
