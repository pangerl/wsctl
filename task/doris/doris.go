// Package doris @Author lanpang
// @Date 2024/8/23 下午7:02:00
// @Desc
package doris

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"vhagar/config"
	"vhagar/libs"
	"vhagar/notify"
	"vhagar/task"

	"github.com/olekukonko/tablewriter"
)

//var doriser *Doris

//func Work() config.Tasker {
//	cfg := config.Config
//	doris := newDoris(cfg)
//	// 初始化数据
//	doris.Gather()
//	return doris
//}

func init() {
	task.Add(taskName, func() task.Tasker {
		return NewDoris(config.Config, task.GetLogger())
	})
}

func (doris *Doris) Name() string {
	return taskName
}

func (doris *Doris) Check() {
	if doris.Config.Global.Report {
		doris.ReportRobot()
		return
	}
	doris.TableRender()
}

func (doris *Doris) TableRender() {
	tabletitle := []string{"BE 节点总数", "BE 可用节点数", "员工统计表", "使用分析表", "客户群统计表"}
	table := tablewriter.NewWriter(task.GetOutputWriter())
	table.SetHeader(tabletitle)
	tabledata := []string{strconv.Itoa(doris.TotalBackendNum), strconv.Itoa(doris.OnlineBackendNum),
		strconv.Itoa(doris.StaffCount), strconv.Itoa(doris.UseAnalyseCount), strconv.Itoa(doris.CustomerGroupCount)}
	table.Append(tabledata)
	caption := fmt.Sprintf("Job失败数: %d.", len(doris.FailedJobs))
	table.SetCaption(true, caption)
	table.Render()
	for _, jobName := range doris.FailedJobs {
		fmt.Println("JobName: ", jobName)
	}
}

func (doris *Doris) Gather() {
	mysqlClinet, err := libs.NewMysqlClient(doris.Config.Doris.DB, "wshoto")
	if err != nil {
		libs.Logger.Errorw("Failed to create mysql client", "err", err)
		return
	}
	defer func() {
		if mysqlClinet != nil {
			err := mysqlClinet.Close()
			if err != nil {
				return
			}
		}
	}()
	doris.MysqlClient = mysqlClinet
	// 获取当前零点时间
	todayTime := task.GetZeroTime(time.Now())
	yesterday := todayTime.AddDate(0, 0, -1)
	yesterdayTime := task.GetZeroTime(yesterday)
	// 失败任务
	failedJobs := selectFailedJob(todayTime.Format("2006-01-02 15:04:05"), doris.MysqlClient)
	doris.FailedJobs = failedJobs
	// 员工统计表
	staffCount := selectStaffCount(yesterdayTime.Format("2006-01-02 15:04:05"), doris.MysqlClient)
	doris.StaffCount = staffCount
	// 使用分析表
	useAnalyseCount := selectUseAnalyseCount(yesterdayTime.Format("2006-01-02 15:04:05"), doris.MysqlClient)
	doris.UseAnalyseCount = useAnalyseCount
	// 客户群统计表
	customerGroupCount := selectCustomerGroupCount(yesterdayTime.Format("2006-01-02 15:04:05"), doris.MysqlClient)
	doris.CustomerGroupCount = customerGroupCount
	// 检查 BE 节点健康
	getBENum(doris)
}

func (doris *Doris) ReportRobot() {
	var builder strings.Builder

	failedJobCount := len(doris.FailedJobs)
	// 组装巡检内容
	builder.WriteString("# Doris 巡检 \n")
	builder.WriteString("**项目名称：**<font color='info'>" + config.Config.Global.ProjectName + "</font>\n")
	builder.WriteString("**巡检时间：**<font color='info'>" + time.Now().Format("2006-01-02") + "</font>\n")
	builder.WriteString("**BE节点总数：**<font color='info'>" + strconv.Itoa(doris.TotalBackendNum) + "</font>\n")
	builder.WriteString("**在线节点数：**<font color='info'>" + strconv.Itoa(doris.OnlineBackendNum) + "</font>\n")

	builder.WriteString("==================\n")

	builder.WriteString("**Job失败数：**<font color='info'>" + strconv.Itoa(failedJobCount) + "</font>\n")

	for _, jobName := range doris.FailedJobs {
		builder.WriteString("> <font color='red'>" + jobName + "</font>\n")
	}
	builder.WriteString("==================\n")
	builder.WriteString("**昨天各表增量数据**\n\n")
	builder.WriteString("**员工统计表：**<font color='info'>" + strconv.Itoa(doris.StaffCount) + "</font>\n")
	builder.WriteString("**使用分析表：**<font color='info'>" + strconv.Itoa(doris.UseAnalyseCount) + "</font>\n")
	builder.WriteString("**客户群统计表：**<font color='info'>" + strconv.Itoa(doris.CustomerGroupCount) + "</font>\n")

	markdown := &notify.WeChatMarkdown{
		MsgType: "markdown",
		Markdown: &notify.Markdown{
			Content: builder.String(),
		},
	}

	notify.Send(markdown, taskName)
}

// 查询失败的job
func selectFailedJob(queryTime string, db *sql.DB) []string {
	// 定义查询语句
	query := `
	SELECT
	   name
	FROM
	   sys_job
	WHERE
	   frequency = 'd'
	   AND status = 1
	   AND name NOT LIKE '%dwd_%'
	   AND name != 'ads_bi_mbr_staff_ptt_mall_statistics'
	   AND name != 'ads_bi_mbr_staff_sales_conversion'
	   AND last_execute_time < ?`
	rows, err := db.Query(query, queryTime)
	if err != nil {
		libs.Logger.Errorw("数据查询失败", "err", err)
		return nil
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			libs.Logger.Errorw("Failed info", "err", err)
		}
	}(rows)

	// 创建一个切片来存储结果
	var failedJobs []string

	// 处理查询结果
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			libs.Logger.Errorw("Failed info", "err", err)
			return nil
		}
		failedJobs = append(failedJobs, name)
	}
	// 输出结果
	//fmt.Println("Failed Job Names:", failedJobs)
	return failedJobs
}

// 查询员工统计表
func selectStaffCount(queryTime string, db *sql.DB) int {
	// 打印 queryTime
	//fmt.Println("queryTime:", queryTime)
	// 定义查询语句
	query := `
	SELECT
		count(1) data_cnt
	from
		ads_bi_mbr_staff_pull_new_d
	where
		ds = ?
		and date_type = 'day';`
	rows := db.QueryRow(query, queryTime)
	// 处理查询结果
	var staffCount int
	err := rows.Scan(&staffCount)
	if err != nil {
		libs.Logger.Errorw("Failed info", "err", err)
		return -1
	}
	return staffCount
}

// 使用分析表
func selectUseAnalyseCount(queryTime string, db *sql.DB) int {
	// 定义查询语句
	query := `
	SELECT
		count(1) data_cnt
	from
		qw_user_use_analyse_d
	where
		ds = ?;`
	rows := db.QueryRow(query, queryTime)
	// 处理查询结果
	var useAnalyseCount int
	err := rows.Scan(&useAnalyseCount)
	// 打印查询语句，拼装完整 sql
	libs.Logger.Infow("query:", "query", query, "queryTime", queryTime)
	if err != nil {
		libs.Logger.Errorw("Failed info", "err", err)
		return -1
	}
	return useAnalyseCount
}

// 客户群统计表
func selectCustomerGroupCount(queryTime string, db *sql.DB) int {
	// 定义查询语句
	query := `
	SELECT
		count(1) data_cnt
	from
		dws_customer_group_st_h
	where
		ds = ?;`
	rows := db.QueryRow(query, queryTime)
	// 处理查询结果
	var customerGroupCount int
	err := rows.Scan(&customerGroupCount)
	if err != nil {
		libs.Logger.Errorw("Failed info", "err", err)
		return -1
	}
	return customerGroupCount
}

func getBENum(doris *Doris) {
	healthUrl := fmt.Sprintf("http://%s:%d/api/health", doris.Config.Doris.Ip, doris.DorisCfg.HttpPort)

	// 发起 HTTP GET 请求
	body := task.DoRequest(healthUrl)
	if body == nil { // 请求失败
		libs.Logger.Errorf("Failed to get response from %s", healthUrl)
		return
	}

	// 解析 JSON 响应
	var response dorisResponse
	if err := json.Unmarshal(body, &response); err != nil {
		libs.Logger.Errorf("Error unmarshalling json: %v", err)
		return
	}

	doris.OnlineBackendNum = response.Data.OnlineBackendNum
	doris.TotalBackendNum = response.Data.TotalBackendNum
}
