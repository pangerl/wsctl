// Package metric @Author lanpang
// @Date 2024/8/15 下午4:32:00
// @Desc
package metric

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func StartMetric(m *Metric) {

	if m.Metric.Wsapp {
		go setprobeHTTPStatusCode(m.Nacos)
	}

	if m.Metric.Rocketmq {
		go setBrokerCount(m.Rocketmq.RocketmqDashboard)
	}

	if m.Metric.Conversation {
		go setMessageCount(m)
	}
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting server at http://%s:%s/metrics\n", getLocalIp(), m.Metric.Port)
	log.Fatal(http.ListenAndServe(":"+m.Metric.Port, nil))
}
