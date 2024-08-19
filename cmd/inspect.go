// Package cmd @Author lanpang
// @Date 2024/8/6 下午4:49:00
// @Desc
package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"time"
	"vhagar/inspect"
	"vhagar/notifier"
)

var (
	rocketmq bool
)

// versionCmd represents the version command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "项目巡检",
	Long:  `获取项目的企业数据，活跃数，会话数`,
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case rocketmq:
			mqTask()
		default:
			log.Println("开始项目巡检")
			// 初始化 inspect 对象
			esclient, _ := inspect.NewESClient(CONFIG.ES)
			dbClient, _ := inspect.NewPGClient(CONFIG.PG)
			defer func() {
				if dbClient != nil {
					dbClient.Close()
				}
				if esclient != nil {
					esclient.Stop()
				}
			}()
			_inspect := inspect.NewInspect(CONFIG.Tenant.Corp, esclient, dbClient, CONFIG.ProjectName)
			// 执行巡检 job
			inspectTask(_inspect, 0)
		}
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().BoolVarP(&rocketmq, "rocketmq", "m", false, "获取 rocketmq broker 信息")

}

func inspectTask(_inspect *inspect.Inspect, duration time.Duration) {
	// 当前时间
	dateNow := time.Now().AddDate(0, 0, 0)
	log.Print("启动企微租户巡检任务")
	//inspect.GetVersion(url)
	for _, corp := range _inspect.Corp {
		//fmt.Println(corp.Corpid)
		if _inspect.DBClient != nil {
			// 获取租户名
			_inspect.SetCorpName(corp.Corpid)
			// 获取用户数
			_inspect.SetUserNum(corp.Corpid)
			// 获取客户群
			_inspect.SetCustomerGroupNum(corp.Corpid)
			// 获取客户群人数
			_inspect.SetCustomerGroupUserNum(corp.Corpid)
		}
		if _inspect.EsClient != nil {
			// 获取客户数
			_inspect.SetCustomerNum(corp.Corpid)
			// 获取活跃数
			_inspect.SetActiveNum(corp.Corpid, dateNow)
			// 获取会话数
			if corp.Convenabled {
				_inspect.SetMessageNum(corp.Corpid, dateNow)
			}
		}
		//fmt.Println(*corp)
	}
	// 发送巡检报告
	markdownList := _inspect.TransformToMarkdown(CONFIG.Inspection.Userlist, dateNow)
	log.Println("任务等待时间", duration)
	time.Sleep(duration)
	for _, markdown := range markdownList {
		for _, robotkey := range CONFIG.Inspection.Robotkey {
			err := notifier.SendWecom(markdown, robotkey, CONFIG.ProxyURL)
			if err != nil {
				return
			}
		}
	}

}

func mqTask() {
	log.Print("启动 rocketmq 巡检任务")
	clusterdata, _ := inspect.GetMQDetail(CONFIG.Rocketmq.RocketmqDashboard)
	markdown := inspect.MQDetailToMarkdown(clusterdata, CONFIG.ProjectName)
	for _, robotkey := range CONFIG.Inspection.Robotkey {
		_ = notifier.SendWecom(markdown, robotkey, CONFIG.ProxyURL)
	}
}
