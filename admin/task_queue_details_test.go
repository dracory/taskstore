package admin

import (
	"log/slog"
	"net/http/httptest"
	"testing"
)

func Test_taskQueueDetails(t *testing.T) {
	// Test taskQueueDetails controller constructor with real SQLite store
	store := setupTestStore(t)
	logger := slog.Default()

	controller := taskQueueDetails(*logger, store)

	if controller == nil {
		t.Error("taskQueueDetails() should return a non-nil controller")
	}
	if controller.store == nil {
		t.Error("taskQueueDetails() should set store")
	}
}

func Test_taskQueueDetailsController_struct_fields(t *testing.T) {
	// Test taskQueueDetailsController struct fields
	store := setupTestStore(t)
	logger := slog.Default()

	controller := taskQueueDetails(*logger, store)

	if controller.logger == (slog.Logger{}) {
		t.Error("taskQueueDetailsController logger should be set")
	}
	if controller.store == nil {
		t.Error("taskQueueDetailsController store should not be nil")
	}
}

func Test_taskQueueDetailsController_with_nil_store(t *testing.T) {
	// Test taskQueueDetailsController with nil store
	logger := slog.Default()

	controller := taskQueueDetails(*logger, nil)

	if controller.store != nil {
		t.Error("taskQueueDetailsController store should be nil")
	}
}

func Test_taskQueueDetailsControllerData_request_field(t *testing.T) {
	// Test taskQueueDetailsControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := taskQueueDetailsControllerData{request: req}

	if data.request == nil {
		t.Error("taskQueueDetailsControllerData request should not be nil")
	}
}

func Test_taskQueueDetailsControllerData_nil_request(t *testing.T) {
	// Test taskQueueDetailsControllerData with nil request
	data := taskQueueDetailsControllerData{}

	if data.request != nil {
		t.Error("taskQueueDetailsControllerData request should be nil")
	}
}
