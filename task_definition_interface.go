package taskstore

import "github.com/dromara/carbon/v2"

type TaskDefinitionInterface interface {
	IsActive() bool
	IsCanceled() bool
	IsSoftDeleted() bool

	GetAlias() string
	SetAlias(alias string) TaskDefinitionInterface

	GetCreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) TaskDefinitionInterface

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

	GetSoftDeletedAt() string
	SoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) TaskDefinitionInterface

	GetStatus() string
	SetStatus(status string) TaskDefinitionInterface

	GetTitle() string
	SetTitle(title string) TaskDefinitionInterface

	GetUpdatedAt() string
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
