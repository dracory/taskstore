package admin

import (
	"log/slog"
	"net/http/httptest"
	"testing"
)

func Test_taskQueueCreate(t *testing.T) {
	// Test taskQueueCreate controller constructor with real SQLite store
	store := setupTestStore(t)
	logger := slog.Default()

	controller := taskQueueCreate(*logger, store)

	if controller == nil {
		t.Error("taskQueueCreate() should return a non-nil controller")
	}
	if controller.store == nil {
		t.Error("taskQueueCreate() should set store")
	}
}

func Test_taskQueueCreateController_struct_fields(t *testing.T) {
	// Test taskQueueCreateController struct fields
	store := setupTestStore(t)
	logger := slog.Default()

	controller := taskQueueCreate(*logger, store)

	if controller.logger == (slog.Logger{}) {
		t.Error("taskQueueCreateController logger should be set")
	}
	if controller.store == nil {
		t.Error("taskQueueCreateController store should not be nil")
	}
}

func Test_taskQueueCreateController_with_nil_store(t *testing.T) {
	// Test taskQueueCreateController with nil store
	logger := slog.Default()

	controller := taskQueueCreate(*logger, nil)

	if controller.store != nil {
		t.Error("taskQueueCreateController store should be nil")
	}
}

func Test_taskQueueCreateControllerData_request_field(t *testing.T) {
	// Test taskQueueCreateControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := taskQueueCreateControllerData{request: req}

	if data.request == nil {
		t.Error("taskQueueCreateControllerData request should not be nil")
	}
}

func Test_taskQueueCreateControllerData_nil_request(t *testing.T) {
	// Test taskQueueCreateControllerData with nil request
	data := taskQueueCreateControllerData{}

	if data.request != nil {
		t.Error("taskQueueCreateControllerData request should be nil")
	}
}

func Test_taskQueueCreateController_with_real_store(t *testing.T) {
	// Test taskQueueCreateController with real SQLite store
	store := setupTestStore(t)
	logger := slog.Default()

	controller := taskQueueCreate(*logger, store)

	// Verify controller has access to store
	if controller.store == nil {
		t.Error("taskQueueCreateController store should not be nil")
	}
}
