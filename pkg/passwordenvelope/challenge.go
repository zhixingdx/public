package passwordenvelope

import (
	"context"

	safe "github.com/zhixingdx/oipivot-api/base/safe/v1"
	user "github.com/zhixingdx/oipivot-api/base/user/v1"
	"github.com/zhixingdx/oipivot-api/fault"
	"github.com/zhixingdx/oipivot/backend/pkg/middleware/meta"
	"google.golang.org/grpc"
)

type Challenges interface {
	SealValue(context.Context, *safe.SealValueRequest, ...grpc.CallOption) (*safe.SealValueResponse, error)
	OpenValue(context.Context, *safe.OpenValueRequest, ...grpc.CallOption) (*safe.OpenValueResponse, error)
}

// Challenge 在认证服务指定的用途下签发一次性挑战，客户端用返回的公钥加密挑战和口令。
func (k *Keys) Challenge(ctx context.Context, remote Challenges, purpose string) (*user.PasswordKeyReply, error) {
	binding, err := binding(ctx)
	if err != nil {
		return nil, err
	}
	ticket, err := remote.SealValue(ctx, &safe.SealValueRequest{Binding: binding, Purpose: purpose, Value: []byte{1}})
	if err != nil {
		return nil, err
	}
	kid, x, y := k.Public()
	return &user.PasswordKeyReply{Kid: kid, X: x, Y: y, Challenge: ticket.Envelope}, nil
}

func (k *Keys) Consume(ctx context.Context, remote Challenges, purpose, envelope string) (string, error) {
	binding, err := binding(ctx)
	if err != nil {
		return "", err
	}
	payload, err := k.Open(envelope)
	if err != nil {
		return "", fault.ErrorBadRequest("口令信封无效")
	}
	_, err = remote.OpenValue(ctx, &safe.OpenValueRequest{Binding: binding, Purpose: purpose, Envelope: payload.Challenge})
	if err != nil {
		return "", err
	}
	return payload.Value, nil
}

func binding(ctx context.Context) (*safe.ClientBinding, error) {
	id, err := meta.GetClientID(ctx)
	if err != nil {
		return nil, err
	}
	return &safe.ClientBinding{ClientId: id, UserAgent: meta.GetUserAgent(ctx)}, nil
}
