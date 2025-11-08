package taskstore

import (
	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

// == CLASS ===================================================================

type task struct {
	dataobject.DataObject
}

var _ TaskInterface = (*task)(nil)

// == CONSTRUCTORS ============================================================

func NewTask() TaskInterface {
	o := &task{}

	o.SetID(uid.HumanUid()).
		SetStatus(TaskStatusActive).
		SetMemo("").
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(sb.MAX_DATETIME)

	// err := o.SetMetas(map[string]string{})

	// if err != nil {
	// 	return o
	// }

	return o
}

func NewTaskFromExistingData(data map[string]string) TaskInterface {
	o := &task{}
	o.Hydrate(data)
	return o
}

// == METHODS =================================================================
func (o *task) IsActive() bool {
	return o.Status() == TaskStatusActive
}

func (o *task) IsCanceled() bool {
	return o.Status() == TaskStatusCanceled
}

func (o *task) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

func (o *task) Alias() string {
	return o.Get(COLUMN_ALIAS)
}

func (o *task) SetAlias(alias string) TaskInterface {
	o.Set(COLUMN_ALIAS, alias)
	return o
}

func (o *task) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *task) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt(), carbon.UTC)
}

func (o *task) SetCreatedAt(createdAt string) TaskInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *task) Description() string {
	return o.Get(COLUMN_DESCRIPTION)
}

func (o *task) SetDescription(description string) TaskInterface {
	o.Set(COLUMN_DESCRIPTION, description)
	return o
}

func (o *task) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *task) SetID(id string) TaskInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (o *task) Memo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *task) SetMemo(memo string) TaskInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

// func (o *task) Metas() (map[string]string, error) {
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

// func (o *task) Meta(name string) string {
// 	metas, err := o.Metas()

// 	if err != nil {
// 		return ""
// 	}

// 	if value, exists := metas[name]; exists {
// 		return value
// 	}

// 	return ""
// }

// func (o *task) SetMeta(name string, value string) error {
// 	return o.UpsertMetas(map[string]string{name: value})
// }

// // SetMetas stores metas as json string
// // Warning: it overwrites any existing metas
// func (o *task) SetMetas(metas map[string]string) error {
// 	mapString, err := utils.ToJSON(metas)
// 	if err != nil {
// 		return err
// 	}
// 	o.Set(COLUMN_METAS, mapString)
// 	return nil
// }

// func (o *task) UpsertMetas(metas map[string]string) error {
// 	currentMetas, err := o.Metas()

// 	if err != nil {
// 		return err
// 	}

// 	for k, v := range metas {
// 		currentMetas[k] = v
// 	}

// 	return o.SetMetas(currentMetas)
// }

func (o *task) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *task) SoftDeletedAt() string {
	return o.Get(COLUMN_DELETED_AT)
}

func (o *task) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt(), carbon.UTC)
}

func (o *task) SetSoftDeletedAt(deletedAt string) TaskInterface {
	o.Set(COLUMN_DELETED_AT, deletedAt)
	return o
}

func (o *task) SetStatus(status string) TaskInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *task) Title() string {
	return o.Get(COLUMN_TITLE)
}

func (o *task) SetTitle(title string) TaskInterface {
	o.Set(COLUMN_TITLE, title)
	return o
}

func (o *task) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *task) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.Get(COLUMN_UPDATED_AT), carbon.UTC)
}

func (o *task) SetUpdatedAt(updatedAt string) TaskInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}
