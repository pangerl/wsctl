// Package cofing @Author lanpang
// @Date 2024/8/1 下午5:52:00
// @Desc
package cofing

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

type Config struct {
	ProjectName string
	Nacos       NacosConfig
}

func PreFunc() {
	fmt.Println("读取配置文件！！！")
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
		var newconfig Config
		if _, err := toml.DecodeFile("config.toml", &newconfig); err != nil {
			fmt.Println("配置文件格式错误", configfile)
			return
		}
		PROJECTNAME = newconfig.ProjectName
		NACOSCONFIG = newconfig.Nacos
		//fmt.Printf("全局信息: %+v\n\n", config.CONFIG.ProjectName)
		//fmt.Printf("全局信息: %+v\n\n", config.CONFIG)
		//fmt.Printf("全局信息: %+v\n\n", config.CONFIG.Nacos.Server)
	}
}
