// Package cofing @Author lanpang
// @Date 2024/8/1 下午2:47:00
// @Desc
package cofing

var (
	CONFIG      Config
	PROJECTNAME string
	NACOSCONFIG NacosConfig
	WATCH       bool
	WRITEFILE   string
)

type NacosConfig struct {
	Server    string
	Username  string
	Password  string
	Namespace string
}
