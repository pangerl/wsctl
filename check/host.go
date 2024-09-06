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
	color := tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor}
	tableColor := []tablewriter.Colors{color, color, color, color}
	for ident, data := range hosts {
		tabledata := []string{ident, formatToPercentage(data.CpuUsageActive),
			formatToPercentage(data.MemUsedPercent), bytesToGB(data.MemTotal)}
		// 异常标红
		if isAlarm(data) {
			table.Rich(tabledata, tableColor)
		}
		table.Append(tabledata)
	}
	table.Render()
}

func isAlarm(host *Host) bool {
	if host.CpuUsageActive > 90 {
		return true
	}
	if host.MemUsedPercent > 85 {
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
			host.CpuUsageActive, _ = strconv.ParseFloat(result.Value[1].(string), 64)
		case key == "mem_used_percent":
			host.MemUsedPercent, _ = strconv.ParseFloat(result.Value[1].(string), 64)
		case key == "mem_total":
			host.MemTotal, _ = strconv.ParseFloat(result.Value[1].(string), 64)
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
