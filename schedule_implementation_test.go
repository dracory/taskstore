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
	assert.NotEmpty(t, schedule.GetID())
	assert.Equal(t, "draft", schedule.GetStatus())
	assert.Equal(t, sb.NULL_DATETIME, schedule.GetStartAt())
	assert.Equal(t, sb.MAX_DATETIME, schedule.GetEndAt())
	assert.Equal(t, sb.NULL_DATETIME, schedule.GetLastRunAt())
	assert.Equal(t, sb.NULL_DATETIME, schedule.GetNextRunAt())
	assert.Equal(t, sb.MAX_DATETIME, schedule.GetSoftDeletedAt())
	assert.NotNil(t, schedule.GetRecurrenceRule())
	assert.NotEmpty(t, schedule.GetCreatedAt())
	assert.NotEmpty(t, schedule.GetUpdatedAt())
}

func TestScheduleGettersAndSetters(t *testing.T) {
	schedule := NewSchedule()

	// Test ID
	schedule.SetID("test-id-123")
	assert.Equal(t, "test-id-123", schedule.GetID())

	// Test Name
	schedule.SetName("Test Schedule")
	assert.Equal(t, "Test Schedule", schedule.GetName())

	// Test Description
	schedule.SetDescription("Test Description")
	assert.Equal(t, "Test Description", schedule.GetDescription())

	// Test Status
	schedule.SetStatus("active")
	assert.Equal(t, "active", schedule.GetStatus())

	// Test QueueName
	schedule.SetQueueName("test-queue")
	assert.Equal(t, "test-queue", schedule.GetQueueName())

	// Test TaskDefinitionID
	schedule.SetTaskDefinitionID("task-def-123")
	assert.Equal(t, "task-def-123", schedule.GetTaskDefinitionID())

	// Test TaskParameters
	params := map[string]any{"key": "value", "count": 42}
	schedule.SetTaskParameters(params)
	assert.Equal(t, params, schedule.GetTaskParameters())

	// Test StartAt
	startAt := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	schedule.SetStartAt(startAt)
	assert.Equal(t, startAt, schedule.GetStartAt())

	// Test EndAt
	endAt := carbon.Now(carbon.UTC).AddDays(30).ToDateTimeString(carbon.UTC)
	schedule.SetEndAt(endAt)
	assert.Equal(t, endAt, schedule.GetEndAt())

	// Test ExecutionCount
	schedule.SetExecutionCount(5)
	assert.Equal(t, 5, schedule.GetExecutionCount())

	// Test MaxExecutionCount
	schedule.SetMaxExecutionCount(10)
	assert.Equal(t, 10, schedule.GetMaxExecutionCount())

	// Test LastRunAt
	lastRunAt := carbon.Now(carbon.UTC).AddMinutes(-5).ToDateTimeString(carbon.UTC)
	schedule.SetLastRunAt(lastRunAt)
	assert.Equal(t, lastRunAt, schedule.GetLastRunAt())

	// Test NextRunAt
	nextRunAt := carbon.Now(carbon.UTC).AddMinutes(5).ToDateTimeString(carbon.UTC)
	schedule.SetNextRunAt(nextRunAt)
	assert.Equal(t, nextRunAt, schedule.GetNextRunAt())

	// Test CreatedAt
	createdAt := carbon.Now(carbon.UTC).AddDays(-1).ToDateTimeString(carbon.UTC)
	schedule.SetCreatedAt(createdAt)
	assert.Equal(t, createdAt, schedule.GetCreatedAt())

	// Test UpdatedAt
	updatedAt := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	schedule.SetUpdatedAt(updatedAt)
	assert.Equal(t, updatedAt, schedule.GetUpdatedAt())

	// Test SoftDeletedAt
	softDeletedAt := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	schedule.SetSoftDeletedAt(softDeletedAt)
	assert.Equal(t, softDeletedAt, schedule.GetSoftDeletedAt())

	// Test RecurrenceRule
	rr := NewRecurrenceRule()
	rr.SetFrequency(FrequencyDaily)
	schedule.SetRecurrenceRule(rr)
	assert.Equal(t, FrequencyDaily, schedule.GetRecurrenceRule().GetFrequency())
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
	assert.Equal(t, schedule.GetID(), unmarshaled.GetID())
	assert.Equal(t, schedule.GetName(), unmarshaled.GetName())
	assert.Equal(t, schedule.GetDescription(), unmarshaled.GetDescription())
	assert.Equal(t, schedule.GetStatus(), unmarshaled.GetStatus())
	assert.Equal(t, schedule.GetQueueName(), unmarshaled.GetQueueName())
	assert.Equal(t, schedule.GetTaskDefinitionID(), unmarshaled.GetTaskDefinitionID())
	assert.Equal(t, schedule.GetExecutionCount(), unmarshaled.GetExecutionCount())
	assert.Equal(t, schedule.GetMaxExecutionCount(), unmarshaled.GetMaxExecutionCount())
	assert.NotNil(t, unmarshaled.GetRecurrenceRule())
	assert.Equal(t, FrequencyDaily, unmarshaled.GetRecurrenceRule().GetFrequency())
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

	assert.Equal(t, "Chained Schedule", schedule.GetName())
	assert.Equal(t, "Chained Description", schedule.GetDescription())
	assert.Equal(t, "active", schedule.GetStatus())
	assert.Equal(t, "default", schedule.GetQueueName())
	assert.Equal(t, "task-123", schedule.GetTaskDefinitionID())
	assert.Equal(t, 5, schedule.GetExecutionCount())
	assert.Equal(t, 10, schedule.GetMaxExecutionCount())
}
