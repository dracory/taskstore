package admin

import (
	"log/slog"
	"net/http/httptest"
	"testing"
)

func Test_taskDefinitionUpdate(t *testing.T) {
	// Test taskDefinitionUpdate controller constructor with real SQLite store
	store := setupTestStore(t)
	logger := slog.Default()

	controller := taskDefinitionUpdate(*logger, store)

	if controller == nil {
		t.Error("taskDefinitionUpdate() should return a non-nil controller")
	}
	if controller.store == nil {
		t.Error("taskDefinitionUpdate() should set store")
	}
}

func Test_taskDefinitionUpdateController_struct_fields(t *testing.T) {
	// Test taskDefinitionUpdateController struct fields
	store := setupTestStore(t)
	logger := slog.Default()

	controller := taskDefinitionUpdate(*logger, store)

	if controller.logger == (slog.Logger{}) {
		t.Error("taskDefinitionUpdateController logger should be set")
	}
	if controller.store == nil {
		t.Error("taskDefinitionUpdateController store should not be nil")
	}
}

func Test_taskDefinitionUpdateController_with_nil_store(t *testing.T) {
	// Test taskDefinitionUpdateController with nil store
	logger := slog.Default()

	controller := taskDefinitionUpdate(*logger, nil)

	if controller.store != nil {
		t.Error("taskDefinitionUpdateController store should be nil")
	}
}

func Test_taskDefinitionUpdateControllerData_request_field(t *testing.T) {
	// Test taskDefinitionUpdateControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := taskDefinitionUpdateControllerData{request: req}

	if data.request == nil {
		t.Error("taskDefinitionUpdateControllerData request should not be nil")
	}
}

func Test_taskDefinitionUpdateControllerData_nil_request(t *testing.T) {
	// Test taskDefinitionUpdateControllerData with nil request
	data := taskDefinitionUpdateControllerData{}

	if data.request != nil {
		t.Error("taskDefinitionUpdateControllerData request should be nil")
	}
}
