package task

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"vhagar/config"
)

var (
	outputWriter io.Writer
	outputFile   *os.File
	once         sync.Once
)

const outputFileName = "task_output.log" // 固定文件名

// GetOutputWriter 返回全局唯一的 io.Writer，写入终端和文件
func GetOutputWriter() io.Writer {
	once.Do(func() {
		file, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			// 打开失败只写终端
			outputWriter = os.Stdout
			return
		}
		outputFile = file
		outputWriter = io.MultiWriter(os.Stdout, file)
	})
	return outputWriter
}

// CloseOutputFile 关闭文件，建议在 main 退出时调用
func CloseOutputFile() {
	if outputFile != nil {
		_ = outputFile.Close()
	}
}

// ClearOutputFile 清空日志文件内容
func ClearOutputFile() error {
	// 关闭当前文件句柄
	if outputFile != nil {
		_ = outputFile.Close()
		outputFile = nil
		outputWriter = nil
	}
	// 以截断方式重新打开
	file, err := os.OpenFile(outputFileName, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

// AISummarize 读取巡检内容并调用 AI 总结
func AISummarize(filename string) (string, error) {
	aiCfg := config.Config.AI
	if !aiCfg.Enable || aiCfg.Provider == "" {
		return "", errors.New("AI 配置不完整或未启用")
	}

	// 获取指定服务商的配置
	providerCfg, exists := aiCfg.Providers[aiCfg.Provider]
	if !exists {
		return "", errors.New("未找到指定的 LLM 服务商配置: " + aiCfg.Provider)
	}

	if providerCfg.ApiKey == "" || providerCfg.ApiUrl == "" || providerCfg.Model == "" {
		return "", errors.New("LLM 服务商配置不完整")
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	prompt := "请对以下巡检内容进行简要总结，突出异常和重点：\n" + string(content)

	// 构造 AI 请求
	body := map[string]interface{}{
		"model": providerCfg.Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
	jsonBody, _ := json.Marshal(body)
	client := &http.Client{}
	request, err := http.NewRequest("POST", providerCfg.ApiUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+providerCfg.ApiKey)

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("AI 接口请求失败，状态码:" + resp.Status)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// 解析 AI 返回
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", errors.New("AI 接口返回格式错误: " + err.Error())
	}
	if len(result.Choices) == 0 {
		return "", errors.New("AI 返回内容为空")
	}
	return result.Choices[0].Message.Content, nil
}
