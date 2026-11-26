// Package initdata 提供初始化数据的通用功能
package initdata

import (
	"context"
	"fmt"

	"github.com/samber/lo"
)

// DefaultSyncBatchSize 同步写库默认分批大小
const DefaultSyncBatchSize = 500

// SyncStats 实体同步统计
type SyncStats struct {
	Added   int32 `json:"added"`
	Deleted int32 `json:"deleted"`
	Updated int32 `json:"updated"`
	Skipped int32 `json:"skipped"`
}

// EntitySyncOptions 实体同步选项（先 diff 收集，再分批批量写入）
type EntitySyncOptions[T any] struct {
	Config     *InitData
	EntityName string
	ExtractKey func(*T) string
	Validate   func(*T) bool
	// IsEqual 比较业务字段是否一致；返回 true 表示跳过
	IsEqual func(local, remote *T) bool
	// QueryLocal 返回本地全量数据，key 为业务唯一键
	QueryLocal func(ctx context.Context) (map[string]*T, error)
	// BatchInsert / BatchUpdate / BatchDelete 分批回调，单次 len 不超过 BatchSize
	BatchInsert func(ctx context.Context, items []*T) error
	BatchUpdate func(ctx context.Context, locals, remotes []*T) error
	BatchDelete func(ctx context.Context, items []*T) error
	// BatchSize 分批大小，≤0 时使用 DefaultSyncBatchSize
	BatchSize int
}

// SyncEntities 按云端数据与本地数据 diff 同步：先收集变更，再分批批量提交
func SyncEntities[T any](ctx context.Context, opts EntitySyncOptions[T]) (*SyncStats, error) {
	stats := &SyncStats{}
	if opts.Config == nil {
		return stats, nil
	}

	batchSize := opts.BatchSize
	if batchSize <= 0 {
		batchSize = DefaultSyncBatchSize
	}

	remoteItems, err := LoadData[T](ctx, opts.Config, opts.EntityName)
	if err != nil {
		return nil, err
	}
	if opts.Validate != nil {
		remoteItems = lo.Filter(remoteItems, func(item *T, _ int) bool {
			return opts.Validate(item)
		})
	}
	remoteItems = lo.Filter(remoteItems, func(item *T, _ int) bool {
		return opts.ExtractKey(item) != ""
	})
	remoteItems = lo.UniqBy(remoteItems, func(item *T) string {
		return opts.ExtractKey(item)
	})

	remoteMap := lo.SliceToMap(remoteItems, func(item *T) (string, *T) {
		return opts.ExtractKey(item), item
	})

	localMap, err := opts.QueryLocal(ctx)
	if err != nil {
		return nil, fmt.Errorf("同步%s失败: 查询本地数据异常: %w", opts.EntityName, err)
	}

	toDelete := make([]*T, 0)
	for key, localItem := range localMap {
		if _, ok := remoteMap[key]; !ok {
			toDelete = append(toDelete, localItem)
		}
	}

	toInsert := make([]*T, 0)
	toUpdateLocal := make([]*T, 0)
	toUpdateRemote := make([]*T, 0)
	for key, remoteItem := range remoteMap {
		localItem, exists := localMap[key]
		if !exists {
			toInsert = append(toInsert, remoteItem)
			continue
		}
		if opts.IsEqual != nil && opts.IsEqual(localItem, remoteItem) {
			stats.Skipped++
			continue
		}
		toUpdateLocal = append(toUpdateLocal, localItem)
		toUpdateRemote = append(toUpdateRemote, remoteItem)
	}

	if len(toDelete) > 0 && opts.BatchDelete != nil {
		if err := forEachBatch(toDelete, batchSize, func(batch []*T) error {
			return opts.BatchDelete(ctx, batch)
		}); err != nil {
			return nil, fmt.Errorf("同步%s失败: 批量删除异常: %w", opts.EntityName, err)
		}
		stats.Deleted = int32(len(toDelete))
	}

	if len(toInsert) > 0 && opts.BatchInsert != nil {
		if err := forEachBatch(toInsert, batchSize, func(batch []*T) error {
			return opts.BatchInsert(ctx, batch)
		}); err != nil {
			return nil, fmt.Errorf("同步%s失败: 批量新增异常: %w", opts.EntityName, err)
		}
		stats.Added = int32(len(toInsert))
	}

	if len(toUpdateLocal) > 0 && opts.BatchUpdate != nil {
		if err := forEachBatchPairs(toUpdateLocal, toUpdateRemote, batchSize, func(locals, remotes []*T) error {
			return opts.BatchUpdate(ctx, locals, remotes)
		}); err != nil {
			return nil, fmt.Errorf("同步%s失败: 批量更新异常: %w", opts.EntityName, err)
		}
		stats.Updated = int32(len(toUpdateLocal))
	}

	return stats, nil
}

func forEachBatch[T any](items []*T, size int, fn func(batch []*T) error) error {
	for i := 0; i < len(items); i += size {
		end := i + size
		if end > len(items) {
			end = len(items)
		}
		if err := fn(items[i:end]); err != nil {
			return err
		}
	}
	return nil
}

func forEachBatchPairs[T any](locals, remotes []*T, size int, fn func(locals, remotes []*T) error) error {
	if len(locals) != len(remotes) {
		return fmt.Errorf("批量更新本地与远端条目数不一致")
	}
	for i := 0; i < len(locals); i += size {
		end := i + size
		if end > len(locals) {
			end = len(locals)
		}
		if err := fn(locals[i:end], remotes[i:end]); err != nil {
			return err
		}
	}
	return nil
}
