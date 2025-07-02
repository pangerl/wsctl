package libs

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func NewRedisClient(cfg RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		return nil, err
	}

	zap.S().Infow("redis 连接成功！")
	return client, nil
}
