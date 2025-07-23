// Package config 配置验证逻辑
// 提供全面的配置验证功能，包括必填项检查、格式验证、业务逻辑验证等
// @Author lanpang
// @Date 2024/12/19
// @Desc 配置验证器，确保配置的完整性和正确性
package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"vhagar/database"
	"vhagar/errors"
)

// ValidationRule 验证规则接口
type ValidationRule interface {
	// Validate 执行验证
	Validate(value interface{}) error
	// GetName 获取规则名称
	GetName() string
	// GetDescription 获取规则描述
	GetDescription() string
}

// ConfigValidator 配置验证器
type ConfigValidator struct {
	rules   map[string][]ValidationRule // 字段验证规则映射
	errors  []ValidationError           // 验证错误列表
	strict  bool                        // 是否严格模式
	context map[string]interface{}      // 验证上下文
}

// ValidationError 验证错误
type ValidationError struct {
	Field      string `json:"field"`      // 字段名
	Value      string `json:"value"`      // 字段值
	Rule       string `json:"rule"`       // 违反的规则
	Message    string `json:"message"`    // 错误消息
	Severity   string `json:"severity"`   // 严重程度 (error, warning, info)
	Suggestion string `json:"suggestion"` // 修复建议
}

// Error 实现 error 接口
func (ve ValidationError) Error() string {
	return fmt.Sprintf("[%s] %s: %s (值: %s)", ve.Severity, ve.Field, ve.Message, ve.Value)
}

// NewConfigValidator 创建新的配置验证器
func NewConfigValidator(strict bool) *ConfigValidator {
	return &ConfigValidator{
		rules:   make(map[string][]ValidationRule),
		errors:  make([]ValidationError, 0),
		strict:  strict,
		context: make(map[string]interface{}),
	}
}

// AddRule 添加验证规则
func (cv *ConfigValidator) AddRule(field string, rule ValidationRule) {
	if cv.rules[field] == nil {
		cv.rules[field] = make([]ValidationRule, 0)
	}
	cv.rules[field] = append(cv.rules[field], rule)
}

// SetContext 设置验证上下文
func (cv *ConfigValidator) SetContext(key string, value interface{}) {
	cv.context[key] = value
}

// ValidateConfig 验证完整配置
func (cv *ConfigValidator) ValidateConfig(config *AppConfig) error {
	cv.errors = cv.errors[:0] // 清空错误列表

	// 设置验证上下文
	cv.SetContext("config", config)

	// 验证全局配置
	cv.validateGlobalConfig(&config.Global)

	// 验证数据库配置
	cv.validateDatabaseConfigs(&config.Database)

	// 验证服务配置
	cv.validateServiceConfigs(&config.Services)

	// 验证监控配置
	cv.validateMetricConfig(&config.Metric)

	// 验证租户配置
	cv.validateTenantConfig(&config.Tenant)

	// 验证定时任务配置
	cv.validateCronConfigs(config.Cron)

	// 验证通知配置
	cv.validateNotifyConfig(&config.Global.Notify)

	// 检查是否有错误
	if len(cv.errors) > 0 {
		return cv.buildValidationError()
	}

	return nil
}

// validateGlobalConfig 验证全局配置
func (cv *ConfigValidator) validateGlobalConfig(global *GlobalConfig) {
	// 验证项目名称
	cv.validateRequired("projectname", global.ProjectName)
	cv.validateStringLength("projectname", global.ProjectName, 1, 100)

	// 验证日志级别
	cv.validateRequired("logLevel", global.LogLevel)
	cv.validateEnum("logLevel", global.LogLevel, []string{"debug", "info", "warn", "error"})

	// 验证代理URL
	if global.ProxyURL != "" {
		cv.validateURL("proxyurl", global.ProxyURL)
	}

	// 验证时间间隔
	cv.validateDuration("interval", global.Interval, time.Second, 24*time.Hour)
	cv.validateDuration("duration", global.Duration, time.Minute, 7*24*time.Hour)
}

// validateDatabaseConfigs 验证数据库配置
func (cv *ConfigValidator) validateDatabaseConfigs(db *DatabaseConfigs) {
	// 验证PostgreSQL配置
	if db.PG.HasValue() {
		cv.validateDatabaseConfig("pg", &db.PG)
	}

	// 验证Elasticsearch配置
	if db.ES.HasValue() {
		cv.validateDatabaseConfig("es", &db.ES)
	}

	// 验证客户数据库配置
	if db.Customer.HasValue() {
		cv.validateDatabaseConfig("customer", &db.Customer)
	}

	// 验证Doris配置
	if db.Doris.HasValue() {
		cv.validateDorisConfig(&db.Doris)
	}

	// 验证Redis配置
	if db.Redis.HasValue() {
		cv.validateRedisConfig(&db.Redis)
	}
}

// validateDatabaseConfig 验证通用数据库配置
func (cv *ConfigValidator) validateDatabaseConfig(prefix string, config *database.Config) {
	cv.validateRequired(prefix+".ip", config.Host)
	cv.validateIP(prefix+".ip", config.Host)
	cv.validatePort(prefix+".port", config.Port)
	cv.validateRequired(prefix+".username", config.Username)
	cv.validateRequired(prefix+".password", config.Password)
}

// validateDorisConfig 验证Doris配置
func (cv *ConfigValidator) validateDorisConfig(config *DorisConfig) {
	cv.validateDatabaseConfig("doris", &config.Config)
	cv.validatePort("doris.httpport", config.HttpPort)
}

// validateRedisConfig 验证Redis配置
func (cv *ConfigValidator) validateRedisConfig(config *database.RedisConfig) {
	cv.validateRequired("redis.addr", config.Addr)
	cv.validateRedisAddr("redis.addr", config.Addr)
	cv.validateRange("redis.db", config.DB, 0, 15)
}

// validateServiceConfigs 验证服务配置
func (cv *ConfigValidator) validateServiceConfigs(services *ServiceConfigs) {
	// 验证AI配置
	cv.validateAIConfig(&services.AI)

	// 验证天气配置
	cv.validateWeatherConfig(&services.Weather)

	// 验证RocketMQ配置
	cv.validateRocketMQConfig(&services.RocketMQ)

	// 验证Nacos配置
	cv.validateNacosConfig(&services.Nacos)
}

// validateAIConfig 验证AI配置
func (cv *ConfigValidator) validateAIConfig(ai *AIConfig) {
	if !ai.Enable {
		return
	}

	cv.validateRequired("ai.provider", ai.Provider)

	// 验证提供商是否存在
	if ai.Provider != "" {
		if _, exists := ai.Providers[ai.Provider]; !exists {
			cv.addError("ai.provider", ai.Provider, "provider_not_found",
				fmt.Sprintf("指定的AI提供商 '%s' 不存在", ai.Provider), "error",
				"请检查providers配置中是否包含该提供商")
		}
	}

	// 验证每个提供商配置
	for name, provider := range ai.Providers {
		prefix := fmt.Sprintf("ai.providers.%s", name)
		cv.validateRequired(prefix+".api_key", provider.ApiKey)
		cv.validateRequired(prefix+".api_url", provider.ApiUrl)
		cv.validateURL(prefix+".api_url", provider.ApiUrl)
		cv.validateRequired(prefix+".model", provider.Model)
	}
}

// validateWeatherConfig 验证天气配置
func (cv *ConfigValidator) validateWeatherConfig(weather *WeatherConfig) {
	if weather.ApiHost != "" {
		cv.validateURL("weather.api_host", weather.ApiHost)
	}
	if weather.ApiKey != "" {
		cv.validateStringLength("weather.api_key", weather.ApiKey, 10, 200)
	}
}

// validateRocketMQConfig 验证RocketMQ配置
func (cv *ConfigValidator) validateRocketMQConfig(rocketmq *RocketMQConfig) {
	if rocketmq.RocketmqDashboard != "" {
		cv.validateURL("rocketmq.rocketmqdashboard", rocketmq.RocketmqDashboard)
	}
}

// validateNacosConfig 验证Nacos配置
func (cv *ConfigValidator) validateNacosConfig(nacos *NacosConfig) {
	if nacos.Server != "" {
		cv.validateURL("nacos.server", nacos.Server)
	}
}

// validateMetricConfig 验证监控配置
func (cv *ConfigValidator) validateMetricConfig(metric *MetricConfig) {
	if metric.Enable {
		cv.validateRequired("metric.port", metric.Port)
		cv.validatePortString("metric.port", metric.Port)
	}
}

// validateTenantConfig 验证租户配置
func (cv *ConfigValidator) validateTenantConfig(tenant *TenantConfig) {
	for i, corp := range tenant.Corp {
		prefix := fmt.Sprintf("tenant.corp[%d]", i)
		cv.validateRequired(prefix+".corpid", corp.Corpid)
		cv.validateStringLength(prefix+".corpid", corp.Corpid, 1, 100)
	}
}

// validateCronConfigs 验证定时任务配置
func (cv *ConfigValidator) validateCronConfigs(crons map[string]CrontabConfig) {
	for name, cron := range crons {
		if cron.Crontab {
			prefix := fmt.Sprintf("cron.%s", name)
			cv.validateRequired(prefix+".scheducron", cron.Scheducron)
			cv.validateCronExpression(prefix+".scheducron", cron.Scheducron)
		}
	}
}

// validateNotifyConfig 验证通知配置
func (cv *ConfigValidator) validateNotifyConfig(notify *NotifyConfig) {
	// 验证默认机器人密钥
	for i, key := range notify.Robotkey {
		cv.validateUUID(fmt.Sprintf("notify.robotkey[%d]", i), key)
	}

	// 验证通知器配置
	for name, notifier := range notify.Notifier {
		for i, key := range notifier.Robotkey {
			cv.validateUUID(fmt.Sprintf("notify.notifier.%s.robotkey[%d]", name, i), key)
		}
	}
}

// 基础验证方法

// validateRequired 验证必填项
func (cv *ConfigValidator) validateRequired(field string, value string) {
	if strings.TrimSpace(value) == "" {
		cv.addError(field, value, "required",
			fmt.Sprintf("字段 '%s' 是必填项", field), "error",
			"请提供有效的值")
	}
}

// validateStringLength 验证字符串长度
func (cv *ConfigValidator) validateStringLength(field string, value string, min, max int) {
	length := len(value)
	if length < min || length > max {
		cv.addError(field, value, "string_length",
			fmt.Sprintf("字段 '%s' 长度必须在 %d-%d 之间，当前长度: %d", field, min, max, length), "error",
			fmt.Sprintf("请确保字符串长度在 %d-%d 之间", min, max))
	}
}

// validateEnum 验证枚举值
func (cv *ConfigValidator) validateEnum(field string, value string, validValues []string) {
	for _, valid := range validValues {
		if value == valid {
			return
		}
	}
	cv.addError(field, value, "enum",
		fmt.Sprintf("字段 '%s' 的值必须是以下之一: %v", field, validValues), "error",
		fmt.Sprintf("请使用有效值: %v", validValues))
}

// validateURL 验证URL格式
func (cv *ConfigValidator) validateURL(field string, value string) {
	if value == "" {
		return
	}

	if _, err := url.Parse(value); err != nil {
		cv.addError(field, value, "url_format",
			fmt.Sprintf("字段 '%s' 不是有效的URL格式", field), "error",
			"请提供有效的URL，如: http://example.com")
	}
}

// validateIP 验证IP地址
func (cv *ConfigValidator) validateIP(field string, value string) {
	if value == "" {
		return
	}

	if net.ParseIP(value) == nil {
		cv.addError(field, value, "ip_format",
			fmt.Sprintf("字段 '%s' 不是有效的IP地址", field), "error",
			"请提供有效的IP地址，如: 192.168.1.1")
	}
}

// validatePort 验证端口号
func (cv *ConfigValidator) validatePort(field string, port int) {
	if port < 1 || port > 65535 {
		cv.addError(field, strconv.Itoa(port), "port_range",
			fmt.Sprintf("字段 '%s' 端口号必须在 1-65535 之间", field), "error",
			"请提供有效的端口号 (1-65535)")
	}
}

// validatePortString 验证端口号字符串
func (cv *ConfigValidator) validatePortString(field string, portStr string) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		cv.addError(field, portStr, "port_format",
			fmt.Sprintf("字段 '%s' 不是有效的端口号", field), "error",
			"请提供数字格式的端口号")
		return
	}
	cv.validatePort(field, port)
}

// validateDuration 验证时间间隔
func (cv *ConfigValidator) validateDuration(field string, duration, min, max time.Duration) {
	if duration < min || duration > max {
		cv.addError(field, duration.String(), "duration_range",
			fmt.Sprintf("字段 '%s' 时间间隔必须在 %s-%s 之间", field, min, max), "error",
			fmt.Sprintf("请设置合理的时间间隔 (%s-%s)", min, max))
	}
}

// validateRange 验证数值范围
func (cv *ConfigValidator) validateRange(field string, value, min, max int) {
	if value < min || value > max {
		cv.addError(field, strconv.Itoa(value), "range",
			fmt.Sprintf("字段 '%s' 值必须在 %d-%d 之间", field, min, max), "error",
			fmt.Sprintf("请设置 %d-%d 之间的值", min, max))
	}
}

// validateRedisAddr 验证Redis地址格式
func (cv *ConfigValidator) validateRedisAddr(field string, addr string) {
	if addr == "" {
		return
	}

	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		cv.addError(field, addr, "redis_addr_format",
			fmt.Sprintf("字段 '%s' Redis地址格式错误", field), "error",
			"请使用格式: host:port，如: localhost:6379")
		return
	}

	// 验证主机部分
	host := parts[0]
	if host == "" {
		cv.addError(field, addr, "redis_host_empty",
			fmt.Sprintf("字段 '%s' Redis主机地址不能为空", field), "error",
			"请提供有效的主机地址")
		return
	}

	// 验证端口部分
	cv.validatePortString(field+".port", parts[1])
}

// validateCronExpression 验证Cron表达式
func (cv *ConfigValidator) validateCronExpression(field string, expr string) {
	if expr == "" {
		return
	}

	// 简单的Cron表达式验证（5个字段）
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		cv.addError(field, expr, "cron_format",
			fmt.Sprintf("字段 '%s' Cron表达式格式错误", field), "error",
			"请使用标准的5字段Cron格式: 分 时 日 月 周")
	}
}

// validateUUID 验证UUID格式
func (cv *ConfigValidator) validateUUID(field string, uuid string) {
	if uuid == "" {
		return
	}

	// UUID格式验证
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !uuidRegex.MatchString(uuid) {
		cv.addError(field, uuid, "uuid_format",
			fmt.Sprintf("字段 '%s' 不是有效的UUID格式", field), "warning",
			"请使用标准的UUID格式: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")
	}
}

// validateFilePath 验证文件路径
func (cv *ConfigValidator) validateFilePath(field string, path string, mustExist bool) {
	if path == "" {
		return
	}

	// 检查路径格式
	if !filepath.IsAbs(path) && !strings.HasPrefix(path, "./") && !strings.HasPrefix(path, "../") {
		cv.addError(field, path, "path_format",
			fmt.Sprintf("字段 '%s' 路径格式可能不正确", field), "warning",
			"建议使用绝对路径或相对路径")
	}

	// 检查文件是否存在
	if mustExist {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			cv.addError(field, path, "file_not_found",
				fmt.Sprintf("字段 '%s' 指定的文件不存在", field), "error",
				"请确保文件路径正确且文件存在")
		}
	}
}

// addError 添加验证错误
func (cv *ConfigValidator) addError(field, value, rule, message, severity, suggestion string) {
	cv.errors = append(cv.errors, ValidationError{
		Field:      field,
		Value:      value,
		Rule:       rule,
		Message:    message,
		Severity:   severity,
		Suggestion: suggestion,
	})
}

// buildValidationError 构建验证错误
func (cv *ConfigValidator) buildValidationError() error {
	var errorMessages []string
	var warningMessages []string

	for _, err := range cv.errors {
		switch err.Severity {
		case "error":
			errorMessages = append(errorMessages, err.Error())
		case "warning":
			warningMessages = append(warningMessages, err.Error())
		}
	}

	// 在严格模式下，警告也被视为错误
	if cv.strict {
		errorMessages = append(errorMessages, warningMessages...)
	}

	if len(errorMessages) > 0 {
		return errors.NewWithDetail(
			errors.ErrCodeConfigInvalid,
			"配置验证失败",
			strings.Join(errorMessages, "; "),
		)
	}

	return nil
}

// GetValidationErrors 获取所有验证错误
func (cv *ConfigValidator) GetValidationErrors() []ValidationError {
	return cv.errors
}

// GetErrorCount 获取错误数量
func (cv *ConfigValidator) GetErrorCount() int {
	count := 0
	for _, err := range cv.errors {
		if err.Severity == "error" {
			count++
		}
	}
	return count
}

// GetWarningCount 获取警告数量
func (cv *ConfigValidator) GetWarningCount() int {
	count := 0
	for _, err := range cv.errors {
		if err.Severity == "warning" {
			count++
		}
	}
	return count
}

// QuickValidate 快速验证配置（静态方法）
func QuickValidate(config *AppConfig) error {
	validator := NewConfigValidator(false)
	return validator.ValidateConfig(config)
}

// StrictValidate 严格验证配置（静态方法）
func StrictValidate(config *AppConfig) error {
	validator := NewConfigValidator(true)
	return validator.ValidateConfig(config)
}
