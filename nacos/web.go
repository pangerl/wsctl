package nacos

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

//var Refreshtime time.Duration

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

func Webserver(nacos *Nacos) {
	fmt.Println("Start Nacos check web")
	gin.SetMode(gin.DebugMode)
	RefreshToken(nacos)
	r := gin.Default()
	v1 := r.Group("/")
	{
		v1.GET("/*route", func(c *gin.Context) {
			result, err := nacos.GetJson("json")
			if err != nil {
				c.JSON(500, []string{})
				return
			}
			c.JSON(200, result)
		})
	}
	err := r.Run(nacos.Webport)
	if err != nil {
		fmt.Println(err)
	}
}

func RefreshToken(nacos *Nacos) {
	if len(nacos.Config.Username) != 0 && len(nacos.Config.Password) != 0 {
		go func() {
			for {
				nacos.WithAuth()
				time.Sleep(time.Second * 3600)
			}
		}()
	}
}

func RefreshNacosInstance(nacos *Nacos, interval time.Duration) {
	go func() {
		for {
			nacos.GetNacosInstance()
			time.Sleep(interval)
		}
	}()
}
