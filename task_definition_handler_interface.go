package taskstore

// TaskHandlerInterface alias is kept for backwards compatibility.
// Deprecated: use TaskDefinitionHandlerInterface instead. Will be removed after 2026-11-30.
type TaskHandlerInterface = TaskDefinitionHandlerInterface

// TaskDefinitionHandlerInterface defines the contract for a task definition
// handler implementation. Handlers provide metadata (alias, title,
// description), implement the task logic in Handle, and support being wired
// to a queued task or executed directly with options.
type TaskDefinitionHandlerInterface interface {
	// =======================================================================
	// Metadata Methods
	// =======================================================================

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

	// =======================================================================
	// Informational Methods
	// =======================================================================

	// HasQueuedTask reports whether the handler is currently associated with a
	// queued task.
	HasQueuedTask() bool

	// LogError records an error message for the handler and either appends it
	// to the queued task details or logs it directly.
	LogError(message string)

	// LogInfo records an informational message for the handler and either
	// appends it to the queued task details or logs it directly.
	LogInfo(message string)

	// LogSuccess records a success message for the handler and either appends
	// it to the queued task details or logs it directly.
	LogSuccess(message string)

	// =======================================================================
	// Accessors (Setters and Getters)
	// =======================================================================

	// GetQueuedTask returns the currently associated queued task, if any.
	GetQueuedTask() TaskQueueInterface

	// SetQueuedTask associates the handler with a queued task instance when
	// invoked as part of background processing.
	SetQueuedTask(queuedTask TaskQueueInterface)

	// GetOptions returns the options map used when the handler is executed
	// directly without an associated queued task.
	GetOptions() map[string]string

	// SetOptions provides key-value options when the handler is executed
	// directly, without an associated queued task.
	SetOptions(options map[string]string)

	// GetOutput returns the current output value for the handler. When a
	// queued task is associated, this typically reflects the queued task's
	// output; otherwise it is a handler-local value.
	GetOutput() string

	// SetOutput stores the output value for the handler. When a queued task is
	// associated, implementations should propagate this value to the queued
	// task's output as well.
	SetOutput(output string)

	// GetLastErrorMessage returns the last error message recorded during
	// handler execution.
	GetLastErrorMessage() string

	// GetLastInfoMessage returns the last informational message recorded during
	// handler execution.
	GetLastInfoMessage() string

	// GetLastSuccessMessage returns the last success message recorded during
	// handler execution.
	GetLastSuccessMessage() string

	// GetParam returns the value of a named parameter for the current
	// execution, reading from the queued task parameters when present or from
	// the handler options otherwise.
	GetParam(paramName string) string

	// GetParamArray returns the named parameter split on semicolons into a
	// slice. If the parameter is missing or empty, it returns an empty slice.
	GetParamArray(paramName string) []string
}
