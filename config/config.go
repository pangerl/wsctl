// Package config @Author lanpang
// @Date 2024/9/11 上午11:15:00
// @Desc
package config

import (
	"fmt"
	"log"
	"os"
	"time"
	"vhagar/libs"

	"github.com/BurntSushi/toml"
	//"vhagar/task/nacos"
)

const VERSION = "v5.0"

var (
	Config *CfgType
)

type CfgType struct {
	Global
	DomainListName  string             `toml:"domainListName"`
	NasDir          string             `toml:"nasDir"`
	VictoriaMetrics string             `toml:"victoriaMetrics"`
	Cron            map[string]Crontab `toml:"cron"`
	Nacos           NacosCfg           `toml:"nacos"`
	Tenant          Tenant             `toml:"tenant"`
	PG              libs.DB            `toml:"pg"`
	ES              libs.DB            `toml:"es"`
	Customer        libs.DB            `toml:"customer"`
	Doris           DorisCfg           `toml:"doris"`
	RocketMQ        RocketMQCfg        `toml:"rocketmq"`
	Metric          MetricCfg          `toml:"metric"`
	Redis           libs.RedisConfig   `toml:"redis"`

	AI      AICfg      `toml:"ai"`
	Weather WeatherCfg `toml:"weather"`
}

// 新增 AI 配置结构体，支持多套 LLM 配置
// ai = { enable = true, provider = "openrouter", providers = { openrouter = { api_key = "sk-xxx", api_url = "https://openrouter.ai/api/v1/chat/completions", model = "gpt-3.5-turbo" }, openai = { api_key = "sk-xxx", api_url = "https://api.openai.com/v1/chat/completions", model = "gpt-3.5-turbo" } } }
type AICfg struct {
	Enable    bool                   `toml:"enable"`
	Provider  string                 `toml:"provider"`
	Providers map[string]ProviderCfg `toml:"providers"`
}

type WeatherCfg struct {
	ApiHost string `toml:"api_host"`
	ApiKey  string `toml:"api_key"`
}

// LLM 服务商配置
type ProviderCfg struct {
	ApiKey string `toml:"api_key"`
	ApiUrl string `toml:"api_url"`
	Model  string `toml:"model"`
}

type Global struct {
	LogLevel    string        `toml:"logLevel"`
	LogToFile   bool          `toml:"logToFile"`
	ProjectName string        `toml:"projectname"`
	ProxyURL    string        `toml:"proxyurl"`
	Notify      Notify        `toml:"notify"`
	Watch       bool          `toml:"watch"`
	Report      bool          `toml:"report"`
	Interval    time.Duration `toml:"interval"`
	Duration    time.Duration `toml:"duration"`
}

type Crontab struct {
	Crontab    bool   `toml:"crontab"`
	Scheducron string `toml:"scheducron"`
}

type Notify struct {
	Robotkey []string            `toml:"robotkey"`
	Userlist []string            `toml:"userlist"`
	Notifier map[string]Notifier `toml:"notifier"`
}

type Notifier struct {
	Robotkey []string `json:"robotkey"`
}

type MetricCfg struct {
	Enable    bool
	Port      string
	HealthApi string
}

func InitConfig(cfgFile string) (*CfgType, error) {
	//configFile := path.Join(configDir, "config.toml")
	Config = &CfgType{}

	log.Printf("读取配置文件 %s \n", cfgFile)
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("Failed Info: 配置文件格式错误 %s", err)
			return
		}
	}()
	if _, err := os.Stat(cfgFile); err != nil {
		if os.IsNotExist(err) {
			//log.Fatalf("读取配置文件 %s 失败，报错：%s", cfgFile, err)
			return nil, fmt.Errorf("configuration file(%s) not found", cfgFile)
		}
	} else {
		if _, err := toml.DecodeFile(cfgFile, Config); err != nil {
			//log.Fatalf("Failed Info: 配置文件格式错误 %s", err)
			return nil, fmt.Errorf("failed to load configs of dir: %s err:%s", cfgFile, err)
		}
		//log.Println(Config.Notify)
	}
	// log.Println("配置文件加载成功", "config", Config)
	return Config, nil
}
