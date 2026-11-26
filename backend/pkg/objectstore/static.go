package objectstore

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	static "github.com/zhixingdx/oipivot-api/base/static/v1"
	"github.com/zhixingdx/oipivot/backend/pkg/utils"
	"github.com/go-kratos/kratos/contrib/otel/v3/tracing"
	"github.com/go-kratos/kratos/v3/middleware/metadata"
	"github.com/go-kratos/kratos/v3/middleware/recovery"
	k "github.com/go-kratos/kratos/v3/transport/grpc"
	"github.com/samber/lo"
	"google.golang.org/grpc"
)

// StaticStore 通过 base-static 服务实现的 ObjectStore
// 上传/下载通过预签名 URL + HTTP 直达 S3，gRPC 只交换凭证
type StaticStore struct {
	plusClient static.BaseStaticManagePlusClient
	httpClient *http.Client
	host       string
	service    string
	bucket     string
}

// StaticConfig StaticStore 配置
type StaticConfig struct {
	// base-static gRPC 地址（如 "discovery:///base-static"）
	Endpoint string
	// 调用方服务名（如 "shared-ai"），用于 GetUploadCredential 的 service 参数
	Service string
	// 可选，指定 bucket
	Bucket string
	// CDN host（用于 Host() 返回值和 URL 匹配判断）
	Host string
}

// NewStaticStore 创建 StaticStore
func NewStaticStore(cfg *StaticConfig) (*StaticStore, error) {
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("StaticStore 配置缺少 staticEndpoint")
	}

	conn, err := k.NewClient(
		context.Background(),
		k.WithEndpoint(cfg.Endpoint),
		k.WithMiddleware(
			recovery.Recovery(),
			metadata.Client(),
			tracing.Client(),
		),
		k.WithTimeout(30*time.Second),
		k.WithOptions(utils.MergeGRPCOptions(
			[]grpc.DialOption{grpc.WithStatsHandler(&tracing.ClientHandler{})},
			utils.GRPCKeepaliveOptions(),
		)...),
	)
	if err != nil {
		return nil, fmt.Errorf("连接 base-static 失败: %w", err)
	}

	host := cfg.Host
	if host != "" && !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		host = "https://" + host
	}
	host = strings.TrimSuffix(host, "/")

	return &StaticStore{
		plusClient: static.NewBaseStaticManagePlusClient(conn),
		httpClient: &http.Client{Timeout: 5 * time.Minute},
		host:       host,
		service:    cfg.Service,
		bucket:     cfg.Bucket,
	}, nil
}

// Upload 通过 base-static 签发的私有凭证上传对象，并返回可持久化的访问引用。
func (s *StaticStore) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	return s.UploadReader(ctx, key, bytes.NewReader(data), int64(len(data)), contentType)
}

// UploadReader 使用 base-static 签发的私有 URL 直接流式上传已知长度内容。
func (s *StaticStore) UploadReader(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error) {
	dir, filename := splitKeyToDirAndFilename(key)
	// 数据库仅保存私有文件引用，这里显式声明为私有上传，禁止落到公开目录策略。
	privateUpload := true

	reply, err := s.plusClient.GetUploadCredential(ctx, &static.GetUploadCredentialRequest{
		Bucket:   s.bucket,
		Service:  s.service,
		Dir:      dir,
		Type:     contentType,
		Filename: filename,
		Private:  lo.ToPtr(privateUpload),
	})
	if err != nil {
		return "", fmt.Errorf("获取上传凭证失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, reply.UploadURL, reader)
	if err != nil {
		return "", fmt.Errorf("创建上传请求失败: %w", err)
	}
	req.Header.Set("Content-Type", contentType)
	if size >= 0 {
		req.ContentLength = size
	}
	applySignedHeadersFromPresignedURL(req)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("上传到 S3 失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("上传到 S3 返回 %d: %s", resp.StatusCode, string(body))
	}

	if reply.AccessURL != "" {
		// base-static 可能为了调用方立即预览而返回带 auth_key 的临时地址；对象存储契约要求只持久化稳定 URL。
		return buildAccessURL(reply.AccessURL, ""), nil
	}
	// 兼容旧 base-static：无 accessURL 时按预签名路径 + 可选 StaticConfig.Host 拼浏览地址
	return buildAccessURL(reply.UploadURL, s.host), nil
}

// applySignedHeadersFromPresignedURL 根据预签名 URL 要求补齐必须请求头。
// 当前主要处理 x-amz-acl：若签名中包含该 header，但客户端未显式设置，则按私有对象补默认值，避免签名校验 403。
func applySignedHeadersFromPresignedURL(req *http.Request) {
	if req == nil || req.URL == nil {
		return
	}
	signedHeaders := req.URL.Query().Get("X-Amz-SignedHeaders")
	if signedHeaders == "" {
		signedHeaders = req.URL.Query().Get("x-amz-signedheaders")
	}
	if signedHeaders == "" {
		return
	}
	for _, header := range strings.Split(strings.ToLower(signedHeaders), ";") {
		switch strings.TrimSpace(header) {
		case "x-amz-acl":
			if req.Header.Get("x-amz-acl") == "" {
				req.Header.Set("x-amz-acl", "private")
			}
		}
	}
}

func (s *StaticStore) Download(ctx context.Context, rawURL string) ([]byte, error) {
	reply, err := s.plusClient.GetPrivateFileAccess(ctx, &static.GetPrivateFileAccessRequest{
		Bucket:     s.bucket,
		Service:    s.service,
		KeysOrURLs: []string{rawURL},
	})
	if err != nil {
		return nil, fmt.Errorf("获取下载凭证失败: %w", err)
	}

	presignedURL, ok := reply.Urls[rawURL]
	if !ok || presignedURL == "" {
		return nil, fmt.Errorf("未获取到预签名 URL: %s", rawURL)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, presignedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建下载请求失败: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("从 S3 下载失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("从 S3 下载返回 %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

func (s *StaticStore) Delete(ctx context.Context, rawURL string) error {
	_, err := s.plusClient.DeleteObjects(ctx, &static.DeleteObjectsRequest{
		Bucket:     s.bucket,
		Service:    s.service,
		KeysOrURLs: []string{rawURL},
	})
	if err != nil {
		return fmt.Errorf("删除对象失败: %w", err)
	}
	return nil
}

func (s *StaticStore) PresignGet(ctx context.Context, rawURL string, expiry time.Duration) (string, error) {
	m, err := s.PresignGetMany(ctx, []string{rawURL}, expiry)
	if err != nil {
		return "", err
	}
	presignedURL, ok := m[rawURL]
	if !ok || presignedURL == "" {
		return "", fmt.Errorf("未获取到预签名 URL: %s", rawURL)
	}
	return presignedURL, nil
}

// PresignGetMany 一次 RPC 批量获取预签名 URL，避免逐条 PresignGet 造成的多次网络往返。
func (s *StaticStore) PresignGetMany(ctx context.Context, urls []string, expiry time.Duration) (map[string]string, error) {
	if len(urls) == 0 {
		return nil, nil
	}
	reply, err := s.plusClient.GetPrivateFileAccess(ctx, &static.GetPrivateFileAccessRequest{
		Bucket:     s.bucket,
		Service:    s.service,
		KeysOrURLs: urls,
		Expires:    lo.ToPtr(int32(expiry.Seconds())),
	})
	if err != nil {
		return nil, fmt.Errorf("获取预签名 URL 失败: %w", err)
	}
	// reply.Urls 本身即为「原始 key/URL → 预签名」映射；base-static 只写入 KeysOrURLs 中成功的项，无额外 key，
	// 与先前按 urls 切片回表再组 map 的结果一致，无需二次遍历请求列表。
	if len(reply.Urls) == 0 {
		return nil, nil
	}
	return reply.Urls, nil
}

func (s *StaticStore) GetSize(ctx context.Context, rawURL string) (int64, error) {
	reply, err := s.plusClient.GetFileMetadata(ctx, &static.GetFileMetadataRequest{
		Bucket:   s.bucket,
		Service:  s.service,
		KeyOrURL: rawURL,
	})
	if err != nil {
		return 0, fmt.Errorf("获取文件大小失败: %w", err)
	}
	return reply.Size, nil
}

func (s *StaticStore) Host() string {
	return s.host
}

// splitKeyToDirAndFilename 将 key 拆分为目录和文件名
// 如 "shared-ai/message/content/2026/03/15/abc.json.gz" → ("shared-ai/message/content/2026/03/15", "abc.json.gz")
func splitKeyToDirAndFilename(key string) (string, string) {
	idx := strings.LastIndex(key, "/")
	if idx < 0 {
		return "", key
	}
	return key[:idx], key[idx+1:]
}

// buildAccessURL 从预签名上传 URL 和 cdnHost 构建干净的访问 URL
// 预签名 URL: https://s3.endpoint.com/tenantDir/service/dir/uuid.ext?X-Amz-...
// cdnHost: https://cdn.example.com
// 返回: https://cdn.example.com/tenantDir/service/dir/uuid.ext
func buildAccessURL(presignedURL, cdnHost string) string {
	parsed, err := url.Parse(presignedURL)
	if err != nil {
		return presignedURL
	}

	path := parsed.Path
	if cdnHost != "" {
		cdnHost = strings.TrimSuffix(cdnHost, "/")
		if !strings.HasPrefix(cdnHost, "http://") && !strings.HasPrefix(cdnHost, "https://") {
			cdnHost = "https://" + cdnHost
		}
		return cdnHost + path
	}

	return parsed.Scheme + "://" + parsed.Host + path
}
