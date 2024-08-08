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

type Config struct {
	ProjectName string
	Nacos       nacos.Config
	Tenant      inspect.Tenant
	PG          inspect.DB
	ES          inspect.DB
}
