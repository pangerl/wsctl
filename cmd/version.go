// Package cmd  @Author lanpang
// @Date 2024/8/1 上午11:25:00
// @Desc

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "v1.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version",
	Long:  `查看版本`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vhagar version: ", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
