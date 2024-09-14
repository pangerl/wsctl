// Package config @Author lanpang
// @Date 2024/9/13 下午3:34:00
// @Desc
package config

import (
	"math/rand"
	"time"
	"vhagar/libs"
)

type DorisCfg struct {
	libs.DB
	HttpPort int `toml:"httpport"`
}

type RocketMQCfg struct {
	RocketmqDashboard string `toml:"rocketmqdashboard"`
	NameServer        string `toml:"nameserver"`
}

type NacosCfg struct {
	Server    string `json:"server"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Namespace string `json:"namespace"`
	Writefile string
}

func GetRandomDuration() time.Duration {
	// 创建一个新的随机数生成器，使用当前时间作为种子
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	// 生成随机数
	randomSeconds := r.Intn(300)
	// 将随机秒数转换为时间.Duration
	duration := time.Duration(randomSeconds) * time.Second
	return duration
}
