package audit

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	public "github.com/zhixingdx/public/api/base/audit/v1"
	"google.golang.org/protobuf/proto"
)

// ConfirmSoftDelete 在源事务确定提交后有界收尾，确认失败不能改变业务提交结果。
func (r *Runtime) ConfirmSoftDelete(parent context.Context, prepared []PreparedRecord) {
	if r == nil || len(prepared) == 0 {
		return
	}
	// 源事务已提交后确认独立于请求取消；租户由 PreparedRecord 冻结，避免入站元数据重复透传。
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	type groupKey struct{ tenant, service, deployment string }
	groups := map[groupKey][]PreparedRecord{}
	order := []groupKey{}
	for _, item := range prepared {
		key := groupKey{item.TenantID, item.SourceService, item.SourceDeploymentID}
		if _, ok := groups[key]; !ok {
			order = append(order, key)
		}
		groups[key] = append(groups[key], item)
	}
	for _, key := range order {
		items := groups[key]
		for _, batch := range confirmBatches(items) {
			client, callCtx, err := r.producer.tenantClient(ctx, batch[0].TenantID)
			if err != nil {
				r.Diagnose("confirm")
				continue
			}
			request := confirmRequest(batch)
			confirmed := false
			for attempt := 0; attempt < 2 && ctx.Err() == nil; attempt++ {
				reply, err := client.ConfirmSoftDelete(callCtx, request)
				if err == nil && validConfirmReceipts(batch, reply.GetData()) {
					confirmed = true
					break
				}
			}
			if !confirmed {
				r.Diagnose("confirm")
			}
		}
	}
}

const (
	confirmMaxRecordIDs    = 100
	confirmMaxRequestBytes = 16 * 1024
)

// confirmBatches 按中央确认 RPC 的条数和编码大小上限拆分同一来源的已提交快照。
func confirmBatches(items []PreparedRecord) [][]PreparedRecord {
	if len(items) == 0 {
		return nil
	}
	batches := make([][]PreparedRecord, 0, (len(items)+confirmMaxRecordIDs-1)/confirmMaxRecordIDs)
	for offset := 0; offset < len(items); {
		end := offset
		for end < len(items) && end-offset < confirmMaxRecordIDs {
			candidate := items[offset : end+1]
			if end > offset && proto.Size(confirmRequest(candidate)) > confirmMaxRequestBytes {
				break
			}
			end++
		}
		batches = append(batches, items[offset:end])
		offset = end
	}
	return batches
}

// confirmRequest 从同一租户和来源的预备快照构建受中央协议约束的确认请求。
func confirmRequest(items []PreparedRecord) *public.ConfirmSoftDeleteRequest {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.RecordID
	}
	return &public.ConfirmSoftDeleteRequest{SourceService: items[0].SourceService, SourceDeploymentID: items[0].SourceDeploymentID, RecordIDs: ids}
}

// validConfirmReceipts 按中央确定性确认身份逐项核对，不把无关成功回执当确认。
func validConfirmReceipts(prepared []PreparedRecord, receipts []*Receipt) bool {
	if len(prepared) != len(receipts) {
		return false
	}
	pending := map[string]string{}
	for _, item := range prepared {
		digest := sha256.Sum256([]byte(item.RecordID))
		pending[fmt.Sprintf("soft-delete-confirm:%x", digest)] = item.SourceDeploymentID
	}
	for _, receipt := range receipts {
		if receipt == nil || pending[receipt.EventID] != receipt.SourceDeploymentID || receipt.StorageTarget != "database" || receipt.RecordID == "" || (receipt.Status != "stored" && receipt.Status != "duplicate") {
			return false
		}
		delete(pending, receipt.EventID)
	}
	return len(pending) == 0
}
