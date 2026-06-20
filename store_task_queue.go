package taskstore

import (
	"context"
	"errors"
	"time"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

func (store *Store) TaskQueueCount(ctx context.Context, options TaskQueueQueryInterface) (int64, error) {
	if options == nil {
		return 0, errors.New("task queue query: cannot be nil")
	}
	if err := options.Validate(); err != nil {
		return 0, err
	}
	q := store.buildTaskQueueQuery(options)
	var count int64
	err := q.Table(store.taskQueueTableName).Count(&count)
	return count, err
}

// TaskQueueCreate creates a queued task
func (store *Store) TaskQueueCreate(ctx context.Context, queue TaskQueueInterface) error {
	if queue == nil {
		return errors.New("taskstore: queue is nil")
	}
	if queue.GetID() == "" {
		time.Sleep(1 * time.Millisecond) // !!! important
		queue.SetID(neatuid.GenerateShortID())
	}
	if queue.GetCreatedAt() == "" {
		queue.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}
	if queue.GetUpdatedAt() == "" {
		queue.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}

	row := map[string]any{
		COLUMN_ID:              queue.GetID(),
		COLUMN_QUEUE_NAME:      queue.GetQueueName(),
		COLUMN_TASK_ID:         queue.GetTaskID(),
		COLUMN_PARAMETERS:      queue.GetParameters(),
		COLUMN_STATUS:          queue.GetStatus(),
		COLUMN_OUTPUT:          queue.GetOutput(),
		COLUMN_DETAILS:         queue.GetDetails(),
		COLUMN_ATTEMPTS:        queue.GetAttempts(),
		COLUMN_STARTED_AT:      queue.GetStartedAt(),
		COLUMN_COMPLETED_AT:    queue.GetCompletedAt(),
		COLUMN_CREATED_AT:      queue.GetCreatedAt(),
		COLUMN_UPDATED_AT:      queue.GetUpdatedAt(),
		COLUMN_SOFT_DELETED_AT: queue.GetSoftDeletedAt(),
	}

	return store.db.Query().Table(store.taskQueueTableName).Create(row)
}

func (store *Store) TaskQueueDelete(ctx context.Context, queue TaskQueueInterface) error {
	if queue == nil {
		return errors.New("queue is nil")
	}
	return store.TaskQueueDeleteByID(ctx, queue.GetID())
}

func (store *Store) TaskQueueDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("queue id is empty")
	}
	_, err := store.db.Query().
		Table(store.taskQueueTableName).
		Where(COLUMN_ID+" = ?", id).
		Delete()
	return err
}

// TaskQueueFail fails a queued task
func (store *Store) TaskQueueFail(ctx context.Context, queue TaskQueueInterface) error {
	queue.SetCompletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	queue.SetStatus(TaskQueueStatusFailed)
	return store.TaskQueueUpdate(ctx, queue)
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

func (store *Store) TaskQueueList(ctx context.Context, options TaskQueueQueryInterface) ([]TaskQueueInterface, error) {
	if options == nil {
		return []TaskQueueInterface{}, errors.New("task queue query: cannot be nil")
	}
	if err := options.Validate(); err != nil {
		return []TaskQueueInterface{}, err
	}
	q := store.buildTaskQueueQuery(options)
	var queues []taskQueue
	if err := q.Table(store.taskQueueTableName).Get(&queues); err != nil {
		return []TaskQueueInterface{}, err
	}
	list := make([]TaskQueueInterface, len(queues))
	for i, que := range queues {
		queue := que
		list[i] = &queue
	}
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

	tx, err := store.db.Query().Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	q := tx.Table(store.taskQueueTableName).
		Where(COLUMN_STATUS+" = ?", TaskQueueStatusQueued).
		Where(COLUMN_QUEUE_NAME+" = ?", queueName).
		OrderBy(COLUMN_CREATED_AT, ASC).
		Limit(1)
	if !store.isSQLite {
		q = q.LockForUpdate()
	}

	var tasks []taskQueue
	if err := q.Find(&tasks); err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, nil
	}
	task := tasks[0]

	now := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	_, err = tx.Table(store.taskQueueTableName).
		Where(COLUMN_ID+" = ?", task.ShortID.ID).
		Update(map[string]any{
			COLUMN_STATUS:     TaskQueueStatusRunning,
			COLUMN_STARTED_AT: now,
			COLUMN_UPDATED_AT: now,
		})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	task.StatusField = TaskQueueStatusRunning
	task.SetStartedAt(now)
	task.SetUpdatedAt(now)

	return &task, nil
}

func (store *Store) TaskQueueProcessNextByQueue(ctx context.Context, queueName string) error {
	queueName = normalizeQueueName(queueName)
	nextQueuedTask, err := store.TaskQueueClaimNext(ctx, queueName)
	if err != nil {
		return err
	}
	if nextQueuedTask == nil {
		return nil
	}
	_, err = store.TaskQueueProcessTask(ctx, nextQueuedTask)
	return err
}

func (store *Store) TaskQueueProcessNextAsyncByQueue(ctx context.Context, queueName string) error {
	queueName = normalizeQueueName(queueName)
	nextQueuedTask, err := store.TaskQueueClaimNext(ctx, queueName)
	if err != nil {
		return err
	}
	if nextQueuedTask == nil {
		return nil
	}
	go func(q TaskQueueInterface) {
		_, err := store.TaskQueueProcessTask(ctx, q)
		if err != nil && store.debugEnabled {
			store.logger.Error("TaskQueueProcessTask error", "error", err)
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
	if queue == nil {
		return errors.New("queue not found")
	}
	return store.TaskQueueSoftDelete(ctx, queue)
}

// TaskQueueSuccess completes a queued task successfully
func (store *Store) TaskQueueSuccess(ctx context.Context, queue TaskQueueInterface) error {
	queue.SetCompletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	queue.SetStatus(TaskQueueStatusSuccess)
	return store.TaskQueueUpdate(ctx, queue)
}

func (store *Store) QueuedTaskForceFail(ctx context.Context, queuedTask TaskQueueInterface, waitMinutes int) error {
	startedAt := queuedTask.GetStartedAt()
	if startedAt == "" || startedAt == NULL_DATETIME {
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

// TaskQueueUpdate updates a queued task
func (store *Store) TaskQueueUpdate(ctx context.Context, queue TaskQueueInterface) error {
	if queue == nil {
		return errors.New("queue is nil")
	}
	queue.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	row := map[string]any{
		COLUMN_QUEUE_NAME:      queue.GetQueueName(),
		COLUMN_TASK_ID:         queue.GetTaskID(),
		COLUMN_PARAMETERS:      queue.GetParameters(),
		COLUMN_STATUS:          queue.GetStatus(),
		COLUMN_OUTPUT:          queue.GetOutput(),
		COLUMN_DETAILS:         queue.GetDetails(),
		COLUMN_ATTEMPTS:        queue.GetAttempts(),
		COLUMN_STARTED_AT:      queue.GetStartedAt(),
		COLUMN_COMPLETED_AT:    queue.GetCompletedAt(),
		COLUMN_UPDATED_AT:      queue.GetUpdatedAt(),
		COLUMN_SOFT_DELETED_AT: queue.GetSoftDeletedAt(),
	}

	_, err := store.db.Query().
		Table(store.taskQueueTableName).
		Where(COLUMN_ID+" = ?", queue.GetID()).
		Update(row)
	return err
}

func (store *Store) buildTaskQueueQuery(options TaskQueueQueryInterface) contractsorm.Query {
	// Use Model() to enable neat's automatic soft delete handling via SoftDeletesMaxDate
	q := store.db.Query().Model(&taskQueue{})

	if options == nil {
		return q
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

	if options.HasTaskID() && options.TaskID() != "" {
		q = q.Where(COLUMN_TASK_ID+" = ?", options.TaskID())
	}

	if options.HasQueueName() && options.QueueName() != "" {
		q = q.Where(COLUMN_QUEUE_NAME+" = ?", options.QueueName())
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
