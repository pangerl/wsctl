// Package metric @Author lanpang
// @Date 2024/8/20 下午6:17:00
// @Desc
package metric

import (
	"vhagar/inspect"
	"vhagar/libs"
	"vhagar/nacos"

	"github.com/olivere/elastic/v7"
)

type Metric struct {
	Corp     []*inspect.Corp
	EsClient *elastic.Client
	Rocketmq libs.Rocketmq
	Metric   Config
	Nacos    nacos.Config
}

type Config struct {
	Port         string `json:"port"`
	Wsapp        bool   `json:"wsapp"`
	Rocketmq     bool   `json:"rocketmq"`
	Conversation bool   `json:"conversation"`
}
