package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"vhagar/config"
	"vhagar/libs"

	"github.com/spf13/cobra"
	"github.com/tomasen/realip"
)

var (
	cfgFile  string
	Hostname string
	port     string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wsctl",
	Short: "微盛运维部署工具",
	Long:  `A longer description that vhagar`,
	Run: func(cmd *cobra.Command, args []string) {
		libs.Logger.Warnw("wsctl go go go！！！")
		startWeb(port)

	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		preFunc()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.toml", "config file")
	rootCmd.Flags().StringVarP(&port, "port", "p", "8099", "web 端口")
}

func preFunc() {
	if !filepath.IsAbs(cfgFile) {
		currentDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		cfgFile = filepath.Join(currentDir, cfgFile)
	}
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		panic("配置文件不存在: " + cfgFile)
	}
	if _, err := config.InitConfig(cfgFile); err != nil {
		panic("初始化配置失败: " + err.Error())
	}
	libs.InitLoggerWithConfig(config.Config.LogLevel, config.Config.LogToFile)
}

func startWeb(port string) {
	Hostname, _ = os.Hostname()
	http.HandleFunc("/", ping)
	libs.Logger.Warnw("启动 Web 服务", "url", fmt.Sprintf("http://%s:%s/", getClientIp(), port))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		libs.Logger.Fatalw("Web 服务启动失败", "err", err)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	responseCode := 200
	ip := getClientIp()
	RealIp := realip.FromRequest(r)
	responseJson := make(map[string]interface{})
	responseJson["ClientIp"] = ip
	responseJson["RequestURI"] = r.RequestURI
	responseJson["Header"] = r.Header
	responseJson["Method"] = r.Method
	responseJson["RealIp"] = RealIp
	djson := make(map[string]interface{})
	if err := r.ParseForm(); err != nil {
		djson["message"] = "Submit json format error"
	}
	responseJson["RequestJson"] = djson
	responseJson["Response_code"] = responseCode
	responseJson["Content-Type"] = r.Header.Get("Content-Type")
	responseJson["Hostname"] = Hostname
	// 输出 json 数据
	bytejson, _ := json.MarshalIndent(&responseJson, "", "  ")
	_, err := fmt.Fprintln(w, string(bytejson))
	if err != nil {
		libs.Logger.Errorw("响应输出失败", "err", err)
		return
	}
}

func getClientIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		libs.Logger.Errorw("获取本机 IP 地址失败", "err", err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//fmt.Println("本机 IP 地址:", ipnet.IP.String())
				return ipnet.IP.String()
			}
		}
	}

	return ""
}
