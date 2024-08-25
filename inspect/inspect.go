// Package inspect @Author lanpang
// @Date 2024/8/7 下午3:43:00
// @Desc
package inspect

import (
	"github.com/olivere/elastic/v7"
	"log"
	"time"
	"vhagar/libs"
)

var isalert = false

func NewTenant(ecfg, pcfg DB, projectname, proxyurl string, corp []*Corp, userlist, robotkey []string) *Tenant {
	log.Println("初始化 Tenant 对象")
	esClient, _ := libs.NewESClient(ecfg)
	pgClient, _ := libs.NewPGClient(pcfg)
	defer func() {
		if pgClient != nil {
			pgClient.Close()
		}
		if esClient != nil {
			esClient.Stop()
		}
	}()

	tenant := &Tenant{
		ProjectName: projectname,
		ProxyURL:    proxyurl,
		Version:     "v4.5",
		Corp:        corp,
		ESClient:    esClient,
		PGClient:    pgClient,
		Userlist:    userlist,
		Robotkey:    robotkey,
	}
	return tenant

}

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
