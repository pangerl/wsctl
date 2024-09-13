// Package rocketmq  @Author lanpang
// @Date 2024/9/10 下午6:13:00
// @Desc
package rocketmq

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"vhagar/config"
	"vhagar/notifier"
)

func Check() {
	cfg := config.Config
	rocketmq := newRocketMQ(cfg)
	initData(rocketmq)
	// 发送巡检报告
	rocketmq.ReportRobot(0)
}

func (rocketmq *RocketMQ) ReportRobot(duration time.Duration) {
	// 发送巡检报告
	markdown := mqDetailMarkdown(rocketmq, rocketmq.ProjectName)
	log.Println("任务等待时间", duration)
	time.Sleep(duration)
	for _, robotkey := range rocketmq.Notifier["rocketmq"].Robotkey {
		_ = notifier.SendWecom(markdown, robotkey, rocketmq.ProxyURL)
	}

}

func initData(rocketmq *RocketMQ) {
	// 获取RocketMQ集群信息
	clusterdata, _ := GetMQDetail(rocketmq.RocketmqDashboard)
	rocketmq.ClusterData = clusterdata
}

func GetMQDetail(mqDashboard string) (result ClusterData, err error) {
	// 第一步：发送HTTP请求到RocketMQ Dashboard接口
	url := mqDashboard + "/cluster/list.query"
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

func mqDetailMarkdown(rocketmq *RocketMQ, ProjectName string) *notifier.WeChatMarkdown {

	data := rocketmq.ClusterData
	var builder strings.Builder
	var brokercount int

	brokercount = GetBrokerCount(data)
	// 组装巡检内容
	builder.WriteString("# RocketMQ 巡检 \n")
	builder.WriteString("**项目名称：**<font color='info'>" + ProjectName + "</font>\n")
	builder.WriteString("**巡检时间：**<font color='info'>" + time.Now().Format("2006-01-02") + "</font>\n")
	builder.WriteString("**巡检内容：**\n\n")
	builder.WriteString("**Broker 健康数：**<font color='info'>" + strconv.Itoa(brokercount) + "</font>\n")
	builder.WriteString("========================\n")

	for brokername, brokerdata := range data.BrokerServer {
		builder.WriteString("## Broker Name：<font color='info'>" + brokername + "</font>\n")
		for role, broker := range brokerdata {
			var produceCount, consumeCount int
			produceCount, _ = convertAndCalculate(broker.MsgPutTotalTodayNow, broker.MsgPutTotalTodayMorning)
			consumeCount, _ = convertAndCalculate(broker.MsgGetTotalTodayNow, broker.MsgGetTotalTodayMorning)
			builder.WriteString("### " + getRole(role) + "\n")
			builder.WriteString("> Broker 版本：<font color='info'>" + broker.BrokerVersionDesc + "</font>\n")
			builder.WriteString("> Broker 地址：<font color='info'>" + data.ClusterInfo.BrokerAddrTable[brokername].BrokerAddrs[role] + "</font>\n")
			builder.WriteString("> 今天生产总数：<font color='info'>" + strconv.Itoa(produceCount) + "</font>\n")
			builder.WriteString("> 今天消费总数：<font color='info'>" + strconv.Itoa(consumeCount) + "</font>\n")
			builder.WriteString("> 运行时间：<font color='info'>" + broker.RunTime + "</font>\n")
			builder.WriteString("> 磁盘使用量：<font color='info'>" + broker.CommitLogDirCapacity + "</font>")
			builder.WriteString("\n\n")
		}
		builder.WriteString("========================\n\n")
	}

	markdown := &notifier.WeChatMarkdown{
		MsgType: "markdown",
		Markdown: &notifier.Markdown{
			Content: builder.String(),
		},
	}

	return markdown
}

func GetBrokerCount(data ClusterData) int {
	var brokercount int
	for _, brokerdata := range data.BrokerServer {
		brokercount += len(brokerdata)
	}
	return brokercount
}

func convertAndCalculate(str1, str2 string) (int, error) {
	num1, err := strconv.Atoi(str1)
	if err != nil {
		return 0, err
	}

	num2, err := strconv.Atoi(str2)
	if err != nil {
		return 0, err
	}

	return num1 - num2, nil
}

func getRole(role string) string {
	if role == "0" {
		return "Master"
	}
	return "Slave"
}
