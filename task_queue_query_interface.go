package taskstore

type TaskQueueQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) TaskQueueQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) TaskQueueQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) TaskQueueQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) TaskQueueQueryInterface

	HasID() bool
	ID() string
	SetID(id string) TaskQueueQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) TaskQueueQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) TaskQueueQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) TaskQueueQueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) TaskQueueQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) TaskQueueQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(withDeleted bool) TaskQueueQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) TaskQueueQueryInterface

	HasStatusIn() bool
	StatusIn() []string
	SetStatusIn(statusIn []string) TaskQueueQueryInterface

	HasTaskID() bool
	TaskID() string
	SetTaskID(taskID string) TaskQueueQueryInterface

	HasQueueName() bool
	QueueName() string
	SetQueueName(queueName string) TaskQueueQueryInterface
}
