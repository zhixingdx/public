package data

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/zhixingdx/oipivot-api/config"
)

// NewUniversalRedisFromConfig 根据配置创建 Redis 客户端，供 Wire 注入；与 GetRedisClient 行为一致。
func NewUniversalRedisFromConfig(c *config.Config) (redis.UniversalClient, error) {
	if c == nil || c.Data == nil || c.Data.Redis == nil {
		return nil, fmt.Errorf("data.redis 未配置")
	}
	return GetRedisClient(c.Data.Redis), nil
}
