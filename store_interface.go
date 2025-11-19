package taskstore

import "context"

type StoreInterface interface {
	AutoMigrate() error
	EnableDebug(debug bool) StoreInterface
	SetErrorHandler(handler func(queueName, taskID string, err error)) StoreInterface
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

	QueueRunDefault(ctx context.Context, processSeconds int, unstuckMinutes int)
	QueueRunSerial(ctx context.Context, queueName string, processSeconds int, unstuckMinutes int)
	QueueRunConcurrent(ctx context.Context, queueName string, processSeconds int, unstuckMinutes int)
	QueueStop()
	QueueStopByName(queueName string)
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
