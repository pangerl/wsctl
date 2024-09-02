// Package inspect @Author lanpang
// @Date 2024/8/6 下午3:50:00
// @Desc
package inspect

import (
	"database/sql"
	"github.com/olivere/elastic/v7"
	"vhagar/libs"
)

// 每日巡检版本
var version = "v4.6"

type Inspect struct {
	ProjectName string
	ProxyURL    string
	Rocketmq    libs.Rocketmq
	Notifier    map[string]Notifier
	Tenant      *Tenant
	Doris       *Doris
}

type Tenant struct {
	Corp     []*Corp
	ESClient *elastic.Client
	PGClient *libs.PGClient
}

type Corp struct {
	Corpid               string
	Convenabled          bool
	CorpName             string
	MessageNum           int64
	UserNum              int
	CustomerNum          int64
	CustomerGroupNum     int
	CustomerGroupUserNum int
	DauNum               int64
	WauNum               int64
	MauNum               int64
}

type Doris struct {
	MysqlClient        *sql.DB
	FailedJobs         []string
	StaffCount         int
	UseAnalyseCount    int
	CustomerGroupCount int
}

type Notifier struct {
	Robotkey []string
	Userlist []string
}
