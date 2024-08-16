// Package metric @Author lanpang
// @Date 2024/8/15 下午4:37:00
// @Desc
package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
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
	prometheus.MustRegister(brokerCount)
	for {
		clusterdata, _ := inspect.GetMQDetail(mqDashboard)
		brokercount := inspect.GetBrokerCount(clusterdata)
		brokerCount.Set(float64(brokercount))
		log.Printf("brokercount: %v", brokercount)
		time.Sleep(30 * time.Second) // 每30秒探测一次
	}

}
