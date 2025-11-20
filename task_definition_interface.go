package taskstore

import "github.com/dromara/carbon/v2"

type TaskDefinitionInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	IsActive() bool
	IsCanceled() bool
	IsSoftDeleted() bool

	Alias() string
	SetAlias(alias string) TaskDefinitionInterface

	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) TaskDefinitionInterface

	Description() string
	SetDescription(description string) TaskDefinitionInterface

	ID() string
	SetID(id string) TaskDefinitionInterface

	Memo() string
	SetMemo(memo string) TaskDefinitionInterface

	IsRecurring() int
	SetIsRecurring(isRecurring int) TaskDefinitionInterface

	RecurrenceRule() string
	SetRecurrenceRule(recurrenceRule string) TaskDefinitionInterface

	SoftDeletedAt() string
	SoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) TaskDefinitionInterface

	Status() string
	SetStatus(status string) TaskDefinitionInterface

	Title() string
	SetTitle(title string) TaskDefinitionInterface

	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) TaskDefinitionInterface
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
