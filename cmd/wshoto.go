// Package cmd  @Author lanpang
// @Date 2024/8/1 上午11:25:00
// @Desc

package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"time"
	"vhagar/metric"
	"vhagar/nacos"
)

var (
	web       bool
	webport   string
	writefile string
	watch     bool
	monitor   bool
)

// versionCmd represents the version command
var wshotoCmd = &cobra.Command{
	Use:   "wshoto",
	Short: "服务健康检查工具",
	Long:  `通过 nacos 的服务注册信息，统计微服务的信息`,
	Run: func(cmd *cobra.Command, args []string) {
		_nacos := nacos.NewNacos(CONFIG.Nacos, web, webport, writefile)
		log.Println("获取nacos认证信息")
		_nacos.WithAuth()
		_nacos.GetNacosInstance()
		switch {
		case web:
			nacos.Webserver(_nacos)
		case writefile != "":
			_nacos.WriteFile()
		case monitor:
			metric.Monitor(_nacos, CONFIG.Rocketmq.RocketmqDashboard)
		default:
			if watch {
				log.Printf("监控模式 刷新时间:%s/次\n", 10*time.Second)
				for {
					_nacos.GetNacosInstance()
					_nacos.TableRender()
					time.Sleep(10 * time.Second)
				}
			}
			_nacos.TableRender()
		}
	},
}

func init() {
	rootCmd.AddCommand(wshotoCmd)
	wshotoCmd.Flags().StringVarP(&writefile, "write", "o", "", "导出json文件, prometheus 自动发现文件路径")
	wshotoCmd.Flags().BoolVarP(&web, "web", "w", false, "开启web api Prometheus http_sd_configs")
	wshotoCmd.Flags().StringVarP(&webport, "port", "p", "8099", "web 端口")
	wshotoCmd.Flags().BoolVarP(&watch, "watch", "d", false, "监控服务，定时刷新")
	wshotoCmd.Flags().BoolVarP(&monitor, "monitor", "m", false, "服务接口探活")

}
