// Package initdata 提供初始化数据的通用功能
/* cSpell:disable */
package initdata

import (
	"context"
	"fmt"

	"github.com/samber/lo"
)

// EntityInitOptions 初始化实体的选项
type EntityInitOptions[T any] struct {
	Config        *InitData                                    // 数据配置（可选，如果提供了 DefaultData 则不需要）
	EntityName    string                                       // 实体名称（用于日志和错误提示）
	ExtractKey    func(*T) string                              // 从数据中提取唯一键
	Validate      func(*T) bool                                // 验证数据有效性（可选）
	QueryExisting func(ctx context.Context, keys []string) ([]string, error) // 查询已存在的实体
	BatchInsert   func(ctx context.Context, data []*T) error    // 批量插入实体
	DefaultData   []*T                                         // 默认数据（可选，如果提供则直接使用，不通过 Config 加载）
}

// InitEntities 初始化实体（通用函数）
// 参数:
//   - ctx: 上下文信息
//   - opts: 初始化选项
// 返回:
//   - error: 错误信息
func InitEntities[T any](ctx context.Context, opts EntityInitOptions[T]) error {
	var data []*T
	var err error

	// 如果提供了默认数据，直接使用；否则从配置加载
	if len(opts.DefaultData) > 0 {
		data = opts.DefaultData
	} else {
		if opts.Config == nil {
			return nil
		}
		data, err = LoadData[T](ctx, opts.Config, opts.EntityName)
		if err != nil {
			return err
		}
		if len(data) == 0 {
			return nil
		}
	}

	// 验证和过滤数据
	if opts.Validate != nil {
		data = lo.Filter(data, func(item *T, _ int) bool {
			return opts.Validate(item)
		})
	}
	if len(data) == 0 {
		return nil
	}

	// 过滤掉唯一键为空的数据
	data = lo.Filter(data, func(item *T, _ int) bool {
		return opts.ExtractKey(item) != ""
	})
	if len(data) == 0 {
		return nil
	}

	// 按唯一键去重
	data = lo.UniqBy(data, func(item *T) string {
		return opts.ExtractKey(item)
	})
	if len(data) == 0 {
		return nil
	}

	// 提取所有唯一键
	keys := lo.Map(data, func(item *T, _ int) string {
		return opts.ExtractKey(item)
	})
	keys = lo.Uniq(keys)

	// 查询已存在的实体
	existingKeys, err := opts.QueryExisting(ctx, keys)
	if err != nil {
		return fmt.Errorf("初始化%s失败: 查询数据库异常: %w", opts.EntityName, err)
	}

	// 构建已存在的键集合
	if len(existingKeys) > 0 {
		existingSet := make(map[string]struct{}, len(existingKeys))
		for _, key := range existingKeys {
			if key != "" {
				existingSet[key] = struct{}{}
			}
		}
		// 过滤掉已存在的实体
		data = lo.Filter(data, func(item *T, _ int) bool {
			_, ok := existingSet[opts.ExtractKey(item)]
			return !ok
		})
	}
	if len(data) == 0 {
		return nil
	}

	// 批量插入
	if err := opts.BatchInsert(ctx, data); err != nil {
		return fmt.Errorf("初始化%s失败: 批量写入数据库异常: %w", opts.EntityName, err)
	}

	fmt.Printf("添加的%s数据: \n%v\n", opts.EntityName, data)
	return nil
}
