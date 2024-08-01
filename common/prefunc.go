// Package common @Author lanpang
// @Date 2024/8/1 下午5:52:00
// @Desc
package common

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"vhagar/config"
)

type NewConfig struct {
	ProjectName string
	Nacos       config.Nacos
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
		var newConfig NewConfig
		if _, err := toml.DecodeFile("config.toml", &newConfig); err != nil {
			fmt.Println("配置文件格式错误", configfile)
			return
		}
		config.PROJECTNAME = newConfig.ProjectName
		fmt.Printf("全局信息: %+v\n\n", newConfig.ProjectName)
		fmt.Printf("全局信息: %+v\n\n", config.PROJECTNAME)
		fmt.Printf("全局信息: %+v\n\n", newConfig.Nacos.Server)
	}
}
