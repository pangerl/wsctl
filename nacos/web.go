package nacos

//
//import (
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"time"
//	"vhagar/cmd"
//)
//
////var Refreshtime time.Duration
//
//func response(c *gin.Context) {
//	if c.Request.RequestURI == "/health" {
//		c.JSON(200, gin.H{"status": true})
//		return
//	}
//	if c.Request.RequestURI == "/favicon.ico" {
//		c.JSON(404, "404")
//		return
//	}
//	result, err := config.NACOS.GetJson("json", true)
//	if err != nil {
//		c.JSON(500, []string{})
//		return
//	}
//	c.JSON(200, result)
//}
//
//func Webserver() {
//	fmt.Println("Start Nacos check web")
//	gin.SetMode(gin.DebugMode)
//	RefreshToken()
//	r := gin.Default()
//	v1 := r.Group("/")
//	{
//		v1.GET("/*route", response)
//	}
//	err := r.Run(cmd.WEBPORT)
//	if err != nil {
//		fmt.Println(err)
//	}
//}
//
//func RefreshToken() {
//	if len(cmd.NACOSCONFIG.Username) != 0 && len(cmd.NACOSCONFIG.Password) != 0 {
//		go func() {
//			for {
//				config.NACOS.WithAuth()
//				time.Sleep(time.Second * 3600)
//			}
//		}()
//	}
//}
