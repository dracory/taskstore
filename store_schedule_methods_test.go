package taskstore

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dromara/carbon/v2"
	_ "modernc.org/sqlite"
)

func TestScheduleCRUD(t *testing.T) {
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

	// Create
	schedule := NewSchedule()
	schedule.SetName("Test Schedule")
	schedule.SetDescription("Test Description")
	schedule.SetQueueName("default")
	schedule.SetTaskDefinitionID("task-1")
	schedule.SetStartAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	err = store.ScheduleCreate(ctx, schedule)
	if err != nil {
		t.Fatal(err)
	}

	// Read
	found, err := store.ScheduleFindByID(ctx, schedule.GetID())
	if err != nil {
		t.Fatal(err)
	}
	if found == nil {
		t.Fatal("expected found schedule to not be nil")
	}
	if schedule.GetName() != found.GetName() {
		t.Errorf("expected name %s, got %s", schedule.GetName(), found.GetName())
	}

	// Update
	found.SetName("Updated Schedule")
	err = store.ScheduleUpdate(ctx, found)
	if err != nil {
		t.Fatal(err)
	}

	updated, err := store.ScheduleFindByID(ctx, schedule.GetID())
	if err != nil {
		t.Fatal(err)
	}
	if updated == nil {
		t.Fatal("expected updated schedule to not be nil")
	}
	if updated.GetName() != "Updated Schedule" {
		t.Errorf("expected name 'Updated Schedule', got %s", updated.GetName())
	}

	// List
	list, err := store.ScheduleList(ctx, NewScheduleQuery())
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 schedule, got %d", len(list))
	}

	// Delete
	err = store.ScheduleDelete(ctx, updated)
	if err != nil {
		t.Fatal(err)
	}

	foundAfterDelete, err := store.ScheduleFindByID(ctx, schedule.GetID())
	if err != nil {
		t.Fatal(err)
	}
	if foundAfterDelete != nil {
		t.Error("expected foundAfterDelete to be nil")
	}
}

func TestScheduleCount(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.GetDB().Close()

	ctx := context.Background()

	// Initially count should be 0
	count, err := store.ScheduleCount(ctx, NewScheduleQuery())
	if err != nil {
		t.Errorf("ScheduleCount() error = %v", err)
	}
	if count != 0 {
		t.Errorf("ScheduleCount() = %v, want 0", count)
	}

	// Create a schedule
	schedule := NewSchedule().
		SetName("Test Schedule").
		SetStatus("active").
		SetQueueName("default").
		SetTaskDefinitionID("task_def_1").
		SetRecurrenceRule(NewRecurrenceRule().SetFrequency(FrequencyDaily))

	err = store.ScheduleCreate(ctx, schedule)
	if err != nil {
		t.Fatal(err)
	}

	// Count should be 1
	count, err = store.ScheduleCount(ctx, NewScheduleQuery())
	if err != nil {
		t.Errorf("ScheduleCount() error = %v", err)
	}
	if count != 1 {
		t.Errorf("ScheduleCount() = %v, want 1", count)
	}
}

func TestScheduleDelete(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.GetDB().Close()

	ctx := context.Background()

	// Create a schedule
	schedule := NewSchedule().
		SetName("Test Schedule").
		SetStatus("active").
		SetQueueName("default").
		SetTaskDefinitionID("task_def_1").
		SetRecurrenceRule(NewRecurrenceRule().SetFrequency(FrequencyDaily))

	err = store.ScheduleCreate(ctx, schedule)
	if err != nil {
		t.Fatal(err)
	}

	// Delete the schedule
	err = store.ScheduleDelete(ctx, schedule)
	if err != nil {
		t.Errorf("ScheduleDelete() error = %v", err)
	}

	// Verify it's deleted
	found, err := store.ScheduleFindByID(ctx, schedule.GetID())
	if err != nil {
		t.Errorf("ScheduleFindByID() error = %v", err)
	}
	if found != nil {
		t.Error("Schedule should be deleted")
	}
}

func TestScheduleSoftDelete(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal(err)
	}
	defer store.GetDB().Close()

	ctx := context.Background()

	// Create a schedule
	schedule := NewSchedule().
		SetName("Test Schedule").
		SetStatus("active").
		SetQueueName("default").
		SetTaskDefinitionID("task_def_1").
		SetRecurrenceRule(NewRecurrenceRule().SetFrequency(FrequencyDaily))

	err = store.ScheduleCreate(ctx, schedule)
	if err != nil {
		t.Fatal(err)
	}

	// Soft delete the schedule
	err = store.ScheduleSoftDelete(ctx, schedule)
	if err != nil {
		t.Errorf("ScheduleSoftDelete() error = %v", err)
	}

	// Verify it's soft deleted (should not appear in normal queries)
	list, err := store.ScheduleList(ctx, NewScheduleQuery())
	if err != nil {
		t.Errorf("ScheduleList() error = %v", err)
	}
	if len(list) != 0 {
		t.Error("Soft deleted schedule should not appear in normal queries")
	}
}

func TestScheduleRun(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Error(err)
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
		t.Error(err)
	}

	ctx := context.Background()

	// Create Task Definition
	taskDef := NewTaskDefinition()
	taskDef.SetAlias("test-task")
	err = store.TaskDefinitionCreate(ctx, taskDef)
	if err != nil {
		t.Error(err)
	}

	// Create Schedule
	schedule := NewSchedule()
	schedule.SetName("Test Schedule")
	schedule.SetStatus("active")
	schedule.SetQueueName("default")
	schedule.SetTaskDefinitionID(taskDef.GetID())
	schedule.SetStartAt(carbon.Now(carbon.UTC).AddMinutes(-1).ToDateTimeString(carbon.UTC))
	schedule.SetNextRunAt(carbon.Now(carbon.UTC).AddMinutes(-1).ToDateTimeString(carbon.UTC))
	schedule.GetRecurrenceRule().SetFrequency(FrequencyMinutely)
	schedule.GetRecurrenceRule().SetInterval(1)

	err = store.ScheduleCreate(ctx, schedule)
	if err != nil {
		t.Error(err)
	}

	// Run Schedule
	err = store.ScheduleRun(ctx)
	if err != nil {
		t.Error(err)
	}

	// Verify Task Enqueued
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

	// Verify Schedule Updated
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
}
