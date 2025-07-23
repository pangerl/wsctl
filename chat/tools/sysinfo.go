package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"time"
	"vhagar/libs"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

// SystemInfo 系统信息响应结构
type SystemInfo struct {
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	Data      any    `json:"data"`
}

// CPUInfo CPU信息结构
type CPUInfo struct {
	UsagePercent  float64 `json:"usage_percent"`
	LogicalCores  int     `json:"logical_cores"`
	PhysicalCores int     `json:"physical_cores,omitempty"`
}

// MemoryInfo 内存信息结构
type MemoryInfo struct {
	Total        uint64  `json:"total_bytes"`
	Available    uint64  `json:"available_bytes"`
	Used         uint64  `json:"used_bytes"`
	UsagePercent float64 `json:"usage_percent"`
	TotalGB      float64 `json:"total_gb"`
	AvailableGB  float64 `json:"available_gb"`
	UsedGB       float64 `json:"used_gb"`
}

// DiskInfo 磁盘信息结构
type DiskInfo struct {
	Path         string  `json:"path"`
	Total        uint64  `json:"total_bytes"`
	Free         uint64  `json:"free_bytes"`
	Used         uint64  `json:"used_bytes"`
	UsagePercent float64 `json:"usage_percent"`
	TotalGB      float64 `json:"total_gb"`
	FreeGB       float64 `json:"free_gb"`
	UsedGB       float64 `json:"used_gb"`
}

// AllSystemInfo 综合系统信息结构
type AllSystemInfo struct {
	CPU    CPUInfo    `json:"cpu"`
	Memory MemoryInfo `json:"memory"`
	Disks  []DiskInfo `json:"disks"`
}

// getCPUInfo 获取CPU信息
func getCPUInfo(details bool) (*CPUInfo, error) {
	// 获取CPU使用率（1秒采样）
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, libs.WrapError(libs.ErrCodeInternalErr, "获取CPU使用率失败", err)
	}

	var usagePercent float64
	if len(percentages) > 0 {
		usagePercent = percentages[0]
	}

	// 获取CPU核心数
	logicalCores := runtime.NumCPU()

	cpuInfo := &CPUInfo{
		UsagePercent: usagePercent,
		LogicalCores: logicalCores,
	}

	// 如果需要详细信息，获取物理核心数
	if details {
		physicalCores, err := cpu.Counts(false)
		if err == nil {
			cpuInfo.PhysicalCores = physicalCores
		}
	}

	return cpuInfo, nil
}

// getMemoryInfo 获取内存信息
func getMemoryInfo() (*MemoryInfo, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, libs.WrapError(libs.ErrCodeInternalErr, "获取内存信息失败", err)
	}

	return &MemoryInfo{
		Total:        vmStat.Total,
		Available:    vmStat.Available,
		Used:         vmStat.Used,
		UsagePercent: vmStat.UsedPercent,
		TotalGB:      float64(vmStat.Total) / (1024 * 1024 * 1024),
		AvailableGB:  float64(vmStat.Available) / (1024 * 1024 * 1024),
		UsedGB:       float64(vmStat.Used) / (1024 * 1024 * 1024),
	}, nil
}

// getDiskInfo 获取磁盘信息
func getDiskInfo(details bool) ([]DiskInfo, error) {
	var diskInfos []DiskInfo

	// 获取磁盘分区
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, libs.WrapError(libs.ErrCodeInternalErr, "获取磁盘分区失败", err)
	}

	for _, partition := range partitions {
		// 跳过特殊文件系统
		if partition.Fstype == "tmpfs" || partition.Fstype == "devtmpfs" {
			continue
		}

		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			// 记录错误但继续处理其他分区
			if libs.Logger != nil {
				libs.Logger.Warnw("获取磁盘使用信息失败", "mountpoint", partition.Mountpoint, "error", err)
			}
			continue
		}

		diskInfo := DiskInfo{
			Path:         partition.Mountpoint,
			Total:        usage.Total,
			Free:         usage.Free,
			Used:         usage.Used,
			UsagePercent: usage.UsedPercent,
			TotalGB:      float64(usage.Total) / (1024 * 1024 * 1024),
			FreeGB:       float64(usage.Free) / (1024 * 1024 * 1024),
			UsedGB:       float64(usage.Used) / (1024 * 1024 * 1024),
		}

		diskInfos = append(diskInfos, diskInfo)

		// 如果不需要详细信息，只返回根分区或第一个分区
		if !details && (partition.Mountpoint == "/" || len(diskInfos) == 1) {
			break
		}
	}

	return diskInfos, nil
}

// CallSystemInfoTool 系统信息工具入口
func CallSystemInfoTool(ctx context.Context, args map[string]any) (string, error) {
	// 参数验证
	infoType, ok := args["type"].(string)
	if !ok || infoType == "" {
		err := libs.NewError(libs.ErrCodeInvalidParam, "缺少或无效的 type 参数")
		libs.LogError(err, "系统信息工具调用")
		return "", err
	}

	// 获取详细信息标志，默认为false
	details := false
	if detailsVal, exists := args["details"]; exists {
		if detailsBool, ok := detailsVal.(bool); ok {
			details = detailsBool
		}
	}

	if libs.Logger != nil {
		libs.Logger.Infow("开始获取系统信息", "type", infoType, "details", details)
	}

	var result SystemInfo
	result.Type = infoType
	result.Timestamp = time.Now().Format(time.RFC3339)

	switch infoType {
	case "cpu":
		cpuInfo, err := getCPUInfo(details)
		if err != nil {
			libs.LogErrorWithFields(err, "CPU信息获取", map[string]interface{}{
				"details": details,
			})
			return "", err
		}
		result.Data = cpuInfo

	case "memory":
		memInfo, err := getMemoryInfo()
		if err != nil {
			libs.LogError(err, "内存信息获取")
			return "", err
		}
		result.Data = memInfo

	case "disk":
		diskInfos, err := getDiskInfo(details)
		if err != nil {
			libs.LogErrorWithFields(err, "磁盘信息获取", map[string]interface{}{
				"details": details,
			})
			return "", err
		}
		result.Data = diskInfos

	case "all":
		// 获取所有系统信息
		cpuInfo, err := getCPUInfo(details)
		if err != nil {
			libs.LogError(err, "CPU信息获取")
			return "", err
		}

		memInfo, err := getMemoryInfo()
		if err != nil {
			libs.LogError(err, "内存信息获取")
			return "", err
		}

		diskInfos, err := getDiskInfo(details)
		if err != nil {
			libs.LogError(err, "磁盘信息获取")
			return "", err
		}

		allInfo := AllSystemInfo{
			CPU:    *cpuInfo,
			Memory: *memInfo,
			Disks:  diskInfos,
		}
		result.Data = allInfo

	default:
		err := libs.NewErrorWithDetail(libs.ErrCodeInvalidParam, "不支持的信息类型",
			fmt.Sprintf("支持的类型: cpu, memory, disk, all，当前类型: %s", infoType))
		libs.LogError(err, "系统信息工具调用")
		return "", err
	}

	// 序列化结果
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		appErr := libs.WrapError(libs.ErrCodeInternalErr, "结果序列化失败", err)
		libs.LogError(appErr, "系统信息工具调用")
		return "", appErr
	}

	if libs.Logger != nil {
		libs.Logger.Infow("系统信息获取完成", "type", infoType, "data_size", len(jsonBytes))
	}
	return string(jsonBytes), nil
}
