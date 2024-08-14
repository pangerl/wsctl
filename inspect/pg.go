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

func NewPGClient(conf DB) (*pgx.Conn, *pgx.Conn, *pgx.Conn) {
	connString1 := connStr(conf, "qv30")
	connString2 := connStr(conf, "user")
	connString3 := connStr(conf, "customer")
	conn1, err := pgx.Connect(context.Background(), connString1)
	conn2, _ := pgx.Connect(context.Background(), connString2)
	conn3, _ := pgx.Connect(context.Background(), connString3)
	if err != nil {
		log.Printf("Failed to create PG client: %s \n", err)
		return nil, nil, nil
	}
	return conn1, conn2, conn3
}

func connStr(conf DB, db string) string {
	scheme := map[bool]string{true: "require", false: "disable"}[conf.Sslmode]
	str := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		conf.Username, conf.Password, conf.Ip, conf.Port, db, scheme)
	return str
}
