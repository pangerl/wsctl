// Package models 业务模型包
// @Author lanpang
// @Date 2024/9/6 上午11:37:00
// @Desc 指标相关业务模型
package models

import (
	"time"
)

// MetricData 指标数据模型
// 用于存储时序指标数据
type MetricData struct {
	Timestamp time.Time              `json:"timestamp"` // 时间戳
	Metrics   map[string]interface{} `json:"metrics"`   // 指标数据
	Tags      map[string]string      `json:"tags"`      // 标签信息
	Source    string                 `json:"source"`    // 数据源
}

// HostMetric 主机指标模型
// 用于存储主机性能指标
type HostMetric struct {
	HostID              string    `json:"host_id"`                // 主机ID
	CPUCores            float64   `json:"cpu_cores"`              // CPU核心数
	CPUUsageActive      float64   `json:"cpu_usage_active"`       // CPU使用率
	MemUsedPercent      float64   `json:"mem_used_percent"`       // 内存使用百分比
	MemTotal            float64   `json:"mem_total"`              // 总内存
	NetBytesRecv        float64   `json:"net_bytes_recv"`         // 网络接收字节数
	NetBytesSent        float64   `json:"net_bytes_sent"`         // 网络发送字节数
	RootDiskUsedPercent float64   `json:"root_disk_used_percent"` // 根磁盘使用百分比
	DataDiskUsedPercent float64   `json:"data_disk_used_percent"` // 数据磁盘使用百分比
	NTPOffsetMs         float64   `json:"ntp_offset_ms"`          // NTP时间偏移（毫秒）
	Timestamp           time.Time `json:"timestamp"`              // 采集时间
}

// MetricCollection 指标集合模型
// 用于批量处理指标数据
type MetricCollection struct {
	CollectionID string            `json:"collection_id"` // 集合ID
	Metrics      []*MetricData     `json:"metrics"`       // 指标列表
	StartTime    time.Time         `json:"start_time"`    // 开始时间
	EndTime      time.Time         `json:"end_time"`      // 结束时间
	Source       string            `json:"source"`        // 数据源
	Tags         map[string]string `json:"tags"`          // 公共标签
}

// AlertRule 告警规则模型
// 用于定义指标告警规则
type AlertRule struct {
	RuleID     string            `json:"rule_id"`     // 规则ID
	Name       string            `json:"name"`        // 规则名称
	MetricName string            `json:"metric_name"` // 指标名称
	Condition  string            `json:"condition"`   // 告警条件
	Threshold  float64           `json:"threshold"`   // 阈值
	Duration   time.Duration     `json:"duration"`    // 持续时间
	Labels     map[string]string `json:"labels"`      // 标签
	Enabled    bool              `json:"enabled"`     // 是否启用
	CreatedAt  time.Time         `json:"created_at"`  // 创建时间
	UpdatedAt  time.Time         `json:"updated_at"`  // 更新时间
}

// MetricQuery 指标查询模型
// 用于查询指标数据
type MetricQuery struct {
	QueryID   string            `json:"query_id"`   // 查询ID
	Query     string            `json:"query"`      // 查询语句
	StartTime time.Time         `json:"start_time"` // 开始时间
	EndTime   time.Time         `json:"end_time"`   // 结束时间
	Step      time.Duration     `json:"step"`       // 步长
	Tags      map[string]string `json:"tags"`       // 查询标签
}

// AddMetric 向指标数据中添加单个指标
func (md *MetricData) AddMetric(name string, value interface{}) {
	if md.Metrics == nil {
		md.Metrics = make(map[string]interface{})
	}
	md.Metrics[name] = value
}

// AddTag 向指标数据中添加标签
func (md *MetricData) AddTag(key, value string) {
	if md.Tags == nil {
		md.Tags = make(map[string]string)
	}
	md.Tags[key] = value
}

// GetMetric 获取指定名称的指标值
func (md *MetricData) GetMetric(name string) (interface{}, bool) {
	if md.Metrics == nil {
		return nil, false
	}
	value, exists := md.Metrics[name]
	return value, exists
}

// GetTag 获取指定标签值
func (md *MetricData) GetTag(key string) (string, bool) {
	if md.Tags == nil {
		return "", false
	}
	value, exists := md.Tags[key]
	return value, exists
}

// IsHealthy 检查主机指标是否健康
func (hm *HostMetric) IsHealthy() bool {
	// CPU使用率小于80%，内存使用率小于90%，磁盘使用率小于85%
	return hm.CPUUsageActive < 80.0 &&
		hm.MemUsedPercent < 90.0 &&
		hm.RootDiskUsedPercent < 85.0 &&
		hm.DataDiskUsedPercent < 85.0
}

// GetCPUUsageLevel 获取CPU使用率等级
func (hm *HostMetric) GetCPUUsageLevel() string {
	if hm.CPUUsageActive < 50.0 {
		return "低"
	} else if hm.CPUUsageActive < 80.0 {
		return "中"
	}
	return "高"
}

// GetMemoryUsageLevel 获取内存使用率等级
func (hm *HostMetric) GetMemoryUsageLevel() string {
	if hm.MemUsedPercent < 60.0 {
		return "低"
	} else if hm.MemUsedPercent < 85.0 {
		return "中"
	}
	return "高"
}

// AddMetric 向集合中添加指标
func (mc *MetricCollection) AddMetric(metric *MetricData) {
	mc.Metrics = append(mc.Metrics, metric)

	// 更新时间范围
	if mc.StartTime.IsZero() || metric.Timestamp.Before(mc.StartTime) {
		mc.StartTime = metric.Timestamp
	}
	if mc.EndTime.IsZero() || metric.Timestamp.After(mc.EndTime) {
		mc.EndTime = metric.Timestamp
	}
}

// GetMetricCount 获取指标数量
func (mc *MetricCollection) GetMetricCount() int {
	return len(mc.Metrics)
}

// GetDuration 获取采集时间跨度
func (mc *MetricCollection) GetDuration() time.Duration {
	if mc.StartTime.IsZero() || mc.EndTime.IsZero() {
		return 0
	}
	return mc.EndTime.Sub(mc.StartTime)
}

// IsEnabled 检查告警规则是否启用
func (ar *AlertRule) IsEnabled() bool {
	return ar.Enabled
}

// UpdateThreshold 更新告警阈值
func (ar *AlertRule) UpdateThreshold(threshold float64) {
	ar.Threshold = threshold
	ar.UpdatedAt = time.Now()
}

// Enable 启用告警规则
func (ar *AlertRule) Enable() {
	ar.Enabled = true
	ar.UpdatedAt = time.Now()
}

// Disable 禁用告警规则
func (ar *AlertRule) Disable() {
	ar.Enabled = false
	ar.UpdatedAt = time.Now()
}

// SetTimeRange 设置查询时间范围
func (mq *MetricQuery) SetTimeRange(start, end time.Time) {
	mq.StartTime = start
	mq.EndTime = end
}

// AddTag 向查询中添加标签过滤条件
func (mq *MetricQuery) AddTag(key, value string) {
	if mq.Tags == nil {
		mq.Tags = make(map[string]string)
	}
	mq.Tags[key] = value
}

// GetDuration 获取查询时间跨度
func (mq *MetricQuery) GetDuration() time.Duration {
	return mq.EndTime.Sub(mq.StartTime)
}
