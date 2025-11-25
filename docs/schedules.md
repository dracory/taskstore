# Schedules

Schedules allow you to automatically enqueue tasks at specified times and intervals using recurrence rules.

## Schedule Fields

- **id** - Unique identifier (generated using `uid.HumanUid()`)
- **name** - Human-readable name for the schedule
- **description** - Optional description of what the schedule does
- **status** - Current status of the schedule:
  - `draft` (default) - Schedule is not active
  - `active` - Schedule will automatically enqueue tasks
  - `inactive` - Schedule is paused
  - `completed` - Schedule has finished (max executions reached or end date passed)
- **recurrence_rule** - Defines when and how often the task runs (see RecurrenceRule below)
- **queue_name** - Which queue to enqueue tasks to
- **task_definition_id** - ID of the task definition to enqueue
- **task_parameters** - Parameters to pass to the enqueued task (map[string]any)
- **start_at** - When the schedule becomes active (default: NULL_DATETIME)
- **end_at** - When the schedule stops (default: MAX_DATETIME = 9999-12-31)
- **execution_count** - How many times this schedule has run
- **max_execution_count** - Maximum number of times to run (0 = unlimited)
- **last_run_at** - When the schedule last enqueued a task (default: NULL_DATETIME)
- **next_run_at** - When the schedule will next enqueue a task (default: NULL_DATETIME)
- **created_at** - When the schedule was created
- **updated_at** - When the schedule was last modified
- **soft_deleted_at** - Soft delete timestamp (default: MAX_DATETIME = not deleted)

## How Schedules Work

### Creating a Schedule

```go
schedule := taskstore.NewSchedule()
schedule.SetName("Daily Report")
schedule.SetDescription("Generate daily reports at 9 AM")
schedule.SetStatus("active")
schedule.SetQueueName("reports")
schedule.SetTaskDefinitionID(reportTaskDef.ID())
schedule.SetTaskParameters(map[string]any{
    "report_type": "daily",
})

// Set recurrence rule for daily at 9 AM
schedule.RecurrenceRule().SetFrequency(taskstore.FrequencyDaily)
schedule.RecurrenceRule().SetInterval(1)
schedule.RecurrenceRule().SetStartsAt("2025-01-01 09:00:00")

// Optional: limit executions
schedule.SetMaxExecutionCount(30) // Run for 30 days

err := store.ScheduleCreate(ctx, schedule)
```

### Processing Schedules

Call `ScheduleRun()` periodically (e.g., every minute) to process due schedules:

```go
err := store.ScheduleRun(ctx)
```

The `ScheduleRun` method:
1. Finds all schedules with `status = "active"` and `next_run_at <= NOW()`
2. For each due schedule:
   - Enqueues the task using `TaskDefinitionEnqueueByAlias`
   - Updates `last_run_at` to current time
   - Increments `execution_count`
   - Calculates and updates `next_run_at` using the recurrence rule
   - Marks as `completed` if:
     - `max_execution_count` is reached, OR
     - `next_run_at > end_at`

### Recurrence Rules

Recurrence rules define when tasks should be enqueued. Supported frequencies:

- **FrequencyNone** - One-time execution at `starts_at`
- **FrequencySecondly** - Every N seconds
- **FrequencyMinutely** - Every N minutes
- **FrequencyHourly** - Every N hours
- **FrequencyDaily** - Every N days
- **FrequencyWeekly** - Every N weeks
- **FrequencyMonthly** - Every N months
- **FrequencyYearly** - Every N years

Example recurrence rules:

```go
// Every 5 minutes
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyMinutely)
rr.SetInterval(5)

// Every day at 9 AM
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyDaily)
rr.SetInterval(1)
rr.SetStartsAt("2025-01-01 09:00:00")

// Every Monday and Friday
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyWeekly)
rr.SetInterval(1)
rr.SetDaysOfWeek([]taskstore.DayOfWeek{
    taskstore.DayOfWeekMonday,
    taskstore.DayOfWeekFriday,
})

### Recurrence Rule Types

Recurrence rules are expressed using the following types:

- **Frequency** (`taskstore.Frequency`) – how often the rule repeats:
  - `FrequencyNone` – one-time execution at `starts_at`
  - `FrequencySecondly`, `FrequencyMinutely`, `FrequencyHourly`
  - `FrequencyDaily`, `FrequencyWeekly`, `FrequencyMonthly`, `FrequencyYearly`
- **DayOfWeek** (`taskstore.DayOfWeek`) – days used for weekly rules:
  - `DayOfWeekMonday`, `DayOfWeekTuesday`, ..., `DayOfWeekSunday`
- **MonthOfYear** (`taskstore.MonthOfYear`) – months used for yearly/monthly rules:
  - `MonthOfYearJanuary` ... `MonthOfYearDecember`

The helper function `taskstore.NextRunAt(rule, now)` is used internally to compute the
next occurrence based on these fields.

### Schedule Helper Methods

`ScheduleInterface` provides a few convenience methods that encapsulate common
checks and updates:

- `HasReachedEndDate()` – `true` if the current time is after `end_at`.
- `HasReachedMaxExecutions()` – `true` if `max_execution_count` is set and
  `execution_count >= max_execution_count`.
- `GetNextOccurrence()` – returns the next run datetime (string) based on the
  recurrence rule, or an error if the rule is invalid.
- `IncrementExecutionCount()` – increments `execution_count` by one.
- `UpdateNextRunAt()` – recalculates and updates `next_run_at` using the
  recurrence rule.
- `UpdateLastRunAt()` – updates `last_run_at` to the current time.
- `IsDue()` – `true` if `next_run_at <= now`.

These helpers are useful when building custom schedulers or debugging schedule
state outside of `ScheduleRun`.

## CRUD Operations

```go
// Create
schedule := taskstore.NewSchedule()
// ... configure schedule ...
err := store.ScheduleCreate(ctx, schedule)

// Read
schedule, err := store.ScheduleFindByID(ctx, scheduleID)

// List
query := taskstore.NewScheduleQuery()
query.SetStatus("active")
query.SetLimit(10)
schedules, err := store.ScheduleList(ctx, query)

// Count
count, err := store.ScheduleCount(ctx, query)

// Update
schedule.SetName("Updated Name")
err := store.ScheduleUpdate(ctx, schedule)

// Soft Delete
err := store.ScheduleSoftDelete(ctx, schedule)
// or
err := store.ScheduleSoftDeleteByID(ctx, scheduleID)

// Hard Delete
err := store.ScheduleDelete(ctx, schedule)
// or
err := store.ScheduleDeleteByID(ctx, scheduleID)
```

## Soft Delete Behavior

- By default, `soft_deleted_at = 9999-12-31` (MAX_DATETIME), meaning NOT deleted
- Soft-deleted records have `soft_deleted_at` set to a datetime in the past
- Queries automatically filter out soft-deleted records using `WHERE soft_deleted_at > NOW()`
- To include soft-deleted records, you would need to modify the query (not currently exposed in the interface)

## Best Practices

1. **Run ScheduleRun() periodically** - Set up a cron job or background worker to call `ScheduleRun()` every minute
2. **Set appropriate intervals** - Don't set intervals too small (e.g., every second) to avoid overwhelming the queue
3. **Use max_execution_count** - For finite schedules, set a maximum to automatically complete them
4. **Set end_at dates** - For time-bound schedules, set an end date
5. **Monitor execution_count** - Track how many times a schedule has run
6. **Use descriptive names** - Make it easy to identify schedules in the database
7. **Test recurrence rules** - Use `NextRunAt()` function to verify your recurrence rule produces expected dates