// Package metric @Author lanpang
// @Date 2024/8/15 下午4:32:00
// @Desc
package metric

import (
	"github.com/olivere/elastic/v7"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"vhagar/inspect"
	"vhagar/nacos"
)

func NewMetric(cfg Config, ncfg nacos.Config, mcfg inspect.Rocketmq, corp []*inspect.Corp, es *elastic.Client) *Metric {
	return &Metric{
		Corp:     corp,
		EsClient: es,
		Rocketmq: mcfg,
		Metric:   cfg,
		Nacos:    ncfg,
	}
}

func (m *Metric) StartMetric() {

	if m.Metric.Wsapp {
		go setprobeHTTPStatusCode(m.Nacos)
	}

	if m.Metric.Rocketmq {
		go setBrokerCount(m.Rocketmq.RocketmqDashboard)
	}

	if m.Metric.Conversation {
		go setMessageCount(m)
	}
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting server at http://%s:%s/metrics\n", getLocalIp(), m.Metric.Port)
	log.Fatal(http.ListenAndServe(":"+m.Metric.Port, nil))
}
