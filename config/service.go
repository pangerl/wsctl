// Package config 服务配置定义
// 包含所有外部服务的配置结构体和相关功能
// @Author lanpang
// @Date 2024/12/19
// @Desc 服务相关配置，包括AI、天气、消息队列、服务发现等
package config

import (
	"fmt"
	"strings"
	"time"

	"vhagar/errors"
)

// ServiceConfig 服务配置接口
// 定义所有服务配置必须实现的方法
type ServiceConfig interface {
	// Validate 验证配置的有效性
	Validate() error
	// IsEnabled 检查服务是否启用
	IsEnabled() bool
}

// AIServiceConfig AI服务详细配置
// 扩展了基本的AIConfig，提供更多配置选项
type AIServiceConfig struct {
	Enable      bool                      `toml:"enable"`      // 是否启用AI功能
	Provider    string                    `toml:"provider"`    // 当前使用的提供商
	Providers   map[string]ProviderConfig `toml:"providers"`   // 提供商配置映射
	Timeout     time.Duration             `toml:"timeout"`     // 请求超时时间
	MaxRetries  int                       `toml:"max_retries"` // 最大重试次数
	Temperature float64                   `toml:"temperature"` // 生成温度参数
	MaxTokens   int                       `toml:"max_tokens"`  // 最大令牌数
}

// Validate 验证AI服务配置
func (a *AIServiceConfig) Validate() error {
	if !a.Enable {
		return nil // 未启用时不需要验证
	}

	if a.Provider == "" {
		return errors.New(errors.ErrCodeConfigInvalid, "AI提供商不能为空")
	}

	provider, exists := a.Providers[a.Provider]
	if !exists {
		return errors.NewWithDetail(
			errors.ErrCodeConfigInvalid,
			"指定的AI提供商不存在",
			fmt.Sprintf("提供商: %s", a.Provider),
		)
	}

	if err := provider.Validate(); err != nil {
		return errors.WrapWithDetail(
			errors.ErrCodeConfigInvalid,
			"AI提供商配置无效",
			fmt.Sprintf("提供商: %s", a.Provider),
			err,
		)
	}

	// 验证参数范围
	if a.Temperature < 0 || a.Temperature > 2 {
		return errors.New(errors.ErrCodeConfigInvalid, "AI温度参数必须在0-2之间")
	}

	if a.MaxTokens < 1 || a.MaxTokens > 32000 {
		return errors.New(errors.ErrCodeConfigInvalid, "AI最大令牌数必须在1-32000之间")
	}

	return nil
}

// IsEnabled 检查AI服务是否启用
func (a *AIServiceConfig) IsEnabled() bool {
	return a.Enable
}

// GetCurrentProvider 获取当前AI提供商配置
func (a *AIServiceConfig) GetCurrentProvider() (ProviderConfig, error) {
	if !a.Enable {
		return ProviderConfig{}, errors.New(errors.ErrCodeConfigInvalid, "AI服务未启用")
	}

	provider, exists := a.Providers[a.Provider]
	if !exists {
		return ProviderConfig{}, errors.NewWithDetail(
			errors.ErrCodeConfigInvalid,
			"AI提供商不存在",
			fmt.Sprintf("提供商: %s", a.Provider),
		)
	}

	return provider, nil
}

// ProviderConfig LLM服务商详细配置
type ProviderConfig struct {
	ApiKey      string            `toml:"api_key"`     // API密钥
	ApiUrl      string            `toml:"api_url"`     // API地址
	Model       string            `toml:"model"`       // 模型名称
	Timeout     time.Duration     `toml:"timeout"`     // 请求超时时间
	MaxRetries  int               `toml:"max_retries"` // 最大重试次数
	Headers     map[string]string `toml:"headers"`     // 自定义请求头
	RateLimit   int               `toml:"rate_limit"`  // 速率限制（每分钟请求数）
	Description string            `toml:"description"` // 提供商描述
}

// Validate 验证提供商配置
func (p *ProviderConfig) Validate() error {
	if p.ApiKey == "" {
		return errors.New(errors.ErrCodeConfigInvalid, "API密钥不能为空")
	}

	if p.ApiUrl == "" {
		return errors.New(errors.ErrCodeConfigInvalid, "API地址不能为空")
	}

	if p.Model == "" {
		return errors.New(errors.ErrCodeConfigInvalid, "模型名称不能为空")
	}

	// 验证URL格式
	if !strings.HasPrefix(p.ApiUrl, "http://") && !strings.HasPrefix(p.ApiUrl, "https://") {
		return errors.New(errors.ErrCodeConfigInvalid, "API地址必须以http://或https://开头")
	}

	return nil
}

// WeatherServiceConfig 天气服务详细配置
type WeatherServiceConfig struct {
	Enable      bool              `toml:"enable"`       // 是否启用天气服务
	Provider    string            `toml:"provider"`     // 天气服务提供商
	ApiHost     string            `toml:"api_host"`     // API主机地址
	ApiKey      string            `toml:"api_key"`      // API密钥
	Timeout     time.Duration     `toml:"timeout"`      // 请求超时时间
	MaxRetries  int               `toml:"max_retries"`  // 最大重试次数
	CacheExpiry time.Duration     `toml:"cache_expiry"` // 缓存过期时间
	Headers     map[string]string `toml:"headers"`      // 自定义请求头
	DefaultCity string            `toml:"default_city"` // 默认城市
}

// Validate 验证天气服务配置
func (w *WeatherServiceConfig) Validate() error {
	if !w.Enable {
		return nil
	}

	if w.ApiHost == "" {
		return errors.New(errors.ErrCodeConfigInvalid, "天气服务API主机地址不能为空")
	}

	if w.ApiKey == "" {
		return errors.New(errors.ErrCodeConfigInvalid, "天气服务API密钥不能为空")
	}

	if !strings.HasPrefix(w.ApiHost, "http://") && !strings.HasPrefix(w.ApiHost, "https://") {
		return errors.New(errors.ErrCodeConfigInvalid, "天气服务API地址必须以http://或https://开头")
	}

	return nil
}

// IsEnabled 检查天气服务是否启用
func (w *WeatherServiceConfig) IsEnabled() bool {
	return w.Enable
}

// RocketMQServiceConfig RocketMQ服务详细配置
type RocketMQServiceConfig struct {
	Enable            bool          `toml:"enable"`            // 是否启用RocketMQ
	RocketmqDashboard string        `toml:"rocketmqdashboard"` // RocketMQ控制台地址
	NameServer        string        `toml:"name_server"`       // NameServer地址
	Username          string        `toml:"username"`          // 用户名
	Password          string        `toml:"password"`          // 密码
	ProducerGroup     string        `toml:"producer_group"`    // 生产者组
	ConsumerGroup     string        `toml:"consumer_group"`    // 消费者组
	DefaultTopic      string        `toml:"default_topic"`     // 默认主题
	DefaultTag        string        `toml:"default_tag"`       // 默认标签
	SendTimeout       time.Duration `toml:"send_timeout"`      // 发送超时时间
	ConsumeTimeout    time.Duration `toml:"consume_timeout"`   // 消费超时时间
	MaxRetries        int           `toml:"max_retries"`       // 最大重试次数
	BatchSize         int           `toml:"batch_size"`        // 批量大小
}

// Validate 验证RocketMQ服务配置
func (r *RocketMQServiceConfig) Validate() error {
	if !r.Enable {
		return nil
	}

	if r.NameServer == "" {
		return errors.New(errors.ErrCodeConfigInvalid, "RocketMQ NameServer地址不能为空")
	}

	if r.ProducerGroup == "" {
		return errors.New(errors.ErrCodeConfigInvalid, "RocketMQ生产者组不能为空")
	}

	if r.ConsumerGroup == "" {
		return errors.New(errors.ErrCodeConfigInvalid, "RocketMQ消费者组不能为空")
	}

	return nil
}

// IsEnabled 检查RocketMQ服务是否启用
func (r *RocketMQServiceConfig) IsEnabled() bool {
	return r.Enable
}

// NacosServiceConfig Nacos服务发现详细配置
type NacosServiceConfig struct {
	Enable     bool          `toml:"enable"`      // 是否启用Nacos
	Server     string        `toml:"server"`      // Nacos服务器地址
	Port       int           `toml:"port"`        // Nacos端口
	Username   string        `toml:"username"`    // 用户名
	Password   string        `toml:"password"`    // 密码
	Namespace  string        `toml:"namespace"`   // 命名空间
	Group      string        `toml:"group"`       // 分组
	Writefile  string        `toml:"writefile"`   // 写入文件路径
	Timeout    time.Duration `toml:"timeout"`     // 连接超时时间
	LogLevel   string        `toml:"log_level"`   // 日志级别
	CacheDir   string        `toml:"cache_dir"`   // 缓存目录
	UpdateTime time.Duration `toml:"update_time"` // 更新间隔
}

// Validate 验证Nacos服务配置
func (n *NacosServiceConfig) Validate() error {
	if !n.Enable {
		return nil
	}

	if n.Server == "" {
		return errors.New(errors.ErrCodeConfigInvalid, "Nacos服务器地址不能为空")
	}

	if n.Port <= 0 || n.Port > 65535 {
		return errors.New(errors.ErrCodeConfigInvalid, "Nacos端口必须在1-65535之间")
	}

	if n.Username == "" || n.Password == "" {
		return errors.New(errors.ErrCodeConfigInvalid, "Nacos用户名和密码不能为空")
	}

	return nil
}

// IsEnabled 检查Nacos服务是否启用
func (n *NacosServiceConfig) IsEnabled() bool {
	return n.Enable
}

// GetServerAddress 获取完整的服务器地址
func (n *NacosServiceConfig) GetServerAddress() string {
	if n.Port > 0 {
		return fmt.Sprintf("%s:%d", n.Server, n.Port)
	}
	return n.Server
}

// MetricServiceConfig 监控服务详细配置
type MetricServiceConfig struct {
	Enable      bool          `toml:"enable"`       // 是否启用监控
	Port        string        `toml:"port"`         // 监控端口
	HealthApi   string        `toml:"health_api"`   // 健康检查API路径
	MetricsPath string        `toml:"metrics_path"` // 指标路径
	Timeout     time.Duration `toml:"timeout"`      // 请求超时时间
	Interval    time.Duration `toml:"interval"`     // 采集间隔
	EnablePprof bool          `toml:"enable_pprof"` // 是否启用pprof
	PprofPrefix string        `toml:"pprof_prefix"` // pprof路径前缀
	BasicAuth   BasicAuth     `toml:"basic_auth"`   // 基础认证配置
	TLS         TLSConfig     `toml:"tls"`          // TLS配置
}

// BasicAuth 基础认证配置
type BasicAuth struct {
	Enable   bool   `toml:"enable"`   // 是否启用基础认证
	Username string `toml:"username"` // 用户名
	Password string `toml:"password"` // 密码
}

// TLSConfig TLS配置
type TLSConfig struct {
	Enable   bool   `toml:"enable"`    // 是否启用TLS
	CertFile string `toml:"cert_file"` // 证书文件路径
	KeyFile  string `toml:"key_file"`  // 私钥文件路径
}

// Validate 验证监控服务配置
func (m *MetricServiceConfig) Validate() error {
	if !m.Enable {
		return nil
	}

	if m.Port == "" {
		return errors.New(errors.ErrCodeConfigInvalid, "监控端口不能为空")
	}

	// 验证基础认证配置
	if m.BasicAuth.Enable {
		if m.BasicAuth.Username == "" || m.BasicAuth.Password == "" {
			return errors.New(errors.ErrCodeConfigInvalid, "启用基础认证时用户名和密码不能为空")
		}
	}

	// 验证TLS配置
	if m.TLS.Enable {
		if m.TLS.CertFile == "" || m.TLS.KeyFile == "" {
			return errors.New(errors.ErrCodeConfigInvalid, "启用TLS时证书文件和私钥文件不能为空")
		}
	}

	return nil
}

// IsEnabled 检查监控服务是否启用
func (m *MetricServiceConfig) IsEnabled() bool {
	return m.Enable
}

// GetListenAddress 获取监听地址
func (m *MetricServiceConfig) GetListenAddress() string {
	return ":" + m.Port
}

// DefaultServiceConfigs 返回默认的服务配置
func DefaultServiceConfigs() map[string]ServiceConfig {
	return map[string]ServiceConfig{
		"ai": &AIServiceConfig{
			Enable:      false,
			Timeout:     30 * time.Second,
			MaxRetries:  3,
			Temperature: 0.7,
			MaxTokens:   2000,
		},
		"weather": &WeatherServiceConfig{
			Enable:      false,
			Timeout:     10 * time.Second,
			MaxRetries:  3,
			CacheExpiry: 30 * time.Minute,
		},
		"rocketmq": &RocketMQServiceConfig{
			Enable:         false,
			SendTimeout:    3 * time.Second,
			ConsumeTimeout: 30 * time.Second,
			MaxRetries:     3,
			BatchSize:      32,
		},
		"nacos": &NacosServiceConfig{
			Enable:     false,
			Port:       8848,
			Timeout:    5 * time.Second,
			LogLevel:   "info",
			UpdateTime: 30 * time.Second,
		},
		"metric": &MetricServiceConfig{
			Enable:      false,
			Port:        "8090",
			HealthApi:   "/health",
			MetricsPath: "/metrics",
			Timeout:     5 * time.Second,
			Interval:    15 * time.Second,
			EnablePprof: false,
			PprofPrefix: "/debug/pprof",
		},
	}
}

// ValidateAllServices 验证所有服务配置
func ValidateAllServices(services map[string]ServiceConfig) error {
	for name, service := range services {
		if err := service.Validate(); err != nil {
			return errors.WrapWithDetail(
				errors.ErrCodeConfigInvalid,
				"服务配置验证失败",
				fmt.Sprintf("服务: %s", name),
				err,
			)
		}
	}
	return nil
}
