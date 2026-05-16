package admin

import (
	"log/slog"
	"net/http/httptest"
	"testing"
)

func Test_taskQueueTaskRestart(t *testing.T) {
	// Test taskQueueTaskRestart controller constructor with real SQLite store
	store := setupTestStore(t)
	logger := slog.Default()

	controller := taskQueueTaskRestart(*logger, store)

	if controller == nil {
		t.Error("taskQueueTaskRestart() should return a non-nil controller")
	}
	if controller.store == nil {
		t.Error("taskQueueTaskRestart() should set store")
	}
}

func Test_taskQueueTaskRestartController_struct_fields(t *testing.T) {
	// Test taskQueueTaskRestartController struct fields
	store := setupTestStore(t)
	logger := slog.Default()

	controller := taskQueueTaskRestart(*logger, store)

	if controller.logger == (slog.Logger{}) {
		t.Error("taskQueueTaskRestartController logger should be set")
	}
	if controller.store == nil {
		t.Error("taskQueueTaskRestartController store should not be nil")
	}
}

func Test_taskQueueTaskRestartController_with_nil_store(t *testing.T) {
	// Test taskQueueTaskRestartController with nil store
	logger := slog.Default()

	controller := taskQueueTaskRestart(*logger, nil)

	if controller.store != nil {
		t.Error("taskQueueTaskRestartController store should be nil")
	}
}

func Test_taskQueueTaskRestartControllerData_request_field(t *testing.T) {
	// Test taskQueueTaskRestartControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := taskQueueTaskRestartControllerData{request: req}

	if data.request == nil {
		t.Error("taskQueueTaskRestartControllerData request should not be nil")
	}
}

func Test_taskQueueTaskRestartControllerData_nil_request(t *testing.T) {
	// Test taskQueueTaskRestartControllerData with nil request
	data := taskQueueTaskRestartControllerData{}

	if data.request != nil {
		t.Error("taskQueueTaskRestartControllerData request should be nil")
	}
}
