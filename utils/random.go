// Package utils provides common utility functions
// @Author lanpang
// @Date 2024/12/19
// @Desc Random number and UUID generation utility functions
package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mathRand "math/rand"
	"time"
)

// GetRandomDuration 生成随机时间间隔
// 迁移自 config/task.go 中的 GetRandomDuration 函数
// 生成 0-300 秒之间的随机时间间隔
func GetRandomDuration() time.Duration {
	// 创建一个新的随机数生成器，使用当前时间作为种子
	source := mathRand.NewSource(time.Now().UnixNano())
	r := mathRand.New(source)
	// 生成随机数
	randomSeconds := r.Intn(300)
	// 将随机秒数转换为时间.Duration
	duration := time.Duration(randomSeconds) * time.Second
	return duration
}

// GetRandomDurationInRange 生成指定范围内的随机时间间隔
// 参数:
//   - min: 最小时间间隔
//   - max: 最大时间间隔
//
// 返回:
//   - time.Duration: 随机时间间隔
func GetRandomDurationInRange(min, max time.Duration) time.Duration {
	if min >= max {
		return min
	}

	source := mathRand.NewSource(time.Now().UnixNano())
	r := mathRand.New(source)

	diff := max - min
	randomDiff := time.Duration(r.Int63n(int64(diff)))
	return min + randomDiff
}

// GetRandomInt 生成指定范围内的随机整数 [min, max)
func GetRandomInt(min, max int) int {
	if min >= max {
		return min
	}

	source := mathRand.NewSource(time.Now().UnixNano())
	r := mathRand.New(source)
	return r.Intn(max-min) + min
}

// GetRandomInt64 生成指定范围内的随机 int64 [min, max)
func GetRandomInt64(min, max int64) int64 {
	if min >= max {
		return min
	}

	source := mathRand.NewSource(time.Now().UnixNano())
	r := mathRand.New(source)
	return r.Int63n(max-min) + min
}

// GetRandomFloat64 生成指定范围内的随机浮点数 [min, max)
func GetRandomFloat64(min, max float64) float64 {
	if min >= max {
		return min
	}

	source := mathRand.NewSource(time.Now().UnixNano())
	r := mathRand.New(source)
	return r.Float64()*(max-min) + min
}

// GetRandomString 生成指定长度的随机字符串
// 字符集包含大小写字母和数字
func GetRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	source := mathRand.NewSource(time.Now().UnixNano())
	r := mathRand.New(source)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}

// GetRandomBytes 生成指定长度的随机字节数组
func GetRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("生成随机字节失败: %w", err)
	}
	return bytes, nil
}

// GetSecureRandomInt 使用加密安全的随机数生成器生成随机整数 [0, max)
func GetSecureRandomInt(max int64) (int64, error) {
	if max <= 0 {
		return 0, fmt.Errorf("max 必须大于 0")
	}

	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return 0, fmt.Errorf("生成安全随机数失败: %w", err)
	}
	return n.Int64(), nil
}

// GetSecureRandomString 使用加密安全的随机数生成器生成随机字符串
func GetSecureRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("生成安全随机字符串失败: %w", err)
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

// GenerateUUID 生成简单的 UUID（不依赖外部库）
func GenerateUUID() (string, error) {
	bytes, err := GetRandomBytes(16)
	if err != nil {
		return "", err
	}

	// 设置版本号 (4) 和变体位
	bytes[6] = (bytes[6] & 0x0f) | 0x40 // Version 4
	bytes[8] = (bytes[8] & 0x3f) | 0x80 // Variant bits

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		bytes[0:4],
		bytes[4:6],
		bytes[6:8],
		bytes[8:10],
		bytes[10:16]), nil
}

// ShuffleSlice 随机打乱切片顺序（Fisher-Yates 算法）
func ShuffleSlice[T any](slice []T) {
	source := mathRand.NewSource(time.Now().UnixNano())
	r := mathRand.New(source)

	for i := len(slice) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// RandomChoice 从切片中随机选择一个元素
func RandomChoice[T any](slice []T) (T, error) {
	var zero T
	if len(slice) == 0 {
		return zero, fmt.Errorf("切片为空")
	}

	source := mathRand.NewSource(time.Now().UnixNano())
	r := mathRand.New(source)

	index := r.Intn(len(slice))
	return slice[index], nil
}

// RandomChoices 从切片中随机选择多个元素（可重复）
func RandomChoices[T any](slice []T, count int) ([]T, error) {
	if len(slice) == 0 {
		return nil, fmt.Errorf("切片为空")
	}
	if count < 0 {
		return nil, fmt.Errorf("count 不能为负数")
	}

	source := mathRand.NewSource(time.Now().UnixNano())
	r := mathRand.New(source)

	result := make([]T, count)
	for i := 0; i < count; i++ {
		index := r.Intn(len(slice))
		result[i] = slice[index]
	}
	return result, nil
}
