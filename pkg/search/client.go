// Package search 封装 OpenSearch 客户端，提供业务语义的搜索操作接口
// 对外不暴露 opensearch-go SDK 细节，上层通过 Client 完成文档索引、删除和搜索
package search

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

// Client OpenSearch 客户端封装
type Client struct {
	client      *opensearchapi.Client
	prefix      string
	serviceName string
}

// NewClient 创建搜索客户端
// prefix：环境前缀（如 "dev-"、"prod-"），用于环境隔离
// serviceName：服务命名空间（如 "base-iam"），用于避免不同服务相同表名冲突；为空时退化为仅加环境前缀
func NewClient(addr string, addrs []string, user, password, prefix, serviceName string) (*Client, error) {
	addresses := addrs
	if len(addresses) == 0 && addr != "" {
		addresses = []string{addr}
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("OpenSearch 地址不能为空")
	}

	osCfg := opensearch.Config{
		Addresses: addresses,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		},
	}
	if user != "" {
		osCfg.Username = user
		osCfg.Password = password
	}

	apiClient, err := opensearchapi.NewClient(opensearchapi.Config{Client: osCfg})
	if err != nil {
		return nil, fmt.Errorf("创建 OpenSearch 客户端失败: %w", err)
	}

	return &Client{client: apiClient, prefix: prefix, serviceName: serviceName}, nil
}

// collectBulkAPIErrors OpenSearch Bulk 在 HTTP 200 时仍可能返回 body 内 errors=true，此处汇总失败项便于排查 mapping/权限等问题
func collectBulkAPIErrors(operation, index string, resp *opensearchapi.BulkResp) error {
	if resp == nil || !resp.Errors {
		return nil
	}
	var b strings.Builder
	failN := 0
	for _, item := range resp.Items {
		for op, d := range item {
			if d.Error != nil {
				failN++
				if b.Len() < 900 {
					fmt.Fprintf(&b, "[%s id=%s] %s: %s; ", op, d.ID, d.Error.Type, d.Error.Reason)
				}
				continue
			}
			// 删除文档不存在时常见 404，不视为失败
			if op == "delete" && d.Status == http.StatusNotFound {
				continue
			}
			if d.Status >= 300 {
				failN++
				if b.Len() < 900 {
					fmt.Fprintf(&b, "[%s id=%s] status=%d; ", op, d.ID, d.Status)
				}
			}
		}
	}
	if failN == 0 {
		return fmt.Errorf("%s 索引 %s Bulk 响应 errors=true 但未解析到失败项", operation, index)
	}
	return fmt.Errorf("%s 索引 %s Bulk 存在 %d 项失败: %s", operation, index, failN, strings.TrimSuffix(strings.TrimSpace(b.String()), ";"))
}

// GetPrefix 返回当前环境索引前缀（如 "dev-"、"prod-"、""）
// 仅用于拼接展示或表达式模板透出；业务代码请优先使用 GetIndexName，后者同时支持多索引/排除语法
func (c *Client) GetPrefix() string {
	return c.prefix
}

// GetServiceName 返回当前服务命名空间（如 "base-iam"），未配置时为空
func (c *Client) GetServiceName() string {
	return c.serviceName
}

// GetIndexName 获取带前缀的索引名。
//
// 拼接规则（按段独立处理；serviceName 为空时仅加 prefix，保持旧行为）：
//   - 段内不含 '-'（视为本服务短表名）：prefix + serviceName + "-" + name
//     例："staff" → "dev-base-iam-staff"
//   - 段内含 '-'（视为完整索引名，形如 svcType-svcName-table 或通配符 logs-*）：只加 prefix
//     例："base-iam-staff" → "dev-base-iam-staff"；"logs-*" → "dev-logs-*"
//
// 多索引/排除语法（OpenSearch multi-index）：
//   - 逗号分隔多个索引；每段独立走上述规则
//   - 以 '-' 开头的段表示排除，剥离首个 '-' 后按上述规则处理，最后补回 '-'
//     例："a,-b" + serviceName="base-iam" → "dev-base-iam-a,-dev-base-iam-b"
//
// 跨服务查询：在段内直接写完整索引名即可绕过自服务命名空间
//
//	例：base-stat 里查 base-iam 的 staff → "base-iam-staff" → "dev-base-iam-staff"
func (c *Client) GetIndexName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return c.prefix
	}
	parts := strings.Split(name, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if strings.HasPrefix(p, "-") {
			rest := strings.TrimSpace(p[1:])
			if rest == "" {
				continue
			}
			out = append(out, "-"+c.buildSingleIndexName(rest))
			continue
		}
		out = append(out, c.buildSingleIndexName(p))
	}
	if len(out) == 0 {
		return c.prefix
	}
	return strings.Join(out, ",")
}

// buildSingleIndexName 按"是否含 -"决定是否补服务命名空间，再统一加环境前缀
// 调用方需保证入参已去除首尾空白且非空、非排除前缀
func (c *Client) buildSingleIndexName(name string) string {
	if c.serviceName == "" || strings.Contains(name, "-") {
		return c.prefix + name
	}
	return c.prefix + c.serviceName + "-" + name
}

// IndexExists 检查索引是否存在
func (c *Client) IndexExists(ctx context.Context, name string) (bool, error) {
	resp, err := c.client.Indices.Exists(ctx, opensearchapi.IndicesExistsReq{
		Indices: []string{name},
	})
	if err != nil {
		// 404 通过 error 返回，表示索引不存在
		if resp != nil && resp.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("检查索引 %s 是否存在失败: %w", name, err)
	}
	return resp.StatusCode == 200, nil
}

// CountDocuments 统计索引内文档条数（match_all），用于增量同步判断空索引
func (c *Client) CountDocuments(ctx context.Context, index string) (int64, error) {
	resp, err := c.client.Indices.Count(ctx, &opensearchapi.IndicesCountReq{
		Indices: []string{index},
		Body:    strings.NewReader(`{"query":{"match_all":{}}}`),
	})
	if err != nil {
		if resp != nil {
			ins := resp.Inspect()
			if ins.Response != nil && ins.Response.StatusCode == http.StatusNotFound {
				return 0, nil
			}
		}
		return 0, fmt.Errorf("统计索引 %s 文档数失败: %w", index, err)
	}
	return int64(resp.Count), nil
}

// CreateIndex 创建索引
// body 为索引的 settings + mappings 定义
func (c *Client) CreateIndex(ctx context.Context, name string, body map[string]any) error {
	bodyReader, err := toJSONReader(body)
	if err != nil {
		return err
	}
	_, err = c.client.Indices.Create(ctx, opensearchapi.IndicesCreateReq{
		Index: name,
		Body:  bodyReader,
	})
	if err != nil {
		return fmt.Errorf("创建索引 %s 失败: %w", name, err)
	}
	return nil
}

// IndexDocument 索引单个文档
func (c *Client) IndexDocument(ctx context.Context, index, id string, doc any) error {
	bodyReader, err := toJSONReader(doc)
	if err != nil {
		return err
	}
	_, err = c.client.Index(ctx, opensearchapi.IndexReq{
		Index:      index,
		DocumentID: id,
		Body:       bodyReader,
	})
	if err != nil {
		return fmt.Errorf("索引文档 %s/%s 失败: %w", index, id, err)
	}
	return nil
}

// IndexDocumentVersioned 使用 external_gte 版本栅栏写入文档；旧版本冲突视为已被更新事实覆盖。
func (c *Client) IndexDocumentVersioned(ctx context.Context, index, id string, version int64, doc any) error {
	if version < 1 {
		return fmt.Errorf("索引文档 %s/%s 的外部版本必须大于零", index, id)
	}
	bodyReader, err := toJSONReader(doc)
	if err != nil {
		return err
	}
	apiVersion := int(version)
	resp, err := c.client.Index(ctx, opensearchapi.IndexReq{
		Index: index, DocumentID: id, Body: bodyReader,
		Params: opensearchapi.IndexParams{Version: &apiVersion, VersionType: "external_gte"},
	})
	if err == nil {
		return nil
	}
	if resp != nil && resp.Inspect().Response != nil && resp.Inspect().Response.StatusCode == http.StatusConflict &&
		strings.Contains(err.Error(), "version_conflict_engine_exception") {
		return nil
	}
	return fmt.Errorf("按外部版本索引文档 %s/%s 失败: %w", index, id, err)
}

// DeleteDocument 删除单个文档（文档不存在时不报错）
func (c *Client) DeleteDocument(ctx context.Context, index, id string) error {
	_, err := c.client.Document.Delete(ctx, opensearchapi.DocumentDeleteReq{
		Index:      index,
		DocumentID: id,
	})
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "404") || strings.Contains(errStr, "not_found") {
			return nil
		}
		return fmt.Errorf("删除文档 %s/%s 失败: %w", index, id, err)
	}
	return nil
}

// UpdateDocument 局部更新单个文档，body 为 OpenSearch update API 请求体。
func (c *Client) UpdateDocument(ctx context.Context, index, id string, body map[string]any) error {
	bodyReader, err := toJSONReader(body)
	if err != nil {
		return err
	}
	_, err = c.client.Update(ctx, opensearchapi.UpdateReq{
		Index:      index,
		DocumentID: id,
		Body:       bodyReader,
	})
	if err != nil {
		return fmt.Errorf("更新文档 %s/%s 失败: %w", index, id, err)
	}
	return nil
}

// SearchResult 搜索结果
type SearchResult struct {
	Total        int64
	Hits         []SearchHit
	MaxScore     float64
	Aggregations json.RawMessage
	// 完整性信号交由调用方按业务契约处理，避免统计接口把部分结果当作完整事实缓存。
	TimedOut     bool
	FailedShards int
	Errors       bool
}

// SearchHit 单条搜索命中
type SearchHit struct {
	ID     string
	Score  float64
	Source json.RawMessage
	Sort   []any
}

// GetDocument 读取单个文档；文档不存在时返回 found=false、source=nil。
func (c *Client) GetDocument(ctx context.Context, index, id string) (json.RawMessage, bool, error) {
	resp, err := c.client.Document.Get(ctx, opensearchapi.DocumentGetReq{
		Index:      index,
		DocumentID: id,
	})
	if err != nil {
		return nil, false, fmt.Errorf("读取文档 %s/%s 失败: %w", index, id, err)
	}
	if !resp.Found {
		return nil, false, nil
	}
	return resp.Source, true, nil
}

// Search 执行搜索
func (c *Client) Search(ctx context.Context, index string, body map[string]any) (*SearchResult, error) {
	bodyReader, err := toJSONReader(body)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Search(ctx, &opensearchapi.SearchReq{
		Indices: []string{index},
		Body:    bodyReader,
	})
	if err != nil {
		return nil, fmt.Errorf("搜索索引 %s 失败: %w", index, err)
	}

	result := &SearchResult{
		Total:        int64(resp.Hits.Total.Value),
		Aggregations: resp.Aggregations,
		TimedOut:     resp.Timeout,
		FailedShards: resp.Shards.Failed,
		Errors:       resp.Errors,
	}
	for _, hit := range resp.Hits.Hits {
		sh := SearchHit{
			ID:     hit.ID,
			Score:  float64(hit.Score),
			Source: hit.Source,
			Sort:   hit.Sort,
		}
		result.Hits = append(result.Hits, sh)
	}
	return result, nil
}

// BulkIndex 批量索引文档
// docs: key 为文档 ID，value 为文档内容
func (c *Client) BulkIndex(ctx context.Context, index string, docs map[string]any) error {
	if len(docs) == 0 {
		return nil
	}
	var bulkBody strings.Builder
	for id, doc := range docs {
		action := map[string]any{
			"index": map[string]any{
				"_index": index,
				"_id":    id,
			},
		}
		actionJSON, err := json.Marshal(action)
		if err != nil {
			return fmt.Errorf("序列化 bulk action 失败: %w", err)
		}
		docJSON, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("序列化文档 %s 失败: %w", id, err)
		}
		bulkBody.Write(actionJSON)
		bulkBody.WriteByte('\n')
		bulkBody.Write(docJSON)
		bulkBody.WriteByte('\n')
	}

	resp, err := c.client.Bulk(ctx, opensearchapi.BulkReq{
		Body: strings.NewReader(bulkBody.String()),
	})
	if err != nil {
		return fmt.Errorf("批量索引到 %s 失败: %w", index, err)
	}
	if apiErr := collectBulkAPIErrors("BulkIndex", index, resp); apiErr != nil {
		return apiErr
	}
	return nil
}

// BulkUpdate 批量局部更新文档，updates 的 key 为文档 ID，value 为 OpenSearch update API 请求体。
func (c *Client) BulkUpdate(ctx context.Context, index string, updates map[string]any) error {
	if len(updates) == 0 {
		return nil
	}
	var bulkBody strings.Builder
	for id, body := range updates {
		action := map[string]any{
			"update": map[string]any{
				"_index": index,
				"_id":    id,
			},
		}
		actionJSON, err := json.Marshal(action)
		if err != nil {
			return fmt.Errorf("序列化 bulk update action 失败: %w", err)
		}
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("序列化更新文档 %s 失败: %w", id, err)
		}
		bulkBody.Write(actionJSON)
		bulkBody.WriteByte('\n')
		bulkBody.Write(bodyJSON)
		bulkBody.WriteByte('\n')
	}

	resp, err := c.client.Bulk(ctx, opensearchapi.BulkReq{
		Body: strings.NewReader(bulkBody.String()),
	})
	if err != nil {
		return fmt.Errorf("批量更新 %s 失败: %w", index, err)
	}
	if apiErr := collectBulkAPIErrors("BulkUpdate", index, resp); apiErr != nil {
		return apiErr
	}
	return nil
}

// BulkDelete 批量删除文档
func (c *Client) BulkDelete(ctx context.Context, index string, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	var bulkBody strings.Builder
	for _, id := range ids {
		action := map[string]any{
			"delete": map[string]any{
				"_index": index,
				"_id":    id,
			},
		}
		actionJSON, err := json.Marshal(action)
		if err != nil {
			return fmt.Errorf("序列化 bulk delete action 失败: %w", err)
		}
		bulkBody.Write(actionJSON)
		bulkBody.WriteByte('\n')
	}

	resp, err := c.client.Bulk(ctx, opensearchapi.BulkReq{
		Body: strings.NewReader(bulkBody.String()),
	})
	if err != nil {
		return fmt.Errorf("批量删除 %s 失败: %w", index, err)
	}
	if apiErr := collectBulkAPIErrors("BulkDelete", index, resp); apiErr != nil {
		return apiErr
	}
	return nil
}

// ScrollAllIDs 滚动查询索引中所有文档的 ID（用于全量同步时清理孤儿数据）
func (c *Client) ScrollAllIDs(ctx context.Context, index string, query map[string]any) ([]string, error) {
	body := map[string]any{
		"size":    1000,
		"_source": false,
	}
	if query != nil {
		body["query"] = query
	}

	bodyReader, err := toJSONReader(body)
	if err != nil {
		return nil, err
	}

	scrollDuration := time.Minute
	resp, err := c.client.Search(ctx, &opensearchapi.SearchReq{
		Indices: []string{index},
		Body:    bodyReader,
		Params: opensearchapi.SearchParams{
			Scroll: scrollDuration,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("scroll 查询 %s 失败: %w", index, err)
	}

	var ids []string
	for _, hit := range resp.Hits.Hits {
		ids = append(ids, hit.ID)
	}

	scrollID := resp.ScrollID
	for scrollID != nil && *scrollID != "" && len(resp.Hits.Hits) > 0 {
		scrollResp, sErr := c.client.Scroll.Get(ctx, opensearchapi.ScrollGetReq{
			ScrollID: *scrollID,
			Params: opensearchapi.ScrollGetParams{
				Scroll: scrollDuration,
			},
		})
		if sErr != nil {
			break
		}
		for _, hit := range scrollResp.Hits.Hits {
			ids = append(ids, hit.ID)
		}
		if scrollResp.ScrollID != nil {
			scrollID = scrollResp.ScrollID
		}
		if len(scrollResp.Hits.Hits) == 0 {
			break
		}
	}

	// 清理 scroll context
	if scrollID != nil && *scrollID != "" {
		_, _ = c.client.Scroll.Delete(ctx, opensearchapi.ScrollDeleteReq{
			ScrollIDs: []string{*scrollID},
		})
	}

	return ids, nil
}

// WalkDocumentIDs 逐页遍历命中 ID，任何滚动页失败都会返回错误并尝试释放服务端游标。
func (c *Client) WalkDocumentIDs(ctx context.Context, index string, query map[string]any, pageSize int, visit func([]string) error) error {
	if pageSize < 1 || pageSize > 10000 || visit == nil {
		return fmt.Errorf("滚动查询 %s 参数无效", index)
	}
	body := map[string]any{"size": pageSize, "_source": false, "sort": []string{"_doc"}}
	if query != nil {
		body["query"] = query
	}
	bodyReader, err := toJSONReader(body)
	if err != nil {
		return err
	}
	scrollDuration := time.Minute
	resp, err := c.client.Search(ctx, &opensearchapi.SearchReq{Indices: []string{index}, Body: bodyReader,
		Params: opensearchapi.SearchParams{Scroll: scrollDuration}})
	if err != nil {
		return fmt.Errorf("scroll 查询 %s 失败: %w", index, err)
	}
	var scrollID string
	if resp.ScrollID != nil {
		scrollID = *resp.ScrollID
	}
	defer func() {
		if scrollID == "" {
			return
		}
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, _ = c.client.Scroll.Delete(cleanupCtx, opensearchapi.ScrollDeleteReq{ScrollIDs: []string{scrollID}})
	}()
	for {
		ids := make([]string, 0, len(resp.Hits.Hits))
		for _, hit := range resp.Hits.Hits {
			if hit.ID != "" {
				ids = append(ids, hit.ID)
			}
		}
		if len(ids) > 0 {
			if err := visit(ids); err != nil {
				return err
			}
		}
		if len(resp.Hits.Hits) == 0 {
			return nil
		}
		if scrollID == "" {
			return fmt.Errorf("scroll 查询 %s 缺少游标", index)
		}
		next, nextErr := c.client.Scroll.Get(ctx, opensearchapi.ScrollGetReq{ScrollID: scrollID,
			Params: opensearchapi.ScrollGetParams{Scroll: scrollDuration}})
		if nextErr != nil {
			return fmt.Errorf("scroll 查询 %s 下一页失败: %w", index, nextErr)
		}
		resp.Hits.Hits = next.Hits.Hits
		if next.ScrollID != nil {
			scrollID = *next.ScrollID
		}
	}
}

const aggSyncWatermark = "sync_watermark"

// MaxAggDateField 在指定 query 下对 date 字段做 max 聚合，得到索引侧「当前作用域」最新文档时间，用作增量同步水位
// query 为 nil 时表示 match_all；无文档或聚合值为空时 ok 为 false
func (c *Client) MaxAggDateField(ctx context.Context, index, field string, query map[string]any) (t time.Time, ok bool, err error) {
	body := map[string]any{
		"size": 0,
		"query": map[string]any{
			"match_all": map[string]any{},
		},
		"aggs": map[string]any{
			aggSyncWatermark: map[string]any{
				"max": map[string]any{"field": field},
			},
		},
	}
	if query != nil {
		body["query"] = query
	}
	bodyReader, mErr := toJSONReader(body)
	if mErr != nil {
		return time.Time{}, false, mErr
	}
	resp, sErr := c.client.Search(ctx, &opensearchapi.SearchReq{
		Indices: []string{index},
		Body:    bodyReader,
	})
	if sErr != nil {
		if resp != nil {
			ins := resp.Inspect()
			if ins.Response != nil && ins.Response.StatusCode == http.StatusNotFound {
				return time.Time{}, false, nil
			}
		}
		return time.Time{}, false, fmt.Errorf("聚合索引 %s 字段 %s 最大值失败: %w", index, field, sErr)
	}
	return parseSyncWatermarkAggregations(resp.Aggregations)
}

func parseSyncWatermarkAggregations(raw json.RawMessage) (t time.Time, ok bool, err error) {
	if len(raw) == 0 {
		return time.Time{}, false, nil
	}
	var aggs map[string]json.RawMessage
	if err = json.Unmarshal(raw, &aggs); err != nil {
		return time.Time{}, false, fmt.Errorf("解析 aggregations JSON 失败: %w", err)
	}
	b, exists := aggs[aggSyncWatermark]
	if !exists {
		return time.Time{}, false, nil
	}
	var node map[string]any
	if err = json.Unmarshal(b, &node); err != nil {
		return time.Time{}, false, fmt.Errorf("解析 sync_watermark 聚合失败: %w", err)
	}
	if vas, isStr := node["value_as_string"].(string); isStr && vas != "" {
		for _, layout := range []string{time.RFC3339Nano, time.RFC3339} {
			if parsed, pErr := time.Parse(layout, vas); pErr == nil {
				return parsed, true, nil
			}
		}
		// 部分集群返回非标准日期串，继续尝试数值 value
	}
	switch v := node["value"].(type) {
	case float64:
		if math.IsNaN(v) {
			return time.Time{}, false, nil
		}
		return time.UnixMilli(int64(v)), true, nil
	case nil:
		return time.Time{}, false, nil
	default:
		return time.Time{}, false, nil
	}
}

// DeleteIndex 删除索引
func (c *Client) DeleteIndex(ctx context.Context, name string) error {
	_, err := c.client.Indices.Delete(ctx, opensearchapi.IndicesDeleteReq{
		Indices: []string{name},
	})
	if err != nil {
		return fmt.Errorf("删除索引 %s 失败: %w", name, err)
	}
	return nil
}

// RefreshIndex 等待索引当前写入可搜索，仅供需要显式刷新语义的运维或验收路径使用。
func (c *Client) RefreshIndex(ctx context.Context, name string) error {
	_, err := c.client.Indices.Refresh(ctx, &opensearchapi.IndicesRefreshReq{Indices: []string{name}})
	if err != nil {
		return fmt.Errorf("刷新索引 %s 失败: %w", name, err)
	}
	return nil
}

// GetIndexFieldTypes 获取索引当前 mapping 中各字段的 type
// 返回 map[fieldName]type，如 {"name": "keyword", "age": "integer", "count.view": "long"}
func (c *Client) GetIndexFieldTypes(ctx context.Context, name string) (map[string]string, error) {
	resp, err := c.client.Indices.Mapping.Get(ctx, &opensearchapi.MappingGetReq{
		Indices: []string{name},
	})
	if err != nil {
		return nil, fmt.Errorf("获取索引 %s mapping 失败: %w", name, err)
	}
	idx, ok := resp.Indices[name]
	if !ok {
		return nil, fmt.Errorf("mapping 响应中未找到索引 %s", name)
	}
	var mappings struct {
		Properties map[string]any `json:"properties"`
	}
	if err := json.Unmarshal(idx.Mappings, &mappings); err != nil {
		return nil, fmt.Errorf("解析索引 %s mapping 失败: %w", name, err)
	}
	return collectMappingFieldTypes(mappings.Properties, ""), nil
}

// IndexUsesStrictMapping 检查索引根 mapping 是否明确禁用动态字段。
func (c *Client) IndexUsesStrictMapping(ctx context.Context, name string) (bool, error) {
	resp, err := c.client.Indices.Mapping.Get(ctx, &opensearchapi.MappingGetReq{Indices: []string{name}})
	if err != nil {
		return false, fmt.Errorf("获取索引 %s mapping 失败: %w", name, err)
	}
	idx, ok := resp.Indices[name]
	if !ok {
		return false, fmt.Errorf("mapping 响应中未找到索引 %s", name)
	}
	var mappings struct {
		Dynamic any `json:"dynamic"`
	}
	if err := json.Unmarshal(idx.Mappings, &mappings); err != nil {
		return false, fmt.Errorf("解析索引 %s mapping 失败: %w", name, err)
	}
	return mappings.Dynamic == "strict", nil
}

// ExtractDesiredFieldTypes 从索引定义中提取期望的字段类型映射
func ExtractDesiredFieldTypes(definition map[string]any) map[string]string {
	result := map[string]string{}
	mappings, _ := definition["mappings"].(map[string]any)
	if mappings == nil {
		return result
	}
	properties, _ := mappings["properties"].(map[string]any)
	if properties == nil {
		return result
	}
	result = collectMappingFieldTypes(properties, "")
	return result
}

func collectMappingFieldTypes(properties map[string]any, prefix string) map[string]string {
	result := map[string]string{}
	for name, prop := range properties {
		propMap, _ := prop.(map[string]any)
		if propMap == nil {
			continue
		}
		fieldName := name
		if prefix != "" {
			fieldName = prefix + "." + name
		}
		if t, ok := propMap["type"].(string); ok {
			result[fieldName] = t
		}
		nestedProperties, _ := propMap["properties"].(map[string]any)
		if len(nestedProperties) > 0 {
			for nestedField, nestedType := range collectMappingFieldTypes(nestedProperties, fieldName) {
				result[nestedField] = nestedType
			}
		}
	}
	return result
}

// PutMappingProperties 向已有索引在线添加新字段（不能修改已有字段的类型）
func (c *Client) PutMappingProperties(ctx context.Context, name string, properties map[string]any) error {
	body := map[string]any{"properties": properties}
	bodyReader, err := toJSONReader(body)
	if err != nil {
		return err
	}
	_, err = c.client.Indices.Mapping.Put(ctx, opensearchapi.MappingPutReq{
		Indices: []string{name},
		Body:    bodyReader,
	})
	if err != nil {
		return fmt.Errorf("更新索引 %s mapping 失败: %w", name, err)
	}
	return nil
}

// MappingDrift 描述单个字段的 mapping 漂移
type MappingDrift struct {
	Field    string
	Expected string
	Actual   string
}

// MappingDiffResult 对比期望 mapping 与实际 mapping 的差异
type MappingDiffResult struct {
	// 新增字段：期望有但实际索引中不存在，可通过 PutMapping 在线添加
	NewFields map[string]string
	// 类型变更：字段已存在但类型不同，只能通过删除重建索引解决
	TypeChanges []MappingDrift
}

// IndexEnsureResult 索引检查结果
type IndexEnsureResult struct {
	Created     []string                  // 新创建的索引
	FieldsAdded map[string][]string       // 在线添加了新字段的索引（index → fieldNames）
	Migrated    []string                  // 类型变更后自动重建的索引
	Drifts      map[string][]MappingDrift // 类型变更但未自动重建的索引
}

// DiffMapping 比较期望 mapping 与实际 mapping，区分新增字段和类型变更
func DiffMapping(desired, actual map[string]string) MappingDiffResult {
	result := MappingDiffResult{NewFields: map[string]string{}}
	for field, expectedType := range desired {
		actualType, exists := actual[field]
		if !exists {
			result.NewFields[field] = expectedType
			continue
		}
		if actualType != expectedType {
			result.TypeChanges = append(result.TypeChanges, MappingDrift{
				Field:    field,
				Expected: expectedType,
				Actual:   actualType,
			})
		}
	}
	return result
}

func toJSONReader(v any) (io.Reader, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("JSON 序列化失败: %w", err)
	}
	return strings.NewReader(string(data)), nil
}
