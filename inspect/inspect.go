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
	markdownList := tenantNotifier(tenant)
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

func RocketmqTask(tenant *Tenant) {
	log.Print("启动 rocketmq 巡检任务")
	clusterdata, _ := GetMQDetail(tenant.Rocketmq.RocketmqDashboard)
	markdown := mqDetailMarkdown(clusterdata, tenant.ProjectName)
	for _, robotkey := range tenant.Robotkey {
		_ = notifier.SendWecom(markdown, robotkey, tenant.ProxyURL)
	}
}

func DorisTask(tenant *Tenant, duration time.Duration) {
	// 获取当前零点时间
	todayTime := getZeroTime(time.Now())
	yesterday := todayTime.AddDate(0, 0, -1)
	yesterdayTime := getZeroTime(yesterday)
	if tenant.MysqlClient != nil {
		// 失败任务
		failedJobs := selectFailedJob(todayTime.String(), tenant.MysqlClient)
		tenant.FailedJobs = failedJobs
		// 员工统计表
		staffCount := selectStaffCount(yesterdayTime.String(), tenant.MysqlClient)
		tenant.StaffCount = staffCount
		// 使用分析表
		useAnalyseCount := selectUseAnalyseCount(yesterdayTime.String(), tenant.MysqlClient)
		tenant.UseAnalyseCount = useAnalyseCount
		// 客户群统计表
		customerGroupCount := selectCustomerGroupCount(yesterdayTime.String(), tenant.MysqlClient)
		tenant.CustomerGroupCount = customerGroupCount
	}
	// 发送巡检报告
	markdown := dorisToMarkdown(tenant)
	log.Println("任务等待时间", duration)
	time.Sleep(duration)
	for _, robotkey := range tenant.DorisRobotkey {
		_ = notifier.SendWecom(markdown, robotkey, tenant.ProxyURL)
	}
}
