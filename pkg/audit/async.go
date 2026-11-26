// Package audit 异步上报（不阻塞业务主路径）
package audit

import (
	"context"
	"time"

	"github.com/zhixingdx/public/pkg/logger/kratoslog"
)

// TryAppendAsync 异步上报审计事件，失败仅记日志。
func TryAppendAsync(logger *log.Helper, c *ProducerClient, tenantID string, events []*Record) {
	if c == nil || len(events) == 0 {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()
		_, err := c.BatchAppend(ctx, tenantID, events)
		if err != nil && logger != nil {
			logger.Warnf("审计事件上报失败 tenantID=%s error_type=%T", tenantID, err)
		}
	}()
}
