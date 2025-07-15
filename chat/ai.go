package chat

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"vhagar/config"
)

// Message 表示一条对话历史
// Role: user/assistant/tool
// ToolName/ToolResult 仅在 tool 消息时使用
// Content: 普通文本或 LLM 回复
// ToolCall: LLM function call 请求

type Message struct {
	Role       string `json:"role"`
	Content    string `json:"content"`
	ToolName   string `json:"tool_name,omitempty"`
	ToolInput  string `json:"tool_input,omitempty"`
	ToolResult string `json:"tool_result,omitempty"`
}

// callLLM 单轮调用大模型，返回回复内容
func callLLM(messages []Message, tools []ToolDef, cfg *config.AICfg) (string, error) {
	prompt := buildPrompt(messages, tools)
	provider, err := NewProvider(cfg)
	if err != nil {
		log.Printf("[AI] provider init error: %v", err)
		return "", err
	}
	req, err := provider.BuildRequest(prompt)
	if err != nil {
		log.Printf("[AI] build request error: %v", err)
		return "", err
	}
	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	duration := time.Since(start)
	if err != nil {
		log.Printf("[AI] http request error: %v, duration: %v", err, duration)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("[AI] bad status: %s, duration: %v", resp.Status, duration)
		return "", fmt.Errorf("AI 接口请求失败，状态码: %s", resp.Status)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[AI] read response error: %v", err)
		return "", err
	}
	result, err := provider.ParseResponse(respBody)
	if err != nil {
		log.Printf("[AI] parse response error: %v", err)
		return "", err
	}
	return result, nil
}

// ChatWithAI 多轮 function calling 工具调用主流程
func ChatWithAI(messages []Message, cfg *config.AICfg) (string, error) {
	maxTurns := 5
	for turn := 0; turn < maxTurns; turn++ {
		result, err := callLLM(messages, getBuiltinTools(), cfg)
		if err != nil {
			return "", err
		}
		// 解析 LLM 回复，判断是否需要调用工具
		toolName, toolInput := parseToolCall(result)
		log.Printf("[AI] toolName: %s, toolInput: %s", toolName, toolInput)
		if toolName != "" {
			toolResult, err := callTool(toolName, toolInput)
			if err != nil {
				return "", err
			}
			messages = append(messages, Message{Role: "tool", ToolName: toolName, ToolInput: toolInput, ToolResult: toolResult})
			continue
		}
		return result, nil
	}
	return "", fmt.Errorf("多轮工具调用超出最大轮数")
}

// buildPrompt 构造带 tools 列表和历史的 prompt，适配 ToolDef 结构
func buildPrompt(messages []Message, tools []ToolDef) string {
	var sb strings.Builder
	sb.WriteString("你可以调用如下工具：\n")
	for _, t := range tools {
		paramStr := ""
		if t.Function.Parameters != nil {
			b, _ := json.Marshal(t.Function.Parameters)
			paramStr = string(b)
		}
		sb.WriteString(fmt.Sprintf("- %s: %s\n参数: %s\n", t.Function.Name, t.Function.Description, paramStr))
	}
	sb.WriteString("\n对话历史：\n")
	for _, m := range messages {
		sb.WriteString(fmt.Sprintf("[%s] %s\n", m.Role, m.Content))
		if m.Role == "tool" {
			sb.WriteString(fmt.Sprintf("[tool_result] %s\n", m.ToolResult))
		}
	}
	sb.WriteString("\n请根据用户需求，决定是否需要调用工具。如需调用，请回复：\nCALL <tool_name> <json参数>\n否则直接回复最终答案。\n")
	return sb.String()
}

// parseToolCall 解析 LLM 回复，判断是否需要调用工具
func parseToolCall(reply string) (string, string) {
	// 约定格式：CALL <tool_name> <json参数>
	if strings.HasPrefix(reply, "CALL ") {
		parts := strings.SplitN(reply, " ", 3)
		if len(parts) == 3 {
			return parts[1], parts[2]
		}
	}
	return "", ""
}

// Summarize 对输入内容进行AI总结，突出异常和重点
func Summarize(content string) (string, error) {
	prompt := "请对以下巡检内容进行简要总结，突出异常和重点：\n" + content
	messages := []Message{{Role: "user", Content: prompt}}
	return ChatWithAI(messages, &config.Config.AI)
}

// Provider 工厂
// provider.go 中有 Provider 接口及各实现
func NewProvider(cfg *config.AICfg) (Provider, error) {
	return newProviderImpl(cfg)
}
