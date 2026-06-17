package taskstore

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dromara/carbon/v2"
	_ "modernc.org/sqlite"
)

func TestRunOnce(t *testing.T) {
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

	t.Log("About to run once")
	t.Log("About to find active schedules")
	schedules, err := store.ScheduleList(ctx, NewScheduleQuery().SetStatus("active"))
	if err != nil {
		t.Fatalf("find active schedules error: %v", err)
	}
	t.Logf("Found %d active schedules", len(schedules))
	for _, s := range schedules {
		t.Logf("Schedule: %s, nextRunAt=%s, isDue=%v", s.GetName(), s.GetNextRunAt(), s.IsDue())
	}

	t.Log("About to call RunOnce")
	err = runner.RunOnce(ctx)
	if err != nil {
		t.Fatalf("run once error: %v", err)
	}
	t.Log("Run once completed")
}
