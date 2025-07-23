# Implementation Plan

- [x] 1. 创建新的包结构和基础框架

  - 创建功能性目录结构（database/, logger/, errors/, utils/, models/, config/）
  - 设置包级别的文档注释和基础文件
  - _Requirements: 2.1, 2.2_

- [x] 1.1 创建 database 包结构

  - 创建 database/ 目录和基础文件
  - 定义数据库配置结构体和接口
  - 添加包级别的文档注释
  - _Requirements: 2.1, 2.2_

- [x] 1.2 创建 logger 包结构

  - 创建 logger/ 目录和基础文件
  - 定义日志配置结构体和接口
  - 添加包级别的文档注释
  - _Requirements: 2.1, 2.2_

- [x] 1.3 创建 errors 包结构

  - 创建 errors/ 目录和基础文件
  - 定义错误处理结构体和常量
  - 添加包级别的文档注释
  - _Requirements: 2.1, 2.2_

- [x] 1.4 创建 utils 包结构

  - 创建 utils/ 目录和基础文件
  - 定义工具函数的接口和结构
  - 添加包级别的文档注释
  - _Requirements: 2.1, 2.2_

- [x] 1.5 创建 models 和 config 包结构

  - 创建 models/ 和 config/ 目录
  - 设置业务模型和配置管理的基础文件结构
  - 添加包级别的文档注释
  - _Requirements: 3.1, 3.2_

- [x] 2. 迁移数据库相关代码

  - 将 libs/ 中的数据库连接代码迁移到 database/
  - 重构数据库配置结构体
  - 更新数据库连接函数和接口
  - _Requirements: 2.1, 2.2, 4.1, 4.2_

- [x] 2.1 迁移数据库配置结构体

  - 将 libs/config.go 中的 DB 和 RedisConfig 迁移到 database/config.go
  - 重命名结构体字段以符合 Go 命名规范
  - 保持 HasValue() 方法的功能
  - _Requirements: 2.1, 4.2_

- [x] 2.2 迁移 MySQL 连接代码

  - 将 libs/mysql.go 迁移到 database/mysql.go
  - 重构 NewMysqlClient 函数，使用新的配置结构体
  - 添加错误处理和日志记录
  - _Requirements: 2.1, 2.2, 4.2_

- [x] 2.3 迁移 PostgreSQL 连接代码

  - 将 libs/pg.go 迁移到 database/postgres.go
  - 重构 PGClienter 结构体和相关函数
  - 使用新的配置结构体和错误处理
  - _Requirements: 2.1, 2.2, 4.2_

- [x] 2.4 迁移 Redis 连接代码

  - 将 libs/redis.go 迁移到 database/redis.go
  - 重构 NewRedisClient 函数
  - 使用新的 RedisConfig 结构体
  - _Requirements: 2.1, 2.2, 4.2_

- [x] 2.5 迁移 Elasticsearch 连接代码

  - 将 libs/es.go 迁移到 database/elasticsearch.go
  - 重构 NewESClient 函数
  - 使用新的配置结构体和错误处理
  - _Requirements: 2.1, 2.2, 4.2_

- [x] 3. 迁移日志和错误处理代码

  - 将日志初始化代码迁移到 logger/
  - 将错误处理代码迁移到 errors/
  - 确保新的错误处理机制正常工作
  - _Requirements: 2.1, 2.2, 4.1, 4.2_

- [x] 3.1 迁移日志处理代码

  - 将 libs/logger.go 迁移到 logger/logger.go
  - 重构 InitLoggerWithConfig 函数
  - 定义日志配置结构体
  - _Requirements: 2.1, 2.2, 4.2_

- [x] 3.2 迁 移错误处理代码

  - 将 libs/errors.go 迁移到 errors/errors.go
  - 保持所有错误码和错误类型定义
  - 确保错误处理函数正常工作
  - _Requirements: 2.1, 2.2, 4.2_

- [x] 4. 迁移工具函数

  - 将通用工具函数迁移到 utils/
  - 按功能分类组织工具函数
  - 添加单元测试覆盖
  - _Requirements: 2.1, 2.3, 6.3_

- [x] 4.1 迁移时间相关工具函数

  - 将 task/utils.go 中的 GetZeroTime 迁移到 utils/time.go
  - 添加其他时间处理工具函数
  - 编写单元测试
  - _Requirements: 2.3, 6.3_

- [x] 4.2 迁移 HTTP 请求工具函数

  - 将 task/utils.go 中的 DoRequest 迁移到 utils/http.go
  - 重构函数以使用新的错误处理机制
  - 编写单元测试
  - _Requirements: 2.3, 6.3_

- [x] 4.3 迁移随机数工具函数

  - 将 config/task.go 中的 GetRandomDuration 迁移到 utils/random.go
  - 添加其他随机数生成工具函数
  - 编写单元测试
  - _Requirements: 2.3, 6.3_

- [x] 4.4 迁移其他工具函数

  - 将 task/utils.go 中的 CallUser 等函数迁移到合适的包中
  - 按功能分类组织工具函数
  - 编写单元测试
  - _Requirements: 2.3, 6.3_

- [x] 5. 重构配置管理

  - 创建新的配置结构体和加载逻辑
  - 分离配置定义和业务模型
  - 实现配置文件加载和验证
  - _Requirements: 1.1, 1.2, 1.3, 4.1, 4.2_

- [x] 5.1 创建主配置结构体

  - 在 config/config.go 中定义新的 Config 结构体
  - 整合所有配置项到统一的结构中
  - 使用新的数据库和服务配置结构体
  - _Requirements: 1.1, 1.2, 4.2_

- [x] 5.2 创建服务配置结构体

  - 在 config/service.go 中定义服务相关配置
  - 包括 AI、Weather、RocketMQ、Nacos 等服务配置
  - 保持与原有配置的兼容性
  - _Requirements: 1.1, 1.2, 4.2_

- [x] 5.3 实现配置加载器

  - 在 config/loader.go 中实现配置文件加载逻辑
  - 支持 TOML 格式配置文件
  - 添加配置验证和错误处理
  - _Requirements: 1.1, 1.4, 4.1, 4.2_

- [x] 5.4 创建配置验证逻辑

  - 实现配置项的完整性检查
  - 添加必填项验证
  - 提供清晰的配置错误信息
  - _Requirements: 1.4, 4.1, 4.2_

- [x] 6. 迁移业务模型

  - 将业务相关的结构体迁移到 models/
  - 按业务域分组组织模型
  - 确保模型的完整性和一致性
  - _Requirements: 3.1, 3.2, 3.3, 4.1, 4.2_

- [x] 6.1 迁移租户相关模型

  - 将 config/tenant.go 中的 Tenant 和 Corp 迁移到 models/tenant.go
  - 重构结构体字段命名以符合 Go 规范
  - 添加业务逻辑方法
  - _Requirements: 3.1, 3.2, 3.3_

- [x] 6.2 迁移任务相关模型

  - 将 config/task.go 中的业务模型迁移到 models/
  - 包括 DorisCfg、RocketMQCfg、NacosCfg 等
  - 区分配置结构体和业务模型
  - _Requirements: 3.1, 3.2, 3.3_

- [x] 6.3 创建指标相关模型

  - 将指标相关的结构体整理到 models/metric.go
  - 包括 MetricCfg 和其他监控相关结构体
  - 添加业务逻辑方法
  - _Requirements: 3.1, 3.2, 3.3_

- [x] 7. 更新所有导入路径

  - 批量更新项目中所有的 vhagar/libs 导入
  - 更新配置相关的导入路径
  - 确保所有文件都能正常编译
  - _Requirements: 4.1, 4.2, 4.3_

- [x] 7.1 更新 task 包中的导入路径

  - 更新 task/ 目录下所有文件的导入路径
  - 将 vhagar/libs 替换为新的包路径
  - 确保功能保持不变
  - _Requirements: 4.1, 4.2_

- [x] 7.2 更新 chat 包中的导入路径

  - 更新 chat/ 目录下所有文件的导入路径
  - 将 vhagar/libs 替换为新的包路径
  - 确保 AI 功能正常工作
  - _Requirements: 4.1, 4.2_

- [x] 7.3 更新 notify 包中的导入路径

  - 更新 notify/ 目录下所有文件的导入路径
  - 将 vhagar/libs 替换为新的包路径
  - 确保通知功能正常工作
  - _Requirements: 4.1, 4.2_

- [x] 7.4 更新 metric 包中的导入路径

  - 更新 metric/ 目录下所有文件的导入路径
  - 将 vhagar/libs 替换为新的包路径
  - 确保监控功能正常工作
  - _Requirements: 4.1, 4.2_

- [x] 7.5 更新 cmd 包中的导入路径

  - 更新 cmd/ 目录下所有文件的导入路径
  - 将 vhagar/libs 替换为新的包路径
  - 确保命令行功能正常工作
  - _Requirements: 4.1, 4.2_

- [x] 8. 创建向后兼容层

  - 在原有的 config 包中创建兼容性入口
  - 重新导出主要的类型和函数
  - 确保现有代码无需修改即可工作
  - _Requirements: 4.1, 4.2, 4.3_

- [x] 8.1 创建配置包兼容入口

  - 在 config/config.go 中重新导出新的配置类型
  - 提供类型别名以保持向后兼容
  - 重新导出主要的配置加载函数
  - _Requirements: 4.1, 4.2, 4.3_

- [x] 8.2 创建 libs 包兼容入口

  - 创建临时的 libs 包入口文件
  - 重新导出所有迁移的类型和函数
  - 添加废弃警告注释
  - _Requirements: 4.1, 4.2, 4.3_

- [x] 9. 编写单元测试

  - 为所有新创建的包编写单元测试
  - 确保测试覆盖率达到要求
  - 验证重构后的功能正确性
  - _Requirements: 6.1, 6.2, 6.3, 6.4_

- [x] 9.1 编写数据库包测试

  - 为 database/ 下的所有函数编写单元测试
  - 测试数据库连接和配置验证功能
  - 使用 mock 对象测试数据库交互
  - _Requirements: 6.1, 6.2, 6.3_

- [x] 9.2 编写日志包测试

  - 为 logger/ 下的函数编写单元测试
  - 测试不同日志级别和输出方式
  - 验证日志格式和文件输出
  - _Requirements: 6.1, 6.2, 6.3_

- [x] 9.3 编写错误处理包测试

  - 为 errors/ 下的函数编写单元测试
  - 测试错误码和错误信息的正确性
  - 验证错误包装和解包功能
  - _Requirements: 6.1, 6.2, 6.3_

- [x] 9.4 编写工具包测试

  - 为 utils/ 下的所有工具函数编写单元测试
  - 测试边界条件和异常情况
  - 确保工具函数的健壮性
  - _Requirements: 6.1, 6.2, 6.3_

- [x] 9.5 编写配置包测试

  - 为 config/ 下的函数编写单元测试
  - 测试配置文件加载和验证功能
  - 使用不同的配置文件格式进行测试
  - _Requirements: 6.1, 6.2, 6.4_

- [x] 9.6 编写模型包测试

  - 为 models/ 下的结构体编写单元测试
  - 测试业务逻辑方法的正确性
  - 验证数据模型的完整性
  - _Requirements: 6.1, 6.2, 6.3_

- [x] 10. 集成测试和验证

  - 运行完整的应用程序测试
  - 验证所有功能模块正常工作
  - 确保性能没有显著下降
  - _Requirements: 4.3, 4.4_

- [x] 10.1 运行现有测试套件

  - 执行项目中所有现有的测试用例
  - 确保所有测试都能通过
  - 修复因重构导致的测试失败
  - _Requirements: 4.3, 4.4_

- [x] 10.2 执行应用程序功能测试

  - 启动应用程序并验证基本功能
  - 测试配置文件加载和解析
  - 验证数据库连接和查询功能
  - _Requirements: 4.3, 4.4_

- [x] 10.3 执行性能基准测试

  - 对比重构前后的性能指标
  - 确保没有引入性能回归
  - 优化发现的性能瓶颈
  - _Requirements: 4.3, 4.4_

- [x] 11. 清理和文档更新

  - 删除旧的 libs 目录
  - 更新项目文档和 README
  - 添加迁移指南和最佳实践
  - _Requirements: 5.1, 5.2, 5.3, 5.4_

- [x] 11.1 清理旧代码

  - 删除原有的 libs/ 目录
  - 移除临时的兼容性代码
  - 清理不再使用的导入和变量
  - _Requirements: 5.1, 5.2_

- [x] 11.2 更新项目文档

  - 更新 README.md 中的项目结构说明
  - 添加新的包结构和使用说明
  - 更新开发者指南和贡献指南
  - _Requirements: 5.1, 5.2, 5.3_

- [x] 11.3 创建迁移指南

  - 编写详细的代码迁移指南
  - 提供新旧代码对比示例
  - 说明最佳实践和注意事项
  - _Requirements: 5.3, 5.4_

- [x] 11.4 添加包级别文档
  - 为所有新创建的包添加详细的文档注释
  - 包括使用示例和 API 说明
  - 确保 godoc 生成的文档清晰易懂
  - _Requirements: 5.2, 5.3, 5.4_
