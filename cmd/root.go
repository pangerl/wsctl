package cmd

import (
	"log"
	"os"
	"path/filepath"
	"vhagar/config"

	"github.com/spf13/cobra"
)

var cfgFile string

// var cfg *config.CfgType

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wsctl",
	Short: "微盛运维部署工具",
	Long:  `A longer description that vhagar`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("wsctl go go go！！！")
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
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.toml", "config file")
}

func preFunc() {
	// 如果提供的是相对路径,将其转换为绝对路径
	if !filepath.IsAbs(cfgFile) {
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("无法获取当前工作目录: %v", err)
		}
		cfgFile = filepath.Join(currentDir, cfgFile)
	}

	// 确保文件存在
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		log.Fatalf("配置文件不存在: %s", cfgFile)
	}

	if _, err := config.InitConfig(cfgFile); err != nil {
		log.Fatalln("F! failed to init config:", err)
	}
}

//func preFunc() {
//	//homedir := "."
//	//configfile := filepath.Join(homedir, "config.toml")
//	log.Printf("读取配置文件 %s \n", cfgFile)
//	defer func() {
//		if err := recover(); err != nil {
//			log.Fatalf("Failed Info: 配置文件格式错误 %s", err)
//		}
//	}()
//	if _, err := os.Stat(cfgFile); err != nil {
//		if os.IsNotExist(err) {
//			log.Printf("读取配置文件 %s 失败，报错：%s", cfgFile, err)
//			createTempConfig()
//		}
//	} else {
//		if _, err := toml.DecodeFile(cfgFile, CONFIG); err != nil {
//
//			log.Fatalf("Failed Info: 配置文件格式错误 %s", err)
//		}
//		log.Println(CONFIG.VictoriaMetrics)
//		//log.Println(CONFIG.Cron["tenant"])
//		//log.Println(CONFIG.Tenant.Scheducron)
//		//fmt.Printf("租户信息: %+v\n", CONFIG.PG)
//		//fmt.Printf("租户信息: %+v\n", CONFIG.ES)
//	}
//}
