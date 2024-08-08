// Package inspect @Author lanpang
// @Date 2024/8/8 下午5:14:00
// @Desc
package inspect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func SendWecom(markdown *WeChatMarkdown, robotKey, proxyurl string) (err error) {

	data, err := json.Marshal(markdown)
	CheckErr(err)
	// var wechatRobotURL string
	wechatRobotURL := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + robotKey

	req, err := http.NewRequest(
		"POST",
		wechatRobotURL,
		bytes.NewBuffer(data))

	CheckErr(err)

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	if proxyurl != "" {
		proxy, err := url.Parse(getproxyurl(proxyurl))
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

	CheckErr(err)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		CheckErr(err)
	}(resp.Body)
	fmt.Println("推送企微机器人 response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	return err
}

func TransformToMarkdown(i *Inspect, users []string, dateNow time.Time) (markdown *WeChatMarkdown, err error) {

	var buffer bytes.Buffer
	var isalert bool
	isalert = false

	// buffer.WriteString("# 每日巡检报告 \n")
	buffer.WriteString(fmt.Sprintf("# 每日巡检报告 %s \n", i.Version))
	buffer.WriteString(fmt.Sprintf("**项目名称：**<font color='info'>%s</font>\n", i.ProjectName))
	buffer.WriteString(fmt.Sprintf("**巡检时间：**<font color='info'>%s</font>\n", dateNow.Format("2006-01-02")))
	buffer.WriteString("**巡检内容：**\n")

	for _, corp := range i.Corp {
		buffer.WriteString(fmt.Sprintf("> 企业名称：<font color='info'>%s</font>\n", corp.CorpName))
		if corp.Convenabled {
			buffer.WriteString(fmt.Sprintf("> 昨天拉取会话数：<font color=%s>%d</font>\n",
				getcolor(corp.MessageNum, false), corp.MessageNum))
			if corp.MessageNum == 0 {
				isalert = true
			}
		}
		buffer.WriteString(fmt.Sprintf("> 员工数统计：<font color='info'>%d</font>\n", corp.UserNum))
		buffer.WriteString(fmt.Sprintf("> 客户数统计：<font color='info'>%d</font>\n", corp.CustomerNum))
		buffer.WriteString(fmt.Sprintf("> 日活数统计：<font color=%s>%d</font>\n", getcolor(corp.DauNum, true), corp.DauNum))
		buffer.WriteString(fmt.Sprintf("> 周活数统计：<font color=%s>%d</font>\n", getcolor(corp.WauNum, true), corp.WauNum))
		buffer.WriteString(fmt.Sprintf("> 月活数统计：<font color=%s>%d</font>\n", getcolor(corp.MauNum, true), corp.MauNum))
		buffer.WriteString("==================\n")
	}

	if isalert {
		buffer.WriteString("\n<font color='red'>**注意！巡检结果异常！**</font>" + calluser(users))
	}

	markdown = &WeChatMarkdown{
		MsgType: "markdown",
		Markdown: &Markdown{
			Content: buffer.String(),
		},
	}
	return
}

func calluser(users []string) string {
	var result string
	if len(users) == 0 {
		return result
	}
	for _, user := range users {
		result += fmt.Sprintf("<@%s>", user)
	}
	return result
}

func getcolor(num int64, ignore bool) string {
	color := "info"
	if num == 0 {
		if ignore {
			color = "warning"
		} else {
			color = "red"
		}
	}
	return color
}
func getproxyurl(proxy string) string {
	return fmt.Sprintf("http://%s", proxy)
}
