// Package inspect @Author lanpang
// @Date 2024/8/7 下午3:43:00
// @Desc
package inspect

import (
	"log"
	"time"
	"vhagar/notifier"
)

var isalert = false

func TenantTask(tenant *Tenant, duration time.Duration) {
	// 当前时间
	dateNow := time.Now()
	log.Print("启动企微租户信息巡检任务")
	for _, corp := range tenant.Corp {
		// fmt.Println(corp.Corpid)
		if tenant.PGClient != nil {
			// 获取租户名
			tenant.SetCorpName(corp.Corpid)
			// 获取用户数
			tenant.SetUserNum(corp.Corpid)
			// 获取客户群
			tenant.SetCustomerGroupNum(corp.Corpid)
			// 获取客户群人数
			tenant.SetCustomerGroupUserNum(corp.Corpid)
		}
		if tenant.ESClient != nil {
			// 获取客户数
			tenant.SetCustomerNum(corp.Corpid)
			// 获取活跃数
			tenant.SetActiveNum(corp.Corpid, dateNow)
			// 获取会话数
			if corp.Convenabled {
				tenant.SetMessageNum(corp.Corpid, dateNow)
			}
		}
		//fmt.Println(*corp)
	}
	// 发送巡检报告
	markdownList := tenantNotifier(tenant, dateNow)
	log.Println("任务等待时间", duration)
	time.Sleep(duration)
	for _, markdown := range markdownList {
		for _, robotkey := range tenant.Robotkey {
			err := notifier.SendWecom(markdown, robotkey, tenant.ProxyURL)
			if err != nil {
				return
			}
		}
	}

}

func RocketmqTask(t *Tenant) {
	log.Print("启动 rocketmq 巡检任务")
	clusterdata, _ := GetMQDetail(t.Rocketmq.RocketmqDashboard)
	markdown := mqDetailToMarkdown(clusterdata, t.ProjectName)
	for _, robotkey := range t.Robotkey {
		_ = notifier.SendWecom(markdown, robotkey, t.ProxyURL)
	}
}
