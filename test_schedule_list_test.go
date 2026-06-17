package taskstore

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestScheduleList(t *testing.T) {
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

	t.Log("About to list schedules")
	schedules, err := store.ScheduleList(ctx, NewScheduleQuery().SetStatus("active"))
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	t.Logf("Found %d schedules", len(schedules))
}
