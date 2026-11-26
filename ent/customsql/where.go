// Package customsql where
package customsql

import (
	"encoding/json"
	"strconv"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/zhixingdx/oipivot/backend/ent/ann"
	"github.com/zhixingdx/oipivot/backend/pkg/tag"
	"github.com/samber/lo"
)

// GetStrSliceAndQuery 获取字符串切片的并查询
func GetStrSliceAndQuery(field string, values ...string) func(s *sql.Selector) {
	if len(values) == 0 {
		return func(_ *sql.Selector) {}
	}

	return func(s *sql.Selector) {
		if isPostgres(s) {
			//nolint:goconst
			str := `["` + strings.Join(values, `","`) + `"]`
			s.Where(sql.P(func(b *sql.Builder) {
				b.Ident(field).WriteString(" @> ").Arg(str)
			}))
		} else {
			str := `"` + strings.Join(values, `","`) + `"`
			s.Where(sql.P(func(b *sql.Builder) {
				b.WriteString("JSON_CONTAINS").Wrap(func(b *sql.Builder) {
					b.Ident(field).Comma()
					//nolint:goconst
					b.WriteString(`JSON_ARRAY(` + str + `)`)
				})
			}))
		}
	}
}

// GetStrSliceOrQuery 获取字符串切片的或查询
func GetStrSliceOrQuery(field string, values ...string) func(s *sql.Selector) {
	if len(values) == 0 {
		return func(_ *sql.Selector) {}
	}
	return func(s *sql.Selector) {
		if isPostgres(s) {
			pArr := lo.Map(values, func(val string, _ int) *sql.Predicate { return getPgJSONStrQuery(field, val) })
			s.Where(sql.Or(pArr...))
			return
		}
		pArr := lo.Map(values, func(val string, _ int) *sql.Predicate { return getMysqlJSONStrQuery(field, val) })
		s.Where(sql.Or(pArr...))

	}
}

// GetStrSlice 获取字符串切片的查询函数
// 支持普通数组查询和分层标签查询
func GetStrSlice(field string, queryType ann.JSONArrQueryType, values ...string) func(s *sql.Selector) {
	if len(values) == 0 {
		return func(_ *sql.Selector) {}
	}
	switch queryType {
	case ann.JSONArrQueryTypeAnd:
		return GetStrSliceAndQuery(field, values...)
	case ann.JSONArrQueryTypeOr, ann.JSONArrQueryTypeDefault:
		return GetStrSliceOrQuery(field, values...)
	default:
		// 处理分层标签查询（Tag 组查询）
		return getStrSliceGroupQuery(field, queryType, values...)
	}
}

// GetStrSliceNotOverlapQuery 获取字符串切片的不重叠查询。
// 用于“隐藏语言”等负向可见性判断：字段为空数组时天然满足不重叠。
func GetStrSliceNotOverlapQuery(field string, values ...string) func(s *sql.Selector) {
	if len(values) == 0 {
		return func(_ *sql.Selector) {}
	}
	return func(s *sql.Selector) {
		if isPostgres(s) {
			pArr := lo.Map(values, func(val string, _ int) *sql.Predicate { return getPgJSONStrQuery(field, val) })
			s.Where(sql.Not(sql.Or(pArr...)))
			return
		}
		s.Where(getMysqlJSONStrNotOverlapQuery(field, values...))
	}
}

// getStrSliceGroupQuery 处理分层标签查询
// 只有分割完 层级相同并且前缀相同的才算同一组
// 使用分隔符进行分割，如 edu::school::1 层级不定，需要分层处理
func getStrSliceGroupQuery(field string, queryType ann.JSONArrQueryType, values ...string) func(s *sql.Selector) {
	if len(values) == 0 {
		return func(_ *sql.Selector) {}
	}
	return func(s *sql.Selector) {
		// 分组 先按长度
		lenArr := lo.GroupBy(values, func(item string) int { return len(strings.Split(item, tag.Separator)) })
		groupPredicate := []*sql.Predicate{}
		for l, v := range lenArr {
			if len(v) == 1 {
				// 只有一层，直接构建 predicate 单独做为一组
				groupPredicate = append(groupPredicate, buildTagPredicate(s, field, queryType, v))
			} else {
				// 按照长度头一分组，层级相同并且前缀相同的才算同一组
				groupMap := lo.GroupBy(v, func(item string) string {
					// 去掉最后一层，按前面的分组
					return strings.Join(strings.Split(item, tag.Separator)[:l-1], tag.Separator)
				})
				for _, tagIDs := range groupMap {
					groupPredicate = append(groupPredicate, buildTagPredicate(s, field, queryType, tagIDs))
				}
			}
		}
		if len(groupPredicate) == 0 {
			return
		}
		// 必须判断为 1 的情况 否则在 or 中会与整体条件变成 or
		if len(groupPredicate) == 1 {
			s.Where(groupPredicate[0])
			return
		}
		if queryType == ann.JSONArrQueryTypeGroupAndItemOr {
			s.Where(sql.And(groupPredicate...))
		} else {
			s.Where(sql.Or(groupPredicate...))
		}
	}
}

func getMysqlJSONStrNotOverlapQuery(field string, values ...string) *sql.Predicate {
	pArr := lo.Map(values, func(value string, _ int) *sql.Predicate {
		return sql.P(func(b *sql.Builder) {
			b.WriteString("JSON_SEARCH").Wrap(func(b *sql.Builder) {
				b.Ident(field).Comma()
				b.Arg("one").Comma()
				b.Arg(value)
			}).WriteString(" IS NULL")
		})
	})
	return sql.And(pArr...)
}

func getMysqlJSONStrQuery(field string, value string) *sql.Predicate {
	return sql.P(func(b *sql.Builder) {
		b.WriteString("JSON_CONTAINS").Wrap(func(b *sql.Builder) {
			b.Ident(field).Comma()
			b.WriteString(`'"` + value + `"'`)
			b.WriteString(`,"$"`)
		})
	})
}

// buildTagPredicate 构建标签查询的 predicate
// 根据数据库类型和查询类型，将多个 tagID 组合成单个 predicate
func buildTagPredicate(s *sql.Selector, field string, queryType ann.JSONArrQueryType, tagIDs []string) *sql.Predicate {
	var tagPredicate []*sql.Predicate
	if isPostgres(s) {
		tagPredicate = lo.Map(tagIDs, func(tagID string, _ int) *sql.Predicate {
			return getPgJSONStrQuery(field, tagID)
		})
	} else {
		tagPredicate = lo.Map(tagIDs, func(tagID string, _ int) *sql.Predicate {
			return getMysqlJSONStrQuery(field, tagID)
		})
	}
	if queryType == ann.JSONArrQueryTypeGroupOrItemAnd || queryType == ann.JSONArrQueryTypeAnd {
		return sql.And(tagPredicate...)
	}
	return sql.Or(tagPredicate...)
}

func getPgJSONStrQuery(field string, value string) *sql.Predicate {
	return sql.P(func(b *sql.Builder) {
		jsonValue, _ := json.Marshal([]string{value})
		b.Ident(field).WriteString(" @> ").Arg(string(jsonValue))
	})
}

// GetIntSliceAndQuery 获取整数切片的并查询
func GetIntSliceAndQuery(field string, values ...int) func(s *sql.Selector) {
	if len(values) == 0 {
		return func(_ *sql.Selector) {}
	}
	return func(s *sql.Selector) {
		if isPostgres(s) {
			str := "[" + strings.Join(lo.Map(values, func(val int, _ int) string { return strconv.Itoa(val) }), ",") + "]"
			s.Where(sql.P(func(b *sql.Builder) {
				b.Ident(field).WriteString(" @> ").Arg(str)
			}))
		} else {
			str := strings.Join(lo.Map(values, func(val int, _ int) string { return strconv.Itoa(val) }), `,`)
			s.Where(sql.P(func(b *sql.Builder) {
				b.WriteString("JSON_CONTAINS").Wrap(func(b *sql.Builder) {
					b.Ident(field).Comma()
					b.WriteString(`JSON_ARRAY(` + str + `)`)
				})
			}))
		}
	}
}

// GetIntSliceOrQuery 获取整数切片的或查询
func GetIntSliceOrQuery(field string, values ...int) func(s *sql.Selector) {
	if len(values) == 0 {
		return func(_ *sql.Selector) {}
	}
	return func(s *sql.Selector) {
		if isPostgres(s) {
			pArr := lo.Map(values, func(val int, _ int) *sql.Predicate { return getPgJSONIntQuery(field, val) })
			s.Where(sql.Or(pArr...))
			return
		}
		pArr := lo.Map(values, func(val int, _ int) *sql.Predicate { return getMyJSONIntQuery(field, val) })
		s.Where(sql.Or(pArr...))
	}
}

func getMyJSONIntQuery(field string, value int) *sql.Predicate {
	return sql.P(func(b *sql.Builder) {
		b.WriteString("JSON_CONTAINS").Wrap(func(b *sql.Builder) {
			b.Ident(field).Comma()
			b.WriteString(`"` + strconv.Itoa(value) + `"`)
			b.WriteString(`,"$"`)
		})
	})
}

func getPgJSONIntQuery(field string, value int) *sql.Predicate {
	return sql.P(func(b *sql.Builder) {
		b.Ident(field).WriteString(" @> ").Arg(`[` + strconv.Itoa(value) + `]`)
	})
}

func isPostgres(s *sql.Selector) bool {
	return s.Dialect() == "postgres"
}
