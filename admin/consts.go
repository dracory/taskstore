package admin

var endpoint = "" // initialized in admin.go

type contextKey string

const keyEndpoint = contextKey("endpoint")

const defaultFavicon = `data:image/x-icon;base64,AAABAAEAEBAQAAEABAAoAQAAFgAAACgAAAAQAAAAIAAAAAEABAAAAAAAgAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAmzKzAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABEQEAAQERAAEAAQABAAEAAQABAQEBEQABAAEREQEAAAERARARAREAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD//wAA//8AAP//AAD//wAA//8AAP//AAD//wAAi6MAALu7AAC6owAAuC8AAIkjAAD//wAA//8AAP//AAD//wAA`

const fieldParameters = "parameters"

const fieldQueueID = "queue_id"
const fieldTaskID = "task_id"
const fieldStatus = "status"
const fieldTitle = "title"
const fieldAlias = "alias"
const fieldDescription = "description"
const fieldDetails = "details"

const fieldFilterQueueID = "filter_queue_id"
const fieldFilterStatus = "filter_status"
const fieldFilterName = "filter_name"
const fieldFilterCreatedFrom = "filter_created_from"
const fieldFilterCreatedTo = "filter_created_to"
const fieldFilterTaskID = "filter_task_id"

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
