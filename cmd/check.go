// Package cmd @Author lanpang
// @Date 2024/8/21 下午2:07:00
// @Desc
package cmd

import (
	"time"
	"vhagar/config"
	"vhagar/task/doris"
	//"vhagar/task/es" // 新增 ES 任务包导入
	"vhagar/task/host"
	"vhagar/task/nacos"
	"vhagar/task/rocketmq"
	"vhagar/task/tenant"

	"github.com/spf13/cobra"
)

var (
	task      string
	report    bool
	watch     bool
	writefile string
	interval  time.Duration
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "检查服务",
	Long:  `支持各种服务的健康检测`,
	Run: func(cmd *cobra.Command, args []string) {
		// 监控模式
		config.Config.Global.Watch = watch
		config.Config.Global.Interval = interval
		config.Config.Global.Report = report
		config.Config.Nacos.Writefile = writefile

		var tasks []config.Tasker

		switch task {
		case "host":
			tasks = append(tasks, host.GetHost())
		case "tenant":
			tasks = append(tasks, tenant.GetTenant())
		case "nacos":
			tasks = append(tasks, nacos.GetNacos())
		case "doris":
			tasks = append(tasks, doris.Work())
		case "rocketmq":
			tasks = append(tasks, rocketmq.GetRocketMQ())
		//case "es": // 新增 ES 检查选项
		//	tasks = append(tasks, es.GetES())
		default:
			// 默认执行所有服务检查
			tasks = []config.Tasker{
				host.GetHost(),
				tenant.GetTenant(),
				nacos.GetNacos(),
				doris.Work(),
				rocketmq.GetRocketMQ(),
				//es.GetES(), // 在默认任务列表中添加 ES 检查
			}
		}

		for _, t := range tasks {
			t.Check()
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVarP(&task, "task", "t", "", "指定要检查的服务 (host, tenant, nacos, doris, rocketmq, es)") // 更新帮助信息
	checkCmd.Flags().BoolVarP(&watch, "watch", "w", false, "监控服务，定时刷新")
	checkCmd.Flags().DurationVarP(&interval, "second", "i", 5*time.Second, "自定义监控服务间隔刷新时间")
	checkCmd.Flags().BoolVarP(&report, "report", "r", false, "上报企微机器人")
	checkCmd.Flags().StringVarP(&writefile, "write", "o", "", "导出json文件, prometheus 自动发现文件路径")
}
