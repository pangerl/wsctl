// Package inspect @Author lanpang
// @Date 2024/8/7 下午3:43:00
// @Desc
package inspect

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"github.com/olivere/elastic/v7"
	"log"
	"time"
)

func NewInspect(corp []*Corp, es *elastic.Client, conn1, conn2 *pgx.Conn) *Inspect {
	return &Inspect{
		Corp:      corp,
		EsClient:  es,
		PgClient1: conn1,
		PgClient2: conn2,
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

func (i *Inspect) SetCustomerNum(corpid string) {
	customernum := searchCustomerNum(i.EsClient, corpid)
	for _, corp := range i.Corp {
		if corp.Corpid == corpid {
			corp.CustomerNum = customernum
			return
		}
	}
}

func (i *Inspect) SetMessageNum(corpid string, dateNow time.Time) {
	messagenum := countMessageNum(i.EsClient, corpid, dateNow)
	for _, corp := range i.Corp {
		if corp.Corpid == corpid {
			corp.MessageNum = messagenum
			return
		}
	}
}

func (i *Inspect) SetCorpName(corpid string) {
	corpName := queryCorpName(i.PgClient1, corpid)
	for _, corp := range i.Corp {
		if corp.Corpid == corpid {
			corp.CorpName = corpName
			return
		}
	}
}

func (i *Inspect) SetActiveNum(corpid string, dateNow time.Time) {
	dateDau := dateNow.AddDate(0, 0, -1)
	dateWau := dateNow.AddDate(0, 0, -7)
	dateMau := dateNow.AddDate(0, -1, 0)
	for _, corp := range i.Corp {
		if corp.Corpid == corpid {
			corp.DauNum = searchActiveNum(i.EsClient, corpid, dateDau, dateNow)
			corp.WauNum = searchActiveNum(i.EsClient, corpid, dateWau, dateNow)
			corp.MauNum = searchActiveNum(i.EsClient, corpid, dateMau, dateNow)
			return
		}
	}
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
	CheckErr(err)
	//fmt.Printf("总客户数: %d\n", searchResult.TotalHits())
	return searchResult.TotalHits()
}

func countMessageNum(client *elastic.Client, corpid string, dateNow time.Time) int64 {
	t := dateNow.AddDate(0, 0, -1)
	startTime := GetZeroTime(t).UnixNano() / 1e6
	endTime := GetZeroTime(dateNow).UnixNano() / 1e6

	// Define the query
	query := elastic.NewBoolQuery().
		Must(elastic.NewRangeQuery("msgtime").
			From(startTime). // from timestamp for yesterday 0:00:00
			To(endTime),     // to timestamp for today 0:00:00
		)
	// Make the count request
	countResult, err := client.Count().
		Index("conversation_" + corpid).
		Query(query).
		Do(context.Background())
	CheckErr(err)
	//fmt.Printf("昨天消息数: %d\n", countResult)
	return countResult
}
func searchActiveNum(client *elastic.Client, corpid string, startDate, endDate time.Time) int64 {
	startTime := GetZeroTime(startDate).UnixNano() / 1e6
	endTime := GetZeroTime(endDate).UnixNano() / 1e6
	// 创建 bool 查询
	query := elastic.NewBoolQuery().
		Must(
			elastic.NewTermsQuery("where.entrance", "001", "002", "006"),
			elastic.NewMatchQuery("who.role", "02"),
			elastic.NewTermQuery("where.corpId.keyword", corpid),
			elastic.NewRangeQuery("when.start").Gte(startTime).Lte(endTime),
		)
	searchResult, err := client.Search().
		Index("text_event_index*"). // 设置索引名
		Query(query).               // 设置查询条件
		Aggregation("dau", elastic.NewCardinalityAggregation().Field("who.id.keyword")).
		Size(0).
		Do(context.Background()) // 执行
	CheckErr(err)
	dauAgg, _ := searchResult.Aggregations["dau"]
	cardinalityAgg := &struct {
		Value int64 `json:"value"`
	}{}
	err = json.Unmarshal(dauAgg, cardinalityAgg)
	CheckErr(err)
	//fmt.Println("活跃数：", cardinalityAgg.Value)
	return cardinalityAgg.Value
}

func queryCorpName(conn *pgx.Conn, corpid string) string {
	var corpName string
	query := "SELECT corp_name FROM qw_base_tenant_corp_info WHERE tenant_id=$1 LIMIT 1"
	err := conn.QueryRow(context.Background(), query, corpid).Scan(&corpName)
	CheckErr(err)
	return corpName
}
