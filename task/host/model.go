// Package host @Author lanpang
// @Date 2024/9/6 上午11:37:00
// @Desc
package host

import (
	"vhagar/config"

	"go.uber.org/zap"
)

const taskName = "host"

type Server struct {
	Config *config.CfgType
	Logger *zap.SugaredLogger
	VmUrl  string
	Hosts  map[string]*Host
}

func NewServer(cfg *config.CfgType, logger *zap.SugaredLogger) *Server {
	return &Server{
		Config: cfg,
		Logger: logger,
		VmUrl:  cfg.VictoriaMetrics,
		Hosts:  make(map[string]*Host),
	}
}

type MetricsResponse struct {
	Status string   `json:"status"`
	Data   Response `json:"data"`
}

type Response struct {
	Result []*MetricData `json:"result"`
}

type MetricData struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}

type Host struct {
	//Ident           string
	cpuCores            float64
	cpuUsageActive      float64
	MemUsedPercent      float64
	MemTotal            float64
	netBytesRecv        float64
	netBytesSent        float64
	rootDiskUsedPercent float64
	dataDiskUsedPercent float64
	ntpOffsetMs         float64
}
