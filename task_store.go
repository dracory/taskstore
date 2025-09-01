package taskstore

import (
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

func (store *Store) TaskCount(options TaskQueryInterface) (int64, error) {
	options.SetCountOnly(true)

	q, _, err := store.taskSelectQuery(options)

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
	mapped, err := db.SelectToMapString(sqlStr, params...)
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

func (store *Store) TaskCreate(task TaskInterface) error {
	task.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	task.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	data := task.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.taskTableName).
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

	_, err := store.db.Exec(sqlStr, params...)

	if err != nil {
		return err
	}

	task.MarkAsNotDirty()

	return nil
}

func (store *Store) TaskDelete(task TaskInterface) error {
	if task == nil {
		return errors.New("task is nil")
	}

	return store.TaskDeleteByID(task.ID())
}

func (store *Store) TaskDeleteByID(id string) error {
	if id == "" {
		return errors.New("task id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.taskTableName).
		Prepared(true).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.Exec(sqlStr, params...)

	return err
}

func (store *Store) TaskFindByAlias(alias string) (task TaskInterface, err error) {
	if alias == "" {
		return nil, errors.New("task id is empty")
	}

	query := TaskQuery().SetAlias(alias).SetLimit(1)

	list, err := store.TaskList(query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *Store) TaskFindByID(id string) (task TaskInterface, err error) {
	if id == "" {
		return nil, errors.New("task id is empty")
	}

	query := TaskQuery().SetID(id).SetLimit(1)

	list, err := store.TaskList(query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *Store) TaskList(query TaskQueryInterface) ([]TaskInterface, error) {
	q, columns, err := store.taskSelectQuery(query)

	if err != nil {
		return []TaskInterface{}, err
	}

	sqlStr, sqlParams, errSql := q.Prepared(true).Select(columns...).ToSQL()

	if errSql != nil {
		return []TaskInterface{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	if store.db == nil {
		return []TaskInterface{}, errors.New("taskstore: database is nil")
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)

	if db == nil {
		return []TaskInterface{}, errors.New("taskstore: database is nil")
	}

	modelMaps, err := db.SelectToMapString(sqlStr, sqlParams...)

	if err != nil {
		return []TaskInterface{}, err
	}

	list := []TaskInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewTaskFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *Store) TaskSoftDelete(task TaskInterface) error {
	if task == nil {
		return errors.New("task is nil")
	}

	task.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.TaskUpdate(task)
}

func (store *Store) TaskSoftDeleteByID(id string) error {
	task, err := store.TaskFindByID(id)

	if err != nil {
		return err
	}

	return store.TaskSoftDelete(task)
}

func (store *Store) TaskUpdate(task TaskInterface) error {
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
		Update(store.taskTableName).
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

	_, err := store.db.Exec(sqlStr, params...)

	task.MarkAsNotDirty()

	return err
}

// TaskEnqueueByAlias finds a task by its alias and appends it to the queue
func (st *Store) TaskEnqueueByAlias(taskAlias string, parameters map[string]interface{}) (QueueInterface, error) {
	task, err := st.TaskFindByAlias(taskAlias)

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

	queuedTask := NewQueue().
		SetTaskID(task.ID()).
		SetAttempts(1).
		SetParameters(parametersStr).
		SetStatus(QueueStatusQueued)

	err = st.QueueCreate(queuedTask)

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

func (store *Store) taskSelectQuery(options TaskQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, []any{}, errors.New("site options cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, []any{}, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.taskTableName)

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
		return q, columns, nil // soft deleted sites requested specifically
	}

	softDeleted := goqu.C(COLUMN_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted), columns, nil
}
