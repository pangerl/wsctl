// Package database Elasticsearch 数据库连接工具
package database

import (
	"context"
	"fmt"
	"strconv"

	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

// NewElasticsearchClient 创建 Elasticsearch 数据库连接
// 参数:
//   - cfg: 数据库配置信息
//
// 返回:
//   - *elastic.Client: Elasticsearch 客户端对象
//   - error: 错误信息
func NewElasticsearchClient(cfg Config) (*elastic.Client, error) {
	// 根据 SSL 模式确定协议方案
	scheme := map[bool]string{true: "https", false: "http"}[cfg.SSLMode]

	// 构建 Elasticsearch URL
	esURL := fmt.Sprintf("%s://%s:%d", scheme, cfg.Host, cfg.Port)

	// 创建 Elasticsearch 客户端配置
	client, err := elastic.NewClient(
		elastic.SetSniff(false),                          // 禁用节点嗅探
		elastic.SetScheme(scheme),                        // 设置协议方案
		elastic.SetURL(esURL),                            // 设置服务器 URL
		elastic.SetBasicAuth(cfg.Username, cfg.Password), // 设置基本认证
		elastic.SetHealthcheck(false))                    // 禁用健康检查

	if err != nil {
		zap.S().Errorw("创建 Elasticsearch 客户端失败",
			"url", esURL,
			"username", cfg.Username,
			"err", err)
		return nil, fmt.Errorf("创建 Elasticsearch 客户端失败: %w", err)
	}

	// 执行 Ping 操作检查连接是否正常
	ctx := context.Background()
	_, _, err = client.Ping(esURL).Do(ctx)
	if err != nil {
		zap.S().Errorw("Elasticsearch 连接测试失败",
			"url", esURL,
			"username", cfg.Username,
			"err", err)
		return nil, fmt.Errorf("Elasticsearch 连接测试失败: %w", err)
	}

	zap.S().Infow("Elasticsearch 连接成功",
		"url", esURL,
		"username", cfg.Username)
	return client, nil
}

// NewElasticsearchClientWithContext 创建带上下文的 Elasticsearch 数据库连接
// 参数:
//   - ctx: 上下文对象
//   - cfg: 数据库配置信息
//
// 返回:
//   - *elastic.Client: Elasticsearch 客户端对象
//   - error: 错误信息
func NewElasticsearchClientWithContext(ctx context.Context, cfg Config) (*elastic.Client, error) {
	// 根据 SSL 模式确定协议方案
	scheme := map[bool]string{true: "https", false: "http"}[cfg.SSLMode]

	// 构建 Elasticsearch URL
	esURL := fmt.Sprintf("%s://%s:%d", scheme, cfg.Host, cfg.Port)

	// 创建 Elasticsearch 客户端配置
	client, err := elastic.NewClient(
		elastic.SetSniff(false),                          // 禁用节点嗅探
		elastic.SetScheme(scheme),                        // 设置协议方案
		elastic.SetURL(esURL),                            // 设置服务器 URL
		elastic.SetBasicAuth(cfg.Username, cfg.Password), // 设置基本认证
		elastic.SetHealthcheck(false))                    // 禁用健康检查

	if err != nil {
		zap.S().Errorw("创建 Elasticsearch 客户端失败",
			"url", esURL,
			"username", cfg.Username,
			"err", err)
		return nil, fmt.Errorf("创建 Elasticsearch 客户端失败: %w", err)
	}

	// 使用提供的上下文执行 Ping 操作检查连接
	_, _, err = client.Ping(esURL).Do(ctx)
	if err != nil {
		zap.S().Errorw("Elasticsearch 连接测试失败",
			"url", esURL,
			"username", cfg.Username,
			"err", err)
		return nil, fmt.Errorf("Elasticsearch 连接测试失败: %w", err)
	}

	zap.S().Infow("Elasticsearch 连接成功",
		"url", esURL,
		"username", cfg.Username)
	return client, nil
}

// buildElasticsearchURL 构建 Elasticsearch URL
// 这是一个辅助函数，用于构建标准的 Elasticsearch URL
// 参数:
//   - cfg: 数据库配置信息
//
// 返回:
//   - string: 格式化的 Elasticsearch URL
func buildElasticsearchURL(cfg Config) string {
	scheme := map[bool]string{true: "https", false: "http"}[cfg.SSLMode]
	return scheme + "://" + cfg.Host + ":" + strconv.Itoa(cfg.Port)
}
