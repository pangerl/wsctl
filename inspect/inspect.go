// Package inspect @Author lanpang
// @Date 2024/8/7 下午3:43:00
// @Desc
package inspect

import (
	"time"

	"github.com/olivere/elastic/v7"
)

var isalert = false

func CurrentMessageNum(client *elastic.Client, corpid string, dateNow time.Time) int64 {
	// 统计今天的会话数
	startTime := getZeroTime(dateNow).UnixNano() / 1e6
	endTime := dateNow.UnixNano() / 1e6
	messagenum, _ := countMessageNum(client, corpid, startTime, endTime)
	return messagenum
}

//func mqTask() {
//	log.Print("启动 rocketmq 巡检任务")
//	clusterdata, _ := inspect.GetMQDetail(CONFIG.Rocketmq.RocketmqDashboard)
//	markdown := inspect.MQDetailToMarkdown(clusterdata, CONFIG.ProjectName)
//	for _, robotkey := range common.CONFIG.Inspection.Robotkey {
//		_ = notifier.SendWecom(markdown, robotkey, common.CONFIG.ProxyURL)
//	}
//}
