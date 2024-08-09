// Package inspect @Author lanpang
// @Date 2024/8/8 下午5:14:00
// @Desc
package inspect

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

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
