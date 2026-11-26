package seewo

const (
	// DefaultTimeout 默认请求超时时间（秒）
	DefaultTimeout = 30
	// DefaultSignType 默认签名方式（HMAC-MD5）
	DefaultSignType = SignMethodHMAC
)

// Config 希沃开放平台 API 客户端配置
type Config struct {
	// API 基础地址（例如：https://openapi.seewo.com）
	API string `yaml:"api"`
	// AppID 应用标识，对应请求头 x-sw-app-id
	AppID string `yaml:"appid"`
	// AppSecret 应用密钥，用于生成签名
	AppSecret string `yaml:"appSecret"`
	// DataPermission 数据权限 ID，对应请求头 x-sw-data-permission
	// 某些接口（如 mdm-basis 系列）需要此字段才能访问
	DataPermission string `yaml:"dataPermission,omitempty"`
	// Timeout 请求超时时间（秒，默认 30）
	Timeout int `yaml:"timeout,omitempty"`
	// SignType 签名方式：hmac 或 md5，默认 hmac
	SignType string `yaml:"signType,omitempty"`
}

// Validate 验证配置是否有效
func (c *Config) Validate() error {
	if c.API == "" {
		return ErrConfigInvalidAPI
	}
	if c.AppID == "" {
		return ErrConfigInvalidAppID
	}
	if c.AppSecret == "" {
		return ErrConfigInvalidAppSecret
	}
	return nil
}

// GetSignType 获取签名方式，默认 hmac
func (c *Config) GetSignType() string {
	if c.SignType == "" {
		return DefaultSignType
	}
	return c.SignType
}

// GetTimeout 获取超时时间（秒），默认 30
func (c *Config) GetTimeout() int {
	if c.Timeout <= 0 {
		return DefaultTimeout
	}
	return c.Timeout
}
