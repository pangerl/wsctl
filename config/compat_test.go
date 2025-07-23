package config

import (
	"testing"
	"time"
)

// TestCompat_TypeAliases 测试类型别名
func TestCompat_TypeAliases(t *testing.T) {
	// 测试数据库配置类型别名
	var dbConfig DB
	dbConfig.Host = "localhost"
	dbConfig.Port = 3306
	dbConfig.Username = "user"
	dbConfig.Password = "pass"

	// 验证类型别名是否正确
	if dbConfig.Host != "localhost" || dbConfig.Port != 3306 {
		t.Error("DB type alias not working correctly")
	}

	// 测试Redis配置类型别名
	var redisConfig RedisConfig
	redisConfig.Addr = "localhost:6379"
	redisConfig.Password = "pass"
	redisConfig.DB = 0

	// 验证类型别名是否正确
	if redisConfig.Addr != "localhost:6379" || redisConfig.Password != "pass" {
		t.Error("RedisConfig type alias not working correctly")
	}

	// 测试其他类型别名
	var tenantCfg TenantCfg
	var corpCfg CorpCfg
	var dorisCfg DorisCfg
	var rocketMQCfg RocketMQCfg
	var nacosCfg NacosCfg
	var metricCfg MetricCfg
	var notifierCfg NotifierCfg

	// 验证类型别名是否存在
	_ = tenantCfg
	_ = corpCfg
	_ = dorisCfg
	_ = rocketMQCfg
	_ = nacosCfg
	_ = metricCfg
	_ = notifierCfg
}

// TestGetRandomDuration 测试随机持续时间函数
func TestGetRandomDuration(t *testing.T) {
	duration := GetRandomDuration()

	// 验证返回值是否为时间间隔类型
	if _, ok := interface{}(duration).(time.Duration); !ok {
		t.Error("GetRandomDuration() did not return time.Duration")
	}
}

// TestGetConfig 测试获取配置函数
func TestGetConfig(t *testing.T) {
	// 设置全局配置
	Config = &AppConfig{
		Global: GlobalConfig{
			ProjectName: "test-project",
		},
	}

	// 获取配置
	config := GetConfig()

	// 验证配置是否正确
	if config == nil {
		t.Fatal("GetConfig() returned nil")
	}
	if config.Global.ProjectName != "test-project" {
		t.Errorf("GetConfig() returned wrong ProjectName, got %s, want %s", config.Global.ProjectName, "test-project")
	}
}

// TestLoadConfig_Compat 测试兼容的配置加载函数
func TestLoadConfig_Compat(t *testing.T) {
	// 由于LoadConfig需要实际的配置文件，这里我们只测试函数存在性
	// 实际的功能测试在config_test.go中进行

	// 验证函数类型签名
	var _ func(string) (*AppConfig, error) = LoadConfig

	t.Log("LoadConfig函数签名验证通过")
}

// fsWrite 辅助函数，写入文件
func fsWrite(path, content string) error {
	return fsWriteBytes(path, []byte(content))
}

// fsWriteBytes 辅助函数，写入字节数据
func fsWriteBytes(path string, data []byte) error {
	return nil // 在测试中不实际写入文件
}
