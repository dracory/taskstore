package taskstore

import "errors"

func TaskQueueQuery() TaskQueueQueryInterface {
	return &taskQueueQuery{
		properties: make(map[string]interface{}),
	}
}

type taskQueueQuery struct {
	properties map[string]interface{}
}

var _ TaskQueueQueryInterface = (*taskQueueQuery)(nil)

func (q *taskQueueQuery) Validate() error {
	if q.HasCreatedAtGte() && q.CreatedAtGte() == "" {
		return errors.New("queue query. created_at_gte cannot be empty")
	}

	if q.HasCreatedAtLte() && q.CreatedAtLte() == "" {
		return errors.New("queue query. created_at_lte cannot be empty")
	}

	if q.HasID() && q.ID() == "" {
		return errors.New("queue query. id cannot be empty")
	}

	if q.HasIDIn() && len(q.IDIn()) < 1 {
		return errors.New("queue query. id_in cannot be empty array")
	}

	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("queue query. limit cannot be negative")
	}

	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("queue query. offset cannot be negative")
	}

	if q.HasStatus() && q.Status() == "" {
		return errors.New("queue query. status cannot be empty")
	}

	if q.HasStatusIn() && len(q.StatusIn()) < 1 {
		return errors.New("queue query. status_in cannot be empty array")
	}

	if q.HasTaskID() && q.TaskID() == "" {
		return errors.New("queue query. task_id cannot be empty")
	}

	return nil
}

func (q *taskQueueQuery) Columns() []string {
	if !q.hasProperty("columns") {
		return []string{}
	}

	return q.properties["columns"].([]string)
}

func (q *taskQueueQuery) SetColumns(columns []string) TaskQueueQueryInterface {
	q.properties["columns"] = columns
	return q
}

func (q *taskQueueQuery) HasCountOnly() bool {
	return q.hasProperty("count_only")
}

func (q *taskQueueQuery) IsCountOnly() bool {
	if q.HasCountOnly() {
		return q.properties["count_only"].(bool)
	}

	return false
}

func (q *taskQueueQuery) SetCountOnly(countOnly bool) TaskQueueQueryInterface {
	q.properties["count_only"] = countOnly
	return q
}

func (q *taskQueueQuery) HasCreatedAtGte() bool {
	return q.hasProperty("created_at_gte")
}

func (q *taskQueueQuery) CreatedAtGte() string {
	return q.properties["created_at_gte"].(string)
}

func (q *taskQueueQuery) SetCreatedAtGte(createdAtGte string) TaskQueueQueryInterface {
	q.properties["created_at_gte"] = createdAtGte
	return q
}

func (q *taskQueueQuery) HasCreatedAtLte() bool {
	return q.hasProperty("created_at_lte")
}

func (q *taskQueueQuery) CreatedAtLte() string {
	return q.properties["created_at_lte"].(string)
}

func (q *taskQueueQuery) SetCreatedAtLte(createdAtLte string) TaskQueueQueryInterface {
	q.properties["created_at_lte"] = createdAtLte
	return q
}

func (q *taskQueueQuery) HasID() bool {
	return q.hasProperty("id")
}

func (q *taskQueueQuery) ID() string {
	return q.properties["id"].(string)
}

func (q *taskQueueQuery) SetID(id string) TaskQueueQueryInterface {
	q.properties["id"] = id
	return q
}

func (q *taskQueueQuery) HasIDIn() bool {
	return q.hasProperty("id_in")
}

func (q *taskQueueQuery) IDIn() []string {
	return q.properties["id_in"].([]string)
}

func (q *taskQueueQuery) SetIDIn(idIn []string) TaskQueueQueryInterface {
	q.properties["id_in"] = idIn
	return q
}

func (q *taskQueueQuery) HasLimit() bool {
	return q.hasProperty("limit")
}

func (q *taskQueueQuery) Limit() int {
	return q.properties["limit"].(int)
}

func (q *taskQueueQuery) SetLimit(limit int) TaskQueueQueryInterface {
	q.properties["limit"] = limit
	return q
}

func (q *taskQueueQuery) HasTaskID() bool {
	return q.hasProperty("task_id")
}

func (q *taskQueueQuery) TaskID() string {
	return q.properties["task_id"].(string)
}

func (q *taskQueueQuery) SetTaskID(taskID string) TaskQueueQueryInterface {
	q.properties["task_id"] = taskID
	return q
}

func (q *taskQueueQuery) HasOffset() bool {
	return q.hasProperty("offset")
}

func (q *taskQueueQuery) Offset() int {
	return q.properties["offset"].(int)
}

func (q *taskQueueQuery) SetOffset(offset int) TaskQueueQueryInterface {
	q.properties["offset"] = offset
	return q
}

func (q *taskQueueQuery) HasOrderBy() bool {
	return q.hasProperty("order_by")
}

func (q *taskQueueQuery) OrderBy() string {
	return q.properties["order_by"].(string)
}

func (q *taskQueueQuery) SetOrderBy(orderBy string) TaskQueueQueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

func (q *taskQueueQuery) HasSoftDeletedIncluded() bool {
	return q.hasProperty("soft_delete_included")
}

func (q *taskQueueQuery) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}
	return q.properties["soft_delete_included"].(bool)
}

func (q *taskQueueQuery) SetSoftDeletedIncluded(softDeleteIncluded bool) TaskQueueQueryInterface {
	q.properties["soft_delete_included"] = softDeleteIncluded
	return q
}

func (q *taskQueueQuery) HasSortOrder() bool {
	return q.hasProperty("sort_order")
}

func (q *taskQueueQuery) SortOrder() string {
	return q.properties["sort_order"].(string)
}

func (q *taskQueueQuery) SetSortOrder(sortOrder string) TaskQueueQueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}

func (q *taskQueueQuery) HasStatus() bool {
	return q.hasProperty("status")
}

func (q *taskQueueQuery) Status() string {
	return q.properties["status"].(string)
}

func (q *taskQueueQuery) SetStatus(status string) TaskQueueQueryInterface {
	q.properties["status"] = status
	return q
}

func (q *taskQueueQuery) HasStatusIn() bool {
	return q.hasProperty("status_in")
}

func (q *taskQueueQuery) StatusIn() []string {
	return q.properties["status_in"].([]string)
}

func (q *taskQueueQuery) SetStatusIn(statusIn []string) TaskQueueQueryInterface {
	q.properties["status_in"] = statusIn
	return q
}

func (q *taskQueueQuery) hasProperty(key string) bool {
	return q.properties[key] != nil
}
