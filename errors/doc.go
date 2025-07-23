/*
Package errors 提供统一的错误处理功能。

本包定义了应用程序错误类型、错误码和错误处理函数，支持错误包装、错误链追踪和结构化日志记录。
通过统一的错误处理机制，可以提高代码的可维护性和问题定位效率。

主要功能：

  - 错误码定义与管理
  - 应用错误创建与包装
  - 错误日志记录
  - HTTP状态码映射

错误码分类：

  - 通用错误 (0-19999): 如内部错误、参数错误、未找到、未授权等
  - AI相关错误 (20000-29999): 如AI服务商未找到、AI请求失败等
  - 工具相关错误 (30000-39999): 如工具未找到、工具调用失败等
  - 配置相关错误 (40000-49999): 如配置未找到、配置无效等
  - 网络相关错误 (50000-59999): 如网络超时、网络请求失败等
  - 数据库相关错误 (60000-69999): 如数据库连接失败、查询失败等

基本用法：

	// 创建新的应用错误
	err := errors.New(errors.ErrCodeInvalidParam, "参数不能为空")

	// 包装已有错误
	originalErr := someFunction()
	if originalErr != nil {
		return errors.Wrap(errors.ErrCodeNetworkFailed, "调用外部API失败", originalErr)
	}

	// 创建带详细信息的错误
	err := errors.NewWithDetail(
		errors.ErrCodeConfigNotFound,
		"配置文件未找到",
		fmt.Sprintf("路径: %s", path),
	)

	// 记录错误日志
	errors.LogError(err, "操作上下文")

	// 带字段的错误日志
	errors.LogErrorWithFields(err, "操作上下文", map[string]any{
		"user_id": 123,
		"action": "login",
	})

	// 检查错误类型
	if appErr, ok := errors.IsAppError(err); ok {
		// 处理应用错误
		code := appErr.Code
		message := appErr.Message
	}

	// 获取HTTP状态码
	statusCode := errors.GetHTTPStatusFromError(err)

预定义错误：

本包预定义了一系列常用错误实例，可以直接使用：

	errors.ErrInternalServer  // 内部服务器错误
	errors.ErrInvalidParam    // 参数错误
	errors.ErrNotFound        // 资源未找到
	errors.ErrUnauthorized    // 未授权访问
	errors.ErrForbidden       // 禁止访问
	errors.ErrConfigNotFound  // 配置文件未找到
	errors.ErrNetworkTimeout  // 网络请求超时

错误链追踪：

本包支持Go 1.13+的错误链功能，可以通过 errors.Unwrap 获取原始错误：

	if appErr, ok := err.(*errors.AppError); ok {
		originalErr := appErr.Unwrap()
		// 处理原始错误
	}
*/
package errors
