/*
Package models 提供业务模型定义。

本包包含应用程序的核心业务模型，如租户、任务和指标等，这些模型用于表示业务实体和业务逻辑。
每个模型都定义了相应的结构体和方法，用于数据存储、业务处理和展示。

主要模型：

  - 租户模型 (Tenant, Corp): 表示系统中的租户和企业
  - 任务模型 (DorisTask, RocketMQTask, NacosTask): 表示各类任务
  - 指标模型 (MetricData, HostMetric): 表示监控指标和主机状态

租户模型：

租户模型用于管理多租户系统中的租户和企业信息：

	// 创建租户
	tenant := &models.Tenant{
		Corps: []*models.Corp{
			{
				CorpID:      "corp1",
				CorpName:    "企业1",
				ConvEnabled: true,
			},
		},
	}

	// 获取特定企业
	corp := tenant.GetCorpByID("corp1")

	// 获取活跃企业列表
	activeCorps := tenant.GetActiveCorps()

	// 获取统计信息
	totalUsers := tenant.GetTotalUsers()
	totalMessages := tenant.GetTotalMessages()

任务模型：

任务模型表示系统中的各类任务，如 Doris 任务、RocketMQ 任务和 Nacos 任务：

	// 创建 Doris 任务
	task := &models.DorisTask{
		TaskID:    "task1",
		Status:    models.TaskStatusPending,
		CreatedAt: time.Now(),
	}

	// 更新任务状态
	task.UpdateStatus(models.TaskStatusRunning)

	// 检查任务状态
	if task.IsRunning() {
		// 处理运行中的任务
	}

	// 获取任务运行时长
	duration := task.GetDuration()

指标模型：

指标模型用于表示系统监控指标和主机状态：

	// 创建指标数据
	metric := &models.MetricData{
		Timestamp: time.Now(),
		Source:    "host1",
	}

	// 添加指标和标签
	metric.AddMetric("cpu_usage", 75.5)
	metric.AddMetric("memory_usage", 80.2)
	metric.AddTag("environment", "production")

	// 创建主机指标
	hostMetric := &models.HostMetric{
		CPUUsageActive:      75.5,
		MemUsedPercent:      80.2,
		RootDiskUsedPercent: 60.0,
		DataDiskUsedPercent: 70.0,
	}

	// 检查主机健康状态
	if hostMetric.IsHealthy() {
		// 主机状态正常
	}

	// 获取CPU使用率级别
	cpuLevel := hostMetric.GetCPUUsageLevel()

注意事项：

  - 模型结构体应该只包含必要的字段，避免冗余
  - 业务逻辑应该通过模型方法实现，而不是散布在代码各处
  - 模型之间的关系应该清晰，避免循环依赖
*/
package models
