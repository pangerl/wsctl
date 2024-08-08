package cmd

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
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
		PreFunc()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func PreFunc() {
	fmt.Println("读取配置文件！")
	homedir := "."
	configfile := filepath.Join(homedir, "config.toml")
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("配置文件格式错误", configfile, err)
			os.Exit(2)
		}
	}()
	if _, err := os.Stat(configfile); err != nil {
		if !os.IsExist(err) {
			fmt.Println("读取配置文件报错", configfile, err)
			return
		}
	} else {
		if _, err := toml.DecodeFile("config.toml", &CONFIG); err != nil {
			fmt.Println("配置文件格式错误", configfile)
			return
		}
		//fmt.Printf("租户信息: %+v\n", CONFIG.PG)
		//fmt.Printf("租户信息: %+v\n", CONFIG.ES)

	}
}
