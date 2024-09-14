// Package config @Author lanpang
// @Date 2024/9/13 下午3:34:00
// @Desc
package config

import "vhagar/libs"

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
