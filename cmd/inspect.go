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
)

// inspectCmd represents the version command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "项目巡检",
	Long:  `获取项目的企业数据，活跃数，会话数`,
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		//case rocketmq:
		//	mqTask()
		default:
			// 执行巡检 job
			tenant := NewTenant(CONFIG)
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
			inspect.TenantTask(tenant, 0)
		}
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().BoolVarP(&rocketmq, "rocketmq", "m", false, "获取 rocketmq broker 信息")

}

func NewTenant(cfg *Config) *inspect.Tenant {
	log.Println("初始化 Tenant 对象")

	tenant := &inspect.Tenant{
		ProjectName: cfg.ProjectName,
		ProxyURL:    cfg.ProxyURL,
		Version:     "v4.5",
		Corp:        cfg.Tenant.Corp,
		Userlist:    cfg.Tenant.Userlist,
		Robotkey:    cfg.Tenant.Robotkey,
	}
	return tenant

}
