// Package cmd  @Author lanpang
// @Date 2024/8/1 上午11:25:00
// @Desc

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const VERSION = "v3.2"

var parrot bool
var orientation string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "查看版本",
	Long:  `查看版本`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vhagar version: ", VERSION)
		if parrot {
			runParrot(orientation)
		}
	},
}

// init 函数在包被导入时自动执行
func init() {
	// 将 versionCmd 命令添加到 rootCmd 中
	rootCmd.AddCommand(versionCmd)
	// 为 versionCmd 命令设置一个布尔类型的标志 parrot
	versionCmd.Flags().BoolVarP(&parrot, "parrot", "p", false, "一只疯狂的鹦鹉")
	versionCmd.Flags().StringVarP(&orientation, "orientation", "o", "regular", "regular or aussie")
}
