// Package objectstore 大字段对象存储抽象
// 提供统一的对象存储接口，Data 层通过此接口上传/下载/删除大字段内容
package objectstore

import (
	"context"
	"io"
	"time"
)

// ObjectStore 对象存储接口
type ObjectStore interface {
	// Upload 上传内容到对象存储，返回不带授权参数的访问 URL
	Upload(ctx context.Context, key string, data []byte, contentType string) (url string, err error)

	// Download 根据不带授权参数的 URL 下载内容
	Download(ctx context.Context, url string) ([]byte, error)

	// Delete 根据不带授权参数的 URL 删除对象
	Delete(ctx context.Context, url string) error

	// PresignGet 根据不带授权参数的 URL 生成预签名访问地址（带临时授权参数）
	PresignGet(ctx context.Context, url string, expiry time.Duration) (presignedURL string, err error)

	// PresignGetMany 批量生成预签名访问地址：key 为入参原始 URL，value 为预签名 URL。
	// 单个 URL 失败时该 key 不会出现在 map 中（与逐条 PresignGet 且忽略错误时的有效集合一致）。
	// 对经由 base-static 的实现应合并为单次 RPC，避免逐条请求造成的 N 次往返。
	PresignGetMany(ctx context.Context, urls []string, expiry time.Duration) (urlToPresigned map[string]string, err error)

	// GetSize 获取对象大小(字节)
	GetSize(ctx context.Context, url string) (int64, error)

	// Host 返回配置的对象存储 Host 前缀，用于判断字段值是否为对象存储 URL
	// 例如: https://bucket.s3.region.com
	Host() string
}

// StreamingUploader 是大对象流式写入的可选扩展，不破坏已有 ObjectStore 实现的兼容性。
type StreamingUploader interface {
	// UploadReader 从 Reader 流式上传内容；size 必须是已知长度，避免实现层把大对象完整缓冲到内存。
	UploadReader(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (url string, err error)
}
