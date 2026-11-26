package tag

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"
)

// Converter 标签转换器 - 用于在 tagIDs（字符串数组）和 groupMap（分组映射）之间进行转换
// tagID 格式: biz::groupID::tagID
// biz 标识行业（必填），groupID 标识资源
type Converter struct {
	biz      string              // 行业标识（必填，1-16 个字符）
	tagIDs   []string            // 标签 ID 列表
	groupMap map[string][]string // 分组映射，key 为 groupID，value 为该组下的 tagID 列表
}

const (
	Separator = "::" // 标签 ID 分隔符
	maxBizLen = 16   // biz 最大长度
)

var (
	// ErrBizEmpty biz 为空错误
	ErrBizEmpty = fmt.Errorf("biz 不能为空")
	// ErrBizTooLong biz 长度超过限制错误
	ErrBizTooLong = fmt.Errorf("biz 长度不能超过 %d 个字符", maxBizLen)
	// ErrBizContainsSeparator biz 包含分隔符错误
	ErrBizContainsSeparator = fmt.Errorf("biz 不能包含分隔符 %s", Separator)
	// ErrGroupIDEmpty groupID 为空错误
	ErrGroupIDEmpty = fmt.Errorf("groupID 不能为空")
	// ErrTagIDEmpty tagID 为空错误
	ErrTagIDEmpty = fmt.Errorf("tagID 不能为空")
	// ErrTagIDFormat 表示外部 tagID 不符合 biz::groupID::tagID。
	ErrTagIDFormat = fmt.Errorf("tagID 格式必须为 biz::groupID::tagID")
)

// ParsedTagID 是严格解析后的行业、分组和权威标签 ID。
type ParsedTagID struct {
	BizCode string
	GroupID string
	TagID   string
}

// ParseTagID 严格解析单个 tagID，任何空段、多段或少段都会返回错误。
func ParseTagID(value string) (ParsedTagID, error) {
	parts := strings.Split(strings.TrimSpace(value), Separator)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return ParsedTagID{}, fmt.Errorf("%w: %q", ErrTagIDFormat, value)
	}
	if len(parts[0]) > maxBizLen {
		return ParsedTagID{}, ErrBizTooLong
	}
	return ParsedTagID{BizCode: parts[0], GroupID: parts[1], TagID: parts[2]}, nil
}

// ParseTagIDs 严格解析并拒绝重复 tagID，供持久化前的 Schema 校验使用。
func ParseTagIDs(values []string) ([]ParsedTagID, error) {
	result := make([]ParsedTagID, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		normalized := strings.TrimSpace(value)
		if _, exists := seen[normalized]; exists {
			return nil, fmt.Errorf("tagID 重复: %s", normalized)
		}
		parsed, err := ParseTagID(normalized)
		if err != nil {
			return nil, err
		}
		seen[normalized] = struct{}{}
		result = append(result, parsed)
	}
	return result, nil
}

// NewConverter 创建新的标签转换器
// 参数:
//   - biz: 行业标识（必填，1-16 个字符，不能包含分隔符 ::）
//
// 返回:
//   - *Converter: 转换器实例
//   - error: 如果 biz 为空、长度超过限制或包含分隔符，返回错误
//
// 示例:
//   - converter, err := NewConverter("retail") // 使用 biz 标识行业
func NewConverter(biz string) (*Converter, error) {
	if biz == "" {
		return nil, ErrBizEmpty
	}
	if len(biz) > maxBizLen {
		return nil, ErrBizTooLong
	}
	if strings.Contains(biz, Separator) {
		return nil, ErrBizContainsSeparator
	}
	return &Converter{
		biz:      biz,
		groupMap: make(map[string][]string),
	}, nil
}

// toTagIDs 将 groupMap 转换为 tagIDs
// 将分组映射转换为格式化的标签 ID 列表
// 格式: biz::groupID::tagID（biz 必填）
// 注意: 会自动过滤掉空的 groupID 和 tagID
func (c *Converter) toTagIDs() {
	if c.groupMap == nil {
		c.tagIDs = make([]string, 0)
		return
	}
	entries := lo.Entries(c.groupMap)
	slices.SortFunc(entries, func(left, right lo.Entry[string, []string]) int {
		return strings.Compare(left.Key, right.Key)
	})
	// 稳定排序让相同业务输入得到相同 tagIDs、投影 hash 和幂等键。
	c.tagIDs = lo.FlatMap(entries, func(entry lo.Entry[string, []string], _ int) []string {
		groupID, groupTagIDs := entry.Key, slices.Clone(entry.Value)
		// 跳过空的 groupID
		if groupID == "" {
			return nil
		}
		slices.Sort(groupTagIDs)
		// 过滤掉空的 tagID 并生成格式化的 tagID
		return lo.FilterMap(groupTagIDs, func(tagID string, _ int) (string, bool) {
			if tagID == "" {
				return "", false
			}
			return strings.Join([]string{c.biz, groupID, tagID}, Separator), true
		})
	})
}

// toGroupMap 将 tagIDs 转换为 groupMap
// 解析标签 ID 列表，提取 groupID 和 tagID，并按照 groupID 分组
// 格式: biz::groupID::tagID（biz 必填）
// 注意: 会自动过滤掉空的 groupID 和 tagID
func (c *Converter) toGroupMap() {
	groupMap := make(map[string][]string)
	for _, tagID := range c.tagIDs {
		parts := strings.Split(tagID, Separator)

		// biz 必填，期望格式为 biz::groupID::tagID，长度为 3
		expectedLen := 3
		if len(parts) != expectedLen {
			continue
		}

		// 验证 biz 是否匹配
		if parts[0] != c.biz {
			continue
		}

		// 提取 groupID 和 tagID，过滤掉空值
		groupID := parts[1]
		tagIDValue := parts[2]
		if groupID == "" || tagIDValue == "" {
			continue
		}
		groupMap[groupID] = append(groupMap[groupID], tagIDValue)
	}
	c.groupMap = groupMap
}

// SetTagIDs 设置标签 ID 列表，并自动转换为 groupMap
// 参数:
//   - tagIDs: 标签 ID 列表
func (c *Converter) SetTagIDs(tagIDs []string) {
	c.tagIDs = tagIDs
	c.toGroupMap()
}

// GetTagIDs 获取标签 ID 列表
// 返回:
//   - []string: 标签 ID 列表
func (c *Converter) GetTagIDs() []string {
	return c.tagIDs
}

// SetGroupMap 设置分组映射，并自动转换为 tagIDs
// 参数:
//   - groupMap: 分组映射，key 为 groupID，value 为该组下的 tagID 列表
//
// 返回:
//   - error: 如果 groupID 或 tagID 为空，返回错误
func (c *Converter) SetGroupMap(groupMap map[string][]string) error {
	// 验证 groupID 和 tagID 不能为空
	for groupID, tagIDs := range groupMap {
		if groupID == "" {
			return ErrGroupIDEmpty
		}
		if slices.Contains(tagIDs, "") {
			return ErrTagIDEmpty
		}
	}
	c.groupMap = groupMap
	c.toTagIDs()
	return nil
}

// GetGroupMap 获取分组映射
// 返回:
//   - map[string][]string: 分组映射
func (c *Converter) GetGroupMap() map[string][]string {
	return c.groupMap
}

// SetGroup 设置指定组的标签 ID 列表
// 参数:
//   - groupID: 组 ID（必填，不能为空）
//   - tagIDs: 标签 ID 列表（tagID 不能为空），如果为空则删除该组
//
// 返回:
//   - error: 如果 groupID 或 tagID 为空，返回错误
func (c *Converter) SetGroup(groupID string, tagIDs []string) error {
	if groupID == "" {
		return ErrGroupIDEmpty
	}
	if c.groupMap == nil {
		c.groupMap = make(map[string][]string)
	}
	if len(tagIDs) == 0 {
		// 如果 tagIDs 为空，删除该组
		delete(c.groupMap, groupID)
	} else {
		// 验证 tagID 不能为空
		if slices.Contains(tagIDs, "") {
			return ErrTagIDEmpty
		}
		c.groupMap[groupID] = tagIDs
	}
	c.toTagIDs()
	return nil
}

// DeleteGroup 删除指定组
// 参数:
//   - groupID: 组 ID
func (c *Converter) DeleteGroup(groupID string) {
	if c.groupMap == nil {
		return
	}
	delete(c.groupMap, groupID)
	c.toTagIDs()
}

// GetGroupTagIDs 获取指定组的标签 ID 列表
// 参数:
//   - groupID: 组 ID
//
// 返回:
//   - []string: 标签 ID 列表，如果组不存在则返回 nil
func (c *Converter) GetGroupTagIDs(groupID string) []string {
	if c.groupMap == nil {
		return nil
	}
	return c.groupMap[groupID]
}

// AddGroupTagID 向指定组添加标签 ID（支持添加多个）
// 如果标签 ID 已存在，则不会重复添加（自动去重）
// 参数:
//   - groupID: 组 ID（必填，不能为空）
//   - tagIDs: 要添加的标签 ID 列表（tagID 不能为空，可变参数）
//
// 返回:
//   - error: 如果 groupID 或 tagID 为空，返回错误
func (c *Converter) AddGroupTagID(groupID string, tagIDs ...string) error {
	if groupID == "" {
		return ErrGroupIDEmpty
	}
	if c.groupMap == nil {
		c.groupMap = make(map[string][]string)
	}
	if len(tagIDs) == 0 {
		return nil
	}
	// 验证 tagID 不能为空
	if slices.Contains(tagIDs, "") {
		return ErrTagIDEmpty
	}
	// 合并现有 tagIDs 和新添加的 tagIDs，然后去重
	existing := c.groupMap[groupID]
	merged := append(existing, tagIDs...)
	c.groupMap[groupID] = lo.Uniq(merged)
	c.toTagIDs()
	return nil
}

// RemoveGroupTagID 从指定组移除标签 ID（支持移除多个）
// 如果移除后组为空，则自动删除该组
// 参数:
//   - groupID: 组 ID
//   - tagIDs: 要移除的标签 ID 列表（可变参数）
func (c *Converter) RemoveGroupTagID(groupID string, tagIDs ...string) {
	if c.groupMap == nil {
		return
	}
	if len(tagIDs) == 0 {
		return
	}
	// 将要移除的 tagIDs 转换为 map 以提高查找效率
	removeSet := lo.SliceToMap(tagIDs, func(tagID string) (string, bool) {
		return tagID, true
	})
	// 过滤掉要移除的 tagIDs
	filtered := lo.Filter(c.groupMap[groupID], func(item string, _ int) bool {
		_, shouldRemove := removeSet[item]
		return !shouldRemove
	})
	// 如果过滤后组为空，删除该组；否则更新组
	if len(filtered) == 0 {
		delete(c.groupMap, groupID)
	} else {
		c.groupMap[groupID] = filtered
	}
	c.toTagIDs()
}
