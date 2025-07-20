# Design Document

## Overview

本设计文档描述了对 vhagar 项目中 config 和 libs 目录的重构方案。通过分析当前代码结构和依赖关系，我们将按照 Go 语言最佳实践重新组织代码，提高项目的可维护性和可读性。

当前问题：
- `libs/` 目录命名不符合 Go 标准实践
- `config/` 目录混合了配置结构体和业务模型
- 结构体定义分散，职责不清晰
- 工具函数缺乏合理的分类组织

## Architecture

### 新的目录结构

```
database/                     # 数据库连接工具
├── mysql.go                 # MySQL 连接工具
├── postgres.go              # PostgreSQL 连接工具  
├── redis.go                 # Redis 连接工具
├── elasticsearch.go         # Elasticsearch 连接工具
└── config.go                # 数据库配置结构体

logger/                       # 日志工具
└── logger.go                # 日志初始化和配置

errors/                       # 错误处理
└── errors.go                # 统一错误处理

utils/                        # 通用工具函数
├── time.go                  # 时间相关工具
├── http.go                  # HTTP 请求工具
└── random.go                # 随机数工具

models/                       # 业务模型
├── tenant.go                # 租户相关模型
└── metric.go                # 指标相关模型

config/                       # 配置管理（重构后）
├── config.go                # 主配置结构体和加载逻辑
├── service.go               # 服务配置（AI, Weather, RocketMQ等）
└── loader.go                # 配置文件加载器
```

### 设计原则

1. **单一职责原则**：每个包只负责一个特定的功能域
2. **依赖倒置**：高层模块不依赖低层模块，都依赖抽象
3. **接口隔离**：定义清晰的接口边界
4. **开闭原则**：对扩展开放，对修改封闭

## Components and Interfaces

### 1. 数据库连接组件 (database/)

```go
// database/config.go
package database

// Config 数据库连接配置
type Config struct {
    Host     string `toml:"host"`
    Port     int    `toml:"port"`
    Username string `toml:"username"`
    Password string `toml:"password"`
    Database string `toml:"database"`
    SSLMode  bool   `toml:"ssl_mode"`
}

// HasValue 检查配置是否完整
func (c Config) HasValue() bool {
    return c.Host != "" && c.Port != 0 && c.Username != "" && c.Password != ""
}

// RedisConfig Redis 连接配置
type RedisConfig struct {
    Addr     string `toml:"addr"`
    Password string `toml:"password"`
    DB       int    `toml:"db"`
}
```

```go
// database/mysql.go
package database

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

// MySQLClient MySQL 客户端接口
type MySQLClient interface {
    Connect() (*sql.DB, error)
    Close() error
}

// NewMySQLClient 创建 MySQL 客户端
func NewMySQLClient(cfg Config, dbName string) (*sql.DB, error) {
    // 实现逻辑
}
```

### 2. 日志组件 (logger/)

```go
// logger/logger.go
package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

// Config 日志配置
type Config struct {
    Level  string `toml:"level"`
    ToFile bool   `toml:"to_file"`
}

// Logger 全局日志实例
var Logger *zap.SugaredLogger

// InitLogger 初始化日志器
func InitLogger(cfg Config) error {
    // 实现逻辑
}
```

### 3. 错误处理组件 (errors/)

```go
// errors/errors.go
package errors

import "fmt"

// ErrorCode 错误码类型
type ErrorCode int

// AppError 应用错误结构
type AppError struct {
    Code    ErrorCode `json:"code"`
    Message string    `json:"message"`
    Detail  string    `json:"detail,omitempty"`
    Cause   error     `json:"-"`
}

// Error 实现 error 接口
func (e *AppError) Error() string {
    if e.Detail != "" {
        return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Detail)
    }
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 支持错误链
func (e *AppError) Unwrap() error {
    return e.Cause
}

// 预定义错误码和错误实例
const (
    ErrCodeSuccess         ErrorCode = 0
    ErrCodeInternalErr     ErrorCode = 10001
    ErrCodeConfigNotFound  ErrorCode = 10002
    ErrCodeConfigInvalid   ErrorCode = 10003
    ErrCodeDBConnFailed    ErrorCode = 10004
    ErrCodeNetworkTimeout  ErrorCode = 10005
    ErrCodeAICallFailed    ErrorCode = 10006
)

// 预定义错误实例
var (
    ErrConfigNotFound = &AppError{Code: ErrCodeConfigNotFound, Message: "配置文件未找到"}
    ErrConfigInvalid  = &AppError{Code: ErrCodeConfigInvalid, Message: "配置文件格式错误"}
    ErrDBConnFailed   = &AppError{Code: ErrCodeDBConnFailed, Message: "数据库连接失败"}
    ErrNetworkTimeout = &AppError{Code: ErrCodeNetworkTimeout, Message: "网络请求超时"}
    ErrAICallFailed   = &AppError{Code: ErrCodeAICallFailed, Message: "AI 服务调用失败"}
)

// New 创建新的应用错误
func New(code ErrorCode, message string) *AppError {
    return &AppError{Code: code, Message: message}
}

// Wrap 包装现有错误
func Wrap(err error, code ErrorCode, message string) *AppError {
    return &AppError{Code: code, Message: message, Cause: err}
}
```

### 4. 配置管理组件 (config/)

```go
// config/config.go
package config

import (
    "time"
    "vhagar/database"
    "vhagar/logger"
)

// Config 主配置结构
type Config struct {
    Global   GlobalConfig           `toml:"global"`
    Database DatabaseConfigs        `toml:"database"`
    Services ServiceConfigs         `toml:"services"`
    Logger   logger.Config          `toml:"logger"`
}

// GlobalConfig 全局配置
type GlobalConfig struct {
    ProjectName string        `toml:"project_name"`
    ProxyURL    string        `toml:"proxy_url"`
    Watch       bool          `toml:"watch"`
    Report      bool          `toml:"report"`
    Interval    time.Duration `toml:"interval"`
    Duration    time.Duration `toml:"duration"`
}

// DatabaseConfigs 数据库配置集合
type DatabaseConfigs struct {
    PostgreSQL    database.Config      `toml:"postgresql"`
    MySQL         database.Config      `toml:"mysql"`
    Redis         database.RedisConfig `toml:"redis"`
    Elasticsearch database.Config      `toml:"elasticsearch"`
    Doris         DorisConfig          `toml:"doris"`
}

// ServiceConfigs 服务配置集合
type ServiceConfigs struct {
    AI       AIConfig       `toml:"ai"`
    Weather  WeatherConfig  `toml:"weather"`
    RocketMQ RocketMQConfig `toml:"rocketmq"`
    Nacos    NacosConfig    `toml:"nacos"`
}

// AIConfig AI 服务配置
type AIConfig struct {
    Provider string `toml:"provider"`
    APIKey   string `toml:"api_key"`
    BaseURL  string `toml:"base_url"`
    Model    string `toml:"model"`
}

// WeatherConfig 天气服务配置
type WeatherConfig struct {
    APIKey  string `toml:"api_key"`
    BaseURL string `toml:"base_url"`
}

// RocketMQConfig RocketMQ 配置
type RocketMQConfig struct {
    NameServer string `toml:"name_server"`
    Topic      string `toml:"topic"`
    Group      string `toml:"group"`
}

// NacosConfig Nacos 配置
type NacosConfig struct {
    ServerAddr string `toml:"server_addr"`
    Namespace  string `toml:"namespace"`
    Group      string `toml:"group"`
}

// DorisConfig Doris 数据库配置
type DorisConfig struct {
    Host     string `toml:"host"`
    Port     int    `toml:"port"`
    Username string `toml:"username"`
    Password string `toml:"password"`
    Database string `toml:"database"`
}
```

### 5. 业务模型组件 (models/)

```go
// models/tenant.go
package models

// Tenant 租户模型
type Tenant struct {
    Corps []*Corp `json:"corps"`
}

// Corp 企业模型
type Corp struct {
    CorpID               string `json:"corp_id"`
    ConvEnabled          bool   `json:"conv_enabled"`
    CorpName             string `json:"corp_name"`
    MessageNum           int64  `json:"message_num"`
    YesterdayMessageNum  int64  `json:"yesterday_message_num"`
    UserNum              int    `json:"user_num"`
    CustomerNum          int64  `json:"customer_num"`
    CustomerGroupNum     int    `json:"customer_group_num"`
    CustomerGroupUserNum int    `json:"customer_group_user_num"`
    DAUNum               int    `json:"dau_num"`
    WAUNum               int    `json:"wau_num"`
    MAUNum               int    `json:"mau_num"`
}
```

```go
// models/metric.go
package models

import "time"

// MetricData 指标数据模型
type MetricData struct {
    Timestamp time.Time              `json:"timestamp"`
    Metrics   map[string]interface{} `json:"metrics"`
    Tags      map[string]string      `json:"tags"`
}

// TaskStatus 任务状态模型
type TaskStatus struct {
    TaskID    string    `json:"task_id"`
    Status    string    `json:"status"`
    StartTime time.Time `json:"start_time"`
    EndTime   time.Time `json:"end_time"`
    Error     string    `json:"error,omitempty"`
}
```

### 6. 工具函数组件 (utils/)

```go
// utils/time.go
package utils

import (
    "time"
)

// TimeFormat 常用时间格式
const (
    DateTimeFormat = "2006-01-02 15:04:05"
    DateFormat     = "2006-01-02"
    TimeFormat     = "15:04:05"
)

// FormatTime 格式化时间
func FormatTime(t time.Time, format string) string {
    return t.Format(format)
}

// ParseTime 解析时间字符串
func ParseTime(timeStr, format string) (time.Time, error) {
    return time.Parse(format, timeStr)
}

// GetTodayRange 获取今天的时间范围
func GetTodayRange() (time.Time, time.Time) {
    now := time.Now()
    start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
    end := start.Add(24 * time.Hour).Add(-time.Nanosecond)
    return start, end
}
```

```go
// utils/http.go
package utils

import (
    "bytes"
    "context"
    "encoding/json"
    "io"
    "net/http"
    "time"
)

// HTTPClient HTTP 客户端配置
type HTTPClient struct {
    Client  *http.Client
    Timeout time.Duration
    Retries int
}

// NewHTTPClient 创建新的 HTTP 客户端
func NewHTTPClient(timeout time.Duration, retries int) *HTTPClient {
    return &HTTPClient{
        Client: &http.Client{
            Timeout: timeout,
        },
        Timeout: timeout,
        Retries: retries,
    }
}

// Get 发送 GET 请求
func (c *HTTPClient) Get(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    for k, v := range headers {
        req.Header.Set(k, v)
    }
    
    return c.doRequest(req)
}

// Post 发送 POST 请求
func (c *HTTPClient) Post(ctx context.Context, url string, data interface{}, headers map[string]string) ([]byte, error) {
    var body io.Reader
    if data != nil {
        jsonData, err := json.Marshal(data)
        if err != nil {
            return nil, err
        }
        body = bytes.NewBuffer(jsonData)
    }
    
    req, err := http.NewRequestWithContext(ctx, "POST", url, body)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    for k, v := range headers {
        req.Header.Set(k, v)
    }
    
    return c.doRequest(req)
}

// doRequest 执行请求（带重试机制）
func (c *HTTPClient) doRequest(req *http.Request) ([]byte, error) {
    var lastErr error
    for i := 0; i <= c.Retries; i++ {
        resp, err := c.Client.Do(req)
        if err != nil {
            lastErr = err
            continue
        }
        defer resp.Body.Close()
        
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            lastErr = err
            continue
        }
        
        if resp.StatusCode >= 200 && resp.StatusCode < 300 {
            return body, nil
        }
        
        lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
    }
    return nil, lastErr
}
```

## Data Models

### 配置数据模型

1. **主配置模型**：整合所有配置项，包括全局配置、数据库配置、服务配置和日志配置
2. **数据库配置模型**：统一的数据库连接配置，支持 MySQL、PostgreSQL、Redis、Elasticsearch 和 Doris
3. **服务配置模型**：各种外部服务配置，包括 AI 服务、天气服务、RocketMQ 和 Nacos
4. **日志配置模型**：日志级别和输出方式配置

### 业务数据模型

1. **租户模型**：企业和用户相关数据，包含企业信息、用户统计、消息统计等
2. **指标模型**：监控和统计数据，包含各种业务指标和性能指标
3. **任务模型**：定时任务相关数据，包含任务配置、执行状态等

### 工具函数模型

1. **时间工具**：时间格式化、解析、计算等功能
2. **HTTP 工具**：HTTP 请求封装、重试机制、超时处理
3. **随机数工具**：各种随机数生成、UUID 生成等

## Error Handling

### 错误分类

1. **配置错误**：配置文件不存在、格式错误等
2. **网络错误**：HTTP 请求失败、超时等
3. **数据库错误**：连接失败、查询错误等
4. **业务错误**：AI 调用失败、工具调用失败等

### 错误处理策略

1. **统一错误码**：使用预定义的错误码
2. **错误包装**：保留原始错误信息
3. **结构化日志**：记录详细的错误上下文
4. **优雅降级**：关键错误不影响其他功能

## Testing Strategy

### 单元测试

1. **配置加载测试**：测试各种配置文件格式
2. **数据库连接测试**：测试各种数据库连接场景
3. **工具函数测试**：测试边界条件和异常情况
4. **错误处理测试**：测试错误码和错误信息

### 集成测试

1. **配置集成测试**：测试完整的配置加载流程
2. **数据库集成测试**：测试实际的数据库连接
3. **日志集成测试**：测试日志输出和文件写入

### 测试文件组织

```
database/
├── mysql.go
├── mysql_test.go
├── postgres.go
├── postgres_test.go
├── redis.go
├── redis_test.go
├── config.go
└── config_test.go

logger/
├── logger.go
└── logger_test.go

errors/
├── errors.go
└── errors_test.go

utils/
├── time.go
├── time_test.go
├── http.go
└── http_test.go

config/
├── config.go
├── config_test.go
├── loader.go
└── loader_test.go

models/
├── tenant.go
├── tenant_test.go
├── metric.go
└── metric_test.go
```

## Migration Plan

### 阶段 1：创建新的包结构
1. 创建 `database/`、`logger/`、`errors/`、`utils/`、`models/` 目录
2. 移动和重构代码文件到对应目录
3. 重构 `config/` 目录结构

### 阶段 2：更新导入路径
1. 批量更新所有 `vhagar/libs` 导入
2. 更新配置相关的导入路径
3. 确保编译通过

### 阶段 3：测试和验证
1. 运行所有现有测试
2. 添加新的单元测试
3. 验证应用程序功能

### 阶段 4：清理和文档
1. 删除旧的 `libs/` 目录
2. 更新 README 和文档
3. 添加包级别的文档注释
4. 代码注释用中文

## Backward Compatibility

为了保持向后兼容性，我们将：

1. **保留原有的 config 包入口**：在 `config/config.go` 中重新导出新的配置结构
2. **提供别名导入**：为常用的类型和函数提供别名
3. **渐进式迁移**：允许新旧代码共存一段时间
4. **清晰的迁移指南**：提供详细的迁移步骤和示例

```go
// config/config.go - 向后兼容入口
package config

import (
    "vhagar/database"
    "vhagar/logger"
    "vhagar/models"
)

// 重新导出主要类型（向后兼容）
type DB = database.Config
type RedisConfig = database.RedisConfig

// 重新导出业务模型
type Tenant = models.Tenant
type Corp = models.Corp

// 全局配置实例
var GlobalConfig *Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
    // 实现配置加载逻辑
    return nil, nil
}
```