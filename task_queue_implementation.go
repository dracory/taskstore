package taskstore

import (
	"encoding/json"

	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

// == CLASS ===================================================================

type taskQueue struct {
	dataobject.DataObject
}

var _ TaskQueueInterface = (*taskQueue)(nil)

// == CONSTRUCTORS ============================================================

// NewTaskQueue creates a new task queue
// If a queue name is provided, it will be used; otherwise DefaultQueueName is used.
func NewTaskQueue(queueName ...string) TaskQueueInterface {
	name := DefaultQueueName
	if len(queueName) > 0 && queueName[0] != "" {
		name = queueName[0]
	}

	o := &taskQueue{}

	o.SetID(uid.HumanUid()).
		SetStatus(TaskQueueStatusQueued).
		SetQueueName(name).
		SetAttempts(0).
		SetOutput("").
		SetDetails("").
		SetParameters("{}").
		SetStartedAt(sb.NULL_DATETIME).
		SetCompletedAt(sb.NULL_DATETIME).
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(sb.MAX_DATETIME)

	return o
}

func NewTaskQueueFromExistingData(data map[string]string) TaskQueueInterface {
	o := &taskQueue{}
	o.Hydrate(data)
	return o
}

// == METHODS =================================================================

func (o *taskQueue) IsCanceled() bool {
	return o.Status() == TaskQueueStatusCanceled
}

func (o *taskQueue) IsDeleted() bool {
	return o.Status() == TaskQueueStatusDeleted
}

func (o *taskQueue) IsFailed() bool {
	return o.Status() == TaskQueueStatusFailed
}

func (o *taskQueue) IsQueued() bool {
	return o.Status() == TaskQueueStatusQueued
}

func (o *taskQueue) IsPaused() bool {
	return o.Status() == TaskQueueStatusPaused
}

func (o *taskQueue) IsRunning() bool {
	return o.Status() == TaskQueueStatusRunning
}

func (o *taskQueue) IsSuccess() bool {
	return o.Status() == TaskQueueStatusSuccess
}

func (o *taskQueue) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

func (o *taskQueue) GetAttempts() int {
	attempts := o.Get(COLUMN_ATTEMPTS)
	return cast.ToInt(attempts)
}

// Attempts alias is kept for backwards compatibility.
// Deprecated: use GetAttempts instead. Will be removed after 2026-11-30.
func (o *taskQueue) Attempts() int {
	return o.GetAttempts()
}

func (o *taskQueue) SetAttempts(attempts int) TaskQueueInterface {
	o.Set(COLUMN_ATTEMPTS, cast.ToString(attempts))
	return o
}

func (o *taskQueue) GetCompletedAt() string {
	return o.Get(COLUMN_COMPLETED_AT)
}

// CompletedAt alias is kept for backwards compatibility.
// Deprecated: use GetCompletedAt instead. Will be removed after 2026-11-30.
func (o *taskQueue) CompletedAt() string {
	return o.GetCompletedAt()
}

func (o *taskQueue) CompletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetCompletedAt(), carbon.UTC)
}

func (o *taskQueue) SetCompletedAt(completedAt string) TaskQueueInterface {
	o.Set(COLUMN_COMPLETED_AT, completedAt)
	return o
}

func (o *taskQueue) GetCreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

// CreatedAt alias is kept for backwards compatibility.
// Deprecated: use GetCreatedAt instead. Will be removed after 2026-11-30.
func (o *taskQueue) CreatedAt() string {
	return o.GetCreatedAt()
}

func (o *taskQueue) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetCreatedAt(), carbon.UTC)
}

func (o *taskQueue) SetCreatedAt(createdAt string) TaskQueueInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *taskQueue) GetID() string {
	return o.Get(COLUMN_ID)
}

// ID alias is kept for backwards compatibility.
// Deprecated: use GetID instead. Will be removed after 2026-11-30.
func (o *taskQueue) ID() string {
	return o.GetID()
}

// AppendDetails appends details to the queued task
// !!! warning does not auto-save it for performance reasons
func (o *taskQueue) AppendDetails(details string) TaskQueueInterface {
	ts := carbon.Now().Format("Y-m-d H:i:s")
	text := o.Details()
	if text != "" {
		text += "\n"
	}
	text += ts + " : " + details
	return o.SetDetails(text)
}

func (o *taskQueue) GetDetails() string {
	return o.Get(COLUMN_DETAILS)
}

// Details alias is kept for backwards compatibility.
// Deprecated: use GetDetails instead. Will be removed after 2026-11-30.
func (o *taskQueue) Details() string {
	return o.GetDetails()
}

func (o *taskQueue) SetDetails(details string) TaskQueueInterface {
	o.Set(COLUMN_DETAILS, details)
	return o
}

func (o *taskQueue) GetQueueName() string {
	return o.Get(COLUMN_QUEUE_NAME)
}

// QueueName alias is kept for backwards compatibility.
// Deprecated: use GetQueueName instead. Will be removed after 2026-11-30.
func (o *taskQueue) QueueName() string {
	return o.GetQueueName()
}

func (o *taskQueue) SetQueueName(queueName string) TaskQueueInterface {
	o.Set(COLUMN_QUEUE_NAME, queueName)
	return o
}

func (o *taskQueue) SetID(id string) TaskQueueInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// func (o *taskQueue) Memo() string {
// 	return o.Get(COLUMN_MEMO)
// }

// func (o *taskQueue) SetMemo(memo string) TaskQueueInterface {
// 	o.Set(COLUMN_MEMO, memo)
// 	return o
// }

// func (o *taskQueue) Metas() (map[string]string, error) {
// 	metasStr := o.Get(COLUMN_METAS)

// 	if metasStr == "" {
// 		metasStr = "{}"
// 	}

// 	metasJson, errJson := utils.FromJSON(metasStr, map[string]string{})
// 	if errJson != nil {
// 		return map[string]string{}, errJson
// 	}

// 	return maputils.MapStringAnyToMapStringString(metasJson.(map[string]any)), nil
// }

// func (o *taskQueue) Meta(name string) string {
// 	metas, err := o.Metas()

// 	if err != nil {
// 		return ""
// 	}

// 	if value, exists := metas[name]; exists {
// 		return value
// 	}

// 	return ""
// }

// func (o *taskQueue) SetMeta(name string, value string) error {
// 	return o.UpsertMetas(map[string]string{name: value})
// }

// // SetMetas stores metas as json string
// // Warning: it overwrites any existing metas
// func (o *taskQueue) SetMetas(metas map[string]string) error {
// 	mapString, err := utils.ToJSON(metas)
// 	if err != nil {
// 		return err
// 	}
// 	o.Set(COLUMN_METAS, mapString)
// 	return nil
// }

// func (o *taskQueue) UpsertMetas(metas map[string]string) error {
// 	currentMetas, err := o.Metas()

// 	if err != nil {
// 		return err
// 	}

// 	for k, v := range metas {
// 		currentMetas[k] = v
// 	}

// 	return o.SetMetas(currentMetas)
// }

func (o *taskQueue) GetOutput() string {
	return o.Get(COLUMN_OUTPUT)
}

// Output alias is kept for backwards compatibility.
// Deprecated: use GetOutput instead. Will be removed after 2026-11-30.
func (o *taskQueue) Output() string {
	return o.GetOutput()
}

func (o *taskQueue) SetOutput(output string) TaskQueueInterface {
	o.Set(COLUMN_OUTPUT, output)
	return o
}

func (o *taskQueue) GetParameters() string {
	return o.Get(COLUMN_PARAMETERS)
}

// Parameters alias is kept for backwards compatibility.
// Deprecated: use GetParameters instead. Will be removed after 2026-11-30.
func (o *taskQueue) Parameters() string {
	return o.GetParameters()
}

func (o *taskQueue) SetParameters(parameters string) TaskQueueInterface {
	o.Set(COLUMN_PARAMETERS, parameters)
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
	parametersJsonBytes, jsonErr := json.Marshal(parameters)
	if jsonErr != nil {
		return o, jsonErr
	}
	parametersJson := string(parametersJsonBytes)
	return o.SetParameters(parametersJson), nil
}

func (o *taskQueue) GetStartedAt() string {
	return o.Get(COLUMN_STARTED_AT)
}

// StartedAt alias is kept for backwards compatibility.
// Deprecated: use GetStartedAt instead. Will be removed after 2026-11-30.
func (o *taskQueue) StartedAt() string {
	return o.GetStartedAt()
}

func (o *taskQueue) StartedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetStartedAt(), carbon.UTC)
}

func (o *taskQueue) SetStartedAt(startedAt string) TaskQueueInterface {
	o.Set(COLUMN_STARTED_AT, startedAt)
	return o
}

func (o *taskQueue) GetStatus() string {
	return o.Get(COLUMN_STATUS)
}

// Status alias is kept for backwards compatibility.
// Deprecated: use GetStatus instead. Will be removed after 2026-11-30.
func (o *taskQueue) Status() string {
	return o.GetStatus()
}

func (o *taskQueue) GetSoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

// SoftDeletedAt alias is kept for backwards compatibility.
// Deprecated: use GetSoftDeletedAt instead. Will be removed after 2026-11-30.
func (o *taskQueue) SoftDeletedAt() string {
	return o.GetSoftDeletedAt()
}

func (o *taskQueue) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetSoftDeletedAt(), carbon.UTC)
}

func (o *taskQueue) SetSoftDeletedAt(deletedAt string) TaskQueueInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, deletedAt)
	return o
}

func (o *taskQueue) SetStatus(status string) TaskQueueInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *taskQueue) GetTaskID() string {
	return o.Get(COLUMN_TASK_ID)
}

// TaskID alias is kept for backwards compatibility.
// Deprecated: use GetTaskID instead. Will be removed after 2026-11-30.
func (o *taskQueue) TaskID() string {
	return o.GetTaskID()
}

func (o *taskQueue) SetTaskID(taskID string) TaskQueueInterface {
	o.Set(COLUMN_TASK_ID, taskID)
	return o
}

func (o *taskQueue) GetUpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

// UpdatedAt alias is kept for backwards compatibility.
// Deprecated: use GetUpdatedAt instead. Will be removed after 2026-11-30.
func (o *taskQueue) UpdatedAt() string {
	return o.GetUpdatedAt()
}

func (o *taskQueue) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.GetUpdatedAt(), carbon.UTC)
}

func (o *taskQueue) SetUpdatedAt(updatedAt string) TaskQueueInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}
