// Package models 业务模型包
// @Author lanpang
// @Date 2024/9/13 下午3:34:00
// @Desc 任务相关业务模型
package models

import (
	"time"
	"vhagar/database"
)

// DorisTask Doris 任务业务模型
// 包含 Doris 数据库连接配置和 HTTP 端口配置
type DorisTask struct {
	DatabaseConfig database.Config `json:"database_config"` // 数据库连接配置
	HTTPPort       int             `json:"http_port"`       // HTTP 端口
	TaskID         string          `json:"task_id"`         // 任务ID
	Status         string          `json:"status"`          // 任务状态
	CreatedAt      time.Time       `json:"created_at"`      // 创建时间
	UpdatedAt      time.Time       `json:"updated_at"`      // 更新时间
}

// RocketMQTask RocketMQ 任务业务模型
// 包含 RocketMQ 连接配置和认证信息
type RocketMQTask struct {
	Dashboard string    `json:"dashboard"`  // RocketMQ 控制台地址
	Username  string    `json:"username"`   // 用户名
	Password  string    `json:"password"`   // 密码
	TaskID    string    `json:"task_id"`    // 任务ID
	Status    string    `json:"status"`     // 任务状态
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// NacosTask Nacos 任务业务模型
// 包含 Nacos 服务配置和文件写入路径
type NacosTask struct {
	Server    string    `json:"server"`     // Nacos 服务器地址
	Username  string    `json:"username"`   // 用户名
	Password  string    `json:"password"`   // 密码
	Namespace string    `json:"namespace"`  // 命名空间
	WriteFile string    `json:"write_file"` // 写入文件路径
	TaskID    string    `json:"task_id"`    // 任务ID
	Status    string    `json:"status"`     // 任务状态
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// TaskStatus 任务状态常量
const (
	TaskStatusPending   = "pending"   // 等待中
	TaskStatusRunning   = "running"   // 运行中
	TaskStatusCompleted = "completed" // 已完成
	TaskStatusFailed    = "failed"    // 失败
	TaskStatusCancelled = "cancelled" // 已取消
)

// IsRunning 检查 Doris 任务是否正在运行
func (d *DorisTask) IsRunning() bool {
	return d.Status == TaskStatusRunning
}

// IsCompleted 检查 Doris 任务是否已完成
func (d *DorisTask) IsCompleted() bool {
	return d.Status == TaskStatusCompleted
}

// UpdateStatus 更新 Doris 任务状态
func (d *DorisTask) UpdateStatus(status string) {
	d.Status = status
	d.UpdatedAt = time.Now()
}

// GetDuration 获取 Doris 任务运行时长
func (d *DorisTask) GetDuration() time.Duration {
	if d.Status == TaskStatusRunning {
		return time.Since(d.CreatedAt)
	}
	return d.UpdatedAt.Sub(d.CreatedAt)
}

// IsRunning 检查 RocketMQ 任务是否正在运行
func (r *RocketMQTask) IsRunning() bool {
	return r.Status == TaskStatusRunning
}

// IsCompleted 检查 RocketMQ 任务是否已完成
func (r *RocketMQTask) IsCompleted() bool {
	return r.Status == TaskStatusCompleted
}

// UpdateStatus 更新 RocketMQ 任务状态
func (r *RocketMQTask) UpdateStatus(status string) {
	r.Status = status
	r.UpdatedAt = time.Now()
}

// GetDuration 获取 RocketMQ 任务运行时长
func (r *RocketMQTask) GetDuration() time.Duration {
	if r.Status == TaskStatusRunning {
		return time.Since(r.CreatedAt)
	}
	return r.UpdatedAt.Sub(r.CreatedAt)
}

// IsRunning 检查 Nacos 任务是否正在运行
func (n *NacosTask) IsRunning() bool {
	return n.Status == TaskStatusRunning
}

// IsCompleted 检查 Nacos 任务是否已完成
func (n *NacosTask) IsCompleted() bool {
	return n.Status == TaskStatusCompleted
}

// UpdateStatus 更新 Nacos 任务状态
func (n *NacosTask) UpdateStatus(status string) {
	n.Status = status
	n.UpdatedAt = time.Now()
}

// GetDuration 获取 Nacos 任务运行时长
func (n *NacosTask) GetDuration() time.Duration {
	if n.Status == TaskStatusRunning {
		return time.Since(n.CreatedAt)
	}
	return n.UpdatedAt.Sub(n.CreatedAt)
}

// HasWriteFile 检查 Nacos 任务是否配置了写入文件
func (n *NacosTask) HasWriteFile() bool {
	return n.WriteFile != ""
}
