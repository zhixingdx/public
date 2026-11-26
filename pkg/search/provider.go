package search

import (
	"strings"

	"github.com/zhixingdx/oipivot-api/config"
)

// NewFromConfig 根据全局配置创建 OpenSearch 客户端；未配置任何有效地址时返回 (nil, nil)。
// 在 Wire 中与 objectstore.NewFromConfig 类似，由组装根注入 *Client。
// 服务命名空间取自 server.name（如 "base-iam"），用于避免不同服务相同表名冲突；
// 未配置时退化为仅加环境前缀的旧行为。
func NewFromConfig(c *config.Config) (*Client, error) {
	if c == nil || c.Data == nil || c.Data.GetSearch() == nil {
		return nil, nil
	}
	s := c.Data.GetSearch()
	if strings.TrimSpace(s.GetAddr()) == "" && !hasSearchAddress(s.GetAddrs()) {
		return nil, nil
	}
	serviceName := ""
	if c.Server != nil {
		serviceName = c.Server.GetName()
	}
	return NewClient(s.GetAddr(), s.GetAddrs(), s.GetUser(), s.GetPassword(), s.GetPrefix(), serviceName)
}

// hasSearchAddress 判断地址列表中是否至少包含一个有效的 OpenSearch 地址。
func hasSearchAddress(addrs []string) bool {
	for _, addr := range addrs {
		if strings.TrimSpace(addr) != "" {
			return true
		}
	}
	return false
}
