// Package cmd @Author lanpang
// @Date 2024/8/21 下午2:07:00
// @Desc
package cmd

import (
	"github.com/spf13/cobra"
	"vhagar/check"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "检查服务",
	Long:  `支持各种服务的健康检测`,
	Run: func(cmd *cobra.Command, args []string) {
		//mqStatus := check.ProbeRocketMq(CONFIG.Rocketmq.NameServer)
		check.HostCheck(CONFIG.VictoriaMetrics)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
