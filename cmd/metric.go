// Package cmd @Author lanpang
// @Date 2024/8/16 下午1:42:00
// @Desc
package cmd

import (
	"github.com/spf13/cobra"
	"vhagar/inspect"
	"vhagar/metric"
)

var metricCmd = &cobra.Command{
	Use:   "metric",
	Short: "监控指标",
	Long:  `监控指标metric`,
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化 metric 对象
		esclient, _ := inspect.NewESClient(CONFIG.ES)
		defer func() {
			if esclient != nil {
				esclient.Stop()
			}
		}()
		m := metric.NewMetric(CONFIG.Metric, CONFIG.Nacos, CONFIG.Rocketmq, CONFIG.Tenant.Corp, esclient)
		m.StartMetric()
	},
}

func init() {
	rootCmd.AddCommand(metricCmd)
}
