// Package cmd  @Author lanpang
// @Date 2024/8/1 上午11:25:00
// @Desc

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
	"vhagar/nacos"
)

var (
	web       bool
	webport   string
	writefile string
	watch     bool
)

// versionCmd represents the version command
var nacosCmd = &cobra.Command{
	Use:   "nacos",
	Short: "nacos",
	Long:  `nacos`,
	Run: func(cmd *cobra.Command, args []string) {
		_nacos := nacos.NewNacos(CONFIG.Nacos, web, webport, writefile)
		fmt.Println("获取nacos认证信息")
		_nacos.WithAuth()
		_nacos.GetNacosInstance()
		switch {
		case web:
			nacos.Webserver(_nacos)
		case writefile != "":
			_nacos.WriteFile()
		default:
			if watch {
				fmt.Printf("监控模式 刷新时间:%s/次\n", 5*time.Second)
				for {
					_nacos.GetNacosInstance()
					_nacos.TableRender()
					time.Sleep(5 * time.Second)
				}
			}
			_nacos.TableRender()
		}
	},
}

func init() {
	rootCmd.AddCommand(nacosCmd)
	nacosCmd.Flags().StringVarP(&writefile, "write", "o", "", "导出json文件, prometheus 自动发现文件路径")
	nacosCmd.Flags().BoolVarP(&web, "web", "w", false, "监控服务")
	nacosCmd.Flags().StringVarP(&webport, "port", "p", ":8099", "web 端口")
	nacosCmd.Flags().BoolVarP(&watch, "watch", "d", false, "监控服务")
}
