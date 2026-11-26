package dict

import "github.com/zhixingdx/oipivot/backend/pkg/lang"

// 云服务商
type CloudProvider string

const (
	// CloudProviderAliyun 阿里云
	CloudProviderAliyun CloudProvider = "aliyun"
	// CloudProviderTencent 腾讯云
	CloudProviderTencent CloudProvider = "tencent"
	// CloudProviderVolcengine 火山引擎
	CloudProviderVolcengine CloudProvider = "volcengine"
	// CloudProviderHuawei 华为云
	CloudProviderHuawei CloudProvider = "huawei"
	// CloudProviderAWS 亚马逊云
	CloudProviderAWS CloudProvider = "aws"
	// CloudProviderAzure 微软云
	CloudProviderAzure CloudProvider = "azure"
)

// CloudProviderValues 返回所有云服务商的值
func CloudProviderValues() []string {
	return []string{
		string(CloudProviderAliyun),
		string(CloudProviderTencent),
		string(CloudProviderVolcengine),
		string(CloudProviderHuawei),
		string(CloudProviderAWS),
		string(CloudProviderAzure),
	}
}

// CloudProviderLabels 返回所有云服务商的标签
func CloudProviderLabels() map[string]string {
	return map[string]string{
		"aliyun":     "aliyun",
		"tencent":    "tencent",
		"volcengine": "volcengine",
		"huawei":     "huawei",
		"aws":        "aws",
		"azure":      "azure",
	}
}

// 多语言支持
var cloudProviderLangMap = map[string]map[string]string{
	"zh": {
		"aliyun":     "阿里云",
		"tencent":    "腾讯云",
		"volcengine": "火山引擎",
		"huawei":     "华为云",
		"aws":        "亚马逊云",
		"azure":      "微软云",
	},
}

// CloudProviderLangLabels 返回指定语言的云服务商标签
func CloudProviderLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(cloudProviderLangMap, langs)
}

// String 返回字符串
func (c CloudProvider) String() string {
	return string(c)
}

// Label 返回标签
func (c CloudProvider) Label() string {
	return CloudProviderLabels()[c.String()]
}

// Values _
func (c CloudProvider) Values() []string {
	return CloudProviderValues()
}

// LangLabel
func (c CloudProvider) LangLabel(langs ...string) (string, bool) {
	labels, b := CloudProviderLangLabels(langs...)
	if !b {
		return c.Label(), false
	}
	return labels[c.String()], true
}

type CloudProduct string

const (
	CloudProductSSL CloudProduct = "ssl"
	// DNS
	CloudProductDNS CloudProduct = "dns"
)

// CloudProviderValues 返回所有云服务商的值
func CloudProductValues() []string {
	return []string{
		string(CloudProductSSL),
		string(CloudProductDNS),
	}
}

// CloudProductLabels 返回所有云商品的标签
func CloudProductLabels() map[string]string {
	return map[string]string{
		"ssl": "ssl",
		"dns": "dns",
	}
}

// 多语言支持
var cloudProductLangMap = map[string]map[string]string{
	"zh": {
		"ssl": "证书",
		"dns": "DNS",
	},
}

// CloudProductLangLabels 返回指定语言的云商品标签
func CloudProductLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(cloudProductLangMap, langs)
}

// String 返回字符串
func (c CloudProduct) String() string {
	return string(c)
}

// Label 返回标签
func (c CloudProduct) Label() string {
	return CloudProductLabels()[c.String()]
}

// Values _
func (c CloudProduct) Values() []string {
	return CloudProductValues()
}

// LangLabel
func (c CloudProduct) LangLabel(langs ...string) (string, bool) {
	labels, b := CloudProductLangLabels(langs...)
	if !b {
		return c.Label(), false
	}
	return labels[c.String()], true
}
