package taskstore

import (
	"time"

	"github.com/dromara/carbon/v2"
)

// parseTime converts a datetime string to time.Time.
// NULL_DATETIME and empty strings are converted to a zero time.Time.
func parseTime(s string) time.Time {
	if s == "" || s == NULL_DATETIME {
		return time.Time{}
	}
	return carbon.Parse(s, carbon.UTC).StdTime()
}
