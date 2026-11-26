// Package auth 定义平台员工认证使用的白名单。
package auth

import (
	"github.com/go-kratos/kratos/v3/middleware"
	iam "github.com/zhixingdx/oipivot-api/base/iam/v1"
	"github.com/zhixingdx/oipivot-api/config"
	"github.com/zhixingdx/oipivot/backend/pkg/lang"
)

// ListConf 白名单和灰名单配置
// 用于配置不需要认证或特殊处理的路径
type ListConf struct {
	White           []string // 白名单路径 - 完全跳过认证
	Gray            []string // 灰名单路径 - 可选认证并跳过权限检查
	PermissionWhite []string // 权限白名单 - 跳过权限检查
}

const TypeIam = "iam"

// NewAuthMiddleware 为生成的服务模板装配平台员工身份。具体行业身份由业务项目接入。
func NewAuthMiddleware(c *config.Auth, list *ListConf, l *lang.CommonLang) middleware.Middleware {
	if c == nil || c.Type != TypeIam || c.Addr == "" {
		panic("平台 HTTP 须配置 IAM 认证地址")
	}
	if list == nil {
		list = new(ListConf)
	}
	client := iam.GetPortalPlusRPC(c.Addr)
	return AuthMiddle(client.GetStaffInfo, client.CheckStaffPermission, list.White, list.Gray, list.PermissionWhite, l)
}
