package taskstore

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// == INTERFACE =================================================================

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

// == CONTEXT INTERFACE ========================================================

// TaskHandlerWithContext is an optional interface that task handlers can implement
// to receive context for cancellation support. This is backward compatible - handlers
// that don't implement this will continue to work using the standard Handle() method.
//
// Example usage:
//
//	type MyHandler struct {
//	    TaskDefinitionHandlerBase
//	}
//
//	func (h *MyHandler) HandleWithContext(ctx context.Context) bool {
//	    select {
//	    case <-ctx.Done():
//	        h.LogInfo("Task cancelled")
//	        return false
//	    case <-time.After(5 * time.Second):
//	        h.LogSuccess("Task completed")
//	        return true
//	    }
//	}
type TaskHandlerWithContext interface {
	TaskDefinitionHandlerInterface
	HandleWithContext(ctx context.Context) bool
}

// == BASE IMPLEMENTATION ======================================================

// TaskHandlerBase alias is kept for backwards compatibility.
// Deprecated: use TaskDefinitionHandlerBase instead. Will be removed after 2026-11-30.
type TaskHandlerBase = TaskDefinitionHandlerBase

// TaskDefinitionHandlerBase provides concurrency-safe shared behavior for task
// definition handlers, including access to the current queued task, parameter
// lookup and logging of error, info and success messages during task
// execution.
type TaskDefinitionHandlerBase struct {
	mu             sync.RWMutex
	queuedTask     TaskQueueInterface // dynamic
	options        map[string]string
	output         string
	errorMessage   string
	infoMessage    string
	successMessage string
}

// GetLastErrorMessage returns the last error message recorded via LogError.
func (handler *TaskDefinitionHandlerBase) GetLastErrorMessage() string {
	handler.mu.RLock()
	defer handler.mu.RUnlock()
	return handler.errorMessage
}

// GetLastInfoMessage returns the last informational message recorded via LogInfo.
func (handler *TaskDefinitionHandlerBase) GetLastInfoMessage() string {
	handler.mu.RLock()
	defer handler.mu.RUnlock()
	return handler.infoMessage
}

// GetLastSuccessMessage returns the last success message recorded via LogSuccess.
func (handler *TaskDefinitionHandlerBase) GetLastSuccessMessage() string {
	handler.mu.RLock()
	defer handler.mu.RUnlock()
	return handler.successMessage
}

// GetQueuedTask returns the currently associated queued task, if any.
func (handler *TaskDefinitionHandlerBase) GetQueuedTask() TaskQueueInterface {
	handler.mu.RLock()
	defer handler.mu.RUnlock()
	return handler.queuedTask
}

// SetQueuedTask associates the handler with a specific queued task instance.
func (handler *TaskDefinitionHandlerBase) SetQueuedTask(queuedTask TaskQueueInterface) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	handler.queuedTask = queuedTask
}

// GetOptions returns the options map used when the handler is executed directly
// without an associated queued task.
func (handler *TaskDefinitionHandlerBase) GetOptions() map[string]string {
	handler.mu.RLock()
	defer handler.mu.RUnlock()
	return handler.options
}

// SetOptions sets the options map used when the handler is executed directly
// without an associated queued task.
func (handler *TaskDefinitionHandlerBase) SetOptions(options map[string]string) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	handler.options = options
}

func (handler *TaskDefinitionHandlerBase) GetOutput() string {
	handler.mu.RLock()
	qt := handler.queuedTask
	output := handler.output
	handler.mu.RUnlock()

	if qt != nil {
		return qt.GetOutput()
	}

	return output
}

func (handler *TaskDefinitionHandlerBase) SetOutput(output string) {
	handler.mu.Lock()
	handler.output = output
	qt := handler.queuedTask
	handler.mu.Unlock()

	if qt != nil {
		qt.SetOutput(output)
	}
}

// HasQueuedTask reports whether the handler is currently associated with a
// queued task.
func (handler *TaskDefinitionHandlerBase) HasQueuedTask() bool {
	handler.mu.RLock()
	defer handler.mu.RUnlock()
	return handler.queuedTask != nil
}

// LogError records an error message for the handler and either appends it to
// the queued task details (when a queued task is present) or prints it using
// fmt.Println.
func (handler *TaskDefinitionHandlerBase) LogError(message string) {
	handler.mu.Lock()
	handler.errorMessage = message
	qt := handler.queuedTask
	handler.mu.Unlock()

	if qt != nil {
		qt.AppendDetails(message)
	} else {
		fmt.Println("ERROR:", message)
	}
}

// LogInfo records an informational message for the handler and either
// appends it to the queued task details (when a queued task is present) or
// prints it using fmt.Println.
func (handler *TaskDefinitionHandlerBase) LogInfo(message string) {
	handler.mu.Lock()
	handler.infoMessage = message
	qt := handler.queuedTask
	handler.mu.Unlock()

	if qt != nil {
		qt.AppendDetails(message)
	} else {
		fmt.Println("INFO:", message)
	}
}

// LogSuccess records a success message for the handler and either appends it
// to the queued task details (when a queued task is present) or prints it
// using fmt.Println.
func (handler *TaskDefinitionHandlerBase) LogSuccess(message string) {
	handler.mu.Lock()
	handler.successMessage = message
	qt := handler.queuedTask
	handler.mu.Unlock()

	if qt != nil {
		qt.AppendDetails(message)
	} else {
		fmt.Println("SUCCESS:", message)
	}
}

// GetParam returns the value of a named parameter for the current execution.
// When a queued task is present it reads from the task's parameter map;
// otherwise it falls back to the handler options. If the parameter is
// missing or the queued task parameters cannot be decoded, an empty string
// is returned.
func (handler *TaskDefinitionHandlerBase) GetParam(paramName string) string {
	handler.mu.RLock()
	qt := handler.queuedTask
	opts := handler.options
	handler.mu.RUnlock()

	if qt != nil {
		parameters, parametersErr := qt.ParametersMap()

		if parametersErr != nil {
			qt.AppendDetails("Parameters JSON incorrect. " + parametersErr.Error())
			return ""
		}

		parameter, parameterExists := parameters[paramName]

		if !parameterExists {
			return ""
		}

		return parameter
	} else {
		return opts[paramName]
	}
}

// GetParamArray returns the named parameter split on semicolons into a slice.
// If the parameter is missing or empty, it returns an empty slice.
func (handler *TaskDefinitionHandlerBase) GetParamArray(paramName string) []string {
	param := handler.GetParam(paramName)

	if param == "" {
		return []string{}
	}

	result := strings.Split(param, ";")
	if result == nil {
		return []string{}
	}

	return result
}

// == DEPRECATED ALIASES =======================================================

// LastErrorMessage alias is kept for backwards compatibility.
// Deprecated: use GetLastErrorMessage instead. Will be removed after 2026-11-30.
func (handler *TaskDefinitionHandlerBase) LastErrorMessage() string {
	return handler.GetLastErrorMessage()
}

// ErrorMessage alias is kept for backwards compatibility.
// Deprecated: use GetLastErrorMessage instead. Will be removed after 2026-11-30.
func (handler *TaskDefinitionHandlerBase) ErrorMessage() string {
	return handler.GetLastErrorMessage()
}

// LastInfoMessage alias is kept for backwards compatibility.
// Deprecated: use GetLastInfoMessage instead. Will be removed after 2026-11-30.
func (handler *TaskDefinitionHandlerBase) LastInfoMessage() string {
	return handler.GetLastInfoMessage()
}

// InfoMessage alias is kept for backwards compatibility.
// Deprecated: use GetLastInfoMessage instead. Will be removed after 2026-11-30.
func (handler *TaskDefinitionHandlerBase) InfoMessage() string {
	return handler.GetLastInfoMessage()
}

// LastSuccessMessage alias is kept for backwards compatibility.
// Deprecated: use GetLastSuccessMessage instead. Will be removed after 2026-11-30.
func (handler *TaskDefinitionHandlerBase) LastSuccessMessage() string {
	return handler.GetLastSuccessMessage()
}

// SuccessMessage alias is kept for backwards compatibility.
// Deprecated: use GetLastSuccessMessage instead. Will be removed after 2026-11-30.
func (handler *TaskDefinitionHandlerBase) SuccessMessage() string {
	return handler.GetLastSuccessMessage()
}

// QueuedTask alias is kept for backwards compatibility.
// Deprecated: use GetQueuedTask instead. Will be removed after 2026-11-30.
func (handler *TaskDefinitionHandlerBase) QueuedTask() TaskQueueInterface {
	return handler.GetQueuedTask()
}

// Options alias is kept for backwards compatibility.
// Deprecated: use GetOptions instead. Will be removed after 2026-11-30.
func (handler *TaskDefinitionHandlerBase) Options() map[string]string {
	return handler.GetOptions()
}
