package cmd

import (
	"vhagar/task"

	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "启动MCP服务",
	Long:  "通过MCP服务执行巡检任务",
	Run: func(cmd *cobra.Command, args []string) {
		task.TaskMCP(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
