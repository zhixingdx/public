package api

import "context"

// OrgAPI 组织相关 API
type OrgAPI struct {
	client Client
}

// NewOrgAPI 创建组织 API 实例
func NewOrgAPI(client Client) *OrgAPI {
	return &OrgAPI{client: client}
}

// ---- 请求/响应结构体 ----

// OrgInfo 组织信息（基础）
// 参考 Java SDK: MdmOrgApiQueryOrgByCodeResult.ResponseBody
type OrgInfo struct {
	UID  string `json:"uid,omitempty"`  // 组织 UID
	Name string `json:"name,omitempty"` // 组织名称
	Code string `json:"code,omitempty"` // 组织编码
}

// OrgDetailInfo 组织详细信息
// 参考 Java SDK: MdmOrgApiQueryOrgByUidResult.JSONResponseBody
type OrgDetailInfo struct {
	UID         string  `json:"uid,omitempty"`         // 组织id
	Name        string  `json:"name,omitempty"`       // 组织名称
	Code        string  `json:"code,omitempty"`        // 组织编码
	SchoolType  int     `json:"schoolType,omitempty"` // 魔方学校类型
	SchoolMotto string  `json:"schoolMotto,omitempty"` // 校训
	Longitude   float64 `json:"longitude,omitempty"`  // 经度
	Latitude    float64 `json:"latitude,omitempty"`    // 纬度
	Address     string  `json:"address,omitempty"`     // 学校地址
	ProvinceID  string  `json:"provinceId,omitempty"`  // 省编码
	Province    string  `json:"province,omitempty"`    // 省名称
	CityID      string  `json:"cityId,omitempty"`      // 市编码
	City        string  `json:"city,omitempty"`        // 市名称
	DistrictID  string  `json:"districtId,omitempty"` // 区县编码
	District    string  `json:"district,omitempty"`    // 区县名称
}

// GetOrgByCodeRequest 通过编码查询组织信息请求
// 参考 Java SDK: MdmOrgApiQueryOrgByCodeParam
// 路径: POST /mdm-basis/mdm-org-api/query-org-by-code
// 请求体结构: {"reqDto": {"code": "xxx"}}
type GetOrgByCodeRequest struct {
	ReqDto GetOrgByCodeReqDto `json:"reqDto"`
}

// GetOrgByCodeReqDto 查询组织请求参数
type GetOrgByCodeReqDto struct {
	Code string `json:"code"` // 组织编码
}

// GetOrgByUIDRequest 通过 UID 查询组织信息请求
// 参考 Java SDK: MdmOrgApiQueryOrgByUidParam
// 路径: POST /mdm-basis/mdm-org-api/query-org-by-uid
type GetOrgByUIDRequest struct {
	ReqDto GetOrgByUIDReqDto `json:"reqDto"`
}

// GetOrgByUIDReqDto 通过 UID 查询组织请求参数
type GetOrgByUIDReqDto struct {
	OrgUID string `json:"orgUid"` // 组织编码
}

// ---- API 方法 ----

// GetOrgByCode 通过编码查询组织信息
// 参考 Java SDK: MdmOrgApiQueryOrgByCodeRequest
// POST /mdm-basis/mdm-org-api/query-org-by-code
// 注意：此接口直接返回 OrgInfo 对象（非 {"data":...} 嵌套）
func (a *OrgAPI) GetOrgByCode(ctx context.Context, req *GetOrgByCodeRequest) (*OrgInfo, error) {
	var resp struct {
		FlatResponse
		OrgInfo
	}
	if err := a.client.PostJSON(ctx, "/mdm-basis/mdm-org-api/query-org-by-code", req, &resp); err != nil {
		return nil, err
	}
	if err := resp.FlatResponse.CheckError(); err != nil {
		return nil, err
	}
	return &resp.OrgInfo, nil
}

// GetOrgByUID 通过 UID 查询组织详细信息
// 参考 Java SDK: MdmOrgApiQueryOrgByUidRequest
// POST /mdm-basis/mdm-org-api/query-org-by-uid
// 注意：此接口直接返回 OrgDetailInfo 对象（非 {"data":...} 嵌套）
func (a *OrgAPI) GetOrgByUID(ctx context.Context, req *GetOrgByUIDRequest) (*OrgDetailInfo, error) {
	var resp struct {
		FlatResponse
		OrgDetailInfo
	}
	if err := a.client.PostJSON(ctx, "/mdm-basis/mdm-org-api/query-org-by-uid", req, &resp); err != nil {
		return nil, err
	}
	if err := resp.FlatResponse.CheckError(); err != nil {
		return nil, err
	}
	return &resp.OrgDetailInfo, nil
}
