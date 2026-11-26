package api

import "context"

// UserAPI 用户相关 API
type UserAPI struct {
	client Client
}

// NewUserAPI 创建用户 API 实例
func NewUserAPI(client Client) *UserAPI {
	return &UserAPI{client: client}
}

// ---- 请求/响应结构体 ----

// UserPhoneInfo 用户手机号信息
// 参考 Java SDK: MdmUserApiGetOrCreateUserUidByPhoneResult.UserPhoneDto
type UserPhoneInfo struct {
	Phone   string `json:"phone,omitempty"`   // 手机号
	UserUID string `json:"userUid,omitempty"` // 用户ID
}

// GetOrCreateUserUIDByPhoneRequest 根据手机号创建或获取用户ID请求
// 参考 Java SDK: MdmUserApiGetOrCreateUserUidByPhoneParam
// 路径: POST /mdm-basis/mdm-user-api/get-or-create-user-uid-by-phone
type GetOrCreateUserUIDByPhoneRequest struct {
	ReqDto GetOrCreateUserUIDByPhoneReqDto `json:"reqDto"`
}

// GetOrCreateUserUIDByPhoneReqDto 根据手机号创建或获取用户ID请求参数
type GetOrCreateUserUIDByPhoneReqDto struct {
	Phones []string `json:"phones"` // 手机号列表
}

// GetOrCreateUserUIDByPhoneResponse 根据手机号创建或获取用户ID响应
// 参考 Java SDK: MdmUserApiGetOrCreateUserUidByPhoneResult.GetOrCreateUserUidByPhoneRespDto
type GetOrCreateUserUIDByPhoneResponse struct {
	Users []UserPhoneInfo `json:"users,omitempty"` // 用户列表
}

// ---- API 方法 ----

// GetOrCreateUserUIDByPhone 根据手机号创建或获取用户ID
// 参考 Java SDK: MdmUserApiGetOrCreateUserUidByPhoneRequest
// POST /mdm-basis/mdm-user-api/get-or-create-user-uid-by-phone
func (a *UserAPI) GetOrCreateUserUIDByPhone(ctx context.Context, req *GetOrCreateUserUIDByPhoneRequest) (*GetOrCreateUserUIDByPhoneResponse, error) {
	var resp Response[GetOrCreateUserUIDByPhoneResponse]
	if err := a.client.PostJSON(ctx, "/mdm-basis/mdm-user-api/get-or-create-user-uid-by-phone", req, &resp); err != nil {
		return nil, err
	}
	if err := resp.CheckError(); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
