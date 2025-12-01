package taskstore

import (
	"github.com/dromara/carbon/v2"
)

type TaskQueueInterface interface {
	// =======================================================================
	// Metadata Methods
	// =======================================================================

	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	// =======================================================================
	// Informational Methods
	// =======================================================================

	IsCanceled() bool
	IsDeleted() bool
	IsFailed() bool
	IsQueued() bool
	IsPaused() bool
	IsRunning() bool
	IsSuccess() bool
	IsSoftDeleted() bool

	// =======================================================================
	// Accessors (Setters and Getters)
	// =======================================================================

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

	QueueName() string
	SetQueueName(queueName string) TaskQueueInterface
}
