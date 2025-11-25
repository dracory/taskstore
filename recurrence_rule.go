package taskstore

import (
	"encoding/json"
	"fmt"

	"github.com/dracory/sb"
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

// NextRunAt calculates the next time a recurrence rule should run, given the
// current time. It respects the rule's start and end times, interval and
// frequency, and returns an error if no further runs exist.
func NextRunAt(rule RecurrenceRuleInterface, now *carbon.Carbon) (*carbon.Carbon, error) {
	startsAt := parseDateTime(rule.GetStartsAt())

	endsAt := parseDateTime(rule.GetEndsAt())

	// If end time has passed, return max datetime to indicate no more runs
	if now.Gt(endsAt) {
		return carbon.Parse(sb.MAX_DATETIME, carbon.UTC), nil
	}

	if interval := rule.GetInterval(); interval <= 0 {
		return nil, fmt.Errorf("interval must be positive")
	}

	if now.Lt(startsAt) {
		return startsAt, nil
	}

	if rule.GetFrequency() == FrequencyNone {
		return startsAt, nil
	}

	freq := frequencyToRRuleFrequency(rule.GetFrequency())

	r, err := rrule.NewRRule(rrule.ROption{
		Freq:     freq,
		Interval: rule.GetInterval(),
		Count:    100,
		Dtstart:  startsAt.StdTime(),
	})

	if err != nil {
		return nil, err
	}

	times := r.Between(now.StdTime(), endsAt.StdTime(), true)

	if len(times) == 0 {
		return nil, fmt.Errorf("no more runs")
	}

	return carbon.Parse(times[0].String(), carbon.UTC), nil
}

func frequencyToRRuleFrequency(frequency Frequency) rrule.Frequency {
	switch frequency {
	case FrequencySecondly:
		return rrule.SECONDLY
	case FrequencyMinutely:
		return rrule.MINUTELY
	case FrequencyHourly:
		return rrule.HOURLY
	case FrequencyDaily:
		return rrule.DAILY
	case FrequencyWeekly:
		return rrule.WEEKLY
	case FrequencyMonthly:
		return rrule.MONTHLY
	case FrequencyYearly:
		return rrule.YEARLY
	default:
		return rrule.MAXYEAR
	}
}

// parseDateTime parses a UTC datetime string into a carbon instance.
func parseDateTime(dateTimeUTC string) *carbon.Carbon {
	return carbon.Parse(dateTimeUTC, carbon.UTC)
}

// NewRecurrenceRule creates a new recurrence rule with default values.
// By default, it has no end time (MAX_DATETIME) and an interval of 1.
func NewRecurrenceRule() RecurrenceRuleInterface {
	r := recurrenceRule{}

	// By default, it does not have an end time
	r.SetEndsAt(sb.MAX_DATETIME)

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

// Clone creates a shallow copy of the recurrence rule.
func (r *recurrenceRule) Clone() RecurrenceRuleInterface {
	return &recurrenceRule{
		frequency:    r.frequency,
		startsAt:     r.startsAt,
		endsAt:       r.endsAt,
		interval:     r.interval,
		daysOfWeek:   r.daysOfWeek,
		daysOfMonth:  r.daysOfMonth,
		monthsOfYear: r.monthsOfYear,
	}
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
