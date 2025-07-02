// Package metric @Author lanpang
// @Date 2024/8/20 下午5:31:00
// @Desc
package metric

import (
	"time"
	"vhagar/config"
	"vhagar/libs"
	"vhagar/task/message"

	"github.com/prometheus/client_golang/prometheus"
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

func setMessageCount() {
	prometheus.MustRegister(messageCount)

	// 初始化 esclient
	esclient, _ := libs.NewESClient(config.Config.ES)
	defer func() {
		if esclient != nil {
			esclient.Stop()
		}
	}()
	if esclient == nil {
		return
	}
	corpList := config.Config.Tenant.Corp
	for {
		dateNow := time.Now()
		for _, corp := range corpList {
			if corp.Convenabled {
				messagenum := message.CurrentMessageNum(esclient, corp.Corpid, dateNow)
				messageCount.WithLabelValues(corp.Corpid).Set(float64(messagenum))
				libs.Logger.Infow("corp messagenum", "corpid", corp.Corpid, "messagenum", messagenum)
			}
		}
		time.Sleep(300 * time.Second)
	}
}
