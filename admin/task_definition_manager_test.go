package admin

import (
	"log/slog"
	"net/http/httptest"
	"testing"
)

func Test_taskDefinitionManager(t *testing.T) {
	// Test taskDefinitionManager controller constructor with real SQLite store
	store := setupTestStore(t)
	layout := setupTestLayout(t)
	logger := slog.Default()

	controller := taskDefinitionManager(*logger, store, layout)

	if controller == nil {
		t.Error("taskDefinitionManager() should return a non-nil controller")
	}
	if controller.store == nil {
		t.Error("taskDefinitionManager() should set store")
	}
	if controller.layout == nil {
		t.Error("taskDefinitionManager() should set layout")
	}
}

func Test_taskDefinitionManagerController_struct_fields(t *testing.T) {
	// Test taskDefinitionManagerController struct fields
	store := setupTestStore(t)
	layout := setupTestLayout(t)
	logger := slog.Default()

	controller := taskDefinitionManager(*logger, store, layout)

	if controller.logger == (slog.Logger{}) {
		t.Error("taskDefinitionManagerController logger should be set")
	}
	if controller.store == nil {
		t.Error("taskDefinitionManagerController store should not be nil")
	}
	if controller.layout == nil {
		t.Error("taskDefinitionManagerController layout should not be nil")
	}
}

func Test_taskDefinitionManagerController_with_nil_store(t *testing.T) {
	// Test taskDefinitionManagerController with nil store
	layout := setupTestLayout(t)
	logger := slog.Default()

	controller := taskDefinitionManager(*logger, nil, layout)

	if controller.store != nil {
		t.Error("taskDefinitionManagerController store should be nil")
	}
}

func Test_taskDefinitionManagerControllerData_request_field(t *testing.T) {
	// Test taskDefinitionManagerControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := taskDefinitionManagerControllerData{request: req}

	if data.request == nil {
		t.Error("taskDefinitionManagerControllerData request should not be nil")
	}
}

func Test_taskDefinitionManagerControllerData_nil_request(t *testing.T) {
	// Test taskDefinitionManagerControllerData with nil request
	data := taskDefinitionManagerControllerData{}

	if data.request != nil {
		t.Error("taskDefinitionManagerControllerData request should be nil")
	}
}
