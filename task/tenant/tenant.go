// Package tenant @Author lanpang
// @Date 2024/8/23 上午11:15:00
// @Desc
package tenant

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"vhagar/config"
	"vhagar/libs"
	"vhagar/notify"
	"vhagar/task"

	"github.com/olekukonko/tablewriter"

	"github.com/jackc/pgx/v5"
	"github.com/olivere/elastic/v7"
)

var isalert = false

func init() {
	task.Add(taskName, func() task.Tasker {
		return newTenant(config.Config)
	})
}

func (tenant *Tenanter) Check() {
	task.EchoPrompt("开始巡检企微租户信息")
	if tenant.Report {
		tenant.ReportRobot(tenant.Duration)
		return
	}
	tenant.TableRender()
}

func (tenant *Tenanter) TableRender() {
	tabletitle := []string{"企业名称", "会话数", "员工数", "客户数", "客户群数", "客户群人数", "日活", "周活", "月活"}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tabletitle)
	//color := tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor}
	//tableColor := []tablewriter.Colors{color, color, color, color, color, color, color, color}
	for _, corp := range tenant.Corp {
		tabledata := []string{corp.CorpName, strconv.FormatInt(corp.MessageNum, 10), strconv.Itoa(corp.UserNum),
			strconv.FormatInt(corp.CustomerNum, 10), strconv.Itoa(corp.CustomerGroupNum), strconv.Itoa(corp.CustomerGroupUserNum),
			strconv.FormatInt(corp.DauNum, 10), strconv.FormatInt(corp.WauNum, 10), strconv.FormatInt(corp.MauNum, 10)}
		table.Append(tabledata)
	}
	table.Render()
}

func (tenant *Tenanter) ReportRobot(duration time.Duration) {
	// 发送巡检报告
	markdownList := tenantRender(tenant)

	for _, markdown := range markdownList {
		notify.Send(markdown, taskName)
	}

}

func (tenant *Tenanter) ReportWshoto() {
	log.Println("推送微盛运营平台")
	// 将 []*Corp 转换为 []any
	var data = make([]any, len(tenant.Corp))
	for i, c := range tenant.Corp {
		data[i] = c
	}
	inspectBody := notify.InspectBody{
		JobType: "tenant",
		Data:    data,
	}
	err := notify.SendWshoto(&inspectBody, tenant.ProxyURL)
	if err != nil {
		return
	}
}

func (tenant *Tenanter) Gather() {
	// 创建ESClient，PGClienter
	esClient, err := libs.NewESClient(config.Config.ES)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return
	}
	pgClient, err := libs.NewPGClienter(config.Config.PG)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return
	}
	if config.Config.Customer.HasValue() {
		log.Println("读取新的customer库")
		conn, err := libs.NewPGClient(config.Config.Customer, "customer")
		if err != nil {
			log.Printf("Failed info: %s \n", err)
			return
		}
		pgClient.Conn["customer"] = conn
	}
	defer func() {
		if pgClient != nil {
			pgClient.Close()
		}
		if esClient != nil {
			esClient.Stop()
		}
	}()
	tenant.PGClient = pgClient
	tenant.ESClient = esClient
	for _, corp := range tenant.Corp {
		tenant.getTenantData(corp)
	}
	log.Print("检查成功")
}

func (tenant *Tenanter) getTenantData(corp *config.Corp) {
	// 当前时间
	dateNow := time.Now()
	if tenant.PGClient != nil {
		// 获取租户名
		tenant.SetCorpName(corp.Corpid)
		// 获取用户数
		tenant.SetUserNum(corp.Corpid)
		// 获取客户群
		tenant.SetCustomerGroupNum(corp.Corpid)
		// 获取客户群人数
		tenant.SetCustomerGroupUserNum(corp.Corpid)
	}
	if tenant.ESClient != nil {
		// 获取客户数
		tenant.SetCustomerNum(corp.Corpid)
		// 获取活跃数
		tenant.SetActiveNum(corp.Corpid, dateNow)
		// 获取会话数
		if corp.Convenabled {
			tenant.SetMessageNum(corp.Corpid, dateNow)
		}
	}
}

func tenantRender(t *Tenanter) []*notify.WeChatMarkdown {

	var inspectList []*notify.WeChatMarkdown
	isalert = false

	headString := headCorpString()

	length := len(t.Corp)
	// 每次返回8个租户的信息
	chunkSize := 8

	for n := 0; n < length; n += chunkSize {
		end := n + chunkSize
		if end > length {
			end = length
		}
		slice := t.Corp[n:end]
		markdown := tenantMarkdown(headString, slice)
		inspectList = append(inspectList, markdown)
	}
	return inspectList
}
func tenantMarkdown(headString string, Corp []*config.Corp) *notify.WeChatMarkdown {
	var builder strings.Builder
	// 添加巡检头文件
	builder.WriteString(headString)
	for _, corp := range Corp {
		// 组装租户巡检信息
		builder.WriteString(generateCorpString(corp))
	}
	if isalert {
		builder.WriteString("\n<font color='red'>**注意！巡检结果异常！**</font>" + task.CallUser(config.Config.Notify.Userlist))
	}
	markdown := &notify.WeChatMarkdown{
		MsgType: "markdown",
		Markdown: &notify.Markdown{
			Content: builder.String(),
		},
	}

	// fmt.Println("调试信息", builder.String())
	return markdown
}
func generateCorpString(corp *config.Corp) string {
	var builder strings.Builder

	builder.WriteString("> 企业名称：<font color='info'>" + corp.CorpName + "</font>\n")
	if corp.Convenabled {
		builder.WriteString("> 昨天拉取会话数：<font color='info'>" + strconv.FormatInt(corp.MessageNum, 10) + "</font>\n")
		if corp.MessageNum <= 0 {
			isalert = true
		}
	}
	builder.WriteString("> 员工人数：<font color='info'>" + strconv.Itoa(corp.UserNum) + "</font>\n")
	builder.WriteString("> 客户人数：<font color='info'>" + strconv.FormatInt(corp.CustomerNum, 10) + "</font>\n")
	builder.WriteString("> 客户群数：<font color='info'>" + strconv.Itoa(corp.CustomerGroupNum) + "</font>\n")
	builder.WriteString("> 客户群人数：<font color='info'>" + strconv.Itoa(corp.CustomerGroupUserNum) + "</font>\n")
	builder.WriteString("> 日活跃数：<font color='info'>" + strconv.FormatInt(corp.DauNum, 10) + "</font>\n")
	builder.WriteString("> 周活跃数：<font color='info'>" + strconv.FormatInt(corp.WauNum, 10) + "</font>\n")
	builder.WriteString("> 月活跃数：<font color='info'>" + strconv.FormatInt(corp.MauNum, 10) + "</font>\n")
	builder.WriteString("==================\n")

	return builder.String()
}
func headCorpString() string {
	var builder strings.Builder
	// 组装巡检内容
	builder.WriteString("# 每日巡检报告 " + version + "\n")
	builder.WriteString("**项目名称：**<font color='info'>" + config.Config.ProjectName + "</font>\n")
	builder.WriteString("**巡检时间：**<font color='info'>" + time.Now().Format("2006-01-02") + "</font>\n")
	builder.WriteString("**巡检内容：**\n")

	return builder.String()
}

// SetCustomerGroupUserNum 设置客户群人数
func (tenant *Tenanter) SetCustomerGroupUserNum(corpid string) {
	customergroupusernum, _ := queryCustomerGroupUserNum(tenant.PGClient.Conn["customer"], corpid)
	for _, corp := range tenant.Corp {
		if corp.Corpid == corpid {
			corp.CustomerGroupUserNum = customergroupusernum
			return
		}
	}
}

// SetCustomerGroupNum 设置客户群数
func (tenant *Tenanter) SetCustomerGroupNum(corpid string) {
	customergroupnum, _ := queryCustomerGroupNum(tenant.PGClient.Conn["customer"], corpid)
	for _, corp := range tenant.Corp {
		if corp.Corpid == corpid {
			corp.CustomerGroupNum = customergroupnum
			return
		}
	}
}

// SetMessageNum 统计昨天的会话数
func (tenant *Tenanter) SetMessageNum(corpid string, dateNow time.Time) {
	date := dateNow.AddDate(0, 0, -1)
	startTime := task.GetZeroTime(date).UnixNano() / 1e6
	endTime := task.GetZeroTime(dateNow).UnixNano() / 1e6
	var orgCorpId = corpid
	if strings.HasPrefix(corpid, "wpIaoBE") {
		orgCorpId, _ = queryOrgCorpId(tenant.PGClient.Conn["qv30"], corpid)
	}
	messagenum, _ := countMessageNum(tenant.ESClient, orgCorpId, startTime, endTime)
	for _, corp := range tenant.Corp {
		if corp.Corpid == corpid {
			corp.MessageNum = messagenum
			return
		}
	}
}

// SetCorpName 设置租户名称
func (tenant *Tenanter) SetCorpName(corpid string) {
	corpName, _ := queryCorpName(tenant.PGClient.Conn["qv30"], corpid)
	for _, corp := range tenant.Corp {
		if corp.Corpid == corpid {
			corp.CorpName = corpName
			return
		}
	}
}

// SetCustomerNum 设置客户数
func (tenant *Tenanter) SetCustomerNum(corpid string) {
	customernum, _ := searchCustomerNum(tenant.ESClient, corpid)
	for _, corp := range tenant.Corp {
		if corp.Corpid == corpid {
			corp.CustomerNum = customernum
			return
		}
	}
}

// SetUserNum 设置员工数
func (tenant *Tenanter) SetUserNum(corpid string) {
	userNum, _ := queryUserNum(tenant.PGClient.Conn["user"], corpid)
	for _, corp := range tenant.Corp {
		if corp.Corpid == corpid {
			corp.UserNum = userNum
			return
		}
	}
}

// SetActiveNum 设置活跃数
func (tenant *Tenanter) SetActiveNum(corpid string, dateNow time.Time) {
	dateDau := dateNow.AddDate(0, 0, -1)
	dateWau := dateNow.AddDate(0, 0, -7)
	dateMau := dateNow.AddDate(0, -1, 0)
	for _, corp := range tenant.Corp {
		if corp.Corpid == corpid {
			corp.DauNum, _ = searchActiveNum(tenant.ESClient, corpid, dateDau, dateNow)
			corp.WauNum, _ = searchActiveNum(tenant.ESClient, corpid, dateWau, dateNow)
			corp.MauNum, _ = searchActiveNum(tenant.ESClient, corpid, dateMau, dateNow)
			return
		}
	}
}

// 客户数
func searchCustomerNum(client *elastic.Client, corpid string) (int64, error) {
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
		log.Printf("Failed info: %s \n", err)
		return -1, err
	}
	//fmt.Printf("总客户数: %d\n", searchResult.TotalHits())
	return searchResult.TotalHits(), nil
}

// 会话数
func countMessageNum(client *elastic.Client, corpid string, startTime, endTime int64) (int64, error) {

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
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return -1, err
	}
	//fmt.Printf("昨天消息数: %d\n", countResult)
	return countResult, nil
}

// 活跃数
func searchActiveNum(client *elastic.Client, corpid string, startDate, endDate time.Time) (int64, error) {
	startTime := task.GetZeroTime(startDate).UnixNano() / 1e6
	endTime := task.GetZeroTime(endDate).UnixNano() / 1e6
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
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return -1, err
	}
	dauAgg, _ := searchResult.Aggregations["dau"]
	cardinalityAgg := &struct {
		Value int64 `json:"value"`
	}{}
	err = json.Unmarshal(dauAgg, cardinalityAgg)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return -1, err
	}
	//fmt.Println("活跃数：", cardinalityAgg.Value)
	return cardinalityAgg.Value, nil
}

// 租户名称
func queryCorpName(conn *pgx.Conn, corpid string) (string, error) {
	var corpName string
	query := "SELECT corp_name FROM qw_base_tenant_corp_info WHERE tenant_id=$1 LIMIT 1"
	err := conn.QueryRow(context.Background(), query, corpid).Scan(&corpName)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return "-1", err
	}
	return corpName, nil
}

// 解密 ID
func queryOrgCorpId(conn *pgx.Conn, corpid string) (string, error) {
	var orgCorpId string
	query := "SELECT org_corp_id FROM qw_base_tenant_corp_info WHERE tenant_id=$1 LIMIT 1"
	err := conn.QueryRow(context.Background(), query, corpid).Scan(&orgCorpId)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return "-1", err
	}
	return orgCorpId, nil
}

// 员工数
func queryUserNum(conn *pgx.Conn, corpid string) (int, error) {
	var userNum int
	query := "SELECT count(*) FROM qw_user WHERE deleted=0 AND tenant_id=$1 LIMIT 1"
	err := conn.QueryRow(context.Background(), query, corpid).Scan(&userNum)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return -1, err
	}
	return userNum, nil
}

// 客户群数
func queryCustomerGroupNum(conn *pgx.Conn, corpid string) (int, error) {
	var customerGroupNum int
	query := "SELECT count(1) FROM co_saas_customer_group WHERE dismiss=false AND tenant_id=$1 AND deleted_at IS NULL"
	err := conn.QueryRow(context.Background(), query, corpid).Scan(&customerGroupNum)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return -1, err
	}
	return customerGroupNum, nil
}

// 客户群人数
func queryCustomerGroupUserNum(conn *pgx.Conn, corpid string) (int, error) {
	var customerGroupUserNum int
	query := "SELECT count(1) FROM co_saas_customer_group_user WHERE type = 2 AND loss = false AND deleted_at IS NULL AND tenant_id=$1"
	err := conn.QueryRow(context.Background(), query, corpid).Scan(&customerGroupUserNum)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return -1, err
	}
	return customerGroupUserNum, nil
}

func CurrentMessageNum(client *elastic.Client, corpid string, dateNow time.Time) int64 {
	// 统计今天的会话数
	startTime := task.GetZeroTime(dateNow).UnixNano() / 1e6
	endTime := dateNow.UnixNano() / 1e6
	messagenum, _ := countMessageNum(client, corpid, startTime, endTime)
	return messagenum
}
