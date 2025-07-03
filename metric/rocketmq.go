// Package metric @Author lanpang
// @Date 2024/8/15 下午4:37:00
// @Desc
package metric

import (
	"time"
	"vhagar/config"
	"vhagar/libs"
	"vhagar/task/rocketmq"

	"github.com/prometheus/client_golang/prometheus"
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
	rocket := rocketmq.NewRocketMQ(config.Config, libs.Logger)
	for {
		rocket.Gather()
		conut := len(rocket.BrokerMap)
		brokerCount.Set(float64(conut))
		libs.Logger.Infow("brokercount", "count", conut)
		time.Sleep(30 * time.Second) // 每30秒探测一次
	}

}
