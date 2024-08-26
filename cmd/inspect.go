// Package cmd @Author lanpang
// @Date 2024/8/6 下午4:49:00
// @Desc
package cmd

import (
	"log"
	"vhagar/inspect"
	"vhagar/libs"

	"github.com/spf13/cobra"
)

var (
	rocketmq bool
	doris    bool
)

// inspectCmd represents the version command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "项目巡检",
	Long:  `获取项目的企业数据，活跃数，会话数`,
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化Tenant
		tenant := NewTenant(CONFIG)
		switch {
		case rocketmq:
			tenant.Rocketmq = CONFIG.Rocketmq
			inspect.RocketmqTask(tenant)
		case doris:
			// 创建 mysqlClinet
			mysqlClinet, _ := libs.NewMysqlClient(CONFIG.Doris, "wshoto")
			defer func() {
				if mysqlClinet != nil {
					err := mysqlClinet.Close()
					if err != nil {
						return
					}
				}
			}()
			tenant.MysqlClient = mysqlClinet
		default:
			// 创建ESClient，PGClient
			esClient, _ := libs.NewESClient(CONFIG.ES)
			pgClient, _ := libs.NewPGClient(CONFIG.PG)
			defer func() {
				if pgClient != nil {
					pgClient.Close()
				}
				if esClient != nil {
					esClient.Stop()
				}
			}()
			tenant.ESClient = esClient
			tenant.PGClient = pgClient
			tenant.Corp = CONFIG.Tenant.Corp
			inspect.TenantTask(tenant, 0)
		}
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().BoolVarP(&rocketmq, "rocketmq", "m", false, "获取 rocketmq broker 信息")
	inspectCmd.Flags().BoolVarP(&doris, "doris", "d", false, "检查 doris 服务")

}

func NewTenant(cfg *Config) *inspect.Tenant {
	log.Println("初始化 Tenant 对象")

	tenant := &inspect.Tenant{
		ProjectName: cfg.ProjectName,
		ProxyURL:    cfg.ProxyURL,
		Version:     "v4.5",
		Userlist:    cfg.Tenant.Userlist,
		Robotkey:    cfg.Tenant.Robotkey,
	}
	return tenant

}
