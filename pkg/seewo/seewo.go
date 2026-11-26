package seewo

import (
	"context"
	"fmt"
	"strings"

	"github.com/zhixingdx/public/pkg/seewo/api"
)

// Seewo 希沃开放平台 SDK 主入口
// 参考 Java SDK: com.seewo.open.sdk.DefaultSeewoClient
type Seewo struct {

	// API 模块
	Org  *api.OrgAPI  // 组织管理
	User *api.UserAPI // 用户管理
}

// NewWithConfig 使用配置对象创建希沃 SDK 实例
func NewWithConfig(config *Config) (*Seewo, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	client := NewClient(config)

	return &Seewo{
		Org:  api.NewOrgAPI(client),
		User: api.NewUserAPI(client),
	}, nil
}

// TestConnection 测试连接
// 使用最简单的查询接口验证配置和 API 连通性
// 使用 GetOrgByCode 接口，即使查询不到数据，只要 API 能正常响应（签名正确、连接正常）即可
func (s *Seewo) TestConnection(ctx context.Context) error {
	// 使用一个测试编码进行查询
	// 如果 API 能正常响应（即使返回数据不存在），说明连接和签名都是正常的
	_, err := s.Org.GetOrgByCode(ctx, &api.GetOrgByCodeRequest{
		ReqDto: api.GetOrgByCodeReqDto{
			Code: "TEST",
		},
	})

	// 如果返回的是连接错误或签名错误，则测试失败
	// 如果返回的是业务错误（如数据不存在），说明连接和签名正常，测试成功
	if err != nil {
		errStr := err.Error()
		errStrLower := strings.ToLower(errStr)

		// 检查是否是连接相关的错误（请求失败、超时、网络错误等）
		if strings.Contains(errStr, "请求失败") || strings.Contains(errStrLower, "connection") ||
			strings.Contains(errStrLower, "timeout") || strings.Contains(errStrLower, "network") ||
			strings.Contains(errStrLower, "dial") {
			return fmt.Errorf("连接测试失败: %w", err)
		}

		// 检查是否是认证/授权错误（401、403、签名错误等）
		if strings.Contains(errStr, "HTTP状态码=401") || strings.Contains(errStr, "HTTP状态码=403") ||
			strings.Contains(errStrLower, "401") || strings.Contains(errStrLower, "403") ||
			strings.Contains(errStrLower, "unauthorized") || strings.Contains(errStrLower, "forbidden") ||
			strings.Contains(errStr, "签名") || strings.Contains(errStr, "认证") {
			return fmt.Errorf("连接测试失败: %w", err)
		}

		// 其他错误（如业务错误：数据不存在、参数错误等）视为连接测试成功
		// 因为这说明 API 能正常响应，只是业务逻辑层面的错误
	}
	return nil
}
