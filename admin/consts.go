package admin

var endpoint = "" // initialized in admin.go

const pathHome = "home"

const pathTaskQueueCreate = "task-queue-create"
const pathTaskQueueDelete = "task-queue-delete"
const pathTaskQueueDetails = "task-queue-details"
const pathTaskQueueManager = "task-queue-manager"
const pathTaskQueueParameters = "task-queue-parameters"
const pathTaskQueueRequeue = "task-queue-requeue"
const pathTaskQueueTaskRestart = "task-queue-task-restart"

// const pathQueueUpdate = "queue-update"

const pathTaskDefinitionCreate = "task-definition-create"
const pathTaskDefinitionManager = "task-definition-manager"
const pathTaskDefinitionUpdate = "task-definition-update"
const pathTaskDefinitionDelete = "task-definition-delete"

const actionModalQueuedTaskFilterShow = "modal-queued-task-filter-show"

// const actionModalQueuedTaskRequeueShow = "modal-queued-task-requeue-show"
// const actionModalQueuedTaskRequeueSubmitted = "modal-queued-task-requeue-submitted"
const actionModalQueuedTaskRestartShow = "modal-queue-task-restart-show"
const actionModalQueuedTaskRestartSubmitted = "modal-queue-task-restart-submitted"

// const actionModalQueuedTaskEnqueueShow = "modal-queued-task-enqueue-show"
// const actionModalQueuedTaskEnqueueSubmitted = "modal-task-enqueue-submitted"
