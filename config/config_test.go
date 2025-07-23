package config

import (
	"testing"
	"time"
	"vhagar/database"
)

// TestDefaultConfig 测试默认配置
func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Global.LogLevel != "info" {
		t.Errorf("Expected default LogLevel 'info', got '%s'", config.Global.LogLevel)
	}

	if config.Global.LogToFile != false {
		t.Errorf("Expected default LogToFile false, got %v", config.Global.LogToFile)
	}

	if config.Global.ProjectName != "vhagar" {
		t.Errorf("Expected default ProjectName 'vhagar', got '%s'", config.Global.ProjectName)
	}

	if config.Global.Watch != true {
		t.Errorf("Expected default Watch true, got %v", config.Global.Watch)
	}

	if config.Global.Report != true {
		t.Errorf("Expected default Report true, got %v", config.Global.Report)
	}

	if config.Global.Interval != 5*time.Minute {
		t.Errorf("Expected default Interval 5 minutes, got %v", config.Global.Interval)
	}

	if config.Global.Duration != time.Hour {
		t.Errorf("Expected default Duration 1 hour, got %v", config.Global.Duration)
	}

	if config.Metric.Enable != false {
		t.Errorf("Expected default Metric.Enable false, got %v", config.Metric.Enable)
	}

	if config.Metric.Port != "8090" {
		t.Errorf("Expected default Metric.Port '8090', got '%s'", config.Metric.Port)
	}
}

// TestAppConfig_Validate 测试配置验证
func TestAppConfig_Validate(t *testing.T) {
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
			name: "无效日志级别",
			config: &AppConfig{
				Global: GlobalConfig{
					LogLevel:    "invalid",
					ProjectName: "test-project",
				},
			},
			expectErr: true,
		},
		{
			name: "缺少项目名称",
			config: &AppConfig{
				Global: GlobalConfig{
					LogLevel:    "info",
					ProjectName: "",
				},
			},
			expectErr: true,
		},
		{
			name: "启用AI但缺少提供商",
			config: &AppConfig{
				Global: GlobalConfig{
					LogLevel:    "info",
					ProjectName: "test-project",
				},
				Services: ServiceConfigs{
					AI: AIConfig{
						Enable:   true,
						Provider: "",
					},
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.expectErr {
				t.Errorf("AppConfig.Validate() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestAppConfig_GetDatabaseConfig 测试获取数据库配置
func TestAppConfig_GetDatabaseConfig(t *testing.T) {
	// 创建测试配置
	config := &AppConfig{
		Database: DatabaseConfigs{
			PG: database.Config{
				Host:     "pg-host",
				Port:     5432,
				Username: "pg-user",
				Password: "pg-pass",
				Database: "pg-db",
			},
			ES: database.Config{
				Host:     "es-host",
				Port:     9200,
				Username: "es-user",
				Password: "es-pass",
			},
			Redis: database.RedisConfig{
				Addr:     "redis-host:6379",
				Password: "redis-pass",
				DB:       0,
			},
		},
	}

	tests := []struct {
		name    string
		dbType  string
		wantNil bool
	}{
		{"PostgreSQL", "pg", false},
		{"PostgreSQL别名", "postgresql", false},
		{"Elasticsearch", "es", false},
		{"Elasticsearch别名", "elasticsearch", false},
		{"Redis", "redis", false},
		{"不存在的类型", "unknown", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.GetDatabaseConfig(tt.dbType)
			if (result == nil) != tt.wantNil {
				t.Errorf("GetDatabaseConfig(%s) = %v, wantNil %v", tt.dbType, result, tt.wantNil)
			}
		})
	}
}

// TestAppConfig_IsAIEnabled 测试AI功能启用状态
func TestAppConfig_IsAIEnabled(t *testing.T) {
	tests := []struct {
		name     string
		aiConfig AIConfig
		want     bool
	}{
		{
			name: "AI启用",
			aiConfig: AIConfig{
				Enable: true,
			},
			want: true,
		},
		{
			name: "AI禁用",
			aiConfig: AIConfig{
				Enable: false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &AppConfig{
				Services: ServiceConfigs{
					AI: tt.aiConfig,
				},
			}
			if got := config.IsAIEnabled(); got != tt.want {
				t.Errorf("IsAIEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestAppConfig_GetAIProvider 测试获取AI提供商配置
func TestAppConfig_GetAIProvider(t *testing.T) {
	tests := []struct {
		name      string
		aiConfig  AIConfig
		wantFound bool
	}{
		{
			name: "有效提供商",
			aiConfig: AIConfig{
				Enable:   true,
				Provider: "test-provider",
				Providers: map[string]ProviderConfig{
					"test-provider": {
						ApiKey: "test-key",
						ApiUrl: "https://test.com",
						Model:  "test-model",
					},
				},
			},
			wantFound: true,
		},
		{
			name: "无效提供商",
			aiConfig: AIConfig{
				Enable:   true,
				Provider: "invalid-provider",
				Providers: map[string]ProviderConfig{
					"test-provider": {
						ApiKey: "test-key",
						ApiUrl: "https://test.com",
						Model:  "test-model",
					},
				},
			},
			wantFound: false,
		},
		{
			name: "AI禁用",
			aiConfig: AIConfig{
				Enable:   false,
				Provider: "test-provider",
				Providers: map[string]ProviderConfig{
					"test-provider": {
						ApiKey: "test-key",
						ApiUrl: "https://test.com",
						Model:  "test-model",
					},
				},
			},
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &AppConfig{
				Services: ServiceConfigs{
					AI: tt.aiConfig,
				},
			}
			_, found := config.GetAIProvider()
			if found != tt.wantFound {
				t.Errorf("GetAIProvider() found = %v, wantFound %v", found, tt.wantFound)
			}
		})
	}
}

// TestAppConfig_IsMetricEnabled 测试监控功能启用状态
func TestAppConfig_IsMetricEnabled(t *testing.T) {
	tests := []struct {
		name         string
		metricConfig MetricConfig
		want         bool
	}{
		{
			name: "监控启用",
			metricConfig: MetricConfig{
				Enable: true,
			},
			want: true,
		},
		{
			name: "监控禁用",
			metricConfig: MetricConfig{
				Enable: false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &AppConfig{
				Metric: tt.metricConfig,
			}
			if got := config.IsMetricEnabled(); got != tt.want {
				t.Errorf("IsMetricEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestAppConfig_GetNotifierConfig 测试获取通知器配置
func TestAppConfig_GetNotifierConfig(t *testing.T) {
	// 创建测试配置
	config := &AppConfig{
		Global: GlobalConfig{
			Notify: NotifyConfig{
				Robotkey: []string{"default-key"},
				Notifier: map[string]NotifierConfig{
					"test-notifier": {
						Robotkey: []string{"test-key"},
					},
				},
			},
		},
	}

	tests := []struct {
		name      string
		notifier  string
		wantFound bool
	}{
		{"存在的通知器", "test-notifier", true},
		{"不存在的通知器", "invalid-notifier", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notifier, found := config.GetNotifierConfig(tt.notifier)
			if found != tt.wantFound {
				t.Errorf("GetNotifierConfig(%s) found = %v, wantFound %v", tt.notifier, found, tt.wantFound)
			}

			if found && len(notifier.Robotkey) == 0 {
				t.Errorf("GetNotifierConfig(%s) returned empty Robotkey", tt.notifier)
			}
		})
	}
}

// TestAppConfig_GetCronConfig 测试获取定时任务配置
func TestAppConfig_GetCronConfig(t *testing.T) {
	// 创建测试配置
	config := &AppConfig{
		Cron: map[string]CrontabConfig{
			"test-cron": {
				Crontab:    true,
				Scheducron: "0 * * * *",
			},
		},
	}

	tests := []struct {
		name      string
		cronName  string
		wantFound bool
	}{
		{"存在的定时任务", "test-cron", true},
		{"不存在的定时任务", "invalid-cron", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cron, found := config.GetCronConfig(tt.cronName)
			if found != tt.wantFound {
				t.Errorf("GetCronConfig(%s) found = %v, wantFound %v", tt.cronName, found, tt.wantFound)
			}

			if found && !cron.Crontab {
				t.Errorf("GetCronConfig(%s) returned disabled cron", tt.cronName)
			}
		})
	}
}
