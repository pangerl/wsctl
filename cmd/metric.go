// Package cmd @Author lanpang
// @Date 2024/8/16 下午1:42:00
// @Desc
package cmd

import (
	"github.com/spf13/cobra"
	"vhagar/metric"
)

var metricCmd = &cobra.Command{
	Use:   "metric",
	Short: "监控指标",
	Long:  `监控指标metric`,
	Run: func(cmd *cobra.Command, args []string) {
		// 启动 metric 服务
		metric.StartMetric()
	},
}

func init() {
	rootCmd.AddCommand(metricCmd)
}
