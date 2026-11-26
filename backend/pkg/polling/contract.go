// Package polling 为业务回调分配有界、公平的推进机会，不拥有业务状态。
package polling

import (
	"context"
	"fmt"
	"strings"
	"time"
)

var ErrAlreadyRunning = NewError("running")

// JobPolicy 描述单轮资源预算；扫描节拍由调用方维护。
type JobPolicy struct {
	JobKey                                                        string
	MinInterval                                                   time.Duration
	MaxItems, MaxConcurrency, MaxGroupConcurrency, MaxGroupVisits int
	CandidateBatchSize                                            int
	DispatchBudget, AttemptTimeout, RecordTimeout                 time.Duration
}

// Validate 在分配内存和访问基础设施前拒绝越界配置，避免单轮资源随配置无界增长。
func (p JobPolicy) Validate() error {
	if strings.TrimSpace(p.JobKey) == "" || len(p.JobKey) > 256 {
		return NewError("job_key")
	}
	for _, bound := range []struct {
		name         string
		value, upper int
	}{
		{"MaxItems", p.MaxItems, 1000}, {"MaxConcurrency", p.MaxConcurrency, 100},
		{"MaxGroupConcurrency", p.MaxGroupConcurrency, 100}, {"MaxGroupVisits", p.MaxGroupVisits, 10000},
		{"CandidateBatchSize", p.CandidateBatchSize, 100},
	} {
		if bound.value < 1 || bound.value > bound.upper {
			return &Error{Reason: "capacity", Field: bound.name, Limit: fmt.Sprint(bound.upper)}
		}
	}
	if p.MaxGroupConcurrency > p.MaxConcurrency || p.MaxConcurrency > p.MaxItems || p.MaxGroupVisits < p.MaxItems {
		return NewError("concurrency")
	}
	if p.MinInterval < 0 || p.MinInterval > 24*time.Hour {
		return NewError("interval")
	}
	for _, bound := range []struct {
		name         string
		value, upper time.Duration
	}{
		{"DispatchBudget", p.DispatchBudget, time.Hour}, {"AttemptTimeout", p.AttemptTimeout, time.Hour}, {"RecordTimeout", p.RecordTimeout, time.Minute},
	} {
		if bound.value < time.Millisecond || bound.value > bound.upper {
			return &Error{Reason: "duration", Field: bound.name, Limit: bound.upper.String()}
		}
	}
	return nil
}

// Candidate 只携带调度投影；业务对象由回调重新读取。
type Candidate struct {
	RecordKey, GroupKey string
	LastCheckedAt       *time.Time
	ReadyAt             time.Time
}

// Scan 固定本轮时间及最小间隔，所有发现和候选查询共用此条件。
type Scan struct {
	Now         time.Time
	MinInterval time.Duration
}

// Source 提供有界、稳定排序的存储查询及独立检查事实写入。
// DiscoverGroups 按键升序返回严格大于 after 的分组；空 after 从头开始。
// ListCandidates 按 COALESCE(lastCheckedAt, readyAt)、记录键升序，排除本轮已经选择的键。
type Source interface {
	DiscoverGroups(context.Context, Scan, string, int) ([]string, error)
	ListCandidates(context.Context, Scan, string, []string, int) ([]Candidate, error)
	RecordAttempt(context.Context, Candidate, time.Time) error
}

// Advance 复核资格并调用已有业务；未实际尝试返回 false。
type Advance func(context.Context, Candidate) (attempted bool, err error)

// Lease 维护同一 Job 的执行权及持久游标；Lost 关闭后不得继续派发。
type Lease interface {
	Lost() <-chan struct{}
	Err() error
	Cursor(context.Context) (string, error)
	SaveCursor(context.Context, string) error
	Release(context.Context) error
}

// Coordinator 由基础设施适配器实现，不把驱动依赖引入核心。
type Coordinator interface {
	Acquire(context.Context, string) (Lease, error)
}

// Report 汇总有界调度事实，不使用分组键作指标标签。
type Report struct {
	DiscoveredGroups, GroupVisits, Dispatched, Attempted, Skipped, CallbackErrors, RecordErrors int
	CandidateQueries, CandidatesFetched                                                         int
	WaitDuration                                                                                time.Duration
	Canceled                                                                                    bool
}

// Runner 每次调用只执行一轮，可复用于不同业务和调度器。
type Runner struct {
	Policy      JobPolicy
	Source      Source
	Advance     Advance
	Coordinator Coordinator
	Now         func() time.Time
}

// now 返回可注入的事实时钟；超时仍由 context 的单调时钟保证。
func (r *Runner) now() time.Time {
	if r.Now != nil {
		return r.Now()
	}
	return time.Now()
}
