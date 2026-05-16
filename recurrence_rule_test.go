package taskstore

import (
	"strings"
	"testing"

	"github.com/dromara/carbon/v2"
	"github.com/teambition/rrule-go"
)

func Test_frequencyToRRuleFrequency(t *testing.T) {
	tests := []struct {
		name      string
		frequency Frequency
		want      rrule.Frequency
	}{
		{
			name:      "secondly frequency",
			frequency: FrequencySecondly,
			want:      rrule.SECONDLY,
		},
		{
			name:      "minutely frequency",
			frequency: FrequencyMinutely,
			want:      rrule.MINUTELY,
		},
		{
			name:      "hourly frequency",
			frequency: FrequencyHourly,
			want:      rrule.HOURLY,
		},
		{
			name:      "daily frequency",
			frequency: FrequencyDaily,
			want:      rrule.DAILY,
		},
		{
			name:      "weekly frequency",
			frequency: FrequencyWeekly,
			want:      rrule.WEEKLY,
		},
		{
			name:      "monthly frequency",
			frequency: FrequencyMonthly,
			want:      rrule.MONTHLY,
		},
		{
			name:      "yearly frequency",
			frequency: FrequencyYearly,
			want:      rrule.YEARLY,
		},
		{
			name:      "unknown frequency defaults to MAXYEAR",
			frequency: FrequencyNone,
			want:      rrule.MAXYEAR,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := frequencyToRRuleFrequency(tt.frequency); got != tt.want {
				t.Errorf("frequencyToRRuleFrequency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseDateTime(t *testing.T) {
	tests := []struct {
		name string
		date string
	}{
		{
			name: "valid datetime",
			date: "2024-01-01 00:00:00",
		},
		{
			name: "datetime with timezone",
			date: "2024-01-01T00:00:00Z",
		},
		{
			name: "empty string",
			date: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseDateTime(tt.date)
			if got == nil {
				t.Error("parseDateTime() should not return nil")
			}
		})
	}
}

func TestRecurrenceRule_Clone(t *testing.T) {
	original := NewRecurrenceRule().
		SetFrequency(FrequencyDaily).
		SetStartsAt("2024-01-01 00:00:00").
		SetEndsAt("2024-12-31 23:59:59").
		SetInterval(2).
		SetDaysOfWeek([]DayOfWeek{DayOfWeekMonday, DayOfWeekFriday}).
		SetDaysOfMonth([]int{1, 15}).
		SetMonthsOfYear([]MonthOfYear{MonthOfYearJanuary, MonthOfYearDecember})

	clone := original.(*recurrenceRule).Clone()

	// Verify clone has same values
	if clone.GetFrequency() != original.GetFrequency() {
		t.Errorf("Clone frequency mismatch")
	}
	if clone.GetStartsAt() != original.GetStartsAt() {
		t.Errorf("Clone startsAt mismatch")
	}
	if clone.GetEndsAt() != original.GetEndsAt() {
		t.Errorf("Clone endsAt mismatch")
	}
	if clone.GetInterval() != original.GetInterval() {
		t.Errorf("Clone interval mismatch")
	}

	// Verify clone is independent
	clone.SetFrequency(FrequencyHourly)
	if original.GetFrequency() == FrequencyHourly {
		t.Error("Modifying clone should not affect original")
	}
}

func TestRecurrenceRule_String(t *testing.T) {
	rule := NewRecurrenceRule().
		SetFrequency(FrequencyDaily).
		SetStartsAt("2024-01-01 00:00:00").
		SetInterval(1)

	str := rule.(*recurrenceRule).String()
	if str == "" {
		t.Error("String() should not return empty string")
	}
}

func TestRecurrenceRule_JSON(t *testing.T) {
	rule := NewRecurrenceRule().
		SetFrequency(FrequencyDaily).
		SetStartsAt("2024-01-01 00:00:00").
		SetEndsAt("2024-12-31 23:59:59").
		SetInterval(2).
		SetDaysOfWeek([]DayOfWeek{DayOfWeekMonday}).
		SetDaysOfMonth([]int{1, 15}).
		SetMonthsOfYear([]MonthOfYear{MonthOfYearJanuary})

	// Test MarshalJSON
	data, err := rule.(*recurrenceRule).MarshalJSON()
	if err != nil {
		t.Errorf("MarshalJSON() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("MarshalJSON() should not return empty data")
	}

	// Test UnmarshalJSON
	var newRule recurrenceRule
	err = newRule.UnmarshalJSON(data)
	if err != nil {
		t.Errorf("UnmarshalJSON() error = %v", err)
	}

	if newRule.GetFrequency() != rule.GetFrequency() {
		t.Error("UnmarshalJSON frequency mismatch")
	}
	if newRule.GetStartsAt() != rule.GetStartsAt() {
		t.Error("UnmarshalJSON startsAt mismatch")
	}
	if newRule.GetInterval() != rule.GetInterval() {
		t.Error("UnmarshalJSON interval mismatch")
	}
}

func TestNextRunAt(t *testing.T) {
	type testCase struct {
		name        string
		rule        RecurrenceRuleInterface
		now         *carbon.Carbon
		expected    *carbon.Carbon
		expectedErr string
	}

	testCases := []testCase{
		{
			name: "Starts in the future",
			rule: NewRecurrenceRule().
				SetFrequency(FrequencyDaily).
				SetStartsAt("2024-10-30T10:00:00Z").
				SetInterval(1),
			now:      carbon.Parse("2024-10-29T00:00:00Z", carbon.UTC),
			expected: carbon.Parse("2024-10-30T10:00:00Z", carbon.UTC),
		},
		{
			name: "Daily recurrence",
			rule: NewRecurrenceRule().
				SetFrequency(FrequencyDaily).
				SetStartsAt("2024-10-28T10:00:00Z").
				SetInterval(1),
			now:      carbon.Parse("2024-10-29T00:00:00Z", carbon.UTC),
			expected: carbon.Parse("2024-10-29T10:00:00Z", carbon.UTC),
		},
		{
			name: "Daily recurrence interval 2",
			rule: NewRecurrenceRule().
				SetFrequency(FrequencyDaily).
				SetStartsAt("2024-10-28T10:00:00Z").
				SetInterval(2),
			now:      carbon.Parse("2024-10-29T00:00:00Z", carbon.UTC),
			expected: carbon.Parse("2024-10-30T10:00:00Z", carbon.UTC),
		},
		{
			name: "Daily recurrence interval 3",
			rule: NewRecurrenceRule().
				SetFrequency(FrequencyDaily).
				SetStartsAt("2024-10-28T10:00:00Z").
				SetInterval(3),
			now:      carbon.Parse("2024-10-29T00:00:00Z", carbon.UTC),
			expected: carbon.Parse("2024-10-31T10:00:00Z", carbon.UTC),
		},
		{
			name: "Weekly recurrence - next week",
			rule: NewRecurrenceRule().
				SetFrequency(FrequencyWeekly).
				SetStartsAt("2024-10-28T10:00:00Z").
				SetInterval(1).
				SetDaysOfWeek([]DayOfWeek{DayOfWeekMonday}),
			now:      carbon.Parse("2024-10-31T00:00:00Z", carbon.UTC),
			expected: carbon.Parse("2024-11-04T10:00:00Z", carbon.UTC),
		},
		// {
		// 	name: "Ends at is before the next run - same day",
		// 	rule: NewRecurrenceRule().
		// 		SetFrequency(FrequencyWeekly).
		// 		SetStartsAt("2024-10-28T10:00:00Z").
		// 		SetEndsAt("2024-10-28T12:00:00Z").
		// 		SetInterval(1).
		// 		SetDaysOfWeek([]DayOfWeek{DayOfWeekMonday}),
		// 	now:      carbon.Parse("2024-10-28T11:00:00Z", carbon.UTC),
		// 	expected: carbon.Parse("2024-10-28T12:00:00Z", carbon.UTC),
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nextRun, err := NextRunAt(tc.rule, tc.now)
			if tc.expectedErr != "" {
				if err == nil {
					t.Errorf("Expected error containing %q, but got no error", tc.expectedErr)
				} else if !strings.Contains(err.Error(), tc.expectedErr) {
					t.Errorf("Expected error containing %q, but got %q", tc.expectedErr, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got %q", err.Error())
				} else if !nextRun.Eq(tc.expected) {
					t.Errorf("Expected %s, but got %s", tc.expected, nextRun)
				}
			}
		})
	}
}
