// Package es @Author lanpang
// @Date 2024/9/19 下午2:48:00
// @Desc
package es

import (
	"vhagar/config"

	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

type ES struct {
	Config   *config.CfgType
	Logger   *zap.SugaredLogger
	ESClient *elastic.Client
	NodeList []*NodeInfo
	Status   string
	// 新增字段
	ClusterJVMUsage  float64
	UnassignedShards int
	TotalDataSize    int64
}

func NewES(cfg *config.CfgType, logger *zap.SugaredLogger) *ES {
	return &ES{
		Config: cfg,
		Logger: logger,
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
