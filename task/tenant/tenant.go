// Package tenant @Author lanpang
// @Date 2024/8/23 上午11:15:00
// @Desc
package tenant

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
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

func init() {
	task.Add(taskName, func() task.Tasker {
		return NewTenanter(config.Config, libs.Logger)
	})
}

func (tenant *Tenanter) Check() {
	if tenant.Config.Report {
		tenant.ReportRobot()
		return
	}
	tenant.TableRender()
}

func (tenant *Tenanter) TableRender() {
	tabletitle := []string{"企业名称", "员工数", "客户数", "客户群数", "客户群人数", "日活", "周活", "月活"}
	table := tablewriter.NewWriter(task.GetOutputWriter())
	table.SetHeader(tabletitle)
	//color := tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor}
	//tableColor := []tablewriter.Colors{color, color, color, color, color, color, color, color}
	for _, corp := range tenant.Corp {
		tabledata := []string{corp.CorpName, strconv.Itoa(corp.UserNum),
			strconv.FormatInt(corp.CustomerNum, 10), strconv.Itoa(corp.CustomerGroupNum), strconv.Itoa(corp.CustomerGroupUserNum),
			strconv.Itoa(corp.DauNum), strconv.Itoa(corp.WauNum), strconv.Itoa(corp.MauNum)}
		table.Append(tabledata)
	}
	table.Render()
}

func (tenant *Tenanter) ReportRobot() {
	// 发送巡检报告
	markdownList := tenantRender(tenant)

	for _, markdown := range markdownList {
		notify.Send(markdown, taskName)
	}

}

func (tenant *Tenanter) ReportWshoto() {
	libs.Logger.Warnw("推送微盛运营平台")
	// 将 []*Corp 转换为 []any
	var data = make([]any, len(tenant.Corp))
	for i, c := range tenant.Corp {
		data[i] = c
	}
	inspectBody := notify.InspectBody{
		JobType: "tenant",
		Data:    data,
	}
	err := notify.SendWshoto(&inspectBody, tenant.Config.ProxyURL)
	if err != nil {
		return
	}
}

func (tenant *Tenanter) Gather() {
	// 创建ESClient，PGClienter
	//esClient, err := libs.NewESClient(config.Config.ES)
	//if err != nil {
	//	log.Printf("Failed info: %s \n", err)
	//	return
	//}
	// 创建 mysqlClinet，PGCliente
	mysqlClinet, err := libs.NewMysqlClient(config.Config.Doris.DB, "wshoto")
	if err != nil {
		libs.Logger.Errorw("Failed to create mysql client", "err", err)
		return
	}
	pgClient, err := libs.NewPGClienter(config.Config.PG)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return
	}
	if config.Config.Customer.HasValue() {
		libs.Logger.Info("读取新的customer库")
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
		if mysqlClinet != nil {
			err := mysqlClinet.Close()
			if err != nil {
				return
			}
		}
	}()
	tenant.PGClient = pgClient
	tenant.MysqlClient = mysqlClinet
	for _, corp := range tenant.Corp {
		tenant.getTenantData(corp)
	}
	libs.Logger.Info("检查成功")
}

func (tenant *Tenanter) getTenantData(corp *config.Corp) {
	// 当前时间
	dateNow := time.Now()
	if tenant.PGClient != nil {
		// 获取租户名
		tenant.SetCorpName(corp.Corpid)
		// 获取用户数
		tenant.SetUserNum(corp.Corpid)
		// 获取客户数
		tenant.SetCustomerNum(corp.Corpid)
		// 获取客户群
		tenant.SetCustomerGroupNum(corp.Corpid)
		// 获取客户群人数
		tenant.SetCustomerGroupUserNum(corp.Corpid)
	}
	if tenant.MysqlClient != nil {
		// 获取活跃数
		tenant.SetActiveNum(corp.Corpid, dateNow)
	}
}

func tenantRender(t *Tenanter) []*notify.WeChatMarkdown {

	var inspectList []*notify.WeChatMarkdown

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
	builder.WriteString("> 员工人数：<font color='info'>" + strconv.Itoa(corp.UserNum) + "</font>\n")
	builder.WriteString("> 客户人数：<font color='info'>" + strconv.FormatInt(corp.CustomerNum, 10) + "</font>\n")
	builder.WriteString("> 客户群数：<font color='info'>" + strconv.Itoa(corp.CustomerGroupNum) + "</font>\n")
	builder.WriteString("> 客户群人数：<font color='info'>" + strconv.Itoa(corp.CustomerGroupUserNum) + "</font>\n")
	builder.WriteString("> 日活跃数：<font color='info'>" + strconv.Itoa(corp.DauNum) + "</font>\n")
	builder.WriteString("> 周活跃数：<font color='info'>" + strconv.Itoa(corp.WauNum) + "</font>\n")
	builder.WriteString("> 月活跃数：<font color='info'>" + strconv.Itoa(corp.MauNum) + "</font>\n")
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
	//customernum, _ := searchCustomerNum(tenant.ESClient, corpid)
	customerNum, _ := queryCustomerNum(tenant.PGClient.Conn["customer"], corpid)
	for _, corp := range tenant.Corp {
		if corp.Corpid == corpid {
			corp.CustomerNum = customerNum
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

func GetZeroTimeBeforeNDays(dateNow time.Time, n int) string {
	now := dateNow
	// 获取今天零点
	todayZero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// 减去 n 天
	targetZero := todayZero.AddDate(0, 0, -n)
	return targetZero.Format("2006-01-02 15:04:05")
}

// SetActiveNum 设置活跃数
func (tenant *Tenanter) SetActiveNum(corpid string, dateNow time.Time) {
	// 获取当前零点时间
	todayTime := task.GetZeroTime(dateNow).Format("2006-01-02 15:04:05")
	dateDau := GetZeroTimeBeforeNDays(dateNow, 1)
	dateWau := GetZeroTimeBeforeNDays(dateNow, 7)
	dateMau := GetZeroTimeBeforeNDays(dateNow, 30)
	for _, corp := range tenant.Corp {
		if corp.Corpid == corpid {
			corp.DauNum = queryActiveNum(corpid, dateDau, todayTime, tenant.MysqlClient)
			corp.WauNum = queryActiveNum(corpid, dateWau, todayTime, tenant.MysqlClient)
			corp.MauNum = queryActiveNum(corpid, dateMau, todayTime, tenant.MysqlClient)
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

// 活跃数（新）
func queryActiveNum(corpid, startDate, endDate string, db *sql.DB) int {
	// 定义查询语句
	query := `
		SELECT COUNT(DISTINCT who_id)
		FROM tracking_event_log_h
		WHERE corpId = ?
		  AND where_entrance IN ('001', '002', '006')
		  AND who_role = '02'
		  AND gmt_create >= ?
		  AND gmt_create < ?;`
	libs.Logger.Infof("queryActiveNum SQL: %s | args: %s, %s, %s", strings.ReplaceAll(query, "\n", " "), corpid, startDate, endDate)
	rows := db.QueryRow(query, corpid, startDate, endDate)
	// 处理查询结果
	var activeNum int
	err := rows.Scan(&activeNum)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return -1
	}
	return activeNum
}

// 活跃数（旧）
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

// 客户数
func queryCustomerNum(conn *pgx.Conn, corpid string) (int64, error) {
	var customerNum int64
	query := "SELECT count(1) FROM co_saas_customer_related WHERE deleted=0 AND tenant_id=$1 LIMIT 1"
	err := conn.QueryRow(context.Background(), query, corpid).Scan(&customerNum)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return -1, err
	}
	return customerNum, nil
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
	query := "SELECT count(1) FROM co_saas_customer_group_user WHERE loss = false AND deleted_at IS NULL AND tenant_id=$1"
	err := conn.QueryRow(context.Background(), query, corpid).Scan(&customerGroupUserNum)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return -1, err
	}
	return customerGroupUserNum, nil
}
