

>vhagar: 瓦格哈尔，冰与火之歌，龙之家族中最大的一条龙。

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


### 生成镜像

```shell
cd vhagar
# 编译打包，指定tag
docker build -t vhagar:v1.0 .
# 查找镜像
docker images|grep vhagar
# 推送到tcr
docker tag 946289afea65 ka-tcr.tencentcloudcr.com/middleware/alarm-go:v4.2
docker push ka-tcr.tencentcloudcr.com/middleware/alarm-go:v4.2
# 离线镜像
docker save -o alarm-go_v4.2 ka-tcr.tencentcloudcr.com/middleware/alarm-go:v4.2
```