package taskstore

import (
	"encoding/json"
	"testing"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSchedule(t *testing.T) {
	schedule := NewSchedule()

	assert.NotNil(t, schedule)
	assert.NotEmpty(t, schedule.ID())
	assert.Equal(t, "draft", schedule.Status())
	assert.Equal(t, sb.NULL_DATETIME, schedule.StartAt())
	assert.Equal(t, sb.MAX_DATETIME, schedule.EndAt())
	assert.Equal(t, sb.NULL_DATETIME, schedule.LastRunAt())
	assert.Equal(t, sb.NULL_DATETIME, schedule.NextRunAt())
	assert.Equal(t, sb.MAX_DATETIME, schedule.SoftDeletedAt())
	assert.NotNil(t, schedule.RecurrenceRule())
	assert.NotEmpty(t, schedule.CreatedAt())
	assert.NotEmpty(t, schedule.UpdatedAt())
}

func TestScheduleGettersAndSetters(t *testing.T) {
	schedule := NewSchedule()

	// Test ID
	schedule.SetID("test-id-123")
	assert.Equal(t, "test-id-123", schedule.ID())

	// Test Name
	schedule.SetName("Test Schedule")
	assert.Equal(t, "Test Schedule", schedule.Name())

	// Test Description
	schedule.SetDescription("Test Description")
	assert.Equal(t, "Test Description", schedule.Description())

	// Test Status
	schedule.SetStatus("active")
	assert.Equal(t, "active", schedule.Status())

	// Test QueueName
	schedule.SetQueueName("test-queue")
	assert.Equal(t, "test-queue", schedule.QueueName())

	// Test TaskDefinitionID
	schedule.SetTaskDefinitionID("task-def-123")
	assert.Equal(t, "task-def-123", schedule.TaskDefinitionID())

	// Test TaskParameters
	params := map[string]any{"key": "value", "count": 42}
	schedule.SetTaskParameters(params)
	assert.Equal(t, params, schedule.TaskParameters())

	// Test StartAt
	startAt := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	schedule.SetStartAt(startAt)
	assert.Equal(t, startAt, schedule.StartAt())

	// Test EndAt
	endAt := carbon.Now(carbon.UTC).AddDays(30).ToDateTimeString(carbon.UTC)
	schedule.SetEndAt(endAt)
	assert.Equal(t, endAt, schedule.EndAt())

	// Test ExecutionCount
	schedule.SetExecutionCount(5)
	assert.Equal(t, 5, schedule.ExecutionCount())

	// Test MaxExecutionCount
	schedule.SetMaxExecutionCount(10)
	assert.Equal(t, 10, schedule.MaxExecutionCount())

	// Test LastRunAt
	lastRunAt := carbon.Now(carbon.UTC).AddMinutes(-5).ToDateTimeString(carbon.UTC)
	schedule.SetLastRunAt(lastRunAt)
	assert.Equal(t, lastRunAt, schedule.LastRunAt())

	// Test NextRunAt
	nextRunAt := carbon.Now(carbon.UTC).AddMinutes(5).ToDateTimeString(carbon.UTC)
	schedule.SetNextRunAt(nextRunAt)
	assert.Equal(t, nextRunAt, schedule.NextRunAt())

	// Test CreatedAt
	createdAt := carbon.Now(carbon.UTC).AddDays(-1).ToDateTimeString(carbon.UTC)
	schedule.SetCreatedAt(createdAt)
	assert.Equal(t, createdAt, schedule.CreatedAt())

	// Test UpdatedAt
	updatedAt := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	schedule.SetUpdatedAt(updatedAt)
	assert.Equal(t, updatedAt, schedule.UpdatedAt())

	// Test SoftDeletedAt
	softDeletedAt := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	schedule.SetSoftDeletedAt(softDeletedAt)
	assert.Equal(t, softDeletedAt, schedule.SoftDeletedAt())

	// Test RecurrenceRule
	rr := NewRecurrenceRule()
	rr.SetFrequency(FrequencyDaily)
	schedule.SetRecurrenceRule(rr)
	assert.Equal(t, FrequencyDaily, schedule.RecurrenceRule().GetFrequency())
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
	require.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)

	// Unmarshal from JSON - create a new schedule and unmarshal into it
	unmarshaled := NewSchedule()
	err = json.Unmarshal(jsonBytes, unmarshaled)
	require.NoError(t, err)

	// Verify fields
	assert.Equal(t, schedule.ID(), unmarshaled.ID())
	assert.Equal(t, schedule.Name(), unmarshaled.Name())
	assert.Equal(t, schedule.Description(), unmarshaled.Description())
	assert.Equal(t, schedule.Status(), unmarshaled.Status())
	assert.Equal(t, schedule.QueueName(), unmarshaled.QueueName())
	assert.Equal(t, schedule.TaskDefinitionID(), unmarshaled.TaskDefinitionID())
	assert.Equal(t, schedule.ExecutionCount(), unmarshaled.ExecutionCount())
	assert.Equal(t, schedule.MaxExecutionCount(), unmarshaled.MaxExecutionCount())
	assert.NotNil(t, unmarshaled.RecurrenceRule())
	assert.Equal(t, FrequencyDaily, unmarshaled.RecurrenceRule().GetFrequency())
}

func TestNewScheduleQuery(t *testing.T) {
	query := NewScheduleQuery()

	assert.NotNil(t, query)
	assert.Equal(t, 10, query.Limit()) // Default limit
	assert.Equal(t, 0, query.Offset())
}

func TestScheduleQueryGettersAndSetters(t *testing.T) {
	query := NewScheduleQuery()

	// Test ID
	query.SetID("test-id-123")
	assert.Equal(t, "test-id-123", query.ID())

	// Test Name
	query.SetName("Test Schedule")
	assert.Equal(t, "Test Schedule", query.Name())

	// Test Status
	query.SetStatus("active")
	assert.Equal(t, "active", query.Status())

	// Test QueueName
	query.SetQueueName("test-queue")
	assert.Equal(t, "test-queue", query.QueueName())

	// Test TaskDefinitionID
	query.SetTaskDefinitionID("task-def-123")
	assert.Equal(t, "task-def-123", query.TaskDefinitionID())

	// Test Limit
	query.SetLimit(50)
	assert.Equal(t, 50, query.Limit())

	// Test Offset
	query.SetOffset(100)
	assert.Equal(t, 100, query.Offset())
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

	assert.Equal(t, "test-id", query.ID())
	assert.Equal(t, "Test Name", query.Name())
	assert.Equal(t, "active", query.Status())
	assert.Equal(t, "default", query.QueueName())
	assert.Equal(t, "task-123", query.TaskDefinitionID())
	assert.Equal(t, 20, query.Limit())
	assert.Equal(t, 40, query.Offset())
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

	assert.Equal(t, "Chained Schedule", schedule.Name())
	assert.Equal(t, "Chained Description", schedule.Description())
	assert.Equal(t, "active", schedule.Status())
	assert.Equal(t, "default", schedule.QueueName())
	assert.Equal(t, "task-123", schedule.TaskDefinitionID())
	assert.Equal(t, 5, schedule.ExecutionCount())
	assert.Equal(t, 10, schedule.MaxExecutionCount())
}
