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
	Nacos       nacos.Config
	Tenant      inspect.Tenant
	PG          libs.DB
	ES          libs.DB
	Doris       inspect.Doris
	Rocketmq    libs.Rocketmq
	Metric      metric.Config
}

//type Crontab struct {
//	TenantJob bool
//	TestJob   bool
//}
