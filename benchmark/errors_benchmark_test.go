package benchmark

import (
	"fmt"
	"testing"
	"vhagar/errors"
)

// 测试错误创建性能
func BenchmarkErrorCreation(b *testing.B) {
	// 重置计时器
	b.ResetTimer()

	// 执行基准测试
	for i := 0; i < b.N; i++ {
		_ = errors.New(errors.ErrCodeInvalidParam, "测试错误")
	}
}

// 测试错误包装性能
func BenchmarkErrorWrapping(b *testing.B) {
	// 创建原始错误
	originalErr := fmt.Errorf("原始错误")

	// 重置计时器
	b.ResetTimer()

	// 执行基准测试
	for i := 0; i < b.N; i++ {
		_ = errors.Wrap(errors.ErrCodeNetworkFailed, "网络错误", originalErr)
	}
}

// 测试错误检查性能
func BenchmarkErrorChecking(b *testing.B) {
	// 创建应用错误
	appErr := errors.New(errors.ErrCodeInvalidParam, "测试错误")
	err := error(appErr)

	// 重置计时器
	b.ResetTimer()

	// 执行基准测试
	for i := 0; i < b.N; i++ {
		_, _ = errors.IsAppError(err)
	}
}

// 测试错误码获取性能
func BenchmarkErrorCodeRetrieval(b *testing.B) {
	// 创建应用错误
	appErr := errors.New(errors.ErrCodeInvalidParam, "测试错误")
	err := error(appErr)

	// 重置计时器
	b.ResetTimer()

	// 执行基准测试
	for i := 0; i < b.N; i++ {
		_ = errors.GetErrorCode(err)
	}
}

// 测试HTTP状态码获取性能
func BenchmarkHTTPStatusRetrieval(b *testing.B) {
	// 创建应用错误
	appErr := errors.New(errors.ErrCodeInvalidParam, "测试错误")
	err := error(appErr)

	// 重置计时器
	b.ResetTimer()

	// 执行基准测试
	for i := 0; i < b.N; i++ {
		_ = errors.GetHTTPStatusFromError(err)
	}
}
