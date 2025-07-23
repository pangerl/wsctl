// Package metric @Author lanpang
// @Date 2024/8/15 下午4:32:00
// @Desc
package metric

import (
	"fmt"
	"net"
	"net/http"
	"vhagar/config"
	"vhagar/logger"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartMetric() {
	cfg := config.Config.Metric
	// 服务健康检查
	go strobeHTTPStatusCode(cfg.HealthApi)
	// rocketmq 指标
	go setBrokerCount()
	// 会话数统计
	go setMessageCount()
	http.Handle("/metrics", promhttp.Handler())
	logger.Logger.Warnw("启动 metrics 服务", "url", fmt.Sprintf("http://%s:%s/metrics", getClientIp(), cfg.Port))
	err := http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		logger.Logger.Fatalw("metrics 服务启动失败", "err", err)
	}
}

func getClientIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logger.Logger.Errorw("获取本机 IP 地址失败", "err", err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//fmt.Println("本机 IP 地址:", ipnet.IP.String())
				return ipnet.IP.String()
			}
		}
	}

	return ""
}
