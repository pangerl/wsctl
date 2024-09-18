package es

import (
	"context"
	"fmt"
	"vhagar/libs"
)

type ES struct{}

func GetES() *ES {
	return &ES{}
}

func (e *ES) Check() {
	fmt.Println("正在检查 Elasticsearch 集群状态...")

	// 创建 Elasticsearch 客户端
	client, err := libs.GetESClient()
	if err != nil {
		fmt.Printf("创建 Elasticsearch 客户端失败: %s\n", err)
		return
	}

	// 检查集群健康状态
	health, err := client.ClusterHealth().Do(context.Background())
	if err != nil {
		fmt.Printf("获取集群健康状态失败: %s\n", err)
		return
	}
	fmt.Printf("集群健康状态: %s\n", health.Status)

	// 检查节点信息
	stats, err := client.NodesStats().Do(context.Background())
	if err != nil {
		fmt.Printf("获取节点统计信息失败: %s\n", err)
		return
	}

	for _, node := range stats.Nodes {
		fmt.Printf("节点 %s:\n", node.Name)
		fmt.Printf("  JVM 堆使用: %d/%d\n", node.JVM.Mem.HeapUsed, node.JVM.Mem.HeapMax)
		fmt.Printf("  磁盘使用: %d/%d\n", node.FS.Total.TotalInBytes-node.FS.Total.AvailableInBytes, node.FS.Total.TotalInBytes)
	}
}
