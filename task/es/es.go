package es

import (
	"context"
	"fmt"
	"log"
	"vhagar/config"
	"vhagar/libs"
	"vhagar/task"
)

const taskName = "es"

func init() {
	task.Add(taskName, func() task.Tasker {
		return newES(config.Config)
	})
}

func (es *ES) Gather() {
	esClient, err := libs.NewESClient(config.Config.ES)
	if err != nil {
		log.Printf("Failed info: %s \n", err)
		return
	}
	defer func() {
		if esClient != nil {
			esClient.Stop()
		}
	}()
	es.ESClient = esClient
}

func (es *ES) Check() {
	fmt.Println("正在检查 Elasticsearch 集群状态...")
	client := es.ESClient

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
		fmt.Printf("IP %s:\n", node.IP)
		fmt.Println(node.Process.Mem)
		fmt.Println(node.OS.CPU, node.OS.Mem)
		fmt.Println(node.JVM.Mem)
		fmt.Printf("  HeapUsed: %s", node.JVM.Mem.HeapUsed)
		fmt.Printf("HeapUsedInBytes %d", node.JVM.Mem.HeapUsedInBytes)

		fmt.Printf("  磁盘使用: %d/%d\n", node.FS.Total.TotalInBytes-node.FS.Total.AvailableInBytes, node.FS.Total.TotalInBytes)
	}
}
