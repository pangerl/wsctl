// Package metric @Author lanpang
// @Date 2024/8/15 下午4:37:00
// @Desc
package metric

import (
	"time"
	"vhagar/config"
	"vhagar/logger"
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

func init() {
	// 只注册一次指标，避免重复注册 panic
	prometheus.MustRegister(brokerCount)
}

func setBrokerCount() {
	// 不再在此注册指标，避免重复注册
	// prometheus.MustRegister(brokerCount)
	rocket := rocketmq.NewRocketMQ(config.Config, logger.Logger)
	for {
		rocket.Gather()
		conut := len(rocket.BrokerMap)
		brokerCount.Set(float64(conut))
		logger.Logger.Infow("brokercount", "count", conut)
		time.Sleep(60 * time.Second) // 每30秒探测一次
	}

}
