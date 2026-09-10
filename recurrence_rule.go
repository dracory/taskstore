package taskstore

import (
	"encoding/json"
	"fmt"

	"github.com/dromara/carbon/v2"
	"github.com/teambition/rrule-go"
)

// Define a string type alias
// Frequency represents how often a schedule recurs (daily, weekly, etc.).
// It is a string-based alias compatible with rrule-go frequencies.
type Frequency string

// Define the constants as strings
const (
	FrequencyNone     Frequency = "none"
	FrequencySecondly Frequency = "secondly"
	FrequencyMinutely Frequency = "minutely"
	FrequencyHourly   Frequency = "hourly"
	FrequencyDaily    Frequency = "daily"
	FrequencyWeekly   Frequency = "weekly"
	FrequencyMonthly  Frequency = "monthly"
	FrequencyYearly   Frequency = "yearly"
)

// DayOfWeek represents a day of the week used in weekly recurrence rules.
type DayOfWeek string

const (
	DayOfWeekMonday    DayOfWeek = "monday"
	DayOfWeekTuesday   DayOfWeek = "tuesday"
	DayOfWeekWednesday DayOfWeek = "wednesday"
	DayOfWeekThursday  DayOfWeek = "thursday"
	DayOfWeekFriday    DayOfWeek = "friday"
	DayOfWeekSaturday  DayOfWeek = "saturday"
	DayOfWeekSunday    DayOfWeek = "sunday"
)

// MonthOfYear represents a month used in yearly or monthly recurrence rules.
type MonthOfYear string

const (
	MonthOfYearJanuary   MonthOfYear = "JANUARY"
	MonthOfYearFebruary  MonthOfYear = "FEBRUARY"
	MonthOfYearMarch     MonthOfYear = "MARCH"
	MonthOfYearApril     MonthOfYear = "APRIL"
	MonthOfYearMay       MonthOfYear = "MAY"
	MonthOfYearJune      MonthOfYear = "JUNE"
	MonthOfYearJuly      MonthOfYear = "JULY"
	MonthOfYearAugust    MonthOfYear = "AUGUST"
	MonthOfYearSeptember MonthOfYear = "SEPTEMBER"
	MonthOfYearOctober   MonthOfYear = "OCTOBER"
	MonthOfYearNovember  MonthOfYear = "NOVEMBER"
	MonthOfYearDecember  MonthOfYear = "DECEMBER"
)

// RecurrenceRuleInterface defines the contract for recurrence rules used by schedules.
// It exposes frequency, start/end times, interval, and optional day/month filters.
type RecurrenceRuleInterface interface {
	// GetFrequency returns how often the rule recurs (e.g. daily, weekly).
	GetFrequency() Frequency

	// SetFrequency sets how often the rule recurs.
	SetFrequency(Frequency) RecurrenceRuleInterface

	// GetStartsAt returns the UTC datetime when the rule becomes active.
	GetStartsAt() string

	// SetStartsAt sets the UTC datetime when the rule becomes active.
	SetStartsAt(dateTimeUTC string) RecurrenceRuleInterface

	// GetEndsAt returns the UTC datetime when the rule stops producing occurrences.
	GetEndsAt() string

	// SetEndsAt sets the UTC datetime when the rule stops producing occurrences.
	SetEndsAt(dateTimeUTC string) RecurrenceRuleInterface

	// GetInterval returns the step interval between occurrences (e.g. every N days).
	GetInterval() int

	// SetInterval sets the step interval between occurrences.
	SetInterval(int) RecurrenceRuleInterface

	// GetDaysOfWeek returns the days of the week the rule applies to (for weekly rules).
	GetDaysOfWeek() []DayOfWeek

	// SetDaysOfWeek sets the days of the week the rule applies to (for weekly rules).
	SetDaysOfWeek([]DayOfWeek) RecurrenceRuleInterface

	// GetDaysOfMonth returns the days of the month the rule applies to.
	GetDaysOfMonth() []int

	// SetDaysOfMonth sets the days of the month the rule applies to.
	SetDaysOfMonth([]int) RecurrenceRuleInterface

	// GetMonthsOfYear returns the months of the year the rule applies to.
	GetMonthsOfYear() []MonthOfYear

	// SetMonthsOfYear sets the months of the year the rule applies to.
	SetMonthsOfYear([]MonthOfYear) RecurrenceRuleInterface
}

// ErrNoMoreRuns is returned by NextRunAt when a recurrence rule has no more
// occurrences in the future. Callers (e.g. UpdateNextRunAt) use this sentinel
// to distinguish "schedule is done" from transient errors, and set
// NextRunAtField to MAX_DATETIME so IsDue() returns false and the runner
// marks the schedule "completed".
var ErrNoMoreRuns = fmt.Errorf("no more runs")

// NextRunAt calculates the next time a recurrence rule should run, given the
// current time.
//
// Business Logic:
//  1. If now is past endsAt, return ErrNoMoreRuns — the schedule's end time
//     has passed.
//  2. If interval is not positive, return an error — the rule is invalid.
//  3. If now is before startsAt, return startsAt — the first run hasn't
//     happened yet.
//  4. If frequency is FrequencyNone (one-time schedule) and now >= startsAt,
//     return ErrNoMoreRuns — the single run has already passed.
//  5. Otherwise, build an rrule and ask it for the first occurrence strictly
//     after now. If none exists (rrule exhausted or next exceeds endsAt),
//     return ErrNoMoreRuns.
//  6. This function never returns MAX_DATETIME — that is a storage-layer
//     sentinel. UpdateNextRunAt translates ErrNoMoreRuns into MAX_DATETIME.
func NextRunAt(rule RecurrenceRuleInterface, now *carbon.Carbon) (*carbon.Carbon, error) {
	startsAt := parseDateTime(rule.GetStartsAt())

	endsAt := parseDateTime(rule.GetEndsAt())

	// If end time has passed, no more runs.
	if now.Gt(endsAt) {
		return nil, ErrNoMoreRuns
	}

	if interval := rule.GetInterval(); interval <= 0 {
		return nil, fmt.Errorf("interval must be positive")
	}

	if now.Lt(startsAt) {
		return startsAt, nil
	}

	if rule.GetFrequency() == FrequencyNone {
		// One-time schedule: if the start time has already passed, there
		// are no more runs. Returning startsAt here would cause IsDue() to
		// return true on every tick, firing the schedule repeatedly.
		if now.Gte(startsAt) {
			return nil, ErrNoMoreRuns
		}
		return startsAt, nil
	}

	freq, err := frequencyToRRuleFrequency(rule.GetFrequency())
	if err != nil {
		return nil, err
	}

	// Count is intentionally left at 0 (infinite). A non-zero Count limits
	// the rrule to that many total occurrences from Dtstart; once exhausted,
	// After() returns a zero time, UpdateNextRunAt() silently keeps the old
	// (past) next_run_at, and the schedule fires on every runner tick.
	r, err := rrule.NewRRule(rrule.ROption{
		Freq:       freq,
		Interval:   rule.GetInterval(),
		Dtstart:    startsAt.StdTime(),
		Byweekday:  daysOfWeekToRRuleWeekdays(rule.GetDaysOfWeek()),
		Bymonthday: rule.GetDaysOfMonth(),
		Bymonth:    monthsOfYearToRRuleMonths(rule.GetMonthsOfYear()),
	})

	if err != nil {
		return nil, err
	}

	// After() with inc=false returns the first occurrence strictly after now.
	// It is more efficient than Between() because it stops at the first match
	// instead of collecting all occurrences up to endsAt.
	next := r.After(now.StdTime(), false)

	if next.IsZero() {
		return nil, ErrNoMoreRuns
	}

	// Honour the rule's end time even though the rrule itself is infinite.
	if next.After(endsAt.StdTime()) {
		return nil, ErrNoMoreRuns
	}

	return carbon.CreateFromStdTime(next), nil
}

func frequencyToRRuleFrequency(frequency Frequency) (rrule.Frequency, error) {
	switch frequency {
	case FrequencySecondly:
		return rrule.SECONDLY, nil
	case FrequencyMinutely:
		return rrule.MINUTELY, nil
	case FrequencyHourly:
		return rrule.HOURLY, nil
	case FrequencyDaily:
		return rrule.DAILY, nil
	case FrequencyWeekly:
		return rrule.WEEKLY, nil
	case FrequencyMonthly:
		return rrule.MONTHLY, nil
	case FrequencyYearly:
		return rrule.YEARLY, nil
	default:
		return rrule.MAXYEAR, fmt.Errorf("unknown frequency: %s", frequency)
	}
}

// dayOfWeekToRRuleWeekday converts a DayOfWeek string to the corresponding
// rrule.Weekday constant. Returns ok=false for unknown day names.
func dayOfWeekToRRuleWeekday(d DayOfWeek) (rrule.Weekday, bool) {
	switch d {
	case DayOfWeekMonday:
		return rrule.MO, true
	case DayOfWeekTuesday:
		return rrule.TU, true
	case DayOfWeekWednesday:
		return rrule.WE, true
	case DayOfWeekThursday:
		return rrule.TH, true
	case DayOfWeekFriday:
		return rrule.FR, true
	case DayOfWeekSaturday:
		return rrule.SA, true
	case DayOfWeekSunday:
		return rrule.SU, true
	default:
		return rrule.Weekday{}, false
	}
}

// daysOfWeekToRRuleWeekdays converts a slice of DayOfWeek strings to rrule
// Weekday values. Unknown day names are silently skipped.
func daysOfWeekToRRuleWeekdays(days []DayOfWeek) []rrule.Weekday {
	if len(days) == 0 {
		return nil
	}
	result := make([]rrule.Weekday, 0, len(days))
	for _, d := range days {
		if wd, ok := dayOfWeekToRRuleWeekday(d); ok {
			result = append(result, wd)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// monthsOfYearToRRuleMonths converts a slice of MonthOfYear strings to rrule
// month numbers (1=January … 12=December). Unknown month names are skipped.
func monthsOfYearToRRuleMonths(months []MonthOfYear) []int {
	if len(months) == 0 {
		return nil
	}
	result := make([]int, 0, len(months))
	for _, m := range months {
		switch m {
		case MonthOfYearJanuary:
			result = append(result, 1)
		case MonthOfYearFebruary:
			result = append(result, 2)
		case MonthOfYearMarch:
			result = append(result, 3)
		case MonthOfYearApril:
			result = append(result, 4)
		case MonthOfYearMay:
			result = append(result, 5)
		case MonthOfYearJune:
			result = append(result, 6)
		case MonthOfYearJuly:
			result = append(result, 7)
		case MonthOfYearAugust:
			result = append(result, 8)
		case MonthOfYearSeptember:
			result = append(result, 9)
		case MonthOfYearOctober:
			result = append(result, 10)
		case MonthOfYearNovember:
			result = append(result, 11)
		case MonthOfYearDecember:
			result = append(result, 12)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// parseDateTime parses a UTC datetime string into a carbon instance.
func parseDateTime(dateTimeUTC string) *carbon.Carbon {
	return carbon.Parse(dateTimeUTC, carbon.UTC)
}

// NewRecurrenceRule creates a new recurrence rule with default values.
// By default, it has no end time (MAX_DATETIME) and an interval of 1.
func NewRecurrenceRule() RecurrenceRuleInterface {
	r := recurrenceRule{}

	// By default, it runs once (no recurrence)
	r.SetFrequency(FrequencyNone)

	// By default, it does not have an end time
	r.SetEndsAt(MAX_DATETIME)

	// By default, the interval is 1
	r.SetInterval(1)

	return &r
}

// recurrenceRule is the concrete implementation of RecurrenceRuleInterface.
// It stores all recurrence fields including frequency, timing, and filters.
type recurrenceRule struct {
	frequency    Frequency
	startsAt     string
	endsAt       string
	interval     int
	daysOfWeek   []DayOfWeek
	daysOfMonth  []int
	monthsOfYear []MonthOfYear
}

// GetFrequency returns how often the rule recurs.
func (r *recurrenceRule) GetFrequency() Frequency {
	return r.frequency
}

// SetFrequency sets how often the rule recurs.
func (r *recurrenceRule) SetFrequency(frequency Frequency) RecurrenceRuleInterface {
	r.frequency = frequency
	return r
}

// GetStartsAt returns the UTC datetime when the rule becomes active.
func (r *recurrenceRule) GetStartsAt() string {
	return r.startsAt
}

// SetStartsAt sets the UTC datetime when the rule becomes active.
func (r *recurrenceRule) SetStartsAt(startsAt string) RecurrenceRuleInterface {
	r.startsAt = startsAt
	return r
}

// GetEndsAt returns the UTC datetime when the rule stops producing occurrences.
func (r *recurrenceRule) GetEndsAt() string {
	return r.endsAt
}

// SetEndsAt sets the UTC datetime when the rule stops producing occurrences.
func (r *recurrenceRule) SetEndsAt(endsAt string) RecurrenceRuleInterface {
	r.endsAt = endsAt
	return r
}

// GetInterval returns the step interval between occurrences.
func (r *recurrenceRule) GetInterval() int {
	return r.interval
}

// SetInterval sets the step interval between occurrences.
func (r *recurrenceRule) SetInterval(interval int) RecurrenceRuleInterface {
	r.interval = interval
	return r
}

// GetDaysOfWeek returns the days of the week the rule applies to.
func (r *recurrenceRule) GetDaysOfWeek() []DayOfWeek {
	return r.daysOfWeek
}

// SetDaysOfWeek sets the days of the week the rule applies to.
func (r *recurrenceRule) SetDaysOfWeek(daysOfWeek []DayOfWeek) RecurrenceRuleInterface {
	r.daysOfWeek = daysOfWeek
	return r
}

// GetDaysOfMonth returns the days of the month the rule applies to.
func (r *recurrenceRule) GetDaysOfMonth() []int {
	return r.daysOfMonth
}

// SetDaysOfMonth sets the days of the month the rule applies to.
func (r *recurrenceRule) SetDaysOfMonth(daysOfMonth []int) RecurrenceRuleInterface {
	r.daysOfMonth = daysOfMonth
	return r
}

// GetMonthsOfYear returns the months of the year the rule applies to.
func (r *recurrenceRule) GetMonthsOfYear() []MonthOfYear {
	return r.monthsOfYear
}

// SetMonthsOfYear sets the months of the year the rule applies to.
func (r *recurrenceRule) SetMonthsOfYear(monthsOfYear []MonthOfYear) RecurrenceRuleInterface {
	r.monthsOfYear = monthsOfYear
	return r
}

// String returns a human-readable representation of the recurrence rule.
func (r *recurrenceRule) String() string {
	return fmt.Sprintf("frequency: %s, startsAt: %s, endsAt: %s, interval: %d, daysOfWeek: %v, daysOfMonth: %v, monthsOfYear: %v",
		r.frequency, r.startsAt, r.endsAt, r.interval, r.daysOfWeek, r.daysOfMonth, r.monthsOfYear)
}

// Clone creates a deep copy of the recurrence rule.
func (r *recurrenceRule) Clone() RecurrenceRuleInterface {
	clone := &recurrenceRule{
		frequency:    r.frequency,
		startsAt:     r.startsAt,
		endsAt:       r.endsAt,
		interval:     r.interval,
		daysOfWeek:   append([]DayOfWeek(nil), r.daysOfWeek...),
		daysOfMonth:  append([]int(nil), r.daysOfMonth...),
		monthsOfYear: append([]MonthOfYear(nil), r.monthsOfYear...),
	}
	return clone
}

// MarshalJSON serializes the recurrence rule into JSON.
func (r *recurrenceRule) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Frequency    Frequency     `json:"frequency"`
		StartsAt     string        `json:"startsAt"`
		EndsAt       string        `json:"endsAt"`
		Interval     int           `json:"interval"`
		DaysOfWeek   []DayOfWeek   `json:"daysOfWeek"`
		DaysOfMonth  []int         `json:"daysOfMonth"`
		MonthsOfYear []MonthOfYear `json:"monthsOfYear"`
	}{
		Frequency:    r.frequency,
		StartsAt:     r.startsAt,
		EndsAt:       r.endsAt,
		Interval:     r.interval,
		DaysOfWeek:   r.daysOfWeek,
		DaysOfMonth:  r.daysOfMonth,
		MonthsOfYear: r.monthsOfYear,
	})
}

// UnmarshalJSON deserializes the recurrence rule from JSON.
func (r *recurrenceRule) UnmarshalJSON(data []byte) error {
	var v struct {
		Frequency    Frequency     `json:"frequency"`
		StartsAt     string        `json:"startsAt"`
		EndsAt       string        `json:"endsAt"`
		Interval     int           `json:"interval"`
		DaysOfWeek   []DayOfWeek   `json:"daysOfWeek"`
		DaysOfMonth  []int         `json:"daysOfMonth"`
		MonthsOfYear []MonthOfYear `json:"monthsOfYear"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*r = recurrenceRule{
		frequency:    v.Frequency,
		startsAt:     v.StartsAt,
		endsAt:       v.EndsAt,
		interval:     v.Interval,
		daysOfWeek:   v.DaysOfWeek,
		daysOfMonth:  v.DaysOfMonth,
		monthsOfYear: v.MonthsOfYear,
	}
	return nil
}
