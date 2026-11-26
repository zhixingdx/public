package workflow

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zhixingdx/public/pkg/logger/kratoslog"
	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
)

// Memo key 常量
const (
	MemoKeyService     = "service"
	MemoKeyType        = "type"
	MemoKeyDescription = "description"
	MemoKeyInterval    = "interval"
	MemoKeyManagedBy   = "managed-by"

	MemoTypeScheduled      = "scheduled"
	MemoManagedByAutoSetup = "auto-setup"
)

// ScheduleConfig 定时任务配置
// 每个字段含义：
//   - ID: Schedule 标识，最终 Schedule ID 为 "{serviceName}-{ID}"
//   - Description: 任务描述，写入 Memo 供 Temporal UI 展示
//   - Workflow: Temporal Workflow 函数引用
//   - Interval: 固定执行间隔（与 Cron 二选一）
//   - Cron: Cron 表达式（与 Interval 二选一）
//   - Args: 传给 Workflow 的参数（可选，与 Temporal ScheduleWorkflowAction.Args 一致）
type ScheduleConfig struct {
	ID          string
	Description string
	Workflow    any
	Interval    time.Duration
	Cron        string
	Args        []interface{}
}

// EnsureSchedules 批量检查并创建 Temporal 定时任务（服务启动时调用）
// 逐个检查 Schedule 是否已存在，不存在则创建，存在则跳过
// 固定策略：Overlap=SKIP, TaskQueue=serviceName, Memo 自动填充
func EnsureSchedules(temporalClient client.Client, serviceName string, schedules []ScheduleConfig, logger *log.Helper) {
	for _, s := range schedules {
		ensureSchedule(temporalClient, serviceName, s, logger)
	}
}

// EnsureSchedulesContext 在调用方截止时间内注册缺失计划，并把依赖故障交还给启动降级流程。
func EnsureSchedulesContext(ctx context.Context, temporalClient client.Client, serviceName string, schedules []ScheduleConfig, logger *log.Helper) error {
	for _, schedule := range schedules {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := ensureScheduleContext(ctx, temporalClient, serviceName, schedule, logger); err != nil {
			return err
		}
	}
	return nil
}

// ensureScheduleContext 仅在 Temporal 明确表示计划不存在时创建，避免依赖故障被误认为缺失。
func ensureScheduleContext(ctx context.Context, temporalClient client.Client, serviceName string, schedule ScheduleConfig, logger *log.Helper) error {
	id := fmt.Sprintf("%s-%s", serviceName, schedule.ID)
	if _, err := temporalClient.ScheduleClient().GetHandle(ctx, id).Describe(ctx); err == nil {
		logger.Infow("msg", "定时任务已存在，跳过创建", "scheduleID", id)
		return nil
	} else {
		var notFound *serviceerror.NotFound
		if !errors.As(err, &notFound) {
			return fmt.Errorf("检查定时任务 %s 失败: %w", id, err)
		}
	}
	options, interval := scheduleOptions(serviceName, schedule)
	if _, err := temporalClient.ScheduleClient().Create(ctx, options); err != nil {
		var alreadyExists *serviceerror.AlreadyExists
		if !errors.As(err, &alreadyExists) {
			return fmt.Errorf("创建定时任务 %s 失败: %w", id, err)
		}
		return nil
	}
	logger.Infow("msg", "创建定时任务成功", "scheduleID", id, "interval", interval)
	return nil
}

// ensureSchedule 单个定时任务的检查与创建
func ensureSchedule(temporalClient client.Client, serviceName string, s ScheduleConfig, logger *log.Helper) {
	scheduleID := fmt.Sprintf("%s-%s", serviceName, s.ID)
	ctx := context.Background()
	handle := temporalClient.ScheduleClient().GetHandle(ctx, scheduleID)

	_, err := handle.Describe(ctx)
	if err == nil {
		logger.Infow("msg", "定时任务已存在，跳过创建", "scheduleID", scheduleID)
		return
	}

	options, intervalStr := scheduleOptions(serviceName, s)
	_, err = temporalClient.ScheduleClient().Create(ctx, options)
	if err != nil {
		logger.Errorw("msg", "创建定时任务失败", "scheduleID", scheduleID, "err", err)
		return
	}
	logger.Infow("msg", "创建定时任务成功", "scheduleID", scheduleID, "interval", intervalStr)
}

// scheduleOptions 保持新旧入口完全一致的 ID、任务队列、重叠策略和 Memo 配置。
func scheduleOptions(serviceName string, s ScheduleConfig) (client.ScheduleOptions, string) {
	scheduleID := fmt.Sprintf("%s-%s", serviceName, s.ID)
	spec := client.ScheduleSpec{}
	intervalStr := ""
	if s.Cron != "" {
		spec.CronExpressions = []string{s.Cron}
		intervalStr = s.Cron
	} else {
		spec.Intervals = []client.ScheduleIntervalSpec{{Every: s.Interval}}
		intervalStr = s.Interval.String()
	}

	action := &client.ScheduleWorkflowAction{
		Workflow:  s.Workflow,
		TaskQueue: serviceName,
	}
	if len(s.Args) > 0 {
		action.Args = s.Args
	}

	return client.ScheduleOptions{
		ID:      scheduleID,
		Spec:    spec,
		Action:  action,
		Overlap: enumspb.SCHEDULE_OVERLAP_POLICY_SKIP,
		Memo: map[string]any{
			MemoKeyService:     serviceName,
			MemoKeyType:        MemoTypeScheduled,
			MemoKeyDescription: s.Description,
			MemoKeyInterval:    intervalStr,
			MemoKeyManagedBy:   MemoManagedByAutoSetup,
		},
	}, intervalStr
}
