// Package cmd @Author lanpang
// @Date 2024/8/6 下午4:49:00
// @Desc
package cmd

import (
	"github.com/spf13/cobra"
	"vhagar/inspect"
)

var (
	rocketmq bool
)

// inspectCmd represents the version command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "项目巡检",
	Long:  `获取项目的企业数据，活跃数，会话数`,
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		//case rocketmq:
		//	mqTask()
		default:
			// 执行巡检 job
			tenant := inspect.NewTenant(
				CONFIG.ES, CONFIG.PG, CONFIG.ProjectName, CONFIG.ProxyURL,
				CONFIG.Tenant.Corp, CONFIG.Tenant.Userlist, CONFIG.Tenant.Robotkey)
			inspect.TenantTask(tenant, 0)
		}
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().BoolVarP(&rocketmq, "rocketmq", "m", false, "获取 rocketmq broker 信息")

}
