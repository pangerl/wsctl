// Package inspect @Author lanpang
// @Date 2024/8/7 下午3:43:00
// @Desc
package inspect

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/olivere/elastic/v7"
	"log"
	"strconv"
	"strings"
	"time"
)

var isalert = false

func NewInspect(corp []*Corp, es *elastic.Client, conn1, conn2 *pgx.Conn, name, version string) *Inspect {
	return &Inspect{
		ProjectName: name,
		Version:     version,
		Corp:        corp,
		EsClient:    es,
		PgClient1:   conn1,
		PgClient2:   conn2,
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

func (i *Inspect) TransformToMarkdown(users []string, dateNow time.Time) (markdown *WeChatMarkdown, err error) {

	var builder strings.Builder
	isalert = false

	// 组装巡检内容
	builder.WriteString("# 每日巡检报告 " + i.Version + "\n")
	builder.WriteString("**项目名称：**<font color='info'>" + i.ProjectName + "</font>\n")
	builder.WriteString("**巡检时间：**<font color='info'>" + dateNow.Format("2006-01-02") + "</font>\n")
	builder.WriteString("**巡检内容：**\n")

	for _, corp := range i.Corp {

		builder.WriteString(generateCorpString(corp))
	}
	if isalert {
		builder.WriteString("\n<font color='red'>**注意！巡检结果异常！**</font>" + calluser(users))
	}

	markdown = &WeChatMarkdown{
		MsgType: "markdown",
		Markdown: &Markdown{
			Content: builder.String(),
		},
	}

	return
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

func (i *Inspect) SetUserNum(corpid string) {
	userNum := queryUserNum(i.PgClient2, corpid)
	for _, corp := range i.Corp {
		if corp.Corpid == corpid {
			corp.UserNum = userNum
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
func queryUserNum(conn *pgx.Conn, corpid string) int {
	var userNum int
	query := "SELECT count(*) from qw_user WHERE deleted=0 AND tenant_id=$1 LIMIT 1"
	err := conn.QueryRow(context.Background(), query, corpid).Scan(&userNum)
	CheckErr(err)
	return userNum
}

func calluser(users []string) string {
	var result string
	if len(users) == 0 {
		return result
	}
	for _, user := range users {
		result += fmt.Sprintf("<@%s>", user)
	}
	return result
}

func getcolor(num int64, ignore bool) string {
	color := "info"
	if num == 0 {
		if ignore {
			color = "warning"
		} else {
			color = "red"
		}
	}
	return color
}

func generateCorpString(corp *Corp) string {
	var builder strings.Builder

	builder.WriteString("> 企业名称：<font color='info'>" + corp.CorpName + "</font>\n")
	if corp.Convenabled {
		builder.WriteString("> 昨天拉取会话数：<font color='info'>" + strconv.FormatInt(corp.MessageNum, 10) + "</font>\n")
		if corp.MessageNum == 0 {
			isalert = true
		}
	}
	builder.WriteString("> 员工数统计：<font color='info'>" + strconv.Itoa(corp.UserNum) + "</font>\n")
	builder.WriteString("> 客户数统计：<font color='info'>" + strconv.FormatInt(corp.CustomerNum, 10) + "</font>\n")
	builder.WriteString("> 日活数统计：<font color='info'>" + strconv.FormatInt(corp.DauNum, 10) + "</font>\n")
	builder.WriteString("> 周活数统计：<font color='info'>" + strconv.FormatInt(corp.WauNum, 10) + "</font>\n")
	builder.WriteString("> 月活数统计：<font color='info'>" + strconv.FormatInt(corp.MauNum, 10) + "</font>\n")
	builder.WriteString("==================\n")

	return builder.String()
}
