// Package message @Author lanpang
// @Date 2024/8/6 下午3:50:00
// @Desc
package message

import (
	"github.com/olivere/elastic/v7"
	"vhagar/config"
	"vhagar/libs"
)

const taskName = "message"

type Tenanter struct {
	config.Global
	NasDir    string
	DirIsExis bool
	Corp      []*config.Corp
	ESClient  *elastic.Client
	PGClient  *libs.PGClienter
}

func newTenant(cfg *config.CfgType) *Tenanter {
	return &Tenanter{
		Global: cfg.Global,
		Corp:   cfg.Tenant.Corp,
		NasDir: cfg.NasDir,
	}
}
