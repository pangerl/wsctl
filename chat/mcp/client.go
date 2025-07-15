package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
)

type WeatherArgs struct {
	Location string `json:"location"`
}

type WeatherResp struct {
	Location string `json:"location"`
	Weather  string `json:"weather"`
	Temp     string `json:"temp"`
	Wind     string `json:"wind"`
	Humidity string `json:"humidity"`
	Time     string `json:"time"`
}

// MCPClient 用于与 mcp-server 通信
// 这里简单实现为每次调用都启动子进程，后续可优化为长连接复用
var mcpMutex sync.Mutex

// ToolMeta 结构体（本地定义，需与 tools.go 保持一致）
type ToolMeta struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
}

type Client struct {
	cmdPath string // MCP server 可执行文件路径
}

func NewClient(cmdPath string) *Client {
	return &Client{cmdPath: cmdPath}
}

func CallWeather(city, country, lang, unit string) (string, error) {
	fmt.Printf("[CallWeather] city=%s, country=%s, lang=%s, unit=%s\n", city, country, lang, unit)
	mcpMutex.Lock()
	defer mcpMutex.Unlock()

	cmd := exec.Command("go", "run", "./chat/mcp/weather/main.go")
	fmt.Printf("[CallWeather] 启动子进程: %v\n", cmd.Args)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Printf("[CallWeather] 获取 stdin 失败: %v\n", err)
		return "", fmt.Errorf("获取 stdin 失败: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("[CallWeather] 获取 stdout 失败: %v\n", err)
		return "", fmt.Errorf("获取 stdout 失败: %v", err)
	}
	if err := cmd.Start(); err != nil {
		fmt.Printf("[CallWeather] 启动 mcp-server 失败: %v\n", err)
		return "", fmt.Errorf("启动 mcp-server 失败: %v", err)
	}
	// 构造 JSON-RPC 2.0 请求
	type mcpRequest struct {
		JSONRPC string      `json:"jsonrpc"`
		Method  string      `json:"method"`
		Params  interface{} `json:"params"`
		Id      string      `json:"id"`
	}
	req := mcpRequest{
		JSONRPC: "2.0",
		Method:  "get_forecast",
		Params:  WeatherArgs{Location: city},
		Id:      "1",
	}
	fmt.Printf("[CallWeather] 写入请求: %+v\n", req)
	enc := json.NewEncoder(stdin)
	if err := enc.Encode(req); err != nil {
		fmt.Printf("[CallWeather] 写入请求失败: %v\n", err)
		return "", fmt.Errorf("写入请求失败: %v", err)
	}
	stdin.Close()
	// 读取响应
	dec := json.NewDecoder(stdout)
	type mcpResponse struct {
		Content struct {
			Text string `json:"text"`
		} `json:"content"`
		Error interface{} `json:"error"`
	}
	var resp mcpResponse
	fmt.Printf("[CallWeather] 开始读取响应...\n")
	if err := dec.Decode(&resp); err != nil {
		fmt.Printf("[CallWeather] 解析响应失败: %v\n", err)
		return "", fmt.Errorf("解析响应失败: %v", err)
	}
	cmd.Wait()
	fmt.Printf("[CallWeather] 响应内容: %+v\n", resp)
	if resp.Error != nil && resp.Error != "" {
		fmt.Printf("[CallWeather] mcp-server 错误: %v\n", resp.Error)
		return "", fmt.Errorf("mcp-server 错误: %v", resp.Error)
	}
	return resp.Content.Text, nil
}

// ListTools 启动 MCP server，发送 tools/list 请求，返回工具列表
func (c *Client) ListTools(ctx context.Context) ([]ToolMeta, error) {
	mcpMutex.Lock()
	defer mcpMutex.Unlock()

	cmd := exec.CommandContext(ctx, "go", "run", c.cmdPath)
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
	// 构造 tools/list 请求
	type mcpRequest struct {
		JSONRPC string      `json:"jsonrpc"`
		Method  string      `json:"method"`
		Params  interface{} `json:"params"`
		Id      string      `json:"id"`
	}
	req := mcpRequest{
		JSONRPC: "2.0",
		Method:  "tools/list",
		Params:  map[string]interface{}{},
		Id:      "1",
	}
	enc := json.NewEncoder(stdin)
	if err := enc.Encode(req); err != nil {
		return nil, fmt.Errorf("写入请求失败: %v", err)
	}
	stdin.Close()
	// 读取响应
	type toolInfo struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Parameters  map[string]interface{} `json:"parameters"`
	}
	type mcpResponse struct {
		Result struct {
			Tools []toolInfo `json:"tools"`
		} `json:"result"`
		Error interface{} `json:"error"`
	}
	dec := json.NewDecoder(stdout)
	var resp mcpResponse
	if err := dec.Decode(&resp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	cmd.Wait()
	if resp.Error != nil && resp.Error != "" {
		return nil, fmt.Errorf("mcp-server 错误: %v", resp.Error)
	}
	var tools []ToolMeta
	for _, t := range resp.Result.Tools {
		tools = append(tools, ToolMeta{
			Name:        t.Name,
			Description: t.Description,
			Parameters:  t.Parameters,
		})
	}
	return tools, nil
}

// CallTool 启动 MCP server，发送 tools/call 请求，带 toolName 和参数，返回结果
func (c *Client) CallTool(ctx context.Context, toolName string, params map[string]interface{}) (string, error) {
	mcpMutex.Lock()
	defer mcpMutex.Unlock()

	cmd := exec.CommandContext(ctx, "go", "run", c.cmdPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("获取 stdin 失败: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("获取 stdout 失败: %v", err)
	}
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("启动 mcp-server 失败: %v", err)
	}
	// 构造 tools/call 请求
	type mcpRequest struct {
		JSONRPC string      `json:"jsonrpc"`
		Method  string      `json:"method"`
		Params  interface{} `json:"params"`
		Id      string      `json:"id"`
	}
	callParams := map[string]interface{}{
		"tool":   toolName,
		"params": params,
	}
	req := mcpRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  callParams,
		Id:      "1",
	}
	enc := json.NewEncoder(stdin)
	if err := enc.Encode(req); err != nil {
		return "", fmt.Errorf("写入请求失败: %v", err)
	}
	stdin.Close()
	// 读取响应
	type mcpResponse struct {
		Result struct {
			Content struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"result"`
		Error interface{} `json:"error"`
	}
	dec := json.NewDecoder(stdout)
	var resp mcpResponse
	if err := dec.Decode(&resp); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}
	cmd.Wait()
	if resp.Error != nil && resp.Error != "" {
		return "", fmt.Errorf("mcp-server 错误: %v", resp.Error)
	}
	return resp.Result.Content.Text, nil
}
