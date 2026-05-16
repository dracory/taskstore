package admin

import (
	"log/slog"
	"net/http/httptest"
	"testing"
)

func Test_taskQueueParameters(t *testing.T) {
	// Test taskQueueParameters controller constructor with real SQLite store
	store := setupTestStore(t)
	logger := slog.Default()

	controller := taskQueueParameters(*logger, store)

	if controller == nil {
		t.Error("taskQueueParameters() should return a non-nil controller")
	}
	if controller.store == nil {
		t.Error("taskQueueParameters() should set store")
	}
}

func Test_taskQueueParametersController_struct_fields(t *testing.T) {
	// Test taskQueueParametersController struct fields
	store := setupTestStore(t)
	logger := slog.Default()

	controller := taskQueueParameters(*logger, store)

	if controller.logger == (slog.Logger{}) {
		t.Error("taskQueueParametersController logger should be set")
	}
	if controller.store == nil {
		t.Error("taskQueueParametersController store should not be nil")
	}
}

func Test_taskQueueParametersController_with_nil_store(t *testing.T) {
	// Test taskQueueParametersController with nil store
	logger := slog.Default()

	controller := taskQueueParameters(*logger, nil)

	if controller.store != nil {
		t.Error("taskQueueParametersController store should be nil")
	}
}

func Test_taskQueueParametersControllerData_request_field(t *testing.T) {
	// Test taskQueueParametersControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := taskQueueParametersControllerData{request: req}

	if data.request == nil {
		t.Error("taskQueueParametersControllerData request should not be nil")
	}
}

func Test_taskQueueParametersControllerData_nil_request(t *testing.T) {
	// Test taskQueueParametersControllerData with nil request
	data := taskQueueParametersControllerData{}

	if data.request != nil {
		t.Error("taskQueueParametersControllerData request should be nil")
	}
}
