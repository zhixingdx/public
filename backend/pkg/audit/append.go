// Package audit 为业务服务提供审计信封与受控批量出站封装。
package audit

import (
	"context"
	"errors"

	pb "github.com/zhixingdx/oipivot-api/base/audit/v1"

	"google.golang.org/protobuf/proto"
)

// Record 复用公开协议的不可变审计信封。
type Record = pb.RecordEnvelope

// Receipt 保留存储目标及接收状态，云 accepted 仅表示 Collector 本跳接收。
type Receipt = pb.RecordReceipt

var (
	// ErrInvalidRecords 表示本地批次缺少身份或同一事件出现冲突，仅供内部诊断。
	ErrInvalidRecords = errors.New("invalid audit records")
	// ErrInvalidReceipt 表示远端回执与提交命令不一致，不能统计为成功。
	ErrInvalidReceipt = errors.New("invalid audit receipt")
)

// BatchAppend 按精确策略键分组后分块，避免数据库与云路由混批及超过服务端 100 条上限。
// 返回的每条回执已校验，发生错误时保留此前已确认的回执；不自动重试或生成新事件 ID。
func (c *ProducerClient) BatchAppend(ctx context.Context, tenant string, records []*Record) ([]*Receipt, error) {
	client, ctx, err := c.tenantClient(ctx, tenant)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, ErrInvalidRecords
	}
	type policyKey struct {
		service                           string
		action                            string
		version                           int32
		category, captureMode, deleteMode string
	}
	groups := make(map[policyKey][]*Record)
	order := make([]policyKey, 0)
	seen := make(map[[2]string]*Record)
	for _, record := range records {
		if record == nil || record.EventID == "" || record.SourceDeploymentID == "" || record.BizCall.GetService() == "" || record.BizCall.GetAction() == "" || record.PolicyVersion <= 0 {
			return nil, ErrInvalidRecords
		}
		id := [2]string{record.SourceDeploymentID, record.EventID}
		if previous, ok := seen[id]; ok {
			if !proto.Equal(previous, record) {
				return nil, ErrInvalidRecords
			}
			continue
		}
		seen[id] = record
		key := policyKey{record.BizCall.Service, record.BizCall.Action, record.PolicyVersion, record.Category, record.CaptureMode, record.DeleteMode}
		if _, ok := groups[key]; !ok {
			order = append(order, key)
		}
		groups[key] = append(groups[key], record)
	}
	var receipts []*Receipt
	for _, key := range order {
		group := groups[key]
		for offset := 0; offset < len(group); offset += 100 {
			items := group[offset:min(offset+100, len(group))]
			reply, err := client.AppendRecords(ctx, &pb.AppendRecordsRequest{Items: items})
			if err != nil {
				return receipts, err
			}
			batchReceipts := reply.GetData()
			if !validReceipts(items, batchReceipts) {
				return receipts, ErrInvalidReceipt
			}
			receipts = append(receipts, batchReceipts...)
		}
	}
	return receipts, nil
}

// validReceipts 拒绝缺失、重复、错位身份及与存储目标不相符的回执，防止假成功。
func validReceipts(records []*Record, receipts []*Receipt) bool {
	if len(records) != len(receipts) {
		return false
	}
	pending := make(map[[2]string]*Record, len(records))
	for _, record := range records {
		pending[[2]string{record.SourceDeploymentID, record.EventID}] = record
	}
	for _, receipt := range receipts {
		if receipt == nil {
			return false
		}
		key := [2]string{receipt.SourceDeploymentID, receipt.EventID}
		r := pending[key]
		if r == nil {
			return false
		}
		delete(pending, key)
		if r.Category == "data_retention" && r.CaptureMode == "resource_snapshot" && r.DeleteMode == "soft" {
			if receipt.StorageTarget != "database" || (receipt.Status != "stored" && receipt.Status != "duplicate") || receipt.RecordID == "" {
				return false
			}
		} else if receipt.StorageTarget != "cloud" || receipt.Status != "accepted" || receipt.RecordID != "" {
			return false
		}
	}
	return true
}
