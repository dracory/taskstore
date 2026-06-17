package taskstore

import "errors"

func NewScheduleQuery() ScheduleQueryInterface {
	return &scheduleQuery{
		properties: make(map[string]interface{}),
	}
}

type scheduleQuery struct {
	properties map[string]interface{}
}

var _ ScheduleQueryInterface = (*scheduleQuery)(nil)

func (q *scheduleQuery) Validate() error {
	if q.HasID() && q.ID() == "" {
		return errors.New("schedule query. id cannot be empty")
	}

	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("schedule query. limit cannot be negative")
	}

	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("schedule query. offset cannot be negative")
	}

	return nil
}

func (q *scheduleQuery) ID() string {
	if !q.hasProperty("id") {
		return ""
	}
	return q.properties["id"].(string)
}

func (q *scheduleQuery) SetID(id string) ScheduleQueryInterface {
	q.properties["id"] = id
	return q
}

func (q *scheduleQuery) HasID() bool {
	return q.hasProperty("id")
}

func (q *scheduleQuery) Name() string {
	if !q.hasProperty("name") {
		return ""
	}
	return q.properties["name"].(string)
}

func (q *scheduleQuery) SetName(name string) ScheduleQueryInterface {
	q.properties["name"] = name
	return q
}

func (q *scheduleQuery) Status() string {
	if !q.hasProperty("status") {
		return ""
	}
	return q.properties["status"].(string)
}

func (q *scheduleQuery) SetStatus(status string) ScheduleQueryInterface {
	q.properties["status"] = status
	return q
}

func (q *scheduleQuery) QueueName() string {
	if !q.hasProperty("queue_name") {
		return ""
	}
	return q.properties["queue_name"].(string)
}

func (q *scheduleQuery) SetQueueName(queueName string) ScheduleQueryInterface {
	q.properties["queue_name"] = queueName
	return q
}

func (q *scheduleQuery) TaskDefinitionID() string {
	if !q.hasProperty("task_definition_id") {
		return ""
	}
	return q.properties["task_definition_id"].(string)
}

func (q *scheduleQuery) SetTaskDefinitionID(taskDefinitionID string) ScheduleQueryInterface {
	q.properties["task_definition_id"] = taskDefinitionID
	return q
}

func (q *scheduleQuery) HasLimit() bool {
	return q.hasProperty("limit")
}

func (q *scheduleQuery) Limit() int {
	if !q.hasProperty("limit") {
		return 0
	}
	return q.properties["limit"].(int)
}

func (q *scheduleQuery) SetLimit(limit int) ScheduleQueryInterface {
	q.properties["limit"] = limit
	return q
}

func (q *scheduleQuery) HasOffset() bool {
	return q.hasProperty("offset")
}

func (q *scheduleQuery) Offset() int {
	if !q.hasProperty("offset") {
		return 0
	}
	return q.properties["offset"].(int)
}

func (q *scheduleQuery) SetOffset(offset int) ScheduleQueryInterface {
	q.properties["offset"] = offset
	return q
}

func (q *scheduleQuery) hasProperty(key string) bool {
	return q.properties[key] != nil
}
