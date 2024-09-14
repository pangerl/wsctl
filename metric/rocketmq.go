// Package metric @Author lanpang
// @Date 2024/8/15 下午4:37:00
// @Desc
package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
	"vhagar/task/rocketmq"
)

var (
	brokerCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "mq_broker_count",
			Help: "Status count for rocketmq broker",
		})
)

func setBrokerCount() {
	prometheus.MustRegister(brokerCount)
	m := rocketmq.GetRocketMQ()
	for {
		m.InitData()
		conut := len(m.BrokerMap)
		brokerCount.Set(float64(conut))
		log.Printf("brokercount: %v", conut)
		time.Sleep(30 * time.Second) // 每30秒探测一次
	}

}
