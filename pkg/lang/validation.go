package lang

import (
	"context"

	"buf.build/go/protovalidate"
)

// GetValidateErrorMsg 将标准校验违规转换为本地化提示，兼容数组元素规则缺省的描述符。
func (c *CommonLang) GetValidateErrorMsg(ctx context.Context, violation *protovalidate.Violation, langs ...string) (string, bool) {
	constraintID := violation.Proto.GetRuleId()
	if constraintID == "" {
		return "", false
	}
	msgKey := validatePrefix + MsgKey(constraintID)
	// 判断是否存在
	if _, ok := dfMsgMap[msgKey]; !ok {
		return "", false
	}
	// 数组元素约束会提供值和规则，但不一定附带字段或规则描述符；错误格式化不能因此把 400 变成 panic。
	fieldDesc, ruleDesc := "", ""
	if violation.FieldDescriptor != nil {
		fieldDesc = string(violation.FieldDescriptor.Name())
	}
	if violation.RuleDescriptor != nil {
		ruleDesc = string(violation.RuleDescriptor.Name())
	}
	templateData := map[string]interface{}{
		"Field":     violation.FieldValue.String(),
		"FieldDesc": fieldDesc,
		"Rule":      violation.RuleValue.String(),
		"RuleDesc":  ruleDesc,
		"Message":   violation.Proto.GetMessage(),
	}
	return c.getMsg(ctx, msgKey, templateData, langs...), true
}
