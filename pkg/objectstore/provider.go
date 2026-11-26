package objectstore

import (
	"github.com/zhixingdx/public/api/config"
)

// NewFromConfig 根据配置创建 ObjectStore 实例
// 在 wire 中使用: wire.NewSet(objectstore.NewFromConfig)
func NewFromConfig(c *config.Config) (ObjectStore, error) {
	if c.Data == nil || c.Data.ObjectStore == nil {
		return nil, nil
	}
	osc := c.Data.ObjectStore

	switch osc.Mode {
	case "static":
		staticEndpoint := resolveStaticEndpoint(c, osc)
		return NewStaticStore(&StaticConfig{
			Endpoint: staticEndpoint,
			Service:  osc.Service,
			Host:     osc.Host,
		})
	default:
		return NewS3Store(&S3Config{
			Provider:  osc.Provider,
			SecretID:  osc.SecretID,
			SecretKey: osc.SecretKey,
			Region:    osc.Region,
			Bucket:    osc.Bucket,
			Endpoint:  osc.Endpoint,
			Host:      osc.Host,
		})
	}
}

// resolveStaticEndpoint 解析 static 模式下的 endpoint。
// 优先使用 data.objectStore.staticEndpoint；
// 若未配置，则回退使用 services["base-static"]。
func resolveStaticEndpoint(c *config.Config, osc *config.ObjectStoreConf) string {
	if osc != nil && osc.StaticEndpoint != "" {
		return osc.StaticEndpoint
	}
	if c != nil && c.Services != nil {
		return c.Services["base-static"]
	}
	return ""
}
