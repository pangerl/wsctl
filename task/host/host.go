// Package host @Author lanpang
// @Date 2024/9/6 上午11:31:00
// @Desc
package host

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"vhagar/config"
	"vhagar/task"

	"github.com/olekukonko/tablewriter"
)

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
		tabledata := []string{ident, formatToPercentage(data.cpuUsageActive),
			formatToPercentage(data.MemUsedPercent), formatBytes(data.MemTotal),
			formatBytes(data.netBytesRecv), formatBytes(data.netBytesSent),
			formatToPercentage(data.rootDiskUsedPercent), formatToPercentage(data.dataDiskUsedPercent)}
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

func init() {
	task.Add(taskName, func() task.Tasker {
		return newServer(config.Config)
	})
}

func (s *Server) ReportRobot() {}

//func GetHost() *Server {
//	cfg := config.Config
//	server := newServer(cfg)
//	// 获取服务器信息
//	initData(server)
//	return server
//}

func (s *Server) Check() {
	task.EchoPrompt("开始巡检服务器状态")
	if s.Report {
		// 发送机器人
		s.ReportRobot()
		return
	}
	s.TableRender()
}

func (s *Server) Gather() {
	// CPU 使用率
	getHostCpuUsageActive(s)
	// 内存 使用率
	getHostMemUsedPercent(s)
	// 内存 大小
	getHostMemTotal(s)
	// 入网流量
	getHostNetBytesRecv(s)
	// 出网流量
	getHostNetBytesSent(s)
	// 系统盘
	getHostRootDiskUsedPercent(s)
	// 数据盘
	getHostDataDiskUsedPercent(s)
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
	if host.cpuUsageActive > 90 {
		return true
	}
	if host.MemUsedPercent > 85 {
		return true
	}
	if host.rootDiskUsedPercent > 80 {
		return true
	}
	if host.dataDiskUsedPercent > 85 {
		return true
	}
	return false
}

func queryVmData(url string) []*MetricData {
	body := task.DoRequest(url)
	var metricsResponse MetricsResponse
	if err := json.Unmarshal(body, &metricsResponse); err != nil {
		return nil
	}
	res := metricsResponse
	if res.Status != "success" {
		log.Printf("查询报错，Status: %s \n", res.Status)
		return nil
	}

	return res.Data.Result
}

//	func (s *Server) getHostData(key ...string) {
//		var url string
//		var baseUrl = s.vmUrl
//		if len(key) == 1 {
//			url = baseUrl + "/api/v1/query?query=" + key[0]
//		}
//		if key[0] == "disk_used_percent" {
//			url = fmt.Sprintf("%s/api/v1/query?query=%s{path='%s'}", baseUrl, key[0], key[1])
//		}
//		if key[0] == "rate" {
//			url = fmt.Sprintf("%s/api/v1/query?query=%s(%s{interface='%s'}[1m])", baseUrl, key[0], key[1], key[2])
//		}
//		//fmt.Println(url)
//		response, _ := queryVictoriaMetrics(url)
//		if response.Status != "success" {
//			log.Printf("查询报错，Status: %s \n", response.Status)
//			return
//		}
//		for _, result := range response.Data.Result {
//			ident := result.Metric["ident"]
//			host := s.getHost(ident)
//			switch {
//			case key[0] == "cpu_usage_active":
//				host.cpuUsageActive, _ = strconv.ParseFloat(result.Value[1].(string), 64)
//			case key[0] == "mem_used_percent":
//				host.MemUsedPercent, _ = strconv.ParseFloat(result.Value[1].(string), 64)
//			case key[0] == "mem_total":
//				host.MemTotal, _ = strconv.ParseFloat(result.Value[1].(string), 64)
//			case key[1] == "net_bytes_recv":
//				host.netBytesRecv, _ = strconv.ParseFloat(result.Value[1].(string), 64)
//			case key[1] == "net_bytes_sent":
//				host.netBytesSent, _ = strconv.ParseFloat(result.Value[1].(string), 64)
//			case key[0] == "disk_used_percent":
//				host.DiskUsedPercent = make(map[string]float64)
//				host.DiskUsedPercent[key[1]], _ = strconv.ParseFloat(result.Value[1].(string), 64)
//			default:
//				fmt.Printf("xxx")
//			}
//		}
//	}
func getHostCpuUsageActive(s *Server) {
	key := "cpu_usage_active"
	setHostData(s, key)
}

func getHostMemUsedPercent(s *Server) {
	key := "mem_used_percent"
	setHostData(s, key)
}

func getHostMemTotal(s *Server) {
	key := "mem_total"
	setHostData(s, key)
}

func getHostNetBytesRecv(s *Server) {
	key := "net_bytes_recv"
	nic := "eth0"
	url := fmt.Sprintf("%s/api/v1/query?query=rate(%s{interface='%s'}[1m])", s.vmUrl, key, nic)
	setHostData(s, key, url)
}

func getHostNetBytesSent(s *Server) {
	key := "net_bytes_sent"
	nic := "eth0"
	url := fmt.Sprintf("%s/api/v1/query?query=rate(%s{interface='%s'}[1m])", s.vmUrl, key, nic)
	setHostData(s, key, url)
}

func getHostRootDiskUsedPercent(s *Server) {
	key := "root_disk_used_percent"
	path := "/"
	url := fmt.Sprintf("%s/api/v1/query?query=disk_used_percent{path='%s'}", s.vmUrl, path)
	setHostData(s, key, url)
}

func getHostDataDiskUsedPercent(s *Server) {
	key := "data_disk_used_percent"
	path := "/data"
	url := fmt.Sprintf("%s/api/v1/query?query=disk_used_percent{path='%s'}", s.vmUrl, path)
	setHostData(s, key, url)
}

func setHostData(s *Server, key ...string) {
	var url string
	if len(key) == 1 {
		url = s.vmUrl + "/api/v1/query?query=" + key[0]
	} else {
		url = key[1]
	}
	fmt.Printf("url: %s\n", url)
	results := queryVmData(url)
	for _, result := range results {
		ident := result.Metric["ident"]
		value := result.Value[1].(string)
		fmt.Printf("ident: %s, key: %s, value: %s\n", ident, key[0], value)
		AddOrUpdateHost(s.Hosts, ident, key[0], value)
	}
}

func AddOrUpdateHost(hosts map[string]*Host, ident, key, value string) {
	if host, exists := hosts[ident]; exists {
		newValue, _ := strconv.ParseFloat(value, 64)
		switch {
		case key == "cpu_usage_active":
			host.cpuUsageActive = newValue
			fmt.Printf("cpu_usage_active: %f\n", newValue)
		case key == "mem_used_percent":
			host.MemUsedPercent = newValue
		case key == "mem_total":
			host.MemTotal = newValue
		case key == "net_bytes_recv":
			host.netBytesRecv = newValue
		case key == "net_bytes_sent":
			host.netBytesSent = newValue
		case key == "root_disk_used_percent":
			host.rootDiskUsedPercent = newValue
		case key == "data_disk_used_percent":
			host.dataDiskUsedPercent = newValue
		default:
			fmt.Printf("xxx")
		}
	} else {
		newHost := Host{}
		hosts[ident] = &newHost
	}
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
