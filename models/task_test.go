package models

import (
	"testing"
	"time"
)

// TestDorisTask_IsRunning 测试Doris任务是否正在运行
func TestDorisTask_IsRunning(t *testing.T) {
	// 创建测试数据
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "运行中状态",
			status:   TaskStatusRunning,
			expected: true,
		},
		{
			name:     "等待中状态",
			status:   TaskStatusPending,
			expected: false,
		},
		{
			name:     "已完成状态",
			status:   TaskStatusCompleted,
			expected: false,
		},
		{
			name:     "失败状态",
			status:   TaskStatusFailed,
			expected: false,
		},
		{
			name:     "已取消状态",
			status:   TaskStatusCancelled,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &DorisTask{
				Status: tt.status,
			}
			if got := task.IsRunning(); got != tt.expected {
				t.Errorf("DorisTask.IsRunning() = %v, 期望 %v", got, tt.expected)
			}
		})
	}
}

// TestDorisTask_IsCompleted 测试Doris任务是否已完成
func TestDorisTask_IsCompleted(t *testing.T) {
	// 创建测试数据
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "运行中状态",
			status:   TaskStatusRunning,
			expected: false,
		},
		{
			name:     "等待中状态",
			status:   TaskStatusPending,
			expected: false,
		},
		{
			name:     "已完成状态",
			status:   TaskStatusCompleted,
			expected: true,
		},
		{
			name:     "失败状态",
			status:   TaskStatusFailed,
			expected: false,
		},
		{
			name:     "已取消状态",
			status:   TaskStatusCancelled,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &DorisTask{
				Status: tt.status,
			}
			if got := task.IsCompleted(); got != tt.expected {
				t.Errorf("DorisTask.IsCompleted() = %v, 期望 %v", got, tt.expected)
			}
		})
	}
}

// TestDorisTask_UpdateStatus 测试更新Doris任务状态
func TestDorisTask_UpdateStatus(t *testing.T) {
	// 创建测试数据
	task := &DorisTask{
		Status:    TaskStatusPending,
		UpdatedAt: time.Now().Add(-1 * time.Hour), // 设置为1小时前
	}

	// 记录原始更新时间
	originalUpdatedAt := task.UpdatedAt

	// 更新状态
	newStatus := TaskStatusRunning
	task.UpdateStatus(newStatus)

	// 验证状态是否更新
	if task.Status != newStatus {
		t.Errorf("UpdateStatus() Status = %v, 期望 %v", task.Status, newStatus)
	}

	// 验证更新时间是否更新
	if !task.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("UpdateStatus() 未更新UpdatedAt字段")
	}
}

// TestDorisTask_GetDuration 测试获取Doris任务运行时长
func TestDorisTask_GetDuration(t *testing.T) {
	// 创建测试数据
	now := time.Now()
	tests := []struct {
		name      string
		status    string
		createdAt time.Time
		updatedAt time.Time
	}{
		{
			name:      "运行中任务",
			status:    TaskStatusRunning,
			createdAt: now.Add(-30 * time.Minute),
			updatedAt: now.Add(-10 * time.Minute), // 这个值在运行中状态下不会被使用
		},
		{
			name:      "已完成任务",
			status:    TaskStatusCompleted,
			createdAt: now.Add(-30 * time.Minute),
			updatedAt: now.Add(-10 * time.Minute),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &DorisTask{
				Status:    tt.status,
				CreatedAt: tt.createdAt,
				UpdatedAt: tt.updatedAt,
			}

			duration := task.GetDuration()

			if tt.status == TaskStatusRunning {
				// 对于运行中的任务，我们只能检查持续时间是否合理
				// 因为time.Since是基于当前时间的
				expectedMinDuration := time.Since(tt.createdAt) - time.Second
				if duration < expectedMinDuration {
					t.Errorf("GetDuration() = %v, 期望至少 %v", duration, expectedMinDuration)
				}
			} else {
				// 对于已完成的任务，我们可以精确检查持续时间
				expectedDuration := tt.updatedAt.Sub(tt.createdAt)
				if duration != expectedDuration {
					t.Errorf("GetDuration() = %v, 期望 %v", duration, expectedDuration)
				}
			}
		})
	}
}

// TestRocketMQTask_IsRunning 测试RocketMQ任务是否正在运行
func TestRocketMQTask_IsRunning(t *testing.T) {
	// 创建测试数据
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "运行中状态",
			status:   TaskStatusRunning,
			expected: true,
		},
		{
			name:     "等待中状态",
			status:   TaskStatusPending,
			expected: false,
		},
		{
			name:     "已完成状态",
			status:   TaskStatusCompleted,
			expected: false,
		},
		{
			name:     "失败状态",
			status:   TaskStatusFailed,
			expected: false,
		},
		{
			name:     "已取消状态",
			status:   TaskStatusCancelled,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &RocketMQTask{
				Status: tt.status,
			}
			if got := task.IsRunning(); got != tt.expected {
				t.Errorf("RocketMQTask.IsRunning() = %v, 期望 %v", got, tt.expected)
			}
		})
	}
}

// TestRocketMQTask_IsCompleted 测试RocketMQ任务是否已完成
func TestRocketMQTask_IsCompleted(t *testing.T) {
	// 创建测试数据
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "运行中状态",
			status:   TaskStatusRunning,
			expected: false,
		},
		{
			name:     "等待中状态",
			status:   TaskStatusPending,
			expected: false,
		},
		{
			name:     "已完成状态",
			status:   TaskStatusCompleted,
			expected: true,
		},
		{
			name:     "失败状态",
			status:   TaskStatusFailed,
			expected: false,
		},
		{
			name:     "已取消状态",
			status:   TaskStatusCancelled,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &RocketMQTask{
				Status: tt.status,
			}
			if got := task.IsCompleted(); got != tt.expected {
				t.Errorf("RocketMQTask.IsCompleted() = %v, 期望 %v", got, tt.expected)
			}
		})
	}
}

// TestRocketMQTask_UpdateStatus 测试更新RocketMQ任务状态
func TestRocketMQTask_UpdateStatus(t *testing.T) {
	// 创建测试数据
	task := &RocketMQTask{
		Status:    TaskStatusPending,
		UpdatedAt: time.Now().Add(-1 * time.Hour), // 设置为1小时前
	}

	// 记录原始更新时间
	originalUpdatedAt := task.UpdatedAt

	// 更新状态
	newStatus := TaskStatusRunning
	task.UpdateStatus(newStatus)

	// 验证状态是否更新
	if task.Status != newStatus {
		t.Errorf("UpdateStatus() Status = %v, 期望 %v", task.Status, newStatus)
	}

	// 验证更新时间是否更新
	if !task.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("UpdateStatus() 未更新UpdatedAt字段")
	}
}

// TestRocketMQTask_GetDuration 测试获取RocketMQ任务运行时长
func TestRocketMQTask_GetDuration(t *testing.T) {
	// 创建测试数据
	now := time.Now()
	tests := []struct {
		name      string
		status    string
		createdAt time.Time
		updatedAt time.Time
	}{
		{
			name:      "运行中任务",
			status:    TaskStatusRunning,
			createdAt: now.Add(-30 * time.Minute),
			updatedAt: now.Add(-10 * time.Minute), // 这个值在运行中状态下不会被使用
		},
		{
			name:      "已完成任务",
			status:    TaskStatusCompleted,
			createdAt: now.Add(-30 * time.Minute),
			updatedAt: now.Add(-10 * time.Minute),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &RocketMQTask{
				Status:    tt.status,
				CreatedAt: tt.createdAt,
				UpdatedAt: tt.updatedAt,
			}

			duration := task.GetDuration()

			if tt.status == TaskStatusRunning {
				// 对于运行中的任务，我们只能检查持续时间是否合理
				expectedMinDuration := time.Since(tt.createdAt) - time.Second
				if duration < expectedMinDuration {
					t.Errorf("GetDuration() = %v, 期望至少 %v", duration, expectedMinDuration)
				}
			} else {
				// 对于已完成的任务，我们可以精确检查持续时间
				expectedDuration := tt.updatedAt.Sub(tt.createdAt)
				if duration != expectedDuration {
					t.Errorf("GetDuration() = %v, 期望 %v", duration, expectedDuration)
				}
			}
		})
	}
}

// TestNacosTask_IsRunning 测试Nacos任务是否正在运行
func TestNacosTask_IsRunning(t *testing.T) {
	// 创建测试数据
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "运行中状态",
			status:   TaskStatusRunning,
			expected: true,
		},
		{
			name:     "等待中状态",
			status:   TaskStatusPending,
			expected: false,
		},
		{
			name:     "已完成状态",
			status:   TaskStatusCompleted,
			expected: false,
		},
		{
			name:     "失败状态",
			status:   TaskStatusFailed,
			expected: false,
		},
		{
			name:     "已取消状态",
			status:   TaskStatusCancelled,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &NacosTask{
				Status: tt.status,
			}
			if got := task.IsRunning(); got != tt.expected {
				t.Errorf("NacosTask.IsRunning() = %v, 期望 %v", got, tt.expected)
			}
		})
	}
}

// TestNacosTask_IsCompleted 测试Nacos任务是否已完成
func TestNacosTask_IsCompleted(t *testing.T) {
	// 创建测试数据
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{
			name:     "运行中状态",
			status:   TaskStatusRunning,
			expected: false,
		},
		{
			name:     "等待中状态",
			status:   TaskStatusPending,
			expected: false,
		},
		{
			name:     "已完成状态",
			status:   TaskStatusCompleted,
			expected: true,
		},
		{
			name:     "失败状态",
			status:   TaskStatusFailed,
			expected: false,
		},
		{
			name:     "已取消状态",
			status:   TaskStatusCancelled,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &NacosTask{
				Status: tt.status,
			}
			if got := task.IsCompleted(); got != tt.expected {
				t.Errorf("NacosTask.IsCompleted() = %v, 期望 %v", got, tt.expected)
			}
		})
	}
}

// TestNacosTask_UpdateStatus 测试更新Nacos任务状态
func TestNacosTask_UpdateStatus(t *testing.T) {
	// 创建测试数据
	task := &NacosTask{
		Status:    TaskStatusPending,
		UpdatedAt: time.Now().Add(-1 * time.Hour), // 设置为1小时前
	}

	// 记录原始更新时间
	originalUpdatedAt := task.UpdatedAt

	// 更新状态
	newStatus := TaskStatusRunning
	task.UpdateStatus(newStatus)

	// 验证状态是否更新
	if task.Status != newStatus {
		t.Errorf("UpdateStatus() Status = %v, 期望 %v", task.Status, newStatus)
	}

	// 验证更新时间是否更新
	if !task.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("UpdateStatus() 未更新UpdatedAt字段")
	}
}

// TestNacosTask_GetDuration 测试获取Nacos任务运行时长
func TestNacosTask_GetDuration(t *testing.T) {
	// 创建测试数据
	now := time.Now()
	tests := []struct {
		name      string
		status    string
		createdAt time.Time
		updatedAt time.Time
	}{
		{
			name:      "运行中任务",
			status:    TaskStatusRunning,
			createdAt: now.Add(-30 * time.Minute),
			updatedAt: now.Add(-10 * time.Minute), // 这个值在运行中状态下不会被使用
		},
		{
			name:      "已完成任务",
			status:    TaskStatusCompleted,
			createdAt: now.Add(-30 * time.Minute),
			updatedAt: now.Add(-10 * time.Minute),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &NacosTask{
				Status:    tt.status,
				CreatedAt: tt.createdAt,
				UpdatedAt: tt.updatedAt,
			}

			duration := task.GetDuration()

			if tt.status == TaskStatusRunning {
				// 对于运行中的任务，我们只能检查持续时间是否合理
				expectedMinDuration := time.Since(tt.createdAt) - time.Second
				if duration < expectedMinDuration {
					t.Errorf("GetDuration() = %v, 期望至少 %v", duration, expectedMinDuration)
				}
			} else {
				// 对于已完成的任务，我们可以精确检查持续时间
				expectedDuration := tt.updatedAt.Sub(tt.createdAt)
				if duration != expectedDuration {
					t.Errorf("GetDuration() = %v, 期望 %v", duration, expectedDuration)
				}
			}
		})
	}
}

// TestNacosTask_HasWriteFile 测试Nacos任务是否配置了写入文件
func TestNacosTask_HasWriteFile(t *testing.T) {
	// 创建测试数据
	tests := []struct {
		name      string
		writeFile string
		expected  bool
	}{
		{
			name:      "有写入文件",
			writeFile: "/tmp/nacos-config.yaml",
			expected:  true,
		},
		{
			name:      "无写入文件",
			writeFile: "",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &NacosTask{
				WriteFile: tt.writeFile,
			}
			if got := task.HasWriteFile(); got != tt.expected {
				t.Errorf("NacosTask.HasWriteFile() = %v, 期望 %v", got, tt.expected)
			}
		})
	}
}
