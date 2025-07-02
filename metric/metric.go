// Package metric @Author lanpang
// @Date 2024/8/15 下午4:32:00
// @Desc
package metric

import (
	"fmt"
	"net/http"
	"vhagar/config"
	"vhagar/libs"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	libs.Logger.Infow("启动 metrics 服务", "url", fmt.Sprintf("http://%s:%s/metrics", getClientIp(), cfg.Port))
	err := http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		libs.Logger.Fatalw("metrics 服务启动失败", "err", err)
	}
}
