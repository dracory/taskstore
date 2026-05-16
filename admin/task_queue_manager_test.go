package admin

import (
	"log/slog"
	"net/http/httptest"
	"testing"
)

func Test_taskQueueManager(t *testing.T) {
	// Test taskQueueManager controller constructor with real SQLite store
	store := setupTestStore(t)
	layout := setupTestLayout(t)
	logger := slog.Default()

	controller := taskQueueManager(*logger, store, layout)

	if controller == nil {
		t.Error("taskQueueManager() should return a non-nil controller")
	}
	if controller.store == nil {
		t.Error("taskQueueManager() should set store")
	}
	if controller.layout == nil {
		t.Error("taskQueueManager() should set layout")
	}
}

func Test_taskQueueManagerController_struct_fields(t *testing.T) {
	// Test taskQueueManagerController struct fields
	store := setupTestStore(t)
	layout := setupTestLayout(t)
	logger := slog.Default()

	controller := taskQueueManager(*logger, store, layout)

	if controller.logger == (slog.Logger{}) {
		t.Error("taskQueueManagerController logger should be set")
	}
	if controller.store == nil {
		t.Error("taskQueueManagerController store should not be nil")
	}
	if controller.layout == nil {
		t.Error("taskQueueManagerController layout should not be nil")
	}
}

func Test_taskQueueManagerController_with_nil_store(t *testing.T) {
	// Test taskQueueManagerController with nil store
	layout := setupTestLayout(t)
	logger := slog.Default()

	controller := taskQueueManager(*logger, nil, layout)

	if controller.store != nil {
		t.Error("taskQueueManagerController store should be nil")
	}
}

func Test_taskQueueManagerControllerData_request_field(t *testing.T) {
	// Test taskQueueManagerControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := taskQueueManagerControllerData{request: req}

	if data.request == nil {
		t.Error("taskQueueManagerControllerData request should not be nil")
	}
}

func Test_taskQueueManagerControllerData_nil_request(t *testing.T) {
	// Test taskQueueManagerControllerData with nil request
	data := taskQueueManagerControllerData{}

	if data.request != nil {
		t.Error("taskQueueManagerControllerData request should be nil")
	}
}
