// Package check @Author lanpang
// @Date 2024/9/6 上午11:31:00
// @Desc
package check

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
)

func tableRender(hosts map[string]*Host) {
	tabletitle := []string{"IP", "CPU使用率", "内存使用率", "内存大小", "系统盘使用率", "数据盘使用率"}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tabletitle)
	color := tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor}
	tableColor := []tablewriter.Colors{color, color, color, color}
	// ident 排序
	identList := ipSort(hosts)
	for _, ident := range identList {
		data := hosts[ident]
		tabledata := []string{ident, formatToPercentage(data.CpuUsageActive),
			formatToPercentage(data.MemUsedPercent), bytesToGB(data.MemTotal),
			formatToPercentage(data.DiskUsedPercent["/"]), formatToPercentage(data.DiskUsedPercent["/data"])}
		// 异常标红
		if isAlarm(data) {
			table.Rich(tabledata, tableColor)
		}
		table.Append(tabledata)
	}
	table.Render()
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

func getHostData(baseUrl string, key ...string) {
	var url string
	if len(key) == 1 {
		url = baseUrl + "/api/v1/query?query=" + key[0]
	} else if len(key) == 2 {
		url = fmt.Sprintf("%s/api/v1/query?query=%s{path='%s'}", baseUrl, key[0], key[1])
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
		host := getHost(ident)
		switch {
		case key[0] == "cpu_usage_active":
			host.CpuUsageActive, _ = strconv.ParseFloat(result.Value[1].(string), 64)
		case key[0] == "mem_used_percent":
			host.MemUsedPercent, _ = strconv.ParseFloat(result.Value[1].(string), 64)
		case key[0] == "mem_total":
			host.MemTotal, _ = strconv.ParseFloat(result.Value[1].(string), 64)
		case key[0] == "disk_used_percent":
			host.DiskUsedPercent = make(map[string]float64)
			host.DiskUsedPercent[key[1]], _ = strconv.ParseFloat(result.Value[1].(string), 64)
		default:
			fmt.Printf("xxx")
		}
	}
}

func getHost(ident string) *Host {
	if host, exists := hosts[ident]; exists {
		return host
	}
	newHost := Host{}
	hosts[ident] = &newHost
	return hosts[ident]
}

func formatToPercentage(value float64) string {
	result := fmt.Sprintf("%.1f%%", value)
	return result
}

func bytesToGB(bytes float64) string {
	gb := bytes / (1024 * 1024 * 1024)
	result := fmt.Sprintf("%.2f GB", gb)
	return result
}
