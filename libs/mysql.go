// Package libs @Author lanpang
// @Date 2024/8/23 上午10:16:00
// @Desc
package libs

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// 连接数据库的编码格式
//var charset string = "utf8"

func NewMysqlClient(conf DB, dbName string) (*sql.DB, error) {
	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		conf.Username, conf.Password, conf.Ip, conf.Port, dbName)
	if dbName == "wshoto" {
		dsn = dsn + "?interpolateParams=true"
	}
	//fmt.Println(dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("dsn格式不正确", err)
		return nil, err
	}
	// 测试连接是否成功
	err = db.Ping()
	if err != nil {
		log.Println("校验失败", err)
		return nil, err
	}
	log.Println("mysql数据库连接成功！")
	return db, nil
}

//func TestJob(tenant *Tenant) {
//	db := tenant.MysqlClient
//	rows, err := db.Query("SELECT id, username, password FROM xxl_job_user;")
//	if err != nil {
//		log.Println("数据查询失败. err:", err)
//	}
//	defer func(rows *sql.Rows) {
//		err := rows.Close()
//		if err != nil {
//			log.Printf("Failed info: %s \n", err)
//		}
//	}(rows)
//	var id int
//	var username, password string
//	for rows.Next() {
//		rows.Scan(&id, &username, &password)
//		fmt.Println(id, username, password)
//	}
//}
