// Package nacos @Author lanpang
// @Date 2024/8/12 下午6:37:00
// @Desc
package nacos

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// 定义 Prometheus 指标
var (
	probeHTTPStatusCode = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "probe_http_status_code",
			Help: "Status code of the HTTP probe for each server instance",
		},
		[]string{"namespace", "service", "ip", "port", "url"},
	)
)

// probeInstance 发送 HTTP 请求并检查返回值
func probeInstance(instance ServerInstance) {
	url := fmt.Sprintf("http://%s:%s/actuator/test", instance.Ip, instance.Port)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error requesting URL %s: %v\n", url, err)
		probeHTTPStatusCode.WithLabelValues(instance.NamespaceName, instance.ServiceName, instance.Ip, instance.Port, url).Set(1)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed info: %s \n", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body from %s: %v\n", url, err)
		probeHTTPStatusCode.WithLabelValues(instance.NamespaceName, instance.ServiceName, instance.Ip, instance.Port, url).Set(1)
		return
	}

	if strings.TrimSpace(string(body)) == "success" {
		probeHTTPStatusCode.WithLabelValues(instance.NamespaceName, instance.ServiceName, instance.Ip, instance.Port, url).Set(0)
	} else {
		probeHTTPStatusCode.WithLabelValues(instance.NamespaceName, instance.ServiceName, instance.Ip, instance.Port, url).Set(1)
	}
}

func Monitor(nacos *Nacos) {
	// 注册 Prometheus 指标
	prometheus.MustRegister(probeHTTPStatusCode)
	// 假设我们有一个 JSON 字符串形式的 HealthInstance
	healthInstances := nacos.Clusterdata.HealthInstance
	// 每500秒刷新服务状态
	var interval time.Duration
	interval = 500 * time.Second
	RefreshNacosInstance(nacos, interval)

	// 设置一个定时器来定期探测每个实例的健康状况
	go func() {
		for {
			for _, instance := range healthInstances {
				probeInstance(instance)
			}
			time.Sleep(30 * time.Second) // 每30秒探测一次
		}
	}()

	// 设置 HTTP 服务器并暴露 /metrics 端点
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting server at http://%s%s/metrics\n", getLocalIp(), nacos.Webport)
	log.Fatal(http.ListenAndServe(nacos.Webport, nil))
}

func getLocalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("获取本机 IP 地址失败:", err)
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
