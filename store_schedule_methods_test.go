package taskstore

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestScheduleCRUD(t *testing.T) {
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

	// Create
	schedule := NewSchedule()
	schedule.SetName("Test Schedule")
	schedule.SetDescription("Test Description")
	schedule.SetQueueName("default")
	schedule.SetTaskDefinitionID("task-1")
	schedule.SetStartAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	err = store.ScheduleCreate(ctx, schedule)
	require.NoError(t, err)

	// Read
	found, err := store.ScheduleFindByID(ctx, schedule.GetID())
	require.NoError(t, err)
	require.NotNil(t, found)
	require.Equal(t, schedule.GetName(), found.GetName())

	// Update
	found.SetName("Updated Schedule")
	err = store.ScheduleUpdate(ctx, found)
	require.NoError(t, err)

	updated, err := store.ScheduleFindByID(ctx, schedule.GetID())
	require.NoError(t, err)
	require.Equal(t, "Updated Schedule", updated.GetName())

	// List
	list, err := store.ScheduleList(ctx, NewScheduleQuery())
	require.NoError(t, err)
	require.Len(t, list, 1)

	// Delete
	err = store.ScheduleDelete(ctx, updated)
	require.NoError(t, err)

	foundAfterDelete, err := store.ScheduleFindByID(ctx, schedule.GetID())
	require.NoError(t, err)
	require.Nil(t, foundAfterDelete)
}

func TestScheduleRun(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	store, err := NewStore(NewStoreOptions{
		TaskDefinitionTableName: "task_definitions",
		TaskQueueTableName:      "task_queue",
		ScheduleTableName:       "schedules",
		DB:                      db,
		AutomigrateEnabled:      true,
	})
	assert.NoError(t, err)

	ctx := context.Background()

	// Create Task Definition
	taskDef := NewTaskDefinition()
	taskDef.SetAlias("test-task")
	err = store.TaskDefinitionCreate(ctx, taskDef)
	assert.NoError(t, err)

	// Create Schedule
	schedule := NewSchedule()
	schedule.SetName("Test Schedule")
	schedule.SetStatus("active")
	schedule.SetQueueName("default")
	schedule.SetTaskDefinitionID(taskDef.ID())
	schedule.SetStartAt(carbon.Now(carbon.UTC).AddMinutes(-1).ToDateTimeString(carbon.UTC))
	schedule.SetNextRunAt(carbon.Now(carbon.UTC).AddMinutes(-1).ToDateTimeString(carbon.UTC))
	schedule.GetRecurrenceRule().SetFrequency(FrequencyMinutely)
	schedule.GetRecurrenceRule().SetInterval(1)

	err = store.ScheduleCreate(ctx, schedule)
	assert.NoError(t, err)

	// Run Schedule
	err = store.ScheduleRun(ctx)
	assert.NoError(t, err)

	// Verify Task Enqueued
	queueList, err := store.TaskQueueList(ctx, TaskQueueQuery())
	assert.NoError(t, err)
	assert.Len(t, queueList, 1)
	assert.Equal(t, taskDef.ID(), queueList[0].TaskID())

	// Verify Schedule Updated
	updatedSchedule, err := store.ScheduleFindByID(ctx, schedule.GetID())
	assert.NoError(t, err)
	assert.NotEqual(t, sb.NULL_DATETIME, updatedSchedule.GetLastRunAt())
	assert.True(t, carbon.Parse(updatedSchedule.GetNextRunAt(), carbon.UTC).Gt(carbon.Now(carbon.UTC)))
}
