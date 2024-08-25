// Package common @Author lanpang
// @Date 2024/8/1 下午2:47:00
// @Desc
package common

import (
	"vhagar/inspect"
	"vhagar/metric"
	"vhagar/nacos"
)

type Config struct {
	ProjectName string
	ProxyURL    string
	Crontab     Crontab
	Nacos       nacos.Config
	Tenant      inspect.Tenant
	PG          inspect.DB
	ES          inspect.DB
	Doris       inspect.DB
	Rocketmq    Rocketmq
	Metric      metric.Config
}

type Crontab struct {
	TenantJob bool
	TestJob   bool
}

type Rocketmq struct {
	RocketmqDashboard string
	NameServer        string
}
