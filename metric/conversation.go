// Package metric @Author lanpang
// @Date 2024/8/20 下午5:31:00
// @Desc
package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	messageCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "message_count",
			Help: "Status count for conversation",
		})
)

//func setMessageCount(ecfg inspect.DB, tenant inspect.Tenant) {
//	prometheus.MustRegister(messageCount)
//	esclient, _ := inspect.NewESClient(ecfg)
//	for {
//		log.Printf("brokercount: %v", brokercount)
//		time.Sleep(30 * time.Second) // 每30秒探测一次
//	}
//}
