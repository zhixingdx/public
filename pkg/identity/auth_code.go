package identity

import (
	"context"
	"crypto/subtle"
	"errors"
	"strings"

	km "github.com/go-kratos/kratos/v3/metadata"
	"github.com/zhixingdx/public/api/apicommon/v1"
	"google.golang.org/grpc/metadata"
)

// AuthCodeTokens 仅用于没有人员身份的认证码 RPC。每个调用服务使用独立凭据。
type AuthCodeTokens map[string]string

func (tokens AuthCodeTokens) Validate() error {
	seen := map[string]bool{}
	for service, token := range tokens {
		if (service != "base-user" && service != "base-iam" && service != "base-safe") || len(token) < 32 || strings.IndexFunc(token, func(r rune) bool { return r <= 32 || r >= 127 }) >= 0 || seen[token] {
			return errors.New("authCodeTokens 需要按认证服务设置独立的至少 32 字节可打印凭据")
		}
		seen[token] = true
	}
	return nil
}

// CodeActor 优先使用可信人员上下文。仅签发和核验登录前验证码时可使用调用服务凭据。
// safe 按业务来源验证认证服务；notify 只接受 safe 代为投递，并保存原认证服务身份。
func (tokens AuthCodeTokens) CodeActor(ctx context.Context, call *apicommon.BizCall, caller string) (Actor, error) {
	md, _ := km.FromServerContext(ctx)
	if len(md.Values("x-md-global-user-id"))+len(md.Values("x-md-global-staff-id")) > 0 {
		return FromContext(ctx)
	}
	if call.GetService() != "base-user" && call.GetService() != "base-iam" {
		return Actor{}, ErrPersonRequired
	}
	switch call.GetAction() {
	case "login", "register", "reset-password":
	default:
		return Actor{}, ErrPersonRequired
	}
	expected := tokens[caller]
	bearer := metadata.ValueFromIncomingContext(ctx, "authorization")
	if expected == "" || len(bearer) != 1 || subtle.ConstantTimeCompare([]byte(bearer[0]), []byte("Bearer "+expected)) != 1 {
		return Actor{}, ErrPersonRequired
	}
	return Actor{Type: "service", ID: call.Service}, nil
}
