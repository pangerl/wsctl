package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"vhagar/chat/mcp"
)

type ToolHandler func(input string) (string, error)
type ToolMeta struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
	Handler     ToolHandler
}

// OpenAI function-calling 兼容结构体
// 具体函数定义
type FunctionDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// 工具定义，外层 type/function
type ToolDef struct {
	Type     string      `json:"type"`
	Function FunctionDef `json:"function"`
}

var toolRegistry = map[string]ToolMeta{}

// RegisterTool 注册工具
func RegisterTool(meta ToolMeta) {
	toolRegistry[meta.Name] = meta
}

// getBuiltinTools 返回 OpenAI 兼容结构
func getBuiltinTools() []ToolDef {
	tools := make([]ToolDef, 0, len(toolRegistry))
	for _, meta := range toolRegistry {
		tools = append(tools, ToolDef{
			Type: "function",
			Function: FunctionDef{
				Name:        meta.Name,
				Description: meta.Description,
				Parameters:  meta.Parameters,
			},
		})
	}
	log.Printf("[AI] get builtin tools: %v", tools)
	return tools
}

// callTool 自动化工具调用
func callTool(toolName, toolInput string) (string, error) {
	meta, ok := toolRegistry[toolName]
	if !ok {
		return "", fmt.Errorf("未知工具: %s", toolName)
	}
	return meta.Handler(toolInput)
}

// MCP 工具初始化
func InitToolsFromMCP(ctx context.Context, mcpClient *mcp.Client) error {
	tools, err := mcpClient.ListTools(ctx)
	if err != nil {
		return err
	}
	for _, t := range tools {
		RegisterTool(ToolMeta{
			Name:        t.Name,
			Description: t.Description,
			Parameters:  t.Parameters,
			Handler:     func(input string) (string, error) { return MCPHandler(ctx, mcpClient, t.Name, input) },
		})
	}
	return nil
}

// MCPHandler 统一转发到 MCP
func MCPHandler(ctx context.Context, mcpClient *mcp.Client, toolName, input string) (string, error) {
	// input 是 JSON 字符串，解析为 map
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(input), &params); err != nil {
		return "", fmt.Errorf("参数解析失败: %v", err)
	}
	return mcpClient.CallTool(ctx, toolName, params)
}
