// Package libs @Author lanpang
// @Date 2024/8/6 下午6:10:00
// @Desc
package libs

import (
	"context"
	"github.com/olivere/elastic/v7"
	"log"
	"strconv"
)

func NewESClient(conf DB) (*elastic.Client, error) {
	scheme := map[bool]string{true: "https", false: "http"}[conf.Sslmode]
	esurl := scheme + "://" + conf.Ip + ":" + strconv.Itoa(conf.Port)
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetScheme(scheme),
		elastic.SetURL(esurl),
		elastic.SetBasicAuth(conf.Username, conf.Password),
		elastic.SetHealthcheck(false))

	if err != nil {
		log.Printf("Failed to create ES client: %s \n", err)
		return nil, err
	}

	// 在创建客户端后立即执行一次Ping操作，检查连接是否正常
	_, _, err = client.Ping(esurl).Do(context.Background())
	if err != nil {
		log.Printf("Failed to connect to ES: %s \n", err)
		return nil, err
	}
	log.Println("ES 连接成功！")
	return client, nil
}
