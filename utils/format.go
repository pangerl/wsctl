// Package utils provides common utility functions
// @Author lanpang
// @Date 2024/12/19
// @Desc String formatting and text processing utility functions
package utils

import (
	"fmt"
	"strings"
)

// CallUser 格式化用户提及字符串
// 迁移自 task/utils.go 中的 CallUser 函数
// 将用户列表格式化为提及格式，如 <@user1><@user2>
func CallUser(users []string) string {
	var result string
	if len(users) == 0 {
		return result
	}
	for _, user := range users {
		result += fmt.Sprintf("<@%s>", user)
	}
	return result
}

// CallUserWithSeparator 格式化用户提及字符串（带分隔符）
// 将用户列表格式化为提及格式，用指定分隔符分隔
func CallUserWithSeparator(users []string, separator string) string {
	if len(users) == 0 {
		return ""
	}

	mentions := make([]string, len(users))
	for i, user := range users {
		mentions[i] = fmt.Sprintf("<@%s>", user)
	}
	return strings.Join(mentions, separator)
}

// FormatUserMention 格式化单个用户提及
func FormatUserMention(user string) string {
	if user == "" {
		return ""
	}
	return fmt.Sprintf("<@%s>", user)
}

// FormatChannelMention 格式化频道提及
func FormatChannelMention(channel string) string {
	if channel == "" {
		return ""
	}
	return fmt.Sprintf("<#%s>", channel)
}

// FormatRoleMention 格式化角色提及
func FormatRoleMention(role string) string {
	if role == "" {
		return ""
	}
	return fmt.Sprintf("<@&%s>", role)
}

// TruncateString 截断字符串到指定长度
// 如果字符串长度超过 maxLen，则截断并添加省略号
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	if maxLen <= 3 {
		return s[:maxLen]
	}

	return s[:maxLen-3] + "..."
}

// PadString 填充字符串到指定长度
// 如果字符串长度小于 length，则在右侧填充空格
func PadString(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(" ", length-len(s))
}

// PadStringLeft 左填充字符串到指定长度
// 如果字符串长度小于 length，则在左侧填充空格
func PadStringLeft(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return strings.Repeat(" ", length-len(s)) + s
}

// CenterString 居中字符串到指定长度
func CenterString(s string, length int) string {
	if len(s) >= length {
		return s
	}

	totalPad := length - len(s)
	leftPad := totalPad / 2
	rightPad := totalPad - leftPad

	return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad)
}

// RemoveEmptyStrings 从字符串切片中移除空字符串
func RemoveEmptyStrings(slice []string) []string {
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if strings.TrimSpace(s) != "" {
			result = append(result, s)
		}
	}
	return result
}

// SplitAndTrim 分割字符串并去除空白
func SplitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// JoinNonEmpty 连接非空字符串
func JoinNonEmpty(slice []string, sep string) string {
	nonEmpty := RemoveEmptyStrings(slice)
	return strings.Join(nonEmpty, sep)
}

// ContainsAny 检查字符串是否包含任意一个子字符串
func ContainsAny(s string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// ContainsAll 检查字符串是否包含所有子字符串
func ContainsAll(s string, substrings []string) bool {
	for _, substr := range substrings {
		if !strings.Contains(s, substr) {
			return false
		}
	}
	return true
}

// ReverseString 反转字符串
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsEmpty 检查字符串是否为空或只包含空白字符
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsNotEmpty 检查字符串是否不为空且不只包含空白字符
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}
