package cache

import (
	"github.com/redis/go-redis/v9"
	"github.com/zhixingdx/public/api/config"
	"github.com/zhixingdx/public/pkg/lang"
	"github.com/zhixingdx/public/pkg/logger/kratoslog"
	"golang.org/x/sync/singleflight"
)

// dataCachePrefixAndExpire 解析 data.cache 的前缀与表缓存过期秒数；未配置 cache 时 prefix 为空、过期默认 86400（1 天）。
func dataCachePrefixAndExpire(d *config.Data) (prefix string, expireSec int32) {
	expireSec = 86400
	if d != nil && d.Cache != nil {
		prefix = d.Cache.Prefix
		expireSec = d.Cache.Expiration
	}
	return prefix, expireSec
}

// NewDataRedisCache 组装 Data 层使用的 Redis Cache（前缀、过期语义与原先 NewData 内联逻辑一致）。
func NewDataRedisCache(c *config.Config, client redis.UniversalClient, common *lang.CommonLang, logger log.Logger) Cache {
	prefix := ""
	if c != nil && c.Data != nil {
		prefix, _ = dataCachePrefixAndExpire(c.Data)
	}
	return NewRedisCache(client, prefix, common, logger)
}

// NewDataSingleFlight 为 Data 层提供 singleflight.Group，与 Table 缓存配合防击穿。
func NewDataSingleFlight() *singleflight.Group {
	return &singleflight.Group{}
}

// NewDataTableCache 基于 redisCache 与 data.cache 配置组装表格缓存。
func NewDataTableCache(c *config.Config, redisCache Cache, sf *singleflight.Group, logger log.Logger) *Table {
	expireSec := int32(86400)
	if c != nil && c.Data != nil {
		_, expireSec = dataCachePrefixAndExpire(c.Data)
	}
	return NewTable(redisCache, sf, expireSec, log.NewHelper(log.With(logger, "layer", "data")))
}
