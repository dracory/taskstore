package taskstore

import "errors"

func TaskDefinitionQuery() TaskDefinitionQueryInterface {
	return &taskDefinitionQuery{
		properties: make(map[string]interface{}),
	}
}

type taskDefinitionQuery struct {
	properties map[string]interface{}
}

var _ TaskDefinitionQueryInterface = (*taskDefinitionQuery)(nil)

func (q *taskDefinitionQuery) Validate() error {
	if q.HasAlias() && q.Alias() == "" {
		return errors.New("task query. alias cannot be empty")
	}

	if q.HasCreatedAtGte() && q.CreatedAtGte() == "" {
		return errors.New("task query. created_at_gte cannot be empty")
	}

	if q.HasCreatedAtLte() && q.CreatedAtLte() == "" {
		return errors.New("task query. created_at_lte cannot be empty")
	}

	if q.HasID() && q.ID() == "" {
		return errors.New("task query. id cannot be empty")
	}

	if q.HasIDIn() && len(q.IDIn()) < 1 {
		return errors.New("task query. id_in cannot be empty array")
	}

	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("task query. limit cannot be negative")
	}

	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("task query. offset cannot be negative")
	}

	if q.HasStatus() && q.Status() == "" {
		return errors.New("task query. status cannot be empty")
	}

	if q.HasStatusIn() && len(q.StatusIn()) < 1 {
		return errors.New("task query. status_in cannot be empty array")
	}

	return nil
}

func (q *taskDefinitionQuery) HasAlias() bool {
	return q.hasProperty("alias")
}

func (q *taskDefinitionQuery) Alias() string {
	if !q.hasProperty("alias") {
		return ""
	}

	return q.properties["alias"].(string)
}

func (q *taskDefinitionQuery) SetAlias(alias string) TaskDefinitionQueryInterface {
	q.properties["alias"] = alias
	return q
}

func (q *taskDefinitionQuery) Columns() []string {
	if !q.hasProperty("columns") {
		return []string{}
	}

	return q.properties["columns"].([]string)
}

func (q *taskDefinitionQuery) SetColumns(columns []string) TaskDefinitionQueryInterface {
	q.properties["columns"] = columns
	return q
}

func (q *taskDefinitionQuery) HasCountOnly() bool {
	return q.hasProperty("count_only")
}

func (q *taskDefinitionQuery) IsCountOnly() bool {
	if q.HasCountOnly() {
		return q.properties["count_only"].(bool)
	}

	return false
}

func (q *taskDefinitionQuery) SetCountOnly(countOnly bool) TaskDefinitionQueryInterface {
	q.properties["count_only"] = countOnly
	return q
}

func (q *taskDefinitionQuery) HasCreatedAtGte() bool {
	return q.hasProperty("created_at_gte")
}

func (q *taskDefinitionQuery) CreatedAtGte() string {
	return q.properties["created_at_gte"].(string)
}

func (q *taskDefinitionQuery) SetCreatedAtGte(createdAtGte string) TaskDefinitionQueryInterface {
	q.properties["created_at_gte"] = createdAtGte
	return q
}

func (q *taskDefinitionQuery) HasCreatedAtLte() bool {
	return q.hasProperty("created_at_lte")
}

func (q *taskDefinitionQuery) CreatedAtLte() string {
	return q.properties["created_at_lte"].(string)
}

func (q *taskDefinitionQuery) SetCreatedAtLte(createdAtLte string) TaskDefinitionQueryInterface {
	q.properties["created_at_lte"] = createdAtLte
	return q
}

func (q *taskDefinitionQuery) HasID() bool {
	return q.hasProperty("id")
}

func (q *taskDefinitionQuery) ID() string {
	return q.properties["id"].(string)
}

func (q *taskDefinitionQuery) SetID(id string) TaskDefinitionQueryInterface {
	q.properties["id"] = id
	return q
}

func (q *taskDefinitionQuery) HasIDIn() bool {
	return q.hasProperty("id_in")
}

func (q *taskDefinitionQuery) IDIn() []string {
	return q.properties["id_in"].([]string)
}

func (q *taskDefinitionQuery) SetIDIn(idIn []string) TaskDefinitionQueryInterface {
	q.properties["id_in"] = idIn
	return q
}

func (q *taskDefinitionQuery) HasLimit() bool {
	return q.hasProperty("limit")
}

func (q *taskDefinitionQuery) Limit() int {
	return q.properties["limit"].(int)
}

func (q *taskDefinitionQuery) SetLimit(limit int) TaskDefinitionQueryInterface {
	q.properties["limit"] = limit
	return q
}

func (q *taskDefinitionQuery) HasOffset() bool {
	return q.hasProperty("offset")
}

func (q *taskDefinitionQuery) Offset() int {
	return q.properties["offset"].(int)
}

func (q *taskDefinitionQuery) SetOffset(offset int) TaskDefinitionQueryInterface {
	q.properties["offset"] = offset
	return q
}

func (q *taskDefinitionQuery) HasOrderBy() bool {
	return q.hasProperty("order_by")
}

func (q *taskDefinitionQuery) OrderBy() string {
	return q.properties["order_by"].(string)
}

func (q *taskDefinitionQuery) SetOrderBy(orderBy string) TaskDefinitionQueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

func (q *taskDefinitionQuery) HasSoftDeletedIncluded() bool {
	return q.hasProperty("soft_delete_included")
}

func (q *taskDefinitionQuery) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}
	return q.properties["soft_delete_included"].(bool)
}

func (q *taskDefinitionQuery) SetSoftDeletedIncluded(softDeleteIncluded bool) TaskDefinitionQueryInterface {
	q.properties["soft_delete_included"] = softDeleteIncluded
	return q
}

func (q *taskDefinitionQuery) HasSortOrder() bool {
	return q.hasProperty("sort_order")
}

func (q *taskDefinitionQuery) SortOrder() string {
	return q.properties["sort_order"].(string)
}

func (q *taskDefinitionQuery) SetSortOrder(sortOrder string) TaskDefinitionQueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}

func (q *taskDefinitionQuery) HasStatus() bool {
	return q.hasProperty("status")
}

func (q *taskDefinitionQuery) Status() string {
	return q.properties["status"].(string)
}

func (q *taskDefinitionQuery) SetStatus(status string) TaskDefinitionQueryInterface {
	q.properties["status"] = status
	return q
}

func (q *taskDefinitionQuery) HasStatusIn() bool {
	return q.hasProperty("status_in")
}

func (q *taskDefinitionQuery) StatusIn() []string {
	return q.properties["status_in"].([]string)
}

func (q *taskDefinitionQuery) SetStatusIn(statusIn []string) TaskDefinitionQueryInterface {
	q.properties["status_in"] = statusIn
	return q
}

func (q *taskDefinitionQuery) hasProperty(key string) bool {
	return q.properties[key] != nil
}
