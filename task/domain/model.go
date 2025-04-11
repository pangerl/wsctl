// Package domain @Author Trae AI
// @Date 2024/8/6 下午3:50:00
// @Desc 域名连通性检测任务
package domain

import (
	"vhagar/config"
)

// 每日巡检版本
var version = config.VERSION

const taskName = "domain"

// Domain 结构体，用于存储域名连通性检测结果
type Domain struct {
	Name    string `json:"name"`    // 域名
	Port    int    `json:"port"`    // 端口
	IsAlive bool   `json:"isAlive"` // 是否连通
}

// Domainer 域名检测任务结构体
type Domainer struct {
	//config.Global
	Domains     []*Domain // 域名列表
	TotalCount  int       // 总域名数
	AliveCount  int       // 连通域名数
	FailedCount int       // 不通域名数
}

func newDomain(cfg *config.CfgType) *Domainer {
	return &Domainer{
		//Global:  cfg.Global,
		Domains: make([]*Domain, 0),
	}
}
