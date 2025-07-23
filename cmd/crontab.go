// Package cmd @Author lanpang
// @Date 2024/8/9 下午4:25:00
// @Desc
package cmd

import (
	"vhagar/config"
	"vhagar/logger"
	"vhagar/metric"
	"vhagar/task"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var crontabCmd = &cobra.Command{
	Use:   "cron",
	Short: "启动定时任务",
	Long: `可自定义周期性运行 task
相关配置见配置文件的 [crontab]
`,
	Run: func(cmd *cobra.Command, args []string) {
		// 启动 metric 服务
		if config.Config.Metric.Enable {
			go metric.StartMetric()
		}
		// 启动定时任务
		logger.Logger.Warnw("启动任务调度")
		crontabJob()
	},
}

func init() {
	rootCmd.AddCommand(crontabCmd)
}

func crontabJob() {
	c := cron.New() //创建一个cron实例
	// 获取等待时间
	duration := config.GetRandomDuration()
	config.Config.Global.Duration = duration
	config.Config.Global.Report = true
	cronCfg := config.Config.Cron
	// 添加任务
	for name, cronJob := range cronCfg {
		// 判断是否是定时任务
		taskName := name
		if cronJob.Crontab {
			logger.Logger.Warnw("添加定时任务", "task", taskName)
			_, err := c.AddFunc(cronJob.Scheducron, func() {
				task.Do(taskName)
			})
			if err != nil {
				logger.Logger.Fatalw("添加定时任务失败", "err", err)
			}
		}
	}
	//启动/关闭
	c.Run()
	defer c.Stop()
}
