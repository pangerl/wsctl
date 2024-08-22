// Package cmd @Author lanpang
// @Date 2024/8/9 下午4:25:00
// @Desc
package cmd

import (
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"log"
	"vhagar/inspect"
)

var crontabCmd = &cobra.Command{
	Use:   "crontab",
	Short: " 启动定时任务",
	Long: `可自定义周期性运行 job
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
	duration := inspect.GetRandomDuration()
	// 每日巡检 job
	if CONFIG.Crontab.Inspectjob {
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
		// 加入定时任务
		_, err := c.AddFunc(CONFIG.Inspection.Scheducron, func() {
			inspectTask(_inspect, duration)
		})
		if err != nil {
			log.Printf("Failed to add crontab job: %s \n", err)
		}
	}
	// 测试任务
	if CONFIG.Crontab.Testjob {
		_, err := c.AddFunc("* * * * *", func() {
			testjob()
		})
		if err != nil {
			log.Printf("Failed to add crontab job: %s \n", err)
		}
	}

	//启动/关闭
	c.Start()
	defer c.Stop()
	select {
	//查询语句，保持程序运行，在这里等同于for{}
	}
}

func testjob() {
	log.Printf("大王叫我来巡山，巡了南山巡北山。。。 \n")
}
