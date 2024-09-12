// Package cmd @Author lanpang
// @Date 2024/8/21 下午2:07:00
// @Desc
package cmd

import (
	"github.com/spf13/cobra"
	"vhagar/task/host"
	"vhagar/task/tenant"
)

var (
	_host   bool
	_tenant bool
	report  bool
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "检查服务",
	Long:  `支持各种服务的健康检测`,
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case _host:
			host.Check()
		case _tenant:
			tenant.Check(report)
		default:
			//
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().BoolVarP(&_host, "svc", "s", false, "检查主机的健康状态")
	checkCmd.Flags().BoolVarP(&_tenant, "tenant", "t", false, "检查企微租户的状态")
	checkCmd.Flags().BoolVarP(&report, "report", "r", false, "上报企微机器人")

}
