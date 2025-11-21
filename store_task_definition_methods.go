package taskstore

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

func (store *Store) TaskDefinitionCount(ctx context.Context, options TaskDefinitionQueryInterface) (int64, error) {
	options.SetCountOnly(true)

	q, _, err := store.taskDefinitionSelectQuery(options)

	if err != nil {
		return -1, err
	}

	sqlStr, params, errSql := q.Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)
	mapped, err := db.SelectToMapString(ctx, sqlStr, params...)
	if err != nil {
		return -1, err
	}

	if len(mapped) < 1 {
		return -1, nil
	}

	countStr := mapped[0]["count"]

	i, err := strconv.ParseInt(countStr, 10, 64)

	if err != nil {
		return -1, err

	}

	return i, nil
}

func (store *Store) TaskDefinitionCreate(ctx context.Context, task TaskDefinitionInterface) error {
	task.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	task.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	data := task.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.taskDefinitionTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	if store.db == nil {
		return errors.New("taskstore: database is nil")
	}

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

	if err != nil {
		return err
	}

	task.MarkAsNotDirty()

	return nil
}

func (store *Store) TaskDefinitionDelete(ctx context.Context, task TaskDefinitionInterface) error {
	if task == nil {
		return errors.New("task is nil")
	}

	return store.TaskDefinitionDeleteByID(ctx, task.ID())
}

func (store *Store) TaskDefinitionDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("task id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.taskDefinitionTableName).
		Prepared(true).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

	return err
}

func (store *Store) TaskDefinitionFindByAlias(ctx context.Context, alias string) (task TaskDefinitionInterface, err error) {
	if alias == "" {
		return nil, errors.New("task id is empty")
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

func (store *Store) TaskDefinitionFindByID(ctx context.Context, id string) (task TaskDefinitionInterface, err error) {
	if id == "" {
		return nil, errors.New("task id is empty")
	}

	query := TaskDefinitionQuery().SetID(id).SetLimit(1)

	list, err := store.TaskDefinitionList(ctx, query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *Store) TaskDefinitionList(ctx context.Context, query TaskDefinitionQueryInterface) ([]TaskDefinitionInterface, error) {
	q, columns, err := store.taskDefinitionSelectQuery(query)

	if err != nil {
		return []TaskDefinitionInterface{}, err
	}

	sqlStr, sqlParams, errSql := q.Prepared(true).Select(columns...).ToSQL()

	if errSql != nil {
		return []TaskDefinitionInterface{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	if store.db == nil {
		return []TaskDefinitionInterface{}, errors.New("taskstore: database is nil")
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)

	if db == nil {
		return []TaskDefinitionInterface{}, errors.New("taskstore: database is nil")
	}

	modelMaps, err := db.SelectToMapString(ctx, sqlStr, sqlParams...)

	if err != nil {
		return []TaskDefinitionInterface{}, err
	}

	list := []TaskDefinitionInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewTaskDefinitionFromExistingData(modelMap)
		list = append(list, model)
	})

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

	return store.TaskDefinitionSoftDelete(ctx, task)
}

func (store *Store) TaskDefinitionUpdate(ctx context.Context, task TaskDefinitionInterface) error {
	if task == nil {
		return errors.New("task is nil")
	}

	task.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := task.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.taskDefinitionTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(task.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	if store.db == nil {
		return errors.New("taskstore: database is nil")
	}

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

	task.MarkAsNotDirty()

	return err
}

// TaskEnqueueByAlias finds a task by its alias and appends it to the queue
func (st *Store) TaskEnqueueByAlias(ctx context.Context, taskAlias string, parameters map[string]interface{}) (TaskQueueInterface, error) {
	task, err := st.TaskDefinitionFindByAlias(ctx, taskAlias)

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
		SetTaskID(task.ID()).
		SetAttempts(0).
		SetParameters(parametersStr).
		SetStatus(TaskQueueStatusQueued)

	err = st.TaskQueueCreate(ctx, queuedTask)

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

func (store *Store) taskDefinitionSelectQuery(options TaskDefinitionQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, []any{}, errors.New("site options cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, []any{}, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.taskDefinitionTableName)

	if options.HasAlias() {
		q = q.Where(goqu.C(COLUMN_ALIAS).Eq(options.Alias()))
	}

	if options.HasCreatedAtGte() && options.HasCreatedAtLte() {
		q = q.Where(
			goqu.C(COLUMN_CREATED_AT).Gte(options.CreatedAtGte()),
			goqu.C(COLUMN_CREATED_AT).Lte(options.CreatedAtLte()),
		)
	} else if options.HasCreatedAtGte() {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Gte(options.CreatedAtGte()))
	} else if options.HasCreatedAtLte() {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Lte(options.CreatedAtLte()))
	}

	if options.HasID() {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	if options.HasIDIn() {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn()))
	}

	if options.HasStatus() {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status()))
	}

	if options.HasStatusIn() {
		q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn()))
	}

	if !options.IsCountOnly() {
		if options.HasLimit() {
			q = q.Limit(uint(options.Limit()))
		}

		if options.HasOffset() {
			q = q.Offset(uint(options.Offset()))
		}
	}

	sortOrder := sb.DESC
	if options.HasSortOrder() {
		sortOrder = options.SortOrder()
	}

	if options.HasOrderBy() {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.Order(goqu.I(options.OrderBy()).Asc())
		} else {
			q = q.Order(goqu.I(options.OrderBy()).Desc())
		}
	}

	columns = []any{}

	for _, column := range options.Columns() {
		columns = append(columns, column)
	}

	if options.SoftDeletedIncluded() {
		return q, columns, nil // soft deleted records requested specifically
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted), columns, nil
}
