package taskstore

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestTaskQueueRunnerRunOnceProcessesQueuedTasks(t *testing.T) {
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

	handler := new(testHandler)
	err = store.TaskHandlerAdd(ctx, handler, true)
	require.NoError(t, err)

	for i := 0; i < 2; i++ {
		_, err = store.TaskDefinitionEnqueueByAlias(ctx, DefaultQueueName, handler.Alias(), map[string]any{})
		require.NoError(t, err)
	}

	queued, err := store.TaskQueueList(ctx, TaskQueueQuery().SetStatus(TaskQueueStatusQueued))
	require.NoError(t, err)
	require.Len(t, queued, 2)

	runner := NewTaskQueueRunner(store, TaskQueueRunnerOptions{IntervalSeconds: 1, UnstuckMinutes: 1, QueueName: DefaultQueueName})

	err = runner.RunOnce(ctx)
	assert.NoError(t, err)

	success, err := store.TaskQueueList(ctx, TaskQueueQuery().SetStatus(TaskQueueStatusSuccess))
	assert.NoError(t, err)
	assert.Len(t, success, 2)

	remainingQueued, err := store.TaskQueueList(ctx, TaskQueueQuery().SetStatus(TaskQueueStatusQueued))
	assert.NoError(t, err)
	assert.Len(t, remainingQueued, 0)
}

func TestTaskQueueRunnerStartStop(t *testing.T) {
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

	runner := NewTaskQueueRunner(store, TaskQueueRunnerOptions{IntervalSeconds: 1, UnstuckMinutes: 1, QueueName: DefaultQueueName})

	runner.Start(ctx)
	time.Sleep(50 * time.Millisecond)
	assert.True(t, runner.IsRunning())

	runner.Stop()
	time.Sleep(50 * time.Millisecond)
	assert.False(t, runner.IsRunning())
}
