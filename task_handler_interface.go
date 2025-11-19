package taskstore

type TaskHandlerInterface interface {
	Alias() string

	Title() string

	Description() string

	Handle() bool

	SetQueuedTask(queuedTask TaskQueueInterface)

	SetOptions(options map[string]string)
}
