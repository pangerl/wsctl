// Package config @Author lanpang
// @Date 2024/8/1 下午2:47:00
// @Desc
package config

var (
	PROJECTNAME string
	NACOS       Nacos
)

type Nacos struct {
	Server   string
	Username string
	Password string
}
