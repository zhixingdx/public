// Package s3url 提供 S3 私有 URL 的自动预签名处理工具
// 支持从字符串/JSON 中提取 URL、清除授权参数（入库）、批量预签名替换（出库）
package s3url

import (
	"encoding/json"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

// s3AuthParams 入库时需从 URL 查询串剥离的鉴权相关参数
// 覆盖：S3 兼容协议、GCS 签名、常见 CDN Type A/sign、CloudFront、Akamai、Azure Blob SAS 等
var s3AuthParams = []string{
	// --- AWS SigV4 / S3 兼容 ---
	"X-Amz-Algorithm",
	"X-Amz-Credential",
	"X-Amz-Date",
	"X-Amz-Expires",
	"X-Amz-Signature",
	"X-Amz-SignedHeaders",
	// AWS SDK / S3 兼容客户端在预签名 URL 上附加的操作参数（非资源标识，出库替换时需与 SigV4 一并剥离）
	"x-id",
	"X-Amz-Security-Token",
	"AWSAccessKeyId",
	"Signature",
	"Expires",
	// --- Google Cloud Storage 查询字符串签名（V4）---
	"X-Goog-Algorithm",
	"X-Goog-Credential",
	"X-Goog-Date",
	"X-Goog-Expires",
	"X-Goog-SignedHeaders",
	"X-Goog-Signature",
	// --- Azure Blob SAS（Shared Access Signature）---
	"sv",
	"sr",
	"st",
	"se",
	"sp",
	"spr",
	"sip",
	// Azure SAS / Fastly URL 签名等共用短名 sig
	"sig",
	// --- CDN Type A（腾讯、阿里可配置、华为等，值形态 timestamp-rand-uid-md5）---
	"auth_key",
	"AUTH_KEY",
	// --- CDN sign / HMAC（百度、火山、自建、部分网宿等）---
	"sign",
	// --- CloudFront 标准查询参数签名 ---
	"Policy",
	"Key-Pair-Id",
	// --- Akamai EdgeAuth / Token Auth ---
	"__hdnea__",
	"hdnea",
	"hdnts",
	"__token__",
	// --- 又拍云融合云下载防盗链 ---
	"_upt",
	// --- 网宿等厂商常见时间+密钥拼接鉴权 ---
	"wsSecret",
	"wsTime",
	// --- 腾讯云 COS ---
	"x-cos-security-token",
	"q-sign-algorithm",
	"q-ak",
	"q-sign-time",
	"q-key-time",
	"q-header-list",
	"q-url-param-list",
	"q-signature",
	// --- 兼容阿里云 OSS 预签名 ---
	"OSSAccessKeyId",
	"SecurityToken",
}

// urlRegexp 匹配 HTTP(S) 对象/文件 URL（query 部分可选）
// 约束：authority 之后必须有路径（至少一个 /），避免把 https://example.com 的 .com 当成扩展名；
// 扩展名：路径中最后一个点之后为「后缀」——首字符为字母或数字（支持 .7z），其余为字母数字与 _+-（后缀内不含 .）。
// 路径段使用贪婪匹配 + 回溯，使「多点」文件名（如 file.name.tar.gz）仍落在最后一层后缀（.gz），而非误把 .name 当扩展名。
// 不含空白、引号、尖括号和 JSON 转义反斜杠；长度受限。
var urlRegexp = regexp.MustCompile(
	`https?://[^/\s"'<>\\` + "`" + `]+/` +
		`[^\s"'<>\\` + "`" + `]+` +
		`\.[a-zA-Z0-9][a-zA-Z0-9_+-]{0,127}` +
		`(?:\?[^\s"'<>\\` + "`" + `]*)?`,
)

// StripAuthParams 清除 URL 查询串中的鉴权相关参数，保留干净的原始 URL
// 用于入库时清理数据，避免存储临时授权信息（具体参数见 s3AuthParams）
func StripAuthParams(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	query := parsed.Query()
	if len(query) == 0 {
		return rawURL
	}
	changed := false
	for _, param := range s3AuthParams {
		if query.Has(param) {
			query.Del(param)
			changed = true
		}
	}
	if !changed {
		return rawURL
	}
	parsed.RawQuery = query.Encode()
	if len(query) == 0 {
		parsed.RawQuery = ""
	}
	return parsed.String()
}

// ExtractURLs 从文本中提取所有匹配的 URL
func ExtractURLs(text string) []string {
	return urlRegexp.FindAllString(text, -1)
}

// MatchPrefix 检查 URL 是否匹配指定的域名/路径前缀列表
// prefixes 为空时不匹配任何 URL（表示未配置预签名范围，不参与处理）
func MatchPrefix(rawURL string, prefixes []string) bool {
	if len(prefixes) == 0 {
		return false
	}
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Hostname() == "" {
		return false
	}
	targetHost := strings.ToLower(parsed.Hostname())
	targetPath := parsed.Path
	if targetPath == "" {
		targetPath = "/"
	}
	for _, prefix := range prefixes {
		host, path, ok := parsePrefix(prefix)
		if !ok {
			continue
		}
		if targetHost != host {
			continue
		}
		if path == "/" || targetPath == path {
			return true
		}
		if strings.HasSuffix(path, "/") && strings.HasPrefix(targetPath, path) {
			return true
		}
		if strings.HasPrefix(targetPath, path+"/") {
			return true
		}
	}
	return false
}

func parsePrefix(prefix string) (host string, path string, ok bool) {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return "", "", false
	}
	if !strings.Contains(prefix, "://") {
		prefix = "https://" + prefix
	}
	parsed, err := url.Parse(prefix)
	if err != nil || parsed.Hostname() == "" {
		return "", "", false
	}
	host = strings.ToLower(parsed.Hostname())
	path = parsed.Path
	if path == "" {
		path = "/"
	}
	return host, path, true
}

// CollectURLsForPresign 从文本中收集需要预签名的 URL（去重）
// 已包含授权参数的 URL 会先清理，再去重
func CollectURLsForPresign(text string, prefixes []string) []string {
	if len(prefixes) == 0 {
		return nil
	}
	urls := ExtractURLs(text)
	seen := make(map[string]bool, len(urls))
	var result []string
	for _, u := range urls {
		clean := StripAuthParams(u)
		if !MatchPrefix(clean, prefixes) {
			continue
		}
		if seen[clean] {
			continue
		}
		seen[clean] = true
		result = append(result, clean)
	}
	return result
}

// ReplaceURLs 在文本中批量替换 URL。
// urlMap 有两种用法：① 预签名出库 key 为已 StripAuthParams 的干净 URL、value 为预签名/CDN URL；
// ② 入库 StripAllAuthParams 时 key 为原文中带鉴权参数的完整 URL、value 为去掉参数后的 URL。
// 对 ExtractURLs 得到的每条完整 URL u，依次尝试 urlMap[u]、urlMap[StripAuthParams(u)]，
// 命中则将 u 整段替换（避免仅用干净 key 做子串替换时，SigV4 与 x-id 等参数顺序导致无法匹配）。
func ReplaceURLs(text string, urlMap map[string]string) string {
	if len(urlMap) == 0 {
		return text
	}
	urls := ExtractURLs(text)
	seen := make(map[string]bool, len(urls))
	unique := make([]string, 0, len(urls))
	for _, u := range urls {
		if !seen[u] {
			seen[u] = true
			unique = append(unique, u)
		}
	}
	sort.Slice(unique, func(i, j int) bool {
		return len(unique[i]) > len(unique[j])
	})
	for _, u := range unique {
		repl, ok := urlMap[u]
		if !ok {
			repl, ok = urlMap[StripAuthParams(u)]
		}
		if !ok || repl == "" || u == repl {
			continue
		}
		text = strings.ReplaceAll(text, u, repl)
	}
	return text
}

// StripAllAuthParams 清除文本中所有 URL 的 S3 授权参数
// 用于入库前批量清理字段内容
func StripAllAuthParams(text string, prefixes []string) string {
	if len(prefixes) == 0 {
		return text
	}
	urls := ExtractURLs(text)
	if len(urls) == 0 {
		return text
	}
	urlMap := make(map[string]string)
	for _, u := range urls {
		if !MatchPrefix(u, prefixes) {
			continue
		}
		clean := StripAuthParams(u)
		if clean != u {
			urlMap[u] = clean
		}
	}
	return ReplaceURLs(text, urlMap)
}

// ProcessJSONFieldForStore 处理 JSON 字段的入库清理
// 将 JSON 序列化为字符串 → 清除授权参数 → 反序列化回原类型
func ProcessJSONFieldForStore(jsonData interface{}, prefixes []string) (interface{}, error) {
	data, err := MarshalNoEscape(jsonData)
	if err != nil {
		return jsonData, err
	}
	text := string(data)
	cleaned := StripAllAuthParams(text, prefixes)
	if cleaned == text {
		return jsonData, nil
	}
	var result interface{}
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return jsonData, err
	}
	return result, nil
}

// ProcessJSONFieldForRead 处理 JSON 字段的出库预签名
// 将 JSON 序列化为字符串 → 替换预签名 URL → 反序列化回原类型
func ProcessJSONFieldForRead(jsonData interface{}, urlMap map[string]string) (interface{}, error) {
	if len(urlMap) == 0 {
		return jsonData, nil
	}
	data, err := MarshalNoEscape(jsonData)
	if err != nil {
		return jsonData, err
	}
	text := string(data)
	replaced := ReplaceURLs(text, urlMap)
	if replaced == text {
		return jsonData, nil
	}
	var result interface{}
	if err := json.Unmarshal([]byte(replaced), &result); err != nil {
		return jsonData, err
	}
	return result, nil
}

// MarshalNoEscape JSON 序列化，不转义 HTML 字符（避免 & < > 被转义导致 URL 匹配失败）
func MarshalNoEscape(v interface{}) ([]byte, error) {
	buf := &strings.Builder{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	s := buf.String()
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	return []byte(s), nil
}
