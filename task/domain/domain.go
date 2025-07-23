// Package domain @Author Trae AI
// @Date 2024/8/23 上午11:15:00
// @Desc 域名连通性检测任务
package domain

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"vhagar/config"
	"vhagar/libs"
	"vhagar/notify"
	"vhagar/task"

	"github.com/olekukonko/tablewriter"
)

var isalert = false

func init() {
	task.Add(taskName, func() task.Tasker {
		return NewDomainer(config.Config, libs.Logger)
	})
}

// Check 实现Tasker接口，检查并展示结果
func (d *Domainer) Check() {
	if d.Config.Report {
		d.ReportRobot()
		return
	}
	d.TableRender()
}

// TableRender 表格方式展示结果
func (d *Domainer) TableRender() {
	tabletitle := []string{"域名", "端口", "连通状态"}
	table := tablewriter.NewWriter(task.GetOutputWriter())
	table.SetHeader(tabletitle)

	for _, domain := range d.Domains {
		status := "不通"
		if domain.IsAlive {
			status = "正常"
		}
		tabledata := []string{domain.Name, strconv.Itoa(domain.Port), status}
		table.Append(tabledata)
	}

	caption := fmt.Sprintf("总共检测 %d 个域名，%d 个正常，%d 个不通", d.TotalCount, d.AliveCount, d.FailedCount)
	table.SetCaption(true, caption)
	table.Render()
}

// ReportRobot 机器人方式发送报告
func (d *Domainer) ReportRobot() {
	// 发送巡检报告
	isalert = d.FailedCount > 0
	if isalert {
		headString := headString()
		markdown := domainMarkdown(headString, d)
		notify.Send(markdown, taskName)
	}
}

// Gather 实现Tasker接口，收集数据
func (d *Domainer) Gather() {
	fileName := d.Config.DomainListName
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		log.Fatalf("配置文件不存在: %s", fileName)
	}
	filePath := filepath.Join(".", fileName)
	domains, err := readDomainListFile(filePath)
	if err != nil {
		d.Logger.Errorf("Failed to read domain list file: %s", err)
		return
	}

	// 测试域名连通性
	d.TotalCount = len(domains)
	d.AliveCount = 0
	d.FailedCount = 0

	// 创建一个map来跟踪每个域名的连通状态
	domainStatusMap := make(map[string]bool)

	for _, domain := range domains {
		isAlive := testConnection(domain.Name, domain.Port)
		domain.IsAlive = isAlive

		// 更新域名状态映射
		// 如果域名已经在映射中且为true，保持true
		// 如果域名不在映射中，添加当前状态
		if currentStatus, exists := domainStatusMap[domain.Name]; exists {
			domainStatusMap[domain.Name] = currentStatus || isAlive
		} else {
			domainStatusMap[domain.Name] = isAlive
		}

		d.Domains = append(d.Domains, domain)
	}

	// 根据域名状态映射更新计数
	for _, isAlive := range domainStatusMap {
		if isAlive {
			d.AliveCount++
		} else {
			d.FailedCount++
		}
	}

	// 更新总域名数为唯一域名的数量
	d.TotalCount = len(domainStatusMap)

	d.Logger.Info("域名连通性检查完成")
}

// readDomainListFile 读取域名列表文件
func readDomainListFile(filePath string) ([]*Domain, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	domains := make([]*Domain, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // 跳过空行和注释行
		}

		// 解析域名和端口
		parts := strings.Fields(line)
		domainName := parts[0]

		// 处理多个端口的情况
		ports := []int{443} // 默认端口
		if len(parts) > 1 {
			ports = make([]int, 0)
			for _, portStr := range parts[1:] {
				port, err := strconv.Atoi(portStr)
				if err != nil {
					libs.Logger.Errorf("Invalid port for domain %s: %s", domainName, portStr)
					continue
				}
				ports = append(ports, port)
			}
		}

		// 为每个端口创建一个域名记录
		for _, port := range ports {
			domain := &Domain{
				Name: domainName,
				Port: port,
			}
			domains = append(domains, domain)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return domains, nil
}

// testConnection 测试域名连通性
func testConnection(domain string, port int) bool {
	address := fmt.Sprintf("%s:%d", domain, port)
	maxRetries := 3
	retryDelay := 1 * time.Second

	// 如果有代理配置，使用代理连接
	if config.Config.ProxyURL != "" {
		dialer := &net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}
		proxyUrl, err := url.Parse(config.Config.ProxyURL)
		if err != nil {
			libs.Logger.Errorf("Invalid proxy URL: %s", err)
			return false
		}
		transport := &http.Transport{
			Proxy:                 http.ProxyURL(proxyUrl),
			DialContext:           dialer.DialContext,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
		client := &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		}

		for i := 0; i < maxRetries; i++ {
			scheme := "http"
			if port == 443 {
				scheme = "https"
			}
			resp, err := client.Get(fmt.Sprintf("%s://%s", scheme, address))
			if err == nil {
				if resp != nil {
					resp.Body.Close()
				}
				return true
			}
			libs.Logger.Errorf("Proxy connection attempt %d failed for %s: %v", i+1, address, err)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
			}
		}
		return false
	}

	// 没有代理配置，直接连接
	for i := 0; i < maxRetries; i++ {
		conn, err := net.DialTimeout("tcp", address, 5*time.Second)
		if err == nil {
			defer func(conn net.Conn) {
				err := conn.Close()
				if err != nil {
					libs.Logger.Errorf("Failed to close connection: %s", err)
				}
			}(conn)
			return true
		}
		libs.Logger.Errorf("Direct connection attempt %d failed for %s: %v", i+1, address, err)
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}
	return false
}

// domainMarkdown 生成Markdown格式的报告
func domainMarkdown(headString string, d *Domainer) *notify.WeChatMarkdown {
	var builder strings.Builder
	// 添加巡检头文件
	builder.WriteString(headString)

	// 添加统计信息
	builder.WriteString(fmt.Sprintf("> **总共检测域名：**<font color='info'>%d</font>\n", d.TotalCount))
	builder.WriteString(fmt.Sprintf("> **正常域名数量：**<font color='info'>%d</font>\n", d.AliveCount))
	builder.WriteString(fmt.Sprintf("> **不通域名数量：**<font color='%s'>%d</font>\n", getColorByStatus(d.FailedCount > 0), d.FailedCount))
	builder.WriteString("==================\n")

	// 如果有不通的域名，列出详情
	if d.FailedCount > 0 {
		builder.WriteString("**不通域名详情：**\n")
		// 创建一个map来跟踪已经显示过的不通域名
		shownDomains := make(map[string]bool)

		for _, domain := range d.Domains {
			if !domain.IsAlive {
				// 如果这个域名还没有显示过，则显示它
				if !shownDomains[domain.Name] {
					builder.WriteString(fmt.Sprintf("> %s:%d\n", domain.Name, domain.Port))
					shownDomains[domain.Name] = true
				} else {
					// 如果已经显示过这个域名，只显示端口
					builder.WriteString(fmt.Sprintf(">   └─ 端口:%d\n", domain.Port))
				}
			}
		}
		builder.WriteString("==================\n")
	}

	if isalert {
		builder.WriteString("\n<font color='red'>**注意！域名连通性检测异常！**</font>" + task.CallUser(config.Config.Notify.Userlist))
	}

	markdown := &notify.WeChatMarkdown{
		MsgType: "markdown",
		Markdown: &notify.Markdown{
			Content: builder.String(),
		},
	}

	return markdown
}

// headString 生成报告头部信息
func headString() string {
	var builder strings.Builder
	// 组装巡检内容
	builder.WriteString("# 域名连通性检测" + "\n")
	builder.WriteString("**项目名称：**<font color='info'>" + config.Config.ProjectName + "</font>\n")
	builder.WriteString("**巡检时间：**<font color='info'>" + time.Now().Format("2006-01-02 15:04:05") + "</font>\n")
	builder.WriteString("**巡检内容：**\n")

	return builder.String()
}

// getColorByStatus 根据状态返回颜色
func getColorByStatus(isError bool) string {
	if isError {
		return "warning"
	}
	return "info"
}
