// Package libs @Author lanpang
// @Date 2024/8/23 上午10:16:00
// @Desc
package libs

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

// 连接数据库的编码格式
//var charset string = "utf8"

func NewMysqlClient(conf DB, dbName string) (*sql.DB, error) {
	// 构建数据库连接字符串，增加连接超时参数（5秒）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		conf.Username, conf.Password, conf.Ip, conf.Port, dbName)
	if dbName == "wshoto" {
		dsn = dsn + "?interpolateParams=true&timeout=5s"
	} else {
		dsn = dsn + "?timeout=5s"
	}
	//fmt.Println(dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		zap.S().Errorw("dsn格式不正确", "err", err)
		return nil, err
	}
	// 测试连接是否成功
	err = db.Ping()
	if err != nil {
		zap.S().Errorw("数据库校验失败", "err", err)
		return nil, err
	}
	zap.S().Infow("mysql数据库连接成功！")
	return db, nil
}
