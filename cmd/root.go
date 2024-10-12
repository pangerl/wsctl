package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomasen/realip"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"vhagar/config"
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
		log.Println("wsctl go go go！！！")
		startWeb(port)

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
	rootCmd.Flags().StringVarP(&port, "port", "p", "8099", "web 端口")
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

func startWeb(port string) {
	Hostname, _ = os.Hostname()
	http.HandleFunc("/", ping)
	log.Printf("Starting server at http://%s:%s/\n", getClientIp(), port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
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
		log.Printf("Failed to start server: %v", err)
		return
	}
}

func getClientIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("获取本机 IP 地址失败:", err)
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
