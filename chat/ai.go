package chat

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"vhagar/config"
)

// ChatWithAI 主流程，调度 provider 并记录日志
func ChatWithAI(input string, cfg *config.AICfg) (string, error) {
	provider, err := NewProvider(cfg)
	if err != nil {
		log.Printf("[AI] provider init error: %v", err)
		return "", err
	}
	req, err := provider.BuildRequest(input)
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
	log.Printf("[AI] request success, model: %s, duration: %v", cfg.Providers[cfg.Provider].Model, duration)
	return result, nil
}

// Summarize 对输入内容进行AI总结，突出异常和重点
func Summarize(content string) (string, error) {
	prompt := "请对以下巡检内容进行简要总结，突出异常和重点：\n" + content
	return ChatWithAI(prompt, &config.Config.AI)
}

// Provider 工厂
// provider.go 中有 Provider 接口及各实现
func NewProvider(cfg *config.AICfg) (Provider, error) {
	return newProviderImpl(cfg)
}
