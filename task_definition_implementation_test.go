package taskstore

import (
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

func TestNewTaskDefinition(t *testing.T) {
	task := NewTaskDefinition()

	if task == nil {
		t.Fatal("NewTaskDefinition: Expected task to be created, got nil")
	}

	if task.ID() == "" {
		t.Error("NewTaskDefinition: Expected ID to be set")
	}

	if task.Status() != TaskDefinitionStatusActive {
		t.Errorf("NewTaskDefinition: Expected status to be %s, got %s", TaskDefinitionStatusActive, task.Status())
	}

	if task.Memo() != "" {
		t.Errorf("NewTaskDefinition: Expected memo to be empty, got %s", task.Memo())
	}

	if task.CreatedAt() == "" {
		t.Error("NewTaskDefinition: Expected CreatedAt to be set")
	}

	if task.UpdatedAt() == "" {
		t.Error("NewTaskDefinition: Expected UpdatedAt to be set")
	}

	if task.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Errorf("NewTaskDefinition: Expected SoftDeletedAt to be %s, got %s", sb.MAX_DATETIME, task.SoftDeletedAt())
	}
}

func TestNewTaskDefinitionFromExistingData(t *testing.T) {
	data := map[string]string{
		COLUMN_ID:              "test-id",
		COLUMN_ALIAS:           "test-alias",
		COLUMN_TITLE:           "Test Title",
		COLUMN_DESCRIPTION:     "Test Description",
		COLUMN_STATUS:          TaskDefinitionStatusCanceled,
		COLUMN_MEMO:            "Test Memo",
		COLUMN_CREATED_AT:      "2023-01-01 12:00:00",
		COLUMN_UPDATED_AT:      "2023-01-02 12:00:00",
		COLUMN_SOFT_DELETED_AT: "2023-01-03 12:00:00",
	}

	task := NewTaskDefinitionFromExistingData(data)

	if task.ID() != "test-id" {
		t.Errorf("NewTaskDefinitionFromExistingData: Expected ID to be 'test-id', got %s", task.ID())
	}

	if task.Alias() != "test-alias" {
		t.Errorf("NewTaskDefinitionFromExistingData: Expected Alias to be 'test-alias', got %s", task.Alias())
	}

	if task.Title() != "Test Title" {
		t.Errorf("NewTaskDefinitionFromExistingData: Expected Title to be 'Test Title', got %s", task.Title())
	}

	if task.Description() != "Test Description" {
		t.Errorf("NewTaskDefinitionFromExistingData: Expected Description to be 'Test Description', got %s", task.Description())
	}

	if task.Status() != TaskDefinitionStatusCanceled {
		t.Errorf("NewTaskDefinitionFromExistingData: Expected Status to be %s, got %s", TaskDefinitionStatusCanceled, task.Status())
	}

	if task.Memo() != "Test Memo" {
		t.Errorf("NewTaskDefinitionFromExistingData: Expected Memo to be 'Test Memo', got %s", task.Memo())
	}
}

func TestTaskDefinition_IsActive(t *testing.T) {
	task := NewTaskDefinition()

	// Test active status
	task.SetStatus(TaskDefinitionStatusActive)
	if !task.IsActive() {
		t.Error("IsActive: Expected task to be active when status is TaskDefinitionStatusActive")
	}

	// Test non-active status
	task.SetStatus(TaskDefinitionStatusCanceled)
	if task.IsActive() {
		t.Error("IsActive: Expected task to not be active when status is TaskDefinitionStatusCanceled")
	}
}

func TestTaskDefinition_IsCanceled(t *testing.T) {
	task := NewTaskDefinition()

	// Test canceled status
	task.SetStatus(TaskDefinitionStatusCanceled)
	if !task.IsCanceled() {
		t.Error("IsCanceled: Expected task to be canceled when status is TaskDefinitionStatusCanceled")
	}

	// Test non-canceled status
	task.SetStatus(TaskDefinitionStatusActive)
	if task.IsCanceled() {
		t.Error("IsCanceled: Expected task to not be canceled when status is TaskDefinitionStatusActive")
	}
}

func TestTaskDefinition_IsSoftDeleted(t *testing.T) {
	task := NewTaskDefinition()

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

func TestTaskDefinition_CreatedAtCarbon(t *testing.T) {
	task := NewTaskDefinition()
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

func TestTaskDefinition_UpdatedAtCarbon(t *testing.T) {
	task := NewTaskDefinition()
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

func TestTaskDefinition_SoftDeletedAtCarbon(t *testing.T) {
	task := NewTaskDefinition()
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

func TestTaskDefinition_SettersAndGetters(t *testing.T) {
	task := NewTaskDefinition()

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
	task.SetStatus(TaskDefinitionStatusCanceled)
	if task.Status() != TaskDefinitionStatusCanceled {
		t.Errorf("Status: Expected %s, got %s", TaskDefinitionStatusCanceled, task.Status())
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

func TestTaskDefinition_ChainedSetters(t *testing.T) {
	task := NewTaskDefinition()

	// Test that setters return the task instance for chaining
	result := task.SetID("test-id").
		SetAlias("test-alias").
		SetTitle("Test Title").
		SetDescription("Test Description").
		SetMemo("Test Memo").
		SetStatus(TaskDefinitionStatusCanceled).
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
	if task.Status() != TaskDefinitionStatusCanceled {
		t.Error("ChainedSetters: Status not set correctly")
	}
}
