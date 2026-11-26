package seewo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// ContentTypeJSON JSON 内容类型
	ContentTypeJSON = "application/json;charset=utf-8"

	// 希沃开放平台请求头常量
	// 参考 Java SDK: com.seewo.open.sdk.DefaultSeewoClient.setParam
	// 源码位置: /private/tmp/com/seewo/open/sdk/DefaultSeewoClient.class

	// HeaderAppID 应用标识
	HeaderAppID = "x-sw-app-id"
	// HeaderReqPath 请求路径
	HeaderReqPath = "x-sw-req-path"
	// HeaderVersion 版本号
	HeaderVersion = "x-sw-version"
	// HeaderTimestamp 毫秒时间戳
	HeaderTimestamp = "x-sw-timestamp"
	// HeaderSignHeaders 参与签名的额外 header 列表（逗号分隔）
	HeaderSignHeaders = "x-sw-sign-headers"
	// HeaderSignType 签名方式
	HeaderSignType = "x-sw-sign-type"
	// HeaderContentMD5 请求体 MD5
	HeaderContentMD5 = "x-sw-content-md5"
	// HeaderSign 签名值
	HeaderSign = "x-sw-sign"
	// HeaderDataPermission 数据权限 ID
	HeaderDataPermission = "x-sw-data-permission"
	// HeaderReqID 请求 ID（响应头）
	HeaderReqID = "x-sw-req-id"
	// HeaderMessage 错误消息（响应头）
	HeaderMessage = "x-sw-message"
	// VersionValue 当前版本号
	VersionValue = "3"
	// HeaderPrefix 希沃系统头前缀
	HeaderPrefix = "x-sw-"
)

// Client 希沃开放平台 API 客户端
type Client struct {
	config     *Config
	httpClient *http.Client
}

// NewClient 创建新的客户端实例
func NewClient(config *Config) *Client {
	timeout := time.Duration(config.GetTimeout()) * time.Second

	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// Get 发送 GET 请求
func (c *Client) Get(ctx context.Context, path string, params map[string]string, result interface{}) error {
	return c.request(ctx, http.MethodGet, path, params, nil, result)
}

// PostJSON 发送 POST JSON 请求
func (c *Client) PostJSON(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.request(ctx, http.MethodPost, path, nil, body, result)
}

// request 发送 HTTP 请求
// 自动处理签名、请求头设置、响应解析
func (c *Client) request(ctx context.Context, method, path string, queryParams map[string]string, body interface{}, result interface{}) error {
	// 1. 序列化请求体
	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("序列化请求体失败: %w", err)
		}
	}

	// 2. 构建请求头（系统头 + 签名）
	headers := c.buildHeaders(path, queryParams, bodyBytes)

	// 3. 构建完整 URL
	requestURL := c.buildURL(path, queryParams)

	// 4. 创建 HTTP 请求
	var reqBody io.Reader
	if len(bodyBytes) > 0 {
		reqBody = bytes.NewBuffer(bodyBytes)
	}
	req, err := http.NewRequestWithContext(ctx, method, requestURL, reqBody)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 5. 设置请求头
	req.Header.Set("Content-Type", ContentTypeJSON)
	req.Header.Set("Accept", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 6. 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRequestFailed, err)
	}
	defer resp.Body.Close()

	// 7. 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应体失败: %w", err)
	}

	// 8. 检查 HTTP 状态码
	// 参考 Java SDK: DefaultSeewoClient.errorDetect
	if resp.StatusCode != http.StatusOK {
		reqID := resp.Header.Get(HeaderReqID)
		message := resp.Header.Get(HeaderMessage)
		if message == "" {
			message = resp.Status
		}
		return fmt.Errorf("%w: HTTP状态码=%d, message=%s, reqId=%s, 响应=%s",
			ErrRequestFailed, resp.StatusCode, message, reqID, string(respBody))
	}

	// 9. 解析响应
	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("%w: %v, 响应内容=%s", ErrInvalidResponse, err, string(respBody))
		}
	}

	return nil
}

// buildHeaders 构建请求头并生成签名
// 参考 Java SDK: com.seewo.open.sdk.DefaultSeewoClient.setParam
func (c *Client) buildHeaders(path string, queryParams map[string]string, bodyBytes []byte) map[string]string {
	headers := make(map[string]string)

	// 系统头
	headers[HeaderAppID] = c.config.AppID
	headers[HeaderReqPath] = path
	headers[HeaderVersion] = VersionValue
	headers[HeaderTimestamp] = fmt.Sprintf("%d", time.Now().UnixMilli())
	headers[HeaderSignType] = c.config.GetSignType()

	// 数据权限 ID（某些接口如 mdm-basis 系列需要此字段）
	if c.config.DataPermission != "" {
		headers[HeaderDataPermission] = c.config.DataPermission
	}

	// 请求体 MD5（如果有请求体）
	if len(bodyBytes) > 0 {
		headers[HeaderContentMD5] = GetContentMD5(bodyBytes)
	}

	// extraHeaders: 非 x-sw- 前缀的自定义 header（目前无）
	// 参考 Java SDK: DefaultSeewoClient.extraHeaders
	// 这里暂时为空，如需添加自定义 header 可扩展
	headers[HeaderSignHeaders] = ""

	// 生成签名
	// 参考 Java SDK: DefaultSeewoClient.sign
	// 合并 queries + headers，移除 x-sw-sign，然后签名
	signParams := make(map[string]string)
	for k, v := range queryParams {
		signParams[k] = v
	}
	for k, v := range headers {
		signParams[k] = v
	}
	// x-sw-sign 不参与签名
	delete(signParams, HeaderSign)

	sign, _ := SignRequest(signParams, c.config.AppSecret, c.config.GetSignType())
	headers[HeaderSign] = sign

	return headers
}

// buildURL 构建请求 URL
func (c *Client) buildURL(path string, queryParams map[string]string) string {
	baseURL := strings.TrimRight(c.config.API, "/")
	requestURL := baseURL + path

	if len(queryParams) == 0 {
		return requestURL
	}

	values := url.Values{}
	for k, v := range queryParams {
		values.Set(k, v)
	}

	return requestURL + "?" + values.Encode()
}
