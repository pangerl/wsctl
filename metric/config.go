// Package metric @Author lanpang
// @Date 2024/8/20 下午6:17:00
// @Desc
package metric

import (
	"github.com/olivere/elastic/v7"
	"vhagar/inspect"
	"vhagar/nacos"
)

type Metric struct {
	Corp     []*inspect.Corp
	EsClient *elastic.Client
	Rocketmq inspect.Rocketmq
	Metric   Config
	Nacos    nacos.Config
}

type Config struct {
	Port         string
	Wsapp        bool
	Rocketmq     bool
	Conversation bool
}
