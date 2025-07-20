// Package database Redis 数据库连接工具
package database

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// NewRedisClient 创建 Redis 数据库连接
// 参数:
//   - cfg: Redis 配置信息
//
// 返回:
//   - *redis.Client: Redis 客户端对象
//   - error: 错误信息
func NewRedisClient(cfg RedisConfig) (*redis.Client, error) {
	// 创建 Redis 客户端配置
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,     // Redis 服务器地址
		Password: cfg.Password, // Redis 密码
		DB:       cfg.DB,       // Redis 数据库编号
	})

	// 测试 Redis 连接是否成功
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		zap.S().Errorw("Redis 连接测试失败",
			"addr", cfg.Addr,
			"db", cfg.DB,
			"err", err)
		return nil, fmt.Errorf("Redis 连接失败: %w", err)
	}

	zap.S().Infow("Redis 连接成功",
		"addr", cfg.Addr,
		"db", cfg.DB)
	return client, nil
}

// NewRedisClientWithContext 创建带上下文的 Redis 数据库连接
// 参数:
//   - ctx: 上下文对象
//   - cfg: Redis 配置信息
//
// 返回:
//   - *redis.Client: Redis 客户端对象
//   - error: 错误信息
func NewRedisClientWithContext(ctx context.Context, cfg RedisConfig) (*redis.Client, error) {
	// 创建 Redis 客户端配置
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,     // Redis 服务器地址
		Password: cfg.Password, // Redis 密码
		DB:       cfg.DB,       // Redis 数据库编号
	})

	// 使用提供的上下文测试 Redis 连接
	_, err := client.Ping(ctx).Result()
	if err != nil {
		zap.S().Errorw("Redis 连接测试失败",
			"addr", cfg.Addr,
			"db", cfg.DB,
			"err", err)
		return nil, fmt.Errorf("Redis 连接失败: %w", err)
	}

	zap.S().Infow("Redis 连接成功",
		"addr", cfg.Addr,
		"db", cfg.DB)
	return client, nil
}
