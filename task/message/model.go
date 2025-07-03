// Package message @Author lanpang
// @Date 2024/8/6 下午3:50:00
// @Desc
package message

import (
	"vhagar/config"
	"vhagar/libs"

	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

const taskName = "message"

type Tenanter struct {
	Config    *config.CfgType
	Logger    *zap.SugaredLogger
	NasDir    string
	DirIsExis bool
	Corp      []*config.Corp
	ESClient  *elastic.Client
	PGClient  *libs.PGClienter
}

func NewTenanter(cfg *config.CfgType, logger *zap.SugaredLogger) *Tenanter {
	return &Tenanter{
		Config: cfg,
		Logger: logger,
		Corp:   cfg.Tenant.Corp,
		NasDir: cfg.NasDir,
	}
}
