// Package rocketmq  @Author lanpang
// @Date 2024/9/10 下午6:13:00
// @Desc
package rocketmq

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"vhagar/config"
	"vhagar/libs"
	"vhagar/notify"
	"vhagar/task"

	"github.com/olekukonko/tablewriter"
)

const (
	roleMaster = "0"
)

func init() {
	task.Add(taskName, func() task.Tasker {
		return NewRocketMQ(config.Config)
	})
}

func (rocketmq *RocketMQ) Check() {
	//task.EchoPrompt("开始巡检 RocketMQ 信息")
	if config.Config.Report {
		rocketmq.ReportRobot()
		return
	}
	rocketmq.TableRender()
}

func (rocketmq *RocketMQ) ReportRobot() {
	brokerList := rocketmq.BrokerMap
	var builder strings.Builder

	// 组装巡检内容
	builder.WriteString("# RocketMQ 巡检 \n")
	builder.WriteString("**项目名称：**<font color='info'>" + config.Config.ProjectName + "</font>\n")
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

	markdown := &notify.WeChatMarkdown{
		MsgType: "markdown",
		Markdown: &notify.Markdown{
			Content: builder.String(),
		},
	}
	notify.Send(markdown, taskName)
}

func (rocketmq *RocketMQ) TableRender() {
	// 输出RocketMQ巡检报告
	tabletitle := []string{"Broker Name", "Role", "Version", "IP", "今天生产总数", "今天消费总数", "运行时间", "磁盘.可用空间/总空间"}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tabletitle)
	table.SetAutoMergeCellsByColumnIndex([]int{0, 0})
	table.SetRowLine(true)
	for _, broker := range rocketmq.BrokerMap {
		tabledata := []string{broker.name, broker.role, broker.version, broker.addr,
			strconv.Itoa(broker.todayProduceCount), strconv.Itoa(broker.todayConsumeCount), broker.runTime, broker.useDisk}
		table.Append(tabledata)
	}
	caption := fmt.Sprintf("Broker 实例数: %d.", len(rocketmq.BrokerMap))
	table.SetCaption(true, caption)
	table.Render()
}

func (rocketmq *RocketMQ) Gather() {
	// 获取RocketMQ集群信息
	clusterdata, _ := GetMQDetail(rocketmq.RocketmqDashboard, rocketmq.Username, rocketmq.Password)
	for brokername, brokerdata := range clusterdata.BrokerServer {
		for role, broker := range brokerdata {
			addr := clusterdata.ClusterInfo.BrokerAddrTable[brokername].BrokerAddrs[role]
			_broker := rocketmq.getBroker(addr)
			_broker.name = brokername
			_broker.role = getRole(role)
			_broker.version = broker.BrokerVersionDesc
			_broker.addr = addr
			_broker.runTime = formatRunTime(broker.RunTime)
			_broker.useDisk = formatUseDisk(broker.CommitLogDirCapacity)
			_broker.todayProduceCount = convertAndCalculate(broker.MsgPutTotalTodayNow, broker.MsgPutTotalTodayMorning)
			_broker.todayConsumeCount = convertAndCalculate(broker.MsgGetTotalTodayNow, broker.MsgGetTotalTodayMorning)
		}
	}
}

func (rocketmq *RocketMQ) getBroker(addr string) *BrokerDetail {
	if broker, exists := rocketmq.BrokerMap[addr]; exists {
		return broker
	}
	newBroker := BrokerDetail{}
	rocketmq.BrokerMap[addr] = &newBroker
	return &newBroker
}

func formatRunTime(runTime string) string {
	cleanedStr := strings.Trim(runTime, "[] ")
	// 使用逗号分割字符串
	items := strings.Split(cleanedStr, ",")
	return items[0]
}

func formatUseDisk(useDisk string) string {
	items := strings.Split(useDisk, ",")
	if len(items) < 2 {
		return useDisk
	}
	totalParts := strings.Split(items[0], ":")
	freeParts := strings.Split(items[1], ":")
	if len(totalParts) < 2 || len(freeParts) < 2 {
		return useDisk
	}
	total := strings.TrimSpace(totalParts[1])
	free := strings.TrimSpace(freeParts[1])
	return free + "/" + total
}

func GetMQDetail(mqDashboard, username, password string) (result ClusterData, err error) {
	// 新增：登录获取 cookie
	loginUrl := mqDashboard + "/login/login.do"
	clusterUrl := mqDashboard + "/cluster/list.query"

	client := &http.Client{}

	// 1. 登录获取 cookie
	loginData := url.Values{}
	loginData.Set("username", username)
	loginData.Set("password", password)

	loginReq, err := http.NewRequest("POST", loginUrl, strings.NewReader(loginData.Encode()))
	if err != nil {
		libs.Logger.Errorf("E! fail to create login request: %v", err)
		return result, err
	}
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	loginResp, err := client.Do(loginReq)
	if err != nil {
		libs.Logger.Errorf("E! fail to login: %v", err)
		return result, err
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode != 200 {
		libs.Logger.Errorf("E! login failed, status: %d", loginResp.StatusCode)
		return result, fmt.Errorf("login failed, status: %d", loginResp.StatusCode)
	}

	// 提取 cookie
	cookies := loginResp.Cookies()
	if len(cookies) == 0 {
		libs.Logger.Errorf("E! login response has no cookies")
		return result, fmt.Errorf("login response has no cookies")
	}

	// 2. 带 cookie 请求 cluster/list.query
	clusterReq, err := http.NewRequest("GET", clusterUrl, nil)
	if err != nil {
		libs.Logger.Errorf("E! fail to create cluster request: %v", err)
		return result, err
	}
	for _, c := range cookies {
		clusterReq.AddCookie(c)
	}

	clusterResp, err := client.Do(clusterReq)
	if err != nil {
		libs.Logger.Errorf("E! fail to request cluster info: %v", err)
		return result, err
	}
	defer clusterResp.Body.Close()

	if clusterResp.StatusCode != 200 {
		libs.Logger.Errorf("E! cluster info failed, status: %d", clusterResp.StatusCode)
		return result, fmt.Errorf("cluster info failed, status: %d", clusterResp.StatusCode)
	}

	body, err := io.ReadAll(clusterResp.Body)
	if err != nil {
		libs.Logger.Errorf("E! fail to read cluster response: %v", err)
		return result, err
	}

	// 解析JSON响应
	var responseData ResponseData
	if err := json.Unmarshal(body, &responseData); err != nil {
		libs.Logger.Errorf("E! fail to unmarshal JSON response: %v", err)
		return result, err
	}
	result = responseData.Data
	return result, nil
}

func convertAndCalculate(str1, str2 string) int {
	num1, err := strconv.Atoi(str1)
	if err != nil {
		return 0
	}
	num2, err := strconv.Atoi(str2)
	if err != nil {
		return 0
	}
	return num1 - num2
}

func getRole(role string) string {
	if role == roleMaster {
		return "Master"
	}
	return "Slave"
}
