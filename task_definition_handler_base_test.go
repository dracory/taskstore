package taskstore

import (
	"testing"
)

func Test_TaskDefinitionHandlerBase(t *testing.T) {
	handler := newTestTaskHandler()

	handler.SetOptions(map[string]string{
		"completeWithSuccess": "yes",
	})

	isOK := handler.Handle()

	if !isOK {
		t.Fatalf("Test 1: Must be successful")
	}

	if handler.SuccessMessage() != "Task forced to succeed." {
		t.Fatalf("Test 1: Message must be 'Task forced to succeed.', but found: %s", handler.SuccessMessage())
	}

	handler2 := newTestTaskHandler()

	handler2.SetOptions(map[string]string{
		"completeWithFail": "yes",
	})

	isOK = handler2.Handle()

	if isOK {
		t.Fatalf("Test 2: Must Fail")
	}

	if handler2.ErrorMessage() != "Task forced to fail." {
		t.Fatalf("Test 2: Message must be 'Task forced to fail.', but found: %s", handler2.ErrorMessage())
	}
}

func Test_TaskDefinitionHandlerBase_GetParamArray(t *testing.T) {
	handler := newTestTaskHandler()

	// Test case 1: Empty input
	handler.SetOptions(map[string]string{})
	result := handler.GetParamArray("paramArray")
	if len(result) != 0 {
		t.Errorf("Test Case 1 Failed: Expected empty array, got %v", result)
	}

	// Test case 2: Valid input (options)
	handler.SetOptions(map[string]string{"paramArray": "value1;value2;value3"})
	result = handler.GetParamArray("paramArray")
	expected := []string{"value1", "value2", "value3"}
	if len(result) != len(expected) {
		t.Errorf("Test Case 2 Failed: Expected %v, got %v", expected, result)
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Test Case 2 Failed: Expected %v, got %v", expected, result)
			break
		}
	}

	// // Test case 3: Valid input (queuedTask)
	queuedTask, err := NewTaskQueue().SetParametersMap(map[string]string{"paramArray": "value4;value5;value6"})
	if err != nil {
		t.Errorf("Test Case 3 Failed: %v", err)
	}
	handler.SetQueuedTask(queuedTask)
	result = handler.GetParamArray("paramArray")
	expected = []string{"value4", "value5", "value6"}
	if len(result) != len(expected) {
		t.Errorf("Test Case 3 Failed: Expected %v, got %v", expected, result)
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Test Case 3 Failed: Expected %v, got %v", expected, result)
			break
		}
	}

	// Test case 4: Valid input (queuedTask)
	queuedTask, err = NewTaskQueue().SetParametersMap(map[string]string{"paramArray": "singleValue"})
	if err != nil {
		t.Errorf("Test Case 4 Failed: %v", err)
	}
	handler.SetQueuedTask(queuedTask)
	result = handler.GetParamArray("paramArray")
	expected = []string{"singleValue"}
	if len(result) != len(expected) {
		t.Errorf("Test Case 4 Failed: Expected %v, got %v", expected, result)
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Test Case 4 Failed: Expected %v, got %v", expected, result)
			break
		}
	}
}

func newTestTaskHandler() *testTaskHandler {
	return &testTaskHandler{}
}

type testTaskHandler struct {
	TaskDefinitionHandlerBase
}

var _ TaskDefinitionHandlerInterface = (*testTaskHandler)(nil) // verify it extends the task interface

func (handler *testTaskHandler) Alias() string {
	return "HelloWorldTaskHandler"
}

func (handler *testTaskHandler) Title() string {
	return "Hello World"
}

func (handler *testTaskHandler) Description() string {
	return "Say hello world"
}

func (handler *testTaskHandler) Handle() bool {
	if handler.GetParam("completeWithSuccess") == "yes" {
		handler.LogSuccess("Task forced to succeed.")
		return true
	}

	if handler.GetParam("completeWithFail") == "yes" {
		handler.LogError("Task forced to fail.")
		return false
	}

	return false
}
