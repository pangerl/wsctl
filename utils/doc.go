/*
Package utils 提供通用工具函数。

本包包含各种通用的工具函数，如字符串处理、HTTP请求、时间处理和随机数生成等，
这些函数可以在应用程序的各个部分重复使用，提高代码复用性和一致性。

主要功能：

  - 字符串处理: 截断、填充、分割等
  - HTTP请求: GET、POST、PUT、DELETE等
  - 时间处理: 格式化、解析、计算等
  - 随机数生成: 随机整数、字符串、持续时间等
  - 用户提及: 格式化用户、频道、角色提及

字符串处理：

	// 截断字符串
	truncated := utils.TruncateString("这是一个很长的字符串", 10)

	// 填充字符串
	padded := utils.PadString("短字符串", 10, ' ')

	// 居中字符串
	centered := utils.CenterString("居中", 10, '-')

	// 移除空字符串
	filtered := utils.RemoveEmptyStrings([]string{"a", "", "b", ""})

	// 分割并修剪
	parts := utils.SplitAndTrim("a, b, c", ",")

HTTP请求：

	// 创建HTTP客户端
	client := utils.NewHTTPClient(5*time.Second, 3)

	// 发送GET请求
	resp, err := client.Get("https://api.example.com/users")

	// 发送POST请求
	resp, err := client.Post("https://api.example.com/users", "application/json", body)

	// 简化的请求函数
	result, err := utils.DoRequest("GET", "https://api.example.com/users", nil)

时间处理：

	// 获取零点时间
	zeroTime := utils.GetZeroTime(time.Now())

	// 格式化时间
	formatted := utils.FormatTime(time.Now(), "yyyy-MM-dd HH:mm:ss")

	// 解析时间
	parsed, err := utils.ParseTime("2025-07-21 12:34:56", "yyyy-MM-dd HH:mm:ss")

	// 获取今天的时间范围
	start, end := utils.GetTodayRange()

	// 检查是否为今天
	if utils.IsToday(someTime) {
		// 处理今天的时间
	}

随机数生成：

	// 生成随机持续时间
	duration := utils.GetRandomDuration(1*time.Minute, 5*time.Minute)

	// 生成随机整数
	num := utils.GetRandomInt(1, 100)

	// 生成随机字符串
	str := utils.GetRandomString(10)

	// 生成UUID
	uuid := utils.GenerateUUID()

用户提及：

	// 格式化用户提及
	mention := utils.FormatUserMention("user123", "张三")

	// 调用用户
	message := utils.CallUser("user123", "张三", "你好")

注意事项：

  - HTTP请求函数会自动处理重试和错误记录
  - 随机数生成函数使用安全的随机源
  - 时间处理函数考虑了时区问题
*/
package utils
