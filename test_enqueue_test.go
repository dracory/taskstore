package taskstore

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestEnqueueByAlias(t *testing.T) {
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

	t.Log("About to enqueue")
	queued, err := store.TaskDefinitionEnqueueByAlias(ctx, "default", "test-task", map[string]any{"key": "value"})
	if err != nil {
		t.Fatalf("enqueue error: %v", err)
	}
	t.Logf("Enqueued task ID: %s", queued.GetID())
}
