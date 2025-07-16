package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"vhagar/config"
)

// Provider 接口，所有 LLM 服务商需实现
// BuildRequest 构造 LLM 请求
// ParseResponse 解析 LLM 响应
type Provider interface {
	BuildRequest(input interface{}, tools []ToolDef) (*http.Request, error)
	ParseResponse(respBody []byte) (string, error)
}

// 通用 JSON 请求构造
func buildJSONRequest(ctx context.Context, method, url string, body interface{}, headers map[string]string) (*http.Request, error) {
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("[AI] buildJSONRequest error: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req, nil
}

// 通用 JSON 响应解析
type jsonUnmarshalTarget interface{}

func parseJSONResponse(data []byte, out jsonUnmarshalTarget, errMsg string) error {
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("%s: %w", errMsg, err)
	}
	return nil
}

// GeminiProvider 适配 Google Gemini API
type GeminiProvider struct {
	ApiUrl string
	ApiKey string
	Model  string
}

func (g *GeminiProvider) BuildRequest(input interface{}, tools []ToolDef) (*http.Request, error) {
	// 兼容原有逻辑，input 只支持 string
	text, _ := input.(string)
	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": text}}},
		},
	}
	headers := map[string]string{"x-goog-api-key": g.ApiKey}
	return buildJSONRequest(context.Background(), "POST", g.ApiUrl, body, headers)
}

func (g *GeminiProvider) ParseResponse(respBody []byte) (string, error) {
	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	err := parseJSONResponse(respBody, &result, "Gemini 返回格式错误")
	if err != nil {
		return "", err
	}
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("Gemini 返回内容为空")
	}
	return result.Candidates[0].Content.Parts[0].Text, nil
}

// OpenAIProvider 适配 OpenAI/OpenRouter API
type OpenAIProvider struct {
	ApiUrl string
	ApiKey string
	Model  string
}

// 支持传入 tools 参数
func (o *OpenAIProvider) BuildRequest(input interface{}, tools []ToolDef) (*http.Request, error) {
	if o.Model == "" {
		return nil, errors.New("LLM 服务商 model 配置不完整")
	}
	// 组装 tools 字段
	var toolsArr []map[string]interface{}
	for _, t := range tools {
		b, _ := json.Marshal(t)
		var m map[string]interface{}
		json.Unmarshal(b, &m)
		toolsArr = append(toolsArr, m)
	}
	// input 直接作为 messages
	body := map[string]interface{}{
		"model":    o.Model,
		"messages": input,
		"tools":    toolsArr,
		"stream":   false,
	}
	log.Printf("[AI] body: %v", body)
	headers := map[string]string{}
	if o.ApiKey != "" {
		headers["Authorization"] = "Bearer " + o.ApiKey
	}
	return buildJSONRequest(context.Background(), "POST", o.ApiUrl, body, headers)
}

func (o *OpenAIProvider) ParseResponse(respBody []byte) (string, error) {
	// 直接返回原始字符串，由 ChatWithAI 负责结构化解析
	return string(respBody), nil
}

// extractGeminiModel 从 Gemini ApiUrl 提取模型名
func extractGeminiModel(apiUrl string) string {
	idx := strings.Index(apiUrl, "/models/")
	if idx == -1 {
		return ""
	}
	rest := strings.TrimSpace(apiUrl[idx+len("/models/"):])
	end := strings.IndexAny(rest, "/:")
	if end == -1 {
		return rest
	}
	return rest[:end]
}

// newProviderImpl 工厂方法，按 provider 类型返回对应实现
func newProviderImpl(cfg *config.AICfg) (Provider, error) {
	if cfg == nil || !cfg.Enable || cfg.Provider == "" {
		return nil, errors.New("AI 配置不完整或未启用")
	}
	providerCfg, exists := cfg.Providers[cfg.Provider]
	if !exists {
		return nil, errors.New("未找到指定的 LLM 服务商配置: " + cfg.Provider)
	}
	if providerCfg.ApiKey == "" || providerCfg.ApiUrl == "" {
		return nil, errors.New("LLM 服务商配置不完整")
	}
	provider := strings.ToLower(cfg.Provider)
	if provider == "gemini" {
		model := extractGeminiModel(providerCfg.ApiUrl)
		return &GeminiProvider{
			ApiUrl: providerCfg.ApiUrl,
			ApiKey: providerCfg.ApiKey,
			Model:  model,
		}, nil
	}
	// 只要不是 gemini，全部走 OpenAIProvider 逻辑
	return &OpenAIProvider{
		ApiUrl: providerCfg.ApiUrl,
		ApiKey: providerCfg.ApiKey,
		Model:  providerCfg.Model,
	}, nil
}
