package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"vhagar/config"
	"vhagar/errors"
	"vhagar/logger"
)

// Tools 变量已移至 tools.go 中的 toolRegistry 统一管理

// callLLM 单轮调用大模型，返回回复内容
func callLLM(ctx context.Context, messages []any) (string, error) {
	req, err := buildRequest(ctx, messages)
	if err != nil {
		return "", err // buildRequest已经处理了错误日志
	}

	aiCfg := config.Config.Services.AI
	fullModelName := aiCfg.Provider + "/" + aiCfg.Providers[aiCfg.Provider].Model

	start := time.Now()
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		appErr := errors.WrapError(errors.ErrCodeNetworkFailed, "AI HTTP请求失败", err)
		errors.LogErrorWithFields(appErr, "AI调用", map[string]interface{}{
			"model":    fullModelName,
			"duration": duration,
		})
		return "", appErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		appErr := errors.NewErrorWithDetail(errors.ErrCodeAIRequestFailed, "AI接口返回错误状态码", resp.Status)
		errors.LogErrorWithFields(appErr, "AI调用", map[string]interface{}{
			"model":       fullModelName,
			"status_code": resp.StatusCode,
			"duration":    duration,
		})
		return "", appErr
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		appErr := errors.WrapError(errors.ErrCodeNetworkFailed, "读取AI响应失败", err)
		errors.LogError(appErr, "AI调用")
		return "", appErr
	}

	result, err := parseResponse(respBody)
	if err != nil {
		return "", err // parseResponse已经处理了错误日志
	}

	fmt.Println("AI调用成功", "model", fullModelName, "duration", duration)
	return result, nil
}

// ChatWithAI 多轮 function calling 工具调用主流程
func ChatWithAI(ctx context.Context, messages []any) (string, error) {
	maxTurns := 5
	logger.Logger.Infow("开始AI对话", "max_turns", maxTurns, "initial_messages", len(messages))

	for turn := 0; turn < maxTurns; turn++ {
		logger.Logger.Infow("AI对话轮次", "turn", turn+1, "messages_count", len(messages))

		result, err := callLLM(ctx, messages)
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
			appErr := errors.WrapError(errors.ErrCodeAIResponseInvalid, "LLM返回格式错误", err)
			errors.LogErrorWithFields(appErr, "AI对话", map[string]interface{}{
				"turn":     turn + 1,
				"response": result,
			})
			return "", appErr
		}

		choice := resp.Choices[0]
		msg := choice.Message
		finishReason := choice.FinishReason

		logger.Logger.Infow("AI响应解析", "turn", turn+1, "finish_reason", finishReason, "tool_calls_count", len(msg.ToolCalls))

		// 2. 根据 finish_reason 处理
		switch finishReason {
		case "stop":
			// 对话完成，返回内容
			logger.Logger.Infow("AI对话完成", "turn", turn+1, "content_length", len(msg.Content))
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
				err := errors.NewError(errors.ErrCodeAIResponseInvalid, "LLM返回tool_calls但内容为空")
				errors.LogError(err, "AI对话")
				return "", err
			}

			for _, tc := range msg.ToolCalls {
				// 解析 arguments
				var args map[string]any
				if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
					logger.Logger.Warnw("工具参数解析失败", "tool", tc.Function.Name, "args", tc.Function.Arguments, "error", err)
					toolResult := "参数解析失败: " + err.Error()
					messages = append(messages, map[string]any{
						"role":         "tool",
						"tool_call_id": tc.ID,
						"name":         tc.Function.Name,
						"content":      toolResult,
					})
					continue
				}

				logger.Logger.Infow("调用工具", "tool", tc.Function.Name, "args", args)
				toolResult, err := CallTool(ctx, tc.Function.Name, args)
				if err != nil {
					errors.LogErrorWithFields(err, "工具调用", map[string]interface{}{
						"tool": tc.Function.Name,
						"args": args,
					})
					toolResult = fmt.Sprintf("工具 %s 调用失败: %v", tc.Function.Name, err)
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
			appErr := errors.NewErrorWithDetail(errors.ErrCodeAIResponseInvalid, "LLM finish_reason异常", finishReason)
			errors.LogErrorWithFields(appErr, "AI对话", map[string]interface{}{
				"turn":          turn + 1,
				"finish_reason": finishReason,
				"response":      result,
			})
			return "", appErr
		}
	}

	err := errors.NewErrorWithDetail(errors.ErrCodeAIRequestFailed, "多轮工具调用超出最大轮数", fmt.Sprintf("max_turns=%d", maxTurns))
	errors.LogError(err, "AI对话")
	return "", err
}

// Summarize 对输入内容进行AI总结，突出异常和重点
func Summarize(ctx context.Context, content string) (string, error) {
	if content == "" {
		err := errors.NewError(errors.ErrCodeInvalidParam, "巡检内容不能为空")
		errors.LogError(err, "AI总结")
		return "", err
	}

	prompt := "请对以下巡检内容进行简要总结，突出异常和重点：\n" + content
	messages := []any{
		map[string]any{
			"role":    "user",
			"content": prompt,
		},
	}

	logger.Logger.Infow("开始AI总结", "content_length", len(content))
	result, err := ChatWithAI(ctx, messages)
	if err != nil {
		errors.LogErrorWithFields(err, "AI总结", map[string]interface{}{
			"content_length": len(content),
		})
		return "", err
	}

	logger.Logger.Infow("AI总结完成", "result_length", len(result))
	return result, nil
}
