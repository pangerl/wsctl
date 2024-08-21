// Package cmd  @Author lanpang
// @Date 2024/8/1 上午11:25:00
// @Desc

package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"time"
	"vhagar/nacos"
)

var (
	writefile string
	watch     bool
	interval  time.Duration
)

// nacosCmd represents the version command
var nacosCmd = &cobra.Command{
	Use:   "nacos",
	Short: "服务健康检查工具",
	Long:  `通过 nacos 的服务注册信息，统计微服务的信息`,
	Run: func(cmd *cobra.Command, args []string) {
		_nacos := nacos.NewNacos(CONFIG.Nacos, writefile)
		log.Println("获取nacos认证信息")
		_nacos.WithAuth()
		_nacos.GetNacosInstance()
		switch {
		//case web:
		//	nacos.Webserver(_nacos)
		case writefile != "":
			_nacos.WriteFile()
		default:
			if watch {
				log.Printf("监控模式 刷新时间:%s/次\n", interval)
				for {
					_nacos.GetNacosInstance()
					_nacos.TableRender()
					time.Sleep(interval)
				}
			}
			_nacos.TableRender()
		}
	},
}

func init() {
	rootCmd.AddCommand(nacosCmd)
	nacosCmd.Flags().StringVarP(&writefile, "write", "o", "", "导出json文件, prometheus 自动发现文件路径")
	nacosCmd.Flags().BoolVarP(&watch, "watch", "w", false, "监控服务，定时刷新")
	nacosCmd.Flags().DurationVarP(&interval, "second", "s", 5*time.Second, " 自定义监控服务间隔刷新时间")
	//nacosCmd.Flags().BoolVarP(&web, "web", "w", false, "开启web api Prometheus http_sd_configs")
	//nacosCmd.Flags().StringVarP(&webport, "port", "p", "8099", "web 端口")
}
