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
var version = "v4.6"

type Tenanter struct {
	config.Global
	Corp     []*config.Corp
	ESClient *elastic.Client
	PGClient *libs.PGClient
}

func newTenant(cfg *config.CfgType) *Tenanter {
	return &Tenanter{
		cfg.Global,
		cfg.Tenant.Corp,
		nil,
		nil,
	}
}
