package seewo

import (
	"crypto/hmac"
	"crypto/md5"
	"fmt"
	"sort"
	"strings"
)

const (
	// SignMethodHMAC HMAC-MD5 签名方式
	SignMethodHMAC = "hmac"
	// SignMethodMD5 MD5 签名方式
	SignMethodMD5 = "md5"
)

// SignRequest 生成请求签名
// 签名算法参考 Java SDK: com.seewo.open.sdk.util.SignUtils.signRequest
// 源码位置: /private/tmp/com/seewo/open/sdk/util/SignUtils.class
//
// 签名流程:
//  1. 对参数 key 进行字典序排序
//  2. 拼接待签名字符串（跳过 key 或 value 为空的参数）
//  3. 根据签名方式计算签名值
//
// 参数:
//   - params: 参与签名的所有参数（headers + queries，不含 x-sw-sign）
//   - secret: AppSecret
//   - signMethod: 签名方式，"hmac" 或 "md5"
//
// 返回:
//   - 签名值（大写十六进制字符串）
func SignRequest(params map[string]string, secret, signMethod string) (string, error) {
	// 1. 获取所有 key 并排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 2. 构建待签名字符串
	var sb strings.Builder

	// MD5 方式先拼接 secret
	if signMethod == SignMethodMD5 {
		sb.WriteString(secret)
	}

	// 遍历排序后的参数，拼接 key + value（跳过空值）
	for _, key := range keys {
		value := params[key]
		if key != "" && value != "" {
			sb.WriteString(key)
			sb.WriteString(value)
		}
	}

	// 3. 计算签名
	var signBytes []byte
	if signMethod == SignMethodHMAC {
		// HMAC-MD5 签名
		signBytes = encryptHMAC(sb.String(), secret)
	} else {
		// MD5 签名：末尾再拼接 secret
		sb.WriteString(secret)
		signBytes = encryptMD5(sb.String())
	}

	// 4. 转为大写十六进制字符串
	return byte2hex(signBytes), nil
}

// encryptHMAC 使用 HMAC-MD5 加密
// 参考 Java SDK: com.seewo.open.sdk.util.SignUtils.encryptHMAC
func encryptHMAC(data, secret string) []byte {
	mac := hmac.New(md5.New, []byte(secret))
	mac.Write([]byte(data))
	return mac.Sum(nil)
}

// encryptMD5 使用 MD5 加密
// 参考 Java SDK: com.seewo.open.sdk.util.SignUtils.encryptMD5
func encryptMD5(data string) []byte {
	hash := md5.Sum([]byte(data))
	return hash[:]
}

// byte2hex 字节数组转大写十六进制字符串
// 参考 Java SDK: com.seewo.open.sdk.util.SignUtils.byte2hex
func byte2hex(bytes []byte) string {
	return strings.ToUpper(fmt.Sprintf("%x", bytes))
}

// GetContentMD5 计算请求体的 MD5 值（大写十六进制）
// 参考 Java SDK: com.seewo.open.sdk.util.SignUtils.getContentMD5
func GetContentMD5(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	hash := md5.Sum(data)
	return strings.ToUpper(fmt.Sprintf("%x", hash[:]))
}
