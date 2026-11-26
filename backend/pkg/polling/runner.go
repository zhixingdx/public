package polling

import (
	"context"
	"errors"
)

// Run 获取执行权，运行有界轮次，并在所有已启动回调收口后释放租约。
func (r *Runner) Run(ctx context.Context) (report Report, err error) {
	if err = r.Policy.Validate(); err != nil {
		return
	}
	if r.Source == nil || r.Advance == nil || r.Coordinator == nil {
		return report, NewError("dependencies")
	}
	lease, err := r.Coordinator.Acquire(ctx, r.Policy.JobKey)
	if err != nil {
		return report, err
	}
	jobCtx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	// 监听协程随本轮结束或失租退出；回调共享取消，不能脱离 Job 生命周期。
	go func() {
		select {
		case <-jobCtx.Done():
		case <-lease.Lost():
			cancel(lease.Err())
		}
	}()
	defer func() {
		releaseCtx, stop := context.WithTimeout(context.WithoutCancel(ctx), r.Policy.RecordTimeout)
		defer stop()
		err = errors.Join(err, lease.Release(releaseCtx))
	}()
	defer func() {
		// 派发器可能先于监听协程观察到失租，收口时同步取消以保留原因和取消统计。
		if leaseLost(lease) {
			cancel(lease.Err())
		}
		report.Canceled = jobCtx.Err() != nil
		if report.Canceled {
			err = errors.Join(err, context.Cause(jobCtx))
		}
	}()
	cursor, err := lease.Cursor(jobCtx)
	if err != nil {
		return report, err
	}
	dispatchCtx, stopDispatch := context.WithTimeout(jobCtx, r.Policy.DispatchBudget)
	defer stopDispatch()
	scan := Scan{Now: r.now(), MinInterval: r.Policy.MinInterval}
	groups, err := r.discover(dispatchCtx, scan, cursor)
	if err != nil {
		return report, err
	}
	report.DiscoveredGroups = len(groups)
	err = r.dispatch(jobCtx, dispatchCtx, lease, scan, groups, &report)
	return report, err
}

// leaseLost 同步读取失租通知，不等待监听协程将信号传播到 context。
func leaseLost(lease Lease) bool {
	select {
	case <-lease.Lost():
		return true
	default:
		return false
	}
}

// discover 有界预取候选组；到尾部回绕，但不把预取位置写作已访问游标。
func (r *Runner) discover(ctx context.Context, scan Scan, cursor string) ([]string, error) {
	limit := r.Policy.MaxGroupVisits
	groups, err := r.Source.DiscoverGroups(ctx, scan, cursor, limit)
	if err != nil {
		return nil, err
	}
	if err = validateGroups(groups, cursor, limit); err != nil {
		return nil, err
	}
	if cursor == "" || len(groups) == limit {
		return groups, nil
	}
	head, err := r.Source.DiscoverGroups(ctx, scan, "", limit-len(groups))
	if err != nil {
		return nil, err
	}
	if err = validateGroups(head, "", limit-len(groups)); err != nil {
		return nil, err
	}
	for _, group := range head {
		if group > cursor {
			break
		}
		groups = append(groups, group)
	}
	return groups, nil
}

// validateGroups 拒绝违反 keyset 契约的适配器，避免重复和无界返回。
func validateGroups(groups []string, after string, limit int) error {
	if len(groups) > limit {
		return NewError("group_limit")
	}
	for _, group := range groups {
		if group == "" || group <= after {
			return NewError("group_order")
		}
		after = group
	}
	return nil
}

// availableGroup 在取候选前应用并发限制，防止满组占据队首。
func availableGroup(groups []groupState, start, limit int) int {
	for offset := 0; offset < len(groups); offset++ {
		i := (start + offset) % len(groups)
		if !groups[i].exhausted && groups[i].inflight < limit {
			return i
		}
	}
	return -1
}

// validateCandidate 防止错误适配器绕过时间、分组与本轮去重契约。
func validateCandidate(c Candidate, group groupState, scan Scan) error {
	if c.RecordKey == "" || c.GroupKey != group.key || c.ReadyAt.IsZero() || c.ReadyAt.After(scan.Now) || (c.LastCheckedAt != nil && c.LastCheckedAt.Add(scan.MinInterval).After(scan.Now)) {
		return NewError("candidate")
	}
	for _, key := range group.excluded {
		if c.RecordKey == key {
			return NewError("duplicate")
		}
	}
	return nil
}

// attempt 将业务调用与记时分离；业务超时后用仍有效的 Job 预算记时。
func (r *Runner) attempt(ctx context.Context, group int, candidate Candidate) completion {
	c := completion{group: group}
	attemptCtx, cancel := context.WithTimeout(ctx, r.Policy.AttemptTimeout)
	c.attempted, c.callbackErr = r.Advance(attemptCtx, candidate)
	cancel()
	if !c.attempted {
		return c
	}
	finished := r.now()
	recordCtx, stop := context.WithTimeout(ctx, r.Policy.RecordTimeout)
	defer stop()
	c.recordErr = r.Source.RecordAttempt(recordCtx, candidate, finished)
	return c
}
