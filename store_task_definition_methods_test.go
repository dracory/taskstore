package taskstore

import (
	"context"
	"strings"
	"testing"
)

func Test_Store_TaskDefinitionCreate(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("TaskDefinitionCreate: Error[%v]", err)
	}

	task := NewTaskDefinition().
		SetAlias("TASK_ALIAS_01").
		SetTitle("TASK_TITLE_01").
		SetDescription("TASK_DESCRIPTION_01")

	query := store.SqlCreateTaskDefinitionTable()
	if strings.Contains(query, "unsupported driver") {
		t.Fatalf("TaskDefinitionCreate: UnExpected Query, received [%v]", query)
	}

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatalf("TaskDefinitionCreate: Table creation error: [%v]", err)
	}

	err = store.TaskDefinitionCreate(context.Background(), task)
	if err != nil {
		t.Fatalf("TaskDefinitionCreate: Error in Creating TaskDefinition: received [%v]", err)
	}
}
