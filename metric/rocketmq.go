// Package metric @Author lanpang
// @Date 2024/8/15 下午4:37:00
// @Desc
package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"vhagar/inspect"
)

var (
	brokerCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "mq_broker_count",
			Help: "Status count for rocketmq broker",
		})
)

func setBrokerCount(mqDashboard string) {
	clusterdata, _ := inspect.GetMQDetail(mqDashboard)
	brokercount := inspect.GetBrokerCount(clusterdata)
	log.Printf("brokercount: %v", brokercount)
	brokerCount.Set(float64(brokercount))
}
