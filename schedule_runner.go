package taskstore

import (
	"context"
	"log"
	"sync/atomic"
	"time"
)

type ScheduleRunnerOptions struct {
	IntervalSeconds int
	Logger          *log.Logger
}

type ScheduleRunnerInterface interface {
	Start(ctx context.Context)
	Stop()
	IsRunning() bool
	RunOnce(ctx context.Context) error
	SetInitialRuns(ctx context.Context) error
}

type scheduleRunner struct {
	store   StoreInterface
	opts    ScheduleRunnerOptions
	running atomic.Bool
	stopCh  chan struct{}
}

func NewScheduleRunner(store StoreInterface, opts ScheduleRunnerOptions) ScheduleRunnerInterface {
	if opts.IntervalSeconds <= 0 {
		opts.IntervalSeconds = 60
	}

	return &scheduleRunner{
		store:  store,
		opts:   opts,
		stopCh: make(chan struct{}, 1),
	}
}

func (r *scheduleRunner) Start(ctx context.Context) {
	if !r.running.CompareAndSwap(false, true) {
		return
	}

	go func() {
		ticker := time.NewTicker(time.Duration(r.opts.IntervalSeconds) * time.Second)
		defer ticker.Stop()
		defer r.running.Store(false)

		for {
			if !r.shouldContinue(ctx) {
				return
			}

			if err := r.RunOnce(ctx); err != nil {
				r.logf("ScheduleRunner: RunOnce error: %v", err)
			}

			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				return
			case <-r.stopCh:
				return
			}
		}
	}()
}

func (r *scheduleRunner) Stop() {
	if !r.running.Load() {
		return
	}

	select {
	case r.stopCh <- struct{}{}:
	default:
	}
}

func (r *scheduleRunner) IsRunning() bool {
	return r.running.Load()
}

func (r *scheduleRunner) RunOnce(ctx context.Context) error {
	schedules, err := r.findActiveSchedulesToBeRun(ctx)
	if err != nil {
		return err
	}

	for _, s := range schedules {
		if err := r.runSchedule(ctx, s); err != nil {
			r.logf("ScheduleRunner: error running schedule %s: %v", s.GetID(), err)
		}
	}

	return nil
}

// SetInitialRuns initializes the NextRunAt field for all active schedules
// that still have NULL_DATETIME (not yet initialized). Call once at startup.
//
// Business Logic:
//  1. Query all schedules with status "active".
//  2. Skip any schedule whose NextRunAt is not NULL_DATETIME (already
//     initialized).
//  3. Skip any schedule that has already reached its max executions —
//     it should have been marked completed already, but if it wasn't
//     (e.g. crash between increment and status update), mark it now.
//  4. For each uninitialized schedule, call GetNextOccurrence:
//     a. If a valid time is returned, store it as NextRunAt.
//     b. If any error is returned (ErrNoMoreRuns or other), set NextRunAt
//     to MAX_DATETIME and mark status "completed". This is consistent
//     with UpdateNextRunAt, which sets MAX_DATETIME for any error.
//     - ErrNoMoreRuns: the schedule is exhausted (e.g. one-time schedule
//     with a past start time, or expired endsAt).
//     - Other errors (e.g. "interval must be positive"): the rule is
//     misconfigured. Setting MAX_DATETIME stops the bleeding.
func (r *scheduleRunner) SetInitialRuns(ctx context.Context) error {
	query := NewScheduleQuery().SetStatus("active")
	schedules, err := r.store.ScheduleList(ctx, query)
	if err != nil {
		return err
	}

	for _, s := range schedules {
		if !isNullDateTime(s.GetNextRunAt()) {
			continue
		}

		// A schedule that has reached max executions but still has
		// NULL_DATETIME as NextRunAt (e.g. crash between increment and
		// status update) should be marked completed immediately.
		if s.HasReachedMaxExecutions() {
			s.SetNextRunAt(MAX_DATETIME)
			s.SetStatus("completed")
			if err := r.store.ScheduleUpdate(ctx, s); err != nil {
				r.logf("ScheduleRunner: error marking completed schedule %s: %v", s.GetID(), err)
			}
			continue
		}

		next, err := s.GetNextOccurrence()
		if err != nil {
			// Any error (ErrNoMoreRuns or misconfigured rule) — set
			// MAX_DATETIME and mark completed, consistent with
			// UpdateNextRunAt's error handling.
			s.SetNextRunAt(MAX_DATETIME)
			s.SetStatus("completed")
			if err := r.store.ScheduleUpdate(ctx, s); err != nil {
				r.logf("ScheduleRunner: error marking completed schedule %s: %v", s.GetID(), err)
			}
			continue
		}

		s.SetNextRunAt(next)
		if err := r.store.ScheduleUpdate(ctx, s); err != nil {
			r.logf("ScheduleRunner: error updating schedule %s: %v", s.GetID(), err)
		}
	}

	return nil
}

func (r *scheduleRunner) shouldContinue(ctx context.Context) bool {
	if ctx != nil && ctx.Err() != nil {
		return false
	}

	return r.running.Load()
}

func (r *scheduleRunner) findActiveSchedules(ctx context.Context) ([]ScheduleInterface, error) {
	query := NewScheduleQuery().SetStatus("active")
	return r.store.ScheduleList(ctx, query)
}

// findActiveSchedulesToBeRun returns all active schedules that are due to
// run now.
//
// Business Logic:
//  1. Query all schedules with status "active".
//  2. For each schedule that has reached its end date or max executions,
//     mark status "completed" and skip.
//  3. For each schedule whose NextRunAt is NULL_DATETIME (not yet
//     initialized), call UpdateNextRunAt. If the result is MAX_DATETIME
//     (exhausted), mark status "completed" and skip.
//  4. Return all remaining schedules where IsDue() returns true.
func (r *scheduleRunner) findActiveSchedulesToBeRun(ctx context.Context) ([]ScheduleInterface, error) {
	schedules, err := r.findActiveSchedules(ctx)
	if err != nil {
		return nil, err
	}

	due := make([]ScheduleInterface, 0, len(schedules))

	for _, s := range schedules {
		// Mark schedules that have reached end or max executions as completed
		if s.HasReachedEndDate() || s.HasReachedMaxExecutions() {
			s.SetStatus("completed")
			if err := r.store.ScheduleUpdate(ctx, s); err != nil {
				r.logf("ScheduleRunner: error marking completed schedule %s: %v", s.GetID(), err)
			}
			continue
		}

		// Initialize next run if needed
		if isNullDateTime(s.GetNextRunAt()) {
			s.UpdateNextRunAt()
			if err := r.store.ScheduleUpdate(ctx, s); err != nil {
				r.logf("ScheduleRunner: error initializing next run for schedule %s: %v", s.GetID(), err)
			}
			// If initialization found no more runs, UpdateNextRunAt set
			// NextRunAt to MAX_DATETIME. Mark as completed and skip.
			if isMaxDateTime(s.GetNextRunAt()) {
				s.SetStatus("completed")
				if err := r.store.ScheduleUpdate(ctx, s); err != nil {
					r.logf("ScheduleRunner: error marking completed schedule %s: %v", s.GetID(), err)
				}
				continue
			}
		}

		if s.IsDue() {
			due = append(due, s)
		}
	}

	return due, nil
}

// runSchedule executes a single schedule if it is due.
//
// Business Logic:
//  1. Double-check termination: if HasReachedEndDate or
//     HasReachedMaxExecutions, mark status "completed" and return.
//  2. If IsDue() is false, return without action.
//  3. Look up the associated task definition; if missing, log and return.
//  4. Enqueue the task via TaskDefinitionEnqueueByAlias.
//  5. Update LastRunAt, increment ExecutionCount, recalculate NextRunAt
//     via UpdateNextRunAt.
//  6. Mark status "completed" if any termination condition holds:
//     HasReachedEndDate, HasReachedMaxExecutions, or NextRunAt is
//     MAX_DATETIME (ErrNoMoreRuns — e.g. a one-time schedule that just
//     fired its only run).
func (r *scheduleRunner) runSchedule(ctx context.Context, s ScheduleInterface) error {
	// Double-check termination conditions
	if s.HasReachedEndDate() || s.HasReachedMaxExecutions() {
		s.SetStatus("completed")
		return r.store.ScheduleUpdate(ctx, s)
	}

	if !s.IsDue() {
		return nil
	}

	taskDef, err := r.store.TaskDefinitionFindByID(ctx, s.GetTaskDefinitionID())
	if err != nil {
		return err
	}

	if taskDef == nil {
		r.logf("ScheduleRunner: task definition not found for schedule %s", s.GetID())
		return nil
	}

	_, err = r.store.TaskDefinitionEnqueueByAlias(ctx, s.GetQueueName(), taskDef.GetAlias(), s.GetTaskParameters())
	if err != nil {
		return err
	}

	s.UpdateLastRunAt()
	s.IncrementExecutionCount()
	s.UpdateNextRunAt()

	if s.HasReachedEndDate() || s.HasReachedMaxExecutions() {
		s.SetStatus("completed")
	} else if isMaxDateTime(s.GetNextRunAt()) {
		// No more occurrences (e.g. one-time schedule that already fired).
		s.SetStatus("completed")
	}

	return r.store.ScheduleUpdate(ctx, s)
}

func (r *scheduleRunner) logf(format string, args ...interface{}) {
	if r.opts.Logger != nil {
		r.opts.Logger.Printf(format, args...)
	}
}
