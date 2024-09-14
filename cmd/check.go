// Package cmd @Author lanpang
// @Date 2024/8/21 下午2:07:00
// @Desc
package cmd

import (
	"github.com/spf13/cobra"
	"time"
	"vhagar/config"
	"vhagar/task/doris"
	"vhagar/task/host"
	"vhagar/task/nacos"
	"vhagar/task/rocketmq"
	"vhagar/task/tenant"
)

var (
	_host     bool
	_tenant   bool
	_nacos    bool
	_doris    bool
	_rocketmq bool
	report    bool
	watch     bool
	writefile string
	interval  time.Duration
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "检查服务",
	Long:  `支持各种服务的健康检测`,
	Run: func(cmd *cobra.Command, args []string) {
		// 监控模式
		config.Config.Global.Watch = watch
		config.Config.Global.Interval = interval
		config.Config.Global.Report = report
		config.Config.Nacos.Writefile = writefile
		switch {
		case _host:
			host.Check()
		case _tenant:
			tenant.Check()
		case _nacos:
			n := nacos.GetNacos()
			n.Check()
		case _doris:
			d := doris.GetDoris()
			d.Check()
		case _rocketmq:
			m := rocketmq.GetRocketMQ()
			m.Check()
		default:
			// 默认执行所有服务检查
			host.Check()
			tenant.Check()
			n := nacos.GetNacos()
			n.Check()
			d := doris.GetDoris()
			d.Check()
			m := rocketmq.GetRocketMQ()
			m.Check()
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().BoolVarP(&_host, "svc", "s", false, "检查主机的健康状态")
	checkCmd.Flags().BoolVarP(&_tenant, "tenant", "t", false, "检查企微租户的状态")
	checkCmd.Flags().BoolVarP(&_doris, "doris", "d", false, "检查doris的状态")
	checkCmd.Flags().BoolVarP(&_rocketmq, "rocketmq", "m", false, "检查rocketmq的状态")
	checkCmd.Flags().BoolVarP(&_nacos, "nacos", "n", false, "检查nacos的服务状态")
	checkCmd.Flags().BoolVarP(&watch, "watch", "w", false, "监控服务，定时刷新")
	checkCmd.Flags().DurationVarP(&interval, "second", "i", 5*time.Second, "自定义监控服务间隔刷新时间")
	checkCmd.Flags().BoolVarP(&report, "report", "r", false, "上报企微机器人")
	checkCmd.Flags().StringVarP(&writefile, "write", "o", "", "导出json文件, prometheus 自动发现文件路径")

}
