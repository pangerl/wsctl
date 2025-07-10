package cmd

import (
	"bufio"
	"fmt"
	"os"
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
		if !aiCfg.Enable {
			fmt.Println("AI 聊天功能未启用，请检查 config.toml 配置。")
			return
		}
		fmt.Println("请输入你的问题，按回车发送，输入 exit 退出：")
		scanner := bufio.NewScanner(os.Stdin)
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
			reply, err := chat.ChatWithAI(input, aiCfg)
			if err != nil {
				fmt.Println("AI 调用出错:", err)
				continue
			}
			fmt.Println("AI:", reply)
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
