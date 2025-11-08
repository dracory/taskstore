package taskstore

import (
	"strings"
	"testing"
)

func Test_Store_TaskCreate(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("QueueCreate: Error[%v]", err)
	}

	task := NewTask().
		SetAlias("TASK_ALIAS_01").
		SetTitle("TASK_TITLE_01").
		SetDescription("TASK_DESCRIPTION_01")

	query := store.SqlCreateTaskTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("QueueCreate: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("QueueCreate: Table creation error: [%v]", err)
	}

	err = store.TaskCreate(task)
	if err != nil {
		t.Fatalf("QueueCreate: Error in Creating Queue: received [%v]", err)
	}
}
