package taskstore

import (
	"context"

	"github.com/dromara/carbon/v2"
)

type TaskQueueInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	IsCanceled() bool
	IsDeleted() bool
	IsFailed() bool
	IsQueued() bool
	IsPaused() bool
	IsRunning() bool
	IsSuccess() bool
	IsSoftDeleted() bool

	Attempts() int
	SetAttempts(attempts int) TaskQueueInterface

	CompletedAt() string
	CompletedAtCarbon() *carbon.Carbon
	SetCompletedAt(completedAt string) TaskQueueInterface

	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) TaskQueueInterface

	Details() string
	AppendDetails(details string) TaskQueueInterface
	SetDetails(details string) TaskQueueInterface

	ID() string
	SetID(id string) TaskQueueInterface

	// Memo() string
	// SetMemo(memo string) TaskQueueInterface

	// Meta(name string) string
	// SetMeta(name string, value string) error
	// Metas() (map[string]string, error)
	// SetMetas(metas map[string]string) error
	// UpsertMetas(metas map[string]string) error

	Output() string
	SetOutput(output string) TaskQueueInterface

	Parameters() string
	SetParameters(parameters string) TaskQueueInterface
	ParametersMap() (map[string]string, error)
	SetParametersMap(parameters map[string]string) (TaskQueueInterface, error)

	SoftDeletedAt() string
	SoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) TaskQueueInterface

	StartedAt() string
	StartedAtCarbon() *carbon.Carbon
	SetStartedAt(startedAt string) TaskQueueInterface

	Status() string
	SetStatus(status string) TaskQueueInterface

	TaskID() string
	SetTaskID(taskID string) TaskQueueInterface

	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) TaskQueueInterface
}

type TaskQueueQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) TaskQueueQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) TaskQueueQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) TaskQueueQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) TaskQueueQueryInterface

	HasID() bool
	ID() string
	SetID(id string) TaskQueueQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) TaskQueueQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) TaskQueueQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) TaskQueueQueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) TaskQueueQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) TaskQueueQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(withDeleted bool) TaskQueueQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) TaskQueueQueryInterface

	HasStatusIn() bool
	StatusIn() []string
	SetStatusIn(statusIn []string) TaskQueueQueryInterface

	HasTaskID() bool
	TaskID() string
	SetTaskID(taskID string) TaskQueueQueryInterface
}

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

type TaskHandlerInterface interface {
	Alias() string

	Title() string

	Description() string

	Handle() bool

	SetQueuedTask(queuedTask TaskQueueInterface)

	SetOptions(options map[string]string)
}

type StoreInterface interface {
	AutoMigrate() error
	EnableDebug(debug bool) StoreInterface
	// Start()
	// Stop()

	TaskQueueCount(options TaskQueueQueryInterface) (int64, error)
	TaskQueueCreate(TaskQueue TaskQueueInterface) error
	TaskQueueDelete(TaskQueue TaskQueueInterface) error
	TaskQueueDeleteByID(id string) error
	TaskQueueFindByID(TaskQueueID string) (TaskQueueInterface, error)
	TaskQueueList(query TaskQueueQueryInterface) ([]TaskQueueInterface, error)
	TaskQueueSoftDelete(TaskQueue TaskQueueInterface) error
	TaskQueueSoftDeleteByID(id string) error
	TaskQueueUpdate(TaskQueue TaskQueueInterface) error

	QueueRunGoroutine(ctx context.Context, processSeconds int, unstuckMinutes int)
	QueuedTaskProcess(queuedTask TaskQueueInterface) (bool, error)

	TaskEnqueueByAlias(alias string, parameters map[string]interface{}) (TaskQueueInterface, error)
	TaskExecuteCli(alias string, args []string) bool

	TaskDefinitionCount(options TaskDefinitionQueryInterface) (int64, error)
	TaskDefinitionCreate(TaskDefinition TaskDefinitionInterface) error
	TaskDefinitionDelete(TaskDefinition TaskDefinitionInterface) error
	TaskDefinitionDeleteByID(id string) error
	TaskDefinitionFindByAlias(alias string) (TaskDefinitionInterface, error)
	TaskDefinitionFindByID(id string) (TaskDefinitionInterface, error)
	TaskDefinitionList(options TaskDefinitionQueryInterface) ([]TaskDefinitionInterface, error)
	TaskDefinitionSoftDelete(TaskDefinition TaskDefinitionInterface) error
	TaskDefinitionSoftDeleteByID(id string) error
	TaskDefinitionUpdate(TaskDefinition TaskDefinitionInterface) error

	TaskHandlerList() []TaskHandlerInterface
	TaskHandlerAdd(taskHandler TaskHandlerInterface, createIfMissing bool) error
}
