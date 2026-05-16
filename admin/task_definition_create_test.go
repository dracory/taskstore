package admin

import (
	"log/slog"
	"net/http/httptest"
	"testing"
)

func Test_taskDefinitionCreate(t *testing.T) {
	// Test taskDefinitionCreate controller constructor with real SQLite store
	store := setupTestStore(t)
	logger := slog.Default()

	controller := taskDefinitionCreate(*logger, store)

	if controller == nil {
		t.Error("taskDefinitionCreate() should return a non-nil controller")
	}
	if controller.store == nil {
		t.Error("taskDefinitionCreate() should set store")
	}
}

func Test_taskDefinitionCreateController_struct_fields(t *testing.T) {
	// Test taskDefinitionCreateController struct fields
	store := setupTestStore(t)
	logger := slog.Default()

	controller := taskDefinitionCreate(*logger, store)

	if controller.logger == (slog.Logger{}) {
		t.Error("taskDefinitionCreateController logger should be set")
	}
	if controller.store == nil {
		t.Error("taskDefinitionCreateController store should not be nil")
	}
}

func Test_taskDefinitionCreateController_with_nil_store(t *testing.T) {
	// Test taskDefinitionCreateController with nil store
	logger := slog.Default()

	controller := taskDefinitionCreate(*logger, nil)

	if controller.store != nil {
		t.Error("taskDefinitionCreateController store should be nil")
	}
}

func Test_taskDefinitionCreateControllerData_request_field(t *testing.T) {
	// Test taskDefinitionCreateControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := taskDefinitionCreateControllerData{request: req}

	if data.request == nil {
		t.Error("taskDefinitionCreateControllerData request should not be nil")
	}
}

func Test_taskDefinitionCreateControllerData_nil_request(t *testing.T) {
	// Test taskDefinitionCreateControllerData with nil request
	data := taskDefinitionCreateControllerData{}

	if data.request != nil {
		t.Error("taskDefinitionCreateControllerData request should be nil")
	}
}
