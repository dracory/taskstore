package taskstore

import (
	"strings"

	"github.com/mingrammer/cfmt"
)

type TaskHandlerBase struct {
	queuedTask     TaskQueueInterface // dynamic
	options        map[string]string
	errorMessage   string
	infoMessage    string
	successMessage string
}

func (handler *TaskHandlerBase) ErrorMessage() string {
	return handler.errorMessage
}

func (handler *TaskHandlerBase) InfoMessage() string {
	return handler.infoMessage
}

func (handler *TaskHandlerBase) SuccessMessage() string {
	return handler.successMessage
}

func (handler *TaskHandlerBase) QueuedTask() TaskQueueInterface {
	return handler.queuedTask
}

func (handler *TaskHandlerBase) SetQueuedTask(queuedTask TaskQueueInterface) {
	handler.queuedTask = queuedTask
}

func (handler *TaskHandlerBase) Options() map[string]string {
	return handler.options
}

func (handler *TaskHandlerBase) SetOptions(options map[string]string) {
	handler.options = options
}

func (handler *TaskHandlerBase) HasQueuedTask() bool {
	return handler.queuedTask != nil
}

func (handler *TaskHandlerBase) LogError(message string) {
	handler.errorMessage = message
	if handler.HasQueuedTask() {
		handler.queuedTask.AppendDetails(message)
	} else {
		cfmt.Errorln(message)
	}
}

func (handler *TaskHandlerBase) LogInfo(message string) {
	handler.infoMessage = message
	if handler.HasQueuedTask() {
		handler.queuedTask.AppendDetails(message)
	} else {
		cfmt.Infoln(message)
	}
}

func (handler *TaskHandlerBase) LogSuccess(message string) {
	handler.successMessage = message
	if handler.HasQueuedTask() {
		handler.queuedTask.AppendDetails(message)
	} else {
		cfmt.Successln(message)
	}
}

func (handler *TaskHandlerBase) GetParam(paramName string) string {
	if handler.queuedTask != nil {
		parameters, parametersErr := handler.queuedTask.ParametersMap()

		if parametersErr != nil {
			handler.queuedTask.AppendDetails("Parameters JSON incorrect. " + parametersErr.Error())
			return ""
		}

		parameter, parameterExists := parameters[paramName]

		if !parameterExists {
			return ""
		}

		return parameter
	} else {
		return handler.options[paramName]
	}
}

func (handler *TaskHandlerBase) GetParamArray(paramName string) []string {
	param := handler.GetParam(paramName)

	if param == "" {
		return []string{}
	}

	return strings.Split(param, ";")
}
