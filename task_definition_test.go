package taskstore

import (
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestNewTaskDefinition(t *testing.T) {
	task := NewTaskDefinition()

	if task == nil {
		t.Fatal("NewTaskDefinition: Expected task to be created, got nil")
	}

	if task.GetID() == "" {
		t.Error("NewTaskDefinition: Expected ID to be set")
	}

	if task.GetStatus() != TaskDefinitionStatusActive {
		t.Errorf("NewTaskDefinition: Expected status to be %s, got %s", TaskDefinitionStatusActive, task.GetStatus())
	}

	if task.GetMemo() != "" {
		t.Errorf("NewTaskDefinition: Expected memo to be empty, got %s", task.GetMemo())
	}

	if task.GetCreatedAt().IsZero() {
		t.Error("NewTaskDefinition: Expected CreatedAt to be set")
	}

	if task.GetUpdatedAt().IsZero() {
		t.Error("NewTaskDefinition: Expected UpdatedAt to be set")
	}

	maxDateTime := carbon.Parse(MAX_DATETIME, carbon.UTC).StdTime()
	if !task.GetSoftDeletedAt().Equal(maxDateTime) {
		t.Errorf("NewTaskDefinition: Expected SoftDeletedAt to be %s, got %s", MAX_DATETIME, task.GetSoftDeletedAt().Format("2006-01-02 15:04:05"))
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

	if task.GetID() != "test-id" {
		t.Errorf("NewTaskDefinitionFromExistingData: Expected ID to be 'test-id', got %s", task.GetID())
	}

	if task.GetAlias() != "test-alias" {
		t.Errorf("NewTaskDefinitionFromExistingData: Expected Alias to be 'test-alias', got %s", task.GetAlias())
	}

	if task.GetTitle() != "Test Title" {
		t.Errorf("NewTaskDefinitionFromExistingData: Expected Title to be 'Test Title', got %s", task.GetTitle())
	}

	if task.GetDescription() != "Test Description" {
		t.Errorf("NewTaskDefinitionFromExistingData: Expected Description to be 'Test Description', got %s", task.GetDescription())
	}

	if task.GetStatus() != TaskDefinitionStatusCanceled {
		t.Errorf("NewTaskDefinitionFromExistingData: Expected Status to be %s, got %s", TaskDefinitionStatusCanceled, task.GetStatus())
	}

	if task.GetMemo() != "Test Memo" {
		t.Errorf("NewTaskDefinitionFromExistingData: Expected Memo to be 'Test Memo', got %s", task.GetMemo())
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
	pastTime := carbon.Now(carbon.UTC).SubHours(1).StdTime()
	task.SetSoftDeletedAt(pastTime)
	if !task.IsSoftDeleted() {
		t.Error("IsSoftDeleted: Expected task to be soft deleted when deleted_at is in the past")
	}

	// Test future deletion time (not yet deleted)
	futureTime := carbon.Now(carbon.UTC).AddHours(1).StdTime()
	task.SetSoftDeletedAt(futureTime)
	if task.IsSoftDeleted() {
		t.Error("IsSoftDeleted: Expected task to not be soft deleted when deleted_at is in the future")
	}
}

func TestTaskDefinition_CreatedAtCarbon(t *testing.T) {
	task := NewTaskDefinition()
	createdAtStr := "2023-01-01 12:00:00"
	task.SetCreatedAt(carbon.Parse(createdAtStr, carbon.UTC).StdTime())

	createdAtCarbon := task.GetCreatedAtCarbon()
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
	task.SetUpdatedAt(carbon.Parse(updatedAtStr, carbon.UTC).StdTime())

	updatedAtCarbon := task.GetUpdatedAtCarbon()
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
	task.SetSoftDeletedAt(carbon.Parse(deletedAtStr, carbon.UTC).StdTime())

	deletedAtCarbon := task.GetSoftDeletedAtCarbon()
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
	if task.GetID() != testID {
		t.Errorf("ID: Expected %s, got %s", testID, task.GetID())
	}

	// Test Alias
	testAlias := "test-alias"
	task.SetAlias(testAlias)
	if task.GetAlias() != testAlias {
		t.Errorf("Alias: Expected %s, got %s", testAlias, task.GetAlias())
	}

	// Test Title
	testTitle := "Test Task Title"
	task.SetTitle(testTitle)
	if task.GetTitle() != testTitle {
		t.Errorf("Title: Expected %s, got %s", testTitle, task.GetTitle())
	}

	// Test Description
	testDescription := "Test task description"
	task.SetDescription(testDescription)
	if task.GetDescription() != testDescription {
		t.Errorf("Description: Expected %s, got %s", testDescription, task.GetDescription())
	}

	// Test Memo
	testMemo := "Test memo"
	task.SetMemo(testMemo)
	if task.GetMemo() != testMemo {
		t.Errorf("Memo: Expected %s, got %s", testMemo, task.GetMemo())
	}

	// Test Status
	task.SetStatus(TaskDefinitionStatusCanceled)
	if task.GetStatus() != TaskDefinitionStatusCanceled {
		t.Errorf("Status: Expected %s, got %s", TaskDefinitionStatusCanceled, task.GetStatus())
	}

	// Test CreatedAt
	testCreatedAt := "2023-01-01 10:00:00"
	task.SetCreatedAt(carbon.Parse(testCreatedAt, carbon.UTC).StdTime())
	if task.GetCreatedAt().Format("2006-01-02 15:04:05") != testCreatedAt {
		t.Errorf("CreatedAt: Expected %s, got %s", testCreatedAt, task.GetCreatedAt().Format("2006-01-02 15:04:05"))
	}

	// Test UpdatedAt
	testUpdatedAt := "2023-01-02 11:00:00"
	task.SetUpdatedAt(carbon.Parse(testUpdatedAt, carbon.UTC).StdTime())
	if task.GetUpdatedAt().Format("2006-01-02 15:04:05") != testUpdatedAt {
		t.Errorf("UpdatedAt: Expected %s, got %s", testUpdatedAt, task.GetUpdatedAt().Format("2006-01-02 15:04:05"))
	}

	// Test SoftDeletedAt
	testDeletedAt := "2023-01-03 12:00:00"
	task.SetSoftDeletedAt(carbon.Parse(testDeletedAt, carbon.UTC).StdTime())
	if task.GetSoftDeletedAt().Format("2006-01-02 15:04:05") != testDeletedAt {
		t.Errorf("SoftDeletedAt: Expected %s, got %s", testDeletedAt, task.GetSoftDeletedAt().Format("2006-01-02 15:04:05"))
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
		SetCreatedAt(carbon.Parse("2023-01-01 10:00:00", carbon.UTC).StdTime()).
		SetUpdatedAt(carbon.Parse("2023-01-02 11:00:00", carbon.UTC).StdTime()).
		SetSoftDeletedAt(carbon.Parse("2023-01-03 12:00:00", carbon.UTC).StdTime())

	if result != task {
		t.Error("ChainedSetters: Expected setters to return the same task instance for chaining")
	}

	// Verify all values were set correctly
	if task.GetID() != "test-id" {
		t.Error("ChainedSetters: ID not set correctly")
	}
	if task.GetAlias() != "test-alias" {
		t.Error("ChainedSetters: Alias not set correctly")
	}
	if task.GetTitle() != "Test Title" {
		t.Error("ChainedSetters: Title not set correctly")
	}
	if task.GetDescription() != "Test Description" {
		t.Error("ChainedSetters: Description not set correctly")
	}
	if task.GetMemo() != "Test Memo" {
		t.Error("ChainedSetters: Memo not set correctly")
	}
	if task.GetStatus() != TaskDefinitionStatusCanceled {
		t.Error("ChainedSetters: Status not set correctly")
	}
}
