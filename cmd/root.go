package cmd

import (
	"fmt"
	"os"
	"vhagar/common"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vhagar",
	Short: "A brief description of vhagar",
	Long:  `A longer description that vhagar`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("程序开始启动！！！")
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		common.PreFunc()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
