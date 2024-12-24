// Package config @Author lanpang
// @Date 2024/9/12 下午5:21:00
// @Desc
package config

type Tenant struct {
	Corp []*Corp
}

type Corp struct {
	Corpid               string `json:"corpid"`
	Convenabled          bool   `json:"convenabled"`
	CorpName             string `json:"corpName"`
	MessageNum           int64  `json:"messageNum"`
	YesterdayMessageNum  int64  `json:"yesterdayMessageNum"`
	UserNum              int    `json:"userNum"`
	CustomerNum          int64  `json:"customerNum"`
	CustomerGroupNum     int    `json:"customerGroupNum"`
	CustomerGroupUserNum int    `json:"customerGroupUserNum"`
	DauNum               int64  `json:"dauNum"`
	WauNum               int64  `json:"wauNum"`
	MauNum               int64  `json:"mauNum"`
}
