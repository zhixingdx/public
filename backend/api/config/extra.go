package config

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/known/structpb"
	json "sigs.k8s.io/json"
)

// DecodeExtra 将扩展配置严格解码为服务声明的 Go 结构体，缺省配置返回零值。
func DecodeExtra[T any](extra *structpb.Struct) (*T, error) {
	var result T
	if extra == nil {
		return &result, nil
	}
	for name, value := range extra.GetFields() {
		if value == nil || value.GetKind() == nil {
			return nil, fmt.Errorf("extra.%s 为空", name)
		}
		if _, null := value.Kind.(*structpb.Value_NullValue); null {
			return nil, fmt.Errorf("extra.%s 不能为 null", name)
		}
	}
	raw, err := extra.MarshalJSON()
	if err != nil {
		return nil, err
	}
	issues, err := json.UnmarshalStrict(raw, &result)
	if err != nil {
		return nil, err
	}
	if err := errors.Join(issues...); err != nil {
		return nil, err
	}
	return &result, nil
}
