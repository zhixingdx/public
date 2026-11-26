// Package initdata 提供初始化数据的通用功能
/* cSpell:disable */
package initdata

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// DefaultRemoteLoadTimeout 远程 JSON 下载超时（含读 body）
const DefaultRemoteLoadTimeout = 60 * time.Second

var remoteLoadHTTPClient = &http.Client{Timeout: DefaultRemoteLoadTimeout}

// LoadData 加载初始化数据（从 URL 或本地数据）
// 参数:
//   - ctx: 上下文信息
//   - cfg: 初始化数据配置
//   - entityName: 实体名称（用于错误提示）
//
// 返回:
//   - []*T: 解析后的数据数组
//   - error: 错误信息
func LoadData[T any](ctx context.Context, cfg *InitData, entityName string) ([]*T, error) {
	var (
		payload []byte
		err     error
		source  = "本地"
	)

	if cfg.URL != "" {
		source = "远程"
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.URL, nil)
		if err != nil {
			return nil, fmt.Errorf("初始化%s失败: 构建请求异常: %w", entityName, err)
		}
		resp, err := remoteLoadHTTPClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("初始化%s失败: 请求远程数据异常: %w", entityName, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("初始化%s失败: 远程返回状态码 %d", entityName, resp.StatusCode)
		}
		payload, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("初始化%s失败: 读取远程响应异常: %w", entityName, err)
		}
	} else if cfg.File != "" {
		source = "本地文件"
		payload, err = os.ReadFile(cfg.File)
		if err != nil {
			return nil, fmt.Errorf("初始化%s失败: 读取本地文件 %s 异常: %w", entityName, cfg.File, err)
		}
	} else {
		if len(cfg.Data) == 0 {
			return nil, nil
		}
		payload, err = json.Marshal(cfg.Data)
		if err != nil {
			return nil, fmt.Errorf("初始化%s失败: 无法序列化本地数据: %w", entityName, err)
		}
	}

	if len(payload) == 0 {
		return nil, nil
	}

	var data []*T
	if err := json.Unmarshal(payload, &data); err != nil {
		var single T
		if singleErr := json.Unmarshal(payload, &single); singleErr != nil {
			return nil, fmt.Errorf("初始化%s失败: 解析%s数据异常: %w", entityName, source, err)
		}
		data = []*T{&single}
	}
	return data, nil
}
