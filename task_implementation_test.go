package taskstore

import (
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

func TestNewTask(t *testing.T) {
	task := NewTask()

	if task == nil {
		t.Fatal("NewTask: Expected task to be created, got nil")
	}

	if task.ID() == "" {
		t.Error("NewTask: Expected ID to be set")
	}

	if task.Status() != TaskStatusActive {
		t.Errorf("NewTask: Expected status to be %s, got %s", TaskStatusActive, task.Status())
	}

	if task.Memo() != "" {
		t.Errorf("NewTask: Expected memo to be empty, got %s", task.Memo())
	}

	if task.CreatedAt() == "" {
		t.Error("NewTask: Expected CreatedAt to be set")
	}

	if task.UpdatedAt() == "" {
		t.Error("NewTask: Expected UpdatedAt to be set")
	}

	if task.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Errorf("NewTask: Expected SoftDeletedAt to be %s, got %s", sb.MAX_DATETIME, task.SoftDeletedAt())
	}
}

func TestNewTaskFromExistingData(t *testing.T) {
	data := map[string]string{
		COLUMN_ID:          "test-id",
		COLUMN_ALIAS:       "test-alias",
		COLUMN_TITLE:       "Test Title",
		COLUMN_DESCRIPTION: "Test Description",
		COLUMN_STATUS:      TaskStatusCanceled,
		COLUMN_MEMO:        "Test Memo",
		COLUMN_CREATED_AT:  "2023-01-01 12:00:00",
		COLUMN_UPDATED_AT:  "2023-01-02 12:00:00",
		COLUMN_DELETED_AT:  "2023-01-03 12:00:00",
	}

	task := NewTaskFromExistingData(data)

	if task.ID() != "test-id" {
		t.Errorf("NewTaskFromExistingData: Expected ID to be 'test-id', got %s", task.ID())
	}

	if task.Alias() != "test-alias" {
		t.Errorf("NewTaskFromExistingData: Expected Alias to be 'test-alias', got %s", task.Alias())
	}

	if task.Title() != "Test Title" {
		t.Errorf("NewTaskFromExistingData: Expected Title to be 'Test Title', got %s", task.Title())
	}

	if task.Description() != "Test Description" {
		t.Errorf("NewTaskFromExistingData: Expected Description to be 'Test Description', got %s", task.Description())
	}

	if task.Status() != TaskStatusCanceled {
		t.Errorf("NewTaskFromExistingData: Expected Status to be %s, got %s", TaskStatusCanceled, task.Status())
	}

	if task.Memo() != "Test Memo" {
		t.Errorf("NewTaskFromExistingData: Expected Memo to be 'Test Memo', got %s", task.Memo())
	}
}

func TestTask_IsActive(t *testing.T) {
	task := NewTask()

	// Test active status
	task.SetStatus(TaskStatusActive)
	if !task.IsActive() {
		t.Error("IsActive: Expected task to be active when status is TaskStatusActive")
	}

	// Test non-active status
	task.SetStatus(TaskStatusCanceled)
	if task.IsActive() {
		t.Error("IsActive: Expected task to not be active when status is TaskStatusCanceled")
	}
}

func TestTask_IsCanceled(t *testing.T) {
	task := NewTask()

	// Test canceled status
	task.SetStatus(TaskStatusCanceled)
	if !task.IsCanceled() {
		t.Error("IsCanceled: Expected task to be canceled when status is TaskStatusCanceled")
	}

	// Test non-canceled status
	task.SetStatus(TaskStatusActive)
	if task.IsCanceled() {
		t.Error("IsCanceled: Expected task to not be canceled when status is TaskStatusActive")
	}
}

func TestTask_IsSoftDeleted(t *testing.T) {
	task := NewTask()

	// Test not soft deleted (default state)
	if task.IsSoftDeleted() {
		t.Error("IsSoftDeleted: Expected new task to not be soft deleted")
	}

	// Test soft deleted
	pastTime := carbon.Now(carbon.UTC).SubHours(1).ToDateTimeString(carbon.UTC)
	task.SetSoftDeletedAt(pastTime)
	if !task.IsSoftDeleted() {
		t.Error("IsSoftDeleted: Expected task to be soft deleted when deleted_at is in the past")
	}

	// Test future deletion time (not yet deleted)
	futureTime := carbon.Now(carbon.UTC).AddHours(1).ToDateTimeString(carbon.UTC)
	task.SetSoftDeletedAt(futureTime)
	if task.IsSoftDeleted() {
		t.Error("IsSoftDeleted: Expected task to not be soft deleted when deleted_at is in the future")
	}
}

func TestTask_CreatedAtCarbon(t *testing.T) {
	task := NewTask()
	createdAtStr := "2023-01-01 12:00:00"
	task.SetCreatedAt(createdAtStr)

	createdAtCarbon := task.CreatedAtCarbon()
	if createdAtCarbon == nil {
		t.Fatal("CreatedAtCarbon: Expected carbon instance, got nil")
	}

	if createdAtCarbon.ToDateTimeString(carbon.UTC) != createdAtStr {
		t.Errorf("CreatedAtCarbon: Expected %s, got %s", createdAtStr, createdAtCarbon.ToDateTimeString(carbon.UTC))
	}
}

func TestTask_UpdatedAtCarbon(t *testing.T) {
	task := NewTask()
	updatedAtStr := "2023-01-02 15:30:45"
	task.SetUpdatedAt(updatedAtStr)

	updatedAtCarbon := task.UpdatedAtCarbon()
	if updatedAtCarbon == nil {
		t.Fatal("UpdatedAtCarbon: Expected carbon instance, got nil")
	}

	if updatedAtCarbon.ToDateTimeString(carbon.UTC) != updatedAtStr {
		t.Errorf("UpdatedAtCarbon: Expected %s, got %s", updatedAtStr, updatedAtCarbon.ToDateTimeString(carbon.UTC))
	}
}

func TestTask_SoftDeletedAtCarbon(t *testing.T) {
	task := NewTask()
	deletedAtStr := "2023-01-03 09:15:30"
	task.SetSoftDeletedAt(deletedAtStr)

	deletedAtCarbon := task.SoftDeletedAtCarbon()
	if deletedAtCarbon == nil {
		t.Fatal("SoftDeletedAtCarbon: Expected carbon instance, got nil")
	}

	if deletedAtCarbon.ToDateTimeString(carbon.UTC) != deletedAtStr {
		t.Errorf("SoftDeletedAtCarbon: Expected %s, got %s", deletedAtStr, deletedAtCarbon.ToDateTimeString(carbon.UTC))
	}
}

func TestTask_SettersAndGetters(t *testing.T) {
	task := NewTask()

	// Test ID
	testID := "test-task-id"
	task.SetID(testID)
	if task.ID() != testID {
		t.Errorf("ID: Expected %s, got %s", testID, task.ID())
	}

	// Test Alias
	testAlias := "test-alias"
	task.SetAlias(testAlias)
	if task.Alias() != testAlias {
		t.Errorf("Alias: Expected %s, got %s", testAlias, task.Alias())
	}

	// Test Title
	testTitle := "Test Task Title"
	task.SetTitle(testTitle)
	if task.Title() != testTitle {
		t.Errorf("Title: Expected %s, got %s", testTitle, task.Title())
	}

	// Test Description
	testDescription := "Test task description"
	task.SetDescription(testDescription)
	if task.Description() != testDescription {
		t.Errorf("Description: Expected %s, got %s", testDescription, task.Description())
	}

	// Test Memo
	testMemo := "Test memo"
	task.SetMemo(testMemo)
	if task.Memo() != testMemo {
		t.Errorf("Memo: Expected %s, got %s", testMemo, task.Memo())
	}

	// Test Status
	task.SetStatus(TaskStatusCanceled)
	if task.Status() != TaskStatusCanceled {
		t.Errorf("Status: Expected %s, got %s", TaskStatusCanceled, task.Status())
	}

	// Test CreatedAt
	testCreatedAt := "2023-01-01 10:00:00"
	task.SetCreatedAt(testCreatedAt)
	if task.CreatedAt() != testCreatedAt {
		t.Errorf("CreatedAt: Expected %s, got %s", testCreatedAt, task.CreatedAt())
	}

	// Test UpdatedAt
	testUpdatedAt := "2023-01-02 11:00:00"
	task.SetUpdatedAt(testUpdatedAt)
	if task.UpdatedAt() != testUpdatedAt {
		t.Errorf("UpdatedAt: Expected %s, got %s", testUpdatedAt, task.UpdatedAt())
	}

	// Test SoftDeletedAt
	testDeletedAt := "2023-01-03 12:00:00"
	task.SetSoftDeletedAt(testDeletedAt)
	if task.SoftDeletedAt() != testDeletedAt {
		t.Errorf("SoftDeletedAt: Expected %s, got %s", testDeletedAt, task.SoftDeletedAt())
	}
}

func TestTask_ChainedSetters(t *testing.T) {
	task := NewTask()

	// Test that setters return the task instance for chaining
	result := task.SetID("test-id").
		SetAlias("test-alias").
		SetTitle("Test Title").
		SetDescription("Test Description").
		SetMemo("Test Memo").
		SetStatus(TaskStatusCanceled).
		SetCreatedAt("2023-01-01 10:00:00").
		SetUpdatedAt("2023-01-02 11:00:00").
		SetSoftDeletedAt("2023-01-03 12:00:00")

	if result != task {
		t.Error("ChainedSetters: Expected setters to return the same task instance for chaining")
	}

	// Verify all values were set correctly
	if task.ID() != "test-id" {
		t.Error("ChainedSetters: ID not set correctly")
	}
	if task.Alias() != "test-alias" {
		t.Error("ChainedSetters: Alias not set correctly")
	}
	if task.Title() != "Test Title" {
		t.Error("ChainedSetters: Title not set correctly")
	}
	if task.Description() != "Test Description" {
		t.Error("ChainedSetters: Description not set correctly")
	}
	if task.Memo() != "Test Memo" {
		t.Error("ChainedSetters: Memo not set correctly")
	}
	if task.Status() != TaskStatusCanceled {
		t.Error("ChainedSetters: Status not set correctly")
	}
}
