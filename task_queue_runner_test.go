package taskstore

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestTaskQueueRunnerRunOnceProcessesQueuedTasks(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.GetDB().Close()

	ctx := context.Background()

	handler := new(testHandler)
	err = store.TaskHandlerAdd(ctx, handler, true)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 2; i++ {
		_, err = store.TaskDefinitionEnqueueByAlias(ctx, DefaultQueueName, handler.Alias(), map[string]any{})
		if err != nil {
			t.Fatal(err)
		}
	}

	queued, err := store.TaskQueueList(ctx, TaskQueueQuery().SetStatus(TaskQueueStatusQueued))
	if err != nil {
		t.Fatal(err)
	}
	if len(queued) != 2 {
		t.Fatalf("expected 2 queued tasks, got %d", len(queued))
	}

	runner := NewTaskQueueRunner(store, TaskQueueRunnerOptions{IntervalSeconds: 1, UnstuckMinutes: 1, QueueName: DefaultQueueName})

	err = runner.RunOnce(ctx)
	if err != nil {
		t.Error(err)
	}

	success, err := store.TaskQueueList(ctx, TaskQueueQuery().SetStatus(TaskQueueStatusSuccess))
	if err != nil {
		t.Error(err)
	}
	if len(success) != 2 {
		t.Errorf("expected 2 success tasks, got %d", len(success))
	}

	remainingQueued, err := store.TaskQueueList(ctx, TaskQueueQuery().SetStatus(TaskQueueStatusQueued))
	if err != nil {
		t.Error(err)
	}
	if len(remainingQueued) != 0 {
		t.Errorf("expected 0 remaining queued tasks, got %d", len(remainingQueued))
	}
}

func TestTaskQueueRunnerStartStop(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.GetDB().Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runner := NewTaskQueueRunner(store, TaskQueueRunnerOptions{IntervalSeconds: 1, UnstuckMinutes: 1, QueueName: DefaultQueueName})

	runner.Start(ctx)
	time.Sleep(50 * time.Millisecond)
	if !runner.IsRunning() {
		t.Error("expected runner to be running")
	}

	runner.Stop()
	time.Sleep(50 * time.Millisecond)
	if runner.IsRunning() {
		t.Error("expected runner to be stopped")
	}
}

// Test 1: Serial Processing (MaxConcurrency = 1, default)
func TestTaskQueueRunner_SerialProcessing(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.GetDB().Close()

	ctx := context.Background()

	// Create a handler that tracks execution times
	var mu sync.Mutex
	var executionTimes []time.Time

	handler := &delayedHandler{
		delay: 50 * time.Millisecond,
		onExecute: func() {
			mu.Lock()
			executionTimes = append(executionTimes, time.Now())
			mu.Unlock()
		},
	}

	err = store.TaskHandlerAdd(ctx, handler, true)
	if err != nil {
		t.Fatal(err)
	}

	// Enqueue 3 tasks
	for i := 0; i < 3; i++ {
		_, err = store.TaskDefinitionEnqueueByAlias(ctx, DefaultQueueName, handler.Alias(), map[string]any{})
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create runner with default MaxConcurrency (should be 1)
	runner := NewTaskQueueRunner(store, TaskQueueRunnerOptions{
		IntervalSeconds: 1,
		UnstuckMinutes:  1,
		QueueName:       DefaultQueueName,
	})

	err = runner.RunOnce(ctx)
	if err != nil {
		t.Error(err)
	}

	// Verify all tasks completed
	success, err := store.TaskQueueList(ctx, TaskQueueQuery().SetStatus(TaskQueueStatusSuccess))
	if err != nil {
		t.Error(err)
	}
	if len(success) != 3 {
		t.Errorf("expected 3 success tasks, got %d", len(success))
	}

	// Verify tasks were executed serially (no overlap)
	mu.Lock()
	defer mu.Unlock()
	if len(executionTimes) != 3 {
		t.Fatalf("expected 3 execution times, got %d", len(executionTimes))
	}

	// Each task should start after the previous one finishes (50ms delay)
	for i := 1; i < len(executionTimes); i++ {
		timeDiff := executionTimes[i].Sub(executionTimes[i-1])
		if timeDiff < 40*time.Millisecond {
			t.Errorf("Tasks should execute serially, got time diff: %v", timeDiff)
		}
	}
}

// Test 2: Concurrent Processing (MaxConcurrency > 1)
func TestTaskQueueRunner_ConcurrentProcessing(t *testing.T) {
	store, err := initStore() // Use in-memory DB now that it's fixed
	if err != nil {
		t.Fatal(err)
	}
	defer store.GetDB().Close()

	ctx := context.Background()

	// Track concurrent execution
	var mu sync.Mutex
	var startTimes []time.Time

	handler := &delayedHandler{
		delay: 100 * time.Millisecond,
		onExecute: func() {
			mu.Lock()
			startTimes = append(startTimes, time.Now())
			mu.Unlock()
		},
	}

	err = store.TaskHandlerAdd(ctx, handler, true)
	if err != nil {
		t.Fatal(err)
	}

	// Enqueue 5 tasks
	for i := 0; i < 5; i++ {
		_, err = store.TaskDefinitionEnqueueByAlias(ctx, DefaultQueueName, handler.Alias(), map[string]any{})
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create runner with MaxConcurrency = 3
	runner := NewTaskQueueRunner(store, TaskQueueRunnerOptions{
		IntervalSeconds: 1,
		UnstuckMinutes:  1,
		QueueName:       DefaultQueueName,
		MaxConcurrency:  3,
	})

	start := time.Now()
	err = runner.RunOnce(ctx)
	if err != nil {
		t.Error(err)
	}
	elapsed := time.Since(start)

	// Wait longer for all DB updates to complete (tasks take 100ms each)
	time.Sleep(3000 * time.Millisecond)

	// Verify all tasks completed
	success, err := store.TaskQueueList(ctx, TaskQueueQuery().SetStatus(TaskQueueStatusSuccess))
	if err != nil {
		t.Error(err)
	}
	if len(success) != 5 {
		// Debug: check other statuses
		queued, _ := store.TaskQueueList(ctx, TaskQueueQuery().SetStatus(TaskQueueStatusQueued))
		running, _ := store.TaskQueueList(ctx, TaskQueueQuery().SetStatus(TaskQueueStatusRunning))
		failed, _ := store.TaskQueueList(ctx, TaskQueueQuery().SetStatus(TaskQueueStatusFailed))
		t.Logf("Success: %d, Queued: %d, Running: %d, Failed: %d", len(success), len(queued), len(running), len(failed))
	}

	// With 5 tasks, 100ms each, and concurrency of 3:
	// - First 3 tasks run in parallel: 100ms
	// - Next 2 tasks run in parallel: 100ms
	// Total should be around 200ms, not 500ms (serial)
	if elapsed >= 350*time.Millisecond {
		t.Errorf("Concurrent execution should be faster than serial, took: %v", elapsed)
	}

	// Verify at least some tasks ran concurrently
	mu.Lock()
	defer mu.Unlock()
	if len(startTimes) != 5 {
		t.Fatalf("expected 5 start times, got %d", len(startTimes))
	}

	// Check if at least 2 tasks started within 50ms of each other (indicating concurrency)
	concurrentStarts := 0
	for i := 1; i < len(startTimes); i++ {
		if startTimes[i].Sub(startTimes[i-1]) < 50*time.Millisecond {
			concurrentStarts++
		}
	}
	if concurrentStarts < 1 {
		t.Error("At least some tasks should run concurrently")
	}
}

// Test 3: Concurrency Limit Enforcement
func TestTaskQueueRunner_ConcurrencyLimitEnforced(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.GetDB().Close()

	ctx := context.Background()

	// Track concurrent execution count
	var currentConcurrent int32
	var maxConcurrent int32

	handler := &delayedHandler{
		delay: 200 * time.Millisecond,
		onExecute: func() {
			current := atomic.AddInt32(&currentConcurrent, 1)
			// Update max if needed
			for {
				max := atomic.LoadInt32(&maxConcurrent)
				if current <= max || atomic.CompareAndSwapInt32(&maxConcurrent, max, current) {
					break
				}
			}
		},
		onComplete: func() {
			atomic.AddInt32(&currentConcurrent, -1)
		},
	}

	err = store.TaskHandlerAdd(ctx, handler, true)
	if err != nil {
		t.Fatal(err)
	}

	// Enqueue 10 tasks
	for i := 0; i < 10; i++ {
		_, err = store.TaskDefinitionEnqueueByAlias(ctx, DefaultQueueName, handler.Alias(), map[string]any{})
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create runner with MaxConcurrency = 2
	runner := NewTaskQueueRunner(store, TaskQueueRunnerOptions{
		IntervalSeconds: 1,
		UnstuckMinutes:  1,
		QueueName:       DefaultQueueName,
		MaxConcurrency:  2,
	})

	err = runner.RunOnce(ctx)
	if err != nil {
		t.Error(err)
	}

	// Verify all tasks completed
	success, err := store.TaskQueueList(ctx, TaskQueueQuery().SetStatus(TaskQueueStatusSuccess))
	if err != nil {
		t.Error(err)
	}
	if len(success) != 10 {
		t.Errorf("expected 10 success tasks, got %d", len(success))
	}

	// Verify max concurrent never exceeded 2
	max := atomic.LoadInt32(&maxConcurrent)
	if max > 2 {
		t.Errorf("Max concurrent tasks should not exceed MaxConcurrency setting, got: %d", max)
	}
}

// Test 4: Graceful Shutdown with In-Flight Tasks
func TestTaskQueueRunner_GracefulShutdownConcurrent(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.GetDB().Close()

	ctx := context.Background()

	// Track task completion
	var completedCount int32

	handler := &delayedHandler{
		delay: 2000 * time.Millisecond,
		onComplete: func() {
			atomic.AddInt32(&completedCount, 1)
		},
	}

	err = store.TaskHandlerAdd(ctx, handler, true)
	if err != nil {
		t.Fatal(err)
	}

	// Enqueue 3 tasks
	for i := 0; i < 3; i++ {
		_, err = store.TaskDefinitionEnqueueByAlias(ctx, DefaultQueueName, handler.Alias(), map[string]any{})
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create runner with MaxConcurrency = 3
	runner := NewTaskQueueRunner(store, TaskQueueRunnerOptions{
		IntervalSeconds: 1,
		UnstuckMinutes:  1,
		QueueName:       DefaultQueueName,
		MaxConcurrency:  3,
	})

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	runner.Start(runCtx)

	// Wait for tasks to start
	time.Sleep(1500 * time.Millisecond)

	// Stop the runner (should wait for in-flight tasks)
	stopStart := time.Now()
	runner.Stop()
	stopDuration := time.Since(stopStart)

	// Verify Stop() waited for tasks to complete
	// Tasks take 1000ms, we waited 500ms before stopping, so Stop should wait ~500ms
	if stopDuration < 400*time.Millisecond {
		t.Errorf("Stop should wait for in-flight tasks, took: %v", stopDuration)
	}

	// Verify all tasks completed successfully
	completed := atomic.LoadInt32(&completedCount)
	if completed != 3 {
		t.Errorf("All in-flight tasks should complete, got: %d", completed)
	}

	success, err := store.TaskQueueList(ctx, TaskQueueQuery().SetStatus(TaskQueueStatusSuccess))
	if err != nil {
		t.Error(err)
	}
	if len(success) != 3 {
		t.Errorf("expected 3 success tasks, got %d", len(success))
	}
}

// Test 5: Default MaxConcurrency
func TestTaskQueueRunner_DefaultMaxConcurrency(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.GetDB().Close()

	ctx := context.Background()

	var mu sync.Mutex
	var executionTimes []time.Time

	handler := &delayedHandler{
		delay: 50 * time.Millisecond,
		onExecute: func() {
			mu.Lock()
			executionTimes = append(executionTimes, time.Now())
			mu.Unlock()
		},
	}

	err = store.TaskHandlerAdd(ctx, handler, true)
	if err != nil {
		t.Fatal(err)
	}

	// Enqueue 3 tasks
	for i := 0; i < 3; i++ {
		_, err = store.TaskDefinitionEnqueueByAlias(ctx, DefaultQueueName, handler.Alias(), map[string]any{})
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create runner WITHOUT specifying MaxConcurrency
	runner := NewTaskQueueRunner(store, TaskQueueRunnerOptions{
		IntervalSeconds: 1,
		UnstuckMinutes:  1,
		QueueName:       DefaultQueueName,
		// MaxConcurrency not set - should default to 1
	})

	err = runner.RunOnce(ctx)
	if err != nil {
		t.Error(err)
	}

	// Verify tasks executed serially
	mu.Lock()
	defer mu.Unlock()
	if len(executionTimes) != 3 {
		t.Fatalf("expected 3 execution times, got %d", len(executionTimes))
	}

	// Verify serial execution (tasks don't overlap)
	for i := 1; i < len(executionTimes); i++ {
		timeDiff := executionTimes[i].Sub(executionTimes[i-1])
		if timeDiff < 40*time.Millisecond {
			t.Errorf("Default should be serial processing, got time diff: %v", timeDiff)
		}
	}
}

// Helper: delayedHandler for testing concurrent execution
type delayedHandler struct {
	TaskDefinitionHandlerBase
	delay      time.Duration
	onExecute  func()
	onComplete func()
}

func (h *delayedHandler) Alias() string {
	return "DelayedHandler"
}

func (h *delayedHandler) Title() string {
	return "Delayed Test Handler"
}

func (h *delayedHandler) Description() string {
	return "Handler with configurable delay for testing"
}

func (h *delayedHandler) Handle() bool {
	if h.onExecute != nil {
		h.onExecute()
	}

	time.Sleep(h.delay)

	if h.onComplete != nil {
		h.onComplete()
	}

	return true
}
