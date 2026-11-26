package polling

import (
	"context"
	"errors"
	"time"
)

type groupState struct {
	key       string
	excluded  []string
	inflight  int
	exhausted bool
	pending   []Candidate
	noMore    bool
}
type completion struct {
	group                  int
	attempted              bool
	callbackErr, recordErr error
}

// roundState 仅由派发协程持有可变状态，回调通过完成通道交还槽位。
type roundState struct {
	groups         []groupState
	done           chan completion
	seen           map[string]bool
	inflight, next int
	stopped        bool
	result         error
	report         *Report
}

// dispatch 轮转分配槽位，满组暂让其他组；所有槽位占用时等待完成事件。
func (r *Runner) dispatch(jobCtx, dispatchCtx context.Context, lease Lease, scan Scan, keys []string, report *Report) error {
	// 按本轮组数摊分预取额度，所有组的候选缓冲合计不超过 MaxItems。
	pageSize := min(r.Policy.CandidateBatchSize, max(1, r.Policy.MaxItems/max(1, len(keys))))
	state := &roundState{groups: make([]groupState, len(keys)), done: make(chan completion, r.Policy.MaxConcurrency), seen: map[string]bool{}, report: report}
	for i, key := range keys {
		state.groups[i].key = key
	}
	for !state.stopped && report.Dispatched < r.Policy.MaxItems && report.GroupVisits < r.Policy.MaxGroupVisits {
		state.drainReady()
		if state.stopped || dispatchCtx.Err() != nil || leaseLost(lease) {
			break
		}
		index := availableGroup(state.groups, state.next, r.Policy.MaxGroupConcurrency)
		if index < 0 || state.inflight >= r.Policy.MaxConcurrency {
			if !state.wait(dispatchCtx) {
				break
			}
			continue
		}
		group := &state.groups[index]
		state.next = (index + 1) % len(state.groups)
		report.GroupVisits++
		candidate, err := r.nextCandidate(dispatchCtx, lease, scan, group, min(pageSize, r.Policy.MaxItems-report.Dispatched), report)
		if err != nil {
			state.result = errors.Join(state.result, err)
			break
		}
		if candidate == nil {
			continue
		}
		if state.seen[candidate.RecordKey] {
			state.result = errors.Join(state.result, NewError("duplicate"))
			break
		}
		// 查询期间也可能记时失败或失租；直接检查租约，避免依赖异步 context 取消的调度时机。
		state.drainReady()
		if state.stopped || dispatchCtx.Err() != nil || leaseLost(lease) {
			break
		}
		state.seen[candidate.RecordKey] = true
		group.excluded = append(group.excluded, candidate.RecordKey)
		group.inflight++
		state.inflight++
		report.Dispatched++
		// 每个协程只运行一次回调并报告完成；结束派发后仍等待全部已启动调用退出。
		go func() { state.done <- r.attempt(jobCtx, index, *candidate) }()
	}
	for state.inflight > 0 {
		state.collect(<-state.done)
	}
	return state.result
}

// nextCandidate 只在实际轮转访问时写游标；批量读取候选不提前消费其他组的机会。
func (r *Runner) nextCandidate(ctx context.Context, lease Lease, scan Scan, group *groupState, limit int, report *Report) (*Candidate, error) {
	if err := lease.SaveCursor(ctx, group.key); err != nil {
		return nil, err
	}
	if len(group.pending) == 0 && !group.noMore {
		report.CandidateQueries++
		rows, err := r.Source.ListCandidates(ctx, scan, group.key, group.excluded, limit)
		if err != nil {
			return nil, err
		}
		if len(rows) > limit {
			return nil, NewError("candidate_limit")
		}
		// 一页必须整体合法且稳定有序，避免先执行部分候选后才发现存储契约错误。
		pageKeys := make(map[string]bool, len(rows))
		for i, c := range rows {
			if err := validateCandidate(c, *group, scan); err != nil {
				return nil, err
			}
			if pageKeys[c.RecordKey] {
				return nil, NewError("duplicate")
			}
			pageKeys[c.RecordKey] = true
			if i > 0 {
				before, after := candidateWaitingSince(rows[i-1]), candidateWaitingSince(c)
				if before.After(after) || (before.Equal(after) && rows[i-1].RecordKey >= c.RecordKey) {
					return nil, NewError("candidate_order")
				}
			}
		}
		report.CandidatesFetched += len(rows)
		// 复制切片所有权，弹出缓存项时不能改写 Source 持有的候选投影。
		group.pending = append([]Candidate(nil), rows...)
		group.noMore = len(rows) < limit
	}
	if len(group.pending) == 0 {
		group.exhausted = true
		return nil, nil
	}
	candidate := group.pending[0]
	// 清除已弹出引用；余下投影只在本轮有效，下一轮从数据库重新查询。
	group.pending[0] = Candidate{}
	group.pending = group.pending[1:]
	return &candidate, nil
}

// candidateWaitingSince 使用统一等待起点校验预取页顺序，未检查记录不永久插队。
func candidateWaitingSince(c Candidate) time.Time {
	if c.LastCheckedAt != nil {
		return *c.LastCheckedAt
	}
	return c.ReadyAt
}

// collect 在派发协程中归并事实；检查事实写入失败后禁止继续派发。
func (s *roundState) collect(c completion) {
	s.groups[c.group].inflight--
	s.inflight--
	if c.attempted {
		s.report.Attempted++
	} else {
		s.report.Skipped++
	}
	if c.callbackErr != nil {
		s.report.CallbackErrors++
	}
	if c.recordErr != nil {
		s.report.RecordErrors++
		s.stopped = true
	}
	s.result = errors.Join(s.result, c.callbackErr, c.recordErr)
}

// drainReady 优先处理已完成通知，防止已知记时失败后仍继续启动业务。
func (s *roundState) drainReady() {
	for s.inflight > 0 {
		select {
		case c := <-s.done:
			s.collect(c)
		default:
			return
		}
	}
}

// wait 等待槽位或派发截止，避免无容量时空转占用 CPU。
func (s *roundState) wait(ctx context.Context) bool {
	if s.inflight == 0 {
		return false
	}
	start := time.Now()
	defer func() { s.report.WaitDuration += time.Since(start) }()
	select {
	case c := <-s.done:
		s.collect(c)
		return !s.stopped
	case <-ctx.Done():
		return false
	}
}
