package taskstore

// TaskHandlerInterface alias is kept for backwards compatibility.
// Deprecated: use TaskDefinitionHandlerInterface instead. Will be removed after 2026-11-30.
type TaskHandlerInterface = TaskDefinitionHandlerInterface

// TaskDefinitionHandlerInterface defines the contract for a task definition
// handler implementation. Handlers provide metadata (alias, title,
// description), implement the task logic in Handle, and support being wired
// to a queued task or executed directly with options.
type TaskDefinitionHandlerInterface interface {
	// Alias returns the unique identifier used to reference the task
	// definition when enqueuing or executing it.
	Alias() string

	// Title returns a short human-readable name for the task.
	Title() string

	// Description returns a longer human-readable description of what the
	// task does.
	Description() string

	// Handle executes the task logic and returns true on success.
	Handle() bool

	// SetQueuedTask associates the handler with a queued task instance when
	// invoked as part of background processing.
	SetQueuedTask(queuedTask TaskQueueInterface)

	// SetOptions provides key-value options when the handler is executed
	// directly, without an associated queued task.
	SetOptions(options map[string]string)
}
