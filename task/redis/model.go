// Package redis @Author lanpang
// @Date 2024/9/20 上午10:26:00
// @Desc
package redis

import (
	"vhagar/config"
	"vhagar/libs"

	"go.uber.org/zap"
)

const taskName = "redis"

type Redis struct {
	Config         *config.CfgType
	Logger         *zap.SugaredLogger
	RedisConfig    libs.RedisConfig
	Version        string
	Role           string
	Slaves         int
	CurrentClients int
	MaxClients     int
	UsedMemory     string
	KeyCount       int
}

func NewRedis(cfg *config.CfgType, logger *zap.SugaredLogger) *Redis {
	return &Redis{
		Config:      cfg,
		Logger:      logger,
		RedisConfig: cfg.Redis,
	}
}
