// Package inspect @Author lanpang
// @Date 2024/8/6 下午6:10:00
// @Desc
package inspect

import (
	"github.com/olivere/elastic/v7"
	"strconv"
)

func NewESClient(conf Db) (*elastic.Client, string) {
	scheme := map[bool]string{true: "https", false: "http"}[conf.Sslmode]
	esurl := scheme + "://" + conf.Ip + ":" + strconv.Itoa(conf.Port)
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetScheme(scheme),
		elastic.SetURL(esurl),
		elastic.SetBasicAuth(conf.Username, conf.Password),
		elastic.SetHealthcheck(false))

	if err != nil {
		panic(err)
	}

	return client, esurl
}
