package config

import (
	"testing"
	"time"
	"vhagar/database"
)

// TestNewConfigValidator 测试创建配置验证器
func TestNewConfigValidator(t *testing.T) {
	validator := NewConfigValidator(false)
	if validator == nil {
		t.Fatal("NewConfigValidator() returned nil")
	}

	if validator.strict {
		t.Error("NewConfigValidator(false) created strict validator")
	}

	strictValidator := NewConfigValidator(true)
	if !strictValidator.strict {
		t.Error("NewConfigValidator(true) did not create strict validator")
	}
}

// TestConfigValidator_AddRule 测试添加验证规则
func TestConfigValidator_AddRule(t *testing.T) {
	validator := NewConfigValidator(false)

	// 创建一个简单的验证规则
	rule := &testValidationRule{
		name:        "test-rule",
		description: "Test validation rule",
	}

	// 添加规则
	validator.AddRule("test-field", rule)

	// 验证规则是否被添加
	if len(validator.rules["test-field"]) != 1 {
		t.Errorf("AddRule() did not add rule, rules count = %d", len(validator.rules["test-field"]))
	}

	// 添加第二个规则
	rule2 := &testValidationRule{
		name:        "test-rule-2",
		description: "Another test validation rule",
	}
	validator.AddRule("test-field", rule2)

	// 验证规则是否被添加
	if len(validator.rules["test-field"]) != 2 {
		t.Errorf("AddRule() did not add second rule, rules count = %d", len(validator.rules["test-field"]))
	}
}

// TestConfigValidator_SetContext 测试设置验证上下文
func TestConfigValidator_SetContext(t *testing.T) {
	validator := NewConfigValidator(false)

	// 设置上下文
	validator.SetContext("test-key", "test-value")

	// 验证上下文是否被设置
	value, exists := validator.context["test-key"]
	if !exists {
		t.Error("SetContext() did not set context key")
	}
	if value != "test-value" {
		t.Errorf("SetContext() set wrong value, got %v, want %v", value, "test-value")
	}
}

// TestQuickValidate 测试快速验证配置
func TestQuickValidate(t *testing.T) {
	tests := []struct {
		name      string
		config    *AppConfig
		expectErr bool
	}{
		{
			name: "有效配置",
			config: &AppConfig{
				Global: GlobalConfig{
					LogLevel:    "info",
					ProjectName: "test-project",
				},
			},
			expectErr: false,
		},
		{
			name: "无效配置",
			config: &AppConfig{
				Global: GlobalConfig{
					LogLevel:    "invalid",
					ProjectName: "",
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := QuickValidate(tt.config)
			if (err != nil) != tt.expectErr {
				t.Errorf("QuickValidate() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestStrictValidate 测试严格验证配置
func TestStrictValidate(t *testing.T) {
	// 创建一个带警告的配置
	config := &AppConfig{
		Global: GlobalConfig{
			LogLevel:    "info",
			ProjectName: "test-project",
			Notify: NotifyConfig{
				Robotkey: []string{"invalid-uuid"}, // 不是有效的UUID格式，会产生警告
			},
		},
	}

	// 非严格模式下，警告不会导致错误
	err := QuickValidate(config)
	if err != nil {
		t.Errorf("QuickValidate() error = %v, want nil", err)
	}

	// 严格模式下，警告会导致错误
	err = StrictValidate(config)
	if err == nil {
		t.Error("StrictValidate() returned nil, want error")
	}
}

// TestConfigValidator_validateRequired 测试必填项验证
func TestConfigValidator_validateRequired(t *testing.T) {
	validator := NewConfigValidator(false)

	// 测试非空值
	validator.validateRequired("test-field", "non-empty")
	if len(validator.errors) != 0 {
		t.Errorf("validateRequired() added error for non-empty value, errors = %v", validator.errors)
	}

	// 测试空值
	validator.validateRequired("test-field", "")
	if len(validator.errors) != 1 {
		t.Errorf("validateRequired() did not add error for empty value, errors count = %d", len(validator.errors))
	}

	// 测试空白值
	validator.validateRequired("test-field-2", "   ")
	if len(validator.errors) != 2 {
		t.Errorf("validateRequired() did not add error for whitespace value, errors count = %d", len(validator.errors))
	}
}

// TestConfigValidator_validateStringLength 测试字符串长度验证
func TestConfigValidator_validateStringLength(t *testing.T) {
	validator := NewConfigValidator(false)

	// 测试有效长度
	validator.validateStringLength("test-field", "valid", 1, 10)
	if len(validator.errors) != 0 {
		t.Errorf("validateStringLength() added error for valid length, errors = %v", validator.errors)
	}

	// 测试过短
	validator.validateStringLength("test-field", "a", 2, 10)
	if len(validator.errors) != 1 {
		t.Errorf("validateStringLength() did not add error for too short value, errors count = %d", len(validator.errors))
	}

	// 测试过长
	validator.validateStringLength("test-field", "too long string", 1, 5)
	if len(validator.errors) != 2 {
		t.Errorf("validateStringLength() did not add error for too long value, errors count = %d", len(validator.errors))
	}
}

// TestConfigValidator_validateEnum 测试枚举值验证
func TestConfigValidator_validateEnum(t *testing.T) {
	validator := NewConfigValidator(false)
	validValues := []string{"a", "b", "c"}

	// 测试有效值
	validator.validateEnum("test-field", "a", validValues)
	if len(validator.errors) != 0 {
		t.Errorf("validateEnum() added error for valid value, errors = %v", validator.errors)
	}

	// 测试无效值
	validator.validateEnum("test-field", "d", validValues)
	if len(validator.errors) != 1 {
		t.Errorf("validateEnum() did not add error for invalid value, errors count = %d", len(validator.errors))
	}
}

// TestConfigValidator_validateURL 测试URL格式验证
func TestConfigValidator_validateURL(t *testing.T) {
	validator := NewConfigValidator(false)

	// 测试有效URL
	validator.validateURL("test-field", "http://example.com")
	if len(validator.errors) != 0 {
		t.Errorf("validateURL() added error for valid URL, errors = %v", validator.errors)
	}

	// 测试无效URL
	validator.validateURL("test-field", "invalid-url")
	if len(validator.errors) != 1 {
		t.Errorf("validateURL() did not add error for invalid URL, errors count = %d", len(validator.errors))
	}

	// 测试空URL
	validator.validateURL("test-field", "")
	if len(validator.errors) != 1 {
		t.Errorf("validateURL() added error for empty URL, errors count = %d", len(validator.errors))
	}
}

// TestConfigValidator_validateIP 测试IP地址验证
func TestConfigValidator_validateIP(t *testing.T) {
	validator := NewConfigValidator(false)

	// 测试有效IP
	validator.validateIP("test-field", "192.168.1.1")
	if len(validator.errors) != 0 {
		t.Errorf("validateIP() added error for valid IP, errors = %v", validator.errors)
	}

	// 测试无效IP
	validator.validateIP("test-field", "invalid-ip")
	if len(validator.errors) != 1 {
		t.Errorf("validateIP() did not add error for invalid IP, errors count = %d", len(validator.errors))
	}

	// 测试空IP
	validator.validateIP("test-field", "")
	if len(validator.errors) != 1 {
		t.Errorf("validateIP() added error for empty IP, errors count = %d", len(validator.errors))
	}
}

// TestConfigValidator_validatePort 测试端口号验证
func TestConfigValidator_validatePort(t *testing.T) {
	validator := NewConfigValidator(false)

	// 测试有效端口
	validator.validatePort("test-field", 8080)
	if len(validator.errors) != 0 {
		t.Errorf("validatePort() added error for valid port, errors = %v", validator.errors)
	}

	// 测试过小的端口
	validator.validatePort("test-field", 0)
	if len(validator.errors) != 1 {
		t.Errorf("validatePort() did not add error for too small port, errors count = %d", len(validator.errors))
	}

	// 测试过大的端口
	validator.validatePort("test-field", 70000)
	if len(validator.errors) != 2 {
		t.Errorf("validatePort() did not add error for too large port, errors count = %d", len(validator.errors))
	}
}

// TestConfigValidator_validateDuration 测试时间间隔验证
func TestConfigValidator_validateDuration(t *testing.T) {
	validator := NewConfigValidator(false)

	// 测试有效时间间隔
	validator.validateDuration("test-field", 5*time.Minute, time.Minute, 10*time.Minute)
	if len(validator.errors) != 0 {
		t.Errorf("validateDuration() added error for valid duration, errors = %v", validator.errors)
	}

	// 测试过小的时间间隔
	validator.validateDuration("test-field", 30*time.Second, time.Minute, 10*time.Minute)
	if len(validator.errors) != 1 {
		t.Errorf("validateDuration() did not add error for too small duration, errors count = %d", len(validator.errors))
	}

	// 测试过大的时间间隔
	validator.validateDuration("test-field", 20*time.Minute, time.Minute, 10*time.Minute)
	if len(validator.errors) != 2 {
		t.Errorf("validateDuration() did not add error for too large duration, errors count = %d", len(validator.errors))
	}
}

// TestConfigValidator_validateRange 测试数值范围验证
func TestConfigValidator_validateRange(t *testing.T) {
	validator := NewConfigValidator(false)

	// 测试有效范围
	validator.validateRange("test-field", 5, 1, 10)
	if len(validator.errors) != 0 {
		t.Errorf("validateRange() added error for valid range, errors = %v", validator.errors)
	}

	// 测试过小的值
	validator.validateRange("test-field", 0, 1, 10)
	if len(validator.errors) != 1 {
		t.Errorf("validateRange() did not add error for too small value, errors count = %d", len(validator.errors))
	}

	// 测试过大的值
	validator.validateRange("test-field", 20, 1, 10)
	if len(validator.errors) != 2 {
		t.Errorf("validateRange() did not add error for too large value, errors count = %d", len(validator.errors))
	}
}

// TestConfigValidator_validateRedisAddr 测试Redis地址验证
func TestConfigValidator_validateRedisAddr(t *testing.T) {
	validator := NewConfigValidator(false)

	// 测试有效地址
	validator.validateRedisAddr("test-field", "localhost:6379")
	if len(validator.errors) != 0 {
		t.Errorf("validateRedisAddr() added error for valid address, errors = %v", validator.errors)
	}

	// 测试无效格式
	validator.validateRedisAddr("test-field", "invalid-format")
	if len(validator.errors) != 1 {
		t.Errorf("validateRedisAddr() did not add error for invalid format, errors count = %d", len(validator.errors))
	}

	// 测试空主机
	validator.validateRedisAddr("test-field", ":6379")
	if len(validator.errors) != 2 {
		t.Errorf("validateRedisAddr() did not add error for empty host, errors count = %d", len(validator.errors))
	}

	// 测试无效端口
	validator.validateRedisAddr("test-field", "localhost:invalid")
	if len(validator.errors) != 3 {
		t.Errorf("validateRedisAddr() did not add error for invalid port, errors count = %d", len(validator.errors))
	}

	// 测试空地址
	validator.validateRedisAddr("test-field", "")
	if len(validator.errors) != 3 {
		t.Errorf("validateRedisAddr() added error for empty address, errors count = %d", len(validator.errors))
	}
}

// TestConfigValidator_validateCronExpression 测试Cron表达式验证
func TestConfigValidator_validateCronExpression(t *testing.T) {
	validator := NewConfigValidator(false)

	// 测试有效表达式
	validator.validateCronExpression("test-field", "0 * * * *")
	if len(validator.errors) != 0 {
		t.Errorf("validateCronExpression() added error for valid expression, errors = %v", validator.errors)
	}

	// 测试无效表达式
	validator.validateCronExpression("test-field", "invalid")
	if len(validator.errors) != 1 {
		t.Errorf("validateCronExpression() did not add error for invalid expression, errors count = %d", len(validator.errors))
	}

	// 测试空表达式
	validator.validateCronExpression("test-field", "")
	if len(validator.errors) != 1 {
		t.Errorf("validateCronExpression() added error for empty expression, errors count = %d", len(validator.errors))
	}
}

// TestConfigValidator_validateUUID 测试UUID格式验证
func TestConfigValidator_validateUUID(t *testing.T) {
	validator := NewConfigValidator(false)

	// 测试有效UUID
	validator.validateUUID("test-field", "123e4567-e89b-12d3-a456-426614174000")
	if len(validator.errors) != 0 {
		t.Errorf("validateUUID() added error for valid UUID, errors = %v", validator.errors)
	}

	// 测试无效UUID
	validator.validateUUID("test-field", "invalid-uuid")
	if len(validator.errors) != 1 {
		t.Errorf("validateUUID() did not add error for invalid UUID, errors count = %d", len(validator.errors))
	}

	// 测试空UUID
	validator.validateUUID("test-field", "")
	if len(validator.errors) != 1 {
		t.Errorf("validateUUID() added error for empty UUID, errors count = %d", len(validator.errors))
	}
}

// TestConfigValidator_validateGlobalConfig 测试全局配置验证
func TestConfigValidator_validateGlobalConfig(t *testing.T) {
	validator := NewConfigValidator(false)

	// 测试有效配置
	validConfig := GlobalConfig{
		LogLevel:    "info",
		ProjectName: "test-project",
		Interval:    5 * time.Minute,
		Duration:    time.Hour,
	}
	validator.validateGlobalConfig(&validConfig)
	if len(validator.errors) != 0 {
		t.Errorf("validateGlobalConfig() added error for valid config, errors = %v", validator.errors)
	}

	// 测试无效配置
	invalidConfig := GlobalConfig{
		LogLevel:    "invalid",
		ProjectName: "",
		Interval:    0,
		Duration:    0,
	}
	validator.errors = nil // 清空错误
	validator.validateGlobalConfig(&invalidConfig)
	if len(validator.errors) == 0 {
		t.Error("validateGlobalConfig() did not add errors for invalid config")
	}
}

// TestConfigValidator_validateDatabaseConfigs 测试数据库配置验证
func TestConfigValidator_validateDatabaseConfigs(t *testing.T) {
	validator := NewConfigValidator(false)

	// 测试有效配置
	validConfig := DatabaseConfigs{
		PG: database.Config{
			Host:     "localhost",
			Port:     5432,
			Username: "postgres",
			Password: "password",
		},
		Redis: database.RedisConfig{
			Addr: "localhost:6379",
		},
	}
	validator.validateDatabaseConfigs(&validConfig)
	if len(validator.errors) != 0 {
		t.Errorf("validateDatabaseConfigs() added error for valid config, errors = %v", validator.errors)
	}

	// 测试无效配置
	invalidConfig := DatabaseConfigs{
		PG: database.Config{
			Host:     "",
			Port:     0,
			Username: "",
			Password: "",
		},
		Redis: database.RedisConfig{
			Addr: "invalid",
		},
	}
	validator.errors = nil // 清空错误
	validator.validateDatabaseConfigs(&invalidConfig)
	if len(validator.errors) == 0 {
		t.Error("validateDatabaseConfigs() did not add errors for invalid config")
	}
}

// TestConfigValidator_GetErrorCount 测试获取错误数量
func TestConfigValidator_GetErrorCount(t *testing.T) {
	validator := NewConfigValidator(false)

	// 添加错误
	validator.addError("field1", "value1", "rule1", "message1", "error", "suggestion1")
	validator.addError("field2", "value2", "rule2", "message2", "warning", "suggestion2")
	validator.addError("field3", "value3", "rule3", "message3", "error", "suggestion3")

	// 验证错误数量
	errorCount := validator.GetErrorCount()
	if errorCount != 2 {
		t.Errorf("GetErrorCount() = %d, want %d", errorCount, 2)
	}
}

// TestConfigValidator_GetWarningCount 测试获取警告数量
func TestConfigValidator_GetWarningCount(t *testing.T) {
	validator := NewConfigValidator(false)

	// 添加警告
	validator.addError("field1", "value1", "rule1", "message1", "warning", "suggestion1")
	validator.addError("field2", "value2", "rule2", "message2", "error", "suggestion2")
	validator.addError("field3", "value3", "rule3", "message3", "warning", "suggestion3")

	// 验证警告数量
	warningCount := validator.GetWarningCount()
	if warningCount != 2 {
		t.Errorf("GetWarningCount() = %d, want %d", warningCount, 2)
	}
}

// TestConfigValidator_GetValidationErrors 测试获取验证错误
func TestConfigValidator_GetValidationErrors(t *testing.T) {
	validator := NewConfigValidator(false)

	// 添加错误
	validator.addError("field1", "value1", "rule1", "message1", "error", "suggestion1")
	validator.addError("field2", "value2", "rule2", "message2", "warning", "suggestion2")

	// 获取错误
	errors := validator.GetValidationErrors()
	if len(errors) != 2 {
		t.Errorf("GetValidationErrors() returned %d errors, want %d", len(errors), 2)
	}

	// 验证错误内容
	if errors[0].Field != "field1" || errors[0].Severity != "error" {
		t.Errorf("First error has wrong content: %+v", errors[0])
	}
	if errors[1].Field != "field2" || errors[1].Severity != "warning" {
		t.Errorf("Second error has wrong content: %+v", errors[1])
	}
}

// TestValidationError_Error 测试验证错误的Error方法
func TestValidationError_Error(t *testing.T) {
	err := ValidationError{
		Field:    "test-field",
		Value:    "test-value",
		Rule:     "test-rule",
		Message:  "Test message",
		Severity: "error",
	}

	expected := "[error] test-field: Test message (值: test-value)"
	if err.Error() != expected {
		t.Errorf("ValidationError.Error() = %q, want %q", err.Error(), expected)
	}
}

// 测试用的验证规则实现
type testValidationRule struct {
	name        string
	description string
}

func (r *testValidationRule) Validate(value interface{}) error {
	return nil
}

func (r *testValidationRule) GetName() string {
	return r.name
}

func (r *testValidationRule) GetDescription() string {
	return r.description
}
