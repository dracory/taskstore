package taskstore

import "github.com/dromara/carbon/v2"

type TaskDefinitionInterface interface {

	// =======================================================================
	// Metadata Methods
	// =======================================================================

	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	// =======================================================================
	// Informational Methods
	// =======================================================================

	IsActive() bool
	IsCanceled() bool
	IsSoftDeleted() bool

	// =======================================================================
	// Accessors (Setters and Getters)
	// =======================================================================

	GetAlias() string

	// Alias alias is kept for backwards compatibility.
	// Deprecated: use GetAlias instead. Will be removed after 2026-11-30.
	Alias() string
	SetAlias(alias string) TaskDefinitionInterface

	GetCreatedAt() string

	// CreatedAt alias is kept for backwards compatibility.
	// Deprecated: use GetCreatedAt instead. Will be removed after 2026-11-30.
	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) TaskDefinitionInterface

	GetDescription() string

	// Description alias is kept for backwards compatibility.
	// Deprecated: use GetDescription instead. Will be removed after 2026-11-30.
	Description() string
	SetDescription(description string) TaskDefinitionInterface

	GetID() string

	// ID alias is kept for backwards compatibility.
	// Deprecated: use GetID instead. Will be removed after 2026-11-30.
	ID() string
	SetID(id string) TaskDefinitionInterface

	GetMemo() string

	// Memo alias is kept for backwards compatibility.
	// Deprecated: use GetMemo instead. Will be removed after 2026-11-30.
	Memo() string
	SetMemo(memo string) TaskDefinitionInterface

	GetIsRecurring() int

	// IsRecurring alias is kept for backwards compatibility.
	// Deprecated: use GetIsRecurring instead. Will be removed after 2026-11-30.
	IsRecurring() int
	SetIsRecurring(isRecurring int) TaskDefinitionInterface

	GetRecurrenceRule() string

	// RecurrenceRule alias is kept for backwards compatibility.
	// Deprecated: use GetRecurrenceRule instead. Will be removed after 2026-11-30.
	RecurrenceRule() string
	SetRecurrenceRule(recurrenceRule string) TaskDefinitionInterface

	GetSoftDeletedAt() string

	// SoftDeletedAt alias is kept for backwards compatibility.
	// Deprecated: use GetSoftDeletedAt instead. Will be removed after 2026-11-30.
	SoftDeletedAt() string
	SoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) TaskDefinitionInterface

	GetStatus() string

	// Status alias is kept for backwards compatibility.
	// Deprecated: use GetStatus instead. Will be removed after 2026-11-30.
	Status() string
	SetStatus(status string) TaskDefinitionInterface

	GetTitle() string

	// Title alias is kept for backwards compatibility.
	// Deprecated: use GetTitle instead. Will be removed after 2026-11-30.
	Title() string
	SetTitle(title string) TaskDefinitionInterface

	GetUpdatedAt() string

	// UpdatedAt alias is kept for backwards compatibility.
	// Deprecated: use GetUpdatedAt instead. Will be removed after 2026-11-30.
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
