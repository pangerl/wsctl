// Package inspect @Author lanpang
// @Date 2024/8/7 下午3:43:00
// @Desc
package inspect

import (
	"log"
	"vhagar/notifier"
)

func RocketmqTask(inspect *Inspect) {
	log.Print("启动 rocketmq 巡检任务")
	clusterdata, _ := GetMQDetail(inspect.Rocketmq.RocketmqDashboard)
	markdown := mqDetailMarkdown(clusterdata, inspect.ProjectName)
	for _, robotkey := range inspect.Notifier["rocketmq"].Robotkey {
		_ = notifier.SendWecom(markdown, robotkey, inspect.ProxyURL)
	}
}
