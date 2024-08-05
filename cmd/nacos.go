// Package cmd  @Author lanpang
// @Date 2024/8/1 上午11:25:00
// @Desc

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"vhagar/nacos"
)

var (
	web       bool
	webport   string
	writefile string
)

// versionCmd represents the version command
var nacosCmd = &cobra.Command{
	Use:   "nacos",
	Short: "nacos",
	Long:  `nacos`,
	Run: func(cmd *cobra.Command, args []string) {
		nacos := nacos.NewNacos(CONFIG.Nacos, web, webport, writefile)
		fmt.Println("获取nacos认证信息")
		nacos.WithAuth()
		fmt.Println("获取注册服务信息")
		nacos.GetNacosInstance()
		switch {
		case web:
			//nacos.Webserver()
		case writefile != "":
			nacos.WriteFile()
		default:
			nacos.TableRender()
			//fmt.Println("x", Nacos)
		}
	},
}

func init() {
	rootCmd.AddCommand(nacosCmd)
	nacosCmd.Flags().StringVarP(&writefile, "write", "o", "", "导出json文件, prometheus 自动发现文件路径")
	nacosCmd.Flags().BoolVarP(&web, "web", "w", false, "监控服务")
	nacosCmd.Flags().StringVarP(&webport, "port", "p", ":8099", "web 端口")

}
