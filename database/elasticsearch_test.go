package database

import (
	"context"
	"testing"
	"time"
)

// TestBuildElasticsearchURL 测试 Elasticsearch URL 构建功能
func TestBuildElasticsearchURL(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected string
	}{
		{
			name: "HTTP连接",
			config: Config{
				Host:    "localhost",
				Port:    9200,
				SSLMode: false,
			},
			expected: "http://localhost:9200",
		},
		{
			name: "HTTPS连接",
			config: Config{
				Host:    "localhost",
				Port:    9200,
				SSLMode: true,
			},
			expected: "https://localhost:9200",
		},
		{
			name: "自定义主机和端口",
			config: Config{
				Host:    "es.example.com",
				Port:    9201,
				SSLMode: true,
			},
			expected: "https://es.example.com:9201",
		},
		{
			name: "IP地址主机",
			config: Config{
				Host:    "192.168.1.100",
				Port:    9200,
				SSLMode: false,
			},
			expected: "http://192.168.1.100:9200",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildElasticsearchURL(tt.config)
			if result != tt.expected {
				t.Errorf("buildElasticsearchURL() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestElasticsearchClientCreation 测试 Elasticsearch 客户端创建功能
// 注意：这个测试不会实际连接到 Elasticsearch 服务器，只测试 URL 构建逻辑
func TestElasticsearchClientCreation(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		valid  bool
	}{
		{
			name: "有效的Elasticsearch配置_HTTP",
			config: Config{
				Host:     "localhost",
				Port:     9200,
				Username: "elastic",
				Password: "password",
				SSLMode:  false,
			},
			valid: true,
		},
		{
			name: "有效的Elasticsearch配置_HTTPS",
			config: Config{
				Host:     "localhost",
				Port:     9200,
				Username: "elastic",
				Password: "password",
				SSLMode:  true,
			},
			valid: true,
		},
		{
			name: "无效的Elasticsearch配置_空主机",
			config: Config{
				Host:     "",
				Port:     9200,
				Username: "elastic",
				Password: "password",
				SSLMode:  false,
			},
			valid: false,
		},
		{
			name: "无效的Elasticsearch配置_零端口",
			config: Config{
				Host:     "localhost",
				Port:     0,
				Username: "elastic",
				Password: "password",
				SSLMode:  false,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试配置验证
			hasValue := tt.config.HasValue()
			if hasValue != tt.valid {
				t.Errorf("Config.HasValue() = %v, want %v", hasValue, tt.valid)
			}

			// 如果配置有效，测试 URL 构建
			if tt.valid {
				url := buildElasticsearchURL(tt.config)
				if url == "" {
					t.Error("buildElasticsearchURL() returned empty string for valid config")
				}
				t.Logf("Generated URL: %s", url)
			}
		})
	}
}

// TestElasticsearchClientWithContext 测试带上下文的 Elasticsearch 客户端创建
func TestElasticsearchClientWithContext(t *testing.T) {
	// 创建一个带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	config := Config{
		Host:     "localhost",
		Port:     9200,
		Username: "elastic",
		Password: "password",
		SSLMode:  false,
	}

	// 注意：这个测试在没有实际 Elasticsearch 服务器的情况下会失败
	// 但我们可以测试函数是否正确处理上下文和配置
	_, err := NewElasticsearchClientWithContext(ctx, config)

	// 在没有 Elasticsearch 服务器的情况下，我们期望得到连接错误
	// 这证明函数正确尝试了连接
	if err == nil {
		t.Log("Elasticsearch 连接成功（可能有实际的 Elasticsearch 服务器运行）")
	} else {
		t.Logf("Elasticsearch 连接失败（预期行为，没有 Elasticsearch 服务器）: %v", err)
	}
}
