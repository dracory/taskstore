package taskstore

import (
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestNextRunAtQuick(t *testing.T) {
	rr := NewRecurrenceRule()
	now := carbon.Now(carbon.UTC)
	past := now.AddMinutes(-1).ToDateTimeString(carbon.UTC)
	rr.SetFrequency(FrequencyMinutely)
	rr.SetInterval(1)
	rr.SetStartsAt(past)

	t.Log("About to calculate next run")
	next, err := NextRunAt(rr, now)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	t.Logf("Next run: %s", next.ToDateTimeString())
}
