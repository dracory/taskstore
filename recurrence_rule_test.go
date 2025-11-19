package taskstore

import (
	"strings"
	"testing"

	"github.com/dromara/carbon/v2"
)

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
