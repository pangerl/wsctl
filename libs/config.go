// Package libs @Author lanpang
// @Date 2024/8/26 下午2:24:00
// @Desc
package libs

import "github.com/jackc/pgx/v5"

type DB struct {
	Ip       string `toml:"ip"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Sslmode  bool   `toml:"sslmode"`
}

type PGClient struct {
	Conn map[string]*pgx.Conn
}
