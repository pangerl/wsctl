// Package libs @Author lanpang
// @Date 2024/8/8 下午1:43:00
// @Desc
package libs

import (
	"context"
	"fmt"
	"net/url"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v5"
)

// Close 关闭所有数据库连接
func (dbClient *PGClienter) Close() {
	for dbName, conn := range dbClient.Conn {
		if err := conn.Close(context.Background()); err != nil {
			zap.S().Errorw("关闭数据库连接失败", "db", dbName, "err", err)
		}
	}
}

func NewPGClienter(conf DB) (*PGClienter, error) {
	clienter := &PGClienter{
		Conn: make(map[string]*pgx.Conn),
	}
	databases := []string{"qv30", "user", "customer"}

	connectDB := func(dbName string) (*pgx.Conn, error) {
		conn, err := NewPGClient(conf, dbName)
		if err != nil {
			if dbName == "user" {
				// 如果连接 'user' 数据库失败，尝试连接 'users' 数据库
				conn, err = NewPGClient(conf, "users")
			}
		}
		return conn, err
	}
	for _, dbName := range databases {
		//conn, err := NewPGClient(conf, dbName)
		conn, err := connectDB(dbName)
		if err != nil {
			return nil, err
		}
		clienter.Conn[dbName] = conn
	}
	return clienter, nil
}

func NewPGClient(conf DB, dbName string) (*pgx.Conn, error) {
	connString := connStr(conf, dbName)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		zap.S().Errorw("连接数据库失败", "db", dbName, "err", err)
		return nil, err
	}
	zap.S().Infow("数据库连接成功！", "db", dbName)
	return conn, nil
}

func connStr(conf DB, db string) string {
	scheme := map[bool]string{true: "require", false: "disable"}[conf.Sslmode]
	// 对密码进行 URL 编码
	encodedPassword := url.QueryEscape(conf.Password)
	str := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		conf.Username, encodedPassword, conf.Ip, conf.Port, db, scheme)
	return str
}
