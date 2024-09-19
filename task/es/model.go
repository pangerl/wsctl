// Package es @Author lanpang
// @Date 2024/9/19 下午2:48:00
// @Desc
package es

import (
	"github.com/olivere/elastic/v7"
	"vhagar/config"
)

type ES struct {
	config.Global
	ESClient *elastic.Client
	NodeList []*NodeInfo
	Status   string
}

func newES(cfg *config.CfgType) *ES {
	return &ES{
		Global: cfg.Global,
	}
}

type NodeInfo struct {
	Name      string
	IP        string
	JVMUsage  float64
	DiskUsage float64
}
