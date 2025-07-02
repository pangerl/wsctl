vhagar
======

高效的运维管理与中间件巡检工具

简介
----
`vhagar` 是一款面向运维场景的综合管理工具，集成了网络测试、定时任务调度、中间件健康巡检与监控指标采集等核心能力，助力提升运维自动化与可观测性。

主要特性
--------
- **网络测试**：支持端口连通性检测，快速定位网络问题
- **定时任务**：内置任务调度器，支持周期性任务自动执行
- **中间件巡检**：一键巡检主流中间件（如 Redis、ES、RocketMQ、Doris、Nacos 等）运行状态
- **监控指标**：采集并展示各类服务与中间件的关键指标
- **统一日志**：集成 zap 日志框架，支持多级别、结构化日志输出

安装与部署
----------
### 1. 二进制包运行

1. 下载最新 [wsctl 二进制文件](https://private-1253767630.cos.ap-shanghai.myqcloud.com/tools/archive/binary_tag/bin/wsctl) 和 [配置模板](https://private-1253767630.cos.ap-shanghai.myqcloud.com/tools/archive/binary_tag/bin/config.toml)
2. 修改配置文件 `config.toml`
3. 赋予执行权限并启动
   ```bash
   chmod +x wsctl
   ./wsctl
   ```

### 2. Docker 部署

1. 配置好本地 `config.toml`
2. 使用如下 docker-compose 配置启动
   ```yaml
   version: "3.8"
   services:
     vhagar:
       image: ka-tcr.tencentcloudcr.com/middleware/vhagar:v1.0
       container_name: vhagar
       ports:
         - "8089:8089"
       volumes:
         - ./config.toml:/app/config.toml
       restart: unless-stopped
   ```
3. 启动服务
   ```bash
   docker-compose up -d
   ```

快速上手
--------
### 启动 Web 服务

默认监听端口为 8099，可通过 `-p` 参数自定义端口：
```bash
./wsctl -p 8888
```
访问 `http://<服务器IP>:8888/` 进入管理界面。

### 查看命令帮助

```bash
./wsctl -h
```

常用命令
--------
- `wsctl crontab`：启动定时任务调度器
- `wsctl metric`：采集并展示监控指标
- `wsctl task`：执行服务巡检任务
- `wsctl version`：查看版本信息

更多命令及参数请通过 `-h` 或 `--help` 查看详细说明。

日志系统
--------
本项目集成 [zap](https://github.com/uber-go/zap) 作为全局日志框架，所有日志输出均通过 zap 统一管理。

- 日志初始化：程序启动时自动完成（见 `main.go`、`libs/logger.go`）
- 推荐调用方式：
  ```go
  libs.Logger.Infow("启动服务", "port", 8080)
  libs.Logger.Errorw("数据库连接失败", "err", err)
  ```
- 日志级别、格式、输出位置可在 `libs/logger.go` 中自定义

调试建议
--------
- 默认开发模式（彩色、详细调用栈），如需生产环境可将 `zap.NewDevelopment()` 改为 `zap.NewProduction()`
- 日志输出可扩展到文件、json 格式等，详见 zap 官方文档

目录结构
--------
```
cmd/         # 命令行入口
config/      # 配置文件与结构体
libs/        # 公共库（日志、数据库、缓存等）
metric/      # 监控指标采集
notify/      # 通知模块
task/        # 各类巡检任务
main.go      # 程序主入口
config.toml  # 配置文件
```

贡献与反馈
----------
如有建议或问题，欢迎提 issue 或 PR！
