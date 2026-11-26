package polling

import "time"

// Options 是接入方可选的容量覆盖项；零值使用公共默认值，不控制计划启停。
type Options struct {
	MaxItems       int `json:"maxItems,omitempty"`
	MaxConcurrency int `json:"maxConcurrency,omitempty"`
}

// DefaultPolicy 提供无需业务逐项填写的有界策略，内部参数与生成 Activity 的预算配套。
func DefaultPolicy(jobKey string) JobPolicy {
	return JobPolicy{JobKey: jobKey, MaxItems: 100, MaxConcurrency: 10, MaxGroupConcurrency: 2, MaxGroupVisits: 1000, CandidateBatchSize: 16, DispatchBudget: 4 * time.Minute, AttemptTimeout: 5 * time.Minute, RecordTimeout: 5 * time.Second}
}

// Policy 应用可选容量覆盖并执行硬边界校验；只有默认并发会随小批次缩小，显式冲突报错。
func (o Options) Policy(jobKey string) (JobPolicy, error) {
	p := DefaultPolicy(jobKey)
	if o.MaxItems != 0 {
		p.MaxItems = o.MaxItems
	}
	if o.MaxConcurrency != 0 {
		p.MaxConcurrency = o.MaxConcurrency
	} else {
		p.MaxConcurrency = min(p.MaxConcurrency, p.MaxItems)
	}
	p.MaxGroupConcurrency = min(p.MaxGroupConcurrency, p.MaxConcurrency)
	return p, p.Validate()
}
