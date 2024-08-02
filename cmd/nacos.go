// Package cmd  @Author lanpang
// @Date 2024/8/1 上午11:25:00
// @Desc

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"vhagar/cofing"
	"vhagar/nacos"
)

var Nacos nacos.Nacos
var conf = &cofing.NACOSCONFIG

// versionCmd represents the version command
var nacosCmd = &cobra.Command{
	Use:   "nacos",
	Short: "nacos",
	Long:  `nacos`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(conf.Username) != 0 && len(conf.Password) != 0 {
			Nacos.WithAuth()
		}
		fmt.Println("获取服务信息")
		Nacos.GetNacosInstance()
		switch {
		case cofing.WATCH:
			//
		case cofing.WRITEFILE != "":
			Nacos.WriteFile()
		default:
			Nacos.TableRender()
			//fmt.Println("x", Nacos)
		}
	},
}

func init() {
	rootCmd.AddCommand(nacosCmd)
	rootCmd.Flags().StringVarP(&cofing.WRITEFILE, "write", "o", "", "导出json文件, prometheus 自动发现文件路径")
	rootCmd.Flags().BoolVarP(&cofing.WATCH, "watch", "w", false, "监控服务")
}
