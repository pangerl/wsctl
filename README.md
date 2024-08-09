

>vhagar: 瓦格哈尔

### 使用 cobra-cli 工具

安装工具：`go install github.com/spf13/cobra-cli@latest`

**创建和初始化项目**

```bash
mkdir vhagar
go mod init cobra
# 项目初始化
cobra-cli init
# 增加功能
cobra-cli add version
# 编译
go build -o vhagar
# 执行
./vhagar
```

### 开发

```bash
# 安装第三方库
go get github.com/BurntSushi/toml
go get github.com/tidwall/gjson
go get github.com/olekukonko/tablewriter
go get github.com/gin-gonic/gin
go get github.com/olivere/elastic/v7
go get github.com/jackc/pgx/v5
go get github.com/robfig/cron/v3
```