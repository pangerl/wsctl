// Package database PostgreSQL 数据库连接工具
package database

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// PGClient PostgreSQL 客户端结构体
// 管理多个 PostgreSQL 数据库连接
type PGClient struct {
	Conn map[string]*pgx.Conn // 数据库连接映射，key 为数据库名称，value 为连接对象
}

// Close 关闭所有 PostgreSQL 数据库连接
// 遍历所有连接并安全关闭，记录关闭失败的错误
func (client *PGClient) Close() {
	for dbName, conn := range client.Conn {
		if err := conn.Close(context.Background()); err != nil {
			zap.S().Errorw("关闭 PostgreSQL 数据库连接失败", "database", dbName, "err", err)
		} else {
			zap.S().Infow("PostgreSQL 数据库连接已关闭", "database", dbName)
		}
	}
}

// NewPGClient 创建新的 PostgreSQL 客户端
// 自动连接到预定义的数据库列表：qv30, user, customer
// 参数:
//   - cfg: 数据库配置信息
//
// 返回:
//   - *PGClient: PostgreSQL 客户端对象
//   - error: 错误信息
func NewPGClient(cfg Config) (*PGClient, error) {
	client := &PGClient{
		Conn: make(map[string]*pgx.Conn),
	}

	// 预定义的数据库列表
	databases := []string{"qv30", "user", "customer"}

	// 连接数据库的内部函数，支持 user -> users 的回退机制
	connectDB := func(dbName string) (*pgx.Conn, error) {
		conn, err := NewPostgreSQLConnection(cfg, dbName)
		if err != nil {
			if dbName == "user" {
				// 如果连接 'user' 数据库失败，尝试连接 'users' 数据库
				zap.S().Warnw("连接 'user' 数据库失败，尝试连接 'users' 数据库", "err", err)
				conn, err = NewPostgreSQLConnection(cfg, "users")
			}
		}
		return conn, err
	}

	// 遍历数据库列表并建立连接
	for _, dbName := range databases {
		conn, err := connectDB(dbName)
		if err != nil {
			// 如果任何一个数据库连接失败，关闭已建立的连接并返回错误
			client.Close()
			return nil, fmt.Errorf("连接数据库 %s 失败: %w", dbName, err)
		}
		client.Conn[dbName] = conn
	}

	zap.S().Infow("PostgreSQL 客户端创建成功", "databases", databases)
	return client, nil
}

// NewPostgreSQLConnection 创建单个 PostgreSQL 数据库连接
// 参数:
//   - cfg: 数据库配置信息
//   - dbName: 数据库名称
//
// 返回:
//   - *pgx.Conn: PostgreSQL 连接对象
//   - error: 错误信息
func NewPostgreSQLConnection(cfg Config, dbName string) (*pgx.Conn, error) {
	// 构建连接字符串
	connString := buildPostgreSQLConnString(cfg, dbName)

	// 建立数据库连接
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		zap.S().Errorw("PostgreSQL 数据库连接失败",
			"database", dbName,
			"host", cfg.Host,
			"port", cfg.Port,
			"err", err)
		return nil, fmt.Errorf("PostgreSQL 连接失败: %w", err)
	}

	zap.S().Infow("PostgreSQL 数据库连接成功",
		"database", dbName,
		"host", cfg.Host,
		"port", cfg.Port)
	return conn, nil
}

// buildPostgreSQLConnString 构建 PostgreSQL 连接字符串
// 参数:
//   - cfg: 数据库配置信息
//   - dbName: 数据库名称
//
// 返回:
//   - string: 格式化的连接字符串
func buildPostgreSQLConnString(cfg Config, dbName string) string {
	// 根据 SSL 模式设置连接参数
	sslMode := map[bool]string{true: "require", false: "disable"}[cfg.SSLMode]

	// 对密码进行 URL 编码以处理特殊字符
	encodedPassword := url.QueryEscape(cfg.Password)

	// 构建 PostgreSQL 连接字符串
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username, encodedPassword, cfg.Host, cfg.Port, dbName, sslMode)

	return connString
}
