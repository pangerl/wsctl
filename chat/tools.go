package chat

import (
	"context"
	"errors"
	"vhagar/chat/tools"
)

type ToolInputSchema struct {
	Type       string         `json:"type"`
	Properties map[string]any `json:"properties,omitempty"`
	Required   []string       `json:"required,omitempty"`
}

type ToolFunction struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  ToolInputSchema `json:"parameters"`
}

type Tool struct {
	Type     string       `json:"type"` // 固定为 "function"
	Function ToolFunction `json:"function"`
}

type ToolOption func(*ToolFunction)

func NewTool(name string, opts ...ToolOption) Tool {
	fn := ToolFunction{
		Name: name,
		Parameters: ToolInputSchema{
			Type:       "object",
			Properties: make(map[string]any),
			Required:   nil, // Will be omitted from JSON if empty
		},
	}
	for _, opt := range opts {
		opt(&fn)
	}
	return Tool{
		Type:     "function",
		Function: fn,
	}
}

// 注册weather工具
func RegisterWeatherTool() Tool {
	return NewTool("weather",
		func(t *ToolFunction) {
			t.Description = "查询天气"
			t.Parameters.Properties["location"] = map[string]any{
				"type":        "string",
				"description": "城市名称或 LocationID 或经纬度",
			}
			t.Parameters.Required = []string{"location"}
		},
	)
}

// 初始化
func init() {
	Tools = append(Tools, RegisterWeatherTool())
}

func callTool(ctx context.Context, name string, args map[string]any) (string, error) {
	switch name {
	case "weather":
		return tools.CallWeatherTool(ctx, args)
	default:
		return "", errors.New("未知工具: " + name)
	}
}
