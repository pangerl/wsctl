package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"vhagar/chat"
	"vhagar/chat/mcp"
	"vhagar/config"

	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "AI 聊天命令",
	Long:  `与 AI 进行基础对话的命令。`,
	Run: func(cmd *cobra.Command, args []string) {
		aiCfg := &config.Config.AI
		if !aiCfg.Enable {
			fmt.Println("AI 聊天功能未启用，请检查 config.toml 配置。")
			return
		}
		// 初始化 MCP 工具
		mcpClient := mcp.NewClient("./chat/mcp/weather-mcp-server")
		if err := chat.InitToolsFromMCP(context.Background(), mcpClient); err != nil {
			fmt.Println("MCP 工具初始化失败:", err)
			return
		}
		// 创建一个 context，当收到 SIGINT (Ctrl+C) 或 SIGTERM 信号时，该 context 会被取消
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop() // 确保在 main 函数退出时停止监听信号

		fmt.Println("请输入你的问题，按回车发送，输入 exit 或按 Ctrl+C 退出：")
		scanner := bufio.NewScanner(os.Stdin)
		var messages []interface{}
		messages = append(messages, map[string]interface{}{
			"role":    "system",
			"content": "你是一个乐于助人、无害的AI助手。",
		})
		for {
			// 检查 context 是否已被取消 (例如用户按了 Ctrl+C)
			if ctx.Err() != nil {
				fmt.Println("\n已退出 chat 模式。")
				break
			}
			fmt.Print("你: ")
			if !scanner.Scan() {
				break
			}
			input := scanner.Text()
			if input == "exit" {
				fmt.Println("已退出 chat 模式。")
				break
			}
			messages = append(messages, map[string]interface{}{
				"role":    "user",
				"content": input,
			})
			reply, err := chat.ChatWithAI(ctx, messages, aiCfg)
			if err != nil {
				fmt.Println("AI 调用出错:", err)
				continue
			}
			fmt.Println("AI:", reply)
			messages = append(messages, map[string]interface{}{
				"role":    "assistant",
				"content": reply,
			})
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
