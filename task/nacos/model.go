// Package nacos @Author lanpang
// @Date 2024/8/2 下午4:41:00
// @Desc
package nacos

import (
	"net/http"
	"vhagar/config"
)

func newNacos(cfg *config.CfgType) *Nacos {
	return &Nacos{
		Global: cfg.Global,
		Config: cfg.Nacos,
	}
}

type Nacos struct {
	config.Global
	Config      config.NacosCfg
	Client      http.Client
	Host        string
	Token       string
	Clusterdata ClusterStatus
}

type ClusterStatus struct {
	HealthInstance   []ServerInstance
	UnHealthInstance []ServerInstance
}

type Service struct {
	Doms  []string `json:"doms"`
	Count int      `json:"count"`
}

type Instance struct {
	GroupName       string      `json:"groupName"`
	Hosts           []Instances `json:"hosts"`
	Dom             string      `json:"dom"`
	Name            string      `json:"name"`
	CacheMillis     int         `json:"cacheMillis"`
	LastRefTime     int64       `json:"lastRefTime"`
	Checksum        string      `json:"checksum"`
	UseSpecifiedURL bool        `json:"useSpecifiedURL"`
	Clusters        string      `json:"clusters"`
	Env             string      `json:"env"`
	Metadata        map[string]interface{}
}

type Instances struct {
	Ip                        string `json:"ip"`
	Port                      int    `json:"port"`
	Valid                     bool   `json:"valid"`
	Healthy                   bool   `json:"healthy"`
	Marked                    bool   `json:"marked"`
	InstanceId                string `json:"instanceId"`
	Metadata                  map[string]string
	Enabled                   bool    `json:"enabled"`
	Weight                    float32 `json:"weight"`
	ClusterName               string  `json:"clusterName"`
	ServiceName               string  `json:"serviceName"`
	Ephemeral                 bool    `json:"ephemeral"`
	InstanceHeartBeatInterval int64   `json:"instanceHeartBeatInterval"`
}
type ServerInstance struct {
	NamespaceName string `json:"namespaceName"`
	ServiceName   string `json:"serviceName"`
	IpAddr        string `json:"ipAddr"`
	Health        string `json:"health"`
	Hostname      string `json:"hostname"`
	Weight        string `json:"weight"`
	Ip            string `json:"ip"`
	Port          string `json:"port"`
	GroupName     string `json:"groupName"`
}

type Nacostarget struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

type Nacosfile struct {
	Data []Nacostarget
}
