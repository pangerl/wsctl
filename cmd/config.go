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
	Proxyurl    string
	Nacos       nacos.Config
	Tenant      inspect.Tenant
	PG          inspect.DB
	ES          inspect.DB
	Inspection  inspect.Config
}

type Inspection struct {
	Scheducron string
	Robotkey   string
}
