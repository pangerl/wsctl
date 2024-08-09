// Package cmd @Author lanpang
// @Date 2024/8/6 下午4:49:00
// @Desc
package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"time"
	"vhagar/inspect"
)

// versionCmd represents the version command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "项目巡检",
	Long:  `获取项目的企业数据，活跃数，会话数`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("开始项目巡检")
		esclient, _ := inspect.NewESClient(CONFIG.ES)
		pgclient1, pgclient2 := inspect.NewPGClient(CONFIG.PG)
		defer func() {
			esclient.Stop()
			err := pgclient1.Close(context.Background())
			if err != nil {
				return
			}
			err = pgclient2.Close(context.Background())
			if err != nil {
				return
			}
		}()

		// 创建 inspect 对象
		_inspect := inspect.NewInspect(CONFIG.Tenant.Corp, esclient, pgclient1, pgclient2, CONFIG.ProjectName, VERSION)
		inspectTask(_inspect)
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}

func inspectTask(_inspect *inspect.Inspect) {
	// 当前时间
	dateNow := time.Now().AddDate(0, 0, 0)
	log.Print("启动巡检任务")
	//inspect.GetVersion(url)
	for _, corp := range _inspect.Corp {
		//fmt.Println(corp.Corpid)
		// 获取租户名
		_inspect.SetCorpName(corp.Corpid)
		// 获取用户数
		_inspect.SetUserNum(corp.Corpid)
		// 获取客户数
		_inspect.SetCustomerNum(corp.Corpid)
		// 获取活跃数
		_inspect.SetActiveNum(corp.Corpid, dateNow)
		// 获取会话数
		if corp.Convenabled {
			_inspect.SetMessageNum(corp.Corpid, dateNow)
		}
		//fmt.Println(*corp)
	}
	// 发送巡检报告
	markdown, _ := _inspect.TransformToMarkdown(CONFIG.Inspection.Userlist, dateNow)
	err := inspect.SendWecom(markdown, CONFIG.Inspection.Robotkey, CONFIG.ProxyURL)
	if err != nil {
		return
	}
}
