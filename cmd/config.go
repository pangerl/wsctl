// Package cmd @Author lanpang
// @Date 2024/8/1 下午2:47:00
// @Desc
package cmd

import (
	"vhagar/inspect"
	"vhagar/libs"
	"vhagar/metric"
	"vhagar/nacos"
)

var (
	CONFIG = &Config{
		ProjectName: "测试项目",
	}
	cfgFile string
)

type Config struct {
	ProjectName string
	ProxyURL    string
	Cron        map[string]crontab
	Notifier    map[string]inspect.Notifier
	Nacos       nacos.Config
	Tenant      inspect.Tenant
	PG          libs.DB
	ES          libs.DB
	Doris       libs.DB
	Rocketmq    libs.Rocketmq
	Metric      metric.Config
}

type crontab struct {
	Crontab    bool
	Scheducron string
}
