// Package cmd @Author lanpang
// @Date 2024/8/6 下午4:49:00
// @Desc
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "项目巡检",
	Long:  `获取项目的企业数据，活跃数，会话数`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vhagar version: ", version)
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
