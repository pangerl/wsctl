// Package notify @Author lanpang
// @Date 2024/9/20 下午4:21:00
// @Desc
package notify

import (
	"time"
	"vhagar/config"
	"vhagar/logger"
)

type Notifier struct {
	Robotkey []string `json:"robotkey"`
}

const wechatRobotURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="

func Send(markdown *WeChatMarkdown, taskName string) {
	logger.Logger.Infow("任务等待时间", "duration", config.Config.Global.Duration)
	time.Sleep(config.Config.Global.Duration)
	robotkey := getRobotkey(taskName)
	//fmt.Println("robotkey", robotkey)
	for _, robotkey := range robotkey {
		err := sendWecom(markdown, robotkey, config.Config.Global.ProxyURL)
		if err != nil {
			logger.Logger.Errorw("发送失败", "err", err)
		}
	}
}

func getRobotkey(taskName string) []string {
	if notifier, ok := config.Config.Global.Notify.Notifier[taskName]; ok {
		return notifier.Robotkey
	}
	return config.Config.Global.Notify.Robotkey
}
