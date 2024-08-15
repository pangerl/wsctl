// Package metric @Author lanpang
// @Date 2024/8/12 下午6:37:00
// @Desc
package metric

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"vhagar/nacos"
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
func probeInstance(instance nacos.ServerInstance) {
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
