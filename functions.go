package taskstore

import (
	"time"

	"github.com/dromara/carbon/v2"
)

// parseTime converts a datetime string to time.Time.
// NULL_DATETIME and empty strings are converted to a zero time.Time.
func parseTime(s string) time.Time {
	if isNullDateTime(s) {
		return time.Time{}
	}
	return carbon.Parse(s, carbon.UTC).StdTime()
}

// isNullDateTime returns true if the given string represents a NULL datetime.
//
// Business Logic:
//  1. NULL_DATETIME means "not yet initialized" — the schedule has never
//     had its NextRunAt calculated (e.g. freshly created, or loaded from
//     DB before the runner processed it).
//  2. SQLite DateTime columns reformat "0002-01-01 00:00:00" to
//     "0002-01-01T00:00:00Z" on read-back, so a raw string comparison
//     against the NULL_DATETIME constant fails after a DB round-trip.
//  3. This helper parses the value and compares against the NULL_DATETIME
//     epoch, making it robust against any separator/timezone formatting
//     the DB layer may apply.
func isNullDateTime(s string) bool {
	if s == "" {
		return true
	}
	// Parse and compare against the NULL_DATETIME epoch; this is robust
	// against any whitespace/separator/timezone formatting the DB layer
	// may apply during a round-trip.
	t := carbon.Parse(s, carbon.UTC)
	if t.IsZero() || t.IsInvalid() {
		return true
	}
	nullEpoch := carbon.Parse(NULL_DATETIME, carbon.UTC)
	return t.Eq(nullEpoch)
}

// isMaxDateTime returns true if the given string represents MAX_DATETIME.
//
// Business Logic:
//  1. MAX_DATETIME means "exhausted" — the recurrence rule returned
//     ErrNoMoreRuns, so there are no future occurrences.
//  2. MAX_DATETIME is reused as the exhaustion sentinel because it already
//     means "infinity" elsewhere in the codebase, and it must not collide
//     with NULL_DATETIME ("not yet initialized") — the two states require
//     opposite runner actions (init vs. mark completed).
//  3. Like isNullDateTime, this helper parses and compares against the
//     MAX_DATETIME epoch, making it robust against DB round-trip
//     reformatting.
func isMaxDateTime(s string) bool {
	if s == "" {
		return false
	}
	t := carbon.Parse(s, carbon.UTC)
	if t.IsZero() || t.IsInvalid() {
		return false
	}
	maxEpoch := carbon.Parse(MAX_DATETIME, carbon.UTC)
	return t.Eq(maxEpoch)
}
