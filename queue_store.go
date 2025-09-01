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

func (store *Store) QueueCount(options QueueQueryInterface) (int64, error) {
	options.SetCountOnly(true)

	q, _, err := store.queueSelectQuery(options)

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

// QueueCreate creates a queued task
func (store *Store) QueueCreate(queue QueueInterface) error {
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
		Insert(store.queueTableName).
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

func (store *Store) QueueDelete(queue QueueInterface) error {
	if queue == nil {
		return errors.New("queue is nil")
	}

	return store.QueueDeleteByID(queue.ID())
}

func (st *Store) QueueDeleteByID(id string) error {
	if id == "" {
		return errors.New("queue id is empty")
	}

	sqlStr, preparedArgs, err := goqu.Dialect(st.dbDriverName).
		From(st.queueTableName).
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

// QueueFail fails a queued task
func (st *Store) QueueFail(queue QueueInterface) error {
	queue.SetCompletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	queue.SetStatus(QueueStatusFailed)
	return st.QueueUpdate(queue)
}

// QueueFindByID finds a Queue by ID
func (store *Store) QueueFindByID(id string) (QueueInterface, error) {
	if id == "" {
		return nil, errors.New("queue id is empty")
	}

	query := QueueQuery().SetID(id).SetLimit(1)

	list, err := store.QueueList(query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *Store) QueueFindRunning(limit int) []QueueInterface {

	runningTasks, errList := store.QueueList(QueueQuery().
		SetStatus(QueueStatusRunning).
		SetLimit(limit).
		SetOrderBy(COLUMN_CREATED_AT).
		SetSortOrder(ASC))

	if errList != nil {
		return nil
	}

	return runningTasks
}

func (store *Store) QueueFindNextQueuedTask() (QueueInterface, error) {
	queuedTasks, errList := store.QueueList(QueueQuery().SetStatus(QueueStatusQueued).
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

func (store *Store) QueueList(query QueueQueryInterface) ([]QueueInterface, error) {
	q, columns, err := store.queueSelectQuery(query)

	if err != nil {
		return []QueueInterface{}, err
	}

	sqlStr, _, errSql := q.Select(columns...).ToSQL()

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	if errSql != nil {
		return []QueueInterface{}, errSql
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)

	if db == nil {
		return []QueueInterface{}, errors.New("queuestore: database is nil")
	}

	modelMaps, err := db.SelectToMapString(sqlStr)

	if err != nil {
		return []QueueInterface{}, err
	}

	list := []QueueInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewQueueFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *Store) QueueProcessNext() error {
	runningTasks := store.QueueFindRunning(1)

	if len(runningTasks) > 0 {
		log.Println("There is already a running task " + runningTasks[0].ID() + " (#" + runningTasks[0].ID() + "). Queue stopped while completed'")
		return nil
	}

	nextQueuedTask, err := store.QueueFindNextQueuedTask()

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

func (store *Store) QueueSoftDelete(queue QueueInterface) error {
	if queue == nil {
		return errors.New("queue is nil")
	}

	queue.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.QueueUpdate(queue)
}

func (store *Store) QueueSoftDeleteByID(id string) error {
	queue, err := store.QueueFindByID(id)

	if err != nil {
		return err
	}

	return store.QueueSoftDelete(queue)
}

// QueueSuccess completes a queued task  successfully
func (st *Store) QueueSuccess(queue QueueInterface) error {
	queue.SetCompletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	queue.SetStatus(QueueStatusSuccess)
	return st.QueueUpdate(queue)
}

func (store *Store) QueuedTaskForceFail(queuedTask QueueInterface, waitMinutes int) error {
	startedAt := queuedTask.StartedAt()

	if startedAt == "" {
		return nil
	}

	minutes := -1 * waitMinutes

	waitTill := queuedTask.StartedAtCarbon().AddMinutes(minutes)

	isOvertime := carbon.Now(carbon.UTC).Gt(waitTill)

	if isOvertime {
		queuedTask.AppendDetails("Failed forcefully after " + cast.ToString(waitMinutes) + " minutes timeout")
		return store.QueueFail(queuedTask)
	}

	return nil
}

// QueueUpdate creates a Queue
func (store *Store) QueueUpdate(queue QueueInterface) error {
	queue.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	dataChanged := queue.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.queueTableName).
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

func (store *Store) queueSelectQuery(options QueueQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, []any{}, errors.New("site options cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, []any{}, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.queueTableName)

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
