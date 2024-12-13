// Package metric @Author lanpang
// @Date 2024/8/12 下午6:37:00
// @Desc
package metric

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
	"vhagar/config"
	"vhagar/task/nacos"

	"github.com/prometheus/client_golang/prometheus"
)

// 定义 Prometheus 指标
var (
	probeHTTPStatusCode = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "probe_success",
			Help: "Status code of the HTTP probe for each server instance",
		},
		[]string{"namespace", "service", "ip", "port", "url"},
	)
)

func setprobeHTTPStatusCode() {
	// 注册 Prometheus 指标
	prometheus.MustRegister(probeHTTPStatusCode)
	// 获取 nacos 服务信息
	n := nacos.NewNacos(config.Config)
	err := n.Init()
	if err != nil {
		log.Printf("初始化 Nacos 服务失败: %v\n", err)
		return
	}
	n.Gather()
	healthInstances := n.Clusterdata.HealthInstance

	// 设置一个定时器来定期探测每个实例的健康状况
	for {
		log.Println("检查服务接口健康状态")
		n.Gather()
		for _, instance := range healthInstances {
			probeInstance(instance)
		}
		time.Sleep(30 * time.Second) // 每30秒探测一次
	}
}

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

func getClientIp() string {
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
