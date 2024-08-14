// Package inspect @Author lanpang
// @Date 2024/8/13 下午2:17:00
// @Desc
package inspect

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"vhagar/notifier"
)

// BrokerData 定义用于解析RocketMQ Dashboard返回数据的结构体
type BrokerData struct {
	RunTime              string `json:"runtime"`
	CommitLogDirCapacity string `json:"commitLogDirCapacity"`
	BrokerVersionDesc    string `json:"brokerVersionDesc"`
}

type Broker struct {
	BrokerName  string            `json:"brokerName"`
	BrokerAddrs map[string]string `json:"brokerAddrs"`
}

type ClusterInfo struct {
	BrokerAddrTable map[string]Broker `json:"brokerAddrTable"`
}

type ClusterData struct {
	BrokerServer map[string]map[string]BrokerData `json:"brokerServer"`
	ClusterInfo  ClusterInfo                      `json:"clusterInfo"`
}

type ResponseData struct {
	Status int         `json:"status"`
	Data   ClusterData `json:"data"`
}

func GetMQDetail() (result ClusterData, err error) {
	// 第一步：发送HTTP请求到RocketMQ Dashboard接口
	url := "http://192.9.253.205:8081/cluster/list.query"
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed info : %v", err)
		}
	}(resp.Body)

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
	}
	// 第二步：解析JSON响应
	var responseData ResponseData
	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Printf("Failed to unmarshal JSON response: %v", err)
	}
	result = responseData.Data

	return result, err
}

func MQDetailToMarkdown(data ClusterData, ProjectName string) *notifier.WeChatMarkdown {

	var builder strings.Builder
	var brokercount int
	// 组装巡检内容
	builder.WriteString("# RocketMQ 巡检 \n")
	builder.WriteString("**项目名称：**<font color='info'>" + ProjectName + "</font>\n")
	builder.WriteString("**巡检时间：**<font color='info'>" + time.Now().Format("2006-01-02") + "</font>\n")
	builder.WriteString("**巡检内容：**\n\n")

	for brokername, brokerdata := range data.BrokerServer {
		builder.WriteString("## Broker Name：<font color='info'>" + brokername + "</font>\n")
		for role, broker := range brokerdata {
			brokercount += 1
			builder.WriteString("### " + getRole(role) + "\n")
			builder.WriteString("> Broker 版本：<font color='info'>" + broker.BrokerVersionDesc + "</font>\n")
			builder.WriteString("> Broker 地址：<font color='info'>" + data.ClusterInfo.BrokerAddrTable[brokername].BrokerAddrs[role] + "</font>\n")
			builder.WriteString("> 运行时间：<font color='info'>" + broker.RunTime + "</font>\n")
			builder.WriteString("> 磁盘使用量：<font color='info'>" + broker.CommitLogDirCapacity + "</font>")
			builder.WriteString("\n\n")
		}
		builder.WriteString("========================\n\n")
	}

	builder.WriteString("**Broker 健康数：**<font color='info'>" + strconv.Itoa(brokercount) + "</font>\n")

	markdown := &notifier.WeChatMarkdown{
		MsgType: "markdown",
		Markdown: &notifier.Markdown{
			Content: builder.String(),
		},
	}

	return markdown
}
