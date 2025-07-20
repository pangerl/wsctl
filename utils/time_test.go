package utils

import (
	"testing"
	"time"
)

func TestGetZeroTime(t *testing.T) {
	// Test with a specific time
	testTime := time.Date(2024, 12, 19, 15, 30, 45, 123456789, time.UTC)
	expected := time.Date(2024, 12, 19, 0, 0, 0, 0, time.UTC)

	result := GetZeroTime(testTime)

	if !result.Equal(expected) {
		t.Errorf("GetZeroTime() = %v, want %v", result, expected)
	}
}

func TestFormatTime(t *testing.T) {
	testTime := time.Date(2024, 12, 19, 15, 30, 45, 0, time.UTC)

	tests := []struct {
		format   string
		expected string
	}{
		{DateTimeFormat, "2024-12-19 15:30:45"},
		{DateFormat, "2024-12-19"},
		{TimeFormat, "15:30:45"},
	}

	for _, test := range tests {
		result := FormatTime(testTime, test.format)
		if result != test.expected {
			t.Errorf("FormatTime(%v, %s) = %s, want %s", testTime, test.format, result, test.expected)
		}
	}
}

func TestParseTime(t *testing.T) {
	tests := []struct {
		timeStr  string
		format   string
		expected time.Time
		hasError bool
	}{
		{"2024-12-19 15:30:45", DateTimeFormat, time.Date(2024, 12, 19, 15, 30, 45, 0, time.UTC), false},
		{"2024-12-19", DateFormat, time.Date(2024, 12, 19, 0, 0, 0, 0, time.UTC), false},
		{"15:30:45", TimeFormat, time.Date(0, 1, 1, 15, 30, 45, 0, time.UTC), false},
		{"invalid", DateTimeFormat, time.Time{}, true},
	}

	for _, test := range tests {
		result, err := ParseTime(test.timeStr, test.format)

		if test.hasError {
			if err == nil {
				t.Errorf("ParseTime(%s, %s) expected error but got none", test.timeStr, test.format)
			}
		} else {
			if err != nil {
				t.Errorf("ParseTime(%s, %s) unexpected error: %v", test.timeStr, test.format, err)
			}
			if !result.Equal(test.expected) {
				t.Errorf("ParseTime(%s, %s) = %v, want %v", test.timeStr, test.format, result, test.expected)
			}
		}
	}
}

func TestGetTodayRange(t *testing.T) {
	start, end := GetTodayRange()

	// Check that start is at 00:00:00
	if start.Hour() != 0 || start.Minute() != 0 || start.Second() != 0 || start.Nanosecond() != 0 {
		t.Errorf("GetTodayRange() start time should be 00:00:00, got %v", start)
	}

	// Check that end is at 23:59:59.999999999
	if end.Hour() != 23 || end.Minute() != 59 || end.Second() != 59 {
		t.Errorf("GetTodayRange() end time should be 23:59:59.999999999, got %v", end)
	}

	// Check that they are on the same day
	startDay := GetZeroTime(start)
	endDay := GetZeroTime(end)
	if !startDay.Equal(endDay) {
		t.Errorf("GetTodayRange() start and end should be on the same day, start: %v, end: %v", startDay, endDay)
	}
}

func TestGetYesterdayRange(t *testing.T) {
	start, end := GetYesterdayRange()
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	// Check that start is yesterday at 00:00:00
	expectedStart := GetZeroTime(yesterday)
	if !start.Equal(expectedStart) {
		t.Errorf("GetYesterdayRange() start = %v, want %v", start, expectedStart)
	}

	// Check that end is yesterday at 23:59:59.999999999
	expectedEnd := expectedStart.Add(24 * time.Hour).Add(-time.Nanosecond)
	if !end.Equal(expectedEnd) {
		t.Errorf("GetYesterdayRange() end = %v, want %v", end, expectedEnd)
	}
}

func TestIsToday(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)

	if !IsToday(today) {
		t.Errorf("IsToday(%v) should be true", today)
	}

	if IsToday(yesterday) {
		t.Errorf("IsToday(%v) should be false", yesterday)
	}

	if IsToday(tomorrow) {
		t.Errorf("IsToday(%v) should be false", tomorrow)
	}
}

func TestDaysBetween(t *testing.T) {
	start := time.Date(2024, 12, 19, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		end      time.Time
		expected int
	}{
		{time.Date(2024, 12, 19, 20, 0, 0, 0, time.UTC), 0}, // Same day
		{time.Date(2024, 12, 20, 5, 0, 0, 0, time.UTC), 1},  // Next day
		{time.Date(2024, 12, 22, 15, 0, 0, 0, time.UTC), 3}, // 3 days later
		{time.Date(2024, 12, 18, 8, 0, 0, 0, time.UTC), -1}, // Previous day
	}

	for _, test := range tests {
		result := DaysBetween(start, test.end)
		if result != test.expected {
			t.Errorf("DaysBetween(%v, %v) = %d, want %d", start, test.end, result, test.expected)
		}
	}
}
