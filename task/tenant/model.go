// Package tenant @Author lanpang
// @Date 2024/8/6 下午3:50:00
// @Desc
package tenant

import (
	"github.com/olivere/elastic/v7"
	"vhagar/config"
	"vhagar/libs"
)

// 每日巡检版本
var version = config.VERSION

const taskName = "tenant"

type Tenanter struct {
	config.Global
	Corp     []*config.Corp
	ESClient *elastic.Client
	PGClient *libs.PGClienter
}

func newTenant(cfg *config.CfgType) *Tenanter {
	return &Tenanter{
		Global: cfg.Global,
		Corp:   cfg.Tenant.Corp,
	}
}
