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
