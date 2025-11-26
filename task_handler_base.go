package taskstore

import (
	"strings"
	"sync"

	"github.com/mingrammer/cfmt"
)

type TaskHandlerBase struct {
	mu             sync.RWMutex
	queuedTask     TaskQueueInterface // dynamic
	options        map[string]string
	errorMessage   string
	infoMessage    string
	successMessage string
}

func (handler *TaskHandlerBase) ErrorMessage() string {
	handler.mu.RLock()
	defer handler.mu.RUnlock()
	return handler.errorMessage
}

func (handler *TaskHandlerBase) InfoMessage() string {
	handler.mu.RLock()
	defer handler.mu.RUnlock()
	return handler.infoMessage
}

func (handler *TaskHandlerBase) SuccessMessage() string {
	handler.mu.RLock()
	defer handler.mu.RUnlock()
	return handler.successMessage
}

func (handler *TaskHandlerBase) QueuedTask() TaskQueueInterface {
	handler.mu.RLock()
	defer handler.mu.RUnlock()
	return handler.queuedTask
}

func (handler *TaskHandlerBase) SetQueuedTask(queuedTask TaskQueueInterface) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	handler.queuedTask = queuedTask
}

func (handler *TaskHandlerBase) Options() map[string]string {
	handler.mu.RLock()
	defer handler.mu.RUnlock()
	return handler.options
}

func (handler *TaskHandlerBase) SetOptions(options map[string]string) {
	handler.mu.Lock()
	defer handler.mu.Unlock()
	handler.options = options
}

func (handler *TaskHandlerBase) HasQueuedTask() bool {
	handler.mu.RLock()
	defer handler.mu.RUnlock()
	return handler.queuedTask != nil
}

func (handler *TaskHandlerBase) LogError(message string) {
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

func (handler *TaskHandlerBase) LogInfo(message string) {
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

func (handler *TaskHandlerBase) LogSuccess(message string) {
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

func (handler *TaskHandlerBase) GetParam(paramName string) string {
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

func (handler *TaskHandlerBase) GetParamArray(paramName string) []string {
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
