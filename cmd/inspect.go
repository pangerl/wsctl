// Package cmd @Author lanpang
// @Date 2024/8/6 下午4:49:00
// @Desc
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
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
		inspect := inspect.NewInspect(CONFIG.Tenant.Corp, esclient)
		fmt.Println(CONFIG.Tenant)
		//inspect.GetVersion(url)
		for _, corp := range inspect.Corp {
			fmt.Println(corp.Corpid)
			inspect.GetCustomerNum(corp.Corpid)
		}
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
