package auth

import (
	"context"

	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/transport"
	"github.com/samber/lo"
	"github.com/zhixingdx/public/api/apicommon/v1"
	iam "github.com/zhixingdx/public/api/base/iam/v1"
	"github.com/zhixingdx/public/api/fault"
	"github.com/zhixingdx/public/pkg/lang"
	"github.com/zhixingdx/public/pkg/middleware/meta"
	"google.golang.org/grpc"
)

func AuthMiddle(
	// 获取用户信息
	getInfo func(ctx context.Context, req *iam.StaffInfoRequest, opts ...grpc.CallOption) (*iam.AuthReply, error),
	// 检查权限
	checkStaffPermission func(ctx context.Context, in *iam.CheckStaffPermissionRequest, opts ...grpc.CallOption) (*apicommon.Bool, error),
	whiteList []string,
	grayList []string,
	permissionWhiteList []string,
	commonLang *lang.CommonLang,
) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			t, b := transport.FromServerContext(ctx)
			if !b {
				return nil, fault.ErrorBadRequest("%s", commonLang.GetRequestErrorMsg(ctx))
			}

			operation := t.Operation()
			if someList(operation, whiteList) {
				return handler(ctx, req)
			}

			token := meta.GetAuthorization(ctx)
			isGray := someList(operation, grayList)
			if token == "" {
				if isGray {
					return handler(ctx, req)
				}
				return nil, fault.ErrorUnauthorized("%s", commonLang.GetPleaseLoginMsg(ctx))
			}

			auth, err := getInfo(ctx, &iam.StaffInfoRequest{Token: token, UserAgent: meta.GetUserAgent(ctx)})
			if err != nil {
				if !isGray {
					return nil, err
				}
				return handler(ctx, req)
			}

			if auth == nil || auth.Staff == nil {
				if isGray {
					return handler(ctx, req)
				}
				return nil, fault.ErrorUnauthorized("%s", commonLang.GetPleaseLoginMsg(ctx))
			}

			staff := auth.Staff
			if !staff.Switch {
				return nil, fault.ErrorForbidden("%s", commonLang.GetAccountDisabledMsg(ctx))
			}
			if staff.GetIsChangePassword() && operation != "/api.base.iam.BaseIamPortalPlus/ChangePassword" {
				fErr := fault.ErrorForbidden("%s", commonLang.GetPleaseChangePasswordMsg(ctx))
				fErr.Metadata = map[string]string{"isChangePassword": "true"}
				return auth, fErr
			}

			ctx = setAuth(ctx, auth)
			ctx = meta.SetStaffID(ctx, staff.Id)
			if isGray {
				// Gray 只负责可选认证和身份注入；业务授权由 Handler 自行完成，不能因附带有效员工 Token 反而触发 RBAC 拒绝。
				return handler(ctx, req)
			}
			if !someList(operation, permissionWhiteList) {
				has, err := checkStaffPermission(ctx, &iam.CheckStaffPermissionRequest{Path: t.Operation(), StaffID: staff.Id})
				if err != nil {
					return nil, err
				}
				if !has.GetValue() {
					return nil, fault.ErrorForbidden("%s", commonLang.GetNoPermissionMsg(ctx))
				}
			}
			return handler(ctx, req)
		}
	}
}

// GetAuth 获取认证信息
func GetAuth(ctx context.Context, commonLang *lang.CommonLang) (*iam.AuthReply, error) {
	// 先尝试服务内取
	val := ctx.Value(staffAuthKey{})
	if val == nil {
		return nil, fault.ErrorUnauthorized("%s", commonLang.GetPleaseLoginMsg(ctx))
	}
	auth, ok := val.(*iam.AuthReply)
	if !ok {
		return nil, fault.ErrorUnauthorized("%s", commonLang.GetPleaseLoginMsg(ctx))
	}
	return auth, nil
}

func GetID(ctx context.Context, commonLang *lang.CommonLang) (string, error) {
	auth, err := GetAuth(ctx, commonLang)
	if err != nil {
		return "", err
	}
	if auth.Staff == nil || auth.Staff.Id == "" {
		return "", fault.ErrorUnauthorized("%s", commonLang.GetPleaseLoginMsg(ctx))
	}
	return auth.Staff.Id, nil
}

func setAuth(ctx context.Context, auth *iam.AuthReply) context.Context {
	return context.WithValue(ctx, staffAuthKey{}, auth)
}

type staffAuthKey struct{}

func someList(operation string, list []string) bool {
	return lo.SomeBy(list, func(white string) bool {
		if operation == white {
			return true
		}
		// 末尾泛匹配
		wLen := len(white)
		oLen := len(operation)
		if wLen > oLen {
			return false
		}
		return wLen > 1 && white[wLen-1] == '*' && operation[:wLen-1] == white[:wLen-1]
	})
}

// SetAuthForTest 注入已认证上下文供跨包测试使用，不应在生产代码调用。
func SetAuthForTest(ctx context.Context, auth *iam.AuthReply) context.Context {
	return setAuth(ctx, auth)
}
