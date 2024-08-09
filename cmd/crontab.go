// Package cmd @Author lanpang
// @Date 2024/8/9 下午4:25:00
// @Desc
package cmd

import (
	"context"
	"github.com/robfig/cron/v3"
	"log"
	"vhagar/inspect"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var crontabCmd = &cobra.Command{
	Use:   "crontab",
	Short: " 任务调度",
	Long:  `定时调度任务`,
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
	// 每日巡检 job
	if CONFIG.Crontab.Inspectjob {
		// 初始化 inspect 对象
		esclient, _ := inspect.NewESClient(CONFIG.ES)
		pgclient1, pgclient2 := inspect.NewPGClient(CONFIG.PG)
		defer func() {
			if pgclient1 != nil {
				err := pgclient1.Close(context.Background())
				if err != nil {
					return
				}
			}
			if pgclient1 != nil {
				err := pgclient2.Close(context.Background())
				if err != nil {
					return
				}
			}
			if esclient != nil {
				esclient.Stop()
			}
		}()
		_inspect := inspect.NewInspect(CONFIG.Tenant.Corp, esclient, pgclient1, pgclient2, CONFIG.ProjectName, VERSION)
		// 加入定时任务
		_, err := c.AddFunc(CONFIG.Inspection.Scheducron, func() {
			inspectTask(_inspect)
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
