// Package metric @Author lanpang
// @Date 2024/8/20 下午5:31:00
// @Desc
package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
	"vhagar/inspect"
)

var (
	messageCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "message_count",
			Help: "Status count for conversation",
		},
		[]string{"corpid"},
	)
)

func setMessageCount(m *Metric) {
	prometheus.MustRegister(messageCount)
	for {
		dateNow := time.Now()
		for _, corp := range m.Corp {
			if corp.Convenabled {
				messagenum := inspect.CurrentMessageNum(m.EsClient, corp.Corpid, dateNow)
				messageCount.WithLabelValues(corp.Corpid).Set(float64(messagenum))
				log.Printf("corp %s messagenum: %v", corp.Corpid, messagenum)
			}

		}
		time.Sleep(300 * time.Second)
	}
}
