package taskstore

import (
	"strings"
	"sync"

	"github.com/mingrammer/cfmt"
)

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
	errorMessage   string
	infoMessage    string
	successMessage string
}

// LastErrorMessage returns the last error message recorded via LogError.
func (handler *TaskDefinitionHandlerBase) LastErrorMessage() string {
	handler.mu.RLock()
	defer handler.mu.RUnlock()
	return handler.errorMessage
}

// ErrorMessage alias is kept for backwards compatibility.
// Deprecated: use LastErrorMessage instead. Will be removed after 2026-11-30.
func (handler *TaskDefinitionHandlerBase) ErrorMessage() string {
	return handler.LastErrorMessage()
}

// LastInfoMessage returns the last informational message recorded via LogInfo.
func (handler *TaskDefinitionHandlerBase) LastInfoMessage() string {
	handler.mu.RLock()
	defer handler.mu.RUnlock()
	return handler.infoMessage
}

// InfoMessage alias is kept for backwards compatibility.
// Deprecated: use LastInfoMessage instead. Will be removed after 2026-11-30.
func (handler *TaskDefinitionHandlerBase) InfoMessage() string {
	return handler.LastInfoMessage()
}

// LastSuccessMessage returns the last success message recorded via LogSuccess.
func (handler *TaskDefinitionHandlerBase) LastSuccessMessage() string {
	handler.mu.RLock()
	defer handler.mu.RUnlock()
	return handler.successMessage
}

// SuccessMessage alias is kept for backwards compatibility.
// Deprecated: use LastSuccessMessage instead. Will be removed after 2026-11-30.
func (handler *TaskDefinitionHandlerBase) SuccessMessage() string {
	return handler.LastSuccessMessage()
}

// QueuedTask returns the currently associated queued task, if any.
func (handler *TaskDefinitionHandlerBase) QueuedTask() TaskQueueInterface {
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

// Options returns the options map used when the handler is executed directly
// without an associated queued task.
func (handler *TaskDefinitionHandlerBase) Options() map[string]string {
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

// HasQueuedTask reports whether the handler is currently associated with a
// queued task.
func (handler *TaskDefinitionHandlerBase) HasQueuedTask() bool {
	handler.mu.RLock()
	defer handler.mu.RUnlock()
	return handler.queuedTask != nil
}

// LogError records an error message for the handler and either appends it to
// the queued task details (when a queued task is present) or prints it using
// cfmt.Errorln.
func (handler *TaskDefinitionHandlerBase) LogError(message string) {
	handler.mu.Lock()
	handler.errorMessage = message
	qt := handler.queuedTask
	handler.mu.Unlock()

	if qt != nil {
		qt.AppendDetails(message)
	} else {
		_, _ = cfmt.Errorln(message)
	}
}

// LogInfo records an informational message for the handler and either
// appends it to the queued task details (when a queued task is present) or
// prints it using cfmt.Infoln.
func (handler *TaskDefinitionHandlerBase) LogInfo(message string) {
	handler.mu.Lock()
	handler.infoMessage = message
	qt := handler.queuedTask
	handler.mu.Unlock()

	if qt != nil {
		qt.AppendDetails(message)
	} else {
		_, _ = cfmt.Infoln(message)
	}
}

// LogSuccess records a success message for the handler and either appends it
// to the queued task details (when a queued task is present) or prints it
// using cfmt.Successln.
func (handler *TaskDefinitionHandlerBase) LogSuccess(message string) {
	handler.mu.Lock()
	handler.successMessage = message
	qt := handler.queuedTask
	handler.mu.Unlock()

	if qt != nil {
		qt.AppendDetails(message)
	} else {
		_, _ = cfmt.Successln(message)
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
