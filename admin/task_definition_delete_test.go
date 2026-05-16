package admin

import (
	"log/slog"
	"net/http/httptest"
	"testing"
)

func Test_taskDefinitionDelete(t *testing.T) {
	// Test taskDefinitionDelete controller constructor with real SQLite store
	store := setupTestStore(t)
	logger := slog.Default()

	controller := taskDefinitionDelete(*logger, store)

	if controller == nil {
		t.Error("taskDefinitionDelete() should return a non-nil controller")
	}
	if controller.store == nil {
		t.Error("taskDefinitionDelete() should set store")
	}
}

func Test_taskDefinitionDeleteController_struct_fields(t *testing.T) {
	// Test taskDefinitionDeleteController struct fields
	store := setupTestStore(t)
	logger := slog.Default()

	controller := taskDefinitionDelete(*logger, store)

	if controller.logger == (slog.Logger{}) {
		t.Error("taskDefinitionDeleteController logger should be set")
	}
	if controller.store == nil {
		t.Error("taskDefinitionDeleteController store should not be nil")
	}
}

func Test_taskDefinitionDeleteController_with_nil_store(t *testing.T) {
	// Test taskDefinitionDeleteController with nil store
	logger := slog.Default()

	controller := taskDefinitionDelete(*logger, nil)

	if controller.store != nil {
		t.Error("taskDefinitionDeleteController store should be nil")
	}
}

func Test_taskDefinitionDeleteControllerData_request_field(t *testing.T) {
	// Test taskDefinitionDeleteControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := taskDefinitionDeleteControllerData{request: req}

	if data.request == nil {
		t.Error("taskDefinitionDeleteControllerData request should not be nil")
	}
}

func Test_taskDefinitionDeleteControllerData_nil_request(t *testing.T) {
	// Test taskDefinitionDeleteControllerData with nil request
	data := taskDefinitionDeleteControllerData{}

	if data.request != nil {
		t.Error("taskDefinitionDeleteControllerData request should be nil")
	}
}
