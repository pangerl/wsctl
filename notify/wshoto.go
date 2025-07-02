// Package notify @Author lanpang
// @Date 2024/9/3 下午2:08:00
// @Desc
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"vhagar/libs"
)

type InspectBody struct {
	JobType string `json:"jobtype"`
	Data    []any  `json:"data"`
}

var domain = "http://10.229.3.2:8088"

func SendWshoto(inspect *InspectBody, proxyURL string) error {

	jsonStr, err := json.Marshal(inspect)
	//fmt.Println("jsonStr：", jsonStr)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	wshotoURL := domain + "/stage-api/project/statistics/push"

	req, err := http.NewRequest("POST", wshotoURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		libs.Logger.Errorw("Failed info", "err", err)
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
		libs.Logger.Errorw("Failed info", "err", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			libs.Logger.Errorw("Failed info", "err", err)
		}
	}(resp.Body)
	libs.Logger.Infow("推送 wshoto response Status", "status", resp.Status)
	//log.Print("response Headers:", resp.Header)
	return nil
}
