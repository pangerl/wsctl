// Package message @Author lanpang
// @Date 2024/8/23 上午11:15:00
// @Desc
package message

import (
	"context"
	"fmt"
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
var ispush = false

func init() {
	task.Add(taskName, func() task.Tasker {
		return newTenant(config.Config)
	})
}

func (tenant *Tenanter) Check() {
	if tenant.Report {
		if ispush {
			tenant.ReportRobot()
		}
		return
	}
	tenant.TableRender()
}

func (tenant *Tenanter) TableRender() {
	tabletitle := []string{"企业名称", "当前会话数", "昨天会话数"}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tabletitle)
	//color := tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor}
	//tableColor := []tablewriter.Colors{color, color, color, color, color, color, color, color}
	for _, corp := range tenant.Corp {
		if corp.Convenabled {
			tabledata := []string{corp.CorpName, strconv.FormatInt(corp.MessageNum, 10), strconv.FormatInt(corp.YesterdayMessageNum, 10)}
			table.Append(tabledata)
		}
	}
	if tenant.NasDir != "" {
		caption := fmt.Sprintf("当天目录创建状态: %s .", strconv.FormatBool(tenant.DirIsExis))
		table.SetCaption(true, caption)
	}
	table.Render()
}

func (tenant *Tenanter) ReportRobot() {
	// 发送巡检报告
	isalert = false
	headString := headCorpString()
	if tenant.NasDir != "" {
		caption := fmt.Sprintf("**数据目录状态: ** <font color='warning'>%s</font>", strconv.FormatBool(tenant.DirIsExis))
		headString += caption + "\n"
	}
	markdown := tenantMarkdown(headString, tenant.Corp)
	notify.Send(markdown, taskName)
}

func (tenant *Tenanter) ReportWshoto() {
	libs.Logger.Infow("推送微盛运营平台")
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
	ispush = false
	// 创建ESClient，PGClienter
	esClient, err := libs.NewESClient(config.Config.ES)
	if err != nil {
		libs.Logger.Errorw("Failed info", "err", err)
		return
	}
	pgClient, err := libs.NewPGClienter(config.Config.PG)
	if err != nil {
		libs.Logger.Errorw("Failed info", "err", err)
		return
	}
	if config.Config.Customer.HasValue() {
		libs.Logger.Info("读取新的customer库")
		conn, err := libs.NewPGClient(config.Config.Customer, "customer")
		if err != nil {
			libs.Logger.Errorw("Failed info", "err", err)
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
		if corp.Convenabled {
			ispush = true
			tenant.getTenantData(corp)
		}
	}
	if tenant.NasDir != "" {
		tenant.DirIsExis = checkDirectoryExistence(tenant.NasDir)
	}
	libs.Logger.Info("检查成功")
}

func (tenant *Tenanter) getTenantData(corp *config.Corp) {
	// 当前时间
	dateNow := time.Now()
	if tenant.PGClient != nil {
		// 获取租户名
		tenant.SetCorpName(corp.Corpid)
	}
	if tenant.ESClient != nil {
		// 获取会话数
		tenant.SetMessageNum(corp.Corpid, dateNow)
		tenant.SetYesterdayMessageNum(corp.Corpid, dateNow)
	}
}

//func tenantRender(t *Tenanter) *notify.WeChatMarkdown {
//
//	return markdown
//}

func tenantMarkdown(headString string, Corp []*config.Corp) *notify.WeChatMarkdown {
	var builder strings.Builder
	// 添加巡检头文件
	builder.WriteString(headString)
	for _, corp := range Corp {
		if corp.Convenabled {
			builder.WriteString(generateCorpString(corp))
		}
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
	builder.WriteString("> 当前拉取会话数：<font color='info'>" + strconv.FormatInt(corp.MessageNum, 10) + "</font>\n")
	builder.WriteString("> 昨天拉取会话数：<font color='info'>" + strconv.FormatInt(corp.YesterdayMessageNum, 10) + "</font>\n")
	if corp.MessageNum <= 0 && corp.YesterdayMessageNum <= 0 {
		isalert = true
	}
	builder.WriteString("==================\n")
	return builder.String()
}
func headCorpString() string {
	var builder strings.Builder
	// 组装巡检内容
	builder.WriteString("# 会话数巡检 \n")
	builder.WriteString("**项目名称：**<font color='info'>" + config.Config.ProjectName + "</font>\n")
	builder.WriteString("**巡检时间：**<font color='info'>" + time.Now().Format("2006-01-02") + "</font>\n")
	builder.WriteString("**巡检内容：**\n")

	return builder.String()
}

// SetMessageNum 统计当前的会话数
func (tenant *Tenanter) SetMessageNum(corpid string, dateNow time.Time) {
	startTime := task.GetZeroTime(dateNow).UnixNano() / 1e6
	endTime := dateNow.UnixNano() / 1e6
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

// SetYesterdayMessageNum 统计昨天的会话数
func (tenant *Tenanter) SetYesterdayMessageNum(corpid string, dateNow time.Time) {
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
			corp.YesterdayMessageNum = messagenum
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
		libs.Logger.Errorw("Failed info", "err", err)
		return -1, err
	}
	//fmt.Printf("昨天消息数: %d\n", countResult)
	return countResult, nil
}

// 租户名称
func queryCorpName(conn *pgx.Conn, corpid string) (string, error) {
	var corpName string
	query := "SELECT corp_name FROM qw_base_tenant_corp_info WHERE tenant_id=$1 LIMIT 1"
	err := conn.QueryRow(context.Background(), query, corpid).Scan(&corpName)
	if err != nil {
		libs.Logger.Errorw("Failed info", "err", err)
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
		libs.Logger.Errorw("Failed info", "err", err)
		return "-1", err
	}
	return orgCorpId, nil
}

func CurrentMessageNum(client *elastic.Client, corpid string, dateNow time.Time) int64 {
	// 统计当前的会话数
	startTime := task.GetZeroTime(dateNow).UnixNano() / 1e6
	endTime := dateNow.UnixNano() / 1e6
	messagenum, _ := countMessageNum(client, corpid, startTime, endTime)
	return messagenum
}

// 检查指定路径下是否存在当天日期命名的目录
func checkDirectoryExistence(path string) bool {
	// 获取当前日期
	currentDate := time.Now().Format("20060102")

	// 拼接出完整的目录路径
	dirPath := fmt.Sprintf("%s/%s", path, currentDate)

	libs.Logger.Info(dirPath)
	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return false // 目录不存在
	}
	return true // 目录存在
}
