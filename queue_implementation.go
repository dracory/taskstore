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

type queue struct {
	dataobject.DataObject
}

var _ QueueInterface = (*queue)(nil)

// == CONSTRUCTORS ============================================================

func NewQueue() QueueInterface {
	o := &queue{}

	o.SetID(uid.HumanUid()).
		SetStatus(QueueStatusQueued).
		// SetMemo("").
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(sb.MAX_DATETIME)

	// err := o.SetMetas(map[string]string{})

	// if err != nil {
	// 	return o
	// }

	return o
}

func NewQueueFromExistingData(data map[string]string) QueueInterface {
	o := &queue{}
	o.Hydrate(data)
	return o
}

// == METHODS =================================================================

func (o *queue) IsCanceled() bool {
	return o.Status() == QueueStatusCanceled
}

func (o *queue) IsDeleted() bool {
	return o.Status() == QueueStatusDeleted
}

func (o *queue) IsFailed() bool {
	return o.Status() == QueueStatusFailed
}

func (o *queue) IsQueued() bool {
	return o.Status() == QueueStatusQueued
}

func (o *queue) IsPaused() bool {
	return o.Status() == QueueStatusPaused
}

func (o *queue) IsRunning() bool {
	return o.Status() == QueueStatusRunning
}

func (o *queue) IsSuccess() bool {
	return o.Status() == QueueStatusSuccess
}

func (o *queue) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

func (o *queue) Attempts() int {
	attempts := o.Get(COLUMN_ATTEMPTS)
	return cast.ToInt(attempts)
}

func (o *queue) SetAttempts(attempts int) QueueInterface {
	o.Set(COLUMN_ATTEMPTS, cast.ToString(attempts))
	return o
}

func (o *queue) CompletedAt() string {
	return o.Get(COLUMN_COMPLETED_AT)
}

func (o *queue) CompletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CompletedAt(), carbon.UTC)
}

func (o *queue) SetCompletedAt(completedAt string) QueueInterface {
	o.Set(COLUMN_COMPLETED_AT, completedAt)
	return o
}

func (o *queue) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *queue) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt(), carbon.UTC)
}

func (o *queue) SetCreatedAt(createdAt string) QueueInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *queue) ID() string {
	return o.Get(COLUMN_ID)
}

// AppendDetails appends details to the queued task
// !!! warning does not auto-save it for performance reasons
func (o *queue) AppendDetails(details string) QueueInterface {
	ts := carbon.Now().Format("Y-m-d H:i:s")
	text := o.Details()
	if text != "" {
		text += "\n"
	}
	text += ts + " : " + details
	return o.SetDetails(text)
}

func (o *queue) Details() string {
	return o.Get(COLUMN_DETAILS)
}

func (o *queue) SetDetails(details string) QueueInterface {
	o.Set(COLUMN_DETAILS, details)
	return o
}

func (o *queue) SetID(id string) QueueInterface {
	o.Set(COLUMN_ID, id)
	return o
}

// func (o *queue) Memo() string {
// 	return o.Get(COLUMN_MEMO)
// }

// func (o *queue) SetMemo(memo string) QueueInterface {
// 	o.Set(COLUMN_MEMO, memo)
// 	return o
// }

// func (o *queue) Metas() (map[string]string, error) {
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

// func (o *queue) Meta(name string) string {
// 	metas, err := o.Metas()

// 	if err != nil {
// 		return ""
// 	}

// 	if value, exists := metas[name]; exists {
// 		return value
// 	}

// 	return ""
// }

// func (o *queue) SetMeta(name string, value string) error {
// 	return o.UpsertMetas(map[string]string{name: value})
// }

// // SetMetas stores metas as json string
// // Warning: it overwrites any existing metas
// func (o *queue) SetMetas(metas map[string]string) error {
// 	mapString, err := utils.ToJSON(metas)
// 	if err != nil {
// 		return err
// 	}
// 	o.Set(COLUMN_METAS, mapString)
// 	return nil
// }

// func (o *queue) UpsertMetas(metas map[string]string) error {
// 	currentMetas, err := o.Metas()

// 	if err != nil {
// 		return err
// 	}

// 	for k, v := range metas {
// 		currentMetas[k] = v
// 	}

// 	return o.SetMetas(currentMetas)
// }

func (o *queue) Output() string {
	return o.Get(COLUMN_OUTPUT)
}

func (o *queue) SetOutput(output string) QueueInterface {
	o.Set(COLUMN_OUTPUT, output)
	return o
}

func (o *queue) Parameters() string {
	return o.Get(COLUMN_PARAMETERS)
}

func (o *queue) SetParameters(parameters string) QueueInterface {
	o.Set(COLUMN_PARAMETERS, parameters)
	return o
}

func (o *queue) ParametersMap() (map[string]string, error) {
	var parameters map[string]string
	jsonErr := json.Unmarshal([]byte(o.Parameters()), &parameters)
	if jsonErr != nil {
		return map[string]string{}, jsonErr
	}
	return parameters, nil
}

func (o *queue) SetParametersMap(parameters map[string]string) (QueueInterface, error) {
	parametersJsonBytes, jsonErr := json.Marshal(parameters)
	if jsonErr != nil {
		return o, jsonErr
	}
	parametersJson := string(parametersJsonBytes)
	return o.SetParameters(parametersJson), nil
}

func (o *queue) StartedAt() string {
	return o.Get(COLUMN_STARTED_AT)
}

func (o *queue) StartedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.StartedAt(), carbon.UTC)
}

func (o *queue) SetStartedAt(startedAt string) QueueInterface {
	o.Set(COLUMN_STARTED_AT, startedAt)
	return o
}

func (o *queue) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *queue) SoftDeletedAt() string {
	return o.Get(COLUMN_DELETED_AT)
}

func (o *queue) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt(), carbon.UTC)
}

func (o *queue) SetSoftDeletedAt(deletedAt string) QueueInterface {
	o.Set(COLUMN_DELETED_AT, deletedAt)
	return o
}

func (o *queue) SetStatus(status string) QueueInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *queue) TaskID() string {
	return o.Get(COLUMN_TASK_ID)
}

func (o *queue) SetTaskID(taskID string) QueueInterface {
	o.Set(COLUMN_TASK_ID, taskID)
	return o
}

func (o *queue) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *queue) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.Get(COLUMN_UPDATED_AT), carbon.UTC)
}

func (o *queue) SetUpdatedAt(updatedAt string) QueueInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}
