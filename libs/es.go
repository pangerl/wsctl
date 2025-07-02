// Package libs @Author lanpang
// @Date 2024/8/6 下午6:10:00
// @Desc
package libs

import (
	"context"
	"strconv"

	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
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
		zap.S().Errorw("创建 ES client 失败", "err", err)
		return nil, err
	}

	// 在创建客户端后立即执行一次Ping操作，检查连接是否正常
	_, _, err = client.Ping(esurl).Do(context.Background())
	if err != nil {
		zap.S().Errorw("连接 ES 失败", "err", err)
		return nil, err
	}
	zap.S().Infow("ES 连接成功！")
	return client, nil
}
