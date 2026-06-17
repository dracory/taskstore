package taskstore

import (
	"encoding/json"
	"time"

	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

// == TYPE =====================================================================

type taskQueue struct {
	orm.ShortID

	QueueNameField   string    `db:"queue_name"`
	TaskIDField      string    `db:"task_id"`
	ParametersField  string    `db:"parameters"`
	StatusField      string    `db:"status"`
	OutputField      string    `db:"output"`
	DetailsField     string    `db:"details"`
	AttemptsField    int       `db:"attempts"`
	StartedAtField   time.Time `db:"started_at"`
	CompletedAtField time.Time `db:"completed_at"`

	CreatedAtField orm.CreatedAt
	UpdatedAtField orm.UpdatedAt
	soft_delete.SoftDeletesMaxDate
}

var _ TaskQueueInterface = (*taskQueue)(nil)

// == CONSTRUCTORS =============================================================

// NewTaskQueue creates a new task queue
// If a queue name is provided, it will be used; otherwise DefaultQueueName is used.
func NewTaskQueue(queueName ...string) TaskQueueInterface {
	name := DefaultQueueName
	if len(queueName) > 0 && queueName[0] != "" {
		name = queueName[0]
	}

	o := &taskQueue{}

	o.SetID(neatuid.GenerateShortID()).
		SetStatus(TaskQueueStatusQueued).
		SetQueueName(name).
		SetAttempts(0).
		SetOutput("").
		SetDetails("").
		SetParameters("{}").
		SetStartedAt(NULL_DATETIME).
		SetCompletedAt(NULL_DATETIME).
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(MAX_DATETIME)

	return o
}

func NewTaskQueueFromExistingData(data map[string]string) TaskQueueInterface {
	o := &taskQueue{}
	o.SetID(data[COLUMN_ID])
	o.SetQueueName(data[COLUMN_QUEUE_NAME])
	o.SetTaskID(data[COLUMN_TASK_ID])
	o.SetParameters(data[COLUMN_PARAMETERS])
	o.SetStatus(data[COLUMN_STATUS])
	o.SetOutput(data[COLUMN_OUTPUT])
	o.SetDetails(data[COLUMN_DETAILS])
	o.SetAttempts(cast.ToInt(data[COLUMN_ATTEMPTS]))
	if v, ok := data[COLUMN_STARTED_AT]; ok {
		o.SetStartedAt(v)
	}
	if v, ok := data[COLUMN_COMPLETED_AT]; ok {
		o.SetCompletedAt(v)
	}
	if v, ok := data[COLUMN_CREATED_AT]; ok {
		o.SetCreatedAt(v)
	}
	if v, ok := data[COLUMN_UPDATED_AT]; ok {
		o.SetUpdatedAt(v)
	}
	if v, ok := data[COLUMN_SOFT_DELETED_AT]; ok {
		o.SetSoftDeletedAt(v)
	}
	return o
}

// == METHODS ==================================================================

func (o *taskQueue) IsCanceled() bool {
	return o.GetStatus() == TaskQueueStatusCanceled
}

func (o *taskQueue) IsDeleted() bool {
	return o.GetStatus() == TaskQueueStatusDeleted
}

func (o *taskQueue) IsFailed() bool {
	return o.GetStatus() == TaskQueueStatusFailed
}

func (o *taskQueue) IsQueued() bool {
	return o.GetStatus() == TaskQueueStatusQueued
}

func (o *taskQueue) IsPaused() bool {
	return o.GetStatus() == TaskQueueStatusPaused
}

func (o *taskQueue) IsRunning() bool {
	return o.GetStatus() == TaskQueueStatusRunning
}

func (o *taskQueue) IsSuccess() bool {
	return o.GetStatus() == TaskQueueStatusSuccess
}

func (o *taskQueue) IsSoftDeleted() bool {
	return o.SoftDeletesMaxDate.IsSoftDeleted()
}

// == SETTERS AND GETTERS ======================================================

func (o *taskQueue) GetAttempts() int {
	return o.AttemptsField
}

func (o *taskQueue) SetAttempts(attempts int) TaskQueueInterface {
	o.AttemptsField = attempts
	return o
}

func (o *taskQueue) GetCompletedAt() string {
	if o.CompletedAtField.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CompletedAtField).ToDateTimeString()
}

func (o *taskQueue) CompletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CompletedAtField)
}

func (o *taskQueue) SetCompletedAt(completedAt string) TaskQueueInterface {
	if completedAt == "" || completedAt == NULL_DATETIME {
		o.CompletedAtField = time.Time{}
		return o
	}
	o.CompletedAtField = carbon.Parse(completedAt, carbon.UTC).StdTime()
	return o
}

func (o *taskQueue) GetCreatedAt() string {
	if o.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString()
}

func (o *taskQueue) CreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt)
}

func (o *taskQueue) SetCreatedAt(createdAt string) TaskQueueInterface {
	if createdAt == "" {
		return o
	}
	o.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

// AppendDetails appends details to the queued task
// !!! warning does not auto-save it for performance reasons
func (o *taskQueue) AppendDetails(details string) TaskQueueInterface {
	ts := carbon.Now().Format("Y-m-d H:i:s")
	text := o.GetDetails()
	if text != "" {
		text += "\n"
	}
	text += ts + " : " + details
	return o.SetDetails(text)
}

func (o *taskQueue) GetDetails() string {
	return o.DetailsField
}

func (o *taskQueue) SetDetails(details string) TaskQueueInterface {
	o.DetailsField = details
	return o
}

func (o *taskQueue) GetQueueName() string {
	return o.QueueNameField
}

func (o *taskQueue) SetQueueName(queueName string) TaskQueueInterface {
	o.QueueNameField = queueName
	return o
}

func (o *taskQueue) GetID() string {
	return o.ShortID.ID
}

func (o *taskQueue) SetID(id string) TaskQueueInterface {
	o.ShortID.ID = id
	return o
}

func (o *taskQueue) GetOutput() string {
	return o.OutputField
}

func (o *taskQueue) SetOutput(output string) TaskQueueInterface {
	o.OutputField = output
	return o
}

func (o *taskQueue) GetParameters() string {
	return o.ParametersField
}

func (o *taskQueue) SetParameters(parameters string) TaskQueueInterface {
	o.ParametersField = parameters
	return o
}

func (o *taskQueue) ParametersMap() (map[string]string, error) {
	// Handle empty string parameters
	if o.GetParameters() == "" {
		return map[string]string{}, nil
	}

	var parameters map[string]string
	jsonErr := json.Unmarshal([]byte(o.GetParameters()), &parameters)
	if jsonErr != nil {
		return map[string]string{}, jsonErr
	}
	return parameters, nil
}

func (o *taskQueue) SetParametersMap(parameters map[string]string) (TaskQueueInterface, error) {
	parametersBytes, jsonErr := json.Marshal(parameters)
	if jsonErr != nil {
		return o, jsonErr
	}
	return o.SetParameters(string(parametersBytes)), nil
}

func (o *taskQueue) GetSoftDeletedAt() string {
	if o.SoftDeletesMaxDate.SoftDeletedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt).ToDateTimeString()
}

func (o *taskQueue) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt)
}

func (o *taskQueue) SetSoftDeletedAt(deletedAt string) TaskQueueInterface {
	if deletedAt == "" {
		return o
	}
	o.SoftDeletesMaxDate.SoftDeletedAt = carbon.Parse(deletedAt, carbon.UTC).StdTime()
	return o
}

func (o *taskQueue) GetStartedAt() string {
	if o.StartedAtField.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.StartedAtField).ToDateTimeString()
}

func (o *taskQueue) StartedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.StartedAtField)
}

func (o *taskQueue) SetStartedAt(startedAt string) TaskQueueInterface {
	if startedAt == "" || startedAt == NULL_DATETIME {
		o.StartedAtField = time.Time{}
		return o
	}
	o.StartedAtField = carbon.Parse(startedAt, carbon.UTC).StdTime()
	return o
}

func (o *taskQueue) GetStatus() string {
	return o.StatusField
}

func (o *taskQueue) SetStatus(status string) TaskQueueInterface {
	o.StatusField = status
	return o
}

func (o *taskQueue) GetTaskID() string {
	return o.TaskIDField
}

func (o *taskQueue) SetTaskID(taskID string) TaskQueueInterface {
	o.TaskIDField = taskID
	return o
}

func (o *taskQueue) GetUpdatedAt() string {
	if o.UpdatedAtField.UpdatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString()
}

func (o *taskQueue) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt)
}

func (o *taskQueue) SetUpdatedAt(updatedAt string) TaskQueueInterface {
	if updatedAt == "" {
		return o
	}
	o.UpdatedAtField.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return o
}
