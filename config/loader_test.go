package config

import (
	"os"
	"path/filepath"
	"testing"
)

// createTestConfig 创建测试配置文件
func createTestConfig(t *testing.T, content string) string {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// 创建配置文件
	configPath := filepath.Join(tempDir, "config.toml")
	err = os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	return configPath
}

// cleanupTestConfig 清理测试配置文件
func cleanupTestConfig(configPath string) {
	os.RemoveAll(filepath.Dir(configPath))
}

// TestLoadConfig 测试加载配置文件
func TestLoadConfig(t *testing.T) {
	// 创建有效的测试配置
	validConfig := `
[global]
projectname = "test-project"
logLevel = "info"

[metric]
enable = true
port = "8090"
`
	configPath := createTestConfig(t, validConfig)
	defer cleanupTestConfig(configPath)

	// 测试加载配置
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// 验证加载的配置
	if config.Global.ProjectName != "test-project" {
		t.Errorf("LoadConfig() ProjectName = %v, want %v", config.Global.ProjectName, "test-project")
	}
	if config.Global.LogLevel != "info" {
		t.Errorf("LoadConfig() LogLevel = %v, want %v", config.Global.LogLevel, "info")
	}
	if !config.Metric.Enable {
		t.Errorf("LoadConfig() Metric.Enable = %v, want %v", config.Metric.Enable, true)
	}
	if config.Metric.Port != "8090" {
		t.Errorf("LoadConfig() Metric.Port = %v, want %v", config.Metric.Port, "8090")
	}
}

// TestLoadConfigWithValidation 测试加载并验证配置文件
func TestLoadConfigWithValidation(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		expectErr bool
	}{
		{
			name: "有效配置",
			content: `
projectname = "test-project"
logLevel = "info"
`,
			expectErr: false,
		},
		{
			name: "无效日志级别",
			content: `
[global]
projectname = "test-project"
logLevel = "invalid"
`,
			expectErr: true,
		},
		{
			name: "缺少项目名称",
			content: `
[global]
logLevel = "info"
`,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := createTestConfig(t, tt.content)
			defer cleanupTestConfig(configPath)

			_, err := LoadConfigWithValidation(configPath)
			if (err != nil) != tt.expectErr {
				t.Errorf("LoadConfigWithValidation() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestConfigLoader_GetConfig 测试获取配置
func TestConfigLoader_GetConfig(t *testing.T) {
	// 创建测试配置
	configContent := `
projectname = "test-project"
logLevel = "info"
`
	configPath := createTestConfig(t, configContent)
	defer cleanupTestConfig(configPath)

	// 创建加载器
	loader, err := NewConfigLoader(configPath, LoaderOptions{
		EnableWatch:    false,
		ValidateConfig: true,
	})
	if err != nil {
		t.Fatalf("NewConfigLoader() error = %v", err)
	}
	defer loader.Close()

	// 获取配置
	config := loader.GetConfig()
	if config == nil {
		t.Fatal("GetConfig() returned nil")
	}

	// 验证配置
	if config.Global.ProjectName != "test-project" {
		t.Errorf("GetConfig() ProjectName = %v, want %v", config.Global.ProjectName, "test-project")
	}
	if config.Global.LogLevel != "info" {
		t.Errorf("GetConfig() LogLevel = %v, want %v", config.Global.LogLevel, "info")
	}
}

// TestConfigLoader_ReloadConfig 测试重新加载配置
func TestConfigLoader_ReloadConfig(t *testing.T) {
	// 创建初始测试配置
	initialConfig := `
projectname = "initial-project"
logLevel = "info"
`
	configPath := createTestConfig(t, initialConfig)
	defer cleanupTestConfig(configPath)

	// 创建加载器
	loader, err := NewConfigLoader(configPath, LoaderOptions{
		EnableWatch:    false,
		ValidateConfig: true,
	})
	if err != nil {
		t.Fatalf("NewConfigLoader() error = %v", err)
	}
	defer loader.Close()

	// 验证初始配置
	config := loader.GetConfig()
	if config.Global.ProjectName != "initial-project" {
		t.Errorf("Initial ProjectName = %v, want %v", config.Global.ProjectName, "initial-project")
	}

	// 更新配置文件
	updatedConfig := `
projectname = "updated-project"
logLevel = "debug"
`
	err = os.WriteFile(configPath, []byte(updatedConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to update config file: %v", err)
	}

	// 重新加载配置
	err = loader.ReloadConfig()
	if err != nil {
		t.Fatalf("ReloadConfig() error = %v", err)
	}

	// 验证更新后的配置
	updatedCfg := loader.GetConfig()
	if updatedCfg.Global.ProjectName != "updated-project" {
		t.Errorf("Updated ProjectName = %v, want %v", updatedCfg.Global.ProjectName, "updated-project")
	}
	if updatedCfg.Global.LogLevel != "debug" {
		t.Errorf("Updated LogLevel = %v, want %v", updatedCfg.Global.LogLevel, "debug")
	}
}

// TestConfigLoader_AddReloadCallback 测试添加重载回调
func TestConfigLoader_AddReloadCallback(t *testing.T) {
	// 创建测试配置
	configContent := `
projectname = "test-project"
logLevel = "info"
`
	configPath := createTestConfig(t, configContent)
	defer cleanupTestConfig(configPath)

	// 创建加载器
	loader, err := NewConfigLoader(configPath, LoaderOptions{
		EnableWatch:    false,
		ValidateConfig: true,
	})
	if err != nil {
		t.Fatalf("NewConfigLoader() error = %v", err)
	}
	defer loader.Close()

	// 创建回调标志
	callbackCalled := false
	var oldCfg, newCfg *AppConfig

	// 添加回调
	loader.AddReloadCallback(func(old, new *AppConfig) error {
		callbackCalled = true
		oldCfg = old
		newCfg = new
		return nil
	})

	// 更新配置文件
	updatedConfig := `
projectname = "updated-project"
logLevel = "debug"
`
	err = os.WriteFile(configPath, []byte(updatedConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to update config file: %v", err)
	}

	// 重新加载配置
	err = loader.ReloadConfig()
	if err != nil {
		t.Fatalf("ReloadConfig() error = %v", err)
	}

	// 验证回调是否被调用
	if !callbackCalled {
		t.Error("Reload callback was not called")
	}

	// 验证回调参数
	if oldCfg == nil {
		t.Error("Callback old config is nil")
	} else if oldCfg.Global.ProjectName != "test-project" {
		t.Errorf("Callback old ProjectName = %v, want %v", oldCfg.Global.ProjectName, "test-project")
	}

	if newCfg == nil {
		t.Error("Callback new config is nil")
	} else if newCfg.Global.ProjectName != "updated-project" {
		t.Errorf("Callback new ProjectName = %v, want %v", newCfg.Global.ProjectName, "updated-project")
	}
}

// TestConfigLoader_GetConfigInfo 测试获取配置信息
func TestConfigLoader_GetConfigInfo(t *testing.T) {
	// 创建测试配置
	configContent := `
projectname = "test-project"
logLevel = "info"
`
	configPath := createTestConfig(t, configContent)
	defer cleanupTestConfig(configPath)

	// 创建加载器
	loader, err := NewConfigLoader(configPath, LoaderOptions{
		EnableWatch:    false,
		ValidateConfig: true,
	})
	if err != nil {
		t.Fatalf("NewConfigLoader() error = %v", err)
	}
	defer loader.Close()

	// 获取配置信息
	info := loader.GetConfigInfo()

	// 验证信息
	if info["path"] != configPath {
		t.Errorf("GetConfigInfo() path = %v, want %v", info["path"], configPath)
	}
	if info["project_name"] != "test-project" {
		t.Errorf("GetConfigInfo() project_name = %v, want %v", info["project_name"], "test-project")
	}
	if info["log_level"] != "info" {
		t.Errorf("GetConfigInfo() log_level = %v, want %v", info["log_level"], "info")
	}
	if info["version"] != VERSION {
		t.Errorf("GetConfigInfo() version = %v, want %v", info["version"], VERSION)
	}
	if info["is_watching"] != false {
		t.Errorf("GetConfigInfo() is_watching = %v, want %v", info["is_watching"], false)
	}
}

// TestValidateConfigFile 测试验证配置文件
func TestValidateConfigFile(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		expectErr bool
	}{
		{
			name: "有效配置",
			content: `
projectname = "test-project"
logLevel = "info"
`,
			expectErr: false,
		},
		{
			name: "无效配置",
			content: `
projectname = ""
logLevel = "invalid"
`,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := createTestConfig(t, tt.content)
			defer cleanupTestConfig(configPath)

			err := ValidateConfigFile(configPath)
			if (err != nil) != tt.expectErr {
				t.Errorf("ValidateConfigFile() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestConfigLoader_ExportConfig 测试导出配置
func TestConfigLoader_ExportConfig(t *testing.T) {
	// 创建测试配置
	configContent := `
projectname = "test-project"
logLevel = "info"
`
	configPath := createTestConfig(t, configContent)
	defer cleanupTestConfig(configPath)

	// 创建加载器
	loader, err := NewConfigLoader(configPath, LoaderOptions{
		EnableWatch:    false,
		ValidateConfig: true,
	})
	if err != nil {
		t.Fatalf("NewConfigLoader() error = %v", err)
	}
	defer loader.Close()

	// 导出配置
	exportPath := filepath.Join(filepath.Dir(configPath), "exported.toml")
	err = loader.ExportConfig(exportPath)
	if err != nil {
		t.Fatalf("ExportConfig() error = %v", err)
	}

	// 验证导出的配置文件存在
	if _, err := os.Stat(exportPath); os.IsNotExist(err) {
		t.Errorf("Exported config file does not exist: %s", exportPath)
	}

	// 加载导出的配置
	exportedLoader, err := NewConfigLoader(exportPath, LoaderOptions{
		EnableWatch:    false,
		ValidateConfig: true,
	})
	if err != nil {
		t.Fatalf("Failed to load exported config: %v", err)
	}
	defer exportedLoader.Close()

	// 验证导出的配置
	exportedCfg := exportedLoader.GetConfig()
	if exportedCfg.Global.ProjectName != "test-project" {
		t.Errorf("Exported ProjectName = %v, want %v", exportedCfg.Global.ProjectName, "test-project")
	}
	if exportedCfg.Global.LogLevel != "info" {
		t.Errorf("Exported LogLevel = %v, want %v", exportedCfg.Global.LogLevel, "info")
	}
}
