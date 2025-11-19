package taskstore

import (
	"github.com/dracory/dataobject"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

// == CLASS ===================================================================

type taskDefinition struct {
	dataobject.DataObject
}

var _ TaskDefinitionInterface = (*taskDefinition)(nil)

// == CONSTRUCTORS ============================================================

func NewTaskDefinition() TaskDefinitionInterface {
	o := &taskDefinition{}

	o.SetID(uid.HumanUid()).
		SetStatus(TaskDefinitionStatusActive).
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

func NewTaskDefinitionFromExistingData(data map[string]string) TaskDefinitionInterface {
	o := &taskDefinition{}
	o.Hydrate(data)
	return o
}

// == METHODS =================================================================
func (o *taskDefinition) IsActive() bool {
	return o.Status() == TaskDefinitionStatusActive
}

func (o *taskDefinition) IsCanceled() bool {
	return o.Status() == TaskDefinitionStatusCanceled
}

func (o *taskDefinition) IsSoftDeleted() bool {
	return o.SoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

func (o *taskDefinition) Alias() string {
	return o.Get(COLUMN_ALIAS)
}

func (o *taskDefinition) SetAlias(alias string) TaskDefinitionInterface {
	o.Set(COLUMN_ALIAS, alias)
	return o
}

func (o *taskDefinition) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *taskDefinition) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.CreatedAt(), carbon.UTC)
}

func (o *taskDefinition) SetCreatedAt(createdAt string) TaskDefinitionInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *taskDefinition) Description() string {
	return o.Get(COLUMN_DESCRIPTION)
}

func (o *taskDefinition) SetDescription(description string) TaskDefinitionInterface {
	o.Set(COLUMN_DESCRIPTION, description)
	return o
}

func (o *taskDefinition) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *taskDefinition) SetID(id string) TaskDefinitionInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (o *taskDefinition) Memo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *taskDefinition) SetMemo(memo string) TaskDefinitionInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

// func (o *taskDefinition) Metas() (map[string]string, error) {
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

// func (o *taskDefinition) Meta(name string) string {
// 	metas, err := o.Metas()

// 	if err != nil {
// 		return ""
// 	}

// 	if value, exists := metas[name]; exists {
// 		return value
// 	}

// 	return ""
// }

// func (o *taskDefinition) SetMeta(name string, value string) error {
// 	return o.UpsertMetas(map[string]string{name: value})
// }

// // SetMetas stores metas as json string
// // Warning: it overwrites any existing metas
// func (o *taskDefinition) SetMetas(metas map[string]string) error {
// 	mapString, err := utils.ToJSON(metas)
// 	if err != nil {
// 		return err
// 	}
// 	o.Set(COLUMN_METAS, mapString)
// 	return nil
// }

// func (o *taskDefinition) UpsertMetas(metas map[string]string) error {
// 	currentMetas, err := o.Metas()

// 	if err != nil {
// 		return err
// 	}

// 	for k, v := range metas {
// 		currentMetas[k] = v
// 	}

// 	return o.SetMetas(currentMetas)
// }

func (o *taskDefinition) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *taskDefinition) SoftDeletedAt() string {
	return o.Get(COLUMN_DELETED_AT)
}

func (o *taskDefinition) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.SoftDeletedAt(), carbon.UTC)
}

func (o *taskDefinition) SetSoftDeletedAt(deletedAt string) TaskDefinitionInterface {
	o.Set(COLUMN_DELETED_AT, deletedAt)
	return o
}

func (o *taskDefinition) SetStatus(status string) TaskDefinitionInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *taskDefinition) Title() string {
	return o.Get(COLUMN_TITLE)
}

func (o *taskDefinition) SetTitle(title string) TaskDefinitionInterface {
	o.Set(COLUMN_TITLE, title)
	return o
}

func (o *taskDefinition) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *taskDefinition) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse(o.Get(COLUMN_UPDATED_AT), carbon.UTC)
}

func (o *taskDefinition) SetUpdatedAt(updatedAt string) TaskDefinitionInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}
