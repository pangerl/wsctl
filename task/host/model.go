// Package host @Author lanpang
// @Date 2024/9/6 上午11:37:00
// @Desc
package host

import "vhagar/config"

type Server struct {
	config.Global
	vmUrl string
	Hosts map[string]*Host
}

func newServer(cfg *config.CfgType) *Server {
	return &Server{
		cfg.Global,
		cfg.VictoriaMetrics,
		make(map[string]*Host),
	}
}

type MetricsResponse struct {
	Status string   `json:"status"`
	Data   Response `json:"data"`
}

type Response struct {
	Result []MetricData `json:"result"`
}

type MetricData struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}

type Host struct {
	//Ident           string
	CpuUsageActive  float64
	MemUsedPercent  float64
	MemTotal        float64
	netBytesRecv    float64
	netBytesSent    float64
	DiskUsedPercent map[string]float64
}
