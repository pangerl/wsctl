// Package utils provides common utility functions
// @Author lanpang
// @Date 2024/12/19
// @Desc Time-related utility functions
package utils

import (
	"time"
)

// TimeFormat 常用时间格式
const (
	DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
	TimeFormat     = "15:04:05"
)

// GetZeroTime 获取指定日期的零点时间
// 将给定时间设置为当天的 00:00:00
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// FormatTime 格式化时间
func FormatTime(t time.Time, format string) string {
	return t.Format(format)
}

// ParseTime 解析时间字符串
func ParseTime(timeStr, format string) (time.Time, error) {
	return time.Parse(format, timeStr)
}

// GetTodayRange 获取今天的时间范围 (00:00:00 到 23:59:59.999999999)
func GetTodayRange() (time.Time, time.Time) {
	now := time.Now()
	start := GetZeroTime(now)
	end := start.Add(24 * time.Hour).Add(-time.Nanosecond)
	return start, end
}

// GetYesterdayRange 获取昨天的时间范围
func GetYesterdayRange() (time.Time, time.Time) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	start := GetZeroTime(yesterday)
	end := start.Add(24 * time.Hour).Add(-time.Nanosecond)
	return start, end
}

// IsToday 判断给定时间是否为今天
func IsToday(t time.Time) bool {
	now := time.Now()
	return GetZeroTime(t).Equal(GetZeroTime(now))
}

// DaysBetween 计算两个日期之间的天数差
func DaysBetween(start, end time.Time) int {
	startDate := GetZeroTime(start)
	endDate := GetZeroTime(end)
	duration := endDate.Sub(startDate)
	return int(duration.Hours() / 24)
}
