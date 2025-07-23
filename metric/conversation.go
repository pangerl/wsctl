// Package metric @Author lanpang
// @Date 2024/8/20 下午5:31:00
// @Desc
package metric

import (
	"time"
	"vhagar/config"
	"vhagar/database"
	"vhagar/logger"
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

func init() {
	// 只注册一次指标，避免重复注册 panic
	prometheus.MustRegister(messageCount)
}

func setMessageCount() {
	// 不再在此注册指标，避免重复注册
	// prometheus.MustRegister(messageCount)

	// 初始化 esclient
	esclient, _ := database.NewElasticsearchClient(config.Config.Database.ES)
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
				logger.Logger.Infow("corp messagenum", "corpid", corp.Corpid, "messagenum", messagenum)
			}
		}
		time.Sleep(300 * time.Second)
	}
}
