// Package dict key
package dict

/* cSpell:disable */

// Key 字典键
type Key string

// String _
func (d Key) String() string {
	return string(d)
}

const (
	// DictGender 性别
	DictGender Key = "gender"
	// DictService 服务
	DictService Key = "service"
	// DictServiceType 服务类型
	DictServiceType Key = "serviceType"
	// DictPlatform 平台
	DictPlatform Key = "platform"
	// DictMiniPlatform 小程序平台
	DictMiniPlatform Key = "mini_platform"
	// DictChannel 渠道
	DictChannel Key = "channel"
	// DictHTTPMethod 请求方式
	DictHTTPMethod Key = "http_method"
	// DictDataType 数据类型
	DictDataType Key = "data_type"
	// DictUseType 用途
	DictUseType Key = "use_type"
	// DictBusiness 业务
	DictBusiness Key = "business"

	// DictOpType 操作类型
	DictOpType Key = "op_type"

	// DictApprovalStatus 审批状态
	DictApprovalStatus Key = "approval_status"

	// DictCloudProvider 云服务商
	DictCloudProvider Key = "cloud_provider"

	// DictCloudProduct 云商品
	DictCloudProduct Key = "cloud_product"
)
