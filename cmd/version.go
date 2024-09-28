// Package cmd  @Author lanpang
// @Date 2024/8/1 上午11:25:00
// @Desc

package cmd

import (
	"fmt"
	"vhagar/libs"

	"github.com/spf13/cobra"
)

const VERSION = "v2.3"

var parrot bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "查看版本",
	Long:  `查看版本`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vhagar version: ", VERSION)
		if parrot {
			libs.RunParrot()
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&parrot, "parrot", "p", false, "一只疯狂的鹦鹉")
}
