package utils

import (
	"strings"
	"testing"
	"time"
)

func TestGetRandomDuration(t *testing.T) {
	duration := GetRandomDuration()

	// Should be between 0 and 300 seconds
	if duration < 0 || duration >= 300*time.Second {
		t.Errorf("GetRandomDuration() = %v, want between 0 and 300 seconds", duration)
	}
}

func TestGetRandomDurationInRange(t *testing.T) {
	min := 5 * time.Second
	max := 10 * time.Second

	duration := GetRandomDurationInRange(min, max)

	if duration < min || duration >= max {
		t.Errorf("GetRandomDurationInRange(%v, %v) = %v, want between %v and %v", min, max, duration, min, max)
	}

	// Test edge case: min >= max
	duration = GetRandomDurationInRange(max, min)
	if duration != max {
		t.Errorf("GetRandomDurationInRange(%v, %v) = %v, want %v", max, min, duration, max)
	}
}

func TestGetRandomInt(t *testing.T) {
	min := 10
	max := 20

	for i := 0; i < 100; i++ {
		result := GetRandomInt(min, max)
		if result < min || result >= max {
			t.Errorf("GetRandomInt(%d, %d) = %d, want between %d and %d", min, max, result, min, max-1)
		}
	}

	// Test edge case: min >= max
	result := GetRandomInt(max, min)
	if result != max {
		t.Errorf("GetRandomInt(%d, %d) = %d, want %d", max, min, result, max)
	}
}

func TestGetRandomInt64(t *testing.T) {
	min := int64(100)
	max := int64(200)

	for i := 0; i < 100; i++ {
		result := GetRandomInt64(min, max)
		if result < min || result >= max {
			t.Errorf("GetRandomInt64(%d, %d) = %d, want between %d and %d", min, max, result, min, max-1)
		}
	}
}

func TestGetRandomFloat64(t *testing.T) {
	min := 1.5
	max := 2.5

	for i := 0; i < 100; i++ {
		result := GetRandomFloat64(min, max)
		if result < min || result >= max {
			t.Errorf("GetRandomFloat64(%f, %f) = %f, want between %f and %f", min, max, result, min, max)
		}
	}
}

func TestGetRandomString(t *testing.T) {
	length := 10
	result := GetRandomString(length)

	if len(result) != length {
		t.Errorf("GetRandomString(%d) length = %d, want %d", length, len(result), length)
	}

	// Check that it only contains valid characters
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for _, char := range result {
		if !strings.ContainsRune(charset, char) {
			t.Errorf("GetRandomString() contains invalid character: %c", char)
		}
	}

	// Test that multiple calls produce different results (with high probability)
	result2 := GetRandomString(length)
	if result == result2 {
		t.Logf("GetRandomString() produced same result twice: %s (this is unlikely but possible)", result)
	}
}

func TestGetRandomBytes(t *testing.T) {
	length := 16
	bytes, err := GetRandomBytes(length)
	if err != nil {
		t.Errorf("GetRandomBytes(%d) error = %v", length, err)
		return
	}

	if len(bytes) != length {
		t.Errorf("GetRandomBytes(%d) length = %d, want %d", length, len(bytes), length)
	}
}

func TestGetSecureRandomInt(t *testing.T) {
	max := int64(100)

	for i := 0; i < 10; i++ {
		result, err := GetSecureRandomInt(max)
		if err != nil {
			t.Errorf("GetSecureRandomInt(%d) error = %v", max, err)
			continue
		}

		if result < 0 || result >= max {
			t.Errorf("GetSecureRandomInt(%d) = %d, want between 0 and %d", max, result, max-1)
		}
	}

	// Test edge case: max <= 0
	_, err := GetSecureRandomInt(0)
	if err == nil {
		t.Error("GetSecureRandomInt(0) expected error, got nil")
	}
}

func TestGetSecureRandomString(t *testing.T) {
	length := 12
	result, err := GetSecureRandomString(length)
	if err != nil {
		t.Errorf("GetSecureRandomString(%d) error = %v", length, err)
		return
	}

	if len(result) != length {
		t.Errorf("GetSecureRandomString(%d) length = %d, want %d", length, len(result), length)
	}

	// Check that it only contains valid characters
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for _, char := range result {
		if !strings.ContainsRune(charset, char) {
			t.Errorf("GetSecureRandomString() contains invalid character: %c", char)
		}
	}
}

func TestGenerateUUID(t *testing.T) {
	uuid, err := GenerateUUID()
	if err != nil {
		t.Errorf("GenerateUUID() error = %v", err)
		return
	}

	// Check UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	if len(uuid) != 36 {
		t.Errorf("GenerateUUID() length = %d, want 36", len(uuid))
	}

	parts := strings.Split(uuid, "-")
	if len(parts) != 5 {
		t.Errorf("GenerateUUID() parts = %d, want 5", len(parts))
	}

	expectedLengths := []int{8, 4, 4, 4, 12}
	for i, part := range parts {
		if len(part) != expectedLengths[i] {
			t.Errorf("GenerateUUID() part %d length = %d, want %d", i, len(part), expectedLengths[i])
		}
	}

	// Test that multiple calls produce different UUIDs
	uuid2, err := GenerateUUID()
	if err != nil {
		t.Errorf("GenerateUUID() second call error = %v", err)
		return
	}

	if uuid == uuid2 {
		t.Errorf("GenerateUUID() produced same UUID twice: %s", uuid)
	}
}

func TestShuffleSlice(t *testing.T) {
	original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	slice := make([]int, len(original))
	copy(slice, original)

	ShuffleSlice(slice)

	// Check that all elements are still present
	if len(slice) != len(original) {
		t.Errorf("ShuffleSlice() changed slice length from %d to %d", len(original), len(slice))
	}

	// Count elements to ensure none are lost
	counts := make(map[int]int)
	for _, v := range slice {
		counts[v]++
	}

	for _, v := range original {
		if counts[v] != 1 {
			t.Errorf("ShuffleSlice() element %d appears %d times, want 1", v, counts[v])
		}
	}

	// Test empty slice
	var empty []int
	ShuffleSlice(empty) // Should not panic
}

func TestRandomChoice(t *testing.T) {
	slice := []string{"a", "b", "c", "d", "e"}

	for i := 0; i < 10; i++ {
		choice, err := RandomChoice(slice)
		if err != nil {
			t.Errorf("RandomChoice() error = %v", err)
			continue
		}

		found := false
		for _, v := range slice {
			if v == choice {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("RandomChoice() = %s, not found in slice", choice)
		}
	}

	// Test empty slice
	var empty []string
	_, err := RandomChoice(empty)
	if err == nil {
		t.Error("RandomChoice() with empty slice expected error, got nil")
	}
}

func TestRandomChoices(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	count := 3

	choices, err := RandomChoices(slice, count)
	if err != nil {
		t.Errorf("RandomChoices() error = %v", err)
		return
	}

	if len(choices) != count {
		t.Errorf("RandomChoices() length = %d, want %d", len(choices), count)
	}

	// Check that all choices are from the original slice
	for _, choice := range choices {
		found := false
		for _, v := range slice {
			if v == choice {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("RandomChoices() contains %d, not found in original slice", choice)
		}
	}

	// Test empty slice
	var empty []int
	_, err = RandomChoices(empty, 1)
	if err == nil {
		t.Error("RandomChoices() with empty slice expected error, got nil")
	}

	// Test negative count
	_, err = RandomChoices(slice, -1)
	if err == nil {
		t.Error("RandomChoices() with negative count expected error, got nil")
	}
}
