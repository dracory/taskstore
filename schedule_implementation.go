package taskstore

import (
	"encoding/json"

	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	"github.com/dromara/carbon/v2"
)

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
}

var _ ScheduleInterface = (*scheduleImplementation)(nil)

// NewSchedule creates a new schedule with default values and a new recurrence rule.
func NewSchedule() ScheduleInterface {
	return &scheduleImplementation{
		StatusField:            "draft",
		RecurrenceRuleField:    "",
		StartAtField:           NULL_DATETIME,
		EndAtField:             MAX_DATETIME,
		LastRunAtField:         NULL_DATETIME,
		NextRunAtField:         NULL_DATETIME,
		ExecutionCountField:    0,
		MaxExecutionCountField: 0,
	}
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
func (s *scheduleImplementation) GetRecurrenceRule() RecurrenceRuleInterface {
	if s.RecurrenceRuleField == "" {
		return NewRecurrenceRule()
	}
	rr := NewRecurrenceRule()
	_ = json.Unmarshal([]byte(s.RecurrenceRuleField), rr)
	return rr
}

// SetRecurrenceRule sets the recurrence rule that defines when the schedule should run.
func (s *scheduleImplementation) SetRecurrenceRule(rule RecurrenceRuleInterface) ScheduleInterface {
	if rule == nil {
		s.RecurrenceRuleField = ""
		return s
	}
	bytes, err := json.Marshal(rule)
	if err != nil {
		s.RecurrenceRuleField = ""
		return s
	}
	s.RecurrenceRuleField = string(bytes)
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
func (s *scheduleImplementation) GetCreatedAt() string {
	if s.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(s.CreatedAtField.CreatedAt).ToDateTimeString()
}

// CreatedAtCarbon returns the created at time of the schedule as a carbon object.
func (s *scheduleImplementation) CreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(s.CreatedAtField.CreatedAt)
}

// SetCreatedAt sets the date and time the schedule was created.
func (s *scheduleImplementation) SetCreatedAt(createdAt string) ScheduleInterface {
	if createdAt == "" {
		return s
	}
	s.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return s
}

// GetUpdatedAt returns the date and time the schedule was last updated.
func (s *scheduleImplementation) GetUpdatedAt() string {
	if s.UpdatedAtField.UpdatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(s.UpdatedAtField.UpdatedAt).ToDateTimeString()
}

// UpdatedAtCarbon returns the updated at time of the schedule as a carbon object.
func (s *scheduleImplementation) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(s.UpdatedAtField.UpdatedAt)
}

// SetUpdatedAt sets the date and time the schedule was last updated.
func (s *scheduleImplementation) SetUpdatedAt(updatedAt string) ScheduleInterface {
	if updatedAt == "" {
		return s
	}
	s.UpdatedAtField.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return s
}

// GetSoftDeletedAt returns the date and time the schedule was soft deleted.
// The default value is max datetime (not soft deleted, 9999-12-31 23:59:59).
// To soft delete a schedule, set softDeletedAt to the current time.
// To unsoft delete a schedule, set softDeletedAt to max datetime.
// A soft deleted schedule is when its in the past.
func (s *scheduleImplementation) GetSoftDeletedAt() string {
	if s.SoftDeletesMaxDate.SoftDeletedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(s.SoftDeletesMaxDate.SoftDeletedAt).ToDateTimeString()
}

// SoftDeletedAtCarbon returns the soft deleted at time of the schedule as a carbon object.
func (s *scheduleImplementation) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(s.SoftDeletesMaxDate.SoftDeletedAt)
}

// SetSoftDeletedAt sets the date and time the schedule was soft deleted.
func (s *scheduleImplementation) SetSoftDeletedAt(softDeletedAt string) ScheduleInterface {
	if softDeletedAt == "" {
		return s
	}
	s.SoftDeletesMaxDate.SoftDeletedAt = carbon.Parse(softDeletedAt, carbon.UTC).StdTime()
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
