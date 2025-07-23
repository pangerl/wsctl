package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"vhagar/config"
	"vhagar/errors"
	"vhagar/logger"
)

// buildRequest 构造 OpenAI 兼容的 HTTP 请求
func buildRequest(ctx context.Context, messages any) (*http.Request, error) {
	cfg := &config.Config.Services.AI
	if cfg == nil || !cfg.Enable || cfg.Provider == "" {
		err := errors.NewError(errors.ErrCodeConfigInvalid, "AI 配置不完整或未启用")
		errors.LogError(err, "AI请求构建")
		return nil, err
	}

	providerCfg, exists := cfg.Providers[cfg.Provider]
	if !exists {
		err := errors.NewErrorWithDetail(errors.ErrCodeAIProviderNotFound, "未找到指定的 LLM 服务商配置", cfg.Provider)
		errors.LogError(err, "AI请求构建")
		return nil, err
	}

	if providerCfg.ApiKey == "" || providerCfg.ApiUrl == "" || providerCfg.Model == "" {
		err := errors.NewError(errors.ErrCodeConfigInvalid, "LLM 服务商配置不完整")
		errors.LogErrorWithFields(err, "AI请求构建", map[string]interface{}{
			"provider":  cfg.Provider,
			"has_key":   providerCfg.ApiKey != "",
			"has_url":   providerCfg.ApiUrl != "",
			"has_model": providerCfg.Model != "",
		})
		return nil, err
	}

	// 组装 tools 字段
	toolsArr := GetToolsForAI()

	// 构造请求体
	body := map[string]any{
		"model":    providerCfg.Model,
		"messages": messages,
		"tools":    toolsArr,
		"stream":   false,
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		appErr := errors.WrapError(errors.ErrCodeAIRequestFailed, "请求体序列化失败", err)
		errors.LogError(appErr, "AI请求构建")
		return nil, appErr
	}

	req, err := http.NewRequestWithContext(ctx, "POST", providerCfg.ApiUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		appErr := errors.WrapError(errors.ErrCodeAIRequestFailed, "HTTP请求创建失败", err)
		errors.LogErrorWithFields(appErr, "AI请求构建", map[string]interface{}{
			"url": providerCfg.ApiUrl,
		})
		return nil, appErr
	}

	req.Header.Set("Content-Type", "application/json")
	if providerCfg.ApiKey != "" {
		req.Header.Set("Authorization", "Bearer "+providerCfg.ApiKey)
	}

	logger.Logger.Infow("AI请求构建完成",
		"provider", cfg.Provider,
		"model", providerCfg.Model,
		"tools_count", len(toolsArr))
	return req, nil
}

// parseResponse 解析 OpenAI 兼容的响应
func parseResponse(respBody []byte) (string, error) {
	if len(respBody) == 0 {
		err := errors.NewError(errors.ErrCodeAIResponseInvalid, "AI响应为空")
		errors.LogError(err, "AI响应解析")
		return "", err
	}

	logger.Logger.Infow("AI响应解析完成", "response_length", len(respBody))
	// 直接返回原始字符串，由 ChatWithAI 负责结构化解析
	return string(respBody), nil
}
