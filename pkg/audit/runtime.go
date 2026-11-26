package audit

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zhixingdx/oipivot-api/apicommon/v1"
	"github.com/zhixingdx/oipivot-api/config"
	log "github.com/zhixingdx/oipivot/backend/pkg/logger/kratoslog"
	"github.com/zhixingdx/oipivot/backend/pkg/middleware/meta"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ErrUnavailable 表示可靠存档依赖不可用，出站原始错误不进入用户消息。
var ErrUnavailable = errors.New("audit archive unavailable")

// DataPolicy 是 Schema 编译出的固定动作与投影配置。
type DataPolicy struct {
	Service             string
	ActionPrefix        string
	PolicyVersion       int32
	UpdatePolicyVersion int32
	Fields              []Field
}

// TenantRecord 冻结真实对象归属；出站租户不从后台 worker 上下文推断。
type TenantRecord struct {
	TenantID string
	Record   *Record
}

// PreparedRecord 保存中央已确认的原存档身份，供源事务提交后确认。
type PreparedRecord struct {
	TenantID           string
	SourceService      string
	SourceDeploymentID string
	RecordID           string
}

// Runtime 统一同步存档和有界尽力发送，由接入服务的生命周期管理。
type Runtime struct {
	producer *ProducerClient
	logger   *log.Helper
	mu       sync.Mutex
	queue    chan TenantRecord
	bytes    int
	closed   bool
	wg       sync.WaitGroup
	stop     context.CancelFunc
	dropped  atomic.Uint64
}

// NewRuntime 仅在显式接入服务创建发送 worker，缺少稳定来源时拒绝启动。
func NewRuntime(c *config.Config, service string, logger *log.Helper) (*Runtime, error) {
	producer, err := NewProducerClient(c, service)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	r := &Runtime{producer: producer, logger: logger, queue: make(chan TenantRecord, 1000), stop: cancel}
	for i := 0; i < 2; i++ {
		r.wg.Add(1)
		go r.sendLoop(ctx)
	}
	return r, nil
}

// UnscopedTenant 返回由部署显式声明的无租户列数据归属。未声明时此类事件无法归属。
func (r *Runtime) UnscopedTenant() string {
	if r == nil {
		return ""
	}
	return r.producer.unscopedTenant
}

// Record 根据源服务上下文冻结事件，不读取或解密业务 DTO。
func (r *Runtime) Record(ctx context.Context, policy DataPolicy, id, operation string) *Record {
	if r == nil {
		return nil
	}
	record := NewMetadataRecord(ctx, r.producer.SourceDeploymentID, &apicommon.BizCall{Service: policy.Service, Action: policy.ActionPrefix + "." + operation, BizId: id}, "data_change")
	record.PolicyVersion = policy.PolicyVersion
	if record.PolicyVersion <= 0 {
		record.PolicyVersion = 1
	}
	if operation == "update" && policy.UpdatePolicyVersion > 0 {
		record.PolicyVersion = policy.UpdatePolicyVersion
	}
	if operation == "create" {
		record.Category = "authorization"
	}
	if operation == "soft_delete.prepare" {
		record.Category = "data_retention"
		record.CaptureMode = "resource_snapshot"
		record.DeleteMode = "soft"
	}
	if operation == "hard_delete" {
		record.DeleteMode = "hard"
	}
	if operation == "update" {
		record.CaptureMode = "change_set"
	}
	record.Outcome = "success"
	if actor, err := meta.GetStaffID(ctx); err == nil && actor != "" {
		record.ActorID = actor
		record.ActorType = "staff"
	} else if actor, err := meta.GetUserID(ctx); err == nil && actor != "" {
		record.ActorID = actor
		record.ActorType = "user"
	} else if actor, err := meta.GetStudentID(ctx); err == nil && actor != "" {
		// 学生入口拥有独立可信身份，不能将其写入误记为后台服务操作；保留已有员工与用户优先级。
		record.ActorID = actor
		record.ActorType = "student"
	} else {
		record.ActorID = policy.Service
		record.ActorType = "service"
	}
	SetSourceMetadata(record, r.producer.SourceEnvironment, r.producer.ServiceVersion)
	return record
}

// RecordBytes 在入队和 RPC 前校验实际信封预算；规范化尺寸以单条硬上限保守约束。
func RecordBytes(record *Record) (int, error) {
	if record == nil {
		return 0, ErrProjection
	}
	raw, err := protojson.Marshal(record)
	if err != nil || len(raw) > 96*1024 {
		return 0, ErrProjection
	}
	if record.CaptureMode == "resource_snapshot" {
		if _, err := CanonicalValue([]byte(record.Snapshot)); err != nil {
			return 0, err
		}
	}
	return len(raw), nil
}

// TryAppend 复制不可变事件且不等待队列空位，容量或编码失败只作安全诊断。
func (r *Runtime) TryAppend(event TenantRecord) {
	if r == nil {
		return
	}
	size, err := RecordBytes(event.Record)
	if err != nil || event.TenantID == "" {
		r.Diagnose("projection")
		return
	}
	event.Record = proto.Clone(event.Record).(*Record)
	stage := ""
	r.mu.Lock()
	switch {
	case r.closed:
		stage = "closed"
	case r.bytes+size > 16*1024*1024:
		stage = "queue_bytes"
	default:
		select {
		case r.queue <- event:
			r.bytes += size
		default:
			stage = "queue_full"
		}
	}
	r.mu.Unlock()
	// 日志沿用统一安全诊断，放在队列锁外，避免日志写入拖住 worker 或 Close。
	if stage != "" {
		r.Diagnose(stage)
	}
}

// Diagnose 仅输出安全阶段分类，避免正文和底层 SQL/RPC 错误进入诊断。
func (r *Runtime) Diagnose(stage string) {
	if r == nil {
		return
	}
	r.dropped.Add(1)
	if r.logger != nil {
		r.logger.Warnw("audit_stage", stage)
	}
}

// Dropped 返回本进程普通事件丢弃及发送失败计数。
func (r *Runtime) Dropped() uint64 {
	if r == nil {
		return 0
	}
	return r.dropped.Load()
}

// sendLoop 在固定 worker 内发送事件，网络等待不会阻塞业务调用。
func (r *Runtime) sendLoop(parent context.Context) {
	defer r.wg.Done()
	for {
		select {
		case <-parent.Done():
			return
		case event, ok := <-r.queue:
			if !ok {
				return
			}
			size, _ := RecordBytes(event.Record)
			r.mu.Lock()
			r.bytes -= size
			r.mu.Unlock()
			ctx, cancel := context.WithTimeout(parent, 4*time.Second)
			_, err := r.producer.BatchAppend(ctx, event.TenantID, []*Record{event.Record})
			cancel()
			if err != nil {
				r.Diagnose("append")
			}
		}
	}
}

// Close 最多等待五秒排空，之后取消发送，且不会向已关闭通道入队。
func (r *Runtime) Close() {
	if r == nil {
		return
	}
	r.mu.Lock()
	if r.closed {
		r.mu.Unlock()
		return
	}
	r.closed = true
	close(r.queue)
	r.mu.Unlock()
	done := make(chan struct{})
	go func() { r.wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
	}
	r.stop()
	r.producer.Close()
	r.wg.Wait()
}

// PrepareSoftDelete 在整个远端阶段共享四秒预算，只接受匹配的数据库持久化回执。
func (r *Runtime) PrepareSoftDelete(parent context.Context, events []TenantRecord) ([]PreparedRecord, error) {
	if r == nil || len(events) == 0 {
		return nil, ErrUnavailable
	}
	deadline := time.Now().Add(4 * time.Second)
	if callerDeadline, ok := parent.Deadline(); ok && callerDeadline.Before(deadline) {
		deadline = callerDeadline
	}
	// 存档租户已由 TenantRecord 冻结；仅继承取消与截止时间，避免入站元数据再次透传造成重复租户头。
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	stopCancel := context.AfterFunc(parent, cancel)
	defer func() {
		stopCancel()
		cancel()
	}()
	type groupKey struct{ tenant, action string }
	groups := map[groupKey][]TenantRecord{}
	order := []groupKey{}
	for _, event := range events {
		if event.TenantID == "" || event.Record == nil || event.Record.CaptureMode != "resource_snapshot" || event.Record.DeleteMode != "soft" {
			return nil, ErrProjection
		}
		if _, err := RecordBytes(event.Record); err != nil {
			return nil, err
		}
		if r.producer.route(event.TenantID) == nil {
			return nil, ErrUnavailable
		}
		key := groupKey{event.TenantID, event.Record.BizCall.GetAction()}
		if _, ok := groups[key]; !ok {
			order = append(order, key)
		}
		groups[key] = append(groups[key], event)
	}
	prepared := []PreparedRecord{}
	for _, key := range order {
		group := groups[key]
		// 二十个单条96KiB上界同时约束原始和中央规范化总尺寸，不复制中央内容校验。
		for offset := 0; offset < len(group); offset += 20 {
			batch := group[offset:min(offset+20, len(group))]
			records := make([]*Record, len(batch))
			for i, item := range batch {
				records[i] = item.Record
			}
			receipts, err := r.producer.BatchAppend(ctx, batch[0].TenantID, records)
			if err != nil || !validReceipts(records, receipts) {
				return nil, ErrUnavailable
			}
			byEvent := map[string]*Record{}
			for _, record := range records {
				byEvent[record.EventID] = record
			}
			for _, receipt := range receipts {
				if receipt.StorageTarget != "database" || receipt.RecordID == "" {
					return nil, ErrInvalidReceipt
				}
				record := byEvent[receipt.EventID]
				prepared = append(prepared, PreparedRecord{TenantID: batch[0].TenantID, SourceService: record.BizCall.Service, SourceDeploymentID: record.SourceDeploymentID, RecordID: receipt.RecordID})
			}
		}
	}
	return prepared, nil
}

// AppendRequired 复用生产者的已校验回执；敏感业务在确认前不能向调用方释放结果。
func (r *Runtime) AppendRequired(ctx context.Context, tenant string, records []*Record) error {
	if r == nil || r.producer == nil || len(records) == 0 {
		return ErrUnavailable
	}
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()
	_, err := r.producer.BatchAppend(ctx, tenant, records)
	return err
}
