// Package cmd @Author lanpang
// @Date 2024/8/1 下午2:47:00
// @Desc
package cmd

import (
	"vhagar/inspect"
	"vhagar/nacos"
)

var (
	CONFIG Config
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
}

type Crontab struct {
	Inspectjob bool
	Testjob    bool
}
