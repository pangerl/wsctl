// Package metric @Author lanpang
// @Date 2024/8/12 下午6:37:00
// @Desc
package metric

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
	"vhagar/config"
	"vhagar/libs"
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

func init() {
	// 只注册一次指标，避免重复注册 panic
	prometheus.MustRegister(probeHTTPStatusCode)
}

func strobeHTTPStatusCode(healthApi string) {
	// 不再在此注册指标，避免重复注册
	// prometheus.MustRegister(probeHTTPStatusCode)
	// 获取 newNacos 服务信息
	newNacos := nacos.NewNacos(config.Config, libs.Logger)
	err := newNacos.Init()
	if err != nil {
		libs.Logger.Errorw("初始化 Nacos 服务失败", "err", err)
		return
	}

	// 设置一个定时器来定期探测每个实例的健康状况
	for {
		libs.Logger.Warnw("检查服务接口健康状态")
		newNacos.Gather()
		healthInstances := newNacos.Clusterdata.HealthInstance

		var wg sync.WaitGroup
		for _, instance := range healthInstances {
			wg.Add(1)
			go func(inst nacos.ServerInstance) {
				defer wg.Done()
				probeInstance(inst, healthApi)
			}(instance)
		}
		wg.Wait()
		time.Sleep(30 * time.Second)
	}
}

// probeInstance 发送 HTTP 请求并检查返回值
func probeInstance(instance nacos.ServerInstance, healthApi string) {
	url := fmt.Sprintf("http://%s:%s%s", instance.Ip, instance.Port, healthApi)
	//libs.Logger.Infow("开始请求", "url", url, "namespace", instance.NamespaceName, "service", instance.ServiceName)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		libs.Logger.Errorw("请求 URL 失败", "url", url, "err", err, "service", instance.ServiceName)
		probeHTTPStatusCode.WithLabelValues(instance.NamespaceName, instance.ServiceName, instance.Ip, instance.Port, url).Set(0)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			libs.Logger.Errorw("请求失败", "err", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		libs.Logger.Errorw("读取响应体失败", "url", url, "err", err)
		probeHTTPStatusCode.WithLabelValues(instance.NamespaceName, instance.ServiceName, instance.Ip, instance.Port, url).Set(0)
		return
	}

	// 解析 JSON 并判断 status 字段
	var result struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(body, &result); err == nil {
		if result.Status == "UP" {
			probeHTTPStatusCode.WithLabelValues(instance.NamespaceName, instance.ServiceName, instance.Ip, instance.Port, url).Set(1)
			return
		}
	}
	probeHTTPStatusCode.WithLabelValues(instance.NamespaceName, instance.ServiceName, instance.Ip, instance.Port, url).Set(0)
}
