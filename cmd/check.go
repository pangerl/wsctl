// Package cmd @Author lanpang
// @Date 2024/8/21 下午2:07:00
// @Desc
package cmd

import (
	"github.com/spf13/cobra"
	"vhagar/task/host"
)

var (
	_host   bool
	_tenant bool
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "检查服务",
	Long:  `支持各种服务的健康检测`,
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case _host:
			host.CheckHost()
		case _tenant:
		//
		default:
			//
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().BoolVarP(&_host, "主机巡检", "s", false, "检查主机的健康状态")
	checkCmd.Flags().BoolVarP(&_tenant, "企微租户巡检", "t", false, "检查企微租户的状态")

}
