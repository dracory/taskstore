package taskstore

import (
	"context"
	"strings"
	"testing"
)

func Test_queuePrependTaskAliasToParameters(t *testing.T) {
	tests := []struct {
		name       string
		alias      string
		parameters map[string]interface{}
		want       map[string]interface{}
	}{
		{
			name:       "prepends alias to empty parameters",
			alias:      "test_alias",
			parameters: map[string]interface{}{},
			want:       map[string]interface{}{"task_alias": "test_alias"},
		},
		{
			name:       "prepends alias to existing parameters",
			alias:      "test_alias",
			parameters: map[string]interface{}{"key1": "value1", "key2": "value2"},
			want:       map[string]interface{}{"task_alias": "test_alias", "key1": "value1", "key2": "value2"},
		},
		{
			name:       "handles complex parameters",
			alias:      "complex_alias",
			parameters: map[string]interface{}{"key1": "value1", "number": 123},
			want:       map[string]interface{}{"task_alias": "complex_alias", "key1": "value1", "number": 123},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := queuePrependTaskAliasToParameters(tt.alias, tt.parameters)
			if len(got) != len(tt.want) {
				t.Errorf("queuePrependTaskAliasToParameters() length = %v, want %v", len(got), len(tt.want))
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("queuePrependTaskAliasToParameters() key %s = %v, want %v", k, got[k], v)
				}
			}
		})
	}
}

func Test_Store_TaskDefinitionCount(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("TaskDefinitionCount: Error[%v]", err)
	}
	defer store.db.Close()

	ctx := context.Background()

	// Create table
	query, err := store.SqlCreateTaskDefinitionTable()
	if err != nil {
		t.Fatalf("SqlCreateTaskDefinitionTable: Error[%v]", err)
	}
	if _, err := store.db.Exec(query); err != nil {
		t.Fatalf("Exec: Error[%v]", err)
	}

	// Initially count should be 0
	count, err := store.TaskDefinitionCount(ctx, TaskDefinitionQuery())
	if err != nil {
		t.Errorf("TaskDefinitionCount() error = %v", err)
	}
	if count != 0 {
		t.Errorf("TaskDefinitionCount() = %v, want 0", count)
	}

	// Create a task definition
	task := NewTaskDefinition().
		SetAlias("TASK_ALIAS_01").
		SetTitle("TASK_TITLE_01").
		SetDescription("TASK_DESCRIPTION_01").
		SetStatus(TaskDefinitionStatusActive)

	err = store.TaskDefinitionCreate(ctx, task)
	if err != nil {
		t.Fatalf("TaskDefinitionCreate: Error[%v]", err)
	}

	// Count should be 1
	count, err = store.TaskDefinitionCount(ctx, TaskDefinitionQuery())
	if err != nil {
		t.Errorf("TaskDefinitionCount() error = %v", err)
	}
	if count != 1 {
		t.Errorf("TaskDefinitionCount() = %v, want 1", count)
	}
}

func Test_Store_TaskDefinitionSoftDelete(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("TaskDefinitionSoftDelete: Error[%v]", err)
	}
	defer store.db.Close()

	ctx := context.Background()

	task := NewTaskDefinition().
		SetAlias("TASK_ALIAS_SOFT_DELETE").
		SetTitle("TASK_TITLE_SOFT_DELETE").
		SetDescription("TASK_DESCRIPTION_SOFT_DELETE").
		SetStatus(TaskDefinitionStatusActive)

	err = store.TaskDefinitionCreate(ctx, task)
	if err != nil {
		t.Fatalf("TaskDefinitionCreate: Error[%v]", err)
	}

	// Soft delete the task
	err = store.TaskDefinitionSoftDelete(ctx, task)
	if err != nil {
		t.Errorf("TaskDefinitionSoftDelete() error = %v", err)
	}

	// Verify it's soft deleted (should not appear in normal queries)
	list, err := store.TaskDefinitionList(ctx, TaskDefinitionQuery())
	if err != nil {
		t.Errorf("TaskDefinitionList() error = %v", err)
	}
	if len(list) != 0 {
		t.Error("Soft deleted task should not appear in normal queries")
	}
}

func Test_Store_TaskDefinitionDelete(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("TaskDefinitionDelete: Error[%v]", err)
	}
	defer store.db.Close()

	ctx := context.Background()

	task := NewTaskDefinition().
		SetAlias("TASK_ALIAS_DELETE").
		SetTitle("TASK_TITLE_DELETE").
		SetDescription("TASK_DESCRIPTION_DELETE").
		SetStatus(TaskDefinitionStatusActive)

	err = store.TaskDefinitionCreate(ctx, task)
	if err != nil {
		t.Fatalf("TaskDefinitionCreate: Error[%v]", err)
	}

	// Delete the task
	err = store.TaskDefinitionDelete(ctx, task)
	if err != nil {
		t.Errorf("TaskDefinitionDelete() error = %v", err)
	}

	// Verify it's deleted
	found, err := store.TaskDefinitionFindByID(ctx, task.GetID())
	if err != nil {
		t.Errorf("TaskDefinitionFindByID() error = %v", err)
	}
	if found != nil {
		t.Error("Task should be deleted")
	}
}

func Test_Store_TaskDefinitionCreate(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("TaskDefinitionCreate: Error[%v]", err)
	}

	task := NewTaskDefinition().
		SetAlias("TASK_ALIAS_01").
		SetTitle("TASK_TITLE_01").
		SetDescription("TASK_DESCRIPTION_01")

	query, err := store.SqlCreateTaskDefinitionTable()
	if err != nil {
		t.Fatalf("TaskDefinitionCreate: Error[%v]", err)
	}
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
