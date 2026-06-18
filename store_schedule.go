package taskstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
)

// ScheduleCount returns the number of schedules that match the given query options.
func (store *Store) ScheduleCount(ctx context.Context, options ScheduleQueryInterface) (int64, error) {
	if options == nil {
		return 0, errors.New("schedule query: cannot be nil")
	}
	q := store.buildScheduleQuery(options)
	var count int64
	err := q.Table(store.scheduleTableName).Count(&count)
	return count, err
}

// ScheduleCreate creates a new schedule record in the store.
func (store *Store) ScheduleCreate(ctx context.Context, schedule ScheduleInterface) error {
	if schedule == nil {
		return errors.New("schedule is nil")
	}
	if schedule.GetCreatedAt() == "" {
		schedule.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}
	if schedule.GetUpdatedAt() == "" {
		schedule.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}
	if schedule.GetSoftDeletedAt() == "" {
		schedule.SetSoftDeletedAt(MAX_DATETIME)
	}

	rrBytes, err := json.Marshal(schedule.GetRecurrenceRule())
	if err != nil {
		return err
	}

	tpBytes, err := json.Marshal(schedule.GetTaskParameters())
	if err != nil {
		return err
	}

	row := map[string]any{
		COLUMN_ID:                  schedule.GetID(),
		COLUMN_NAME:                schedule.GetName(),
		COLUMN_DESCRIPTION:         schedule.GetDescription(),
		COLUMN_STATUS:              schedule.GetStatus(),
		COLUMN_RECURRENCE_RULE:     string(rrBytes),
		COLUMN_QUEUE_NAME:          schedule.GetQueueName(),
		COLUMN_TASK_DEFINITION_ID:  schedule.GetTaskDefinitionID(),
		COLUMN_PARAMETERS:          string(tpBytes),
		COLUMN_START_AT:            schedule.GetStartAt(),
		COLUMN_END_AT:              schedule.GetEndAt(),
		COLUMN_EXECUTION_COUNT:     schedule.GetExecutionCount(),
		COLUMN_MAX_EXECUTION_COUNT: schedule.GetMaxExecutionCount(),
		COLUMN_LAST_RUN_AT:         schedule.GetLastRunAt(),
		COLUMN_NEXT_RUN_AT:         schedule.GetNextRunAt(),
		COLUMN_CREATED_AT:          schedule.CreatedAtCarbon().StdTime(),
		COLUMN_UPDATED_AT:          schedule.UpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT:     schedule.SoftDeletedAtCarbon().StdTime(),
	}

	return store.db.Query().Table(store.scheduleTableName).Create(row)
}

// ScheduleDelete deletes the given schedule from the store.
func (store *Store) ScheduleDelete(ctx context.Context, schedule ScheduleInterface) error {
	if schedule == nil {
		return errors.New("schedule is nil")
	}
	return store.ScheduleDeleteByID(ctx, schedule.GetID())
}

// ScheduleDeleteByID deletes the schedule with the given ID from the store.
func (store *Store) ScheduleDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("schedule id is empty")
	}
	_, err := store.db.Query().
		Table(store.scheduleTableName).
		Where(COLUMN_ID+" = ?", id).
		Delete()
	return err
}

// ScheduleFindByID finds a schedule by its ID.
func (store *Store) ScheduleFindByID(ctx context.Context, id string) (ScheduleInterface, error) {
	if id == "" {
		return nil, errors.New("schedule id is empty")
	}
	q := store.db.Query().Table(store.scheduleTableName).
		Where(COLUMN_ID+" = ?", id)
	q = q.Where(COLUMN_SOFT_DELETED_AT+" = ?", carbon.Parse(MAX_DATETIME, carbon.UTC).StdTime())

	var schedule scheduleImplementation
	if err := q.First(&schedule); err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "no rows found" {
			return nil, nil
		}
		return nil, err
	}
	return &schedule, nil
}

// ScheduleList returns a list of schedules that match the given query options.
func (store *Store) ScheduleList(ctx context.Context, options ScheduleQueryInterface) ([]ScheduleInterface, error) {
	if options == nil {
		return []ScheduleInterface{}, errors.New("schedule query: cannot be nil")
	}
	q := store.buildScheduleQuery(options)
	var schedules []scheduleImplementation
	if err := q.Table(store.scheduleTableName).Get(&schedules); err != nil {
		return []ScheduleInterface{}, err
	}
	list := make([]ScheduleInterface, len(schedules))
	for i, s := range schedules {
		sched := s
		list[i] = &sched
	}
	return list, nil
}

// ScheduleSoftDelete marks the given schedule as soft-deleted.
func (store *Store) ScheduleSoftDelete(ctx context.Context, schedule ScheduleInterface) error {
	if schedule == nil {
		return errors.New("schedule is nil")
	}
	schedule.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	return store.ScheduleUpdate(ctx, schedule)
}

// ScheduleSoftDeleteByID marks the schedule with the given ID as soft-deleted.
func (store *Store) ScheduleSoftDeleteByID(ctx context.Context, id string) error {
	schedule, err := store.ScheduleFindByID(ctx, id)
	if err != nil {
		return err
	}
	if schedule == nil {
		return errors.New("schedule not found")
	}
	return store.ScheduleSoftDelete(ctx, schedule)
}

// ScheduleUpdate updates an existing schedule record in the store.
func (store *Store) ScheduleUpdate(ctx context.Context, schedule ScheduleInterface) error {
	if schedule == nil {
		return errors.New("schedule is nil")
	}
	schedule.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	rrBytes, err := json.Marshal(schedule.GetRecurrenceRule())
	if err != nil {
		return err
	}

	tpBytes, err := json.Marshal(schedule.GetTaskParameters())
	if err != nil {
		return err
	}

	row := map[string]any{
		COLUMN_NAME:                schedule.GetName(),
		COLUMN_DESCRIPTION:         schedule.GetDescription(),
		COLUMN_STATUS:              schedule.GetStatus(),
		COLUMN_RECURRENCE_RULE:     string(rrBytes),
		COLUMN_QUEUE_NAME:          schedule.GetQueueName(),
		COLUMN_TASK_DEFINITION_ID:  schedule.GetTaskDefinitionID(),
		COLUMN_PARAMETERS:          string(tpBytes),
		COLUMN_START_AT:            schedule.GetStartAt(),
		COLUMN_END_AT:              schedule.GetEndAt(),
		COLUMN_EXECUTION_COUNT:     schedule.GetExecutionCount(),
		COLUMN_MAX_EXECUTION_COUNT: schedule.GetMaxExecutionCount(),
		COLUMN_LAST_RUN_AT:         schedule.GetLastRunAt(),
		COLUMN_NEXT_RUN_AT:         schedule.GetNextRunAt(),
		COLUMN_UPDATED_AT:          schedule.UpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT:     schedule.SoftDeletedAtCarbon().StdTime(),
	}

	_, err = store.db.Query().
		Table(store.scheduleTableName).
		Where(COLUMN_ID+" = ?", schedule.GetID()).
		Update(row)
	return err
}

// ScheduleRun runs all active schedules that are due.
func (store *Store) ScheduleRun(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	query := NewScheduleQuery().
		SetStatus("active").
		SetLimit(100)

	schedules, err := store.ScheduleList(ctx, query)
	if err != nil {
		return err
	}

	for _, schedule := range schedules {
		if !schedule.IsDue() {
			continue
		}

		_, processErr := store.TaskDefinitionEnqueueByAlias(
			ctx,
			schedule.GetQueueName(),
			store.GetTaskDefinitionAliasByID(ctx, schedule.GetTaskDefinitionID()),
			schedule.GetTaskParameters(),
		)

		if processErr != nil {
			if store.debugEnabled {
				store.logger.Error("ScheduleRun: failed to enqueue task", "schedule_id", schedule.GetID(), "error", processErr)
			}
			continue
		}

		schedule.IncrementExecutionCount()
		schedule.UpdateLastRunAt()
		schedule.UpdateNextRunAt()

		if err := store.ScheduleUpdate(ctx, schedule); err != nil {
			if store.debugEnabled {
				store.logger.Error("ScheduleRun: failed to update schedule", "schedule_id", schedule.GetID(), "error", err)
			}
		}
	}

	return nil
}

// GetTaskDefinitionAliasByID finds a task definition by ID and returns its alias.
// Returns an empty string if not found.
func (store *Store) GetTaskDefinitionAliasByID(ctx context.Context, id string) string {
	task, err := store.TaskDefinitionFindByID(ctx, id)
	if err != nil || task == nil {
		return ""
	}
	return task.GetAlias()
}

func (store *Store) buildScheduleQuery(options ScheduleQueryInterface) contractsorm.Query {
	q := store.db.Query()

	if options == nil {
		return q
	}

	if options.ID() != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID())
	}

	if options.Name() != "" {
		q = q.Where(COLUMN_NAME+" = ?", options.Name())
	}

	if options.Status() != "" {
		q = q.Where(COLUMN_STATUS+" = ?", options.Status())
	}

	if options.QueueName() != "" {
		q = q.Where(COLUMN_QUEUE_NAME+" = ?", options.QueueName())
	}

	if options.TaskDefinitionID() != "" {
		q = q.Where(COLUMN_TASK_DEFINITION_ID+" = ?", options.TaskDefinitionID())
	}

	if options.Limit() > 0 {
		q = q.Limit(options.Limit())
	}

	if options.Offset() > 0 {
		q = q.Offset(options.Offset())
	}

	// Default: exclude soft-deleted records
	q = q.Where(COLUMN_SOFT_DELETED_AT+" = ?", carbon.Parse(MAX_DATETIME, carbon.UTC).StdTime())

	return q
}
