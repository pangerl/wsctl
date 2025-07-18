package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"vhagar/config"
)

// buildRequest 构造 OpenAI 兼容的 HTTP 请求
func buildRequest(ctx context.Context, messages any, tools []ToolDef) (*http.Request, error) {
	cfg := &config.Config.AI
	if cfg == nil || !cfg.Enable || cfg.Provider == "" {
		return nil, errors.New("AI 配置不完整或未启用")
	}

	providerCfg, exists := cfg.Providers[cfg.Provider]
	if !exists {
		return nil, errors.New("未找到指定的 LLM 服务商配置: " + cfg.Provider)
	}

	if providerCfg.ApiKey == "" || providerCfg.ApiUrl == "" || providerCfg.Model == "" {
		return nil, errors.New("LLM 服务商配置不完整")
	}

	// 组装 tools 字段
	var toolsArr []map[string]any
	for _, t := range tools {
		b, _ := json.Marshal(t)
		var m map[string]any
		json.Unmarshal(b, &m)
		toolsArr = append(toolsArr, m)
	}

	// 构造请求体
	body := map[string]any{
		"model":    providerCfg.Model,
		"messages": messages,
		"tools":    toolsArr,
		"stream":   false,
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", providerCfg.ApiUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("[AI] buildRequest error: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if providerCfg.ApiKey != "" {
		req.Header.Set("Authorization", "Bearer "+providerCfg.ApiKey)
	}

	log.Printf("[AI] request body: %v", body)
	return req, nil
}

// parseResponse 解析 OpenAI 兼容的响应
func parseResponse(respBody []byte) (string, error) {
	// 直接返回原始字符串，由 ChatWithAI 负责结构化解析
	return string(respBody), nil
}
