package taskstore

import (
	"encoding/json"
	"time"

	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
)

// ScheduleInterface defines the contract for a schedule, including its
// identity, metadata, recurrence rule, timing fields, execution limits,
// soft-delete semantics, and helper methods for evaluating schedule state.
type ScheduleInterface interface {
	// GetID the unique identifier of the schedule
	GetID() string

	// SetID sets the unique identifier of the schedule
	SetID(string) ScheduleInterface

	// GetName the name of the schedule
	GetName() string

	// SetName sets the name of the schedule
	SetName(string) ScheduleInterface

	// GetDescription the description of the schedule
	GetDescription() string

	// SetDescription sets the description of the schedule
	SetDescription(string) ScheduleInterface

	// GetStatus the status of the schedule
	// Valid values are "draft" (default), "active", "inactive"
	GetStatus() string

	// SetStatus sets the status of the schedule
	SetStatus(string) ScheduleInterface

	// GetRecurrenceRule the recurrence rule that defines when the schedule should run
	GetRecurrenceRule() RecurrenceRuleInterface

	// SetRecurrenceRule sets the recurrence rule that defines when the schedule should run
	SetRecurrenceRule(RecurrenceRuleInterface) ScheduleInterface

	// GetQueueName the name of the queue that this schedule is associated with
	GetQueueName() string

	// SetQueueName sets the name of the queue that this schedule is associated with
	SetQueueName(string) ScheduleInterface

	// GetTaskDefinitionID the unique identifier of the task definition
	// that this schedule is associated with
	GetTaskDefinitionID() string

	// SetTaskDefinitionID sets the unique identifier of the task definition
	// that this schedule is associated with
	SetTaskDefinitionID(string) ScheduleInterface

	// GetTaskParameters the parameters to be passed to the task definition
	// when it is executed
	GetTaskParameters() map[string]any

	// SetTaskParameters sets the parameters to be passed to the task definition
	// when it is executed
	SetTaskParameters(map[string]any) ScheduleInterface

	// GetStartAt the start date and time of the schedule
	GetStartAt() string

	// SetStartAt sets the start date and time of the schedule
	// If startAt is not set, the schedule will start at the current time
	SetStartAt(string) ScheduleInterface

	// GetEndAt the end date and time of the schedule
	// The default value is the maximum datetime (never expires)
	GetEndAt() string

	// SetEndAt sets the end date and time of the schedule
	SetEndAt(string) ScheduleInterface

	// GetExecutionCount the number of times the schedule has been executed
	GetExecutionCount() int

	// SetExecutionCount sets the number of times the schedule has been executed
	SetExecutionCount(int) ScheduleInterface

	// GetMaxExecutionCount the maximum number of times the schedule is allowed to be executed
	// The default value is int max (no limit)
	// To execute only once, set maxExecutionCount to 1
	GetMaxExecutionCount() int

	// SetMaxExecutionCount sets the maximum number of times the schedule is allowed to be executed
	SetMaxExecutionCount(int) ScheduleInterface

	// GetLastRunAt the last date and time the schedule was executed
	GetLastRunAt() string

	// SetLastRunAt sets the last date and time the schedule was executed
	SetLastRunAt(string) ScheduleInterface

	// GetNextRunAt the next date and time the schedule is scheduled to run
	GetNextRunAt() string

	// SetNextRunAt sets the next date and time the schedule is scheduled to run
	SetNextRunAt(string) ScheduleInterface

	// GetCreatedAt the date and time the schedule was created
	GetCreatedAt() time.Time

	// GetCreatedAtCarbon returns the created at time as a carbon object
	GetCreatedAtCarbon() *carbon.Carbon

	// SetCreatedAt sets the date and time the schedule was created
	SetCreatedAt(time.Time) ScheduleInterface

	// GetUpdatedAt the date and time the schedule was last updated
	GetUpdatedAt() time.Time

	// GetUpdatedAtCarbon returns the updated at time as a carbon object
	GetUpdatedAtCarbon() *carbon.Carbon

	// SetUpdatedAt sets the date and time the schedule was last updated
	SetUpdatedAt(time.Time) ScheduleInterface

	// GetSoftDeletedAt the date and time the schedule was soft deleted
	// The default value is max datetime (not soft deleted, 9999-12-31 23:59:59)
	// To soft delete a schedule, set softDeletedAt to the current time
	// To unsoft delete a schedule, set softDeletedAt to max datetime
	// A soft deleted schedule is when its in the past
	GetSoftDeletedAt() time.Time

	// GetSoftDeletedAtCarbon returns the soft deleted at time as a carbon object
	GetSoftDeletedAtCarbon() *carbon.Carbon

	// SetSoftDeletedAt sets the date and time the schedule was soft deleted
	SetSoftDeletedAt(time.Time) ScheduleInterface

	// HasReachedEndDate returns true if the schedule has reached its end date
	HasReachedEndDate() bool

	// HasReachedMaxExecutions returns true if the schedule has reached its maximum number of executions
	HasReachedMaxExecutions() bool

	// GetNextOccurrence returns the next occurrence of the schedule
	// if invalid recurrence rule, returns error
	GetNextOccurrence() (string, error)

	// IncrementExecutionCount increments the execution count of the schedule by one
	IncrementExecutionCount() ScheduleInterface

	// UpdateNextRunAt calculates the next run at of the schedule and updates it
	UpdateNextRunAt() ScheduleInterface

	// UpdateLastRunAt updates the last run at of the schedule with current time
	UpdateLastRunAt() ScheduleInterface

	// IsDue returns true if the schedule is due to run
	IsDue() bool
}

// ScheduleQueryInterface defines the query parameters used to filter and
// paginate schedules when listing or counting them.
type ScheduleQueryInterface interface {
	// ID the unique identifier of the schedule to filter by
	ID() string

	// SetID sets the unique identifier of the schedule to filter by
	SetID(string) ScheduleQueryInterface

	// Name the name of the schedule to filter by
	Name() string

	// SetName sets the name of the schedule to filter by
	SetName(string) ScheduleQueryInterface

	// Status the status of the schedule to filter by
	Status() string

	// SetStatus sets the status of the schedule to filter by
	SetStatus(string) ScheduleQueryInterface

	// QueueName the name of the queue that schedules are associated with to filter by
	QueueName() string

	// SetQueueName sets the name of the queue that schedules are associated with to filter by
	SetQueueName(string) ScheduleQueryInterface

	// TaskDefinitionID the unique identifier of the task definition that schedules are associated with to filter by
	TaskDefinitionID() string

	// SetTaskDefinitionID sets the unique identifier of the task definition that schedules are associated with to filter by
	SetTaskDefinitionID(string) ScheduleQueryInterface

	// Limit the maximum number of schedules to return
	Limit() int

	// SetLimit sets the maximum number of schedules to return
	SetLimit(int) ScheduleQueryInterface

	// Offset the number of schedules to skip before starting to return results
	Offset() int

	// SetOffset sets the number of schedules to skip before starting to return results
	SetOffset(int) ScheduleQueryInterface
}

// scheduleImplementation is the concrete implementation of ScheduleInterface.
// It stores all schedule fields including timing, recurrence, and metadata.
type scheduleImplementation struct {
	orm.ShortID

	NameField              string `db:"name"`
	DescriptionField       string `db:"description"`
	StatusField            string `db:"status"`
	RecurrenceRuleField    string `db:"recurrence_rule"`
	QueueNameField         string `db:"queue_name"`
	TaskDefinitionIDField  string `db:"task_definition_id"`
	ParametersField        string `db:"parameters"`
	StartAtField           string `db:"start_at"`
	EndAtField             string `db:"end_at"`
	ExecutionCountField    int    `db:"execution_count"`
	MaxExecutionCountField int    `db:"max_execution_count"`
	LastRunAtField         string `db:"last_run_at"`
	NextRunAtField         string `db:"next_run_at"`

	CreatedAtField orm.CreatedAt
	UpdatedAtField orm.UpdatedAt
	soft_delete.SoftDeletesMaxDate

	// cached recurrence rule to allow mutation via GetRecurrenceRule()
	cachedRecurrenceRule RecurrenceRuleInterface
}

var _ ScheduleInterface = (*scheduleImplementation)(nil)

// NewSchedule creates a new schedule with default values and a new recurrence rule.
func NewSchedule() ScheduleInterface {
	o := &scheduleImplementation{}
	o.SetID(neatuid.GenerateShortID())
	o.SetStatus("draft")
	o.SetRecurrenceRule(NewRecurrenceRule())
	o.SetStartAt(NULL_DATETIME)
	o.SetEndAt(MAX_DATETIME)
	o.SetLastRunAt(NULL_DATETIME)
	o.SetNextRunAt(NULL_DATETIME)
	o.SetExecutionCount(0)
	o.SetMaxExecutionCount(0)
	o.SetCreatedAt(carbon.Now(carbon.UTC).StdTime())
	o.SetUpdatedAt(carbon.Now(carbon.UTC).StdTime())
	o.SetSoftDeletedAt(carbon.Parse(MAX_DATETIME, carbon.UTC).StdTime())
	return o
}

// GetID returns the unique identifier of the schedule.
func (s *scheduleImplementation) GetID() string {
	return s.ShortID.ID
}

// SetID sets the unique identifier of the schedule.
func (s *scheduleImplementation) SetID(id string) ScheduleInterface {
	s.ShortID.ID = id
	return s
}

// GetName returns the name of the schedule.
func (s *scheduleImplementation) GetName() string {
	return s.NameField
}

// SetName sets the name of the schedule.
func (s *scheduleImplementation) SetName(name string) ScheduleInterface {
	s.NameField = name
	return s
}

// GetDescription returns the description of the schedule.
func (s *scheduleImplementation) GetDescription() string {
	return s.DescriptionField
}

// SetDescription sets the description of the schedule.
func (s *scheduleImplementation) SetDescription(description string) ScheduleInterface {
	s.DescriptionField = description
	return s
}

// GetStatus returns the status of the schedule.
func (s *scheduleImplementation) GetStatus() string {
	return s.StatusField
}

// SetStatus sets the status of the schedule.
func (s *scheduleImplementation) SetStatus(status string) ScheduleInterface {
	s.StatusField = status
	return s
}

// GetRecurrenceRule returns the recurrence rule that defines when the schedule should run.
// The returned instance is cached and mutable; changes made to it will be reflected
// in subsequent calls to GetRecurrenceRule and when the schedule is persisted.
func (s *scheduleImplementation) GetRecurrenceRule() RecurrenceRuleInterface {
	if s.cachedRecurrenceRule != nil {
		return s.cachedRecurrenceRule
	}
	if s.RecurrenceRuleField == "" {
		s.cachedRecurrenceRule = NewRecurrenceRule()
		return s.cachedRecurrenceRule
	}
	rr := NewRecurrenceRule()
	_ = json.Unmarshal([]byte(s.RecurrenceRuleField), rr)
	s.cachedRecurrenceRule = rr
	return s.cachedRecurrenceRule
}

// SetRecurrenceRule sets the recurrence rule that defines when the schedule should run.
func (s *scheduleImplementation) SetRecurrenceRule(rule RecurrenceRuleInterface) ScheduleInterface {
	if rule == nil {
		s.RecurrenceRuleField = ""
		s.cachedRecurrenceRule = nil
		return s
	}
	bytes, err := json.Marshal(rule)
	if err != nil {
		s.RecurrenceRuleField = ""
		s.cachedRecurrenceRule = nil
		return s
	}
	s.RecurrenceRuleField = string(bytes)
	s.cachedRecurrenceRule = rule
	return s
}

// GetQueueName returns the name of the queue that this schedule is associated with.
func (s *scheduleImplementation) GetQueueName() string {
	return s.QueueNameField
}

// SetQueueName sets the name of the queue that this schedule is associated with.
func (s *scheduleImplementation) SetQueueName(queueName string) ScheduleInterface {
	s.QueueNameField = queueName
	return s
}

// GetTaskDefinitionID returns the unique identifier of the task definition that this schedule is associated with.
func (s *scheduleImplementation) GetTaskDefinitionID() string {
	return s.TaskDefinitionIDField
}

// SetTaskDefinitionID sets the unique identifier of the task definition that this schedule is associated with.
func (s *scheduleImplementation) SetTaskDefinitionID(taskDefinitionID string) ScheduleInterface {
	s.TaskDefinitionIDField = taskDefinitionID
	return s
}

// GetTaskParameters returns the parameters to be passed to the task definition when it is executed.
func (s *scheduleImplementation) GetTaskParameters() map[string]any {
	if s.ParametersField == "" {
		return map[string]any{}
	}
	var params map[string]any
	_ = json.Unmarshal([]byte(s.ParametersField), &params)
	return params
}

// SetTaskParameters sets the parameters to be passed to the task definition when it is executed.
func (s *scheduleImplementation) SetTaskParameters(parameters map[string]any) ScheduleInterface {
	if parameters == nil {
		s.ParametersField = "{}"
		return s
	}
	bytes, err := json.Marshal(parameters)
	if err != nil {
		s.ParametersField = "{}"
		return s
	}
	s.ParametersField = string(bytes)
	return s
}

// GetStartAt returns the start date and time of the schedule.
func (s *scheduleImplementation) GetStartAt() string {
	return s.StartAtField
}

// SetStartAt sets the start date and time of the schedule.
// If startAt is not set, the schedule will start at the current time.
func (s *scheduleImplementation) SetStartAt(startAt string) ScheduleInterface {
	s.StartAtField = startAt
	return s
}

// GetEndAt returns the end date and time of the schedule.
// The default value is the maximum datetime (never expires).
func (s *scheduleImplementation) GetEndAt() string {
	return s.EndAtField
}

// SetEndAt sets the end date and time of the schedule.
func (s *scheduleImplementation) SetEndAt(endAt string) ScheduleInterface {
	s.EndAtField = endAt
	return s
}

// GetExecutionCount returns the number of times the schedule has been executed.
func (s *scheduleImplementation) GetExecutionCount() int {
	return s.ExecutionCountField
}

// SetExecutionCount sets the number of times the schedule has been executed.
func (s *scheduleImplementation) SetExecutionCount(count int) ScheduleInterface {
	s.ExecutionCountField = count
	return s
}

// GetMaxExecutionCount returns the maximum number of times the schedule is allowed to be executed.
// The default value is int max (no limit). To execute only once, set maxExecutionCount to 1.
func (s *scheduleImplementation) GetMaxExecutionCount() int {
	return s.MaxExecutionCountField
}

// SetMaxExecutionCount sets the maximum number of times the schedule is allowed to be executed.
func (s *scheduleImplementation) SetMaxExecutionCount(count int) ScheduleInterface {
	s.MaxExecutionCountField = count
	return s
}

// GetLastRunAt returns the last date and time the schedule was executed.
func (s *scheduleImplementation) GetLastRunAt() string {
	return s.LastRunAtField
}

// SetLastRunAt sets the last date and time the schedule was executed.
func (s *scheduleImplementation) SetLastRunAt(lastRunAt string) ScheduleInterface {
	s.LastRunAtField = lastRunAt
	return s
}

// GetNextRunAt returns the next date and time the schedule is scheduled to run.
func (s *scheduleImplementation) GetNextRunAt() string {
	return s.NextRunAtField
}

// SetNextRunAt sets the next date and time the schedule is scheduled to run.
func (s *scheduleImplementation) SetNextRunAt(nextRunAt string) ScheduleInterface {
	s.NextRunAtField = nextRunAt
	return s
}

// GetCreatedAt returns the date and time the schedule was created.
func (s *scheduleImplementation) GetCreatedAt() time.Time {
	return s.CreatedAtField.CreatedAt
}

// GetCreatedAtCarbon returns the created at time of the schedule as a carbon object.
func (s *scheduleImplementation) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(s.CreatedAtField.CreatedAt)
}

// SetCreatedAt sets the date and time the schedule was created.
func (s *scheduleImplementation) SetCreatedAt(createdAt time.Time) ScheduleInterface {
	s.CreatedAtField.CreatedAt = createdAt
	return s
}

// GetUpdatedAt returns the date and time the schedule was last updated.
func (s *scheduleImplementation) GetUpdatedAt() time.Time {
	return s.UpdatedAtField.UpdatedAt
}

// GetUpdatedAtCarbon returns the updated at time of the schedule as a carbon object.
func (s *scheduleImplementation) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(s.UpdatedAtField.UpdatedAt)
}

// SetUpdatedAt sets the date and time the schedule was last updated.
func (s *scheduleImplementation) SetUpdatedAt(updatedAt time.Time) ScheduleInterface {
	s.UpdatedAtField.UpdatedAt = updatedAt
	return s
}

// GetSoftDeletedAt returns the date and time the schedule was soft deleted.
// The default value is max datetime (not soft deleted, 9999-12-31 23:59:59).
// To soft delete a schedule, set softDeletedAt to the current time.
// To unsoft delete a schedule, set softDeletedAt to max datetime.
// A soft deleted schedule is when its in the past.
func (s *scheduleImplementation) GetSoftDeletedAt() time.Time {
	return s.SoftDeletesMaxDate.SoftDeletedAt
}

// GetSoftDeletedAtCarbon returns the soft deleted at time of the schedule as a carbon object.
func (s *scheduleImplementation) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(s.SoftDeletesMaxDate.SoftDeletedAt)
}

// SetSoftDeletedAt sets the date and time the schedule was soft deleted.
func (s *scheduleImplementation) SetSoftDeletedAt(softDeletedAt time.Time) ScheduleInterface {
	s.SoftDeletesMaxDate.SoftDeletedAt = softDeletedAt
	return s
}

// HasReachedEndDate returns true if the schedule has reached its end date
func (s *scheduleImplementation) HasReachedEndDate() bool {
	endAt := carbon.Parse(s.EndAtField, carbon.UTC)
	now := carbon.Now(carbon.UTC)

	return now.Gt(endAt)
}

// HasReachedMaxExecutions returns true if the schedule has reached its maximum number of executions
func (s *scheduleImplementation) HasReachedMaxExecutions() bool {
	if s.MaxExecutionCountField <= 0 {
		return false
	}
	return s.ExecutionCountField >= s.MaxExecutionCountField
}

// GetNextOccurrence returns the next occurrence of the schedule
// if invalid recurrence rule, returns error
func (s *scheduleImplementation) GetNextOccurrence() (string, error) {
	rule := s.GetRecurrenceRule()
	if rule == nil {
		return "", nil
	}
	nextRunAt, err := NextRunAt(rule, carbon.Now(carbon.UTC))
	if err != nil {
		return "", err
	}
	return nextRunAt.ToDateTimeString(carbon.UTC), nil
}

// IncrementExecutionCount increments the execution count of the schedule by one
func (s *scheduleImplementation) IncrementExecutionCount() ScheduleInterface {
	s.ExecutionCountField++
	return s
}

// UpdateNextRunAt calculates the next run at of the schedule and updates it
func (s *scheduleImplementation) UpdateNextRunAt() ScheduleInterface {
	nextRunAt, err := s.GetNextOccurrence()
	if err != nil {
		return s
	}
	s.NextRunAtField = nextRunAt
	return s
}

// UpdateLastRunAt updates the last run at of the schedule with current time
func (s *scheduleImplementation) UpdateLastRunAt() ScheduleInterface {
	s.LastRunAtField = carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	return s
}

// IsDue returns true if the schedule is due to run
func (s *scheduleImplementation) IsDue() bool {
	if s.StatusField != "active" {
		return false
	}

	if s.HasReachedEndDate() {
		return false
	}

	if s.HasReachedMaxExecutions() {
		return false
	}

	if s.NextRunAtField == NULL_DATETIME || s.NextRunAtField == "" {
		return false
	}

	nextRunAt := carbon.Parse(s.NextRunAtField, carbon.UTC)
	now := carbon.Now(carbon.UTC)

	if now.Gte(nextRunAt) {
		return true
	}

	return false
}

// UnmarshalJSON implements custom JSON unmarshaling to reset cached fields.
func (s *scheduleImplementation) UnmarshalJSON(data []byte) error {
	type alias scheduleImplementation
	if err := json.Unmarshal(data, (*alias)(s)); err != nil {
		return err
	}
	s.cachedRecurrenceRule = nil
	return nil
}
