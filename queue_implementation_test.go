package taskstore

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
)

func TestNewQueue(t *testing.T) {
	queue := NewQueue()

	if queue == nil {
		t.Fatal("NewQueue: Expected queue to be created, got nil")
	}

	if queue.ID() == "" {
		t.Error("NewQueue: Expected ID to be set")
	}

	if queue.Status() != QueueStatusQueued {
		t.Errorf("NewQueue: Expected status to be %s, got %s", QueueStatusQueued, queue.Status())
	}

	if queue.CreatedAt() == "" {
		t.Error("NewQueue: Expected CreatedAt to be set")
	}

	if queue.UpdatedAt() == "" {
		t.Error("NewQueue: Expected UpdatedAt to be set")
	}

	if queue.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Errorf("NewQueue: Expected SoftDeletedAt to be %s, got %s", sb.MAX_DATETIME, queue.SoftDeletedAt())
	}
}

func TestNewQueueFromExistingData(t *testing.T) {
	data := map[string]string{
		COLUMN_ID:           "test-queue-id",
		COLUMN_TASK_ID:      "test-task-id",
		COLUMN_STATUS:       QueueStatusRunning,
		COLUMN_ATTEMPTS:     "3",
		COLUMN_PARAMETERS:   `{"key":"value"}`,
		COLUMN_OUTPUT:       "test output",
		COLUMN_DETAILS:      "test details",
		COLUMN_CREATED_AT:   "2023-01-01 12:00:00",
		COLUMN_UPDATED_AT:   "2023-01-02 12:00:00",
		COLUMN_STARTED_AT:   "2023-01-01 12:30:00",
		COLUMN_COMPLETED_AT: "2023-01-01 13:00:00",
		COLUMN_DELETED_AT:   "2023-01-03 12:00:00",
	}

	queue := NewQueueFromExistingData(data)

	if queue.ID() != "test-queue-id" {
		t.Errorf("NewQueueFromExistingData: Expected ID to be 'test-queue-id', got %s", queue.ID())
	}

	if queue.TaskID() != "test-task-id" {
		t.Errorf("NewQueueFromExistingData: Expected TaskID to be 'test-task-id', got %s", queue.TaskID())
	}

	if queue.Status() != QueueStatusRunning {
		t.Errorf("NewQueueFromExistingData: Expected Status to be %s, got %s", QueueStatusRunning, queue.Status())
	}

	if queue.Attempts() != 3 {
		t.Errorf("NewQueueFromExistingData: Expected Attempts to be 3, got %d", queue.Attempts())
	}

	if queue.Parameters() != `{"key":"value"}` {
		t.Errorf("NewQueueFromExistingData: Expected Parameters to be '{\"key\":\"value\"}', got %s", queue.Parameters())
	}

	if queue.Output() != "test output" {
		t.Errorf("NewQueueFromExistingData: Expected Output to be 'test output', got %s", queue.Output())
	}

	if queue.Details() != "test details" {
		t.Errorf("NewQueueFromExistingData: Expected Details to be 'test details', got %s", queue.Details())
	}
}

func TestQueue_StatusCheckers(t *testing.T) {
	queue := NewQueue()

	// Test IsCanceled
	queue.SetStatus(QueueStatusCanceled)
	if !queue.IsCanceled() {
		t.Error("IsCanceled: Expected queue to be canceled when status is QueueStatusCanceled")
	}
	if queue.IsDeleted() || queue.IsFailed() || queue.IsQueued() || queue.IsPaused() || queue.IsRunning() || queue.IsSuccess() {
		t.Error("IsCanceled: Expected other status checkers to return false")
	}

	// Test IsDeleted
	queue.SetStatus(QueueStatusDeleted)
	if !queue.IsDeleted() {
		t.Error("IsDeleted: Expected queue to be deleted when status is QueueStatusDeleted")
	}
	if queue.IsCanceled() || queue.IsFailed() || queue.IsQueued() || queue.IsPaused() || queue.IsRunning() || queue.IsSuccess() {
		t.Error("IsDeleted: Expected other status checkers to return false")
	}

	// Test IsFailed
	queue.SetStatus(QueueStatusFailed)
	if !queue.IsFailed() {
		t.Error("IsFailed: Expected queue to be failed when status is QueueStatusFailed")
	}
	if queue.IsCanceled() || queue.IsDeleted() || queue.IsQueued() || queue.IsPaused() || queue.IsRunning() || queue.IsSuccess() {
		t.Error("IsFailed: Expected other status checkers to return false")
	}

	// Test IsQueued
	queue.SetStatus(QueueStatusQueued)
	if !queue.IsQueued() {
		t.Error("IsQueued: Expected queue to be queued when status is QueueStatusQueued")
	}
	if queue.IsCanceled() || queue.IsDeleted() || queue.IsFailed() || queue.IsPaused() || queue.IsRunning() || queue.IsSuccess() {
		t.Error("IsQueued: Expected other status checkers to return false")
	}

	// Test IsPaused
	queue.SetStatus(QueueStatusPaused)
	if !queue.IsPaused() {
		t.Error("IsPaused: Expected queue to be paused when status is QueueStatusPaused")
	}
	if queue.IsCanceled() || queue.IsDeleted() || queue.IsFailed() || queue.IsQueued() || queue.IsRunning() || queue.IsSuccess() {
		t.Error("IsPaused: Expected other status checkers to return false")
	}

	// Test IsRunning
	queue.SetStatus(QueueStatusRunning)
	if !queue.IsRunning() {
		t.Error("IsRunning: Expected queue to be running when status is QueueStatusRunning")
	}
	if queue.IsCanceled() || queue.IsDeleted() || queue.IsFailed() || queue.IsQueued() || queue.IsPaused() || queue.IsSuccess() {
		t.Error("IsRunning: Expected other status checkers to return false")
	}

	// Test IsSuccess
	queue.SetStatus(QueueStatusSuccess)
	if !queue.IsSuccess() {
		t.Error("IsSuccess: Expected queue to be success when status is QueueStatusSuccess")
	}
	if queue.IsCanceled() || queue.IsDeleted() || queue.IsFailed() || queue.IsQueued() || queue.IsPaused() || queue.IsRunning() {
		t.Error("IsSuccess: Expected other status checkers to return false")
	}
}

func TestQueue_IsSoftDeleted(t *testing.T) {
	queue := NewQueue()

	// Test not soft deleted (default state)
	if queue.IsSoftDeleted() {
		t.Error("IsSoftDeleted: Expected new queue to not be soft deleted")
	}

	// Test soft deleted
	pastTime := carbon.Now(carbon.UTC).SubHours(1).ToDateTimeString(carbon.UTC)
	queue.SetSoftDeletedAt(pastTime)
	if !queue.IsSoftDeleted() {
		t.Error("IsSoftDeleted: Expected queue to be soft deleted when deleted_at is in the past")
	}

	// Test future deletion time (not yet deleted)
	futureTime := carbon.Now(carbon.UTC).AddHours(1).ToDateTimeString(carbon.UTC)
	queue.SetSoftDeletedAt(futureTime)
	if queue.IsSoftDeleted() {
		t.Error("IsSoftDeleted: Expected queue to not be soft deleted when deleted_at is in the future")
	}
}

func TestQueue_AppendDetails(t *testing.T) {
	queue := NewQueue()

	// Test appending to empty details
	queue.AppendDetails("First detail")
	details := queue.Details()
	if !strings.Contains(details, "First detail") {
		t.Error("AppendDetails: Expected details to contain 'First detail'")
	}
	if !strings.Contains(details, ":") {
		t.Error("AppendDetails: Expected details to contain timestamp separator ':'")
	}

	// Test appending to existing details
	queue.AppendDetails("Second detail")
	details = queue.Details()
	if !strings.Contains(details, "First detail") {
		t.Error("AppendDetails: Expected details to still contain 'First detail'")
	}
	if !strings.Contains(details, "Second detail") {
		t.Error("AppendDetails: Expected details to contain 'Second detail'")
	}
	if !strings.Contains(details, "\n") {
		t.Error("AppendDetails: Expected details to contain newline separator")
	}

	// Verify the format includes timestamps
	lines := strings.Split(details, "\n")
	if len(lines) < 2 {
		t.Error("AppendDetails: Expected at least 2 lines in details")
	}
	for _, line := range lines {
		if line != "" && !strings.Contains(line, ":") {
			t.Errorf("AppendDetails: Expected each line to contain timestamp, got: %s", line)
		}
	}
}

func TestQueue_ParametersMap(t *testing.T) {
	queue := NewQueue()

	// Test with valid JSON parameters
	testParams := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	jsonBytes, _ := json.Marshal(testParams)
	queue.SetParameters(string(jsonBytes))

	retrievedParams, err := queue.ParametersMap()
	if err != nil {
		t.Fatalf("ParametersMap: Unexpected error: %v", err)
	}

	if len(retrievedParams) != len(testParams) {
		t.Errorf("ParametersMap: Expected %d parameters, got %d", len(testParams), len(retrievedParams))
	}

	for key, expectedValue := range testParams {
		if actualValue, exists := retrievedParams[key]; !exists {
			t.Errorf("ParametersMap: Expected key '%s' to exist", key)
		} else if actualValue != expectedValue {
			t.Errorf("ParametersMap: Expected value '%s' for key '%s', got '%s'", expectedValue, key, actualValue)
		}
	}

	// Test with invalid JSON
	queue.SetParameters("invalid json")
	_, err = queue.ParametersMap()
	if err == nil {
		t.Error("ParametersMap: Expected error for invalid JSON, got nil")
	}

	// Test with empty parameters
	queue.SetParameters("")
	_, err = queue.ParametersMap()
	if err == nil {
		t.Error("ParametersMap: Expected error for empty parameters, got nil")
	}
}

func TestQueue_SetParametersMap(t *testing.T) {
	queue := NewQueue()

	testParams := map[string]string{
		"param1": "value1",
		"param2": "value2",
		"param3": "value3",
	}

	result, err := queue.SetParametersMap(testParams)
	if err != nil {
		t.Fatalf("SetParametersMap: Unexpected error: %v", err)
	}

	if result != queue {
		t.Error("SetParametersMap: Expected method to return the same queue instance")
	}

	// Verify the parameters were set correctly
	retrievedParams, err := queue.ParametersMap()
	if err != nil {
		t.Fatalf("SetParametersMap: Error retrieving parameters: %v", err)
	}

	if len(retrievedParams) != len(testParams) {
		t.Errorf("SetParametersMap: Expected %d parameters, got %d", len(testParams), len(retrievedParams))
	}

	for key, expectedValue := range testParams {
		if actualValue, exists := retrievedParams[key]; !exists {
			t.Errorf("SetParametersMap: Expected key '%s' to exist", key)
		} else if actualValue != expectedValue {
			t.Errorf("SetParametersMap: Expected value '%s' for key '%s', got '%s'", expectedValue, key, actualValue)
		}
	}

	// Verify the JSON string is valid
	parametersJSON := queue.Parameters()
	var jsonTest map[string]string
	err = json.Unmarshal([]byte(parametersJSON), &jsonTest)
	if err != nil {
		t.Errorf("SetParametersMap: Generated invalid JSON: %v", err)
	}
}

func TestQueue_CarbonMethods(t *testing.T) {
	queue := NewQueue()

	// Test CreatedAtCarbon
	createdAtStr := "2023-01-01 12:00:00"
	queue.SetCreatedAt(createdAtStr)
	createdAtCarbon := queue.CreatedAtCarbon()
	if createdAtCarbon == nil {
		t.Fatal("CreatedAtCarbon: Expected carbon instance, got nil")
	}
	if createdAtCarbon.ToDateTimeString(carbon.UTC) != createdAtStr {
		t.Errorf("CreatedAtCarbon: Expected %s, got %s", createdAtStr, createdAtCarbon.ToDateTimeString(carbon.UTC))
	}

	// Test UpdatedAtCarbon
	updatedAtStr := "2023-01-02 15:30:45"
	queue.SetUpdatedAt(updatedAtStr)
	updatedAtCarbon := queue.UpdatedAtCarbon()
	if updatedAtCarbon == nil {
		t.Fatal("UpdatedAtCarbon: Expected carbon instance, got nil")
	}
	if updatedAtCarbon.ToDateTimeString(carbon.UTC) != updatedAtStr {
		t.Errorf("UpdatedAtCarbon: Expected %s, got %s", updatedAtStr, updatedAtCarbon.ToDateTimeString(carbon.UTC))
	}

	// Test StartedAtCarbon
	startedAtStr := "2023-01-01 12:30:00"
	queue.SetStartedAt(startedAtStr)
	startedAtCarbon := queue.StartedAtCarbon()
	if startedAtCarbon == nil {
		t.Fatal("StartedAtCarbon: Expected carbon instance, got nil")
	}
	if startedAtCarbon.ToDateTimeString(carbon.UTC) != startedAtStr {
		t.Errorf("StartedAtCarbon: Expected %s, got %s", startedAtStr, startedAtCarbon.ToDateTimeString(carbon.UTC))
	}

	// Test CompletedAtCarbon
	completedAtStr := "2023-01-01 13:00:00"
	queue.SetCompletedAt(completedAtStr)
	completedAtCarbon := queue.CompletedAtCarbon()
	if completedAtCarbon == nil {
		t.Fatal("CompletedAtCarbon: Expected carbon instance, got nil")
	}
	if completedAtCarbon.ToDateTimeString(carbon.UTC) != completedAtStr {
		t.Errorf("CompletedAtCarbon: Expected %s, got %s", completedAtStr, completedAtCarbon.ToDateTimeString(carbon.UTC))
	}

	// Test SoftDeletedAtCarbon
	deletedAtStr := "2023-01-03 09:15:30"
	queue.SetSoftDeletedAt(deletedAtStr)
	deletedAtCarbon := queue.SoftDeletedAtCarbon()
	if deletedAtCarbon == nil {
		t.Fatal("SoftDeletedAtCarbon: Expected carbon instance, got nil")
	}
	if deletedAtCarbon.ToDateTimeString(carbon.UTC) != deletedAtStr {
		t.Errorf("SoftDeletedAtCarbon: Expected %s, got %s", deletedAtStr, deletedAtCarbon.ToDateTimeString(carbon.UTC))
	}
}

func TestQueue_AttemptsHandling(t *testing.T) {
	queue := NewQueue()

	// Test setting and getting attempts
	queue.SetAttempts(5)
	if queue.Attempts() != 5 {
		t.Errorf("Attempts: Expected 5, got %d", queue.Attempts())
	}

	// Test with zero attempts
	queue.SetAttempts(0)
	if queue.Attempts() != 0 {
		t.Errorf("Attempts: Expected 0, got %d", queue.Attempts())
	}

	// Test with negative attempts (edge case)
	queue.SetAttempts(-1)
	if queue.Attempts() != -1 {
		t.Errorf("Attempts: Expected -1, got %d", queue.Attempts())
	}
}

func TestQueue_SettersAndGetters(t *testing.T) {
	queue := NewQueue()

	// Test ID
	testID := "test-queue-id"
	queue.SetID(testID)
	if queue.ID() != testID {
		t.Errorf("ID: Expected %s, got %s", testID, queue.ID())
	}

	// Test TaskID
	testTaskID := "test-task-id"
	queue.SetTaskID(testTaskID)
	if queue.TaskID() != testTaskID {
		t.Errorf("TaskID: Expected %s, got %s", testTaskID, queue.TaskID())
	}

	// Test Status
	queue.SetStatus(QueueStatusRunning)
	if queue.Status() != QueueStatusRunning {
		t.Errorf("Status: Expected %s, got %s", QueueStatusRunning, queue.Status())
	}

	// Test Output
	testOutput := "Test output message"
	queue.SetOutput(testOutput)
	if queue.Output() != testOutput {
		t.Errorf("Output: Expected %s, got %s", testOutput, queue.Output())
	}

	// Test Details
	testDetails := "Test details message"
	queue.SetDetails(testDetails)
	if queue.Details() != testDetails {
		t.Errorf("Details: Expected %s, got %s", testDetails, queue.Details())
	}

	// Test Parameters
	testParameters := `{"key":"value"}`
	queue.SetParameters(testParameters)
	if queue.Parameters() != testParameters {
		t.Errorf("Parameters: Expected %s, got %s", testParameters, queue.Parameters())
	}

	// Test timestamps
	testCreatedAt := "2023-01-01 10:00:00"
	queue.SetCreatedAt(testCreatedAt)
	if queue.CreatedAt() != testCreatedAt {
		t.Errorf("CreatedAt: Expected %s, got %s", testCreatedAt, queue.CreatedAt())
	}

	testUpdatedAt := "2023-01-02 11:00:00"
	queue.SetUpdatedAt(testUpdatedAt)
	if queue.UpdatedAt() != testUpdatedAt {
		t.Errorf("UpdatedAt: Expected %s, got %s", testUpdatedAt, queue.UpdatedAt())
	}

	testStartedAt := "2023-01-01 10:30:00"
	queue.SetStartedAt(testStartedAt)
	if queue.StartedAt() != testStartedAt {
		t.Errorf("StartedAt: Expected %s, got %s", testStartedAt, queue.StartedAt())
	}

	testCompletedAt := "2023-01-01 11:30:00"
	queue.SetCompletedAt(testCompletedAt)
	if queue.CompletedAt() != testCompletedAt {
		t.Errorf("CompletedAt: Expected %s, got %s", testCompletedAt, queue.CompletedAt())
	}

	testDeletedAt := "2023-01-03 12:00:00"
	queue.SetSoftDeletedAt(testDeletedAt)
	if queue.SoftDeletedAt() != testDeletedAt {
		t.Errorf("SoftDeletedAt: Expected %s, got %s", testDeletedAt, queue.SoftDeletedAt())
	}
}

func TestQueue_ChainedSetters(t *testing.T) {
	queue := NewQueue()

	// Test that setters return the queue instance for chaining
	result := queue.SetID("test-id").
		SetTaskID("test-task-id").
		SetStatus(QueueStatusRunning).
		SetAttempts(3).
		SetOutput("test output").
		SetDetails("test details").
		SetParameters(`{"key":"value"}`).
		SetCreatedAt("2023-01-01 10:00:00").
		SetUpdatedAt("2023-01-02 11:00:00").
		SetStartedAt("2023-01-01 10:30:00").
		SetCompletedAt("2023-01-01 11:30:00").
		SetSoftDeletedAt("2023-01-03 12:00:00")

	if result != queue {
		t.Error("ChainedSetters: Expected setters to return the same queue instance for chaining")
	}

	// Verify all values were set correctly
	if queue.ID() != "test-id" {
		t.Error("ChainedSetters: ID not set correctly")
	}
	if queue.TaskID() != "test-task-id" {
		t.Error("ChainedSetters: TaskID not set correctly")
	}
	if queue.Status() != QueueStatusRunning {
		t.Error("ChainedSetters: Status not set correctly")
	}
	if queue.Attempts() != 3 {
		t.Error("ChainedSetters: Attempts not set correctly")
	}
	if queue.Output() != "test output" {
		t.Error("ChainedSetters: Output not set correctly")
	}
	if queue.Details() != "test details" {
		t.Error("ChainedSetters: Details not set correctly")
	}
	if queue.Parameters() != `{"key":"value"}` {
		t.Error("ChainedSetters: Parameters not set correctly")
	}
}
