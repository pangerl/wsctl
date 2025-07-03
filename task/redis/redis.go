package redis

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"vhagar/config"
	"vhagar/libs"
	"vhagar/notify"
	"vhagar/task"

	"os"

	"github.com/olekukonko/tablewriter"
)

func init() {
	task.Add(taskName, func() task.Tasker {
		return NewRedis(config.Config, libs.Logger)
	})
}

func (redis *Redis) Check() {
	//task.EchoPrompt("开始巡检 Redis 状态信息")
	redis.Gather()
	if config.Config.Report {
		// 发送机器人
		redis.ReportRobot()
		return
	}
	redis.TableRender()
}

func (redis *Redis) Gather() {
	redisClient, err := libs.NewRedisClient(redis.Config.Redis)
	if err != nil {
		log.Println("Failed to create redis client. err:", err)
		redis.Logger.Errorw("Failed to create redis client", "err", err)
		return
	}

	defer func() {
		if redisClient != nil {
			err := redisClient.Close()
			if err != nil {
				return
			}
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := redisClient.Info(ctx).Result()
	if err != nil {
		log.Printf("无法获取 Redis 信息: %v", err)
		redis.Logger.Error("无法获取 Redis 信息: %v", err)
		return
	}

	infoMap := parseRedisInfo(info)

	redis.Version = infoMap["redis_version"]
	redis.Role = infoMap["role"]
	redis.Slaves, _ = strconv.Atoi(infoMap["connected_slaves"])
	redis.CurrentClients, _ = strconv.Atoi(infoMap["connected_clients"])
	redis.MaxClients, _ = strconv.Atoi(infoMap["maxclients"])
	redis.UsedMemory = formatMemory(infoMap["used_memory"])
	redis.KeyCount, _ = strconv.Atoi(strings.Split(strings.Split(infoMap["db0"], ",")[0], "=")[1])
}

func (redis *Redis) TableRender() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"属性", "值"})
	table.SetBorder(false)
	table.AppendBulk([][]string{
		{"版本", redis.Version},
		{"角色", redis.Role},
		{"从节点数", strconv.Itoa(redis.Slaves)},
		{"当前连接数", strconv.Itoa(redis.CurrentClients)},
		{"最大连接数", strconv.Itoa(redis.MaxClients)},
		{"使用内存", redis.UsedMemory},
		{"键数量", strconv.Itoa(redis.KeyCount)},
	})
	table.Render()
}

func (redis *Redis) Name() string {
	return "Redis"
}

func parseRedisInfo(info string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(info, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return result
}

func formatMemory(memoryStr string) string {
	memory, _ := strconv.ParseInt(memoryStr, 10, 64)
	if memory < 1024 {
		return fmt.Sprintf("%d B", memory)
	} else if memory < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(memory)/1024)
	} else if memory < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(memory)/(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB", float64(memory)/(1024*1024*1024))
	}
}

func (redis *Redis) ReportRobot() {
	var builder strings.Builder
	// 组装巡检内容
	builder.WriteString("# Redis 巡检 \n")
	builder.WriteString("**项目名称：**<font color='info'>" + config.Config.ProjectName + "</font>\n")
	builder.WriteString("**巡检时间：**<font color='info'>" + time.Now().Format("2006-01-02") + "</font>\n")
	builder.WriteString("**版本：**<font color='info'>" + redis.Version + "</font>\n")
	builder.WriteString("**角色：**<font color='info'>" + redis.Role + "</font>\n")
	builder.WriteString("**从节点数：**<font color='info'>" + strconv.Itoa(redis.Slaves) + "</font>\n")
	builder.WriteString("**当前连接数：**<font color='info'>" + strconv.Itoa(redis.CurrentClients) + "</font>\n")
	builder.WriteString("**最大连接数：**<font color='info'>" + strconv.Itoa(redis.MaxClients) + "</font>\n")
	builder.WriteString("**使用内存：**<font color='info'>" + redis.UsedMemory + "</font>\n")
	builder.WriteString("**键数量：**<font color='info'>" + strconv.Itoa(redis.KeyCount) + "</font>\n")

	markdown := &notify.WeChatMarkdown{
		MsgType: "markdown",
		Markdown: &notify.Markdown{
			Content: builder.String(),
		},
	}

	notify.Send(markdown, taskName)
}
