package api

import "context"

// Client 客户端接口，避免循环导入
// 所有 API 模块通过此接口调用底层 HTTP 客户端
type Client interface {
	// Get 发送 GET 请求
	Get(ctx context.Context, path string, params map[string]string, result interface{}) error
	// PostJSON 发送 POST JSON 请求
	PostJSON(ctx context.Context, path string, body interface{}, result interface{}) error
}
