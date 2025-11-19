package taskstore

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func (store *Store) TaskQueueCount(options TaskQueueQueryInterface) (int64, error) {
	options.SetCountOnly(true)

	q, _, err := store.taskQueueSelectQuery(options)

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

// TaskQueueCreate creates a queued task
func (store *Store) TaskQueueCreate(queue TaskQueueInterface) error {
	if queue.ID() == "" {
		time.Sleep(1 * time.Millisecond) // !!! important
		queue.SetID(uid.MicroUid())
	}
	if queue.CreatedAt() == "" {
		queue.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}
	if queue.UpdatedAt() == "" {
		queue.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	data := queue.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.taskQueueTableName).
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

	queue.MarkAsNotDirty()

	return nil
}

func (store *Store) TaskQueueDelete(queue TaskQueueInterface) error {
	if queue == nil {
		return errors.New("queue is nil")
	}

	return store.TaskQueueDeleteByID(queue.ID())
}

func (st *Store) TaskQueueDeleteByID(id string) error {
	if id == "" {
		return errors.New("queue id is empty")
	}

	sqlStr, preparedArgs, err := goqu.Dialect(st.dbDriverName).
		From(st.taskQueueTableName).
		Prepared(true).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		Delete().
		ToSQL()

	if err != nil {
		return err
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	_, err = st.db.Exec(sqlStr, preparedArgs...)

	return err
}

// TaskQueueFail fails a queued task
func (st *Store) TaskQueueFail(queue TaskQueueInterface) error {
	queue.SetCompletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	queue.SetStatus(TaskQueueStatusFailed)
	return st.TaskQueueUpdate(queue)
}

// TaskQueueFindByID finds a Queue by ID
func (store *Store) TaskQueueFindByID(id string) (TaskQueueInterface, error) {
	if id == "" {
		return nil, errors.New("queue id is empty")
	}

	query := TaskQueueQuery().SetID(id).SetLimit(1)

	list, err := store.TaskQueueList(query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *Store) TaskQueueFindRunning(limit int) []TaskQueueInterface {
	return store.TaskQueueFindRunningByQueue("", limit)
}

func (store *Store) TaskQueueFindNextQueuedTask() (TaskQueueInterface, error) {
	return store.TaskQueueFindNextQueuedTaskByQueue("")
}

func (store *Store) TaskQueueList(query TaskQueueQueryInterface) ([]TaskQueueInterface, error) {
	q, columns, err := store.taskQueueSelectQuery(query)

	if err != nil {
		return []TaskQueueInterface{}, err
	}

	sqlStr, _, errSql := q.Select(columns...).ToSQL()

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	if errSql != nil {
		return []TaskQueueInterface{}, errSql
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)

	if db == nil {
		return []TaskQueueInterface{}, errors.New("queuestore: database is nil")
	}

	modelMaps, err := db.SelectToMapString(sqlStr)

	if err != nil {
		return []TaskQueueInterface{}, err
	}

	list := []TaskQueueInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewTaskQueueFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *Store) TaskQueueProcessNext() error {
	return store.TaskQueueProcessNextByQueue("")
}

func normalizeQueueName(queueName string) string {
	if queueName == "" {
		return DefaultQueueName
	}
	return queueName
}

func (store *Store) TaskQueueFindRunningByQueue(queueName string, limit int) []TaskQueueInterface {
	queueName = normalizeQueueName(queueName)

	runningTasks, errList := store.TaskQueueList(TaskQueueQuery().
		SetStatus(TaskQueueStatusRunning).
		SetQueueName(queueName).
		SetLimit(limit).
		SetOrderBy(COLUMN_CREATED_AT).
		SetSortOrder(ASC))

	if errList != nil {
		return nil
	}

	return runningTasks
}

func (store *Store) TaskQueueFindNextQueuedTaskByQueue(queueName string) (TaskQueueInterface, error) {
	queueName = normalizeQueueName(queueName)

	queuedTasks, errList := store.TaskQueueList(TaskQueueQuery().SetStatus(TaskQueueStatusQueued).
		SetQueueName(queueName).
		SetLimit(1).
		SetOrderBy(COLUMN_CREATED_AT).
		SetSortOrder(ASC))

	if errList != nil {
		return nil, errList
	}

	if len(queuedTasks) < 1 {
		return nil, nil
	}

	return queuedTasks[0], nil
}

func (store *Store) TaskQueueProcessNextByQueue(queueName string) error {
	queueName = normalizeQueueName(queueName)

	runningTasks := store.TaskQueueFindRunningByQueue(queueName, 1)

	if len(runningTasks) > 0 {
		log.Println("There is already a running task " + runningTasks[0].ID() + " (#" + runningTasks[0].ID() + "). Queue stopped while completed'")
		return nil
	}

	nextQueuedTask, err := store.TaskQueueFindNextQueuedTaskByQueue(queueName)

	if err != nil {
		return err
	}

	if nextQueuedTask == nil {
		// DEBUG log.Println("No queued tasks")
		return nil
	}

	_, err = store.QueuedTaskProcess(nextQueuedTask)

	return err
}

func (store *Store) TaskQueueProcessNextAsyncByQueue(queueName string) error {
	queueName = normalizeQueueName(queueName)

	nextQueuedTask, err := store.TaskQueueFindNextQueuedTaskByQueue(queueName)

	if err != nil {
		return err
	}

	if nextQueuedTask == nil {
		// DEBUG log.Println("No queued tasks")
		return nil
	}

	go func(q TaskQueueInterface) {
		_, err := store.QueuedTaskProcess(q)
		if err != nil && store.debugEnabled {
			log.Println("QueuedTaskProcess error:", err)
		}
	}(nextQueuedTask)

	return nil
}

func (store *Store) TaskQueueSoftDelete(queue TaskQueueInterface) error {
	if queue == nil {
		return errors.New("queue is nil")
	}

	queue.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.TaskQueueUpdate(queue)
}

func (store *Store) TaskQueueSoftDeleteByID(id string) error {
	queue, err := store.TaskQueueFindByID(id)

	if err != nil {
		return err
	}

	return store.TaskQueueSoftDelete(queue)
}

// TaskQueueSuccess completes a queued task  successfully
func (st *Store) TaskQueueSuccess(queue TaskQueueInterface) error {
	queue.SetCompletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	queue.SetStatus(TaskQueueStatusSuccess)
	return st.TaskQueueUpdate(queue)
}

func (store *Store) QueuedTaskForceFail(queuedTask TaskQueueInterface, waitMinutes int) error {
	startedAt := queuedTask.StartedAt()

	if startedAt == "" {
		return nil
	}

	minutes := -1 * waitMinutes

	waitTill := queuedTask.StartedAtCarbon().AddMinutes(minutes)

	isOvertime := carbon.Now(carbon.UTC).Gt(waitTill)

	if isOvertime {
		queuedTask.AppendDetails("Failed forcefully after " + cast.ToString(waitMinutes) + " minutes timeout")
		return store.TaskQueueFail(queuedTask)
	}

	return nil
}

// TaskQueueUpdate creates a Queue
func (store *Store) TaskQueueUpdate(queue TaskQueueInterface) error {
	queue.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	dataChanged := queue.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.taskQueueTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(queue.ID())).
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

	queue.MarkAsNotDirty()

	return err
}

func (store *Store) taskQueueSelectQuery(options TaskQueueQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, []any{}, errors.New("site options cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, []any{}, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.taskQueueTableName)

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

	if options.HasTaskID() {
		q = q.Where(goqu.C(COLUMN_TASK_ID).Eq(options.TaskID()))
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
