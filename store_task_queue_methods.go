package taskstore

import (
	"context"
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

func (store *Store) TaskQueueCount(ctx context.Context, options TaskQueueQueryInterface) (int64, error) {
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

// TaskQueueCreate creates a queued task
func (store *Store) TaskQueueCreate(ctx context.Context, queue TaskQueueInterface) error {
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

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

	if err != nil {
		return err
	}

	queue.MarkAsNotDirty()

	return nil
}

func (store *Store) TaskQueueDelete(ctx context.Context, queue TaskQueueInterface) error {
	if queue == nil {
		return errors.New("queue is nil")
	}

	return store.TaskQueueDeleteByID(ctx, queue.ID())
}

func (st *Store) TaskQueueDeleteByID(ctx context.Context, id string) error {
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

	_, err = st.db.ExecContext(ctx, sqlStr, preparedArgs...)

	return err
}

// TaskQueueFail fails a queued task
func (st *Store) TaskQueueFail(ctx context.Context, queue TaskQueueInterface) error {
	queue.SetCompletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	queue.SetStatus(TaskQueueStatusFailed)
	return st.TaskQueueUpdate(ctx, queue)
}

// TaskQueueFindByID finds a Queue by ID
func (store *Store) TaskQueueFindByID(ctx context.Context, id string) (TaskQueueInterface, error) {
	if id == "" {
		return nil, errors.New("queue id is empty")
	}

	query := TaskQueueQuery().SetID(id).SetLimit(1)

	list, err := store.TaskQueueList(ctx, query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *Store) TaskQueueFindRunning(ctx context.Context, limit int) []TaskQueueInterface {
	return store.TaskQueueFindRunningByQueue(ctx, DefaultQueueName, limit)
}

func (store *Store) TaskQueueFindNextQueuedTask(ctx context.Context) (TaskQueueInterface, error) {
	return store.TaskQueueFindNextQueuedTaskByQueue(ctx, DefaultQueueName)
}

func (store *Store) TaskQueueList(ctx context.Context, query TaskQueueQueryInterface) ([]TaskQueueInterface, error) {
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

	modelMaps, err := db.SelectToMapString(ctx, sqlStr)

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

func (store *Store) TaskQueueProcessNext(ctx context.Context) error {
	return store.TaskQueueProcessNextByQueue(ctx, DefaultQueueName)
}

func normalizeQueueName(queueName string) string {
	if queueName == "" {
		return DefaultQueueName
	}
	return queueName
}

func (store *Store) TaskQueueFindRunningByQueue(ctx context.Context, queueName string, limit int) []TaskQueueInterface {
	queueName = normalizeQueueName(queueName)

	runningTasks, errList := store.TaskQueueList(ctx, TaskQueueQuery().
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

func (store *Store) TaskQueueFindNextQueuedTaskByQueue(ctx context.Context, queueName string) (TaskQueueInterface, error) {
	queueName = normalizeQueueName(queueName)

	queuedTasks, errList := store.TaskQueueList(ctx, TaskQueueQuery().SetStatus(TaskQueueStatusQueued).
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

// TaskQueueClaimNext atomically claims the next queued task for processing.
// It uses SELECT FOR UPDATE within a transaction to prevent race conditions
// where multiple workers might try to process the same task.
//
// Returns:
//   - TaskQueueInterface: The claimed task (status updated to "running")
//   - error: Any error that occurred during the operation
//
// Returns (nil, nil) if no tasks are available to claim.
func (store *Store) TaskQueueClaimNext(ctx context.Context, queueName string) (TaskQueueInterface, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	queueName = normalizeQueueName(queueName)

	// Start a database transaction
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Will be a no-op if committed

	// SELECT FOR UPDATE query to lock the row
	// Note: This works across SQLite (3.35+), MySQL, and PostgreSQL
	// SELECT FOR UPDATE query to lock the row
	// Note: This works across SQLite (3.35+), MySQL, and PostgreSQL
	selectSQL := `
		SELECT *
		FROM ` + store.taskQueueTableName + `
		WHERE ` + COLUMN_STATUS + ` = ? 
		  AND ` + COLUMN_QUEUE_NAME + ` = ?
		ORDER BY ` + COLUMN_CREATED_AT + ` ASC
		LIMIT 1`

	params := []interface{}{TaskQueueStatusQueued, queueName}

	if store.dbDriverName != "sqlite" {
		// MySQL and PostgreSQL support FOR UPDATE
		// Note: SKIP LOCKED removed for MySQL 5.7 compatibility (only available in MySQL 8.0+)
		selectSQL += " FOR UPDATE"
	}

	if store.debugEnabled {
		log.Println("TaskQueueClaimNext SELECT:", selectSQL)
	}

	// Execute SELECT query
	rows, err := tx.QueryContext(ctx, selectSQL, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		// No tasks available - this is normal
		return nil, nil
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	taskData := make(map[string]string)
	for i, col := range columns {
		val := values[i]
		taskData[col] = cast.ToString(val)
	}

	id := taskData[COLUMN_ID]

	// Update status to "running" within the same transaction
	updateSQL := `
		UPDATE ` + store.taskQueueTableName + `
		SET ` + COLUMN_STATUS + ` = ?, 
		    ` + COLUMN_STARTED_AT + ` = ?,
		    ` + COLUMN_UPDATED_AT + ` = ?
		WHERE ` + COLUMN_ID + ` = ?`

	now := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	_, err = tx.ExecContext(ctx, updateSQL, TaskQueueStatusRunning, now, now, id)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Create task object from data
	task := NewTaskQueueFromExistingData(taskData)
	// Update the task object to reflect the new status
	task.SetStatus(TaskQueueStatusRunning)
	task.SetStartedAt(now)
	task.SetUpdatedAt(now)
	task.MarkAsNotDirty() // Since we just updated it in DB

	return task, nil
}

func (store *Store) TaskQueueProcessNextByQueue(ctx context.Context, queueName string) error {
	queueName = normalizeQueueName(queueName)

	// Atomically claim the next task
	// Note: Old implementation checked for running tasks which was too restrictive
	// The atomic claim handles concurrency properly
	nextQueuedTask, err := store.TaskQueueClaimNext(ctx, queueName)

	if err != nil {
		return err
	}

	if nextQueuedTask == nil {
		// No tasks available
		return nil
	}

	// Process the claimed task synchronously
	_, err = store.TaskQueueProcessTask(ctx, nextQueuedTask)

	return err
}

func (store *Store) TaskQueueProcessNextAsyncByQueue(ctx context.Context, queueName string) error {
	queueName = normalizeQueueName(queueName)

	// Atomically claim the next task (fixes race condition)
	nextQueuedTask, err := store.TaskQueueClaimNext(ctx, queueName)

	if err != nil {
		return err
	}

	if nextQueuedTask == nil {
		// No tasks available - this is normal
		return nil
	}

	// Spawn goroutine to process the claimed task
	go func(q TaskQueueInterface) {
		_, err := store.TaskQueueProcessTask(ctx, q)
		if err != nil && store.debugEnabled {
			log.Println("TaskQueueProcessTask error:", err)
		}
	}(nextQueuedTask)

	return nil
}

func (store *Store) TaskQueueSoftDelete(ctx context.Context, queue TaskQueueInterface) error {
	if queue == nil {
		return errors.New("queue is nil")
	}

	queue.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.TaskQueueUpdate(ctx, queue)
}

func (store *Store) TaskQueueSoftDeleteByID(ctx context.Context, id string) error {
	queue, err := store.TaskQueueFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.TaskQueueSoftDelete(ctx, queue)
}

// TaskQueueSuccess completes a queued task  successfully
func (st *Store) TaskQueueSuccess(ctx context.Context, queue TaskQueueInterface) error {
	queue.SetCompletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	queue.SetStatus(TaskQueueStatusSuccess)
	return st.TaskQueueUpdate(ctx, queue)
}

func (store *Store) QueuedTaskForceFail(ctx context.Context, queuedTask TaskQueueInterface, waitMinutes int) error {
	startedAt := queuedTask.StartedAt()

	// Skip tasks that haven't actually started yet
	// This includes empty strings and NULL_DATETIME values
	if startedAt == "" || startedAt == sb.NULL_DATETIME {
		return nil
	}

	minutes := -1 * waitMinutes

	waitTill := queuedTask.StartedAtCarbon().AddMinutes(minutes)

	isOvertime := carbon.Now(carbon.UTC).Gt(waitTill)

	if isOvertime {
		queuedTask.AppendDetails("Failed forcefully after " + cast.ToString(waitMinutes) + " minutes timeout")
		return store.TaskQueueFail(ctx, queuedTask)
	}

	return nil
}

// TaskQueueUpdate creates a Queue
func (store *Store) TaskQueueUpdate(ctx context.Context, queue TaskQueueInterface) error {
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

	_, err := store.db.ExecContext(ctx, sqlStr, params...)

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

	if options.HasQueueName() {
		q = q.Where(goqu.C(COLUMN_QUEUE_NAME).Eq(options.QueueName()))
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

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted), columns, nil
}
