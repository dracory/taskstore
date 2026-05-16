package admin

import (
	"log/slog"
	"net/http/httptest"
	"testing"
)

func Test_taskQueueDelete(t *testing.T) {
	// Test taskQueueDelete controller constructor with real SQLite store
	store := setupTestStore(t)
	logger := slog.Default()

	controller := taskQueueDelete(*logger, store)

	if controller == nil {
		t.Error("taskQueueDelete() should return a non-nil controller")
	}
	if controller.store == nil {
		t.Error("taskQueueDelete() should set store")
	}
}

func Test_taskQueueDeleteController_struct_fields(t *testing.T) {
	// Test taskQueueDeleteController struct fields
	store := setupTestStore(t)
	logger := slog.Default()

	controller := taskQueueDelete(*logger, store)

	if controller.logger == (slog.Logger{}) {
		t.Error("taskQueueDeleteController logger should be set")
	}
	if controller.store == nil {
		t.Error("taskQueueDeleteController store should not be nil")
	}
}

func Test_taskQueueDeleteController_with_nil_store(t *testing.T) {
	// Test taskQueueDeleteController with nil store
	logger := slog.Default()

	controller := taskQueueDelete(*logger, nil)

	if controller.store != nil {
		t.Error("taskQueueDeleteController store should be nil")
	}
}

func Test_taskQueueDeleteControllerData_request_field(t *testing.T) {
	// Test taskQueueDeleteControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := taskQueueDeleteControllerData{request: req}

	if data.request == nil {
		t.Error("taskQueueDeleteControllerData request should not be nil")
	}
}

func Test_taskQueueDeleteControllerData_nil_request(t *testing.T) {
	// Test taskQueueDeleteControllerData with nil request
	data := taskQueueDeleteControllerData{}

	if data.request != nil {
		t.Error("taskQueueDeleteControllerData request should be nil")
	}
}
