package taskstore

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/dromara/carbon/v2"
	_ "modernc.org/sqlite"
)

func TestScheduleRunnerRunOnceEnqueuesDueSchedule(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store, err := NewStore(NewStoreOptions{
		TaskDefinitionTableName: "task_definitions",
		TaskQueueTableName:      "task_queue",
		ScheduleTableName:       "schedules",
		DB:                      db,
		AutomigrateEnabled:      true,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	taskDef := NewTaskDefinition()
	taskDef.SetAlias("test-task")
	err = store.TaskDefinitionCreate(ctx, taskDef)
	if err != nil {
		t.Fatal(err)
	}

	schedule := NewSchedule()
	schedule.SetName("Test Schedule")
	schedule.SetStatus("active")
	schedule.SetQueueName("default")
	schedule.SetTaskDefinitionID(taskDef.GetID())

	now := carbon.Now(carbon.UTC)
	past := now.AddMinutes(-1).ToDateTimeString(carbon.UTC)

	schedule.SetStartAt(past)
	schedule.SetNextRunAt(past)
	schedule.GetRecurrenceRule().SetFrequency(FrequencyMinutely)
	schedule.GetRecurrenceRule().SetInterval(1)
	schedule.GetRecurrenceRule().SetStartsAt(past)

	err = store.ScheduleCreate(ctx, schedule)
	if err != nil {
		t.Fatal(err)
	}

	runner := NewScheduleRunner(store, ScheduleRunnerOptions{IntervalSeconds: 1})

	err = runner.RunOnce(ctx)
	if err != nil {
		t.Error(err)
	}

	queueList, err := store.TaskQueueList(ctx, TaskQueueQuery())
	if err != nil {
		t.Error(err)
	}
	if len(queueList) != 1 {
		t.Errorf("expected 1 queued task, got %d", len(queueList))
	}
	if taskDef.GetID() != queueList[0].GetTaskID() {
		t.Errorf("expected task ID %s, got %s", taskDef.GetID(), queueList[0].GetTaskID())
	}

	updatedSchedule, err := store.ScheduleFindByID(ctx, schedule.GetID())
	if err != nil {
		t.Error(err)
	}
	if updatedSchedule == nil {
		t.Fatal("expected updatedSchedule to not be nil")
	}
	if updatedSchedule.GetLastRunAt() == NULL_DATETIME {
		t.Error("expected LastRunAt to not be NULL_DATETIME")
	}
	if !carbon.Parse(updatedSchedule.GetNextRunAt(), carbon.UTC).Gt(carbon.Now(carbon.UTC)) {
		t.Error("expected NextRunAt to be in the future")
	}
	if updatedSchedule.GetExecutionCount() != 1 {
		t.Errorf("expected execution count 1, got %d", updatedSchedule.GetExecutionCount())
	}
	if updatedSchedule.GetStatus() != "active" {
		t.Errorf("expected status 'active', got %s", updatedSchedule.GetStatus())
	}
}

func TestScheduleRunnerSetInitialRuns(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store, err := NewStore(NewStoreOptions{
		TaskDefinitionTableName: "task_definitions",
		TaskQueueTableName:      "task_queue",
		ScheduleTableName:       "schedules",
		DB:                      db,
		AutomigrateEnabled:      true,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	schedule := NewSchedule()
	schedule.SetName("Initial Runs")
	schedule.SetStatus("active")
	schedule.SetQueueName("default")

	startAt := carbon.Now(carbon.UTC).AddMinutes(5).ToDateTimeString(carbon.UTC)
	schedule.SetStartAt(startAt)
	schedule.SetNextRunAt(NULL_DATETIME)
	schedule.GetRecurrenceRule().SetFrequency(FrequencyDaily)
	schedule.GetRecurrenceRule().SetInterval(1)
	schedule.GetRecurrenceRule().SetStartsAt(startAt)

	err = store.ScheduleCreate(ctx, schedule)
	if err != nil {
		t.Fatal(err)
	}

	runner := NewScheduleRunner(store, ScheduleRunnerOptions{IntervalSeconds: 1})

	err = runner.SetInitialRuns(ctx)
	if err != nil {
		t.Error(err)
	}

	updatedSchedule, err := store.ScheduleFindByID(ctx, schedule.GetID())
	if err != nil {
		t.Error(err)
	}
	if updatedSchedule == nil {
		t.Fatal("expected updatedSchedule to not be nil")
	}
	if updatedSchedule.GetNextRunAt() == NULL_DATETIME {
		t.Error("expected NextRunAt to not be NULL_DATETIME")
	}
	if updatedSchedule.GetStatus() != "active" {
		t.Errorf("expected status 'active', got %s", updatedSchedule.GetStatus())
	}
}

func TestScheduleRunner_shouldContinue(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.GetDB().Close()

	runner := NewScheduleRunner(store, ScheduleRunnerOptions{IntervalSeconds: 1})

	// Test with running state
	tests := []struct {
		name     string
		running  bool
		ctx      context.Context
		expected bool
	}{
		{
			name:     "running with valid context",
			running:  true,
			ctx:      context.Background(),
			expected: true,
		},
		{
			name:     "not running",
			running:  false,
			ctx:      context.Background(),
			expected: false,
		},
		{
			name:    "running with cancelled context",
			running: true,
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			expected: false,
		},
		{
			name:     "running with nil context",
			running:  true,
			ctx:      nil,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := runner.(*scheduleRunner)
			r.running.Store(tt.running)
			got := r.shouldContinue(tt.ctx)
			if got != tt.expected {
				t.Errorf("shouldContinue() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestScheduleRunner_logf(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.GetDB().Close()

	// Test with nil logger (should not panic)
	runner := NewScheduleRunner(store, ScheduleRunnerOptions{IntervalSeconds: 1})
	r := runner.(*scheduleRunner)
	r.logf("test message %s", "arg") // Should not panic

	// Test with logger (should not panic)
	logger := log.New(&testLoggerWriter{}, "", 0)
	runnerWithLogger := NewScheduleRunner(store, ScheduleRunnerOptions{IntervalSeconds: 1, Logger: logger})
	rWithLogger := runnerWithLogger.(*scheduleRunner)
	rWithLogger.logf("test message %s", "arg") // Should not panic
}

// testLoggerWriter is a simple io.Writer for testing
type testLoggerWriter struct{}

func (w *testLoggerWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func TestScheduleRunnerStartStop(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store, err := NewStore(NewStoreOptions{
		TaskDefinitionTableName: "task_definitions",
		TaskQueueTableName:      "task_queue",
		ScheduleTableName:       "schedules",
		DB:                      db,
		AutomigrateEnabled:      true,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runner := NewScheduleRunner(store, ScheduleRunnerOptions{IntervalSeconds: 1})

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
