// Package lang bundle common
package lang

import (
	"context"
	stderrors "errors"
	"strings"

	kerrors "github.com/go-kratos/kratos/v3/errors"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/zhixingdx/oipivot-api/config"
	"github.com/zhixingdx/oipivot-api/fault"
	"github.com/zhixingdx/oipivot/backend/pkg/logger/kratoslog"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CommonLang 公共语言包
type CommonLang struct {
	Bundle *i18n.Bundle
	log    *log.Helper
}

// GetCommonBundle 获取公共语言包
func GetCommonBundle(langMapConfig map[string]*config.Lang) *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	LoadMessage(bundle, langMap)
	LoadMessageForConfig(bundle, langMapConfig)
	return bundle
}

// NewCommonLang 创建公共语言包
func NewCommonLang(langMapConfig map[string]*config.Lang, logger log.Logger) *CommonLang {
	return &CommonLang{
		Bundle: GetCommonBundle(langMapConfig),
		log:    log.NewHelper(log.With(logger, "layer", "commonLang")),
	}
}

// HandleMetadataError 处理 metadata 错误
func (c *CommonLang) HandleMetadataError(ctx context.Context, err *kerrors.Error, langs ...string) error {
	if err == nil {
		return nil
	}
	if err.Metadata == nil {
		return fault.ErrorBadRequest("%s", c.GetMetadataFailMsg(ctx, langs...))
	}
	source, sOk := err.Metadata["source"]
	target, tOk := err.Metadata["target"]

	if !sOk || !tOk {
		return fault.ErrorBadRequest("%s", c.GetMetadataFailMsg(ctx, langs...))
	}
	return fault.ErrorBadRequest("%s", c.GetMetadataConversionFailMsg(ctx, source, target, langs...))
}

// GetLocalizer _
func (c *CommonLang) getLocalizer(ctx context.Context, langs ...string) *i18n.Localizer {
	if len(langs) == 0 {
		return i18n.NewLocalizer(c.Bundle, GetAcceptArr(ctx)...)
	}
	return i18n.NewLocalizer(c.Bundle, langs...)
}

// getMsg _
func (c *CommonLang) getMsg(ctx context.Context, key MsgKey, templateData interface{}, langs ...string) string {
	localizer := c.getLocalizer(ctx, langs...)
	v, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:      key.String(),
		TemplateData:   templateData,
		DefaultMessage: dfMsgMap[key],
	})
	if err != nil {
		c.log.Error("GetMsg", err)
		return HandleError(key, dfMsgMap[key], templateData)
	}
	return v
}

// GetParameterErrorMsg 获取参数错误消息
func (c *CommonLang) GetParameterErrorMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, parameterErrorKey, nil, langs...)
}

// GetMetadataFailMsg 获取获取元数据失败消息
func (c *CommonLang) GetMetadataFailMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, getMetadataFailKey, nil, langs...)
}

// GetMetadataConversionFailMsg 获取获取元数据转换失败消息
func (c *CommonLang) GetMetadataConversionFailMsg(ctx context.Context, source string, target string, langs ...string) string {
	return c.getMsg(ctx, getMetadataConversionFailKey, map[string]string{"Source": source, "Target": target}, langs...)
}

// GetMissingMetadataMsg 获取缺少元数据消息
func (c *CommonLang) GetMissingMetadataMsg(ctx context.Context, name string, langs ...string) string {
	return c.getMsg(ctx, missingMetadataKey, map[string]string{"Name": c.userParameterName(ctx, name, langs...)}, langs...)
}

// GetEncryptFailMsg 获取加密失败消息
func (c *CommonLang) GetEncryptFailMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, encryptFailKey, nil, langs...)
}

// GetDecryptFailMsg 获取解密失败消息
func (c *CommonLang) GetDecryptFailMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, decryptFailKey, nil, langs...)
}

// GetParamCanNotEmptyMsg 获取参数不能为空消息
func (c *CommonLang) GetParamCanNotEmptyMsg(ctx context.Context, name string, langs ...string) string {
	return c.getMsg(ctx, paramCanNotEmptyKey, map[string]string{"Name": c.userParameterName(ctx, name, langs...)}, langs...)
}

// userParameterName 把内部单位协议字段转换成当前语言的用户术语，其他参数名保持原样。
func (c *CommonLang) userParameterName(ctx context.Context, name string, langs ...string) string {
	switch strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "")) {
	case "tenantid":
		return c.getMsg(ctx, organizationIDParamKey, nil, langs...)
	case "tenantids", "tenantidin":
		return c.getMsg(ctx, organizationIDsParamKey, nil, langs...)
	default:
		return name
	}
}

// GetNotEditableMsg 获取数据不可编辑消息
func (c *CommonLang) GetNotEditableMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, notEditableKey, nil, langs...)
}

// GetRequestErrorMsg 获取请求错误消息
func (c *CommonLang) GetRequestErrorMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, requestErrorKey, nil, langs...)
}

// GetPleaseLoginMsg 获取请登录消息
func (c *CommonLang) GetPleaseLoginMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, pleaseLoginKey, nil, langs...)
}

// GetNeedReLoginMsg 获取需要重新登录消息
func (c *CommonLang) GetNeedReLoginMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, needReLoginKey, nil, langs...)
}

// GetAccountDisabledMsg 获取账号已禁用消息
func (c *CommonLang) GetAccountDisabledMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, accountDisabledKey, nil, langs...)
}

// GetPleaseChangePasswordMsg 获取请修改密码消息
func (c *CommonLang) GetPleaseChangePasswordMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, pleaseChangePasswordKey, nil, langs...)
}

// GetPasswordErrorMsg 获取密码错误消息
func (c *CommonLang) GetPasswordErrorMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, passwordErrorKey, nil, langs...)
}

// GetNoPermissionMsg 获取没有权限消息
func (c *CommonLang) GetNoPermissionMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, noPermissionKey, nil, langs...)
}

// GetNotConfiguredMsg 获取未配置消息
func (c *CommonLang) GetNotConfiguredMsg(ctx context.Context, name string, langs ...string) string {
	return c.getMsg(ctx, notConfiguredKey, map[string]string{"Name": name}, langs...)
}

// GetCanNotEmptyMsg 获取不能为空消息
func (c *CommonLang) GetCanNotEmptyMsg(ctx context.Context, name string, langs ...string) string {
	return c.getMsg(ctx, canNotEmptyKey, map[string]string{"Name": name}, langs...)
}

// GetExpiredMsg 获取已过期消息
func (c *CommonLang) GetExpiredMsg(ctx context.Context, name string, langs ...string) string {
	return c.getMsg(ctx, expiredKey, map[string]string{"Name": name}, langs...)
}

// GetDataAbnormalMsg 获取数据异常消息
func (c *CommonLang) GetDataAbnormalMsg(ctx context.Context, name string, langs ...string) string {
	return c.getMsg(ctx, dataAbnormalKey, map[string]string{"Name": name}, langs...)
}

// GetDataFoundMsg 获取数据未找到消息
func (c *CommonLang) GetDataNotFoundMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, dataNotFoundKey, nil, langs...)
}

// GetDataDuplicateMsg 获取数据重复消息
func (c *CommonLang) GetDataDuplicateKeyMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, dataDuplicateKey, nil, langs...)
}

// GetDataConstraintMsg 获取数据约束消息
func (c *CommonLang) GetDataConstraintMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, dataConstraintKey, nil, langs...)
}

// GetDataNotLoadedMsg 获取数据未加载消息
func (c *CommonLang) GetDataNotLoadedMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, dataNotLoadedKey, nil, langs...)
}

// GetDataNotSingularMsg 获取数据不是单数消息
func (c *CommonLang) GetDataNotSingularMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, dataNotSingularKey, nil, langs...)
}

// GetDataValidationErrorMsg 获取数据验证错误消息
func (c *CommonLang) GetDataValidationErrorMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, dataValidationErrorKey, nil, langs...)
}

// GetEntValidationFieldDetailMsg 根据 Ent 等 ORM 校验叶子文案返回字段级多语言说明（用于 fault metadata 值）。
// 无法安全映射时返回空字符串，调用方应省略 metadata 中的该字段说明。
func (c *CommonLang) GetEntValidationFieldDetailMsg(ctx context.Context, leaf string, langs ...string) string {
	kind, safe := classifyEntValidationLeaf(leaf)
	switch kind {
	case entValidationDetailMissingField:
		return c.getMsg(ctx, entValidateMissingFieldKey, nil, langs...)
	case entValidationDetailMissingEdge:
		return c.getMsg(ctx, entValidateMissingEdgeKey, nil, langs...)
	case entValidationDetailInvalidQueryField:
		return c.getMsg(ctx, entValidateInvalidQueryFieldKey, nil, langs...)
	case entValidationDetailMaxLen:
		return c.getMsg(ctx, entValidateMaxLenKey, nil, langs...)
	case entValidationDetailMinLen:
		return c.getMsg(ctx, entValidateMinLenKey, nil, langs...)
	case entValidationDetailMaxRuneLen:
		return c.getMsg(ctx, entValidateMaxRuneLenKey, nil, langs...)
	case entValidationDetailMinRuneLen:
		return c.getMsg(ctx, entValidateMinRuneLenKey, nil, langs...)
	case entValidationDetailOutOfRange:
		return c.getMsg(ctx, entValidateOutOfRangeKey, nil, langs...)
	case entValidationDetailFormatMismatch:
		return c.getMsg(ctx, entValidateFormatMismatchKey, nil, langs...)
	case entValidationDetailPassthrough:
		if safe == "" {
			return ""
		}
		return c.getMsg(ctx, entValidateDetailKey, map[string]string{"Detail": safe}, langs...)
	default:
		return ""
	}
}

// GetDataErrorMsg 获取数据错误消息
func (c *CommonLang) GetDataErrorMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, dataErrorKey, nil, langs...)
}

// GetInitSeedFailedMsg 获取种子数据初始化失败消息（Reason 为底层错误说明，供日志与 API 展示）
func (c *CommonLang) GetInitSeedFailedMsg(ctx context.Context, reason string, langs ...string) string {
	return c.getMsg(ctx, initSeedFailedKey, map[string]string{"Reason": reason}, langs...)
}

// GetVersionedResourceBadRequestMsg 返回通用版本资源参数或 payload 契约错误文案。
func (c *CommonLang) GetVersionedResourceBadRequestMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, versionedResourceBadRequestKey, nil, langs...)
}

// GetVersionedResourceNotFoundMsg 返回通用版本资源不存在文案，并保持跨租户场景不泄漏资源信息。
func (c *CommonLang) GetVersionedResourceNotFoundMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, versionedResourceNotFoundKey, nil, langs...)
}

// GetVersionedResourceConflictMsg 返回通用版本资源并发或状态冲突文案。
func (c *CommonLang) GetVersionedResourceConflictMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, versionedResourceConflictKey, nil, langs...)
}

// GetResourceMetadataBadRequestMsg 返回通用实体元数据请求错误文案。
func (c *CommonLang) GetResourceMetadataBadRequestMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, resourceMetadataBadRequestKey, nil, langs...)
}

// GetResourceMetadataNotFoundMsg 返回通用实体元数据所属资源不存在文案。
func (c *CommonLang) GetResourceMetadataNotFoundMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, resourceMetadataNotFoundKey, nil, langs...)
}

// GetResourceMetadataConflictMsg 返回通用实体元数据状态冲突文案。
func (c *CommonLang) GetResourceMetadataConflictMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, resourceMetadataConflictKey, nil, langs...)
}

// GetResourceMetadataLimitExceededMsg 返回通用实体元数据处理限制文案。
func (c *CommonLang) GetResourceMetadataLimitExceededMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, resourceMetadataLimitExceededKey, nil, langs...)
}

// GetResourceMetadataDataAbnormalMsg 返回通用实体元数据存储异常文案。
func (c *CommonLang) GetResourceMetadataDataAbnormalMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, resourceMetadataDataAbnormalKey, nil, langs...)
}

// GetResourceMetadataUnavailableMsg 返回通用实体元数据服务不可用文案。
func (c *CommonLang) GetResourceMetadataUnavailableMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, resourceMetadataUnavailableKey, nil, langs...)
}

// GetResourceMetadataDeadlineExceededMsg 返回通用实体元数据请求超时文案。
func (c *CommonLang) GetResourceMetadataDeadlineExceededMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, resourceMetadataDeadlineExceededKey, nil, langs...)
}

// GetResourceMetadataCanceledMsg 返回通用实体元数据请求取消文案。
func (c *CommonLang) GetResourceMetadataCanceledMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, resourceMetadataCanceledKey, nil, langs...)
}

// HandleResourceMetadataError 保留实体元数据错误码，以公共多语言文案替换底层错误正文。
func (c *CommonLang) HandleResourceMetadataError(ctx context.Context, err error, langs ...string) error {
	if err == nil {
		return nil
	}
	if stderrors.Is(err, context.DeadlineExceeded) {
		return status.Error(codes.DeadlineExceeded, c.GetResourceMetadataDeadlineExceededMsg(ctx, langs...))
	}
	if stderrors.Is(err, context.Canceled) {
		return status.Error(codes.Canceled, c.GetResourceMetadataCanceledMsg(ctx, langs...))
	}
	var kratosError *kerrors.Error
	if stderrors.As(err, &kratosError) {
		switch kratosError.Code {
		case 400:
			return status.Error(codes.InvalidArgument, c.GetResourceMetadataBadRequestMsg(ctx, langs...))
		case 404:
			return status.Error(codes.NotFound, c.GetResourceMetadataNotFoundMsg(ctx, langs...))
		case 409:
			return status.Error(codes.Aborted, c.GetResourceMetadataConflictMsg(ctx, langs...))
		case 429:
			return status.Error(codes.ResourceExhausted, c.GetResourceMetadataLimitExceededMsg(ctx, langs...))
		default:
			return status.Error(codes.Internal, c.GetResourceMetadataUnavailableMsg(ctx, langs...))
		}
	}
	code := status.Code(err)
	switch code {
	case codes.InvalidArgument:
		return status.Error(code, c.GetResourceMetadataBadRequestMsg(ctx, langs...))
	case codes.NotFound:
		return status.Error(code, c.GetResourceMetadataNotFoundMsg(ctx, langs...))
	case codes.AlreadyExists, codes.Aborted, codes.FailedPrecondition:
		return status.Error(code, c.GetResourceMetadataConflictMsg(ctx, langs...))
	case codes.ResourceExhausted:
		return status.Error(code, c.GetResourceMetadataLimitExceededMsg(ctx, langs...))
	case codes.DataLoss:
		return status.Error(code, c.GetResourceMetadataDataAbnormalMsg(ctx, langs...))
	case codes.DeadlineExceeded:
		return status.Error(code, c.GetResourceMetadataDeadlineExceededMsg(ctx, langs...))
	case codes.Canceled:
		return status.Error(code, c.GetResourceMetadataCanceledMsg(ctx, langs...))
	case codes.Unavailable:
		return status.Error(code, c.GetResourceMetadataUnavailableMsg(ctx, langs...))
	default:
		return status.Error(codes.Internal, c.GetResourceMetadataUnavailableMsg(ctx, langs...))
	}
}

// GetTreeParentCanNotSelfMsg 获取树父节点不能是自己的消息
func (c *CommonLang) GetTreeParentCanNotSelfMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, treeParentCanNotSelfKey, nil, langs...)
}

// GetTreeParentCanNotDescendantMsg 获取树父节点不能是子节点的消息
func (c *CommonLang) GetTreeParentCanNotDescendantMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, treeParentCanNotDescendantKey, nil, langs...)
}

// GetTreeSortedSiblingMustIncludeSelfMsg 获取树排序列表必须包含当前节点的消息
func (c *CommonLang) GetTreeSortedSiblingMustIncludeSelfMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, treeSortedSiblingMustIncludeSelfKey, nil, langs...)
}

// GetTreeSiblingNotFoundMsg 获取树兄弟节点不存在消息
func (c *CommonLang) GetTreeSiblingNotFoundMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, treeSiblingNotFoundKey, nil, langs...)
}

// GetTreeSiblingParentMismatchMsg 获取树兄弟节点父级不一致消息
func (c *CommonLang) GetTreeSiblingParentMismatchMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, treeSiblingParentMismatchKey, nil, langs...)
}

// GetTreeSiblingIDEmptyMsg 获取树兄弟节点 ID 不能为空消息
func (c *CommonLang) GetTreeSiblingIDEmptyMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, treeSiblingIDEmptyKey, nil, langs...)
}

// GetTreeSiblingIDDuplicateMsg 获取树兄弟节点 ID 不能重复消息
func (c *CommonLang) GetTreeSiblingIDDuplicateMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, treeSiblingIDDuplicateKey, nil, langs...)
}

// GetCacheNotFoundMsg 获取缓存未找到消息
func (c *CommonLang) GetCacheNotFoundMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, cacheNotFoundKey, nil, langs...)
}

// GetCachePreMatchGetFailMsg 获取缓存前置匹配获取失败消息
func (c *CommonLang) GetCachePreMatchGetFailMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, cachePreMatchGetFailKey, nil, langs...)
}

// GetCacheSetFailMsg 获取缓存设置失败消息
func (c *CommonLang) GetCacheSetFailMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, cacheSetFailKey, nil, langs...)
}

// GetCacheMSetFailMsg 获取缓存批量设置失败消息
func (c *CommonLang) GetCacheMSetFailMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, cacheMSetFailKey, nil, langs...)
}

// GetCacheDelFailMsg 获取缓存删除失败消息
func (c *CommonLang) GetCacheDelFailMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, cacheDelFailKey, nil, langs...)
}

// GetCachePreMatchDelFailMsg 获取缓存前置匹配删除失败消息
func (c *CommonLang) GetCachePreMatchDelFailMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, cachePreMatchDelFailKey, nil, langs...)
}

// GetCacheFlushFailMsg 获取缓存清空失败消息
func (c *CommonLang) GetCacheFlushFailMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, cacheFlushFailKey, nil, langs...)
}

// GetCacheMGetFailMsg 获取缓存批量获取失败消息
func (c *CommonLang) GetCacheMGetFailMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, cacheMGetFailKey, nil, langs...)
}

// GetCacheMDelFailMsg 获取缓存批量删除失败消息
func (c *CommonLang) GetCacheMDelFailMsg(ctx context.Context, langs ...string) string {
	return c.getMsg(ctx, cacheMDelFailKey, nil, langs...)
}
