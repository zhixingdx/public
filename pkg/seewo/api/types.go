package api

import (
	"fmt"

	"github.com/zhixingdx/oipivot-api/fault"
)

// Response 希沃开放平台通用响应结构
// 参考 Java SDK: com.seewo.open.sdk.DefaultSeewoClient.errorDetect
// 源码位置: /private/tmp/com/seewo/open/sdk/DefaultSeewoClient.class
//
// 错误响应示例:
//
//	{"exception": "ErrorType", "message": "错误描述"}
//
// 成功响应示例:
//
//	{"data": {...}}  或  [...]
type Response[T any] struct {
	// Exception 错误标识，非空表示业务错误
	Exception string `json:"exception,omitempty"`
	// Message 错误消息
	Message string `json:"message,omitempty"`
	// Data 业务数据
	Data T `json:"data,omitempty"`
}

// CheckError 校验接口返回是否有错误
func (r *Response[T]) CheckError() error {
	if r.Exception != "" {
		msg := r.Message
		if msg == "" {
			msg = "未知错误"
		}
		return fault.ErrorBadGateway("希沃开放平台接口错误: exception=%s, message=%s", r.Exception, msg)
	}
	return nil
}

// ListResponse 列表响应结构（部分接口返回数组）
type ListResponse[T any] struct {
	// Exception 错误标识
	Exception string `json:"exception,omitempty"`
	// Message 错误消息
	Message string `json:"message,omitempty"`
	// Data 业务数据列表
	Data []T `json:"data,omitempty"`
}

// CheckError 校验接口返回是否有错误
func (r *ListResponse[T]) CheckError() error {
	if r.Exception != "" {
		msg := r.Message
		if msg == "" {
			msg = "未知错误"
		}
		return fault.ErrorBadGateway("希沃开放平台接口错误: exception=%s, message=%s", r.Exception, msg)
	}
	return nil
}

// FlatResponse 扁平响应结构（数据直接在顶层，不包在 data 字段中）
// 某些接口（如 mdm-basis 系列）成功时直接返回数据对象 {"uid":"...","name":"..."}
// 错误时仍返回 {"exception":"...","message":"..."}
// 使用方式：嵌入目标结构体，如 FlatResponse[OrgInfo]
type FlatResponse struct {
	// Exception 错误标识，非空表示业务错误
	Exception string `json:"exception,omitempty"`
	// Message 错误消息
	Message string `json:"message,omitempty"`
}

// CheckError 校验接口返回是否有错误
func (r *FlatResponse) CheckError() error {
	if r.Exception != "" {
		msg := r.Message
		if msg == "" {
			msg = "未知错误"
		}
		return fault.ErrorBadGateway("希沃开放平台接口错误: exception=%s, message=%s", r.Exception, msg)
	}
	return nil
}

// PagedData 分页数据结构
// 某些接口（如 mdm-building-room-api）返回带分页的数据
// 响应格式: {"code":"...", "message":"...", "data":{"total":N, "pageable":{...}, "content":[...]}}
type PagedData[T any] struct {
	Total    int64    `json:"total"`
	Pageable Pageable `json:"pageable"`
	Content  []T      `json:"content"`
}

// Pageable 分页信息
type Pageable struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

// PagedResponse 带分页的通用响应结构
// 响应格式: {"code":"...", "message":"...", "data":{"total":N, "pageable":{...}, "content":[...]}}
type PagedResponse[T any] struct {
	Code    string       `json:"code,omitempty"`
	Message string       `json:"message,omitempty"`
	Data    PagedData[T] `json:"data"`
}

// CheckError 校验接口返回是否有错误
// code 为空或 "0" 表示成功
func (r *PagedResponse[T]) CheckError() error {
	if r.Code != "" && r.Code != "0" {
		msg := r.Message
		if msg == "" {
			msg = "未知错误"
		}
		return fault.ErrorBadGateway("希沃开放平台接口错误: code=%s, message=%s", r.Code, msg)
	}
	return nil
}

// RawResponse 原始响应结构（用于不确定 data 类型的情况）
type RawResponse struct {
	Exception string      `json:"exception,omitempty"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

// CheckError 校验接口返回是否有错误
func (r *RawResponse) CheckError() error {
	if r.Exception != "" {
		msg := r.Message
		if msg == "" {
			msg = "未知错误"
		}
		return fault.ErrorBadGateway("希沃开放平台接口错误: exception=%s, message=%s", r.Exception, msg)
	}
	return nil
}

// SeewoError 希沃 API 错误
type SeewoError struct {
	StatusCode int    // HTTP 状态码
	Exception  string // 错误类型
	Message    string // 错误消息
	ReqID      string // 请求 ID（x-sw-req-id）
}

// Error 实现 error 接口
func (e *SeewoError) Error() string {
	return fmt.Sprintf("希沃API错误: statusCode=%d, exception=%s, message=%s, reqId=%s",
		e.StatusCode, e.Exception, e.Message, e.ReqID)
}
