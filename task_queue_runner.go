package taskstore

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type TaskQueueRunnerOptions struct {
	IntervalSeconds int
	UnstuckMinutes  int
	QueueName       string
	Logger          *log.Logger
	MaxConcurrency  int // 0 or 1 = serial, >1 = concurrent (default: 1)
}

type TaskQueueRunnerInterface interface {
	Start(ctx context.Context)
	Stop()
	IsRunning() bool
	RunOnce(ctx context.Context) error
}

type taskQueueRunner struct {
	store     StoreInterface
	opts      TaskQueueRunnerOptions
	running   atomic.Bool
	stopCh    chan struct{}
	taskWg    sync.WaitGroup // Tracks spawned task goroutines
	semaphore chan struct{}  // Concurrency limiter
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

	// Default MaxConcurrency to 1 (serial) if not specified
	if opts.MaxConcurrency <= 0 {
		opts.MaxConcurrency = 1
	}

	return &taskQueueRunner{
		store:     store,
		opts:      opts,
		stopCh:    make(chan struct{}, 1),
		semaphore: make(chan struct{}, opts.MaxConcurrency),
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

	// Wait for all spawned task goroutines to complete
	r.taskWg.Wait()
}

func (r *taskQueueRunner) IsRunning() bool {
	return r.running.Load()
}

func (r *taskQueueRunner) RunOnce(ctx context.Context) error {
	if r.opts.MaxConcurrency == 1 {
		return r.runOnceSerial(ctx)
	}
	return r.runOnceConcurrent(ctx)
}

// runOnceSerial processes tasks one at a time (original behavior)
func (r *taskQueueRunner) runOnceSerial(ctx context.Context) error {
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

// runOnceConcurrent processes multiple tasks concurrently up to MaxConcurrency limit
func (r *taskQueueRunner) runOnceConcurrent(ctx context.Context) error {
	queueName := normalizeQueueName(r.opts.QueueName)

	// Defer waiting for all spawned goroutines to complete
	defer r.taskWg.Wait()

	for {
		if ctx != nil && ctx.Err() != nil {
			return ctx.Err()
		}

		// Claim next task (don't hold semaphore during claim)
		queuedTask, err := r.store.TaskQueueClaimNext(ctx, queueName)
		if err != nil {
			return err
		}

		if queuedTask == nil {
			return nil // Will wait for spawned goroutines due to defer
		}

		// Acquire semaphore slot (blocks if at max concurrency)
		select {
		case r.semaphore <- struct{}{}:
			// Got a slot, proceed
		case <-ctx.Done():
			return ctx.Err()
		}

		// Track the goroutine
		r.taskWg.Add(1)

		// Spawn goroutine to process the task
		go func(task TaskQueueInterface) {
			defer func() {
				<-r.semaphore   // Release semaphore slot
				r.taskWg.Done() // Mark goroutine as complete
			}()

			_, processErr := r.store.TaskQueueProcessTask(ctx, task)
			if processErr != nil {
				r.logf("TaskQueueRunner: error processing task %s: %v", task.ID(), processErr)
			}
		}(queuedTask)
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
