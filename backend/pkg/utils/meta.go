package utils

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/go-kratos/kratos/v3/metadata"
	"github.com/samber/lo"
)

func GetMeta(ctx context.Context, sensitive ...string) string {
	metaStr := ""
	sensitiveKeys := lo.If(len(sensitive) > 0, sensitive).Else([]string{"auth", "token", "sign"})
	if md, ok := metadata.FromServerContext(ctx); ok {
		metaMap := map[string][]string{}
		md.Range(func(k string, v []string) bool {
			isSensitive := lo.SomeBy(sensitiveKeys, func(s string) bool {
				return strings.Contains(strings.ToLower(k), strings.ToLower(s))
			})
			if !isSensitive {
				metaMap[k] = v
			}
			return true
		})
		metaByte, mErr := json.Marshal(metaMap)
		if mErr == nil {
			metaStr = string(metaByte)
		}
	}
	return metaStr
}
