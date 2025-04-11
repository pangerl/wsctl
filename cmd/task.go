// Package cmd @Author lanpang
// @Date 2024/8/21 下午2:07:00
// @Desc
package cmd

import (
	"time"
	"vhagar/config"
	"vhagar/task"
	_ "vhagar/task/domain"
	_ "vhagar/task/doris"
	_ "vhagar/task/es"
	_ "vhagar/task/host"
	_ "vhagar/task/nacos"
	_ "vhagar/task/redis"
	_ "vhagar/task/rocketmq"
	_ "vhagar/task/tenant"

	"github.com/spf13/cobra"
)

var (
	_task     string
	report    bool
	watch     bool
	writefile string
	interval  time.Duration
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "检查服务",
	Long:  `支持各种服务的健康检测`,
	Run: func(cmd *cobra.Command, args []string) {
		if _task != "" {
			task.Do(_task)
		} else {
			for name := range task.Creators {
				task.Do(name)
			}
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		setEnv()
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.Flags().StringVarP(&_task, "task", "t", "", "指定要检查的服务 (host, tenant, nacos, doris, rocketmq, es, redis，domain)") // 更新帮助信息
	taskCmd.Flags().BoolVarP(&watch, "watch", "w", false, "nacos服务，定时刷新")
	taskCmd.Flags().DurationVarP(&interval, "second", "i", 5*time.Second, "自定义监控服务间隔刷新时间")
	taskCmd.Flags().BoolVarP(&report, "report", "r", false, "上报企微机器人")
	taskCmd.Flags().StringVarP(&writefile, "write", "o", "", "导出json文件, prometheus 自动发现文件路径")
}

func setEnv() {
	config.Config.Global.Watch = watch
	config.Config.Global.Interval = interval
	config.Config.Global.Report = report
	config.Config.Nacos.Writefile = writefile
}
