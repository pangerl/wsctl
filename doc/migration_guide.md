# 配置库重构迁移指南

本文档提供了从旧的 `vhagar/libs` 包迁移到新的模块化包结构的详细指南。

## 概述

为了提高代码的可维护性和模块化程度，我们对项目的核心库进行了重构，将原来集中在 `libs` 包中的功能拆分到了多个专门的包中：

- `config`: 配置管理
- `database`: 数据库连接工具
- `errors`: 错误处理
- `logger`: 日志系统
- `models`: 业务模型
- `utils`: 通用工具函数

## 导入路径变更

### 旧的导入方式

```go
import "vhagar/libs"
```

### 新的导入方式

根据需要导入特定的包：

```go
import (
    "vhagar/config"    // 配置管理
    "vhagar/database"  // 数据库连接
    "vhagar/errors"    // 错误处理
    "vhagar/logger"    // 日志系统
    "vhagar/models"    // 业务模型
    "vhagar/utils"     // 通用工具函数
)
```

## 主要变更对照表

### 配置相关

| 旧代码 | 新代码 |
|--------|--------|
| `libs.Config` | `config.Config` |
| `libs.InitConfig(path)` | `config.InitConfig(path)` |
| `libs.Config.LogLevel` | `config.Config.Global.LogLevel` |
| `libs.Config.LogToFile` | `config.Config.Global.LogToFile` |
| `libs.Config.Tenant` | `config.Config.Tenant` |
| `libs.Config.DB` | `config.Config.Database.PG` |
| `libs.Config.ES` | `config.Config.Database.ES` |
| `libs.Config.Redis` | `config.Config.Database.Redis` |

### 数据库相关

| 旧代码 | 新代码 |
|--------|--------|
| `libs.DB` | `database.Config` |
| `libs.RedisConfig` | `database.RedisConfig` |
| `libs.NewMysqlClient(cfg)` | `database.NewMySQLClient(cfg, dbName)` |
| `libs.NewPGClient(cfg)` | `database.NewPostgreSQLClient(cfg)` |
| `libs.NewRedisClient(cfg)` | `database.NewRedisClient(cfg)` |
| `libs.NewESClient(cfg)` | `database.NewESClient(cfg)` |
| `libs.PGClienter` | `database.PostgreSQLClient` |

### 错误处理

| 旧代码 | 新代码 |
|--------|--------|
| `libs.ErrorCode` | `errors.ErrorCode` |
| `libs.AppError` | `errors.AppError` |
| `libs.NewError(code, msg)` | `errors.New(code, msg)` |
| `libs.WrapError(code, msg, err)` | `errors.Wrap(code, msg, err)` |
| `libs.LogError(err, ctx)` | `errors.LogError(err, ctx)` |
| `libs.ErrCodeInvalidParam` | `errors.ErrCodeInvalidParam` |
| `libs.ErrInvalidParam` | `errors.ErrInvalidParam` |

### 日志系统

| 旧代码 | 新代码 |
|--------|--------|
| `libs.Logger` | `logger.Logger` |
| `libs.InitLogger(cfg)` | `logger.InitLogger(cfg)` |
| `libs.InitLoggerWithConfig(cfg)` | `logger.InitLogger(cfg)` |

### 业务模型

| 旧代码 | 新代码 |
|--------|--------|
| `config.Corp` | `models.Corp` 或 `config.CorpConfig` |
| `config.Tenant` | `models.Tenant` |
| `config.DorisCfg` | `models.DorisTask` |
| `config.RocketMQCfg` | `models.RocketMQTask` |
| `config.NacosCfg` | `models.NacosTask` |

### 工具函数

| 旧代码 | 新代码 |
|--------|--------|
| `task.GetZeroTime()` | `utils.GetZeroTime()` |
| `task.DoRequest()` | `utils.DoRequest()` |
| `config.GetRandomDuration()` | `utils.GetRandomDuration()` |
| `task.CallUser()` | `utils.CallUser()` |

## 代码示例对比

### 配置加载

**旧代码**:
```go
import "vhagar/libs"

func main() {
    if _, err := libs.InitConfig("config.toml"); err != nil {
        panic(err)
    }
    
    // 访问配置
    dbConfig := libs.Config.DB
    redisConfig := libs.Config.Redis
}
```

**新代码**:
```go
import "vhagar/config"

func main() {
    if _, err := config.InitConfig("config.toml"); err != nil {
        panic(err)
    }
    
    // 访问配置
    dbConfig := config.Config.Database.PG
    redisConfig := config.Config.Database.Redis
}
```

### 数据库连接

**旧代码**:
```go
import "vhagar/libs"

func connectDB() {
    dbConfig := libs.Config.DB
    client, err := libs.NewPGClient(dbConfig)
    if err != nil {
        libs.LogError(err, "连接数据库失败")
        return
    }
    
    // 使用客户端
    defer client.Close()
}
```

**新代码**:
```go
import (
    "vhagar/config"
    "vhagar/database"
    "vhagar/errors"
)

func connectDB() {
    dbConfig := config.Config.Database.PG
    client, err := database.NewPostgreSQLClient(dbConfig)
    if err != nil {
        errors.LogError(err, "连接数据库失败")
        return
    }
    
    // 使用客户端
    defer client.Close()
}
```

### 错误处理

**旧代码**:
```go
import "vhagar/libs"

func doSomething(param string) error {
    if param == "" {
        return libs.NewError(libs.ErrCodeInvalidParam, "参数不能为空")
    }
    
    result, err := callExternalAPI(param)
    if err != nil {
        return libs.WrapError(libs.ErrCodeNetworkFailed, "调用外部API失败", err)
    }
    
    return nil
}
```

**新代码**:
```go
import "vhagar/errors"

func doSomething(param string) error {
    if param == "" {
        return errors.New(errors.ErrCodeInvalidParam, "参数不能为空")
    }
    
    result, err := callExternalAPI(param)
    if err != nil {
        return errors.Wrap(errors.ErrCodeNetworkFailed, "调用外部API失败", err)
    }
    
    return nil
}
```

### 日志记录

**旧代码**:
```go
import "vhagar/libs"

func init() {
    libs.InitLogger(libs.LoggerConfig{
        Level:  "info",
        ToFile: true,
    })
}

func process() {
    libs.Logger.Infow("开始处理", "time", time.Now())
    
    // 处理逻辑
    
    libs.Logger.Infow("处理完成", "duration", time.Since(start))
}
```

**新代码**:
```go
import "vhagar/logger"

func init() {
    logger.InitLogger(logger.Config{
        Level:  "info",
        ToFile: true,
    })
}

func process() {
    logger.Logger.Infow("开始处理", "time", time.Now())
    
    // 处理逻辑
    
    logger.Logger.Infow("处理完成", "duration", time.Since(start))
}
```

## 兼容性说明

为了确保平滑迁移，我们提供了临时的兼容层，允许旧代码继续使用 `vhagar/libs` 导入路径。但这个兼容层将在未来版本中移除，因此建议尽快更新代码以使用新的包结构。

## 最佳实践

1. **按需导入**：只导入你需要的特定包，而不是整个 `libs` 包
2. **使用新的错误处理**：利用 `errors` 包提供的结构化错误处理机制
3. **配置访问**：通过 `config.Config` 访问配置，注意字段路径的变化
4. **数据库连接**：使用 `database` 包中的专门客户端函数
5. **日志记录**：使用 `logger` 包进行日志记录，支持结构化日志

## 常见问题

### Q: 如何快速找到需要更新的代码？

A: 可以使用以下命令查找所有引用 `vhagar/libs` 的文件：

```bash
grep -r "vhagar/libs" --include="*.go" .
```

### Q: 配置结构发生了变化，如何确保我的代码正确访问配置？

A: 查看 `config/config.go` 文件中的 `AppConfig` 结构体定义，了解新的配置结构。大多数配置项已经按功能分组到不同的子结构中。

### Q: 我在哪里可以找到完整的错误码列表？

A: 查看 `errors/errors.go` 文件，其中定义了所有的错误码常量。

## 帮助与支持

如果你在迁移过程中遇到任何问题，请联系项目维护团队获取帮助。