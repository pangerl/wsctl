// Package database MySQL 数据库连接工具
package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL 驱动
	"go.uber.org/zap"
)

// NewMySQLClient 创建 MySQL 数据库连接
// 参数:
//   - cfg: 数据库配置信息
//   - dbName: 数据库名称
//
// 返回:
//   - *sql.DB: 数据库连接对象
//   - error: 错误信息
func NewMySQLClient(cfg Config, dbName string) (*sql.DB, error) {
	// 构建数据库连接字符串，增加连接超时参数（5秒）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, dbName)

	// 根据数据库名称添加特定的连接参数
	if dbName == "wshoto" {
		dsn = dsn + "?interpolateParams=true&timeout=5s"
	} else {
		dsn = dsn + "?timeout=5s"
	}

	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		zap.S().Errorw("MySQL DSN 格式不正确", "err", err, "dsn", dsn)
		return nil, fmt.Errorf("创建 MySQL 连接失败: %w", err)
	}

	// 测试连接是否成功
	err = db.Ping()
	if err != nil {
		zap.S().Errorw("MySQL 数据库连接校验失败", "err", err, "host", cfg.Host, "port", cfg.Port, "database", dbName)
		return nil, fmt.Errorf("MySQL 数据库连接校验失败: %w", err)
	}

	zap.S().Infow("MySQL 数据库连接成功", "host", cfg.Host, "port", cfg.Port, "database", dbName)
	return db, nil
}
