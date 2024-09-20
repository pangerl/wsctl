// Package es @Author lanpang
// @Date 2024/9/19 下午2:48:00
// @Desc
package es

import (
	"vhagar/config"

	"github.com/olivere/elastic/v7"
)

type ES struct {
	config.Global
	ESClient *elastic.Client
	NodeList []*NodeInfo
	Status   string
	// 新增字段
	ClusterJVMUsage  float64
	UnassignedShards int
	TotalDataSize    int64
}

func newES(cfg *config.CfgType) *ES {
	return &ES{
		Global: cfg.Global,
	}
}

type NodeInfo struct {
	Name        string
	IP          string
	JVMUsage    float64
	DiskUsage   float64
	LoadAverage float64
	IndexCount  int
	Shards      int
	DataSize    int64
}
