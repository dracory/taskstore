package taskstore

import (
	"context"

	"github.com/dromara/carbon/v2"
)

type QueueInterface interface {
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
	SetAttempts(attempts int) QueueInterface

	CompletedAt() string
	CompletedAtCarbon() *carbon.Carbon
	SetCompletedAt(completedAt string) QueueInterface

	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) QueueInterface

	Details() string
	AppendDetails(details string) QueueInterface
	SetDetails(details string) QueueInterface

	ID() string
	SetID(id string) QueueInterface

	// Memo() string
	// SetMemo(memo string) QueueInterface

	// Meta(name string) string
	// SetMeta(name string, value string) error
	// Metas() (map[string]string, error)
	// SetMetas(metas map[string]string) error
	// UpsertMetas(metas map[string]string) error

	Output() string
	SetOutput(output string) QueueInterface

	Parameters() string
	SetParameters(parameters string) QueueInterface
	ParametersMap() (map[string]string, error)
	SetParametersMap(parameters map[string]string) (QueueInterface, error)

	SoftDeletedAt() string
	SoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) QueueInterface

	StartedAt() string
	StartedAtCarbon() *carbon.Carbon
	SetStartedAt(startedAt string) QueueInterface

	Status() string
	SetStatus(status string) QueueInterface

	TaskID() string
	SetTaskID(taskID string) QueueInterface

	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) QueueInterface
}

type QueueQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) QueueQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) QueueQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) QueueQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) QueueQueryInterface

	HasID() bool
	ID() string
	SetID(id string) QueueQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) QueueQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) QueueQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) QueueQueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) QueueQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) QueueQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(withDeleted bool) QueueQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) QueueQueryInterface

	HasStatusIn() bool
	StatusIn() []string
	SetStatusIn(statusIn []string) QueueQueryInterface

	HasTaskID() bool
	TaskID() string
	SetTaskID(taskID string) QueueQueryInterface
}

type TaskInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	IsActive() bool
	IsCanceled() bool
	IsSoftDeleted() bool

	Alias() string
	SetAlias(alias string) TaskInterface

	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) TaskInterface

	Description() string
	SetDescription(description string) TaskInterface

	ID() string
	SetID(id string) TaskInterface

	Memo() string
	SetMemo(memo string) TaskInterface

	SoftDeletedAt() string
	SoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) TaskInterface

	Status() string
	SetStatus(status string) TaskInterface

	Title() string
	SetTitle(title string) TaskInterface

	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) TaskInterface
}

type TaskQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) TaskQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) TaskQueryInterface

	HasAlias() bool
	Alias() string
	SetAlias(alias string) TaskQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) TaskQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) TaskQueryInterface

	HasID() bool
	ID() string
	SetID(id string) TaskQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) TaskQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) TaskQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) TaskQueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) TaskQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) TaskQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(withDeleted bool) TaskQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) TaskQueryInterface

	HasStatusIn() bool
	StatusIn() []string
	SetStatusIn(statusIn []string) TaskQueryInterface
}

type TaskHandlerInterface interface {
	Alias() string

	Title() string

	Description() string

	Handle() bool

	SetQueuedTask(queuedTask QueueInterface)

	SetOptions(options map[string]string)
}

type StoreInterface interface {
	AutoMigrate() error
	EnableDebug(debug bool) StoreInterface
	// Start()
	// Stop()

	QueueCount(options QueueQueryInterface) (int64, error)
	QueueCreate(Queue QueueInterface) error
	QueueDelete(Queue QueueInterface) error
	QueueDeleteByID(id string) error
	QueueFindByID(QueueID string) (QueueInterface, error)
	QueueList(query QueueQueryInterface) ([]QueueInterface, error)
	QueueSoftDelete(Queue QueueInterface) error
	QueueSoftDeleteByID(id string) error
	QueueUpdate(Queue QueueInterface) error

	QueueRunGoroutine(ctx context.Context, processSeconds int, unstuckMinutes int)
	QueuedTaskProcess(queuedTask QueueInterface) (bool, error)

	TaskEnqueueByAlias(alias string, parameters map[string]interface{}) (QueueInterface, error)
	TaskExecuteCli(alias string, args []string) bool

	TaskCount(options TaskQueryInterface) (int64, error)
	TaskCreate(Task TaskInterface) error
	TaskDelete(Task TaskInterface) error
	TaskDeleteByID(id string) error
	TaskFindByAlias(alias string) (TaskInterface, error)
	TaskFindByID(id string) (TaskInterface, error)
	TaskList(options TaskQueryInterface) ([]TaskInterface, error)
	TaskSoftDelete(Task TaskInterface) error
	TaskSoftDeleteByID(id string) error
	TaskUpdate(Task TaskInterface) error

	TaskHandlerList() []TaskHandlerInterface
	TaskHandlerAdd(taskHandler TaskHandlerInterface, createIfMissing bool) error
}
