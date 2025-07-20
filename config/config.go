// Package config 提供应用程序的配置管理功能
// 包含配置文件加载、解析、验证和访问的完整实现
// @Author lanpang
// @Date 2024/12/19
// @Desc 重构后的配置管理包，整合所有配置项到统一结构中
package config

import (
	"fmt"
	"os"
	"time"

	"vhagar/database"
	"vhagar/errors"
	"vhagar/logger"

	"github.com/BurntSushi/toml"
)

// VERSION 应用程序版本号
const VERSION = "v5.0"

// 全局配置实例
var (
	Config *AppConfig
)

// AppConfig 主应用程序配置结构体
// 整合了所有配置项到统一的结构中
type AppConfig struct {
	// 全局配置
	Global GlobalConfig `toml:",inline"` // 内联全局配置字段

	// 应用特定配置
	DomainListName  string `toml:"domainListName"`  // 出网域名检测列表文件名
	NasDir          string `toml:"nasDir"`          // 会话存档文件路径
	VictoriaMetrics string `toml:"victoriaMetrics"` // VM数据库地址

	// 定时任务配置
	Cron map[string]CrontabConfig `toml:"cron"`

	// 租户配置
	Tenant TenantConfig `toml:"tenant"`

	// 数据库配置
	Database DatabaseConfigs `toml:",inline"` // 内联数据库配置

	// 外部服务配置
	Services ServiceConfigs `toml:",inline"` // 内联服务配置

	// 监控配置
	Metric MetricConfig `toml:"metric"`
}

// GlobalConfig 全局应用程序设置
type GlobalConfig struct {
	LogLevel    string        `toml:"logLevel"`    // 日志级别 (debug, info, warn, error)
	LogToFile   bool          `toml:"logToFile"`   // 日志是否保存到本地文件
	ProjectName string        `toml:"projectname"` // 项目名称
	ProxyURL    string        `toml:"proxyurl"`    // 网络代理URL
	Watch       bool          `toml:"watch"`       // 启用配置文件变更监控
	Report      bool          `toml:"report"`      // 启用报告功能
	Interval    time.Duration `toml:"interval"`    // 周期性任务的默认间隔
	Duration    time.Duration `toml:"duration"`    // 操作的默认持续时间
	Notify      NotifyConfig  `toml:"notify"`      // 通知配置
}

// CrontabConfig 定时任务配置
type CrontabConfig struct {
	Crontab    bool   `toml:"crontab"`    // 是否启用定时任务
	Scheducron string `toml:"scheducron"` // Cron 表达式
}

// NotifyConfig 通知配置
type NotifyConfig struct {
	Robotkey []string                  `toml:"robotkey"` // 默认机器人密钥列表
	Userlist []string                  `toml:"userlist"` // 默认告警@人列表
	Notifier map[string]NotifierConfig `toml:"notifier"` // 特定通知器配置
}

// NotifierConfig 通知器配置
type NotifierConfig struct {
	Robotkey []string `json:"robotkey"` // 机器人密钥列表
}

// TenantConfig 租户配置
type TenantConfig struct {
	Corp []*CorpConfig `toml:"corp"` // 企业配置列表
}

// CorpConfig 企业配置
type CorpConfig struct {
	Corpid      string `json:"corpid"`      // 企业ID
	Convenabled bool   `json:"convenabled"` // 是否开通会话存档功能
}

// DatabaseConfigs 数据库配置集合
type DatabaseConfigs struct {
	PG       database.Config      `toml:"pg"`       // PostgreSQL配置
	ES       database.Config      `toml:"es"`       // Elasticsearch配置
	Customer database.Config      `toml:"customer"` // 客户数据库配置
	Doris    DorisConfig          `toml:"doris"`    // Doris配置
	Redis    database.RedisConfig `toml:"redis"`    // Redis配置
}

// DorisConfig Doris数据库配置
type DorisConfig struct {
	database.Config
	HttpPort int `toml:"httpport"` // HTTP端口
}

// ServiceConfigs 外部服务配置集合
type ServiceConfigs struct {
	AI       AIConfig       `toml:"ai"`       // AI服务配置
	Weather  WeatherConfig  `toml:"weather"`  // 天气服务配置
	RocketMQ RocketMQConfig `toml:"rocketmq"` // RocketMQ配置
	Nacos    NacosConfig    `toml:"nacos"`    // Nacos配置
}

// AIConfig AI服务配置
// 支持多套 LLM 配置
type AIConfig struct {
	Enable    bool                      `toml:"enable"`    // 是否启用AI功能
	Provider  string                    `toml:"provider"`  // 当前使用的提供商
	Providers map[string]ProviderConfig `toml:"providers"` // 提供商配置映射
}

// ProviderConfig LLM服务商配置
type ProviderConfig struct {
	ApiKey string `toml:"api_key"` // API密钥
	ApiUrl string `toml:"api_url"` // API地址
	Model  string `toml:"model"`   // 模型名称
}

// WeatherConfig 天气服务配置
type WeatherConfig struct {
	ApiHost string `toml:"api_host"` // API主机地址
	ApiKey  string `toml:"api_key"`  // API密钥
}

// RocketMQConfig RocketMQ配置
type RocketMQConfig struct {
	RocketmqDashboard string `toml:"rocketmqdashboard"` // RocketMQ控制台地址
	Username          string `json:"username"`          // 用户名
	Password          string `json:"password"`          // 密码
}

// NacosConfig Nacos服务发现配置
type NacosConfig struct {
	Server    string `json:"server"`    // Nacos服务器地址
	Username  string `json:"username"`  // 用户名
	Password  string `json:"password"`  // 密码
	Namespace string `json:"namespace"` // 命名空间
	Writefile string // 写入文件路径
}

// MetricConfig 监控配置
type MetricConfig struct {
	Enable    bool   `toml:"enable"`    // 是否启用监控
	Port      string `toml:"port"`      // 监控端口
	HealthApi string `toml:"healthApi"` // 健康检查API路径
}

// InitConfig 初始化配置文件
// 参数:
//   - cfgFile: 配置文件路径
//
// 返回:
//   - *AppConfig: 配置对象
//   - error: 错误信息
func InitConfig(cfgFile string) (*AppConfig, error) {
	Config = &AppConfig{}

	// 检查配置文件是否存在
	if _, err := os.Stat(cfgFile); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.NewWithDetail(
				errors.ErrCodeConfigNotFound,
				"配置文件未找到",
				fmt.Sprintf("配置文件路径: %s", cfgFile),
			)
		}
		return nil, errors.Wrap(errors.ErrCodeConfigInvalid, "检查配置文件失败", err)
	}

	// 解析配置文件
	if _, err := toml.DecodeFile(cfgFile, Config); err != nil {
		return nil, errors.WrapWithDetail(
			errors.ErrCodeConfigInvalid,
			"配置文件格式错误",
			fmt.Sprintf("文件: %s", cfgFile),
			err,
		)
	}

	// 验证配置
	if err := Config.Validate(); err != nil {
		return nil, errors.Wrap(errors.ErrCodeConfigInvalid, "配置验证失败", err)
	}

	// 初始化日志器
	loggerConfig := logger.Config{
		Level:  Config.LogLevel,
		ToFile: Config.LogToFile,
		Format: "console",
	}
	if err := logger.InitLogger(loggerConfig); err != nil {
		return nil, errors.Wrap(errors.ErrCodeInternalErr, "初始化日志器失败", err)
	}

	log := logger.GetLogger()
	log.Infow("配置文件加载成功", "file", cfgFile, "version", VERSION)

	return Config, nil
}

// Validate 验证配置的完整性和正确性
func (c *AppConfig) Validate() error {
	// 验证项目名称
	if c.ProjectName == "" {
		return errors.New(errors.ErrCodeConfigInvalid, "项目名称不能为空")
	}

	// 验证日志级别
	validLogLevels := []string{"debug", "info", "warn", "error"}
	isValidLevel := false
	for _, level := range validLogLevels {
		if c.LogLevel == level {
			isValidLevel = true
			break
		}
	}
	if !isValidLevel {
		return errors.NewWithDetail(
			errors.ErrCodeConfigInvalid,
			"无效的日志级别",
			fmt.Sprintf("支持的级别: %v", validLogLevels),
		)
	}

	// 验证AI配置
	if c.Services.AI.Enable {
		if c.Services.AI.Provider == "" {
			return errors.New(errors.ErrCodeConfigInvalid, "启用AI时必须指定提供商")
		}

		provider, exists := c.Services.AI.Providers[c.Services.AI.Provider]
		if !exists {
			return errors.NewWithDetail(
				errors.ErrCodeConfigInvalid,
				"指定的AI提供商不存在",
				fmt.Sprintf("提供商: %s", c.Services.AI.Provider),
			)
		}

		if provider.ApiKey == "" || provider.ApiUrl == "" {
			return errors.New(errors.ErrCodeConfigInvalid, "AI提供商配置不完整")
		}
	}

	return nil
}

// GetDatabaseConfig 获取指定类型的数据库配置
func (c *AppConfig) GetDatabaseConfig(dbType string) any {
	switch dbType {
	case "pg", "postgresql":
		return c.Database.PG
	case "es", "elasticsearch":
		return c.Database.ES
	case "customer":
		return c.Database.Customer
	case "doris":
		return c.Database.Doris
	case "redis":
		return c.Database.Redis
	default:
		return nil
	}
}

// IsAIEnabled 检查AI功能是否启用
func (c *AppConfig) IsAIEnabled() bool {
	return c.Services.AI.Enable
}

// GetAIProvider 获取当前AI提供商配置
func (c *AppConfig) GetAIProvider() (ProviderConfig, bool) {
	if !c.IsAIEnabled() {
		return ProviderConfig{}, false
	}

	provider, exists := c.Services.AI.Providers[c.Services.AI.Provider]
	return provider, exists
}

// IsMetricEnabled 检查监控功能是否启用
func (c *AppConfig) IsMetricEnabled() bool {
	return c.Metric.Enable
}

// GetNotifierConfig 获取指定通知器的配置
func (c *AppConfig) GetNotifierConfig(name string) (NotifierConfig, bool) {
	notifier, exists := c.Notify.Notifier[name]
	return notifier, exists
}

// GetCronConfig 获取指定定时任务的配置
func (c *AppConfig) GetCronConfig(name string) (CrontabConfig, bool) {
	cron, exists := c.Cron[name]
	return cron, exists
}

// DefaultConfig 返回默认配置
func DefaultConfig() *AppConfig {
	return &AppConfig{
		Global: GlobalConfig{
			LogLevel:    "info",
			LogToFile:   false,
			ProjectName: "vhagar",
			Watch:       true,
			Report:      true,
			Interval:    5 * time.Minute,
			Duration:    time.Hour,
		},
		Metric: MetricConfig{
			Enable: false,
			Port:   "8090",
		},
	}
}
