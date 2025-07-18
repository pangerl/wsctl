package chat

import (
	"context"
	"fmt"
	"log"
)

type ToolHandler func(ctx context.Context, params map[string]any) (string, error)
type ToolMeta struct {
	Name        string
	Description string
	Parameters  map[string]any
	Handler     ToolHandler
}

// OpenAI function-calling 兼容结构体
// 具体函数定义
type FunctionDef struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
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
	log.Printf("[Tools] 注册工具成功: %s", meta.Name)
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
	log.Printf("[Tools] 获取内置工具列表: %d 个工具", len(tools))
	return tools
}

// callTool 自动化工具调用
func callTool(ctx context.Context, toolName string, params map[string]any) (string, error) {
	meta, ok := toolRegistry[toolName]
	if !ok {
		return "", fmt.Errorf("未知工具: %s", toolName)
	}

	log.Printf("[Tools] 开始调用工具: %s", toolName)
	result, err := meta.Handler(ctx, params)
	if err != nil {
		log.Printf("[Tools] 工具调用失败: %s, 错误: %v", toolName, err)
		return "", fmt.Errorf("工具 %s 调用失败: %w", toolName, err)
	}
	log.Printf("[Tools] 工具调用成功: %s", toolName)
	return result, nil
}

// GetTools 获取所有注册的工具
func GetTools(ctx context.Context, mcpClient string) ([]ToolMeta, error) {
	log.Printf("[Tools] 获取所有注册的工具")
	if len(toolRegistry) == 0 {
		return nil, fmt.Errorf("没有注册任何工具")
	}

	tools := make([]ToolMeta, 0, len(toolRegistry))
	for _, meta := range toolRegistry {
		tools = append(tools, meta)
	}

	log.Printf("[Tools] 共获取到 %d 个工具", len(tools))
	return tools, nil
}

// MCP 工具初始化
func InitToolsFromMCP(ctx context.Context, mcpClient string) error {
	log.Printf("[Tools] 开始从 MCP 初始化工具")
	// 待实现
	tools, err := GetTools(ctx, mcpClient)
	if err != nil {
		return fmt.Errorf("获取 MCP 工具列表失败: %w", err)
	}

	for _, t := range tools {
		tool := t
		RegisterTool(ToolMeta{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  tool.Parameters,
			Handler: func(ctx context.Context, params map[string]any) (string, error) {
				// handler 接收调用时的 ctx
				return MCPHandler(ctx, mcpClient, tool.Name, params)
			},
		})
	}

	log.Printf("[Tools] MCP 工具初始化完成，共加载 %d 个工具", len(tools))
	return nil
}

// MCPHandler 统一转发到 MCP
// func MCPHandler(ctx context.Context, mcpClient *mcp.Client, toolName string, params map[string]any) (string, error) {
// 	log.Printf("[Tools] MCP 调用开始: %s, 参数: %v", toolName, params)
// 	result, err := mcpClient.CallTool(ctx, toolName, params)
// 	if err != nil {
// 		return "", fmt.Errorf("MCP 工具 %s 调用失败: %w", toolName, err)
// 	}
// 	log.Printf("[Tools] MCP 调用成功: %s", toolName)
// 	return result, nil
// }

func MCPHandler(ctx context.Context, mcpClient, toolName string, params map[string]any) (string, error) {
	return "", fmt.Errorf("MCPHandler 未实现")
}
