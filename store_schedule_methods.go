package taskstore

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/database"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

// ScheduleCount returns the number of schedules that match the given query options.
func (store *Store) ScheduleCount(ctx context.Context, options ScheduleQueryInterface) (int64, error) {
	options.SetLimit(1)
	q, _, err := store.scheduleSelectQuery(options)

	if err != nil {
		return -1, err
	}

	sqlStr, params, errSql := q.Prepared(true).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	queryable := database.NewQueryableContext(ctx, store.db)
	mapped, err := database.SelectToMapString(queryable, sqlStr, params...)
	if err != nil {
		return -1, err
	}

	if len(mapped) < 1 {
		return -1, nil
	}

	// Parse the count from the result
	countStr := mapped[0]["count"]
	i, err := strconv.ParseInt(countStr, 10, 64)

	if err != nil {
		return -1, err
	}

	return i, nil
}

// ScheduleCreate creates a new schedule record in the store.
func (store *Store) ScheduleCreate(ctx context.Context, schedule ScheduleInterface) error {
	// Set the created and updated timestamps
	schedule.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	schedule.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	// Prepare the data to be inserted
	data := map[string]interface{}{
		COLUMN_ID:                  schedule.GetID(),
		COLUMN_NAME:                schedule.GetName(),
		COLUMN_DESCRIPTION:         schedule.GetDescription(),
		COLUMN_STATUS:              schedule.GetStatus(),
		COLUMN_QUEUE_NAME:          schedule.GetQueueName(),
		COLUMN_TASK_DEFINITION_ID:  schedule.GetTaskDefinitionID(),
		COLUMN_START_AT:            schedule.GetStartAt(),
		COLUMN_END_AT:              schedule.GetEndAt(),
		COLUMN_EXECUTION_COUNT:     schedule.GetExecutionCount(),
		COLUMN_MAX_EXECUTION_COUNT: schedule.GetMaxExecutionCount(),
		COLUMN_LAST_RUN_AT:         schedule.GetLastRunAt(),
		COLUMN_NEXT_RUN_AT:         schedule.GetNextRunAt(),
		COLUMN_CREATED_AT:          schedule.GetCreatedAt(),
		COLUMN_UPDATED_AT:          schedule.GetUpdatedAt(),
		COLUMN_SOFT_DELETED_AT:     schedule.GetSoftDeletedAt(),
	}

	// Marshal the recurrence rule and task parameters
	rrBytes, err := json.Marshal(schedule.GetRecurrenceRule())
	if err != nil {
		return err
	}
	data[COLUMN_RECURRENCE_RULE] = string(rrBytes)

	// Marshal TaskParameters
	tpBytes, err := json.Marshal(schedule.GetTaskParameters())
	if err != nil {
		return err
	}
	data[COLUMN_PARAMETERS] = string(tpBytes)

	// Prepare the insert query
	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.scheduleTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	// Log the SQL query if debug is enabled
	if store.debugEnabled {
		log.Println(sqlStr)
	}

	// Execute the insert query
	if store.db == nil {
		return errors.New("taskstore: database is nil")
	}

	_, err = store.db.ExecContext(ctx, sqlStr, params...)

	return err
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

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.scheduleTableName).
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

// ScheduleFindByID finds a schedule by its ID.
func (store *Store) ScheduleFindByID(ctx context.Context, id string) (ScheduleInterface, error) {
	if id == "" {
		return nil, errors.New("schedule id is empty")
	}

	query := NewScheduleQuery().SetID(id).SetLimit(1)

	list, err := store.ScheduleList(ctx, query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// ScheduleList returns a list of schedules that match the given query options.
func (store *Store) ScheduleList(ctx context.Context, options ScheduleQueryInterface) ([]ScheduleInterface, error) {
	q, columns, err := store.scheduleSelectQuery(options)

	if err != nil {
		return []ScheduleInterface{}, err
	}

	sqlStr, sqlParams, errSql := q.Prepared(true).Select(columns...).ToSQL()

	if errSql != nil {
		return []ScheduleInterface{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	if store.db == nil {
		return []ScheduleInterface{}, errors.New("taskstore: database is nil")
	}

	queryable := database.NewQueryableContext(ctx, store.db)
	modelMaps, err := database.SelectToMapString(queryable, sqlStr, sqlParams...)

	if err != nil {
		return []ScheduleInterface{}, err
	}

	list := []ScheduleInterface{}

	for _, modelMap := range modelMaps {
		model := NewSchedule()
		model.SetID(modelMap[COLUMN_ID])
		model.SetName(modelMap[COLUMN_NAME])
		model.SetDescription(modelMap[COLUMN_DESCRIPTION])
		model.SetStatus(modelMap[COLUMN_STATUS])
		model.SetQueueName(modelMap[COLUMN_QUEUE_NAME])
		model.SetTaskDefinitionID(modelMap[COLUMN_TASK_DEFINITION_ID])
		model.SetStartAt(modelMap[COLUMN_START_AT])
		model.SetEndAt(modelMap[COLUMN_END_AT])
		model.SetExecutionCount(cast.ToInt(modelMap[COLUMN_EXECUTION_COUNT]))
		model.SetMaxExecutionCount(cast.ToInt(modelMap[COLUMN_MAX_EXECUTION_COUNT]))
		model.SetLastRunAt(modelMap[COLUMN_LAST_RUN_AT])
		model.SetNextRunAt(modelMap[COLUMN_NEXT_RUN_AT])
		model.SetCreatedAt(modelMap[COLUMN_CREATED_AT])
		model.SetUpdatedAt(modelMap[COLUMN_UPDATED_AT])
		model.SetSoftDeletedAt(modelMap[COLUMN_SOFT_DELETED_AT])

		// Unmarshal RecurrenceRule
		if rrStr, ok := modelMap[COLUMN_RECURRENCE_RULE]; ok && rrStr != "" {
			rr := NewRecurrenceRule()
			if err := json.Unmarshal([]byte(rrStr), rr); err == nil {
				model.SetRecurrenceRule(rr)
			}
		}

		// Unmarshal TaskParameters
		if tpStr, ok := modelMap[COLUMN_PARAMETERS]; ok && tpStr != "" {
			var tp map[string]any
			if err := json.Unmarshal([]byte(tpStr), &tp); err == nil {
				model.SetTaskParameters(tp)
			}
		}

		list = append(list, model)
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

	return store.ScheduleSoftDelete(ctx, schedule)
}

// ScheduleUpdate updates an existing schedule record in the store.
func (store *Store) ScheduleUpdate(ctx context.Context, schedule ScheduleInterface) error {
	if schedule == nil {
		return errors.New("schedule is nil")
	}

	schedule.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	data := map[string]interface{}{
		COLUMN_NAME:                schedule.GetName(),
		COLUMN_DESCRIPTION:         schedule.GetDescription(),
		COLUMN_STATUS:              schedule.GetStatus(),
		COLUMN_QUEUE_NAME:          schedule.GetQueueName(),
		COLUMN_TASK_DEFINITION_ID:  schedule.GetTaskDefinitionID(),
		COLUMN_START_AT:            schedule.GetStartAt(),
		COLUMN_END_AT:              schedule.GetEndAt(),
		COLUMN_EXECUTION_COUNT:     schedule.GetExecutionCount(),
		COLUMN_MAX_EXECUTION_COUNT: schedule.GetMaxExecutionCount(),
		COLUMN_LAST_RUN_AT:         schedule.GetLastRunAt(),
		COLUMN_NEXT_RUN_AT:         schedule.GetNextRunAt(),
		COLUMN_UPDATED_AT:          schedule.GetUpdatedAt(),
		COLUMN_SOFT_DELETED_AT:     schedule.GetSoftDeletedAt(),
	}

	// Marshal RecurrenceRule
	rrBytes, err := json.Marshal(schedule.GetRecurrenceRule())
	if err != nil {
		return err
	}
	data[COLUMN_RECURRENCE_RULE] = string(rrBytes)

	// Marshal TaskParameters
	tpBytes, err := json.Marshal(schedule.GetTaskParameters())
	if err != nil {
		return err
	}
	data[COLUMN_PARAMETERS] = string(tpBytes)

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.scheduleTableName).
		Prepared(true).
		Set(data).
		Where(goqu.C(COLUMN_ID).Eq(schedule.GetID())).
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

	_, err = store.db.ExecContext(ctx, sqlStr, params...)

	return err
}

// ScheduleRun scans for due schedules and enqueues their associated tasks.
func (store *Store) ScheduleRun(ctx context.Context) error {
	// Find active schedules that are due
	now := carbon.Now(carbon.UTC)

	// TODO: Implement query to find due schedules
	// For now, we'll just list all active schedules and check in memory (not efficient but simple for start)
	// In production, this should be a DB query

	query := NewScheduleQuery().SetStatus("active")
	schedules, err := store.ScheduleList(ctx, query)
	if err != nil {
		return err
	}

	if store.debugEnabled {
		log.Println("Found schedules:", len(schedules))
	}

	for _, schedule := range schedules {
		// Check if due
		nextRunAt := carbon.Parse(schedule.GetNextRunAt(), carbon.UTC)
		if store.debugEnabled {
			log.Println("Schedule:", schedule.GetID(), "NextRunAt:", nextRunAt, "Now:", now)
		}
		if nextRunAt.Lt(now) || nextRunAt.Eq(now) {
			// Enqueue task
			// Let's fetch TaskDefinition
			taskDef, err := store.TaskDefinitionFindByID(ctx, schedule.GetTaskDefinitionID())
			if err != nil {
				log.Println("Error finding task definition for schedule", schedule.GetID(), err)
				continue
			}
			if taskDef == nil {
				log.Println("Task definition not found for schedule", schedule.GetID())
				continue
			}

			_, err = store.TaskDefinitionEnqueueByAlias(ctx, schedule.GetQueueName(), taskDef.Alias(), schedule.GetTaskParameters())
			if err != nil {
				log.Println("Error enqueuing task for schedule", schedule.GetID(), err)
				continue
			}

			// Update schedule
			schedule.SetLastRunAt(now.ToDateTimeString(carbon.UTC))
			schedule.SetExecutionCount(schedule.GetExecutionCount() + 1)

			// Calculate next run
			nextRun, err := NextRunAt(schedule.GetRecurrenceRule(), now)
			if err != nil {
				log.Println("Error calculating next run for schedule", schedule.GetID(), err)
				// Disable schedule?
				continue
			}
			schedule.SetNextRunAt(nextRun.ToDateTimeString(carbon.UTC))

			// Check max execution count
			if schedule.GetMaxExecutionCount() > 0 && schedule.GetExecutionCount() >= schedule.GetMaxExecutionCount() {
				schedule.SetStatus("completed")
			}

			// Check end date
			endAt := carbon.Parse(schedule.GetEndAt(), carbon.UTC)
			if nextRun.Gt(endAt) {
				schedule.SetStatus("completed")
			}

			err = store.ScheduleUpdate(ctx, schedule)
			if err != nil {
				log.Println("Error updating schedule", schedule.GetID(), err)
			}
		}
	}

	return nil
}

func (store *Store) scheduleSelectQuery(options ScheduleQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, []any{}, errors.New("options cannot be nil")
	}

	q := goqu.Dialect(store.dbDriverName).From(store.scheduleTableName)

	if options.ID() != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	if options.Name() != "" {
		q = q.Where(goqu.C(COLUMN_NAME).Eq(options.Name()))
	}

	if options.Status() != "" {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status()))
	}

	if options.QueueName() != "" {
		q = q.Where(goqu.C(COLUMN_QUEUE_NAME).Eq(options.QueueName()))
	}

	if options.TaskDefinitionID() != "" {
		q = q.Where(goqu.C(COLUMN_TASK_DEFINITION_ID).Eq(options.TaskDefinitionID()))
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	q = q.Where(softDeleted)

	if options.Limit() > 0 {
		q = q.Limit(uint(options.Limit()))
	}

	if options.Offset() > 0 {
		q = q.Offset(uint(options.Offset()))
	}

	columns = []any{
		COLUMN_ID,
		COLUMN_NAME,
		COLUMN_DESCRIPTION,
		COLUMN_STATUS,
		COLUMN_RECURRENCE_RULE,
		COLUMN_QUEUE_NAME,
		COLUMN_TASK_DEFINITION_ID,
		COLUMN_PARAMETERS,
		COLUMN_START_AT,
		COLUMN_END_AT,
		COLUMN_EXECUTION_COUNT,
		COLUMN_MAX_EXECUTION_COUNT,
		COLUMN_LAST_RUN_AT,
		COLUMN_NEXT_RUN_AT,
		COLUMN_CREATED_AT,
		COLUMN_UPDATED_AT,
		COLUMN_SOFT_DELETED_AT,
	}

	return q, columns, nil
}
