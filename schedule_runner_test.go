package taskstore

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestScheduleRunnerRunOnceEnqueuesDueSchedule(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	store, err := NewStore(NewStoreOptions{
		TaskDefinitionTableName: "task_definitions",
		TaskQueueTableName:      "task_queue",
		ScheduleTableName:       "schedules",
		DB:                      db,
		AutomigrateEnabled:      true,
	})
	require.NoError(t, err)

	ctx := context.Background()

	taskDef := NewTaskDefinition()
	taskDef.SetAlias("test-task")
	err = store.TaskDefinitionCreate(ctx, taskDef)
	require.NoError(t, err)

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
	require.NoError(t, err)

	runner := NewScheduleRunner(store, ScheduleRunnerOptions{IntervalSeconds: 1})

	err = runner.RunOnce(ctx)
	assert.NoError(t, err)

	queueList, err := store.TaskQueueList(ctx, TaskQueueQuery())
	assert.NoError(t, err)
	assert.Len(t, queueList, 1)
	assert.Equal(t, taskDef.ID(), queueList[0].TaskID())

	updatedSchedule, err := store.ScheduleFindByID(ctx, schedule.GetID())
	assert.NoError(t, err)
	assert.NotEqual(t, sb.NULL_DATETIME, updatedSchedule.GetLastRunAt())
	assert.True(t, carbon.Parse(updatedSchedule.GetNextRunAt(), carbon.UTC).Gt(carbon.Now(carbon.UTC)))
	assert.Equal(t, 1, updatedSchedule.GetExecutionCount())
	assert.Equal(t, "active", updatedSchedule.GetStatus())
}

func TestScheduleRunnerSetInitialRuns(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	store, err := NewStore(NewStoreOptions{
		TaskDefinitionTableName: "task_definitions",
		TaskQueueTableName:      "task_queue",
		ScheduleTableName:       "schedules",
		DB:                      db,
		AutomigrateEnabled:      true,
	})
	require.NoError(t, err)

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
	require.NoError(t, err)

	runner := NewScheduleRunner(store, ScheduleRunnerOptions{IntervalSeconds: 1})

	err = runner.SetInitialRuns(ctx)
	assert.NoError(t, err)

	updatedSchedule, err := store.ScheduleFindByID(ctx, schedule.GetID())
	assert.NoError(t, err)
	assert.NotEqual(t, sb.NULL_DATETIME, updatedSchedule.GetNextRunAt())
	assert.Equal(t, "active", updatedSchedule.GetStatus())
}

func TestScheduleRunnerStartStop(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	store, err := NewStore(NewStoreOptions{
		TaskDefinitionTableName: "task_definitions",
		TaskQueueTableName:      "task_queue",
		ScheduleTableName:       "schedules",
		DB:                      db,
		AutomigrateEnabled:      true,
	})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runner := NewScheduleRunner(store, ScheduleRunnerOptions{IntervalSeconds: 1})

	runner.Start(ctx)
	time.Sleep(50 * time.Millisecond)
	assert.True(t, runner.IsRunning())

	runner.Stop()
	time.Sleep(50 * time.Millisecond)
	assert.False(t, runner.IsRunning())
}
