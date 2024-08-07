// Package inspect @Author lanpang
// @Date 2024/8/7 下午3:43:00
// @Desc
package inspect

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
)

func NewInspect(corp []Corp, es *elastic.Client) *Inspect {
	return &Inspect{
		Corp:     corp,
		EsClient: es,
	}
}

func (i *Inspect) GetVersion(url string) {
	//查看es当前版本
	version, err := i.EsClient.ElasticsearchVersion(url)
	if err != nil {
		log.Println("查询es版本错误", err)
	}
	log.Println("Elasticsearch version: ", version)
}

func (i *Inspect) GetCustomerNum(corpid string) {
	searchCustomerNum(i.EsClient, corpid)
}

func searchCustomerNum(client *elastic.Client, corpid string) int64 {
	// 创建 bool 查询
	query := elastic.NewBoolQuery().
		Filter(
			elastic.NewTermQuery("tenantId", corpid),
			elastic.NewTermQuery("relatedHiddenAt", 0),
			elastic.NewTermQuery("relatedDelAt", 0),
		)
	searchResult, err := client.Search().
		Index("customer_related_1"). // 设置索引名
		Query(query).                // 设置查询条件
		TrackTotalHits(true).
		Do(context.Background()) // 执行
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("总客户数: %d\n", searchResult.TotalHits())
	return searchResult.TotalHits()
}
