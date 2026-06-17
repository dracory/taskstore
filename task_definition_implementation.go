package taskstore

import (
	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

// == TYPE ====================================================================

type taskDefinition struct {
	orm.ShortID

	StatusField         string `db:"status"`
	AliasField          string `db:"alias"`
	TitleField          string `db:"title"`
	DescriptionField    string `db:"description"`
	MemoField           string `db:"memo"`
	IsRecurringField    int    `db:"is_recurring"`
	RecurrenceRuleField string `db:"recurrence_rule"`

	CreatedAtField orm.CreatedAt
	UpdatedAtField orm.UpdatedAt
	soft_delete.SoftDeletesMaxDate
}

var _ TaskDefinitionInterface = (*taskDefinition)(nil)

// == CONSTRUCTORS ============================================================

func NewTaskDefinition() TaskDefinitionInterface {
	o := &taskDefinition{}
	o.SetID(neatuid.GenerateShortID()).
		SetStatus(TaskDefinitionStatusActive).
		SetAlias("").
		SetTitle("").
		SetDescription("").
		SetIsRecurring(0).
		SetRecurrenceRule("").
		SetMemo("").
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(MAX_DATETIME)

	return o
}

func NewTaskDefinitionFromExistingData(data map[string]string) TaskDefinitionInterface {
	o := &taskDefinition{}
	o.SetID(data[COLUMN_ID])
	o.SetStatus(data[COLUMN_STATUS])
	o.SetAlias(data[COLUMN_ALIAS])
	o.SetTitle(data[COLUMN_TITLE])
	o.SetDescription(data[COLUMN_DESCRIPTION])
	o.SetMemo(data[COLUMN_MEMO])
	o.SetIsRecurring(cast.ToInt(data[COLUMN_IS_RECURRING]))
	o.SetRecurrenceRule(data[COLUMN_RECURRENCE_RULE])
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

// == METHODS =================================================================

func (o *taskDefinition) IsActive() bool {
	return o.GetStatus() == TaskDefinitionStatusActive
}

func (o *taskDefinition) IsCanceled() bool {
	return o.GetStatus() == TaskDefinitionStatusCanceled
}

func (o *taskDefinition) IsSoftDeleted() bool {
	return o.SoftDeletesMaxDate.IsSoftDeleted()
}

// == SETTERS AND GETTERS =====================================================

func (o *taskDefinition) GetAlias() string {
	return o.AliasField
}

func (o *taskDefinition) SetAlias(alias string) TaskDefinitionInterface {
	o.AliasField = alias
	return o
}

func (o *taskDefinition) GetCreatedAt() string {
	if o.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString()
}

func (o *taskDefinition) CreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt)
}

func (o *taskDefinition) SetCreatedAt(createdAt string) TaskDefinitionInterface {
	if createdAt == "" {
		return o
	}
	o.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

func (o *taskDefinition) GetDescription() string {
	return o.DescriptionField
}

func (o *taskDefinition) SetDescription(description string) TaskDefinitionInterface {
	o.DescriptionField = description
	return o
}

func (o *taskDefinition) GetID() string {
	return o.ShortID.ID
}

func (o *taskDefinition) SetID(id string) TaskDefinitionInterface {
	o.ShortID.ID = id
	return o
}

func (o *taskDefinition) GetMemo() string {
	return o.MemoField
}

func (o *taskDefinition) SetMemo(memo string) TaskDefinitionInterface {
	o.MemoField = memo
	return o
}

func (o *taskDefinition) GetIsRecurring() int {
	return o.IsRecurringField
}

func (o *taskDefinition) SetIsRecurring(isRecurring int) TaskDefinitionInterface {
	o.IsRecurringField = isRecurring
	return o
}

func (o *taskDefinition) GetRecurrenceRule() string {
	return o.RecurrenceRuleField
}

func (o *taskDefinition) SetRecurrenceRule(recurrenceRule string) TaskDefinitionInterface {
	o.RecurrenceRuleField = recurrenceRule
	return o
}

func (o *taskDefinition) GetSoftDeletedAt() string {
	if o.SoftDeletesMaxDate.SoftDeletedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt).ToDateTimeString()
}

func (o *taskDefinition) SoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt)
}

func (o *taskDefinition) SetSoftDeletedAt(deletedAt string) TaskDefinitionInterface {
	if deletedAt == "" {
		return o
	}
	o.SoftDeletesMaxDate.SoftDeletedAt = carbon.Parse(deletedAt, carbon.UTC).StdTime()
	return o
}

func (o *taskDefinition) GetStatus() string {
	return o.StatusField
}

func (o *taskDefinition) SetStatus(status string) TaskDefinitionInterface {
	o.StatusField = status
	return o
}

func (o *taskDefinition) GetTitle() string {
	return o.TitleField
}

func (o *taskDefinition) SetTitle(title string) TaskDefinitionInterface {
	o.TitleField = title
	return o
}

func (o *taskDefinition) GetUpdatedAt() string {
	if o.UpdatedAtField.UpdatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt).ToDateTimeString()
}

func (o *taskDefinition) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt)
}

func (o *taskDefinition) SetUpdatedAt(updatedAt string) TaskDefinitionInterface {
	if updatedAt == "" {
		return o
	}
	o.UpdatedAtField.UpdatedAt = carbon.Parse(updatedAt, carbon.UTC).StdTime()
	return o
}
