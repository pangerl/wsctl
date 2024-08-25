// Package inspect @Author lanpang
// @Date 2024/8/6 下午3:50:00
// @Desc
package inspect

import (
	"github.com/olivere/elastic/v7"
	"vhagar/libs"
)

type Tenant struct {
	ProjectName string
	Version     string
	ProxyURL    string
	Corp        []*Corp
	Scheducron  string
	Robotkey    []string
	Userlist    []string
	ESClient    *elastic.Client
	PGClient    *libs.PGClient
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
type DB struct {
	Ip       string
	Port     int
	Username string
	Password string
	Sslmode  bool
}
