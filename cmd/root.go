package cmd

import (
	"embed"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/tomasen/realip"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"vhagar/config"
)

var (
	cfgFile  string
	Hostname string
	//go:embed templates/*.tmpl
	tmpl embed.FS
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wsctl",
	Short: "微盛运维部署工具",
	Long:  `A longer description that vhagar`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("wsctl go go go！！！")
		log.Print("启动调试 web 服务")
		startWeb()

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

func startWeb() {
	Hostname, _ = os.Hostname()
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	t, _ := template.ParseFS(tmpl, "templates/*.tmpl")
	r.SetHTMLTemplate(t)
	v1 := r.Group("/ping")
	v1.Any("/*router", response)
	err := r.Run(":" + config.Config.Port)
	if err != nil {
		log.Printf("Failed to start server: %v", err)
		return
	}
}

func response(c *gin.Context) {
	responseCode := 200
	format := c.DefaultQuery("format", "json")
	httpCode := c.DefaultQuery("http_code", "200")
	if value, err := strconv.Atoi(httpCode); err == nil {
		responseCode = value
	}
	ip := c.ClientIP()
	djson := make(map[string]interface{})
	contentType := c.GetHeader("Content-Type")
	if err := c.ShouldBindJSON(&djson); err != nil && contentType == "application/json" {
		djson["message"] = "Submit json format error"
	}
	RealIp := realip.FromRequest(c.Request)
	responseJson := make(map[string]interface{})
	responseJson["ClientIp"] = ip
	responseJson["RequestURI"] = c.Request.RequestURI
	responseJson["Header"] = c.Request.Header
	responseJson["Method"] = c.Request.Method
	responseJson["RealIp"] = RealIp
	responseJson["RequestJson"] = djson
	responseJson["RequestPostForm"] = c.Request.PostForm
	responseJson["Response_code"] = responseCode
	responseJson["Content-Type"] = c.Request.Header.Get("Content-Type")
	responseJson["Hostname"] = Hostname
	bytejson, _ := json.MarshalIndent(&djson, "", "  ")
	log.Printf("\n============================================================================\n"+
		"Header:%s\n"+
		"IP:%s\n"+
		"X-Forwarded-For:%s\n"+
		"X-Real-Ip:%s\n"+
		"X-Forwarded-Host:%s\n"+
		"RemoteAddr:%s\n"+
		"Content-Type:%s\n"+
		"RequestJson::%s\n"+
		"RequestPostForm::%s\n",
		c.Request.Header,
		c.ClientIP(),
		c.Request.Header.Get("X-Forwarded-For"),
		c.Request.Header.Get("X-Real-Ip"),
		c.Request.Header.Get("X-Forwarded-Host:"),
		c.Request.RemoteAddr,
		c.Request.Header.Get("Content-Type"),
		string(bytejson),
		c.Request.PostForm)
	if format == "json" {
		c.JSON(responseCode, responseJson)
	} else {
		c.HTML(responseCode, "index.tmpl", gin.H{
			"response_json": responseJson,
			"Header":        c.Request.Header,
		})
	}
}
