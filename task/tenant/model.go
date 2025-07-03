// Package tenant @Author lanpang
// @Date 2024/8/6 下午3:50:00
// @Desc
package tenant

import (
	"database/sql"
	"vhagar/config"
	"vhagar/libs"

	"go.uber.org/zap"
)

// 每日巡检版本
var version = config.VERSION

const taskName = "tenant"

type Tenanter struct {
	Config *config.CfgType
	Logger *zap.SugaredLogger
	Corp   []*config.Corp
	//ESClient *elastic.Client
	MysqlClient *sql.DB
	PGClient    *libs.PGClienter
}

func NewTenanter(cfg *config.CfgType, logger *zap.SugaredLogger) *Tenanter {
	return &Tenanter{
		Config: cfg,
		Logger: logger,
		Corp:   cfg.Tenant.Corp,
	}
}
