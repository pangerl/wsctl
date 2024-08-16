// Package metric @Author lanpang
// @Date 2024/8/15 下午4:32:00
// @Desc
package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"vhagar/inspect"
	"vhagar/nacos"
)

type Metric struct {
	Port     string
	Wsapp    bool
	Rocketmq bool
}

var (
	randomNumber = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "random_number",
		Help: "Randomly generated number",
	})
)

func StartMetric(cfg Metric, ncfg nacos.Config, mcfg inspect.Rocketmq) {

	if cfg.Wsapp {
		go setprobeHTTPStatusCode(ncfg)
	}

	if cfg.Rocketmq {
		go setBrokerCount(mcfg.RocketmqDashboard)
	}
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting server at http://%s:%s/metrics\n", getLocalIp(), cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
