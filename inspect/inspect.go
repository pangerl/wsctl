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

func TenantTask(inspect *Inspect, duration time.Duration) {
	tenant := inspect.Tenant
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
	markdownList := tenantNotifier(tenant, inspect.ProjectName, inspect.Notifier["tenant"].Userlist)
	log.Println("任务等待时间", duration)
	time.Sleep(duration)
	for _, markdown := range markdownList {
		for _, robotkey := range inspect.Notifier["tenant"].Robotkey {
			err := notifier.SendWecom(markdown, robotkey, inspect.ProxyURL)
			if err != nil {
				return
			}
		}
	}

}

func RocketmqTask(inspect *Inspect) {
	log.Print("启动 rocketmq 巡检任务")
	clusterdata, _ := GetMQDetail(inspect.Rocketmq.RocketmqDashboard)
	markdown := mqDetailMarkdown(clusterdata, inspect.ProjectName)
	for _, robotkey := range inspect.Notifier["rocketmq"].Robotkey {
		_ = notifier.SendWecom(markdown, robotkey, inspect.ProxyURL)
	}
}

func DorisTask(inspect *Inspect, duration time.Duration) {
	// 获取当前零点时间
	todayTime := getZeroTime(time.Now())
	yesterday := todayTime.AddDate(0, 0, -1)
	yesterdayTime := getZeroTime(yesterday)
	if inspect.Doris.MysqlClient != nil {
		// 失败任务
		failedJobs := selectFailedJob(todayTime.String(), inspect.Doris.MysqlClient)
		inspect.Doris.FailedJobs = failedJobs
		// 员工统计表
		staffCount := selectStaffCount(yesterdayTime.String(), inspect.Doris.MysqlClient)
		inspect.Doris.StaffCount = staffCount
		// 使用分析表
		useAnalyseCount := selectUseAnalyseCount(yesterdayTime.String(), inspect.Doris.MysqlClient)
		inspect.Doris.UseAnalyseCount = useAnalyseCount
		// 客户群统计表
		customerGroupCount := selectCustomerGroupCount(yesterdayTime.String(), inspect.Doris.MysqlClient)
		inspect.Doris.CustomerGroupCount = customerGroupCount
	}
	// 发送巡检报告
	markdown := dorisToMarkdown(inspect.Doris, inspect.ProjectName)
	log.Println("任务等待时间", duration)
	time.Sleep(duration)
	for _, robotkey := range inspect.Notifier["doris"].Robotkey {
		_ = notifier.SendWecom(markdown, robotkey, inspect.ProxyURL)
	}
}
