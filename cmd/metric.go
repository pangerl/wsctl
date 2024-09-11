// Package cmd @Author lanpang
// @Date 2024/8/16 下午1:42:00
// @Desc
package cmd

//var metricCmd = &cobra.Command{
//	Use:   "metric",
//	Short: "监控指标",
//	Long:  `监控指标metric`,
//	Run: func(cmd *cobra.Command, args []string) {
//		// 初始化 metric 对象
//		metrics := NewMetric(CONFIG)
//		// 初始化 esclient
//		esclient, _ := libs.NewESClient(CONFIG.ES)
//		defer func() {
//			if esclient != nil {
//				esclient.Stop()
//			}
//		}()
//		metrics.EsClient = esclient
//		// 启动 metric 服务
//		metric.StartMetric(metrics)
//	},
//}
//
//func init() {
//	rootCmd.AddCommand(metricCmd)
//}
//
//func NewMetric(cfg *Config) *metric.Metric {
//	return &metric.Metric{
//		Corp:     cfg.Tenant.Corp,
//		Rocketmq: cfg.Rocketmq,
//		Metric:   cfg.Metric,
//		Nacos:    cfg.Nacos,
//	}
//}
