// Package inspect @Author lanpang
// @Date 2024/8/8 下午1:43:00
// @Desc
package inspect

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

type DBClient struct {
	conn map[string]*pgx.Conn
}

// Close 关闭所有数据库连接
func (dbClient *DBClient) Close() {
	for dbName, conn := range dbClient.conn {
		if err := conn.Close(context.Background()); err != nil {
			log.Printf("Failed to close connection for database %s: %v", dbName, err)
		}
	}
}

func NewPGClient(conf DB) (*DBClient, error) {
	dbClient := &DBClient{
		conn: make(map[string]*pgx.Conn),
	}
	databases := []string{"qv30", "user", "customer"}
	for _, dbName := range databases {
		connString := connStr(conf, dbName)
		conn, err := pgx.Connect(context.Background(), connString)
		if err != nil {
			log.Printf("Failed to connect to database %s: %s\n", dbName, err)
			return nil, err
		}
		dbClient.conn[dbName] = conn
	}

	return dbClient, nil
}

func connStr(conf DB, db string) string {
	scheme := map[bool]string{true: "require", false: "disable"}[conf.Sslmode]
	str := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		conf.Username, conf.Password, conf.Ip, conf.Port, db, scheme)
	return str
}
