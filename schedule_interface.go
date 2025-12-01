package taskstore

// ScheduleInterface defines the contract for a schedule, including its
// identity, metadata, recurrence rule, timing fields, execution limits,
// soft-delete semantics, and helper methods for evaluating schedule state.
type ScheduleInterface interface {
	// =======================================================================
	// Accessors (Setters and Getters)
	// =======================================================================

	// GetID the unique identifier of the schedule
	GetID() string

	// SetID sets the unique identifier of the schedule
	SetID(string) ScheduleInterface

	// Name the name of the schedule
	GetName() string

	// SetName sets the name of the schedule
	SetName(string) ScheduleInterface

	// Description the description of the schedule
	GetDescription() string

	// SetDescription sets the description of the schedule
	SetDescription(string) ScheduleInterface

	// Status the status of the schedule
	// Valid values are "draft" (default), "active", "inactive"
	GetStatus() string

	// SetStatus sets the status of the schedule
	SetStatus(string) ScheduleInterface

	// RecurrenceRule the recurrence rule that defines when the schedule should run
	GetRecurrenceRule() RecurrenceRuleInterface

	// SetRecurrenceRule sets the recurrence rule that defines when the schedule should run
	SetRecurrenceRule(RecurrenceRuleInterface) ScheduleInterface

	// QueueName the name of the queue that this schedule is associated with
	GetQueueName() string

	// SetQueueName sets the name of the queue that this schedule is associated with
	SetQueueName(string) ScheduleInterface

	// TaskDefinitionID the unique identifier of the task definition
	// that this schedule is associated with
	GetTaskDefinitionID() string

	// SetTaskDefinitionID sets the unique identifier of the task definition
	// that this schedule is associated with
	SetTaskDefinitionID(string) ScheduleInterface

	// TaskParameters the parameters to be passed to the task definition
	// when it is executed
	GetTaskParameters() map[string]any

	// SetTaskParameters sets the parameters to be passed to the task definition
	// when it is executed
	SetTaskParameters(map[string]any) ScheduleInterface

	// StartAt the start date and time of the schedule
	GetStartAt() string

	// SetStartAt sets the start date and time of the schedule
	// If startAt is not set, the schedule will start at the current time
	SetStartAt(string) ScheduleInterface

	// EndAt the end date and time of the schedule
	// The default value is the maximum datetime (never expires)
	GetEndAt() string

	// SetEndAt sets the end date and time of the schedule
	SetEndAt(string) ScheduleInterface

	// ExecutionCount the number of times the schedule has been executed
	GetExecutionCount() int

	// SetExecutionCount sets the number of times the schedule has been executed
	SetExecutionCount(int) ScheduleInterface

	// MaxExecutionCount the maximum number of times the schedule is allowed to be executed
	// The default value is int max (no limit)
	// To execute only once, set maxExecutionCount to 1
	GetMaxExecutionCount() int

	// SetMaxExecutionCount sets the maximum number of times the schedule is allowed to be executed
	SetMaxExecutionCount(int) ScheduleInterface

	// LastRunAt the last date and time the schedule was executed
	GetLastRunAt() string

	// SetLastRunAt sets the last date and time the schedule was executed
	SetLastRunAt(string) ScheduleInterface

	// NextRunAt the next date and time the schedule is scheduled to run
	GetNextRunAt() string

	// SetNextRunAt sets the next date and time the schedule is scheduled to run
	SetNextRunAt(string) ScheduleInterface

	// CreatedAt the date and time the schedule was created
	GetCreatedAt() string

	// SetCreatedAt sets the date and time the schedule was created
	SetCreatedAt(string) ScheduleInterface

	// UpdatedAt the date and time the schedule was last updated
	GetUpdatedAt() string

	// SetUpdatedAt sets the date and time the schedule was last updated
	SetUpdatedAt(string) ScheduleInterface

	// SoftDeletedAt the date and time the schedule was soft deleted
	// The default value is max datetime (not soft deleted, 9999-12-31 23:59:59)
	// To soft delete a schedule, set softDeletedAt to the current time
	// To unsoft delete a schedule, set softDeletedAt to max datetime
	// A soft deleted schedule is when its in the past
	GetSoftDeletedAt() string

	// SetSoftDeletedAt sets the date and time the schedule was soft deleted
	SetSoftDeletedAt(string) ScheduleInterface

	// =======================================================================
	// Informational Methods
	// =======================================================================

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
