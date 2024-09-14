// Package host @Author lanpang
// @Date 2024/9/6 上午11:31:00
// @Desc
package host

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"vhagar/config"
	"vhagar/task"
)

//var hosts = make(map[string]*Host)

// var cfg = &config.Config

// TableRender 输出表格
func (s *Server) TableRender() {
	hosts := s.Hosts
	tabletitle := []string{"IP", "CPU使用率", "内存使用率", "内存大小", "入网流量", "出网流量", "系统盘使用率", "数据盘使用率"}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tabletitle)
	color := tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor}
	tableColor := []tablewriter.Colors{color, color, color, color, color, color, color, color}
	// ident 排序
	identList := ipSort(hosts)
	// 异常计数
	var alarmNum int
	for _, ident := range identList {
		data := hosts[ident]
		tabledata := []string{ident, formatToPercentage(data.CpuUsageActive),
			formatToPercentage(data.MemUsedPercent), formatBytes(data.MemTotal),
			formatBytes(data.netBytesRecv), formatBytes(data.netBytesSent),
			formatToPercentage(data.DiskUsedPercent["/"]), formatToPercentage(data.DiskUsedPercent["/data"])}
		// 异常标红
		if isAlarm(data) {
			table.Rich(tabledata, tableColor)
			alarmNum += 1
		}
		table.Append(tabledata)
	}
	identNum := len(hosts)
	caption := fmt.Sprintf("服务器计数: %d, 巡检异常计数: %d.", identNum, alarmNum)
	table.SetCaption(true, caption)
	table.Render()
}

func (s *Server) ReportRobot() {}

func Check() {
	task.EchoPrompt("开始巡检服务器状态")
	cfg := config.Config
	server := newServer(cfg)
	// 获取服务器信息
	initData(server)
	server.TableRender()
}

func initData(s *Server) {
	// CPU 使用率
	s.getHostData("cpu_usage_active")
	// 内存 使用率
	s.getHostData("mem_used_percent")
	// 内存 大小
	s.getHostData("mem_total")
	// 入网流量
	s.getHostData("rate", "net_bytes_recv", "eth0")
	// 出网流量
	s.getHostData("rate", "net_bytes_sent", "eth0")
	// 系统盘
	s.getHostData("disk_used_percent", "/")
	// 数据盘
	s.getHostData("disk_used_percent", "/data")
}

func ipSort(hosts map[string]*Host) []string {
	keys := make([]string, 0, len(hosts))
	for k := range hosts {
		keys = append(keys, k)
	}

	// 按 IP 地址排序
	sort.Slice(keys, func(i, j int) bool {
		ip1 := net.ParseIP(keys[i])
		ip2 := net.ParseIP(keys[j])
		return ip1.String() < ip2.String() // 进行字符串比较，即可实现排序
	})

	return keys
}

func isAlarm(host *Host) bool {
	if host.CpuUsageActive > 90 {
		return true
	}
	if host.MemUsedPercent > 85 {
		return true
	}
	if host.DiskUsedPercent["/"] > 80 {
		return true
	}
	if host.DiskUsedPercent["/data"] > 85 {
		return true
	}
	return false
}

func queryVictoriaMetrics(url string) (*MetricsResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed info: %s \n", err)
			return
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to query VictoriaMetrics: %s", resp.Status)
	}

	var metricsResponse MetricsResponse
	if err := json.NewDecoder(resp.Body).Decode(&metricsResponse); err != nil {
		return nil, err
	}

	return &metricsResponse, nil
}

func (s *Server) getHostData(key ...string) {
	var url string
	var baseUrl = s.vmUrl
	if len(key) == 1 {
		url = baseUrl + "/api/v1/query?query=" + key[0]
	} else if len(key) == 2 {
		url = fmt.Sprintf("%s/api/v1/query?query=%s{path='%s'}", baseUrl, key[0], key[1])
	} else if len(key) == 3 {
		url = fmt.Sprintf("%s/api/v1/query?query=%s(%s{interface='%s'}[1m])", baseUrl, key[0], key[1], key[2])
	} else {
		log.Printf("不支持的参数")
		return
	}
	//fmt.Println(url)
	response, _ := queryVictoriaMetrics(url)
	if response.Status != "success" {
		log.Printf("查询报错，Status: %s \n", response.Status)
		return
	}
	for _, result := range response.Data.Result {
		ident := result.Metric["ident"]
		host := s.getHost(ident)
		switch {
		case key[0] == "cpu_usage_active":
			host.CpuUsageActive, _ = strconv.ParseFloat(result.Value[1].(string), 64)
		case key[0] == "mem_used_percent":
			host.MemUsedPercent, _ = strconv.ParseFloat(result.Value[1].(string), 64)
		case key[0] == "mem_total":
			host.MemTotal, _ = strconv.ParseFloat(result.Value[1].(string), 64)
		case key[1] == "net_bytes_recv":
			host.netBytesRecv, _ = strconv.ParseFloat(result.Value[1].(string), 64)
		case key[1] == "net_bytes_sent":
			host.netBytesSent, _ = strconv.ParseFloat(result.Value[1].(string), 64)
		case key[0] == "disk_used_percent":
			host.DiskUsedPercent = make(map[string]float64)
			host.DiskUsedPercent[key[1]], _ = strconv.ParseFloat(result.Value[1].(string), 64)
		default:
			fmt.Printf("xxx")
		}
	}
}

func (s *Server) getHost(ident string) *Host {
	if host, exists := s.Hosts[ident]; exists {
		return host
	}
	newHost := Host{}
	s.Hosts[ident] = &newHost
	return s.Hosts[ident]
}

func formatToPercentage(value float64) string {
	result := fmt.Sprintf("%.1f%%", value)
	return result
}

func formatBytes(bytes float64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%.2f B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.2f KB", bytes/1024)
	} else if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", bytes/(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB", bytes/(1024*1024*1024))
	}
}
