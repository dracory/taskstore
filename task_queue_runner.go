package taskstore

import (
	"context"
	"log"
	"sync/atomic"
	"time"
)

type TaskQueueRunnerOptions struct {
	IntervalSeconds int
	UnstuckMinutes  int
	QueueName       string
	Logger          *log.Logger
}

type TaskQueueRunnerInterface interface {
	Start(ctx context.Context)
	Stop()
	IsRunning() bool
	RunOnce(ctx context.Context) error
}

type taskQueueRunner struct {
	store   StoreInterface
	opts    TaskQueueRunnerOptions
	running atomic.Bool
	stopCh  chan struct{}
}

func NewTaskQueueRunner(store StoreInterface, opts TaskQueueRunnerOptions) TaskQueueRunnerInterface {
	if opts.IntervalSeconds <= 0 {
		opts.IntervalSeconds = 10
	}

	if opts.UnstuckMinutes <= 0 {
		opts.UnstuckMinutes = 1
	}

	if opts.QueueName == "" {
		opts.QueueName = DefaultQueueName
	}

	return &taskQueueRunner{
		store:  store,
		opts:   opts,
		stopCh: make(chan struct{}, 1),
	}
}

func (r *taskQueueRunner) Start(ctx context.Context) {
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
				r.logf("TaskQueueRunner: RunOnce error: %v", err)
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

func (r *taskQueueRunner) Stop() {
	if !r.running.Load() {
		return
	}

	select {
	case r.stopCh <- struct{}{}:
	default:
	}
}

func (r *taskQueueRunner) IsRunning() bool {
	return r.running.Load()
}

func (r *taskQueueRunner) RunOnce(ctx context.Context) error {
	queueName := normalizeQueueName(r.opts.QueueName)

	for {
		if ctx != nil && ctx.Err() != nil {
			return ctx.Err()
		}

		queuedTask, err := r.store.TaskQueueClaimNext(ctx, queueName)
		if err != nil {
			return err
		}

		if queuedTask == nil {
			return nil
		}

		_, err = r.store.TaskQueueProcessTask(ctx, queuedTask)
		if err != nil {
			r.logf("TaskQueueRunner: error processing task %s: %v", queuedTask.ID(), err)
		}
	}
}

func (r *taskQueueRunner) shouldContinue(ctx context.Context) bool {
	if ctx != nil && ctx.Err() != nil {
		return false
	}

	return r.running.Load()
}

func (r *taskQueueRunner) logf(format string, args ...interface{}) {
	if r.opts.Logger != nil {
		r.opts.Logger.Printf(format, args...)
	}
}
