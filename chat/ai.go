package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// callLLM 单轮调用大模型，返回回复内容
func callLLM(ctx context.Context, messages []any, tools []ToolDef) (string, error) {
	req, err := buildRequest(ctx, messages, tools)
	if err != nil {
		log.Printf("[AI] build request error: %v", err)
		return "", err
	}

	start := time.Now()
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
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

	result, err := parseResponse(respBody)
	if err != nil {
		log.Printf("[AI] parse response error: %v", err)
		return "", err
	}
	return result, nil
}

// ChatWithAI 多轮 function calling 工具调用主流程
func ChatWithAI(ctx context.Context, messages []any) (string, error) {
	// 检查大模型的配置
	maxTurns := 5
	for turn := 0; turn < maxTurns; turn++ {
		result, err := callLLM(ctx, messages, getBuiltinTools())
		log.Printf("[AI] result: %s", result)
		if err != nil {
			return "", err
		}
		// 1. 解析 result，提取 message 和 finish_reason
		var resp struct {
			Choices []struct {
				FinishReason string `json:"finish_reason"`
				Message      struct {
					Role      string `json:"role"`
					Content   string `json:"content"`
					ToolCalls []struct {
						ID       string `json:"id"`
						Type     string `json:"type"`
						Function struct {
							Name      string `json:"name"`
							Arguments string `json:"arguments"`
						} `json:"function"`
					} `json:"tool_calls"`
				} `json:"message"`
			} `json:"choices"`
		}
		err = json.Unmarshal([]byte(result), &resp)
		if err != nil || len(resp.Choices) == 0 {
			return "", errors.New("LLM 返回格式错误")
		}
		choice := resp.Choices[0]
		msg := choice.Message
		finishReason := choice.FinishReason
		// 2. 根据 finish_reason 处理
		switch finishReason {
		case "stop":
			// 对话完成，返回内容
			return msg.Content, nil
		case "tool_calls":
			assistantMessage := map[string]any{
				"role":       msg.Role,
				"tool_calls": msg.ToolCalls,
			}
			if msg.Content != "" {
				assistantMessage["content"] = msg.Content
			}
			// 记录 message
			messages = append(messages, assistantMessage)
			if len(msg.ToolCalls) == 0 {
				return "", errors.New("LLM 返回 tool_calls 但内容为空")
			}
			for _, tc := range msg.ToolCalls {
				// 解析 arguments
				var args map[string]any
				if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
					toolResult := "参数解析失败: " + err.Error()
					messages = append(messages, map[string]any{
						"role":         "tool",
						"tool_call_id": tc.ID,
						"name":         tc.Function.Name,
						"content":      toolResult,
					})
					continue
				}
				log.Printf("[AI] %s args: %v", tc.Function.Name, args)
				toolResult, err := callTool(ctx, tc.Function.Name, args)
				if err != nil {
					log.Printf("[AI] 工具 %s 调用失败: %v", tc.Function.Name, err)
					toolResult = fmt.Sprintf("工具 %s 调用失败", tc.Function.Name)
				}
				messages = append(messages, map[string]any{
					"role":         "tool",
					"tool_call_id": tc.ID,
					"name":         tc.Function.Name,
					"content":      toolResult,
				})
			}
			// 继续下一轮
			continue
		default:
			log.Printf("[AI] finish_reason 异常: %s, result: %s", finishReason, result)
			return "", fmt.Errorf("LLM finish_reason 异常: %s", finishReason)
		}
	}
	return "", errors.New("多轮工具调用超出最大轮数")
}

// Summarize 对输入内容进行AI总结，突出异常和重点
func Summarize(ctx context.Context, content string) (string, error) {
	prompt := "请对以下巡检内容进行简要总结，突出异常和重点：\n" + content
	messages := []any{
		map[string]any{
			"role":    "user",
			"content": prompt,
		},
	}
	return ChatWithAI(ctx, messages)
}
