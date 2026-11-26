package seewo

import "github.com/zhixingdx/oipivot-api/fault"

var (
	// ErrConfigInvalidAPI 配置中 API 地址无效
	ErrConfigInvalidAPI = fault.ErrorBadRequest("希沃开放平台配置中 API 地址不能为空")
	// ErrConfigInvalidAppID 配置中 AppID 无效
	ErrConfigInvalidAppID = fault.ErrorBadRequest("希沃开放平台配置中 AppID 不能为空")
	// ErrConfigInvalidAppSecret 配置中 AppSecret 无效
	ErrConfigInvalidAppSecret = fault.ErrorBadRequest("希沃开放平台配置中 AppSecret 不能为空")
	// ErrRequestFailed 请求失败
	ErrRequestFailed = fault.ErrorBadGateway("希沃开放平台请求失败")
	// ErrInvalidResponse 响应格式无效
	ErrInvalidResponse = fault.ErrorBadGateway("希沃开放平台响应格式无效")
	// ErrSignFailed 签名生成失败
	ErrSignFailed = fault.ErrorInternalServerError("希沃开放平台签名生成失败")
	// ErrAPIResponseFailed 接口返回业务错误
	ErrAPIResponseFailed = fault.ErrorBadGateway("希沃开放平台接口返回失败")
)
