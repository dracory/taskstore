package taskstore

import (
	"time"

	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

// == INTERFACE =================================================================

type TaskDefinitionInterface interface {
	IsActive() bool
	IsCanceled() bool
	IsSoftDeleted() bool

	GetAlias() string
	SetAlias(alias string) TaskDefinitionInterface

	GetCreatedAt() time.Time
	GetCreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt time.Time) TaskDefinitionInterface

	GetDescription() string
	SetDescription(description string) TaskDefinitionInterface

	GetID() string
	SetID(id string) TaskDefinitionInterface

	GetMemo() string
	SetMemo(memo string) TaskDefinitionInterface

	GetIsRecurring() int
	SetIsRecurring(isRecurring int) TaskDefinitionInterface

	GetRecurrenceRule() string
	SetRecurrenceRule(recurrenceRule string) TaskDefinitionInterface

	GetSoftDeletedAt() time.Time
	GetSoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt time.Time) TaskDefinitionInterface

	GetStatus() string
	SetStatus(status string) TaskDefinitionInterface

	GetTitle() string
	SetTitle(title string) TaskDefinitionInterface

	GetUpdatedAt() time.Time
	GetUpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt time.Time) TaskDefinitionInterface
}

type TaskDefinitionQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) TaskDefinitionQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) TaskDefinitionQueryInterface

	HasAlias() bool
	Alias() string
	SetAlias(alias string) TaskDefinitionQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) TaskDefinitionQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) TaskDefinitionQueryInterface

	HasID() bool
	ID() string
	SetID(id string) TaskDefinitionQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) TaskDefinitionQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) TaskDefinitionQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) TaskDefinitionQueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) TaskDefinitionQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) TaskDefinitionQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(withDeleted bool) TaskDefinitionQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) TaskDefinitionQueryInterface

	HasStatusIn() bool
	StatusIn() []string
	SetStatusIn(statusIn []string) TaskDefinitionQueryInterface
}

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
		SetCreatedAt(carbon.Now(carbon.UTC).StdTime()).
		SetUpdatedAt(carbon.Now(carbon.UTC).StdTime()).
		SetSoftDeletedAt(carbon.Parse(MAX_DATETIME, carbon.UTC).StdTime())

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
		o.SetCreatedAt(parseTime(v))
	}
	if v, ok := data[COLUMN_UPDATED_AT]; ok {
		o.SetUpdatedAt(parseTime(v))
	}
	if v, ok := data[COLUMN_SOFT_DELETED_AT]; ok {
		o.SetSoftDeletedAt(parseTime(v))
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

func (o *taskDefinition) GetCreatedAt() time.Time {
	return o.CreatedAtField.CreatedAt
}

func (o *taskDefinition) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt)
}

func (o *taskDefinition) SetCreatedAt(createdAt time.Time) TaskDefinitionInterface {
	o.CreatedAtField.CreatedAt = createdAt
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

func (o *taskDefinition) GetSoftDeletedAt() time.Time {
	return o.SoftDeletesMaxDate.SoftDeletedAt
}

func (o *taskDefinition) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.SoftDeletesMaxDate.SoftDeletedAt)
}

func (o *taskDefinition) SetSoftDeletedAt(deletedAt time.Time) TaskDefinitionInterface {
	o.SoftDeletesMaxDate.SoftDeletedAt = deletedAt
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

func (o *taskDefinition) GetUpdatedAt() time.Time {
	return o.UpdatedAtField.UpdatedAt
}

func (o *taskDefinition) GetUpdatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.UpdatedAtField.UpdatedAt)
}

func (o *taskDefinition) SetUpdatedAt(updatedAt time.Time) TaskDefinitionInterface {
	o.UpdatedAtField.UpdatedAt = updatedAt
	return o
}
