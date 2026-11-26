/* cSpell:disable */
package dict

import "github.com/zhixingdx/oipivot/backend/pkg/lang"

type GrantType string

// 1. 授权码(Authorization Code):
// 最常见的授权类型，适用于浏览器应用和需要用户交互的场景。流程是：客户端重定向用户到Hydra 进行认证，用户同意授权后，Hydra 返回一个授权码，客户端使用该授权码向Hydra 换取访问令牌。
// 2. 隐式(Implicit):
// 适用于不适合使用授权码流程的场景，例如，直接在浏览器中运行的客户端(如JavaScript 应用程序)。流程是：客户端重定向用户到Hydra 进行认证，Hydra 直接返回访问令牌。
// 3. 密码(Resource Owner Password Credentials):
// 适用于信任的客户端，例如，公司内部的应用。客户端直接使用用户的用户名和密码向Hydra 换取访问令牌。
// 4. 客户端凭证(Client Credentials):
// 适用于服务器到服务器的交互，例如，一个服务需要访问另一个服务的数据。客户端直接使用自己的客户端ID 和密钥向Hydra 换取访问令牌。
// 5. 设备码(Device Code):
// 适用于需要用户交互的场景，例如，一个设备需要访问另一个设备的数据。客户端直接使用自己的客户端ID 和密钥向Hydra 换取访问令牌。

// 授权类型
const (
	GrantTypeAuthorizationCode GrantType = "authorization_code"
	GrantTypeImplicit          GrantType = "implicit"
	GrantTypePassword          GrantType = "password"
	GrantTypeClientCredentials GrantType = "client_credentials"
	GrantTypeDeviceCode        GrantType = "device_code"
)

// GrantTypeValues 返回所有授权类型的值
func GrantTypeValues() []string {
	return []string{
		string(GrantTypeAuthorizationCode),
		string(GrantTypeImplicit),
		string(GrantTypePassword),
		string(GrantTypeClientCredentials),
		string(GrantTypeDeviceCode),
	}
}

// GrantTypeLabels 返回所有授权类型的标签
func GrantTypeLabels() map[string]string {
	return map[string]string{
		"authorization_code": "Authorization Code",
		"implicit":           "Implicit",
		"password":           "Password",
		"client_credentials": "Client Credentials",
		"device_code":        "Device Code",
	}
}

// 多语言支持
var grantTypeLangMap = map[string]map[string]string{
	"zh": {
		"authorization_code": "授权码",
		"implicit":           "隐式",
		"password":           "密码",
		"client_credentials": "客户端凭证",
		"device_code":        "设备码",
	},
	"ru": {
		"authorization_code": "Авторизация",
		"implicit":           "Неявное",
		"password":           "Пароль",
		"client_credentials": "Клиентские учетные данные",
		"device_code":        "Код устройства",
	},
}

func GrantTypeLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(grantTypeLangMap, langs)
}

// String 返回字符串
func (g GrantType) String() string {
	return string(g)
}

// Label 返回标签
func (g GrantType) Label() string {
	return GrantTypeLabels()[g.String()]
}

// Values 返回所有授权类型的值
func (g GrantType) Values() []string {
	return GrantTypeValues()
}

// LangLabel 返回指定语言的授权类型标签
func (g GrantType) LangLabel(langs ...string) (string, bool) {
	labels, b := GrantTypeLangLabels(langs...)
	if !b {
		return g.Label(), false
	}
	return labels[g.String()], true
}
