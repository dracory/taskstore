package admin

import (
	"log/slog"
	"net/http/httptest"
	"testing"
)

func Test_taskQueueRequeue(t *testing.T) {
	// Test taskQueueRequeue controller constructor with real SQLite store
	store := setupTestStore(t)
	logger := slog.Default()

	controller := taskQueueRequeue(*logger, store)

	if controller == nil {
		t.Error("taskQueueRequeue() should return a non-nil controller")
	}
	if controller.store == nil {
		t.Error("taskQueueRequeue() should set store")
	}
}

func Test_taskQueueRequeueController_struct_fields(t *testing.T) {
	// Test taskQueueRequeueController struct fields
	store := setupTestStore(t)
	logger := slog.Default()

	controller := taskQueueRequeue(*logger, store)

	if controller.logger == (slog.Logger{}) {
		t.Error("taskQueueRequeueController logger should be set")
	}
	if controller.store == nil {
		t.Error("taskQueueRequeueController store should not be nil")
	}
}

func Test_taskQueueRequeueController_with_nil_store(t *testing.T) {
	// Test taskQueueRequeueController with nil store
	logger := slog.Default()

	controller := taskQueueRequeue(*logger, nil)

	if controller.store != nil {
		t.Error("taskQueueRequeueController store should be nil")
	}
}

func Test_taskQueueRequeueControllerData_request_field(t *testing.T) {
	// Test taskQueueRequeueControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := taskQueueRequeueControllerData{request: req}

	if data.request == nil {
		t.Error("taskQueueRequeueControllerData request should not be nil")
	}
}

func Test_taskQueueRequeueControllerData_nil_request(t *testing.T) {
	// Test taskQueueRequeueControllerData with nil request
	data := taskQueueRequeueControllerData{}

	if data.request != nil {
		t.Error("taskQueueRequeueControllerData request should be nil")
	}
}
