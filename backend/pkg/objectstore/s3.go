package objectstore

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Config S3 对象存储配置
type S3Config struct {
	// 云服务商: aws / tencent / aliyun / huawei / volcengine
	Provider  string
	SecretID  string
	SecretKey string
	Region    string
	Bucket    string
	// 自定义 endpoint（可选，为空则按 provider+region+bucket 拼接）
	Endpoint string
	// CDN host（可选，用于构建对外访问 URL）
	Host string
}

// S3Store 基于 AWS S3 兼容 API 的对象存储实现
type S3Store struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string
	host          string
}

// NewS3Store 创建 S3 对象存储实例
func NewS3Store(cfg *S3Config) (*S3Store, error) {
	if cfg.SecretID == "" || cfg.SecretKey == "" {
		return nil, fmt.Errorf("S3 配置缺少 SecretID 或 SecretKey")
	}
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("S3 配置缺少 Bucket")
	}

	endpoint := cfg.Endpoint
	if endpoint == "" {
		endpoint = getEndpoint(cfg.Provider, cfg.Region, cfg.Bucket)
	}
	endpoint = strings.TrimSuffix(endpoint, "/")

	usePathStyle := shouldUsePathStyle(endpoint)

	signingEndpoint := endpoint
	if !strings.HasPrefix(signingEndpoint, "http://") && !strings.HasPrefix(signingEndpoint, "https://") {
		signingEndpoint = "https://" + signingEndpoint
	}

	client := s3.New(s3.Options{
		BaseEndpoint: aws.String(signingEndpoint),
		Region:       cfg.Region,
		Credentials:  aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(cfg.SecretID, cfg.SecretKey, "")),
		UsePathStyle: usePathStyle,
	})

	host := cfg.Host
	if host == "" {
		host = signingEndpoint
	}
	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		host = "https://" + host
	}
	host = strings.TrimSuffix(host, "/")

	return &S3Store{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucket:        cfg.Bucket,
		host:          host,
	}, nil
}

// Upload 将内存字节上传到 S3，并复用流式上传实现保持返回语义一致。
func (s *S3Store) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	return s.UploadReader(ctx, key, bytes.NewReader(data), int64(len(data)), contentType)
}

// UploadReader 使用已知长度 Reader 直接写入 S3；调用方应提供可重读文件以支持 SDK 重试。
func (s *S3Store) UploadReader(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
	}
	if size >= 0 {
		input.ContentLength = aws.Int64(size)
	}
	if _, err := s.client.PutObject(ctx, input); err != nil {
		return "", fmt.Errorf("上传对象失败 key=%s: %w", key, err)
	}
	return s.host + "/" + key, nil
}

func (s *S3Store) extractKey(rawURL string) string {
	key := extractKeyFromURL(rawURL)
	if hp := extractKeyFromURL(s.host); hp != "" {
		key = strings.TrimPrefix(key, hp+"/")
	}
	return key
}

func (s *S3Store) Download(ctx context.Context, rawURL string) ([]byte, error) {
	key := s.extractKey(rawURL)
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	output, err := s.client.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("下载对象失败 key=%s: %w", key, err)
	}
	defer output.Body.Close()
	return io.ReadAll(output.Body)
}

func (s *S3Store) Delete(ctx context.Context, rawURL string) error {
	key := s.extractKey(rawURL)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	if _, err := s.client.DeleteObject(ctx, input); err != nil {
		return fmt.Errorf("删除对象失败 key=%s: %w", key, err)
	}
	return nil
}

func (s *S3Store) PresignGet(ctx context.Context, rawURL string, expiry time.Duration) (string, error) {
	key := s.extractKey(rawURL)
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	output, err := s.presignClient.PresignGetObject(ctx, input, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", fmt.Errorf("生成预签名 URL 失败 key=%s: %w", key, err)
	}
	return output.URL, nil
}

func (s *S3Store) PresignGetMany(ctx context.Context, urls []string, expiry time.Duration) (map[string]string, error) {
	if len(urls) == 0 {
		return nil, nil
	}
	out := make(map[string]string, len(urls))
	for _, u := range urls {
		if err := ctx.Err(); err != nil {
			return out, err
		}
		presigned, err := s.PresignGet(ctx, u, expiry)
		if err == nil && presigned != "" {
			out[u] = presigned
		}
	}
	return out, nil
}

func (s *S3Store) GetSize(ctx context.Context, rawURL string) (int64, error) {
	key := s.extractKey(rawURL)
	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	output, err := s.client.HeadObject(ctx, input)
	if err != nil {
		return 0, fmt.Errorf("获取对象大小失败 key=%s: %w", key, err)
	}
	if output.ContentLength != nil {
		return *output.ContentLength, nil
	}
	return 0, nil
}

func (s *S3Store) Host() string {
	return s.host
}

// --- 内部工具函数（从 base/static 复用逻辑） ---

func getEndpoint(provider, region, bucket string) string {
	switch strings.ToLower(provider) {
	case "tencent":
		return "https://" + bucket + ".cos." + region + ".myqcloud.com"
	case "aliyun":
		return "https://oss-" + region + ".aliyuncs.com"
	case "huawei":
		return "https://obs." + region + ".myhuaweicloud.com"
	case "volcengine":
		return "https://tos." + region + ".volces.com"
	default:
		return "https://s3." + region + ".amazonaws.com"
	}
}

func shouldUsePathStyle(endpoint string) bool {
	if endpoint == "" {
		return false
	}
	return !strings.Contains(endpoint, "amazonaws.com") &&
		!strings.Contains(endpoint, "myqcloud.com") &&
		!strings.Contains(endpoint, "aliyuncs.com") &&
		!strings.Contains(endpoint, "myhuaweicloud.com") &&
		!strings.Contains(endpoint, "volces.com")
}

func extractKeyFromURL(rawURL string) string {
	if strings.HasPrefix(rawURL, "http://") || strings.HasPrefix(rawURL, "https://") {
		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			return rawURL
		}
		return strings.TrimPrefix(parsedURL.Path, "/")
	}
	return rawURL
}
