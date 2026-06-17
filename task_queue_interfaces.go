package taskstore

import (
	"github.com/dromara/carbon/v2"
)

type TaskQueueInterface interface {
	IsCanceled() bool
	IsDeleted() bool
	IsFailed() bool
	IsQueued() bool
	IsPaused() bool
	IsRunning() bool
	IsSuccess() bool
	IsSoftDeleted() bool

	GetAttempts() int
	SetAttempts(attempts int) TaskQueueInterface

	GetCompletedAt() string
	CompletedAtCarbon() *carbon.Carbon
	SetCompletedAt(completedAt string) TaskQueueInterface

	GetCreatedAt() string
	CreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) TaskQueueInterface

	GetDetails() string
	AppendDetails(details string) TaskQueueInterface
	SetDetails(details string) TaskQueueInterface

	GetID() string
	SetID(id string) TaskQueueInterface

	GetOutput() string
	SetOutput(output string) TaskQueueInterface

	GetParameters() string
	SetParameters(parameters string) TaskQueueInterface
	ParametersMap() (map[string]string, error)
	SetParametersMap(parameters map[string]string) (TaskQueueInterface, error)

	GetSoftDeletedAt() string
	SoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) TaskQueueInterface

	GetStartedAt() string
	StartedAtCarbon() *carbon.Carbon
	SetStartedAt(startedAt string) TaskQueueInterface

	GetStatus() string
	SetStatus(status string) TaskQueueInterface

	GetTaskID() string
	SetTaskID(taskID string) TaskQueueInterface

	GetUpdatedAt() string
	UpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) TaskQueueInterface

	GetQueueName() string
	SetQueueName(queueName string) TaskQueueInterface
}
