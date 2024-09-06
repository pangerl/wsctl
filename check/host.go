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
	"net/http"
	"os"
	"strconv"
)

func tableRender(hosts map[string]*Host) {
	tabletitle := []string{"IP", "CPU使用率", "内存使用率", "内存大小"}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tabletitle)
	for ident, data := range hosts {
		tabledata := []string{ident, FormatToPercentage(data.CpuUsageActive),
			FormatToPercentage(data.MemUsedPercent), data.MemTotal}
		table.Append(tabledata)
	}
	table.Render()
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

func getHostData(baseUrl, key string) {
	url := baseUrl + "/api/v1/query?query=" + key
	response, _ := queryVictoriaMetrics(url)
	if response.Status != "success" {
		log.Printf("查询报错，Status: %s \n", response.Status)
		return
	}
	for _, result := range response.Data.Result {
		ident := result.Metric["ident"]
		host := getHost(ident)
		switch {
		case key == "cpu_usage_active":
			host.CpuUsageActive = result.Value[1].(string)
		case key == "mem_used_percent":
			host.MemUsedPercent = result.Value[1].(string)
		case key == "mem_total":
			host.MemTotal = result.Value[1].(string)
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

func FormatToPercentage(original string) string {
	// 将字符串转换为浮点数
	value, err := strconv.ParseFloat(original, 64)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return ""
	}

	// 格式化浮点数，只保留小数点后1位，并添加百分号
	result := fmt.Sprintf("%.1f%%", value)

	return result
}
