// Package cmd @Author lanpang
// @Date 2024/8/1 下午2:47:00
// @Desc
package cmd

import (
	"vhagar/inspect"
	"vhagar/metric"
	"vhagar/nacos"
)

var (
	CONFIG = Config{
		ProjectName: "测试项目",
	}
)

const VERSION = "v1.0"

type Config struct {
	ProjectName string
	ProxyURL    string
	Crontab     Crontab
	Nacos       nacos.Config
	Tenant      inspect.Tenant
	PG          inspect.DB
	ES          inspect.DB
	Inspection  inspect.Config
	Rocketmq    inspect.Rocketmq
	Metric      metric.Config
}

type Crontab struct {
	Inspectjob bool
	Testjob    bool
}
