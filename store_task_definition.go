package taskstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
)

func (store *Store) TaskDefinitionCount(ctx context.Context, options TaskDefinitionQueryInterface) (int64, error) {
	if options == nil {
		return 0, errors.New("task definition query: cannot be nil")
	}
	if err := options.Validate(); err != nil {
		return 0, err
	}
	q := store.buildTaskDefinitionQuery(options)
	var count int64
	err := q.Table(store.taskDefinitionTableName).Count(&count)
	return count, err
}

func (store *Store) TaskDefinitionCreate(ctx context.Context, task TaskDefinitionInterface) error {
	if task == nil {
		return errors.New("taskstore: task is nil")
	}

	if task.GetCreatedAt() == "" {
		task.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}
	if task.GetUpdatedAt() == "" {
		task.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}
	if task.GetSoftDeletedAt() == "" {
		task.SetSoftDeletedAt(MAX_DATETIME)
	}

	row := map[string]any{
		COLUMN_ID:              task.GetID(),
		COLUMN_STATUS:          task.GetStatus(),
		COLUMN_ALIAS:           task.GetAlias(),
		COLUMN_TITLE:           task.GetTitle(),
		COLUMN_DESCRIPTION:     task.GetDescription(),
		COLUMN_MEMO:            task.GetMemo(),
		COLUMN_IS_RECURRING:    task.GetIsRecurring(),
		COLUMN_RECURRENCE_RULE: task.GetRecurrenceRule(),
		COLUMN_CREATED_AT:      task.GetCreatedAt(),
		COLUMN_UPDATED_AT:      task.GetUpdatedAt(),
		COLUMN_SOFT_DELETED_AT: task.GetSoftDeletedAt(),
	}

	return store.db.Query().Table(store.taskDefinitionTableName).Create(row)
}

func (store *Store) TaskDefinitionDelete(ctx context.Context, task TaskDefinitionInterface) error {
	if task == nil {
		return errors.New("task is nil")
	}
	return store.TaskDefinitionDeleteByID(ctx, task.GetID())
}

func (store *Store) TaskDefinitionDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("task id is empty")
	}
	_, err := store.db.Query().
		Table(store.taskDefinitionTableName).
		Where(COLUMN_ID+" = ?", id).
		Delete()
	return err
}

func (store *Store) TaskDefinitionFindByAlias(ctx context.Context, alias string) (TaskDefinitionInterface, error) {
	if alias == "" {
		return nil, errors.New("task alias is empty")
	}
	query := TaskDefinitionQuery().SetAlias(alias).SetLimit(1)
	list, err := store.TaskDefinitionList(ctx, query)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

func (store *Store) TaskDefinitionFindByID(ctx context.Context, id string) (TaskDefinitionInterface, error) {
	if id == "" {
		return nil, errors.New("task id is empty")
	}
	q := store.db.Query().Model(&taskDefinition{}).Table(store.taskDefinitionTableName).
		Where(COLUMN_ID+" = ?", id)

	var task taskDefinition
	if err := q.First(&task); err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "no rows found" {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (store *Store) TaskDefinitionList(ctx context.Context, options TaskDefinitionQueryInterface) ([]TaskDefinitionInterface, error) {
	if options == nil {
		return []TaskDefinitionInterface{}, errors.New("task definition query: cannot be nil")
	}
	if err := options.Validate(); err != nil {
		return []TaskDefinitionInterface{}, err
	}
	q := store.buildTaskDefinitionQuery(options)
	var tasks []taskDefinition
	if err := q.Table(store.taskDefinitionTableName).Get(&tasks); err != nil {
		return []TaskDefinitionInterface{}, err
	}
	list := make([]TaskDefinitionInterface, len(tasks))
	for i, t := range tasks {
		task := t
		list[i] = &task
	}
	return list, nil
}

func (store *Store) TaskDefinitionSoftDelete(ctx context.Context, task TaskDefinitionInterface) error {
	if task == nil {
		return errors.New("task is nil")
	}
	task.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	return store.TaskDefinitionUpdate(ctx, task)
}

func (store *Store) TaskDefinitionSoftDeleteByID(ctx context.Context, id string) error {
	task, err := store.TaskDefinitionFindByID(ctx, id)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found")
	}
	return store.TaskDefinitionSoftDelete(ctx, task)
}

func (store *Store) TaskDefinitionUpdate(ctx context.Context, task TaskDefinitionInterface) error {
	if task == nil {
		return errors.New("task is nil")
	}
	task.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	row := map[string]any{
		COLUMN_STATUS:          task.GetStatus(),
		COLUMN_ALIAS:           task.GetAlias(),
		COLUMN_TITLE:           task.GetTitle(),
		COLUMN_DESCRIPTION:     task.GetDescription(),
		COLUMN_MEMO:            task.GetMemo(),
		COLUMN_IS_RECURRING:    task.GetIsRecurring(),
		COLUMN_RECURRENCE_RULE: task.GetRecurrenceRule(),
		COLUMN_UPDATED_AT:      task.GetUpdatedAt(),
		COLUMN_SOFT_DELETED_AT: task.GetSoftDeletedAt(),
	}

	_, err := store.db.Query().
		Table(store.taskDefinitionTableName).
		Where(COLUMN_ID+" = ?", task.GetID()).
		Update(row)
	return err
}

// TaskDefinitionEnqueueByAlias finds a task by its alias and appends it to the queue
func (store *Store) TaskDefinitionEnqueueByAlias(
	ctx context.Context,
	queueName string,
	taskAlias string,
	parameters map[string]any,
) (TaskQueueInterface, error) {
	task, err := store.TaskDefinitionFindByAlias(ctx, taskAlias)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("task with alias '" + taskAlias + "' not found")
	}

	parameters = queuePrependTaskAliasToParameters(taskAlias, parameters)

	parametersBytes, jsonErr := json.Marshal(parameters)
	if jsonErr != nil {
		return nil, errors.New("parameters json marshal error")
	}

	parametersStr := string(parametersBytes)

	queuedTask := NewTaskQueue().
		SetQueueName(queueName).
		SetTaskID(task.GetID()).
		SetAttempts(0).
		SetParameters(parametersStr).
		SetStatus(TaskQueueStatusQueued)

	err = store.TaskQueueCreate(ctx, queuedTask)
	if err != nil {
		return queuedTask, err
	}

	return queuedTask, err
}

// queuePrependTaskAliasToParameters prepends a task alias to the queue parameters so that its easy to distinguish
func queuePrependTaskAliasToParameters(alias string, parameters map[string]interface{}) map[string]interface{} {
	copiedParameters := map[string]interface{}{
		"task_alias": alias,
	}
	for index, element := range parameters {
		copiedParameters[index] = element
	}
	return copiedParameters
}

func (store *Store) buildTaskDefinitionQuery(options TaskDefinitionQueryInterface) contractsorm.Query {
	// Use Model() to enable neat's automatic soft delete handling via SoftDeletesMaxDate
	q := store.db.Query().Model(&taskDefinition{})

	if options == nil {
		return q
	}

	if options.HasAlias() && options.Alias() != "" {
		q = q.Where(COLUMN_ALIAS+" = ?", options.Alias())
	}

	if options.HasCreatedAtGte() && options.CreatedAtGte() != "" {
		q = q.Where(COLUMN_CREATED_AT+" >= ?", options.CreatedAtGte())
		q = q.Where(COLUMN_CREATED_AT+" >= ?", options.CreatedAtGte())
	}

	if options.HasCreatedAtLte() && options.CreatedAtLte() != "" {
		q = q.Where(COLUMN_CREATED_AT+" <= ?", options.CreatedAtLte())
		q = q.Where(COLUMN_CREATED_AT+" <= ?", options.CreatedAtLte())
	}

	if options.HasID() && options.ID() != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID())
	}

	if options.HasIDIn() && len(options.IDIn()) > 0 {
		args := make([]any, len(options.IDIn()))
		for i, id := range options.IDIn() {
			args[i] = id
		}
		q = q.WhereIn(COLUMN_ID, args)
	}

	if options.HasStatus() && options.Status() != "" {
		q = q.Where(COLUMN_STATUS+" = ?", options.Status())
	}

	if options.HasStatusIn() && len(options.StatusIn()) > 0 {
		args := make([]any, len(options.StatusIn()))
		for i, status := range options.StatusIn() {
			args[i] = status
		}
		q = q.WhereIn(COLUMN_STATUS, args)
	}

	if options.HasLimit() && options.Limit() > 0 {
		q = q.Limit(options.Limit())
	}

	if options.HasOffset() && options.Offset() > 0 {
		q = q.Offset(options.Offset())
	}

	if options.HasOrderBy() && options.OrderBy() != "" {
		sortOrder := DESC
		if options.HasSortOrder() && options.SortOrder() != "" {
			sortOrder = options.SortOrder()
		}
		q = q.OrderBy(options.OrderBy(), sortOrder)
	}

	// Handle soft delete filtering via neat's automatic handling (SoftDeletesMaxDate)
	if options.HasSoftDeletedIncluded() && options.SoftDeletedIncluded() {
		q = q.WithSoftDeleted()
	}

	return q
}
