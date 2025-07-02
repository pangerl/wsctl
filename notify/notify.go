// Package notify @Author lanpang
// @Date 2024/9/20 下午4:21:00
// @Desc
package notify

import (
	"time"
	"vhagar/config"
	"vhagar/libs"
)

type Notifier struct {
	Robotkey []string `json:"robotkey"`
}

const wechatRobotURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="

func Send(markdown *WeChatMarkdown, taskName string) {
	libs.Logger.Infow("任务等待时间", "duration", config.Config.Duration)
	time.Sleep(config.Config.Duration)
	robotkey := getRobotkey(taskName)
	//fmt.Println("robotkey", robotkey)
	for _, robotkey := range robotkey {
		err := sendWecom(markdown, robotkey, config.Config.ProxyURL)
		if err != nil {
			libs.Logger.Errorw("发送失败", "err", err)
		}
	}
}

func getRobotkey(taskName string) []string {
	if notifier, ok := config.Config.Notify.Notifier[taskName]; ok {
		return notifier.Robotkey
	}
	return config.Config.Notify.Robotkey
}
