// Package inspect @Author lanpang
// @Date 2024/8/6 下午3:50:00
// @Desc
package inspect

import (
	"github.com/jackc/pgx/v5"
	"github.com/olivere/elastic/v7"
)

type Inspect struct {
	Corp      []*Corp
	EsClient  *elastic.Client
	PgClient1 *pgx.Conn
	PgClient2 *pgx.Conn
}

type Tenant struct {
	Corp []*Corp
}

type Corp struct {
	Corpid      string
	Convenabled bool
	CorpName    string
	MessageNum  int64
	UserNum     int
	CustomerNum int64
	DauNum      int64
	WauNum      int64
	MauNum      int64
}

type DB struct {
	Ip       string
	Port     int
	Username string
	Password string
	Sslmode  bool
}
