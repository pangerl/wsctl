package chat

import (
	"context"
	"encoding/json"
	"vhagar/chat/tools"
	"vhagar/libs"
)

// ToolInputSchema 工具输入参数结构
type ToolInputSchema struct {
	Type       string         `json:"type"`
	Properties map[string]any `json:"properties,omitempty"`
	Required   []string       `json:"required,omitempty"`
}

// ToolFunction 工具函数定义
type ToolFunction struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  ToolInputSchema `json:"parameters"`
}

// Tool 工具结构
type Tool struct {
	Type     string       `json:"type"` // 固定为 "function"
	Function ToolFunction `json:"function"`
}

// ToolMeta 工具元数据
type ToolMeta struct {
	Name        string                                                           `json:"name"`
	Description string                                                           `json:"description"`
	Handler     func(ctx context.Context, params map[string]any) (string, error) `json:"-"`
}

// ToolOption 工具配置选项
type ToolOption func(*ToolFunction)

// 全局工具注册表
var toolRegistry = make(map[string]ToolMeta)

// NewTool 创建新工具
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

// RegisterTool 注册工具到全局注册表
func RegisterTool(meta ToolMeta) error {
	if meta.Name == "" {
		return libs.WrapError(libs.ErrCodeToolRegFailed, "工具注册失败", libs.NewError(libs.ErrCodeInvalidParam, "工具名称不能为空"))
	}
	if meta.Handler == nil {
		return libs.WrapError(libs.ErrCodeToolRegFailed, "工具注册失败", libs.NewError(libs.ErrCodeInvalidParam, "工具处理函数不能为空"))
	}

	toolRegistry[meta.Name] = meta
	// 只有在Logger已初始化时才记录日志
	if libs.Logger != nil {
		libs.Logger.Infow("工具注册成功", "name", meta.Name, "description", meta.Description)
	}
	return nil
}

// GetToolsForAI 获取用于AI请求的工具数组
func GetToolsForAI() []map[string]any {
	var toolsArr []map[string]any

	for _, meta := range toolRegistry {
		// 为每个工具构建Tool结构
		tool := NewTool(meta.Name, func(t *ToolFunction) {
			t.Description = meta.Description
			// 根据工具类型设置参数
			if meta.Name == "weather" {
				t.Parameters.Properties["location"] = map[string]any{
					"type":        "string",
					"description": "城市名称或 LocationID 或经纬度",
				}
				t.Parameters.Required = []string{"location"}
			} else if meta.Name == "sysinfo" {
				t.Parameters.Properties["type"] = map[string]any{
					"type":        "string",
					"description": "系统信息类型",
					"enum":        []string{"cpu", "memory", "disk", "all"},
				}
				t.Parameters.Properties["details"] = map[string]any{
					"type":        "boolean",
					"description": "是否返回详细信息，默认为false",
				}
				t.Parameters.Required = []string{"type"}
			}
		})

		// 序列化工具为map
		b, err := json.Marshal(tool)
		if err != nil {
			if libs.Logger != nil {
				libs.Logger.Warnw("工具序列化失败", "tool", meta.Name, "error", err)
			}
			continue
		}
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			if libs.Logger != nil {
				libs.Logger.Warnw("工具反序列化失败", "tool", meta.Name, "error", err)
			}
			continue
		}
		toolsArr = append(toolsArr, m)
	}

	if libs.Logger != nil {
		libs.Logger.Infow("构建AI工具数组", "count", len(toolsArr))
	}
	return toolsArr
}

// CallTool 调用指定工具
func CallTool(ctx context.Context, toolName string, params map[string]any) (string, error) {
	meta, ok := toolRegistry[toolName]
	if !ok {
		err := libs.NewErrorWithDetail(libs.ErrCodeToolNotFound, "工具未找到", toolName)
		if libs.Logger != nil {
			libs.LogError(err, "工具调用")
		}
		return "", err
	}

	if libs.Logger != nil {
		libs.Logger.Infow("开始调用工具", "name", toolName, "params", params)
	}
	result, err := meta.Handler(ctx, params)
	if err != nil {
		appErr := libs.WrapError(libs.ErrCodeToolCallFailed, "工具调用失败", err)
		if libs.Logger != nil {
			libs.LogErrorWithFields(appErr, "工具调用", map[string]interface{}{
				"tool_name": toolName,
				"params":    params,
			})
		}
		return "", appErr
	}

	if libs.Logger != nil {
		libs.Logger.Infow("工具调用成功", "name", toolName)
	}
	return result, nil
}

// 初始化工具注册
func init() {
	// 注册天气工具
	if err := RegisterTool(ToolMeta{
		Name:        "weather",
		Description: "查询天气信息，支持城市名称、LocationID或经纬度",
		Handler:     tools.CallWeatherTool,
	}); err != nil {
		// 在init阶段，Logger可能还未初始化，所以使用panic而不是LogError
		panic("工具系统初始化失败: " + err.Error())
	}


	// 在init阶段不记录日志，避免Logger未初始化的问题
	// 工具注册成功的日志会在RegisterTool中记录（如果Logger已初始化）
}
