package taskstore

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dracory/sb"
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
	schedule.SetTaskDefinitionID(taskDef.ID())

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
	if taskDef.ID() != queueList[0].TaskID() {
		t.Errorf("expected task ID %s, got %s", taskDef.ID(), queueList[0].TaskID())
	}

	updatedSchedule, err := store.ScheduleFindByID(ctx, schedule.GetID())
	if err != nil {
		t.Error(err)
	}
	if updatedSchedule == nil {
		t.Fatal("expected updatedSchedule to not be nil")
	}
	if updatedSchedule.GetLastRunAt() == sb.NULL_DATETIME {
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
	schedule.SetNextRunAt(sb.NULL_DATETIME)
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
	if updatedSchedule.GetNextRunAt() == sb.NULL_DATETIME {
		t.Error("expected NextRunAt to not be NULL_DATETIME")
	}
	if updatedSchedule.GetStatus() != "active" {
		t.Errorf("expected status 'active', got %s", updatedSchedule.GetStatus())
	}
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
