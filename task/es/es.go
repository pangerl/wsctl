package es

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"vhagar/config"
	"vhagar/libs"
	"vhagar/notifier"
	"vhagar/task"

	"github.com/olekukonko/tablewriter"
)

const taskName = "es"

func init() {
	task.Add(taskName, func() task.Tasker {
		return newES(config.Config)
	})
}

func (es *ES) Gather() {
	esClient, err := libs.NewESClient(config.Config.ES)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return
	}
	defer func() {
		if esClient != nil {
			esClient.Stop()
		}
	}()
	es.ESClient = esClient
	es.getESInfo()
}

func (es *ES) Check() {
	task.EchoPrompt("开始巡检 ES 状态信息")
	if es.Report {
		// 发送机器人
		es.ReportRobot(es.Global.Duration)
		return
	}
	es.TableRender()
}

func (es *ES) getESInfo() {

	// 检查集群健康状态
	health, err := es.ESClient.ClusterHealth().Do(context.Background())
	if err != nil {
		fmt.Printf("获取集群健康状态失败: %s\n", err)
		return
	}
	es.Status = health.Status

	// 获取节点统计信息
	stats, err := es.ESClient.NodesStats().Do(context.Background())
	if err != nil {
		log.Printf("Failed to get node stats: %s\n", err)
		return
	}

	// 填充 NodeList
	for _, node := range stats.Nodes {
		nodeInfo := &NodeInfo{
			Name:      node.Name,
			IP:        node.IP,
			JVMUsage:  float64(node.JVM.Mem.HeapUsedInBytes) / float64(node.JVM.Mem.HeapMaxInBytes) * 100,
			DiskUsage: float64(node.FS.Total.TotalInBytes-node.FS.Total.AvailableInBytes) / float64(node.FS.Total.TotalInBytes) * 100,
		}
		es.NodeList = append(es.NodeList, nodeInfo)
	}
}

func (es *ES) TableRender() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"节点名称", "IP地址", "JVM堆内存使用(%)", "磁盘使用(%)"})

	for _, node := range es.NodeList {
		row := []string{
			node.Name,
			node.IP,
			strconv.FormatFloat(node.JVMUsage, 'f', 2, 64),
			strconv.FormatFloat(node.DiskUsage, 'f', 2, 64),
		}
		table.Append(row)
	}

	table.SetColumnSeparator("|")
	table.SetCenterSeparator("+")

	caption := fmt.Sprintf("ES 集群状态: %s.", es.Status)
	table.SetCaption(true, caption)
	table.Render()
}

func (es *ES) ReportRobot(duration time.Duration) {
	// 发送巡检报告
	markdown := esRender(es, es.ProjectName)
	log.Println("任务等待时间", duration)
	time.Sleep(duration)
	for _, robotkey := range es.Notifier["doris"].Robotkey {
		_ = notifier.SendWecom(markdown, robotkey, es.ProxyURL)
	}

}

func esRender(es *ES, name string) *notifier.WeChatMarkdown {
	var builder strings.Builder
	builder.WriteString("# ES 巡检 \n")
	builder.WriteString("**项目名称：**<font color='info'>" + name + "</font>\n")
	builder.WriteString("**巡检时间：**<font color='info'>" + time.Now().Format("2006-01-02") + "</font>\n")
	builder.WriteString("**集群状态：<font color='info'>" + es.Status + "</font>**\n")

	for _, node := range es.NodeList {
		builder.WriteString("==================\n")
		builder.WriteString("## 节点名称：<font color='info'>" + node.Name + "</font>\n")
		builder.WriteString("**IP地址：**<font color='info'>" + node.IP + "</font>\n")
		builder.WriteString(fmt.Sprintf("**JVM堆内存使用：**<font color='info'> %.2f%% </font>\n", node.JVMUsage))
		builder.WriteString(fmt.Sprintf("**磁盘使用：**<font color='info'> %.2f%% </font>\n", node.DiskUsage))
		builder.WriteString("\n")
	}

	// 添加警告信息
	warnings := es.generateWarnings()
	if len(warnings) > 0 {
		builder.WriteString("警告:\n")
		for _, warning := range warnings {
			builder.WriteString(fmt.Sprintf("- %s\n", warning))
		}
	}

	markdown := &notifier.WeChatMarkdown{
		MsgType: "markdown",
		Markdown: &notifier.Markdown{
			Content: builder.String(),
		},
	}

	return markdown
}

func (es *ES) generateWarnings() []string {
	var warnings []string

	for _, node := range es.NodeList {
		if node.JVMUsage > 80 {
			warnings = append(warnings, fmt.Sprintf("节点 %s JVM堆内存使用率高: %.2f%%", node.Name, node.JVMUsage))
		}
		if node.DiskUsage > 80 {
			warnings = append(warnings, fmt.Sprintf("节点 %s 磁盘使用率高: %.2f%%", node.Name, node.DiskUsage))
		}
	}

	if strings.ToLower(es.Status) != "green" {
		warnings = append(warnings, fmt.Sprintf("集群状态不佳: %s", es.Status))
	}

	return warnings
}
