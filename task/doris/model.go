// Package doris @Author lanpang
// @Date 2024/9/13 下午3:38:00
// @Desc
package doris

import (
	"database/sql"
	"vhagar/config"

	"go.uber.org/zap"
)

const taskName = "doris"

type Doris struct {
	Config *config.CfgType
	Logger *zap.SugaredLogger
	config.DorisCfg
	MysqlClient        *sql.DB
	FailedJobs         []string
	StaffCount         int
	UseAnalyseCount    int
	CustomerGroupCount int
	OnlineBackendNum   int
	TotalBackendNum    int
}

type dorisResponse struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data struct {
		OnlineBackendNum int `json:"online_backend_num"`
		TotalBackendNum  int `json:"total_backend_num"`
	} `json:"data"`
	Count int `json:"count"`
}

func NewDoris(cfg *config.CfgType, logger *zap.SugaredLogger) *Doris {
	return &Doris{
		Config:   cfg,
		Logger:   logger,
		DorisCfg: cfg.Doris,
	}
}
