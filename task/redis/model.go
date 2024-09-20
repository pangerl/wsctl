// Package redis @Author lanpang
// @Date 2024/9/20 上午10:26:00
// @Desc
package redis

import (
	"vhagar/config"
	"vhagar/libs"
)

const taskName = "redis"

type Redis struct {
	config.Global
	Config         libs.RedisConfig
	Version        string
	Role           string
	Slaves         int
	CurrentClients int
	MaxClients     int
	UsedMemory     string
	KeyCount       int
}

func NewRedis(cfg *config.CfgType) *Redis {
	return &Redis{
		Global: cfg.Global,
		Config: cfg.Redis,
	}
}
