package lang

import (
	"context"
	"errors"

	kerrors "github.com/go-kratos/kratos/v3/errors"
	"github.com/zhixingdx/oipivot-api/fault"
	"github.com/zhixingdx/oipivot/backend/pkg/polling"
)

// HandlePollingError 在服务边界本地化调度错误；保留 cause，使互斥、取消和业务错误仍可识别。
func (c *CommonLang) HandlePollingError(ctx context.Context, err error, langs ...string) error {
	if err == nil {
		return nil
	}
	// 逐项处理合并错误，避免只翻译首个错误或丢失其余失败原因。
	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		causes := joined.Unwrap()
		localized := make([]error, 0, len(causes))
		for _, cause := range causes {
			localized = append(localized, c.HandlePollingError(ctx, cause, langs...))
		}
		return errors.Join(localized...)
	}
	var existing *kerrors.Error
	if errors.As(err, &existing) {
		return err
	}
	reason := "unavailable"
	values := map[string]string{}
	var detail *polling.Error
	if errors.As(err, &detail) {
		reason = detail.Reason
		values["Field"] = detail.Field
		values["Limit"] = detail.Limit
	}
	if errors.Is(err, context.Canceled) {
		reason = "canceled"
	} else if errors.Is(err, context.DeadlineExceeded) {
		reason = "deadline"
	}
	key := MsgKey("polling_" + reason)
	// 未知原因不把内部诊断键或驱动正文直接作为用户文案。
	if _, ok := dfMsgMap[key]; !ok {
		key = MsgKey("polling_unavailable")
	}
	message := c.getMsg(ctx, key, values, langs...)
	switch reason {
	case "running":
		return fault.ErrorConflict("%s", message).WithCause(err)
	case "deadline", "canceled":
		return fault.ErrorRequestTimeout("%s", message).WithCause(err)
	case "job_key", "capacity", "concurrency", "interval", "duration", "eligibility", "scan", "request", "fact", "coordinator":
		return fault.ErrorBadRequest("%s", message).WithCause(err)
	default:
		return fault.ErrorInternalServerError("%s", message).WithCause(err)
	}
}
