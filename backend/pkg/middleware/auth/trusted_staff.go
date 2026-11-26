package auth

import (
	"context"
	"errors"
	"strings"

	directory "github.com/zhixingdx/oipivot-api/base/tenant/v1"

	kerrors "github.com/go-kratos/kratos/v3/errors"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/transport"
	"github.com/zhixingdx/oipivot-api/apicommon/v1"
	iam "github.com/zhixingdx/oipivot-api/base/iam/v1"
	"github.com/zhixingdx/oipivot-api/fault"
	"github.com/zhixingdx/oipivot/backend/pkg/lang"
	"github.com/zhixingdx/oipivot/backend/pkg/middleware/meta"
	"google.golang.org/grpc"
)

// TrustedStaffDependencies 复用 IAM 既有 RPC 契约，由服务装配提供客户端方法，避免业务执行器自行实现权限。
type TrustedStaffDependencies struct {
	GetStaffCache        func(context.Context, *iam.OneStaffRequest, ...grpc.CallOption) (*iam.StaffDetailReply, error)
	ResolveTenantRealm   func(context.Context, *directory.OneTenantRequest, ...grpc.CallOption) (*apicommon.Str, error)
	CheckStaffPermission func(context.Context, *iam.CheckStaffPermissionRequest, ...grpc.CallOption) (*apicommon.Bool, error)
}

// NewTrustedStaffMiddleware 为受信内网 RPC 复用员工状态和操作权限；它不认证网络来源，只能装配在已有可信边界内。
// 无身份请求始终拒绝，系统调用兼容由装配层选择入口；不得用客户端自报标记绕过此检查。
func NewTrustedStaffMiddleware(deps TrustedStaffDependencies, list *ListConf, l *lang.CommonLang) middleware.Middleware {
	if deps.GetStaffCache == nil || deps.ResolveTenantRealm == nil || deps.CheckStaffPermission == nil || l == nil {
		panic("可信员工鉴权依赖不完整")
	}
	permissions := []string{}
	if list != nil {
		// 白名单仅影响操作权限；此入口的职责是严格验证主体，不允许匿名执行。
		permissions = append(permissions, list.White...)
		permissions = append(permissions, list.Gray...)
		permissions = append(permissions, list.PermissionWhite...)
	}
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			t, ok := transport.FromServerContext(ctx)
			if !ok || strings.TrimSpace(t.Operation()) == "" {
				return nil, fault.ErrorBadRequest("%s", l.GetRequestErrorMsg(ctx))
			}
			tenantID, tenantErr := meta.GetValue(ctx, "tenantID")
			staffID, staffErr := meta.GetStaffID(ctx)
			if tenantErr != nil || staffErr != nil || strings.TrimSpace(tenantID) == "" || strings.TrimSpace(staffID) == "" {
				return nil, fault.ErrorUnauthorized("%s", l.GetPleaseLoginMsg(ctx))
			}
			realm, realmErr := meta.GetRealm(ctx)
			if realmErr != nil {
				return nil, fault.ErrorForbidden("%s", l.GetNoPermissionMsg(ctx))
			}
			authority, err := deps.ResolveTenantRealm(ctx, &directory.OneTenantRequest{Id: tenantID})
			if err != nil {
				return nil, trustedStaffDependencyError(ctx, err, l)
			}
			if authority == nil || strings.TrimSpace(authority.Value) == "" {
				return nil, fault.ErrorForbidden("%s", l.GetNoPermissionMsg(ctx))
			}
			expectedRealm, normalizeErr := meta.NormalizeRealm(authority.Value)
			if normalizeErr != nil || expectedRealm != realm {
				return nil, fault.ErrorForbidden("%s", l.GetNoPermissionMsg(ctx))
			}
			staff, err := deps.GetStaffCache(ctx, &iam.OneStaffRequest{Id: staffID, TenantID: tenantID})
			if err != nil {
				return nil, trustedStaffDependencyError(ctx, err, l)
			}
			if staff == nil || staff.Id != staffID || staff.TenantID != tenantID {
				return nil, fault.ErrorUnauthorized("%s", l.GetPleaseLoginMsg(ctx))
			}
			if !staff.Switch {
				return nil, fault.ErrorForbidden("%s", l.GetAccountDisabledMsg(ctx))
			}
			if staff.GetIsChangePassword() && t.Operation() != "/api.base.iam.BaseIamPortalPlus/ChangePassword" {
				denied := fault.ErrorForbidden("%s", l.GetPleaseChangePasswordMsg(ctx))
				denied.Metadata = map[string]string{"isChangePassword": "true"}
				return nil, denied
			}
			if !trustedStaffOperationMatches(t.Operation(), permissions) {
				allowed, err := deps.CheckStaffPermission(ctx, &iam.CheckStaffPermissionRequest{StaffID: staffID, Path: t.Operation()})
				if err != nil {
					return nil, trustedStaffDependencyError(ctx, err, l)
				}
				if allowed == nil || !allowed.Value {
					return nil, fault.ErrorForbidden("%s", l.GetNoPermissionMsg(ctx))
				}
			}
			return next(ctx, req)
		}
	}
}

// trustedStaffOperationMatches 沿用统一 IAM 列表精确匹配及末尾通配语义，不从请求参数推导操作项。
func trustedStaffOperationMatches(operation string, entries []string) bool {
	for _, entry := range entries {
		if operation == entry || len(entry) > 1 && strings.HasSuffix(entry, "*") && strings.HasPrefix(operation, strings.TrimSuffix(entry, "*")) {
			return true
		}
	}
	return false
}

// trustedStaffDependencyError 保留取消与授权分类，基础设施失败关闭且不向用户泄露 RPC、SQL 或网络原文。
func trustedStaffDependencyError(ctx context.Context, err error, l *lang.CommonLang) error {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	switch kerrors.Code(err) {
	case 401, 404:
		return fault.ErrorUnauthorized("%s", l.GetPleaseLoginMsg(ctx))
	case 403:
		return fault.ErrorForbidden("%s", l.GetNoPermissionMsg(ctx))
	default:
		return fault.ErrorServiceUnavailable("%s", l.GetRequestErrorMsg(ctx))
	}
}
