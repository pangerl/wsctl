// Package doris @Author lanpang
// @Date 2024/9/13 下午3:38:00
// @Desc
package doris

import (
	"database/sql"
	"vhagar/config"
)

const taskName = "doris"

type Doris struct {
	//config.Global
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

func newDoris(cfg *config.CfgType) *Doris {
	return &Doris{
		//Global:   cfg.Global,
		DorisCfg: cfg.Doris,
	}
}
