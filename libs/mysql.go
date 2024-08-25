// Package libs @Author lanpang
// @Date 2024/8/23 上午10:16:00
// @Desc
package libs

import (
	"database/sql"
	"fmt"
	"log"
	"vhagar/inspect"
)

// 连接数据库的编码格式
var charset string = "utf8"

func NewMysqlClient(conf inspect.DB, dbName string) (*sql.DB, error) {
	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		conf.Username, conf.Password, conf.Ip, conf.Port, dbName, charset)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("dsn格式不正确,err", err)
		return nil, err
	}
	// 测试连接是否成功
	err = db.Ping()
	if err != nil {
		log.Println("校验失败,err", err)
		return nil, err
	}
	log.Println("数据库连接成功！")
	return db, nil
}
