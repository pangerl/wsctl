
vhagar
======

### wsctl 微盛运维管理工具

特性
------

* 企微租户巡检
* 中间件监控

安装
------
#### 二进制文件

```bash
# 修改配置文件
vim config.toml
# 调试模式，单次运行
[root@localhost vhagar]# ./wsctl inspect
2024/08/19 18:07:41 Info: 读取配置文件 config.toml
2024/08/19 18:07:41 开始项目巡检
2024/08/19 18:07:41 启动企微租户巡检任务
2024/08/19 18:07:44 任务等待时间 0s
2024/08/19 18:07:44 推送企微机器人 response Status:200 OK
2024/08/19 18:07:44 推送企微机器人 response Status:200 OK
# 调度模式，定时任务
[root@localhost vhagar]# ./wsctl crontab
2024/08/19 18:08:40 Info: 读取配置文件 config.toml
2024/08/19 18:08:40 启动任务调度
# 后台运行
nohup ./wsctl crontab > /dev/null 2>&1 &
```

#### 源码编译

```bash
git clone https://github.com/pangerl/vhagar.git
cd vhagar
go build -o wsctl
```
#### Docker
```bash
# 修改配置文件
vim config.toml
# 后台启动
docker-compose up -d
```

快速上手
------

### 生成模板配置

```bash
[root@localhost vhagar]# ./wsctl
2024/09/03 11:22:18 读取配置文件 config.toml
2024/09/03 11:22:18 config.toml 文件不存在，创建模板配置文件
2024/09/03 11:22:18 config.toml.tml：创建成功
2024/09/03 11:22:18 wsctl go go go！！！
```

### 查看帮助

```bash
[root@localhost vhagar]# ./wsctl -h
A longer description that vhagar

Usage:
  wsctl [flags]
  wsctl [command]

Available Commands:
  check       检查服务
  completion  Generate the autocompletion script for the specified shell
  crontab      启动定时任务
  help        Help about any command
  inspect     项目巡检
  metric      监控指标
  nacos       服务健康检查工具
  version     查看版本

Flags:
  -c, --config string   config file (default "config.toml")
  -h, --help            help for wsctl

Use "wsctl [command] --help" for more information about a command.
```

### 项目巡检
