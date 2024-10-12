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

func StartMetric(port string) {
	cfg := config.Config.Metric
	if cfg.Wsapp {
		go setprobeHTTPStatusCode()
	}

	if cfg.Rocketmq {
		go setBrokerCount()
	}

	if cfg.Conversation {
		go setMessageCount()
	}
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting server at http://%s:%s/metrics\n", getLocalIp(), port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
