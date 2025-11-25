# Recurrence Rules

Recurrence rules define when and how often scheduled tasks should run. They are based on the iCalendar RFC 5545 standard and use the `rrule-go` library internally.

## Overview

A recurrence rule consists of:
- **Frequency** - How often the task repeats (daily, weekly, monthly, etc.)
- **Interval** - The multiplier for the frequency (e.g., every 2 days, every 3 weeks)
- **Starts At** - When the recurrence begins
- **Ends At** - When the recurrence stops (optional)
- **Days of Week** - Specific days for weekly recurrence (optional)
- **Days of Month** - Specific days for monthly recurrence (optional)
- **Months of Year** - Specific months for yearly recurrence (optional)

## Frequency Types

### FrequencyNone
One-time execution at the `starts_at` time. No recurrence.

```go
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyNone)
rr.SetStartsAt("2025-12-25 09:00:00")
// Runs once on December 25, 2025 at 9:00 AM
```

### FrequencySecondly
Repeats every N seconds.

```go
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencySecondly)
rr.SetInterval(30)
// Runs every 30 seconds
```

⚠️ **Warning**: Use with caution. Very frequent schedules can overwhelm your queue.

### FrequencyMinutely
Repeats every N minutes.

```go
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyMinutely)
rr.SetInterval(15)
// Runs every 15 minutes
```

### FrequencyHourly
Repeats every N hours.

```go
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyHourly)
rr.SetInterval(2)
rr.SetStartsAt("2025-01-01 08:00:00")
// Runs every 2 hours starting at 8:00 AM (8:00, 10:00, 12:00, etc.)
```

### FrequencyDaily
Repeats every N days.

```go
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyDaily)
rr.SetInterval(1)
rr.SetStartsAt("2025-01-01 09:00:00")
// Runs every day at 9:00 AM
```

### FrequencyWeekly
Repeats every N weeks, optionally on specific days.

```go
// Every week on Monday and Friday
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyWeekly)
rr.SetInterval(1)
rr.SetDaysOfWeek([]taskstore.DayOfWeek{
    taskstore.DayOfWeekMonday,
    taskstore.DayOfWeekFriday,
})
rr.SetStartsAt("2025-01-01 09:00:00")

// Every 2 weeks on Wednesday
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyWeekly)
rr.SetInterval(2)
rr.SetDaysOfWeek([]taskstore.DayOfWeek{
    taskstore.DayOfWeekWednesday,
})
```

**Available Days**:
- `DayOfWeekMonday`
- `DayOfWeekTuesday`
- `DayOfWeekWednesday`
- `DayOfWeekThursday`
- `DayOfWeekFriday`
- `DayOfWeekSaturday`
- `DayOfWeekSunday`

### FrequencyMonthly
Repeats every N months, optionally on specific days.

```go
// Every month on the 1st and 15th
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyMonthly)
rr.SetInterval(1)
rr.SetDaysOfMonth([]int{1, 15})
rr.SetStartsAt("2025-01-01 09:00:00")

// Every 3 months on the 10th
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyMonthly)
rr.SetInterval(3)
rr.SetDaysOfMonth([]int{10})
```

### FrequencyYearly
Repeats every N years, optionally in specific months.

```go
// Every year on January 1st
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyYearly)
rr.SetInterval(1)
rr.SetMonthsOfYear([]taskstore.MonthOfYear{
    taskstore.MonthOfYearJanuary,
})
rr.SetStartsAt("2025-01-01 00:00:00")

// Every 2 years in March and September
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyYearly)
rr.SetInterval(2)
rr.SetMonthsOfYear([]taskstore.MonthOfYear{
    taskstore.MonthOfYearMarch,
    taskstore.MonthOfYearSeptember,
})
```

**Available Months**:
- `MonthOfYearJanuary` through `MonthOfYearDecember`

## How NextRunAt Works

The `NextRunAt()` function calculates when a schedule should run next based on the recurrence rule and current time.

### Algorithm

1. **Check if end time has passed**
   - If `now > ends_at`, return `MAX_DATETIME` (no more runs)

2. **Validate interval**
   - If `interval <= 0`, return error

3. **Check if before start time**
   - If `now < starts_at`, return `starts_at`

4. **Handle FrequencyNone**
   - Return `starts_at` (one-time execution)

5. **Calculate next occurrence**
   - Convert frequency to rrule format
   - Create rrule with frequency, interval, and start time
   - Find next occurrence between `now` and `ends_at`
   - Return the next occurrence time

### Example Calculation

```go
// Schedule: Daily at 9:00 AM
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyDaily)
rr.SetInterval(1)
rr.SetStartsAt("2025-01-01 09:00:00")

now := carbon.Parse("2025-01-15 14:30:00", carbon.UTC)
nextRun, err := taskstore.NextRunAt(rr, now)
// nextRun = "2025-01-16 09:00:00" (next day at 9:00 AM)
```

## Common Patterns

### Business Days Only (Monday-Friday)

```go
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyWeekly)
rr.SetInterval(1)
rr.SetDaysOfWeek([]taskstore.DayOfWeek{
    taskstore.DayOfWeekMonday,
    taskstore.DayOfWeekTuesday,
    taskstore.DayOfWeekWednesday,
    taskstore.DayOfWeekThursday,
    taskstore.DayOfWeekFriday,
})
rr.SetStartsAt("2025-01-01 09:00:00")
```

### First Day of Every Month

```go
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyMonthly)
rr.SetInterval(1)
rr.SetDaysOfMonth([]int{1})
rr.SetStartsAt("2025-01-01 00:00:00")
```

### Quarterly Reports (Every 3 Months)

```go
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyMonthly)
rr.SetInterval(3)
rr.SetDaysOfMonth([]int{1})
rr.SetStartsAt("2025-01-01 09:00:00")
```

### Every 6 Hours

```go
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyHourly)
rr.SetInterval(6)
rr.SetStartsAt("2025-01-01 00:00:00")
// Runs at 00:00, 06:00, 12:00, 18:00
```

### Weekend Only (Saturday and Sunday)

```go
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyWeekly)
rr.SetInterval(1)
rr.SetDaysOfWeek([]taskstore.DayOfWeek{
    taskstore.DayOfWeekSaturday,
    taskstore.DayOfWeekSunday,
})
rr.SetStartsAt("2025-01-01 10:00:00")
```

## Time Zones

⚠️ **Important**: All times are stored and processed in UTC.

- Always provide times in UTC format: `"2025-01-01 09:00:00"`
- The system does not handle time zone conversions
- If you need local time, convert to UTC before setting `starts_at`

Example:
```go
// If you want 9 AM EST (UTC-5), use 2 PM UTC
rr.SetStartsAt("2025-01-01 14:00:00") // 9 AM EST = 2 PM UTC
```

## Limits and Constraints

### Default Values
- **Interval**: Default is `1` (set in `NewRecurrenceRule()`)
- **Ends At**: Default is `MAX_DATETIME` (9999-12-31) - no end date

### Constraints
- **Interval must be positive**: `interval > 0`
- **Starts At**: Must be a valid datetime string
- **Ends At**: Must be after `starts_at`
- **Days of Month**: Valid range is 1-31
- **Count Limit**: Internal rrule uses `Count: 100` to limit calculations

## JSON Serialization

Recurrence rules are automatically serialized to/from JSON when saving schedules:

```json
{
  "frequency": "daily",
  "startsAt": "2025-01-01 09:00:00",
  "endsAt": "9999-12-31 23:59:59",
  "interval": 1,
  "daysOfWeek": [],
  "daysOfMonth": [],
  "monthsOfYear": []
}
```

## Testing Recurrence Rules

You can test your recurrence rules before creating a schedule:

```go
rr := taskstore.NewRecurrenceRule()
rr.SetFrequency(taskstore.FrequencyDaily)
rr.SetInterval(1)
rr.SetStartsAt("2025-01-01 09:00:00")

// Test what the next run would be
now := carbon.Now(carbon.UTC)
nextRun, err := taskstore.NextRunAt(rr, now)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Next run: %s\n", nextRun.ToDateTimeString())
```

## Best Practices

1. **Use appropriate frequencies** - Don't use `FrequencySecondly` unless absolutely necessary
2. **Set end dates for finite schedules** - Use `SetEndsAt()` to prevent infinite execution
3. **Align start times** - Set `starts_at` to when you want the first execution
4. **Test your rules** - Use `NextRunAt()` to verify the schedule produces expected times
5. **Consider time zones** - Remember all times are UTC
6. **Use intervals wisely** - `interval: 2` with `FrequencyDaily` means every 2 days, not twice a day
7. **Combine with max_execution_count** - For extra safety, set a maximum number of executions on the schedule

## Troubleshooting

### Schedule not running
- Check that `status = "active"`
- Verify `next_run_at <= NOW()`
- Ensure `ScheduleRun()` is being called periodically
- Check that `interval > 0`

### Unexpected run times
- Verify `starts_at` is in UTC
- Check that `interval` is set correctly
- Ensure `days_of_week` or `days_of_month` are set as intended

### Schedule completed too early
- Check `ends_at` value
- Verify `max_execution_count` on the schedule
- Ensure `NextRunAt` isn't returning `MAX_DATETIME`
