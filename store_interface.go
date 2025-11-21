package taskstore

import "context"

type StoreInterface interface {
	AutoMigrate() error
	EnableDebug(debug bool) StoreInterface
	SetErrorHandler(handler func(queueName, taskID string, err error)) StoreInterface

	TaskQueueCount(ctx context.Context, options TaskQueueQueryInterface) (int64, error)
	TaskQueueCreate(ctx context.Context, TaskQueue TaskQueueInterface) error
	TaskQueueDelete(ctx context.Context, TaskQueue TaskQueueInterface) error
	TaskQueueDeleteByID(ctx context.Context, id string) error
	TaskQueueFindByID(ctx context.Context, TaskQueueID string) (TaskQueueInterface, error)
	TaskQueueList(ctx context.Context, query TaskQueueQueryInterface) ([]TaskQueueInterface, error)
	TaskQueueSoftDelete(ctx context.Context, TaskQueue TaskQueueInterface) error
	TaskQueueSoftDeleteByID(ctx context.Context, id string) error
	TaskQueueUpdate(ctx context.Context, TaskQueue TaskQueueInterface) error

	QueueRunDefault(ctx context.Context, processSeconds int, unstuckMinutes int)
	QueueRunSerial(ctx context.Context, queueName string, processSeconds int, unstuckMinutes int)
	QueueRunConcurrent(ctx context.Context, queueName string, processSeconds int, unstuckMinutes int)
	QueueStop()
	QueueStopByName(queueName string)
	QueuedTaskProcess(ctx context.Context, queuedTask TaskQueueInterface) (bool, error)

	TaskEnqueueByAlias(ctx context.Context, alias string, parameters map[string]interface{}) (TaskQueueInterface, error)
	TaskExecuteCli(alias string, args []string) bool

	TaskDefinitionCount(ctx context.Context, options TaskDefinitionQueryInterface) (int64, error)
	TaskDefinitionCreate(ctx context.Context, TaskDefinition TaskDefinitionInterface) error
	TaskDefinitionDelete(ctx context.Context, TaskDefinition TaskDefinitionInterface) error
	TaskDefinitionDeleteByID(ctx context.Context, id string) error
	TaskDefinitionFindByAlias(ctx context.Context, alias string) (TaskDefinitionInterface, error)
	TaskDefinitionFindByID(ctx context.Context, id string) (TaskDefinitionInterface, error)
	TaskDefinitionList(ctx context.Context, options TaskDefinitionQueryInterface) ([]TaskDefinitionInterface, error)
	TaskDefinitionSoftDelete(ctx context.Context, TaskDefinition TaskDefinitionInterface) error
	TaskDefinitionSoftDeleteByID(ctx context.Context, id string) error
	TaskDefinitionUpdate(ctx context.Context, TaskDefinition TaskDefinitionInterface) error

	TaskHandlerList() []TaskHandlerInterface
	TaskHandlerAdd(ctx context.Context, taskHandler TaskHandlerInterface, createIfMissing bool) error
}
