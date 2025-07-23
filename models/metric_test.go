package models

import (
	"testing"
	"time"
)

// TestMetricData_AddMetric 测试向指标数据中添加单个指标
func TestMetricData_AddMetric(t *testing.T) {
	// 创建测试数据
	md := &MetricData{
		Timestamp: time.Now(),
		Source:    "test",
	}

	// 添加指标
	md.AddMetric("cpu_usage", 75.5)
	md.AddMetric("memory_usage", 80.2)
	md.AddMetric("disk_usage", "high")

	// 验证指标是否被添加
	if len(md.Metrics) != 3 {
		t.Errorf("AddMetric() added %d metrics, want %d", len(md.Metrics), 3)
	}

	// 验证指标值
	if value, ok := md.Metrics["cpu_usage"]; !ok || value != 75.5 {
		t.Errorf("AddMetric() cpu_usage = %v, want %v", value, 75.5)
	}
	if value, ok := md.Metrics["memory_usage"]; !ok || value != 80.2 {
		t.Errorf("AddMetric() memory_usage = %v, want %v", value, 80.2)
	}
	if value, ok := md.Metrics["disk_usage"]; !ok || value != "high" {
		t.Errorf("AddMetric() disk_usage = %v, want %v", value, "high")
	}
}

// TestMetricData_AddTag 测试向指标数据中添加标签
func TestMetricData_AddTag(t *testing.T) {
	// 创建测试数据
	md := &MetricData{
		Timestamp: time.Now(),
		Source:    "test",
	}

	// 添加标签
	md.AddTag("host", "server1")
	md.AddTag("environment", "production")
	md.AddTag("region", "us-west")

	// 验证标签是否被添加
	if len(md.Tags) != 3 {
		t.Errorf("AddTag() added %d tags, want %d", len(md.Tags), 3)
	}

	// 验证标签值
	if value, ok := md.Tags["host"]; !ok || value != "server1" {
		t.Errorf("AddTag() host = %v, want %v", value, "server1")
	}
	if value, ok := md.Tags["environment"]; !ok || value != "production" {
		t.Errorf("AddTag() environment = %v, want %v", value, "production")
	}
	if value, ok := md.Tags["region"]; !ok || value != "us-west" {
		t.Errorf("AddTag() region = %v, want %v", value, "us-west")
	}
}

// TestMetricData_GetMetric 测试获取指定名称的指标值
func TestMetricData_GetMetric(t *testing.T) {
	// 创建测试数据
	md := &MetricData{
		Timestamp: time.Now(),
		Source:    "test",
		Metrics: map[string]interface{}{
			"cpu_usage":    75.5,
			"memory_usage": 80.2,
		},
	}

	// 测试获取存在的指标
	value, exists := md.GetMetric("cpu_usage")
	if !exists {
		t.Error("GetMetric() returned false for existing metric")
	}
	if value != 75.5 {
		t.Errorf("GetMetric() = %v, want %v", value, 75.5)
	}

	// 测试获取不存在的指标
	value, exists = md.GetMetric("disk_usage")
	if exists {
		t.Errorf("GetMetric() returned true for non-existing metric: %v", value)
	}

	// 测试 nil Metrics
	md.Metrics = nil
	value, exists = md.GetMetric("cpu_usage")
	if exists {
		t.Errorf("GetMetric() returned true for nil Metrics: %v", value)
	}
}

// TestMetricData_GetTag 测试获取指定标签值
func TestMetricData_GetTag(t *testing.T) {
	// 创建测试数据
	md := &MetricData{
		Timestamp: time.Now(),
		Source:    "test",
		Tags: map[string]string{
			"host":        "server1",
			"environment": "production",
		},
	}

	// 测试获取存在的标签
	value, exists := md.GetTag("host")
	if !exists {
		t.Error("GetTag() returned false for existing tag")
	}
	if value != "server1" {
		t.Errorf("GetTag() = %v, want %v", value, "server1")
	}

	// 测试获取不存在的标签
	value, exists = md.GetTag("region")
	if exists {
		t.Errorf("GetTag() returned true for non-existing tag: %v", value)
	}

	// 测试 nil Tags
	md.Tags = nil
	value, exists = md.GetTag("host")
	if exists {
		t.Errorf("GetTag() returned true for nil Tags: %v", value)
	}
}

// TestHostMetric_IsHealthy 测试主机指标是否健康
func TestHostMetric_IsHealthy(t *testing.T) {
	// 创建测试数据
	tests := []struct {
		name                string
		cpuUsage            float64
		memUsage            float64
		rootDiskUsage       float64
		dataDiskUsage       float64
		expectedHealthState bool
	}{
		{
			name:                "健康状态",
			cpuUsage:            50.0,
			memUsage:            60.0,
			rootDiskUsage:       70.0,
			dataDiskUsage:       70.0,
			expectedHealthState: true,
		},
		{
			name:                "CPU使用率过高",
			cpuUsage:            85.0,
			memUsage:            60.0,
			rootDiskUsage:       70.0,
			dataDiskUsage:       70.0,
			expectedHealthState: false,
		},
		{
			name:                "内存使用率过高",
			cpuUsage:            50.0,
			memUsage:            95.0,
			rootDiskUsage:       70.0,
			dataDiskUsage:       70.0,
			expectedHealthState: false,
		},
		{
			name:                "根磁盘使用率过高",
			cpuUsage:            50.0,
			memUsage:            60.0,
			rootDiskUsage:       90.0,
			dataDiskUsage:       70.0,
			expectedHealthState: false,
		},
		{
			name:                "数据磁盘使用率过高",
			cpuUsage:            50.0,
			memUsage:            60.0,
			rootDiskUsage:       70.0,
			dataDiskUsage:       90.0,
			expectedHealthState: false,
		},
		{
			name:                "多项指标过高",
			cpuUsage:            85.0,
			memUsage:            95.0,
			rootDiskUsage:       90.0,
			dataDiskUsage:       90.0,
			expectedHealthState: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hm := &HostMetric{
				CPUUsageActive:      tt.cpuUsage,
				MemUsedPercent:      tt.memUsage,
				RootDiskUsedPercent: tt.rootDiskUsage,
				DataDiskUsedPercent: tt.dataDiskUsage,
			}

			isHealthy := hm.IsHealthy()
			if isHealthy != tt.expectedHealthState {
				t.Errorf("IsHealthy() = %v, want %v", isHealthy, tt.expectedHealthState)
			}
		})
	}
}

// TestHostMetric_GetCPUUsageLevel 测试获取CPU使用率等级
func TestHostMetric_GetCPUUsageLevel(t *testing.T) {
	// 创建测试数据
	tests := []struct {
		name          string
		cpuUsage      float64
		expectedLevel string
	}{
		{
			name:          "低使用率",
			cpuUsage:      30.0,
			expectedLevel: "低",
		},
		{
			name:          "中使用率",
			cpuUsage:      60.0,
			expectedLevel: "中",
		},
		{
			name:          "高使用率",
			cpuUsage:      90.0,
			expectedLevel: "高",
		},
		{
			name:          "边界值-低",
			cpuUsage:      49.9,
			expectedLevel: "低",
		},
		{
			name:          "边界值-中",
			cpuUsage:      50.0,
			expectedLevel: "中",
		},
		{
			name:          "边界值-高",
			cpuUsage:      80.0,
			expectedLevel: "高",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hm := &HostMetric{
				CPUUsageActive: tt.cpuUsage,
			}

			level := hm.GetCPUUsageLevel()
			if level != tt.expectedLevel {
				t.Errorf("GetCPUUsageLevel() = %v, want %v", level, tt.expectedLevel)
			}
		})
	}
}

// TestHostMetric_GetMemoryUsageLevel 测试获取内存使用率等级
func TestHostMetric_GetMemoryUsageLevel(t *testing.T) {
	// 创建测试数据
	tests := []struct {
		name          string
		memUsage      float64
		expectedLevel string
	}{
		{
			name:          "低使用率",
			memUsage:      40.0,
			expectedLevel: "低",
		},
		{
			name:          "中使用率",
			memUsage:      70.0,
			expectedLevel: "中",
		},
		{
			name:          "高使用率",
			memUsage:      90.0,
			expectedLevel: "高",
		},
		{
			name:          "边界值-低",
			memUsage:      59.9,
			expectedLevel: "低",
		},
		{
			name:          "边界值-中",
			memUsage:      60.0,
			expectedLevel: "中",
		},
		{
			name:          "边界值-高",
			memUsage:      85.0,
			expectedLevel: "高",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hm := &HostMetric{
				MemUsedPercent: tt.memUsage,
			}

			level := hm.GetMemoryUsageLevel()
			if level != tt.expectedLevel {
				t.Errorf("GetMemoryUsageLevel() = %v, want %v", level, tt.expectedLevel)
			}
		})
	}
}

// TestMetricCollection_AddMetric 测试向集合中添加指标
func TestMetricCollection_AddMetric(t *testing.T) {
	// 创建测试数据
	mc := &MetricCollection{
		CollectionID: "test-collection",
		Source:       "test",
	}

	// 创建指标
	now := time.Now()
	metric1 := &MetricData{
		Timestamp: now.Add(-10 * time.Minute),
		Source:    "test",
		Metrics: map[string]interface{}{
			"cpu_usage": 50.0,
		},
	}
	metric2 := &MetricData{
		Timestamp: now,
		Source:    "test",
		Metrics: map[string]interface{}{
			"memory_usage": 60.0,
		},
	}
	metric3 := &MetricData{
		Timestamp: now.Add(10 * time.Minute),
		Source:    "test",
		Metrics: map[string]interface{}{
			"disk_usage": 70.0,
		},
	}

	// 添加指标
	mc.AddMetric(metric1)
	mc.AddMetric(metric2)
	mc.AddMetric(metric3)

	// 验证指标数量
	if len(mc.Metrics) != 3 {
		t.Errorf("AddMetric() added %d metrics, want %d", len(mc.Metrics), 3)
	}

	// 验证时间范围
	if !mc.StartTime.Equal(metric1.Timestamp) {
		t.Errorf("AddMetric() StartTime = %v, want %v", mc.StartTime, metric1.Timestamp)
	}
	if !mc.EndTime.Equal(metric3.Timestamp) {
		t.Errorf("AddMetric() EndTime = %v, want %v", mc.EndTime, metric3.Timestamp)
	}
}

// TestMetricCollection_GetMetricCount 测试获取指标数量
func TestMetricCollection_GetMetricCount(t *testing.T) {
	// 创建测试数据
	mc := &MetricCollection{
		CollectionID: "test-collection",
		Source:       "test",
		Metrics: []*MetricData{
			{
				Timestamp: time.Now(),
				Source:    "test",
				Metrics: map[string]interface{}{
					"cpu_usage": 50.0,
				},
			},
			{
				Timestamp: time.Now(),
				Source:    "test",
				Metrics: map[string]interface{}{
					"memory_usage": 60.0,
				},
			},
		},
	}

	// 获取指标数量
	count := mc.GetMetricCount()
	if count != 2 {
		t.Errorf("GetMetricCount() = %d, want %d", count, 2)
	}

	// 测试空集合
	mc.Metrics = nil
	count = mc.GetMetricCount()
	if count != 0 {
		t.Errorf("GetMetricCount() = %d, want %d", count, 0)
	}
}

// TestMetricCollection_GetDuration 测试获取采集时间跨度
func TestMetricCollection_GetDuration(t *testing.T) {
	// 创建测试数据
	now := time.Now()
	mc := &MetricCollection{
		CollectionID: "test-collection",
		Source:       "test",
		StartTime:    now.Add(-30 * time.Minute),
		EndTime:      now,
	}

	// 获取时间跨度
	duration := mc.GetDuration()
	expectedDuration := 30 * time.Minute
	if duration != expectedDuration {
		t.Errorf("GetDuration() = %v, want %v", duration, expectedDuration)
	}

	// 测试零时间
	mc.StartTime = time.Time{}
	duration = mc.GetDuration()
	if duration != 0 {
		t.Errorf("GetDuration() = %v, want %v", duration, 0)
	}
}

// TestAlertRule_IsEnabled 测试检查告警规则是否启用
func TestAlertRule_IsEnabled(t *testing.T) {
	// 创建测试数据
	enabledRule := &AlertRule{
		RuleID:  "rule1",
		Name:    "CPU告警",
		Enabled: true,
	}
	disabledRule := &AlertRule{
		RuleID:  "rule2",
		Name:    "内存告警",
		Enabled: false,
	}

	// 测试启用的规则
	if !enabledRule.IsEnabled() {
		t.Errorf("IsEnabled() = false for enabled rule")
	}

	// 测试禁用的规则
	if disabledRule.IsEnabled() {
		t.Errorf("IsEnabled() = true for disabled rule")
	}
}

// TestAlertRule_UpdateThreshold 测试更新告警阈值
func TestAlertRule_UpdateThreshold(t *testing.T) {
	// 创建测试数据
	rule := &AlertRule{
		RuleID:    "rule1",
		Name:      "CPU告警",
		Threshold: 80.0,
		UpdatedAt: time.Now().Add(-24 * time.Hour), // 设置为一天前
	}

	// 记录原始更新时间
	originalUpdatedAt := rule.UpdatedAt

	// 更新阈值
	newThreshold := 90.0
	rule.UpdateThreshold(newThreshold)

	// 验证阈值是否更新
	if rule.Threshold != newThreshold {
		t.Errorf("UpdateThreshold() Threshold = %v, want %v", rule.Threshold, newThreshold)
	}

	// 验证更新时间是否更新
	if !rule.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("UpdateThreshold() did not update UpdatedAt")
	}
}

// TestAlertRule_Enable 测试启用告警规则
func TestAlertRule_Enable(t *testing.T) {
	// 创建测试数据
	rule := &AlertRule{
		RuleID:    "rule1",
		Name:      "CPU告警",
		Enabled:   false,
		UpdatedAt: time.Now().Add(-24 * time.Hour), // 设置为一天前
	}

	// 记录原始更新时间
	originalUpdatedAt := rule.UpdatedAt

	// 启用规则
	rule.Enable()

	// 验证是否启用
	if !rule.Enabled {
		t.Errorf("Enable() did not enable the rule")
	}

	// 验证更新时间是否更新
	if !rule.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("Enable() did not update UpdatedAt")
	}
}

// TestAlertRule_Disable 测试禁用告警规则
func TestAlertRule_Disable(t *testing.T) {
	// 创建测试数据
	rule := &AlertRule{
		RuleID:    "rule1",
		Name:      "CPU告警",
		Enabled:   true,
		UpdatedAt: time.Now().Add(-24 * time.Hour), // 设置为一天前
	}

	// 记录原始更新时间
	originalUpdatedAt := rule.UpdatedAt

	// 禁用规则
	rule.Disable()

	// 验证是否禁用
	if rule.Enabled {
		t.Errorf("Disable() did not disable the rule")
	}

	// 验证更新时间是否更新
	if !rule.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("Disable() did not update UpdatedAt")
	}
}

// TestMetricQuery_SetTimeRange 测试设置查询时间范围
func TestMetricQuery_SetTimeRange(t *testing.T) {
	// 创建测试数据
	query := &MetricQuery{
		QueryID: "query1",
		Query:   "cpu_usage > 80",
	}

	// 设置时间范围
	start := time.Now().Add(-1 * time.Hour)
	end := time.Now()
	query.SetTimeRange(start, end)

	// 验证时间范围是否设置正确
	if !query.StartTime.Equal(start) {
		t.Errorf("SetTimeRange() StartTime = %v, want %v", query.StartTime, start)
	}
	if !query.EndTime.Equal(end) {
		t.Errorf("SetTimeRange() EndTime = %v, want %v", query.EndTime, end)
	}
}

// TestMetricQuery_AddTag 测试向查询中添加标签过滤条件
func TestMetricQuery_AddTag(t *testing.T) {
	// 创建测试数据
	query := &MetricQuery{
		QueryID: "query1",
		Query:   "cpu_usage > 80",
	}

	// 添加标签
	query.AddTag("host", "server1")
	query.AddTag("environment", "production")

	// 验证标签是否添加成功
	if len(query.Tags) != 2 {
		t.Errorf("AddTag() added %d tags, want %d", len(query.Tags), 2)
	}
	if value, ok := query.Tags["host"]; !ok || value != "server1" {
		t.Errorf("AddTag() host = %v, want %v", value, "server1")
	}
	if value, ok := query.Tags["environment"]; !ok || value != "production" {
		t.Errorf("AddTag() environment = %v, want %v", value, "production")
	}
}

// TestMetricQuery_GetDuration 测试获取查询时间跨度
func TestMetricQuery_GetDuration(t *testing.T) {
	// 创建测试数据
	now := time.Now()
	query := &MetricQuery{
		QueryID:   "query1",
		Query:     "cpu_usage > 80",
		StartTime: now.Add(-30 * time.Minute),
		EndTime:   now,
	}

	// 获取时间跨度
	duration := query.GetDuration()
	expectedDuration := 30 * time.Minute
	if duration != expectedDuration {
		t.Errorf("GetDuration() = %v, want %v", duration, expectedDuration)
	}
}
