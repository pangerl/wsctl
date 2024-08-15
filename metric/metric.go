// Package metric @Author lanpang
// @Date 2024/8/15 下午4:32:00
// @Desc
package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
	"vhagar/nacos"
)

func Monitor(nacos *nacos.Nacos, mqDashboard string) {
	// 注册 Prometheus 指标
	prometheus.MustRegister(probeHTTPStatusCode)
	prometheus.MustRegister(brokerCount)
	healthInstances := nacos.Clusterdata.HealthInstance

	// 设置一个定时器来定期探测每个实例的健康状况
	go func() {
		for {
			for _, instance := range healthInstances {
				probeInstance(instance)
			}
			setBrokerCount(mqDashboard)
			time.Sleep(30 * time.Second) // 每30秒探测一次
		}
	}()

	// 设置 HTTP 服务器并暴露 /metrics 端点
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting server at http://%s%s/metrics\n", getLocalIp(), nacos.Webport)
	log.Fatal(http.ListenAndServe(nacos.Webport, nil))
}
