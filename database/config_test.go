package database

import (
	"testing"
)

// TestConfig_HasValue 测试数据库配置验证功能
// 这个测试已经在 mysql_test.go 中定义了，这里不重复

// TestConfig_Fields 测试数据库配置字段的访问和修改
func TestConfig_Fields(t *testing.T) {
	// 创建一个测试配置
	cfg := Config{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "password",
		Database: "testdb",
		SSLMode:  true,
	}

	// 测试字段访问
	if cfg.Host != "localhost" {
		t.Errorf("Config.Host = %v, want %v", cfg.Host, "localhost")
	}
	if cfg.Port != 3306 {
		t.Errorf("Config.Port = %v, want %v", cfg.Port, 3306)
	}
	if cfg.Username != "root" {
		t.Errorf("Config.Username = %v, want %v", cfg.Username, "root")
	}
	if cfg.Password != "password" {
		t.Errorf("Config.Password = %v, want %v", cfg.Password, "password")
	}
	if cfg.Database != "testdb" {
		t.Errorf("Config.Database = %v, want %v", cfg.Database, "testdb")
	}
	if !cfg.SSLMode {
		t.Errorf("Config.SSLMode = %v, want %v", cfg.SSLMode, true)
	}

	// 测试字段修改
	cfg.Host = "newhost"
	cfg.Port = 5432
	cfg.Username = "newuser"
	cfg.Password = "newpass"
	cfg.Database = "newdb"
	cfg.SSLMode = false

	// 验证修改后的值
	if cfg.Host != "newhost" {
		t.Errorf("Config.Host = %v, want %v", cfg.Host, "newhost")
	}
	if cfg.Port != 5432 {
		t.Errorf("Config.Port = %v, want %v", cfg.Port, 5432)
	}
	if cfg.Username != "newuser" {
		t.Errorf("Config.Username = %v, want %v", cfg.Username, "newuser")
	}
	if cfg.Password != "newpass" {
		t.Errorf("Config.Password = %v, want %v", cfg.Password, "newpass")
	}
	if cfg.Database != "newdb" {
		t.Errorf("Config.Database = %v, want %v", cfg.Database, "newdb")
	}
	if cfg.SSLMode {
		t.Errorf("Config.SSLMode = %v, want %v", cfg.SSLMode, false)
	}
}

// TestRedisConfig_Fields 测试 Redis 配置字段的访问和修改
func TestRedisConfig_Fields(t *testing.T) {
	// 创建一个测试配置
	cfg := RedisConfig{
		Addr:     "localhost:6379",
		Password: "password",
		DB:       1,
	}

	// 测试字段访问
	if cfg.Addr != "localhost:6379" {
		t.Errorf("RedisConfig.Addr = %v, want %v", cfg.Addr, "localhost:6379")
	}
	if cfg.Password != "password" {
		t.Errorf("RedisConfig.Password = %v, want %v", cfg.Password, "password")
	}
	if cfg.DB != 1 {
		t.Errorf("RedisConfig.DB = %v, want %v", cfg.DB, 1)
	}

	// 测试字段修改
	cfg.Addr = "newhost:6380"
	cfg.Password = "newpass"
	cfg.DB = 2

	// 验证修改后的值
	if cfg.Addr != "newhost:6380" {
		t.Errorf("RedisConfig.Addr = %v, want %v", cfg.Addr, "newhost:6380")
	}
	if cfg.Password != "newpass" {
		t.Errorf("RedisConfig.Password = %v, want %v", cfg.Password, "newpass")
	}
	if cfg.DB != 2 {
		t.Errorf("RedisConfig.DB = %v, want %v", cfg.DB, 2)
	}
}

// TestRedisConfig_HasValue_EdgeCases 测试 Redis 配置验证的边缘情况
func TestRedisConfig_HasValue_EdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		config RedisConfig
		want   bool
	}{
		{
			name: "空地址但有密码",
			config: RedisConfig{
				Addr:     "",
				Password: "password",
				DB:       0,
			},
			want: false,
		},
		{
			name: "只有地址，无密码",
			config: RedisConfig{
				Addr:     "localhost:6379",
				Password: "",
				DB:       0,
			},
			want: true,
		},
		{
			name: "只有地址，有密码",
			config: RedisConfig{
				Addr:     "localhost:6379",
				Password: "password",
				DB:       0,
			},
			want: true,
		},
		{
			name: "只有地址，有密码，非零DB",
			config: RedisConfig{
				Addr:     "localhost:6379",
				Password: "password",
				DB:       1,
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
