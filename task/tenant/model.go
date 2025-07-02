// Package tenant @Author lanpang
// @Date 2024/8/6 下午3:50:00
// @Desc
package tenant

import (
	"database/sql"
	"vhagar/config"
	"vhagar/libs"
)

// 每日巡检版本
var version = config.VERSION

const taskName = "tenant"

type Tenanter struct {
	config.Global
	Corp []*config.Corp
	//ESClient *elastic.Client
	MysqlClient *sql.DB
	PGClient    *libs.PGClienter
}

func newTenant(cfg *config.CfgType) *Tenanter {
	return &Tenanter{
		Global: cfg.Global,
		Corp:   cfg.Tenant.Corp,
	}
}
