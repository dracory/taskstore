package taskstore

import (
	"context"
	"testing"
)

func Test_Store_TaskHandlerAdd(t *testing.T) {

	handler := new(testHandler)
	handler2 := new(testHandler2)

	store, err := initStore()
	if err != nil {
		t.Fatal("TaskHandlerAdd: Error in Store init: received ", "[", err, "]")
	}

	query := store.SqlCreateTaskDefinitionTable()

	_, err = store.db.Exec(query)
	if err != nil {
		t.Fatal("TaskHandlerAdd: Table creation error: ", "[", err, "]")
	}

	err = store.TaskHandlerAdd(context.Background(), handler, true)
	if err != nil {
		t.Fatal("TaskHandlerAdd: Error in adding handler: received ", "[", err, "]")
	}

	err = store.TaskHandlerAdd(context.Background(), handler, true)
	if err != nil {
		t.Fatal("TaskHandlerAdd: Error in adding handler: received ", "[", err, "]")
	}

	tasksNumber, err := store.TaskDefinitionCount(context.Background(), TaskDefinitionQuery())

	if err != nil {
		t.Fatal("TaskHandlerAdd: Error in counting tasks: received ", "[", err, "]")
	}

	if tasksNumber != 1 {
		t.Fatal("TaskHandlerAdd: Error in counting tasks: expected ", "[", 1, "], received ", "[", tasksNumber, "]")
	}

	err = store.TaskHandlerAdd(context.Background(), handler2, true)
	if err != nil {
		t.Fatal("TaskHandlerAdd: Error in adding handler: received ", "[", err, "]")
	}

	tasksNumber, err = store.TaskDefinitionCount(context.Background(), TaskDefinitionQuery())

	if err != nil {
		t.Fatal("TaskHandlerAdd: Error in counting tasks: received ", "[", err, "]")
	}

	if tasksNumber != 2 {
		t.Fatal("TaskHandlerAdd: Error in counting tasks: expected ", "[", 2, "], received ", "[", tasksNumber, "]")
	}

}

type testHandler struct {
	TaskHandlerBase
}

func (h *testHandler) Alias() string {
	return "TestHandlerAlias"
}

func (h *testHandler) Title() string {
	return "Test Handler Title"
}

func (h *testHandler) Description() string {
	return "Test Handler Description"
}

func (h *testHandler) Handle() bool {
	return true
}

var _ TaskHandlerInterface = (*testHandler)(nil)

type testHandler2 struct {
	TaskHandlerBase
}

func (h *testHandler2) Alias() string {
	return "TestHandlerAlias2"
}

func (h *testHandler2) Title() string {
	return "Test Handler Title 2"
}

func (h *testHandler2) Description() string {
	return "Test Handler Description 2"
}

func (h *testHandler2) Handle() bool {
	return true
}

var _ TaskHandlerInterface = (*testHandler2)(nil)
