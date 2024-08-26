// Package libs @Author lanpang
// @Date 2024/8/26 下午2:24:00
// @Desc
package libs

import "github.com/jackc/pgx/v5"

type DB struct {
	Ip       string
	Port     int
	Username string
	Password string
	Sslmode  bool
}

type PGClient struct {
	Conn map[string]*pgx.Conn
}

type Rocketmq struct {
	RocketmqDashboard string
	NameServer        string
}
