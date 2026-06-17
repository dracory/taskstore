package taskstore

import (
	"context"
	"testing"
)

func Test_Store_TaskHandlerList(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatalf("TaskHandlerList: Error[%v]", err)
	}
	defer store.GetDB().Close()

	ctx := context.Background()

	// Initially empty
	handlers := store.TaskHandlerList()
	if len(handlers) != 0 {
		t.Errorf("TaskHandlerList() should return empty list initially, got %d handlers", len(handlers))
	}

	// Add a handler
	handler := newTestTaskHandler()
	err = store.TaskHandlerAdd(ctx, handler, true)
	if err != nil {
		t.Fatalf("TaskHandlerAdd: Error[%v]", err)
	}

	// Should have one handler
	handlers = store.TaskHandlerList()
	if len(handlers) != 1 {
		t.Errorf("TaskHandlerList() should return 1 handler, got %d", len(handlers))
	}

	// Add another handler
	handler2 := &testHandler2{}
	err = store.TaskHandlerAdd(ctx, handler2, true)
	if err != nil {
		t.Fatalf("TaskHandlerAdd: Error[%v]", err)
	}

	// Should have two handlers
	handlers = store.TaskHandlerList()
	if len(handlers) != 2 {
		t.Errorf("TaskHandlerList() should return 2 handlers, got %d", len(handlers))
	}
}

func Test_Store_TaskHandlerAdd(t *testing.T) {

	handler := new(testHandler)
	handler2 := new(testHandler2)

	store, err := initStore()
	if err != nil {
		t.Fatal("TaskHandlerAdd: Error in Store init: received ", "[", err, "]")
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
	TaskDefinitionHandlerBase
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

var _ TaskDefinitionHandlerInterface = (*testHandler)(nil)

type testHandler2 struct {
	TaskDefinitionHandlerBase
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

var _ TaskDefinitionHandlerInterface = (*testHandler2)(nil)
