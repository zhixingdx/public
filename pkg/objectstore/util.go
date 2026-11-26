package objectstore

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Compress 使用 gzip 压缩数据
func Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("创建 gzip writer 失败: %w", err)
	}
	if _, err = w.Write(data); err != nil {
		return nil, fmt.Errorf("gzip 写入失败: %w", err)
	}
	if err = w.Close(); err != nil {
		return nil, fmt.Errorf("gzip 关闭失败: %w", err)
	}
	return buf.Bytes(), nil
}

// Decompress 解压 gzip 数据
func Decompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("创建 gzip reader 失败: %w", err)
	}
	defer r.Close()
	result, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("gzip 读取失败: %w", err)
	}
	return result, nil
}

// DetectContentType 推断内容类型，返回文件扩展名和 MIME 类型
// configured 不为空时直接使用配置值
func DetectContentType(content string, configured string) (ext string, mime string) {
	if configured != "" {
		return contentTypeToExtMime(configured)
	}
	trimmed := strings.TrimSpace(content)
	switch {
	case strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "["):
		return ".json", "application/json"
	case strings.HasPrefix(trimmed, "<"):
		return ".html", "text/html"
	case strings.HasPrefix(trimmed, "#") || strings.Contains(trimmed[:min(200, len(trimmed))], "\n##"):
		return ".md", "text/markdown"
	default:
		return ".txt", "text/plain"
	}
}

// contentTypeToExtMime 将配置中的内容类型归一化为扩展名和标准 MIME。
func contentTypeToExtMime(ct string) (string, string) {
	switch strings.ToLower(strings.TrimSpace(ct)) {
	case "json":
		return ".json", "application/json"
	case "jsonl", "ndjson", "application/jsonl", "application/x-ndjson", "application/ndjson":
		return ".jsonl", "application/x-ndjson"
	case "html":
		return ".html", "text/html"
	case "md":
		return ".md", "text/markdown"
	case "txt":
		return ".txt", "text/plain"
	default:
		return ".txt", "text/plain"
	}
}

// BuildObjectKey 构建对象存储 key
// 格式: {dir}/{date}/{uuid}{ext}
func BuildObjectKey(dir, ext string) string {
	date := time.Now().Format("2006/01/02")
	return fmt.Sprintf("%s/%s/%s%s", dir, date, uuid.New().String(), ext)
}

// BuildDefaultDir 构建默认目录前缀
// 格式: {serviceName}/{tableName}/{fieldName}
func BuildDefaultDir(serviceName, tableName, fieldName string) string {
	return fmt.Sprintf("%s/%s/%s", serviceName, tableName, fieldName)
}

// PackRef 打包引用：将 URL 和文件大小打包为 "size|url" 格式存入 DB
func PackRef(url string, size int64) string {
	return fmt.Sprintf("%d|%s", size, url)
}

// UnpackRef 解包引用：从 "size|url" 格式提取 URL 和文件大小
// 只切第一个 "|"，URL 部分完整保留
func UnpackRef(value string) (url string, size int64) {
	parts := strings.SplitN(value, "|", 2)
	if len(parts) == 2 {
		n, _ := strconv.ParseInt(parts[0], 10, 64)
		return parts[1], n
	}
	return value, 0
}

// IsObjectStoreURL 判断值是否为对象存储 URL（兼容 packed ref 格式）
func IsObjectStoreURL(value string, host string) bool {
	if value == "" {
		return false
	}
	rawURL, _ := UnpackRef(value)
	if rawURL == "" {
		return false
	}
	if host == "" {
		u, err := url.Parse(rawURL)
		return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
	}
	normalizedHost := strings.TrimSuffix(host, "/")
	if !strings.HasPrefix(normalizedHost, "http://") && !strings.HasPrefix(normalizedHost, "https://") {
		normalizedHost = "https://" + normalizedHost
	}
	return strings.HasPrefix(strings.TrimSuffix(rawURL, "/"), normalizedHost)
}
