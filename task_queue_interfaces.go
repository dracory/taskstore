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

	GetAttempts() int

	// Attempts alias is kept for backwards compatibility.
	// Deprecated: use GetAttempts instead. Will be removed after 2026-11-30.
	Attempts() int
	SetAttempts(attempts int) TaskQueueInterface

	GetCompletedAt() string

	// CompletedAt alias is kept for backwards compatibility.
	// Deprecated: use GetCompletedAt instead. Will be removed after 2026-11-30.
	CompletedAt() string
	CompletedAtCarbon() *carbon.Carbon
	SetCompletedAt(completedAt string) TaskQueueInterface

	GetCreatedAt() string

	// CreatedAt alias is kept for backwards compatibility.
	// Deprecated: use GetCreatedAt instead. Will be removed after 2026-11-30.
	CreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) TaskQueueInterface

	GetDetails() string

	// Details alias is kept for backwards compatibility.
	// Deprecated: use GetDetails instead. Will be removed after 2026-11-30.
	Details() string
	AppendDetails(details string) TaskQueueInterface
	SetDetails(details string) TaskQueueInterface

	GetID() string

	// ID alias is kept for backwards compatibility.
	// Deprecated: use GetID instead. Will be removed after 2026-11-30.
	ID() string
	SetID(id string) TaskQueueInterface

	// Memo() string
	// SetMemo(memo string) TaskQueueInterface

	// Meta(name string) string
	// SetMeta(name string, value string) error
	// Metas() (map[string]string, error)
	// SetMetas(metas map[string]string) error
	// UpsertMetas(metas map[string]string) error

	GetOutput() string

	// Output alias is kept for backwards compatibility.
	// Deprecated: use GetOutput instead. Will be removed after 2026-11-30.
	Output() string
	SetOutput(output string) TaskQueueInterface

	GetParameters() string

	// Parameters alias is kept for backwards compatibility.
	// Deprecated: use GetParameters instead. Will be removed after 2026-11-30.
	Parameters() string
	SetParameters(parameters string) TaskQueueInterface
	ParametersMap() (map[string]string, error)
	SetParametersMap(parameters map[string]string) (TaskQueueInterface, error)

	GetSoftDeletedAt() string

	// SoftDeletedAt alias is kept for backwards compatibility.
	// Deprecated: use GetSoftDeletedAt instead. Will be removed after 2026-11-30.
	SoftDeletedAt() string
	SoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) TaskQueueInterface

	GetStartedAt() string

	// StartedAt alias is kept for backwards compatibility.
	// Deprecated: use GetStartedAt instead. Will be removed after 2026-11-30.
	StartedAt() string
	StartedAtCarbon() *carbon.Carbon
	SetStartedAt(startedAt string) TaskQueueInterface

	GetStatus() string

	// Status alias is kept for backwards compatibility.
	// Deprecated: use GetStatus instead. Will be removed after 2026-11-30.
	Status() string
	SetStatus(status string) TaskQueueInterface

	GetTaskID() string

	// TaskID alias is kept for backwards compatibility.
	// Deprecated: use GetTaskID instead. Will be removed after 2026-11-30.
	TaskID() string
	SetTaskID(taskID string) TaskQueueInterface

	GetUpdatedAt() string

	// UpdatedAt alias is kept for backwards compatibility.
	// Deprecated: use GetUpdatedAt instead. Will be removed after 2026-11-30.
	UpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) TaskQueueInterface

	GetQueueName() string

	// QueueName alias is kept for backwards compatibility.
	// Deprecated: use GetQueueName instead. Will be removed after 2026-11-30.
	QueueName() string
	SetQueueName(queueName string) TaskQueueInterface
}
