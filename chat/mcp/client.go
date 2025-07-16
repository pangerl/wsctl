package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
)

// ToolMeta 结构体，添加了 json 标签以便直接解析。
type ToolMeta struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Annotations map[string]interface{} `json:"annotations"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// 统一的 JSON-RPC 请求结构
type mcpRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      string      `json:"id"`
}

// Client 用于与 mcp-server 通信。
// 建议：cmdPath 应为预编译的二进制文件路径，而非 go 源文件。
type Client struct {
	cmdPath string
}

func NewClient(cmdPath string) *Client {
	return &Client{cmdPath: cmdPath}
}

// ListTools 启动 MCP server，发送 tools/list 请求，返回工具列表。
func (c *Client) ListTools(ctx context.Context) ([]ToolMeta, error) {
	// 使用预编译的二进制文件，而不是 "go run"
	cmd := exec.CommandContext(ctx, c.cmdPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("获取 stdin 失败: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("获取 stdout 失败: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("启动 mcp-server 失败: %v", err)
	}

	req := mcpRequest{
		JSONRPC: "2.0",
		Method:  "tools/list",
		Params:  map[string]interface{}{},
		Id:      "1",
	}
	if err := json.NewEncoder(stdin).Encode(req); err != nil {
		// 在写入失败后，最好尝试终止子进程
		cmd.Process.Kill()
		cmd.Wait()
		return nil, fmt.Errorf("写入请求失败: %v", err)
	}
	// 写入完成后立即关闭 stdin，这对很多子进程是必要的信号
	stdin.Close()

	type mcpListResponse struct {
		Result struct {
			Tools []ToolMeta `json:"tools"`
		} `json:"result"`
		Error interface{} `json:"error"`
	}

	var resp mcpListResponse
	if err := json.NewDecoder(stdout).Decode(&resp); err != nil {
		cmd.Wait()
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("mcp-server 进程退出时发生错误: %v", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("mcp-server 错误: %v", resp.Error)
	}

	return resp.Result.Tools, nil
}

// CallTool 启动 MCP server，发送 tools/call 请求，返回结果。
func (c *Client) CallTool(ctx context.Context, toolName string, params map[string]interface{}) (string, error) {
	cmd := exec.CommandContext(ctx, c.cmdPath)
	cmd.Env = os.Environ()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("获取 stdin 失败: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("获取 stdout 失败: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("获取 stderr 失败: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("启动 mcp-server 失败: %v", err)
	}

	req := mcpRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name":      toolName,
			"arguments": params,
		},
		Id: "1",
	}

	if err := json.NewEncoder(stdin).Encode(req); err != nil {
		cmd.Process.Kill()
		cmd.Wait()
		return "", fmt.Errorf("写入请求失败: %v", err)
	}
	stdin.Close()

	// 并发读取 stdout 和 stderr，避免因一个管道满了而阻塞另一个
	var stdoutBytes, stderrBytes []byte
	var wg sync.WaitGroup
	var stdoutErr, stderrErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		stdoutBytes, stdoutErr = io.ReadAll(stdout)
	}()
	go func() {
		defer wg.Done()
		stderrBytes, stderrErr = io.ReadAll(stderr)
	}()
	wg.Wait()

	if len(stderrBytes) > 0 {
		log.Printf("[CallTool] 子进程标准错误: %s", string(stderrBytes))
	}

	if stdoutErr != nil {
		cmd.Wait()
		return "", fmt.Errorf("读取子进程 stdout 失败: %w", stdoutErr)
	}
	if stderrErr != nil {
		cmd.Wait()
		return "", fmt.Errorf("读取子进程 stderr 失败: %w", stderrErr)
	}

	// 等待进程结束并检查退出状态
	if err := cmd.Wait(); err != nil {
		log.Printf("mcp-server 进程退出时发生错误: %v, stderr: %s", err, string(stderrBytes))
	}

	if len(stdoutBytes) == 0 {
		return "", fmt.Errorf("mcp-server 响应为空, stderr: %s", string(stderrBytes))
	}

	log.Printf("[CallTool] 子进程标准输出: %s", string(stdoutBytes))

	type mcpCallResponse struct {
		Result struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"result"`
		Error interface{} `json:"error"`
	}

	var resp mcpCallResponse
	if err := json.Unmarshal(stdoutBytes, &resp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w, 响应原文: %s", err, string(stdoutBytes))
	}

	if resp.Error != nil {
		return "", fmt.Errorf("mcp-server 错误: %v", resp.Error)
	}

	if len(resp.Result.Content) == 0 {
		return "", fmt.Errorf("mcp-server 响应内容为空")
	}

	// 假定我们只需要第一个内容块的文本
	return resp.Result.Content[0].Text, nil
}
