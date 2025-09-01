package taskstore

import (
	"encoding/json"
	"fmt"

	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/teambition/rrule-go"
)

// Define a string type alias
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

type RecurrenceRule interface {
	GetFrequency() Frequency
	SetFrequency(Frequency) RecurrenceRule

	GetStartsAt() string
	SetStartsAt(dateTimeUTC string) RecurrenceRule

	GetEndsAt() string
	SetEndsAt(dateTimeUTC string) RecurrenceRule

	GetInterval() int
	SetInterval(int) RecurrenceRule

	GetDaysOfWeek() []DayOfWeek
	SetDaysOfWeek([]DayOfWeek) RecurrenceRule

	GetDaysOfMonth() []int
	SetDaysOfMonth([]int) RecurrenceRule

	GetMonthsOfYear() []MonthOfYear
	SetMonthsOfYear([]MonthOfYear) RecurrenceRule
}

func NextRunAt(rule RecurrenceRule, now *carbon.Carbon) (*carbon.Carbon, error) {
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

// func NextRunAt(rule RecurrenceRule, now carbon.Carbon) (carbon.Carbon, error) {
// 	startsAt := parseDateTime(rule.GetStartsAt())

// 	endsAt := parseDateTime(rule.GetEndsAt())

// 	// If end time has passed, return max datetime to indicate no more runs
// 	if now.Gt(endsAt) {
// 		return carbon.Parse(sb.MAX_DATETIME, carbon.UTC), nil
// 	}

// 	if interval := rule.GetInterval(); interval <= 0 {
// 		return carbon.Carbon{}, fmt.Errorf("interval must be positive")
// 	}

// 	if now.Lt(startsAt) {
// 		return startsAt, nil
// 	}

// 	nextRun, err := calculateNextRun(rule, now)

// 	if err != nil {
// 		return carbon.Carbon{}, err
// 	}

// 	// If next run is after end time, return max datetime to indicate no more runs
// 	if nextRun.Gt(endsAt) {
// 		return carbon.Parse(sb.MAX_DATETIME, carbon.UTC), nil
// 	}

// 	return nextRun, nil
// }

// func calculateNextRun(rule RecurrenceRule, now carbon.Carbon) (carbon.Carbon, error) {
// 	interval := rule.GetInterval()

// 	switch rule.GetFrequency() {
// 	case FrequencySecondly:
// 		return now.AddSeconds(interval), nil
// 	case FrequencyMinutely:
// 		return now.AddMinutes(interval), nil
// 	case FrequencyHourly:
// 		return now.AddHours(interval), nil
// 	case FrequencyDaily:
// 		return calculateDailyNextRun(rule, now)
// 	case FrequencyWeekly:
// 		return calculateWeeklyNextRun(rule, now)
// 	case FrequencyMonthly:
// 		return now.AddMonths(interval), nil
// 	case FrequencyYearly:
// 		return now.AddYears(interval), nil
// 	default:
// 		return carbon.Carbon{}, fmt.Errorf("unknown frequency")
// 	}
// }

// func calculateWeeklyNextRun(rule RecurrenceRule, now carbon.Carbon) (carbon.Carbon, error) {
// 	daysOfWeek := rule.GetDaysOfWeek()
// 	startsAt := parseDateTime(rule.GetStartsAt())

// 	daysToAdd := calculateWeeklyDaysToAdd(now, daysOfWeek)

// 	if daysToAdd == 0 && now.Hour()*60+now.Minute() >= startsAt.Hour()*60+startsAt.Minute() {
// 		daysToAdd = 7
// 	}

// 	return startsAt.AddDays(daysToAdd + (rule.GetInterval()-1)*7), nil
// }

// func calculateDailyNextRun(rule RecurrenceRule, now carbon.Carbon) (carbon.Carbon, error) {
// 	startsAt := parseDateTime(rule.GetStartsAt())

// 	diffDays := cast.ToInt(now.DiffInDays(startsAt)) + 1

// 	// If starts in the future, return startsAt
// 	if diffDays <= 0 {
// 		return startsAt, nil
// 	}

// 	interval := rule.GetInterval()

// 	if interval <= 0 {
// 		return carbon.Carbon{}, fmt.Errorf("interval must be positive")
// 	}

// 	daysToAdd := (diffDays / interval) * interval
// 	if diffDays%interval != 0 {
// 		daysToAdd += interval
// 	}

// 	return startsAt.AddDays(daysToAdd), nil
// }

// func calculateWeeklyDaysToAdd(now carbon.Carbon, daysOfWeek []DayOfWeek) int {
// 	dayOfWeek := now.DayOfWeek()

// 	if len(daysOfWeek) == 0 {
// 		return (7 - dayOfWeek) % 7
// 	}

// 	sort.Slice(daysOfWeek, func(i, j int) bool {
// 		return dayOfWeekToInt(daysOfWeek[i]) < dayOfWeekToInt(daysOfWeek[j])
// 	})

// 	nextDayOfWeek := -1
// 	for _, day := range daysOfWeek {
// 		dayInt := dayOfWeekToInt(day)
// 		if dayInt > dayOfWeek {
// 			nextDayOfWeek = dayInt
// 			break
// 		}
// 	}

// 	if nextDayOfWeek == -1 {
// 		nextDayOfWeek = dayOfWeekToInt(daysOfWeek[0])
// 	}

// 	return (nextDayOfWeek - int(dayOfWeek) + 7) % 7
// }

func parseDateTime(dateTimeUTC string) *carbon.Carbon {
	return carbon.Parse(dateTimeUTC, carbon.UTC)
}

// func UNUSED_dayOfWeekToInt(day DayOfWeek) int {
// 	switch day {
// 	case DayOfWeekSunday:
// 		return 0
// 	case DayOfWeekMonday:
// 		return 1
// 	case DayOfWeekTuesday:
// 		return 2
// 	case DayOfWeekWednesday:
// 		return 3
// 	case DayOfWeekThursday:
// 		return 4
// 	case DayOfWeekFriday:
// 		return 5
// 	case DayOfWeekSaturday:
// 		return 6
// 	default:
// 		return 0
// 	}
// }

// func UNUSED_monthOfYearToInt(month MonthOfYear) int {
// 	switch month {
// 	case MonthOfYearJanuary:
// 		return 1
// 	case MonthOfYearFebruary:
// 		return 2
// 	case MonthOfYearMarch:
// 		return 3
// 	case MonthOfYearApril:
// 		return 4
// 	case MonthOfYearMay:
// 		return 5
// 	case MonthOfYearJune:
// 		return 6
// 	case MonthOfYearJuly:
// 		return 7
// 	case MonthOfYearAugust:
// 		return 8
// 	case MonthOfYearSeptember:
// 		return 9
// 	case MonthOfYearOctober:
// 		return 10
// 	case MonthOfYearNovember:
// 		return 11
// 	case MonthOfYearDecember:
// 		return 12
// 	default:
// 		return 0
// 	}
// }

func NewRecurrenceRule() RecurrenceRule {
	r := recurrenceRule{}

	// By default, it does not have an end time
	r.SetEndsAt(sb.MAX_DATETIME)

	// By default, the interval is 1
	r.SetInterval(1)

	return &r
}

type recurrenceRule struct {
	frequency    Frequency
	startsAt     string
	endsAt       string
	interval     int
	daysOfWeek   []DayOfWeek
	daysOfMonth  []int
	monthsOfYear []MonthOfYear
}

func (r *recurrenceRule) GetFrequency() Frequency {
	return r.frequency
}

func (r *recurrenceRule) SetFrequency(frequency Frequency) RecurrenceRule {
	r.frequency = frequency
	return r
}

func (r *recurrenceRule) GetStartsAt() string {
	return r.startsAt
}

func (r *recurrenceRule) SetStartsAt(startsAt string) RecurrenceRule {
	r.startsAt = startsAt
	return r
}

func (r *recurrenceRule) GetEndsAt() string {
	return r.endsAt
}

func (r *recurrenceRule) SetEndsAt(endsAt string) RecurrenceRule {
	r.endsAt = endsAt
	return r
}

func (r *recurrenceRule) GetInterval() int {
	return r.interval
}

func (r *recurrenceRule) SetInterval(interval int) RecurrenceRule {
	r.interval = interval
	return r
}

func (r *recurrenceRule) GetDaysOfWeek() []DayOfWeek {
	return r.daysOfWeek
}

func (r *recurrenceRule) SetDaysOfWeek(daysOfWeek []DayOfWeek) RecurrenceRule {
	r.daysOfWeek = daysOfWeek
	return r
}

func (r *recurrenceRule) GetDaysOfMonth() []int {
	return r.daysOfMonth
}

func (r *recurrenceRule) SetDaysOfMonth(daysOfMonth []int) RecurrenceRule {
	r.daysOfMonth = daysOfMonth
	return r
}

func (r *recurrenceRule) GetMonthsOfYear() []MonthOfYear {
	return r.monthsOfYear
}

func (r *recurrenceRule) SetMonthsOfYear(monthsOfYear []MonthOfYear) RecurrenceRule {
	r.monthsOfYear = monthsOfYear
	return r
}

func (r *recurrenceRule) String() string {
	return fmt.Sprintf("frequency: %s, startsAt: %s, endsAt: %s, interval: %d, daysOfWeek: %v, daysOfMonth: %v, monthsOfYear: %v",
		r.frequency, r.startsAt, r.endsAt, r.interval, r.daysOfWeek, r.daysOfMonth, r.monthsOfYear)
}

func (r *recurrenceRule) Clone() RecurrenceRule {
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
