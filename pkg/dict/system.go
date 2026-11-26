// Package dict system
package dict

/* cSpell:disable */

import (
	"context"

	"github.com/zhixingdx/public/pkg/lang"
	"github.com/zhixingdx/public/pkg/middleware/meta"
)

// Platform 平台
type Platform string

func GetMetaPlatform(ctx context.Context) string {
	v, err := meta.GetValue(ctx, "platform")
	if err != nil {
		return ""
	}
	return v
}

const (
	// PlatformAndroid 安卓
	PlatformAndroid Platform = "android"
	// PlatformIOS 苹果
	PlatformIOS Platform = "ios"
	// PlatformMac Mac
	PlatformMac Platform = "mac"
	// PlatformWindows Windows
	PlatformWindows Platform = "windows"
	// PlatformLinux Linux
	PlatformLinux Platform = "linux"
	// PlatformPad 平板
	PlatformPad Platform = "Pad"
	// PlatformPOS 点单机
	PlatformPOS Platform = "POS"
	// PlatformPDA PDA
	PlatformPDA Platform = "PDA"
	// PlatformOther 其他
	PlatformOther Platform = "other"
)

// PlatformValues _
func PlatformValues() []string {
	return []string{
		string(PlatformAndroid),
		string(PlatformIOS),
		string(PlatformMac),
		string(PlatformWindows),
		string(PlatformLinux),
		string(PlatformPad),
		string(PlatformPOS),
		string(PlatformPDA),
		string(PlatformOther),
	}
}

// 多语言支持
var platformLangMap = map[string]map[string]string{
	"zh": {
		"android": "安卓",
		"ios":     "苹果",
		"mac":     "Mac",
		"windows": "Windows",
		"linux":   "Linux",
		"Pad":     "平板",
		"POS":     "点单机",
		"PDA":     "PDA",
		"other":   "其他",
	},
	"ru": {
		"android": "андроид",
		"ios":     "яблоко",
		"mac":     "мак",
		"windows": "окна",
		"linux":   "линукс",
		"Pad":     "планшет",
		"POS":     "пункт заказа",
		"PDA":     "КПК",
		"other":   "другие",
	},
}

// PlatformLangLabels _
func PlatformLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(platformLangMap, langs)
}

// PlatformLabels _
func PlatformLabels() map[string]string {
	return map[string]string{
		"android": "android",
		"ios":     "ios",
		"mac":     "mac",
		"windows": "windows",
		"linux":   "linux",
		"Pad":     "Pad",
		"POS":     "POS",
		"PDA":     "PDA",
		"other":   "other",
	}
}

// String _
func (p Platform) String() string {
	return string(p)
}

// Label _
func (p Platform) Label() string {
	return PlatformLabels()[p.String()]
}

// LangLabels _
func (p Platform) LangLabel(langs ...string) (string, bool) {
	labels, ok := PlatformLangLabels(langs...)
	if !ok {
		return p.Label(), ok
	}
	return labels[p.String()], ok
}

// Values _
func (p Platform) Values() []string {
	return PlatformValues()
}

// 小程序平台
type MiniPlatform string

func GetMetaMiniPlatform(ctx context.Context) string {
	v, err := meta.GetValue(ctx, "mini-platform")
	if err != nil {
		return ""
	}
	return v
}

const (
	// MiniPlatformWechat 微信小程序
	MiniPlatformWechat MiniPlatform = "wechat"
	// MiniPlatformJD 京东小程序
	MiniPlatformJD MiniPlatform = "jd"
	// MiniPlatformBaidu 百度智能小程序
	MiniPlatformBaidu MiniPlatform = "baidu"
	// MiniPlatformAlipay 支付宝小程序
	MiniPlatformAlipay MiniPlatform = "alipay"
	// MiniPlatformDouyin 抖音小程序
	MiniPlatformDouyin MiniPlatform = "douyin"
	// MiniPlatformQQ QQ 小程序
	MiniPlatformQQ MiniPlatform = "qq"
	// MiniPlatformDingTalk 钉钉小程序
	MiniPlatformDingTalk MiniPlatform = "dingtalk"
	// MiniPlatformWeCom 企业微信小程序
	MiniPlatformWeCom MiniPlatform = "wecom"
	// MiniPlatformAlipayIOT 支付宝 IOT 小程序
	MiniPlatformAlipayIOT MiniPlatform = "alipayiot"
	// MiniPlatformFeishu 飞书小程序
	MiniPlatformFeishu MiniPlatform = "feishu"
	// MiniPlatformOther 其他
	MiniPlatformOther MiniPlatform = "other"
)

// Values _
func MiniPlatformValues() []string {
	return []string{
		string(MiniPlatformWechat),
		string(MiniPlatformJD),
		string(MiniPlatformBaidu),
		string(MiniPlatformAlipay),
		string(MiniPlatformDouyin),
		string(MiniPlatformQQ),
		string(MiniPlatformDingTalk),
		string(MiniPlatformWeCom),
		string(MiniPlatformAlipayIOT),
		string(MiniPlatformFeishu),
		string(MiniPlatformOther),
	}
}

// MiniPlatformLabels _
func MiniPlatformLabels() map[string]string {
	return map[string]string{
		"wechat":    "wechat",
		"jd":        "jd",
		"baidu":     "baidu",
		"alipay":    "alipay",
		"douyin":    "douyin",
		"qq":        "QQ",
		"dingtalk":  "dingtalk",
		"wecom":     "wecom",
		"alipayiot": "alipayiot",
		"feishu":    "feishu",
		"other":     "other",
	}
}

// 多语言支持
var miniPlatformLangMap = map[string]map[string]string{
	"zh": {
		"wechat":    "微信",
		"jd":        "京东",
		"baidu":     "百度",
		"alipay":    "支付宝",
		"douyin":    "抖音",
		"qq":        "QQ",
		"dingtalk":  "钉钉",
		"wecom":     "企业微信",
		"alipayiot": "支付宝 IOT",
		"feishu":    "飞书",
		"other":     "其他",
	},
	"ru": {
		"wechat":    "вэчат",
		"jd":        "джон джон",
		"baidu":     "байду",
		"alipay":    "альпай",
		"douyin":    "доуин",
		"qq":        "QQ",
		"dingtalk":  "дингталк",
		"wecom":     "вэком",
		"alipayiot": "альпай IOT",
		"feishu":    "фейшу",
		"other":     "другие",
	},
}

// MiniPlatformLangLabels _
func MiniPlatformLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(miniPlatformLangMap, langs)
}

// String _
func (m MiniPlatform) String() string {
	return string(m)
}

// Label _
func (m MiniPlatform) Label() string {
	return MiniPlatformLabels()[m.String()]
}

// LangLabel _
func (m MiniPlatform) LangLabel(langs ...string) (string, bool) {
	labels, ok := MiniPlatformLangLabels(langs...)
	if !ok {
		return m.Label(), ok
	}
	return labels[m.String()], ok
}

// Values _
func (MiniPlatform) Values() []string {
	return MiniPlatformValues()
}

// Channel 渠道
type Channel string

func GetMetaChannel(ctx context.Context) string {
	v, err := meta.GetValue(ctx, "channel")
	if err != nil {
		return ""
	}
	return v
}

const (
	// ChannelWeb 网页
	ChannelWeb Channel = "web"
	// ChannelH5 H5
	ChannelH5 Channel = "h5"
	// ChannelApp APP
	ChannelApp Channel = "app"
	// ChannelMiniprogram 小程序
	ChannelMiniprogram Channel = "miniprogram"
	// ChannelDesktop 桌面
	ChannelDesktop Channel = "desktop"
	// ChannelEmbed 嵌入式
	ChannelEmbed Channel = "embed"
)

// Values _
func ChannelValues() []string {
	return []string{
		string(ChannelWeb),
		string(ChannelH5),
		string(ChannelApp),
		string(ChannelMiniprogram),
		string(ChannelDesktop),
		string(ChannelEmbed),
	}
}

// ChannelLabels _
func ChannelLabels() map[string]string {
	return map[string]string{
		"web":         "web",
		"h5":          "H5",
		"app":         "APP",
		"miniprogram": "miniprogram",
		"desktop":     "desktop",
		"embed":       "embed",
	}
}

// 多语言支持
var channelLangMap = map[string]map[string]string{
	"zh": {
		"web":         "网页",
		"h5":          "H5",
		"app":         "APP",
		"miniprogram": "小程序",
		"desktop":     "桌面应用",
		"embed":       "嵌入式",
	},
	"ru": {
		"web":         "веб-страница",
		"h5":          "H5",
		"app":         "Приложение",
		"miniprogram": "мини-программа",
		"desktop":     "Приложение для настольного компьютера",
		"embed":       "встроенный",
	},
}

// ChannelLangLabels _
func ChannelLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(channelLangMap, langs)
}

// String _
func (c Channel) String() string {
	return string(c)
}

// Label _
func (c Channel) Label() string {
	return ChannelLabels()[c.String()]
}

// LangLabel _
func (c Channel) LangLabel(langs ...string) (string, bool) {
	labels, ok := ChannelLangLabels(langs...)
	if !ok {
		return c.Label(), ok
	}
	return labels[c.String()], ok
}

// Values _
func (Channel) Values() []string {
	return ChannelValues()
}

// HTTPMethod 请求方式
type HTTPMethod string

// HTTPMethod values.
const (
	HTTPMethodGet     HTTPMethod = "GET"
	HTTPMethodPost    HTTPMethod = "POST"
	HTTPMethodPut     HTTPMethod = "PUT"
	HTTPMethodDelete  HTTPMethod = "DELETE"
	HTTPMethodPatch   HTTPMethod = "PATCH"
	HTTPMethodOptions HTTPMethod = "OPTIONS"
	HTTPMethodHead    HTTPMethod = "HEAD"
)

// HTTPMethodLabels _
func HTTPMethodLabels() map[string]string {
	return map[string]string{
		"GET":     "GET",
		"POST":    "POST",
		"PUT":     "PUT",
		"DELETE":  "DELETE",
		"PATCH":   "PATCH",
		"OPTIONS": "OPTIONS",
		"HEAD":    "HEAD",
	}
}

// HTTPMethodValues http method 值
func HTTPMethodValues() []string {
	return []string{
		string(HTTPMethodGet),
		string(HTTPMethodPost),
		string(HTTPMethodPut),
		string(HTTPMethodDelete),
		string(HTTPMethodPatch),
		string(HTTPMethodOptions),
		string(HTTPMethodHead),
	}
}

// String _
func (m HTTPMethod) String() string {
	return string(m)
}

// Label _
func (m HTTPMethod) Label() string {
	return HTTPMethodLabels()[m.String()]
}

// Values _
func (HTTPMethod) Values() []string {
	return HTTPMethodValues()
}

// DataType 数据类型
type DataType string

const (
	// DataTypeString 字符串
	DataTypeString DataType = "string"
	// DataTypeNumber 数字
	DataTypeNumber DataType = "number"
	// DataTypeBoolean 布尔
	DataTypeBoolean DataType = "boolean"
	// DataTypeArray 数组
	DataTypeArray DataType = "array"
	// DataTypeObject 对象
	DataTypeObject DataType = "object"
	// DataTypeNull 空
	DataTypeNull DataType = "null"
)

// DataTypeValues _
func DataTypeValues() []string {
	return []string{
		string(DataTypeString),
		string(DataTypeNumber),
		string(DataTypeBoolean),
		string(DataTypeArray),
		string(DataTypeObject),
		string(DataTypeNull),
	}
}

// DataTypeLabels _
func DataTypeLabels() map[string]string {
	return map[string]string{
		"string":  "string",
		"number":  "number",
		"boolean": "boolean",
		"array":   "array",
		"object":  "object",
		"null":    "null",
	}
}

// 多语言支持
var dataTypeLangMap = map[string]map[string]string{
	"zh": {
		"string":  "字符串",
		"number":  "数字",
		"boolean": "布尔",
		"array":   "数组",
		"object":  "对象",
		"null":    "空",
	},
	"ru": {
		"string":  "строка",
		"number":  "число",
		"boolean": "логический",
		"array":   "массив",
		"object":  "объект",
		"null":    "нуль",
	},
}

// DataTypeLangLabels _
func DataTypeLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(dataTypeLangMap, langs)
}

// String _
func (d DataType) String() string {
	return string(d)
}

// Label _
func (d DataType) Label() string {
	return DataTypeLabels()[d.String()]
}

// LangLabel _
func (d DataType) LangLabel(langs ...string) (string, bool) {
	labels, ok := DataTypeLangLabels(langs...)
	if !ok {
		return d.Label(), ok
	}
	return labels[d.String()], ok
}

// Values _
func (DataType) Values() []string {
	return DataTypeValues()
}

// UseType 用于区分 saas 系统里不同的使用者
type UseType string

func GetMetaUseType(ctx context.Context) string {
	v, err := meta.GetValue(ctx, "use-type")
	if err != nil {
		return ""
	}
	return v
}

const (
	// UseTypeAdmin 管理员 整套系统的管理员
	UseTypeAdmin UseType = "admin"
	// UseTypeTenant 租户 租户的管理/运营人员
	UseTypeTenant UseType = "tenant"
	// UseTypeServicer 服务人员 为租户提供服务的人员 如 骑手/配送员/服务员
	UseTypeServicer UseType = "servicer"
	// UseTypeClient 客户 C 端用户 使用租户提供的服务的人
	UseTypeClient UseType = "client"
	// UseTypeOther 其他
	UseTypeOther UseType = "other"
)

// UseTypeValues _
func UseTypeValues() []string {
	return []string{
		string(UseTypeAdmin),
		string(UseTypeTenant),
		string(UseTypeServicer),
		string(UseTypeClient),
		string(UseTypeOther),
	}
}

// UseTypeLabels 返回使用主体类型的协议标签。
func UseTypeLabels() map[string]string {
	return map[string]string{
		"admin":    "admin",
		"tenant":   "tenant",
		"servicer": "servicer",
		"client":   "client",
		"other":    "other",
	}
}

// 多语言支持
var useTypeLangMap = map[string]map[string]string{
	"zh": {
		"admin":    "管理员",
		"tenant":   "单位",
		"servicer": "服务人员",
		"client":   "客户",
		"other":    "其他",
	},
	"ru": {
		"admin":    "администратор",
		"tenant":   "организация",
		"servicer": "обслуживающий персонал",
		"client":   "клиент",
		"other":    "другие",
	},
}

// UseTypeLangLabels _
func UseTypeLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(useTypeLangMap, langs)
}

// String _
func (u UseType) String() string {
	return string(u)
}

// Label _
func (u UseType) Label() string {
	return UseTypeLabels()[u.String()]
}

// LangLabel _
func (u UseType) LangLabel(langs ...string) (string, bool) {
	labels, ok := UseTypeLangLabels(langs...)
	if !ok {
		return u.Label(), ok
	}
	return labels[u.String()], ok
}

// Values _
func (UseType) Values() []string {
	return UseTypeValues()
}

// Business 业务
type Business string

func GetMetaBusiness(ctx context.Context) string {
	v, err := meta.GetValue(ctx, "business")
	if err != nil {
		return ""
	}
	return v
}

const (
	// 教育
	BusinessEducation Business = "education"
	// BusinessFood 餐饮
	BusinessFood Business = "food"
	// BusinessRetail 零售
	BusinessRetail Business = "retail"
	// BusinessHotel 酒店
	BusinessHotel Business = "hotel"
	// BusinessBeauty 美容
	BusinessBeauty Business = "beauty"
	// BusinessOther 其他
	BusinessOther Business = "other"
)

// BusinessValues _
func BusinessValues() []string {
	return []string{
		string(BusinessEducation),
		string(BusinessFood),
		string(BusinessRetail),
		string(BusinessHotel),
		string(BusinessBeauty),
		string(BusinessOther),
	}
}

// BusinessLabels _
func BusinessLabels() map[string]string {
	return map[string]string{
		"education": "教育",
		"food":      "餐饮",
		"retail":    "零售",
		"hotel":     "酒店",
		"beauty":    "美容",
		"other":     "其他",
	}
}

// 多语言支持
var businessLangMap = map[string]map[string]string{
	"zh": {
		"education": "教育",
		"food":      "餐饮",
		"retail":    "零售",
		"hotel":     "酒店",
		"beauty":    "美容",
		"other":     "其他",
	},
	"ru": {
		"education": "образование",
		"food":      "питание",
		"retail":    "розничная торговля",
		"hotel":     "отель",
		"beauty":    "красота",
		"other":     "другие",
	},
}

// BusinessLangLabels _
func BusinessLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(businessLangMap, langs)
}

// String _
func (b Business) String() string {
	return string(b)
}

// Label _
func (b Business) Label() string {
	return BusinessLabels()[b.String()]
}

// LangLabel _
func (b Business) LangLabel(langs ...string) (string, bool) {
	labels, ok := BusinessLangLabels(langs...)
	if !ok {
		return b.Label(), ok
	}
	return labels[b.String()], ok
}

// Values _
func (Business) Values() []string {
	return BusinessValues()
}
