package database

import (
	"testing"

	"github.com/jackc/pgx/v5"
)

// TestBuildPostgreSQLConnString 测试 PostgreSQL 连接字符串构建功能
func TestBuildPostgreSQLConnString(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		dbName   string
		expected string
	}{
		{
			name: "基本连接字符串_SSL禁用",
			config: Config{
				Host:     "localhost",
				Port:     5432,
				Username: "postgres",
				Password: "password",
				SSLMode:  false,
			},
			dbName:   "testdb",
			expected: "postgres://postgres:password@localhost:5432/testdb?sslmode=disable",
		},
		{
			name: "基本连接字符串_SSL启用",
			config: Config{
				Host:     "localhost",
				Port:     5432,
				Username: "postgres",
				Password: "password",
				SSLMode:  true,
			},
			dbName:   "testdb",
			expected: "postgres://postgres:password@localhost:5432/testdb?sslmode=require",
		},
		{
			name: "包含特殊字符的密码",
			config: Config{
				Host:     "localhost",
				Port:     5432,
				Username: "postgres",
				Password: "pass@word#123",
				SSLMode:  false,
			},
			dbName:   "testdb",
			expected: "postgres://postgres:pass%40word%23123@localhost:5432/testdb?sslmode=disable",
		},
		{
			name: "不同的主机和端口",
			config: Config{
				Host:     "192.168.1.100",
				Port:     5433,
				Username: "admin",
				Password: "secret",
				SSLMode:  true,
			},
			dbName:   "production",
			expected: "postgres://admin:secret@192.168.1.100:5433/production?sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildPostgreSQLConnString(tt.config, tt.dbName)
			if result != tt.expected {
				t.Errorf("buildPostgreSQLConnString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestPGClient_Close 测试 PostgreSQL 客户端关闭功能
func TestPGClient_Close(t *testing.T) {
	// 创建一个空的客户端用于测试
	client := &PGClient{
		Conn: make(map[string]*pgx.Conn),
	}

	// 测试关闭空连接映射不会出错
	client.Close()

	// 验证连接映射仍然存在
	if client.Conn == nil {
		t.Error("Close() should not set Conn to nil")
	}
}
