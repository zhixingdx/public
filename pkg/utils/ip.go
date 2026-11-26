package utils

import (
	"context"

	"github.com/zhixingdx/public/pkg/middleware/meta"
)

// 获取客户端 IP
func GetClientIP(ctx context.Context) string {
	header, is := meta.GetRequestHeader(ctx)
	if !is {
		return ""
	}
	ip := header.Get("X-Real-IP")
	if ip == "" {
		ip = header.Get("X-Forwarded-For")
	}
	return ip
}
