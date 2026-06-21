package taskstore

import (
	"encoding/json"
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestNewSchedule(t *testing.T) {
	schedule := NewSchedule()

	if schedule == nil {
		t.Fatal("expected schedule to not be nil")
	}
	if schedule.GetID() == "" {
		t.Error("expected ID to not be empty")
	}
	if schedule.GetStatus() != "draft" {
		t.Errorf("expected status 'draft', got %s", schedule.GetStatus())
	}
	if schedule.GetStartAt() != NULL_DATETIME {
		t.Error("expected StartAt to be NULL_DATETIME")
	}
	if schedule.GetEndAt() != MAX_DATETIME {
		t.Error("expected EndAt to be MAX_DATETIME")
	}
	if schedule.GetLastRunAt() != NULL_DATETIME {
		t.Error("expected LastRunAt to be NULL_DATETIME")
	}
	if schedule.GetNextRunAt() != NULL_DATETIME {
		t.Error("expected NextRunAt to be NULL_DATETIME")
	}
	maxDateTime := carbon.Parse(MAX_DATETIME, carbon.UTC).StdTime()
	if !schedule.GetSoftDeletedAt().Equal(maxDateTime) {
		t.Error("expected SoftDeletedAt to be MAX_DATETIME")
	}
	if schedule.GetRecurrenceRule() == nil {
		t.Error("expected RecurrenceRule to not be nil")
	}
	if schedule.GetCreatedAt().IsZero() {
		t.Error("expected CreatedAt to not be empty")
	}
	if schedule.GetUpdatedAt().IsZero() {
		t.Error("expected UpdatedAt to not be empty")
	}
}

func TestScheduleGettersAndSetters(t *testing.T) {
	schedule := NewSchedule()

	// Test ID
	schedule.SetID("test-id-123")
	if schedule.GetID() != "test-id-123" {
		t.Errorf("expected ID 'test-id-123', got %s", schedule.GetID())
	}

	// Test Name
	schedule.SetName("Test Schedule")
	if schedule.GetName() != "Test Schedule" {
		t.Errorf("expected name 'Test Schedule', got %s", schedule.GetName())
	}

	// Test Description
	schedule.SetDescription("Test Description")
	if schedule.GetDescription() != "Test Description" {
		t.Errorf("expected description 'Test Description', got %s", schedule.GetDescription())
	}

	// Test Status
	schedule.SetStatus("active")
	if schedule.GetStatus() != "active" {
		t.Errorf("expected status 'active', got %s", schedule.GetStatus())
	}

	// Test QueueName
	schedule.SetQueueName("test-queue")
	if schedule.GetQueueName() != "test-queue" {
		t.Errorf("expected queue name 'test-queue', got %s", schedule.GetQueueName())
	}

	// Test TaskDefinitionID
	schedule.SetTaskDefinitionID("task-def-123")
	if schedule.GetTaskDefinitionID() != "task-def-123" {
		t.Errorf("expected task definition ID 'task-def-123', got %s", schedule.GetTaskDefinitionID())
	}

	// Test TaskParameters
	params := map[string]any{"key": "value", "count": float64(42)}
	schedule.SetTaskParameters(params)
	gotParams := schedule.GetTaskParameters()
	if gotParams["key"] != params["key"] || gotParams["count"] != params["count"] {
		t.Error("expected parameters to match")
	}

	// Test StartAt
	startAt := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	schedule.SetStartAt(startAt)
	if schedule.GetStartAt() != startAt {
		t.Errorf("expected start at %s, got %s", startAt, schedule.GetStartAt())
	}

	// Test EndAt
	endAt := carbon.Now(carbon.UTC).AddDays(30).ToDateTimeString(carbon.UTC)
	schedule.SetEndAt(endAt)
	if schedule.GetEndAt() != endAt {
		t.Errorf("expected end at %s, got %s", endAt, schedule.GetEndAt())
	}

	// Test ExecutionCount
	schedule.SetExecutionCount(5)
	if schedule.GetExecutionCount() != 5 {
		t.Errorf("expected execution count 5, got %d", schedule.GetExecutionCount())
	}

	// Test MaxExecutionCount
	schedule.SetMaxExecutionCount(10)
	if schedule.GetMaxExecutionCount() != 10 {
		t.Errorf("expected max execution count 10, got %d", schedule.GetMaxExecutionCount())
	}

	// Test LastRunAt
	lastRunAt := carbon.Now(carbon.UTC).AddMinutes(-5).ToDateTimeString(carbon.UTC)
	schedule.SetLastRunAt(lastRunAt)
	if schedule.GetLastRunAt() != lastRunAt {
		t.Errorf("expected last run at %s, got %s", lastRunAt, schedule.GetLastRunAt())
	}

	// Test NextRunAt
	nextRunAt := carbon.Now(carbon.UTC).AddMinutes(5).ToDateTimeString(carbon.UTC)
	schedule.SetNextRunAt(nextRunAt)
	if schedule.GetNextRunAt() != nextRunAt {
		t.Errorf("expected next run at %s, got %s", nextRunAt, schedule.GetNextRunAt())
	}

	// Test CreatedAt
	createdAt := carbon.Now(carbon.UTC).AddDays(-1).StdTime()
	schedule.SetCreatedAt(createdAt)
	if !schedule.GetCreatedAt().Equal(createdAt) {
		t.Errorf("expected created at %s, got %s", createdAt.Format("2006-01-02 15:04:05"), schedule.GetCreatedAt().Format("2006-01-02 15:04:05"))
	}

	// Test UpdatedAt
	updatedAt := carbon.Now(carbon.UTC).StdTime()
	schedule.SetUpdatedAt(updatedAt)
	if !schedule.GetUpdatedAt().Equal(updatedAt) {
		t.Errorf("expected updated at %s, got %s", updatedAt.Format("2006-01-02 15:04:05"), schedule.GetUpdatedAt().Format("2006-01-02 15:04:05"))
	}

	// Test SoftDeletedAt
	softDeletedAt := carbon.Now(carbon.UTC).StdTime()
	schedule.SetSoftDeletedAt(softDeletedAt)
	if !schedule.GetSoftDeletedAt().Equal(softDeletedAt) {
		t.Errorf("expected soft deleted at %s, got %s", softDeletedAt.Format("2006-01-02 15:04:05"), schedule.GetSoftDeletedAt().Format("2006-01-02 15:04:05"))
	}

	// Test RecurrenceRule
	rr := NewRecurrenceRule()
	rr.SetFrequency(FrequencyDaily)
	schedule.SetRecurrenceRule(rr)
	if schedule.GetRecurrenceRule().GetFrequency() != FrequencyDaily {
		t.Error("expected frequency to be Daily")
	}
}

func TestScheduleJSONMarshaling(t *testing.T) {
	schedule := NewSchedule()
	schedule.SetName("Test Schedule")
	schedule.SetDescription("Test Description")
	schedule.SetStatus("active")
	schedule.SetQueueName("default")
	schedule.SetTaskDefinitionID("task-123")
	schedule.SetTaskParameters(map[string]any{"key": "value"})
	schedule.SetExecutionCount(3)
	schedule.SetMaxExecutionCount(10)

	rr := NewRecurrenceRule()
	rr.SetFrequency(FrequencyDaily)
	rr.SetInterval(1)
	schedule.SetRecurrenceRule(rr)

	// Marshal to JSON
	jsonBytes, err := json.Marshal(schedule)
	if err != nil {
		t.Fatal(err)
	}
	if len(jsonBytes) == 0 {
		t.Error("expected JSON bytes to not be empty")
	}

	// Unmarshal from JSON - create a new schedule and unmarshal into it
	unmarshaled := NewSchedule()
	err = json.Unmarshal(jsonBytes, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	// Verify fields
	if schedule.GetID() != unmarshaled.GetID() {
		t.Errorf("expected ID %s, got %s", schedule.GetID(), unmarshaled.GetID())
	}
	if schedule.GetName() != unmarshaled.GetName() {
		t.Errorf("expected name %s, got %s", schedule.GetName(), unmarshaled.GetName())
	}
	if schedule.GetDescription() != unmarshaled.GetDescription() {
		t.Errorf("expected description %s, got %s", schedule.GetDescription(), unmarshaled.GetDescription())
	}
	if schedule.GetStatus() != unmarshaled.GetStatus() {
		t.Errorf("expected status %s, got %s", schedule.GetStatus(), unmarshaled.GetStatus())
	}
	if schedule.GetQueueName() != unmarshaled.GetQueueName() {
		t.Errorf("expected queue name %s, got %s", schedule.GetQueueName(), unmarshaled.GetQueueName())
	}
	if schedule.GetTaskDefinitionID() != unmarshaled.GetTaskDefinitionID() {
		t.Errorf("expected task definition ID %s, got %s", schedule.GetTaskDefinitionID(), unmarshaled.GetTaskDefinitionID())
	}
	if schedule.GetExecutionCount() != unmarshaled.GetExecutionCount() {
		t.Errorf("expected execution count %d, got %d", schedule.GetExecutionCount(), unmarshaled.GetExecutionCount())
	}
	if schedule.GetMaxExecutionCount() != unmarshaled.GetMaxExecutionCount() {
		t.Errorf("expected max execution count %d, got %d", schedule.GetMaxExecutionCount(), unmarshaled.GetMaxExecutionCount())
	}
	if unmarshaled.GetRecurrenceRule() == nil {
		t.Error("expected RecurrenceRule to not be nil")
	}
	if unmarshaled.GetRecurrenceRule().GetFrequency() != FrequencyDaily {
		t.Error("expected frequency to be Daily")
	}
}

func TestNewScheduleQuery(t *testing.T) {
	query := NewScheduleQuery()

	if query == nil {
		t.Fatal("expected query to not be nil")
	}
	if query.Limit() != 10 {
		t.Errorf("expected default limit 10, got %d", query.Limit())
	}
	if query.Offset() != 0 {
		t.Errorf("expected offset 0, got %d", query.Offset())
	}
}

func TestScheduleQueryGettersAndSetters(t *testing.T) {
	query := NewScheduleQuery()

	// Test ID
	query.SetID("test-id-123")
	if query.ID() != "test-id-123" {
		t.Errorf("expected ID 'test-id-123', got %s", query.ID())
	}

	// Test Name
	query.SetName("Test Schedule")
	if query.Name() != "Test Schedule" {
		t.Errorf("expected name 'Test Schedule', got %s", query.Name())
	}

	// Test Status
	query.SetStatus("active")
	if query.Status() != "active" {
		t.Errorf("expected status 'active', got %s", query.Status())
	}

	// Test QueueName
	query.SetQueueName("test-queue")
	if query.QueueName() != "test-queue" {
		t.Errorf("expected queue name 'test-queue', got %s", query.QueueName())
	}

	// Test TaskDefinitionID
	query.SetTaskDefinitionID("task-def-123")
	if query.TaskDefinitionID() != "task-def-123" {
		t.Errorf("expected task definition ID 'task-def-123', got %s", query.TaskDefinitionID())
	}

	// Test Limit
	query.SetLimit(50)
	if query.Limit() != 50 {
		t.Errorf("expected limit 50, got %d", query.Limit())
	}

	// Test Offset
	query.SetOffset(100)
	if query.Offset() != 100 {
		t.Errorf("expected offset 100, got %d", query.Offset())
	}
}

func TestScheduleQueryChaining(t *testing.T) {
	query := NewScheduleQuery().
		SetID("test-id").
		SetName("Test Name").
		SetStatus("active").
		SetQueueName("default").
		SetTaskDefinitionID("task-123").
		SetLimit(20).
		SetOffset(40)

	if query.ID() != "test-id" {
		t.Errorf("expected ID 'test-id', got %s", query.ID())
	}
	if query.Name() != "Test Name" {
		t.Errorf("expected name 'Test Name', got %s", query.Name())
	}
	if query.Status() != "active" {
		t.Errorf("expected status 'active', got %s", query.Status())
	}
	if query.QueueName() != "default" {
		t.Errorf("expected queue name 'default', got %s", query.QueueName())
	}
	if query.TaskDefinitionID() != "task-123" {
		t.Errorf("expected task definition ID 'task-123', got %s", query.TaskDefinitionID())
	}
	if query.Limit() != 20 {
		t.Errorf("expected limit 20, got %d", query.Limit())
	}
	if query.Offset() != 40 {
		t.Errorf("expected offset 40, got %d", query.Offset())
	}
}

func TestScheduleSetterChaining(t *testing.T) {
	schedule := NewSchedule().
		SetName("Chained Schedule").
		SetDescription("Chained Description").
		SetStatus("active").
		SetQueueName("default").
		SetTaskDefinitionID("task-123").
		SetExecutionCount(5).
		SetMaxExecutionCount(10)

	if schedule.GetName() != "Chained Schedule" {
		t.Errorf("expected name 'Chained Schedule', got %s", schedule.GetName())
	}
	if schedule.GetDescription() != "Chained Description" {
		t.Errorf("expected description 'Chained Description', got %s", schedule.GetDescription())
	}
	if schedule.GetStatus() != "active" {
		t.Errorf("expected status 'active', got %s", schedule.GetStatus())
	}
	if schedule.GetQueueName() != "default" {
		t.Errorf("expected queue name 'default', got %s", schedule.GetQueueName())
	}
	if schedule.GetTaskDefinitionID() != "task-123" {
		t.Errorf("expected task definition ID 'task-123', got %s", schedule.GetTaskDefinitionID())
	}
	if schedule.GetExecutionCount() != 5 {
		t.Errorf("expected execution count 5, got %d", schedule.GetExecutionCount())
	}
	if schedule.GetMaxExecutionCount() != 10 {
		t.Errorf("expected max execution count 10, got %d", schedule.GetMaxExecutionCount())
	}
}
