// Package inspect @Author lanpang
// @Date 2024/8/13 下午2:17:00
// @Desc
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type WeChatMarkdown struct {
	MsgType  string    `json:"msgtype"`
	Markdown *Markdown `json:"markdown"`
}

type Markdown struct {
	Content string `json:"content"`
}

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
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	// 第二步：解析JSON响应
	var responseData ResponseData
	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Fatalf("Failed to unmarshal JSON response: %v", err)
	}
	result = responseData.Data

	return result, err
}

func MQDetailToMarkdown(data ClusterData, ProjectName string) *WeChatMarkdown {

	var builder strings.Builder
	// 组装巡检内容
	builder.WriteString("# RocketMQ 巡检报告 \n")
	builder.WriteString("**项目名称：**<font color='info'>" + ProjectName + "</font>\n")
	builder.WriteString("**巡检时间：**<font color='info'>" + time.Now().Format("2006-01-02") + "</font>\n")
	builder.WriteString("**巡检内容：**\n")

	for brokername, brokerdata := range data.BrokerServer {
		builder.WriteString("> Broker 名称：<font color='info'>" + brokername + "</font>\n")
		for role, broker := range brokerdata {
			builder.WriteString("> Broker 角色：<font color='info'>" + role + "</font>\n")
			builder.WriteString("> Broker 版本：<font color='info'>" + broker.BrokerVersionDesc + "</font>\n")
			builder.WriteString("> Broker 地址：<font color='info'>" + data.ClusterInfo.BrokerAddrTable[brokername].BrokerAddrs[role] + "</font>\n")
			builder.WriteString("> 运行时间：<font color='info'>" + broker.RunTime + "</font>\n")
			builder.WriteString("> 磁盘使用量：<font color='info'>" + broker.CommitLogDirCapacity + "</font>\n")
		}
		builder.WriteString("==================\n")
	}

	markdown := &WeChatMarkdown{
		MsgType: "markdown",
		Markdown: &Markdown{
			Content: builder.String(),
		},
	}

	return markdown
}
func SendWecom(markdown *WeChatMarkdown, robotKey, proxyURL string) error {

	jsonStr, _ := json.Marshal(markdown)
	wechatRobotURL := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + robotKey

	req, err := http.NewRequest("POST", wechatRobotURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return err
		}
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed info: %s \n", err)
		}
	}(resp.Body)
	log.Print("推送企微机器人 response Status:", resp.Status)
	//log.Print("response Headers:", resp.Header)
	return nil
}

func main() {
	clusterdata, _ := GetMQDetail()
	markdown := MQDetailToMarkdown(clusterdata, "千金药业")
	_ = SendWecom(markdown, "8ab989f5-d86f-4c99-b95f-5d23a30eb351", "")
}
