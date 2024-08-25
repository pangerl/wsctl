// Package metric @Author lanpang
// @Date 2024/8/20 下午6:17:00
// @Desc
package metric

import (
	"vhagar/inspect"
	"vhagar/nacos"

	"github.com/olivere/elastic/v7"
)

type Metric struct {
	Corp     []*inspect.Corp
	EsClient *elastic.Client
	Rocketmq Rocketmq
	Metric   Config
	Nacos    nacos.Config
}

type Config struct {
	Port         string
	Wsapp        bool
	Rocketmq     bool
	Conversation bool
}

type Rocketmq struct {
	RocketmqDashboard string
	NameServer        string
}
