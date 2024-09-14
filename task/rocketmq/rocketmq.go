// Package rocketmq  @Author lanpang
// @Date 2024/9/10 下午6:13:00
// @Desc
package rocketmq

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"vhagar/config"
	"vhagar/notifier"
	"vhagar/task"
)

func GetRocketMQ() *RocketMQ {
	cfg := config.Config
	rocketmq := newRocketMQ(cfg)
	rocketmq.InitData()
	return rocketmq
}

func (rocketmq *RocketMQ) Check() {
	task.EchoPrompt("开始巡检 RocketMQ 信息")
	if rocketmq.Report {
		rocketmq.ReportRobot(rocketmq.Global.Duration)
		return
	}
	rocketmq.TableRender()
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

func (rocketmq *RocketMQ) TableRender() {
	// 输出RocketMQ巡检报告
	tabletitle := []string{"Broker Name", "Role", "Version", "IP", "今天生产总数", "今天消费总数", "运行时间", "磁盘可用空间"}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tabletitle)
	table.SetAutoMergeCellsByColumnIndex([]int{0, 0})
	table.SetRowLine(true)
	for _, broker := range rocketmq.BrokerList {
		tabledata := []string{broker.name, broker.role, broker.version, broker.addr,
			strconv.Itoa(broker.todayProduceCount), strconv.Itoa(broker.todayConsumeCount), broker.runTime, broker.useDisk}
		table.Append(tabledata)
	}
	caption := fmt.Sprintf("Broker 实例数: %d.", len(rocketmq.BrokerList))
	table.SetCaption(true, caption)
	table.Render()
}

func (rocketmq *RocketMQ) InitData() {
	// 获取RocketMQ集群信息
	clusterdata, _ := GetMQDetail(rocketmq.RocketmqDashboard)
	for brokername, brokerdata := range clusterdata.BrokerServer {
		for role, broker := range brokerdata {
			rocketmq.BrokerList = append(rocketmq.BrokerList, &BrokerDetail{
				name:              brokername,
				role:              getRole(role),
				version:           broker.BrokerVersionDesc,
				addr:              clusterdata.ClusterInfo.BrokerAddrTable[brokername].BrokerAddrs[role],
				runTime:           formatRunTime(broker.RunTime),
				useDisk:           formatUseDisk(broker.CommitLogDirCapacity),
				todayProduceCount: convertAndCalculate(broker.MsgPutTotalTodayNow, broker.MsgPutTotalTodayMorning),
				todayConsumeCount: convertAndCalculate(broker.MsgGetTotalTodayNow, broker.MsgGetTotalTodayMorning),
			})
		}
	}
}

func formatRunTime(runTime string) string {
	cleanedStr := strings.Trim(runTime, "[] ")
	// 使用逗号分割字符串
	items := strings.Split(cleanedStr, ",")
	return items[0]
}

func formatUseDisk(useDisk string) string {
	// 使用逗号分割字符串
	items := strings.Split(useDisk, ",")
	total := strings.TrimSpace(strings.Split(items[0], ":")[1])
	free := strings.TrimSpace(strings.Split(items[1], ":")[1])
	return free + "/" + total
}

func GetMQDetail(mqDashboard string) (result ClusterData, err error) {
	// 第一步：发送HTTP请求到RocketMQ Dashboard接口
	url := mqDashboard + "/cluster/list.query"
	body := task.DoRequest(url)
	// 第二步：解析JSON响应
	var responseData ResponseData
	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Printf("E! fail to unmarshal JSON response: %v", err)
	}
	result = responseData.Data

	return result, err
}

func mqDetailMarkdown(rocketmq *RocketMQ, ProjectName string) *notifier.WeChatMarkdown {

	brokerList := rocketmq.BrokerList
	var builder strings.Builder

	// 组装巡检内容
	builder.WriteString("# RocketMQ 巡检 \n")
	builder.WriteString("**项目名称：**<font color='info'>" + ProjectName + "</font>\n")
	builder.WriteString("**巡检时间：**<font color='info'>" + time.Now().Format("2006-01-02") + "</font>\n")
	builder.WriteString("**巡检内容：**\n\n")
	builder.WriteString("**Broker 健康数：**<font color='info'>" + strconv.Itoa(len(brokerList)) + "</font>\n")
	builder.WriteString("========================\n")
	for _, broker := range brokerList {
		builder.WriteString("## Broker Name：<font color='info'>" + broker.name + "</font>\n")
		builder.WriteString("### " + broker.role + "\n")
		builder.WriteString("> Broker 版本：<font color='info'>" + broker.version + "</font>\n")
		builder.WriteString("> Broker 地址：<font color='info'>" + broker.addr + "</font>\n")
		builder.WriteString("> 今天生产总数：<font color='info'>" + strconv.Itoa(broker.todayProduceCount) + "</font>\n")
		builder.WriteString("> 今天消费总数：<font color='info'>" + strconv.Itoa(broker.todayConsumeCount) + "</font>\n")
		builder.WriteString("> 运行时间：<font color='info'>" + broker.runTime + "</font>\n")
		builder.WriteString("> 磁盘可用空间：<font color='info'>" + broker.useDisk + "</font>")
		builder.WriteString("\n\n")
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

func convertAndCalculate(str1, str2 string) int {
	num1, err := strconv.Atoi(str1)
	if err != nil {
		log.Println("E! fail to str to int", err)
		return -1
	}
	num2, err := strconv.Atoi(str2)
	if err != nil {
		log.Println("E! fail to str to int", err)
		return -1
	}

	return num1 - num2
}

func getRole(role string) string {
	if role == "0" {
		return "Master"
	}
	return "Slave"
}
