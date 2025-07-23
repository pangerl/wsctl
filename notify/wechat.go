// Package notify @Author lanpang
// @Date 2024/8/8 下午5:14:00
// @Desc
package notify

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"vhagar/logger"
)

type WeChatMarkdown struct {
	MsgType  string    `json:"msgtype"`
	Markdown *Markdown `json:"markdown"`
}

type Markdown struct {
	Content string `json:"content"`
}

func sendWecom(markdown *WeChatMarkdown, robotKey, proxyURL string) error {
	jsonStr, _ := json.Marshal(markdown)
	robotURL := wechatRobotURL + robotKey

	req, err := http.NewRequest("POST", robotURL, bytes.NewBuffer(jsonStr))
	if err != nil {
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
		logger.Logger.Errorw("推送企微机器人失败", "err", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Logger.Errorw("Failed info", "err", err)
		}
	}(resp.Body)
	logger.Logger.Warnw("推送企微机器人 response Status", "status", resp.Status)
	//log.Print("response Headers:", resp.Header)
	return nil
}
