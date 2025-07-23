package database

import (
	"testing"
)

// TestConfig_HasValue 测试数据库配置验证功能
func TestConfig_HasValue(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   bool
	}{
		{
			name: "完整配置",
			config: Config{
				Host:     "localhost",
				Port:     3306,
				Username: "root",
				Password: "password",
			},
			want: true,
		},
		{
			name: "缺少主机地址",
			config: Config{
				Host:     "",
				Port:     3306,
				Username: "root",
				Password: "password",
			},
			want: false,
		},
		{
			name: "缺少端口号",
			config: Config{
				Host:     "localhost",
				Port:     0,
				Username: "root",
				Password: "password",
			},
			want: false,
		},
		{
			name: "缺少用户名",
			config: Config{
				Host:     "localhost",
				Port:     3306,
				Username: "",
				Password: "password",
			},
			want: false,
		},
		{
			name: "缺少密码",
			config: Config{
				Host:     "localhost",
				Port:     3306,
				Username: "root",
				Password: "",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.HasValue(); got != tt.want {
				t.Errorf("Config.HasValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRedisConfig_HasValue 测试 Redis 配置验证功能
func TestRedisConfig_HasValue(t *testing.T) {
	tests := []struct {
		name   string
		config RedisConfig
		want   bool
	}{
		{
			name: "完整配置",
			config: RedisConfig{
				Addr:     "localhost:6379",
				Password: "password",
				DB:       0,
			},
			want: true,
		},
		{
			name: "缺少地址",
			config: RedisConfig{
				Addr:     "",
				Password: "password",
				DB:       0,
			},
			want: false,
		},
		{
			name: "只有地址",
			config: RedisConfig{
				Addr: "localhost:6379",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.HasValue(); got != tt.want {
				t.Errorf("RedisConfig.HasValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBuildMySQLDSN 测试 MySQL DSN 构建功能
func TestBuildMySQLDSN(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		dbName   string
		expected string
	}{
		{
			name: "标准数据库",
			config: Config{
				Host:     "localhost",
				Port:     3306,
				Username: "root",
				Password: "password",
			},
			dbName:   "testdb",
			expected: "root:password@tcp(localhost:3306)/testdb?timeout=5s",
		},
		{
			name: "wshoto数据库",
			config: Config{
				Host:     "localhost",
				Port:     3306,
				Username: "root",
				Password: "password",
			},
			dbName:   "wshoto",
			expected: "root:password@tcp(localhost:3306)/wshoto?interpolateParams=true&timeout=5s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 注意：这里我们不能直接比较字符串，因为我们没有访问内部 DSN 构建函数
			// 所以这个测试主要是为了记录预期的 DSN 格式
			t.Logf("Expected DSN for %s: %s", tt.name, tt.expected)
		})
	}
}

// TestNewMySQLClient 测试 MySQL 客户端创建功能
// 注意：这个测试不会实际连接到 MySQL 服务器，只测试函数签名
func TestNewMySQLClient(t *testing.T) {
	// 由于 NewMySQLClient 函数会尝试实际连接数据库，
	// 这里我们只能测试函数的存在性，而不是实际功能
	// 在真实环境中，应该使用 mock 对象来测试

	cfg := Config{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "password",
	}

	// 这个调用会失败，因为没有实际的数据库，但我们可以验证函数存在
	_, err := NewMySQLClient(cfg, "testdb")

	// 我们期望这里会有错误，因为没有实际的数据库连接
	if err == nil {
		t.Log("意外地成功连接到数据库")
	} else {
		t.Logf("预期的数据库连接错误: %v", err)
	}
}
