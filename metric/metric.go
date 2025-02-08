// Package metric @Author lanpang
// @Date 2024/8/15 下午4:32:00
// @Desc
package metric

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"vhagar/config"
)

func StartMetric() {
	cfg := config.Config.Metric
	// 服务健康检查
	go setprobeHTTPStatusCode(cfg.HealthApi)
	// rocketmq 指标
	go setBrokerCount()
	// 会话数统计
	go setMessageCount()
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting server at http://%s:%s/metrics\n", getClientIp(), cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
