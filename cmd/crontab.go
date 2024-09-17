// Package cmd @Author lanpang
// @Date 2024/8/9 下午4:25:00
// @Desc
package cmd

import (
	"log"
	"vhagar/config"
	"vhagar/task/doris"
	"vhagar/task/tenant"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var crontabCmd = &cobra.Command{
	Use:   "cron",
	Short: " 启动定时任务",
	Long: `可自定义周期性运行 task
相关配置见配置文件的 [crontab]
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("启动任务调度")
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
	// 租户巡检 task
	if config.Config.Cron["tenant"].Crontab {
		// 加入定时任务
		_, err := c.AddFunc(config.Config.Cron["tenant"].Scheducron, func() {
			tenant.GetTenant().Check()
		})
		if err != nil {
			log.Fatalf("Failed to add crontab task: %s \n", err)
		}
	}
	//  doris 巡检 task
	if config.Config.Cron["doris"].Crontab {
		// 加入定时任务
		_, err := c.AddFunc(config.Config.Cron["doris"].Scheducron, func() {
			d := doris.GetDoris()
			d.Check()
		})
		if err != nil {
			log.Fatalf("Failed to add crontab task: %s \n", err)
		}
	}

	//启动/关闭
	c.Start()
	defer c.Stop()
	select {
	//查询语句，保持程序运行，在这里等同于for{}
	}
}

//func testjob() {
//	log.Printf("大王叫我来巡山，巡了南山巡北山。。。 \n")
//}
