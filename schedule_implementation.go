package taskstore

import (
	"encoding/json"

	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

// scheduleImplementation is the concrete implementation of ScheduleInterface.
// It stores all schedule fields including timing, recurrence, and metadata.
type scheduleImplementation struct {
	id                string
	name              string
	description       string
	status            string
	recurrenceRule    RecurrenceRuleInterface
	queueName         string
	taskDefinitionID  string
	taskParameters    map[string]any
	startAt           string
	endAt             string
	executionCount    int
	maxExecutionCount int
	lastRunAt         string
	nextRunAt         string
	createdAt         string
	updatedAt         string
	softDeletedAt     string
}

var _ ScheduleInterface = (*scheduleImplementation)(nil)

// NewSchedule creates a new schedule with default values and a new recurrence rule.
func NewSchedule() ScheduleInterface {
	return &scheduleImplementation{
		id:             uid.HumanUid(),
		status:         "draft",
		recurrenceRule: NewRecurrenceRule(),
		startAt:        sb.NULL_DATETIME,
		endAt:          sb.MAX_DATETIME,
		lastRunAt:      sb.NULL_DATETIME,
		nextRunAt:      sb.NULL_DATETIME,
		createdAt:      carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC),
		updatedAt:      carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC),
		softDeletedAt:  sb.MAX_DATETIME,
	}
}

// GetID returns the unique identifier of the schedule.
func (s *scheduleImplementation) GetID() string {
	return s.id
}

// SetID sets the unique identifier of the schedule.
func (s *scheduleImplementation) SetID(id string) ScheduleInterface {
	s.id = id
	return s
}

// GetName returns the name of the schedule.
func (s *scheduleImplementation) GetName() string {
	return s.name
}

// SetName sets the name of the schedule.
func (s *scheduleImplementation) SetName(name string) ScheduleInterface {
	s.name = name
	return s
}

// Description returns the description of the schedule.
func (s *scheduleImplementation) GetDescription() string {
	return s.description
}

// SetDescription sets the description of the schedule.
func (s *scheduleImplementation) SetDescription(description string) ScheduleInterface {
	s.description = description
	return s
}

// Status returns the status of the schedule.
func (s *scheduleImplementation) GetStatus() string {
	return s.status
}

// SetStatus sets the status of the schedule.
func (s *scheduleImplementation) SetStatus(status string) ScheduleInterface {
	s.status = status
	return s
}

// RecurrenceRule returns the recurrence rule that defines when the schedule should run.
func (s *scheduleImplementation) GetRecurrenceRule() RecurrenceRuleInterface {
	return s.recurrenceRule
}

// SetRecurrenceRule sets the recurrence rule that defines when the schedule should run.
func (s *scheduleImplementation) SetRecurrenceRule(rule RecurrenceRuleInterface) ScheduleInterface {
	s.recurrenceRule = rule
	return s
}

// QueueName returns the name of the queue that this schedule is associated with.
func (s *scheduleImplementation) GetQueueName() string {
	return s.queueName
}

// SetQueueName sets the name of the queue that this schedule is associated with.
func (s *scheduleImplementation) SetQueueName(queueName string) ScheduleInterface {
	s.queueName = queueName
	return s
}

// TaskDefinitionID returns the unique identifier of the task definition that this schedule is associated with.
func (s *scheduleImplementation) GetTaskDefinitionID() string {
	return s.taskDefinitionID
}

// SetTaskDefinitionID sets the unique identifier of the task definition that this schedule is associated with.
func (s *scheduleImplementation) SetTaskDefinitionID(taskDefinitionID string) ScheduleInterface {
	s.taskDefinitionID = taskDefinitionID
	return s
}

// TaskParameters returns the parameters to be passed to the task definition when it is executed.
func (s *scheduleImplementation) GetTaskParameters() map[string]any {
	return s.taskParameters
}

// SetTaskParameters sets the parameters to be passed to the task definition when it is executed.
func (s *scheduleImplementation) SetTaskParameters(parameters map[string]any) ScheduleInterface {
	s.taskParameters = parameters
	return s
}

// StartAt returns the start date and time of the schedule.
func (s *scheduleImplementation) GetStartAt() string {
	return s.startAt
}

// SetStartAt sets the start date and time of the schedule.
// If startAt is not set, the schedule will start at the current time.
func (s *scheduleImplementation) SetStartAt(startAt string) ScheduleInterface {
	s.startAt = startAt
	return s
}

// EndAt returns the end date and time of the schedule.
// The default value is the maximum datetime (never expires).
func (s *scheduleImplementation) GetEndAt() string {
	return s.endAt
}

// SetEndAt sets the end date and time of the schedule.
func (s *scheduleImplementation) SetEndAt(endAt string) ScheduleInterface {
	s.endAt = endAt
	return s
}

// ExecutionCount returns the number of times the schedule has been executed.
func (s *scheduleImplementation) GetExecutionCount() int {
	return s.executionCount
}

// SetExecutionCount sets the number of times the schedule has been executed.
func (s *scheduleImplementation) SetExecutionCount(count int) ScheduleInterface {
	s.executionCount = count
	return s
}

// MaxExecutionCount returns the maximum number of times the schedule is allowed to be executed.
// The default value is int max (no limit). To execute only once, set maxExecutionCount to 1.
func (s *scheduleImplementation) GetMaxExecutionCount() int {
	return s.maxExecutionCount
}

// SetMaxExecutionCount sets the maximum number of times the schedule is allowed to be executed.
func (s *scheduleImplementation) SetMaxExecutionCount(count int) ScheduleInterface {
	s.maxExecutionCount = count
	return s
}

// LastRunAt returns the last date and time the schedule was executed.
func (s *scheduleImplementation) GetLastRunAt() string {
	return s.lastRunAt
}

// SetLastRunAt sets the last date and time the schedule was executed.
func (s *scheduleImplementation) SetLastRunAt(lastRunAt string) ScheduleInterface {
	s.lastRunAt = lastRunAt
	return s
}

// NextRunAt returns the next date and time the schedule is scheduled to run.
func (s *scheduleImplementation) GetNextRunAt() string {
	return s.nextRunAt
}

// SetNextRunAt sets the next date and time the schedule is scheduled to run.
func (s *scheduleImplementation) SetNextRunAt(nextRunAt string) ScheduleInterface {
	s.nextRunAt = nextRunAt
	return s
}

// CreatedAt returns the date and time the schedule was created.
func (s *scheduleImplementation) GetCreatedAt() string {
	return s.createdAt
}

// SetCreatedAt sets the date and time the schedule was created.
func (s *scheduleImplementation) SetCreatedAt(createdAt string) ScheduleInterface {
	s.createdAt = createdAt
	return s
}

// UpdatedAt returns the date and time the schedule was last updated.
func (s *scheduleImplementation) GetUpdatedAt() string {
	return s.updatedAt
}

// SetUpdatedAt sets the date and time the schedule was last updated.
func (s *scheduleImplementation) SetUpdatedAt(updatedAt string) ScheduleInterface {
	s.updatedAt = updatedAt
	return s
}

// SoftDeletedAt returns the date and time the schedule was soft deleted.
// The default value is max datetime (not soft deleted, 9999-12-31 23:59:59).
// To soft delete a schedule, set softDeletedAt to the current time.
// To unsoft delete a schedule, set softDeletedAt to max datetime.
// A soft deleted schedule is when its in the past.
func (s *scheduleImplementation) GetSoftDeletedAt() string {
	return s.softDeletedAt
}

// SetSoftDeletedAt sets the date and time the schedule was soft deleted.
func (s *scheduleImplementation) SetSoftDeletedAt(softDeletedAt string) ScheduleInterface {
	s.softDeletedAt = softDeletedAt
	return s
}

// HasReachedEndDate returns true if the schedule has reached its end date
func (s *scheduleImplementation) HasReachedEndDate() bool {
	endAt := carbon.Parse(s.endAt, carbon.UTC)
	now := carbon.Now(carbon.UTC)

	return now.Gt(endAt)
}

// HasReachedMaxExecutions returns true if the schedule has reached its maximum number of executions
func (s *scheduleImplementation) HasReachedMaxExecutions() bool {
	if s.maxExecutionCount <= 0 {
		return false
	}

	return s.executionCount >= s.maxExecutionCount
}

// GetNextOccurrence returns the next occurrence of the schedule
// if invalid recurrence rule, returns error
func (s *scheduleImplementation) GetNextOccurrence() (string, error) {
	now := carbon.Now(carbon.UTC)

	next, err := NextRunAt(s.recurrenceRule, now)
	if err != nil {
		return s.nextRunAt, err
	}

	return next.ToDateTimeString(carbon.UTC), nil
}

// IncrementExecutionCount increments the execution count of the schedule by one
func (s *scheduleImplementation) IncrementExecutionCount() ScheduleInterface {
	s.executionCount++
	return s
}

// UpdateNextRunAt calculates the next run at of the schedule and updates it
func (s *scheduleImplementation) UpdateNextRunAt() ScheduleInterface {
	next, err := s.GetNextOccurrence()
	if err != nil {
		return s
	}

	s.nextRunAt = next
	return s
}

// UpdateLastRunAt updates the last run at of the schedule with current time
func (s *scheduleImplementation) UpdateLastRunAt() ScheduleInterface {
	s.lastRunAt = carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	return s
}

// IsDue returns true if the schedule is due to run
func (s *scheduleImplementation) IsDue() bool {
	nextRunAt := carbon.Parse(s.nextRunAt, carbon.UTC)
	now := carbon.Now(carbon.UTC)

	return nextRunAt.Lt(now) || nextRunAt.Eq(now)
}

// MarshalJSON custom marshaler to handle all fields including RecurrenceRule interface
func (s *scheduleImplementation) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":                s.id,
		"name":              s.name,
		"description":       s.description,
		"status":            s.status,
		"recurrenceRule":    s.recurrenceRule,
		"queueName":         s.queueName,
		"taskDefinitionID":  s.taskDefinitionID,
		"taskParameters":    s.taskParameters,
		"startAt":           s.startAt,
		"endAt":             s.endAt,
		"executionCount":    s.executionCount,
		"maxExecutionCount": s.maxExecutionCount,
		"lastRunAt":         s.lastRunAt,
		"nextRunAt":         s.nextRunAt,
		"createdAt":         s.createdAt,
		"updatedAt":         s.updatedAt,
		"softDeletedAt":     s.softDeletedAt,
	})
}

// UnmarshalJSON custom unmarshaler to handle RecurrenceRule interface
func (s *scheduleImplementation) UnmarshalJSON(data []byte) error {
	var aux struct {
		ID                string          `json:"id"`
		Name              string          `json:"name"`
		Description       string          `json:"description"`
		Status            string          `json:"status"`
		RecurrenceRule    json.RawMessage `json:"recurrenceRule"`
		QueueName         string          `json:"queueName"`
		TaskDefinitionID  string          `json:"taskDefinitionID"`
		TaskParameters    map[string]any  `json:"taskParameters"`
		StartAt           string          `json:"startAt"`
		EndAt             string          `json:"endAt"`
		ExecutionCount    int             `json:"executionCount"`
		MaxExecutionCount int             `json:"maxExecutionCount"`
		LastRunAt         string          `json:"lastRunAt"`
		NextRunAt         string          `json:"nextRunAt"`
		CreatedAt         string          `json:"createdAt"`
		UpdatedAt         string          `json:"updatedAt"`
		SoftDeletedAt     string          `json:"softDeletedAt"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	s.id = aux.ID
	s.name = aux.Name
	s.description = aux.Description
	s.status = aux.Status
	s.queueName = aux.QueueName
	s.taskDefinitionID = aux.TaskDefinitionID
	s.taskParameters = aux.TaskParameters
	s.startAt = aux.StartAt
	s.endAt = aux.EndAt
	s.executionCount = aux.ExecutionCount
	s.maxExecutionCount = aux.MaxExecutionCount
	s.lastRunAt = aux.LastRunAt
	s.nextRunAt = aux.NextRunAt
	s.createdAt = aux.CreatedAt
	s.updatedAt = aux.UpdatedAt
	s.softDeletedAt = aux.SoftDeletedAt

	if aux.RecurrenceRule != nil {
		rule := NewRecurrenceRule()
		if err := json.Unmarshal(aux.RecurrenceRule, rule); err != nil {
			return err
		}
		s.recurrenceRule = rule
	}

	return nil
}

type ScheduleQuery struct {
	id               string
	name             string
	status           string
	queueName        string
	taskDefinitionID string
	limit            int
	offset           int
}

var _ ScheduleQueryInterface = (*ScheduleQuery)(nil)

func NewScheduleQuery() ScheduleQueryInterface {
	return &ScheduleQuery{
		limit: 10,
	}
}

func (q *ScheduleQuery) ID() string {
	return q.id
}

func (q *ScheduleQuery) SetID(id string) ScheduleQueryInterface {
	q.id = id
	return q
}

func (q *ScheduleQuery) Name() string {
	return q.name
}

func (q *ScheduleQuery) SetName(name string) ScheduleQueryInterface {
	q.name = name
	return q
}

func (q *ScheduleQuery) Status() string {
	return q.status
}

func (q *ScheduleQuery) SetStatus(status string) ScheduleQueryInterface {
	q.status = status
	return q
}

func (q *ScheduleQuery) QueueName() string {
	return q.queueName
}

func (q *ScheduleQuery) SetQueueName(queueName string) ScheduleQueryInterface {
	q.queueName = queueName
	return q
}

func (q *ScheduleQuery) TaskDefinitionID() string {
	return q.taskDefinitionID
}

func (q *ScheduleQuery) SetTaskDefinitionID(taskDefinitionID string) ScheduleQueryInterface {
	q.taskDefinitionID = taskDefinitionID
	return q
}

func (q *ScheduleQuery) Limit() int {
	return q.limit
}

func (q *ScheduleQuery) SetLimit(limit int) ScheduleQueryInterface {
	q.limit = limit
	return q
}

func (q *ScheduleQuery) Offset() int {
	return q.offset
}

func (q *ScheduleQuery) SetOffset(offset int) ScheduleQueryInterface {
	q.offset = offset
	return q
}
