// Package initdata 提供初始化数据的通用功能
// 支持从 URL 或本地配置加载数据，并提供统一的实体初始化逻辑
/* cSpell:disable */
package initdata

import (
	structpb "google.golang.org/protobuf/types/known/structpb"
)

// InitData 初始化的数据配置
type InitData struct {
	// URL 如果有填，则从 URL 获取数据
	URL string `json:"url"`
	// File 如果有填且 URL 为空，则从本地 JSON 文件获取数据
	File string `json:"file"`
	// Data 数组
	// 如果 url/file 为空 data 不为空，则从 data 获取 以 biz 数据为主
	Data []*structpb.Struct `json:"data"`
}
