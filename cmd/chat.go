package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"vhagar/chat"
	"vhagar/config"

	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "AI 聊天命令",
	Long:  `与 AI 进行基础对话的命令。`,
	Run: func(cmd *cobra.Command, args []string) {
		aiCfg := &config.Config.AI
		if !aiCfg.Enable || aiCfg.Provider == "" {
			fmt.Println("AI 聊天功能未启用，请检查 config.toml 配置。")
			return
		}
		// 初始化 MCP 工具
		// 创建一个 context，当收到 SIGINT (Ctrl+C) 或 SIGTERM 信号时，该 context 会被取消
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop() // 确保在 main 函数退出时停止监听信号

		fmt.Println("请输入你的问题，按回车发送，输入 exit 或按 Ctrl+C 退出：")
		scanner := bufio.NewScanner(os.Stdin)
		var messages []any
		messages = append(messages, map[string]any{
			"role":    "system",
			"content": "你是一个乐于助人、无害的AI助手。",
		})
		for {
			fmt.Print("你: ")
			if !scanner.Scan() {
				break
			}
			input := scanner.Text()
			if input == "exit" {
				fmt.Println("已退出 chat 模式。")
				break
			}
			messages = append(messages, map[string]any{
				"role":    "user",
				"content": input,
			})
			reply, err := chat.ChatWithAI(ctx, messages)
			if err != nil {
				// 检查是否是用户取消操作
				if ctx.Err() != nil {
					fmt.Println("\n操作已取消，退出 chat 模式。")
					break
				}
				fmt.Println("AI 调用出错:", err)
				continue
			}
			fmt.Println("AI:", reply)
			messages = append(messages, map[string]any{
				"role":    "assistant",
				"content": reply,
			})
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
