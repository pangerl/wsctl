package cmd

import (
	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wsctl",
	Short: "ws运维部署工具",
	Long:  `A longer description that vhagar`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("程序开始启动！！！")
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		preFunc()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func preFunc() {
	homedir := "."
	configfile := filepath.Join(homedir, "config.toml")
	log.Printf("读取配置文件 %s \n", configfile)
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("Failed Info: 配置文件格式错误 %s", err)
		}
	}()
	if _, err := os.Stat(configfile); err != nil {
		if !os.IsExist(err) {
			log.Fatalf("Failed Info: 读取配置文件报错 %s", err)
		}
	} else {
		if _, err := toml.DecodeFile("config.toml", &CONFIG); err != nil {
			log.Fatalf("Failed Info: 配置文件格式错误 %s", err)
		}
		//fmt.Println(CONFIG.Crontab)
		//fmt.Printf("租户信息: %+v\n", CONFIG.PG)
		//fmt.Printf("租户信息: %+v\n", CONFIG.ES)
	}
}
