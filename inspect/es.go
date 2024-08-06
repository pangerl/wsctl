// Package inspect @Author lanpang
// @Date 2024/8/6 下午6:10:00
// @Desc
package inspect

import "github.com/olivere/elastic/v7"

func NewESClient() *elastic.Client {
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL("http://localhost:9200"),
		elastic.SetBasicAuth("user", "secret"))

	if err != nil {
		panic(err)
	}
	return client
}
