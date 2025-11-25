package taskstore

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"github.com/dracory/sb"
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
			r.logf("ScheduleRunner: error running schedule %s: %v", s.ID(), err)
		}
	}

	return nil
}

func (r *scheduleRunner) SetInitialRuns(ctx context.Context) error {
	query := NewScheduleQuery().SetStatus("active")
	schedules, err := r.store.ScheduleList(ctx, query)
	if err != nil {
		return err
	}

	for _, s := range schedules {
		if s.NextRunAt() != sb.NULL_DATETIME {
			continue
		}

		next, err := s.GetNextOccurrence()
		if err != nil {
			r.logf("ScheduleRunner: error calculating initial next run for schedule %s: %v", s.ID(), err)
			continue
		}

		s.SetNextRunAt(next)
		if err := r.store.ScheduleUpdate(ctx, s); err != nil {
			r.logf("ScheduleRunner: error updating schedule %s: %v", s.ID(), err)
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
				r.logf("ScheduleRunner: error marking completed schedule %s: %v", s.ID(), err)
			}
			continue
		}

		// Initialize next run if needed
		if s.NextRunAt() == sb.NULL_DATETIME {
			s.UpdateNextRunAt()
			if err := r.store.ScheduleUpdate(ctx, s); err != nil {
				r.logf("ScheduleRunner: error initializing next run for schedule %s: %v", s.ID(), err)
			}
		}

		if s.IsDue() {
			due = append(due, s)
		}
	}

	return due, nil
}

func (r *scheduleRunner) runSchedule(ctx context.Context, s ScheduleInterface) error {
	// Double-check termination conditions
	if s.HasReachedEndDate() || s.HasReachedMaxExecutions() {
		s.SetStatus("completed")
		return r.store.ScheduleUpdate(ctx, s)
	}

	if !s.IsDue() {
		return nil
	}

	taskDef, err := r.store.TaskDefinitionFindByID(ctx, s.TaskDefinitionID())
	if err != nil {
		return err
	}

	if taskDef == nil {
		r.logf("ScheduleRunner: task definition not found for schedule %s", s.ID())
		return nil
	}

	_, err = r.store.TaskDefinitionEnqueueByAlias(ctx, s.QueueName(), taskDef.Alias(), s.TaskParameters())
	if err != nil {
		return err
	}

	s.UpdateLastRunAt()
	s.IncrementExecutionCount()
	s.UpdateNextRunAt()

	if s.HasReachedEndDate() || s.HasReachedMaxExecutions() {
		s.SetStatus("completed")
	}

	return r.store.ScheduleUpdate(ctx, s)
}

func (r *scheduleRunner) logf(format string, args ...interface{}) {
	if r.opts.Logger != nil {
		r.opts.Logger.Printf(format, args...)
	}
}
