package taskstore

// type ScheduleDefinition interface {
// 	GetID() string
// 	SetID(string) ScheduleDefinition

// 	GetRecurrenceRule() string
// 	SetRecurrenceRule(dateTimeUtc string) ScheduleDefinition

// 	GetStartsAt() string
// 	SetStartsAt(dateTimeUtc string)

// 	GetEndsAt() string
// 	SetEndsAt(string)

// 	IsValid() bool
// 	GetNextRunTime(string) (string, error)
// }

// type ScheduleInterface interface {
// 	ScheduleDefinitionID() string
// 	SetScheduleDefinition(ScheduleDefinition)
// }

// type RecurrenceRule struct {
// 	// Required
// 	Frequency string

// 	// UNTIL and COUNT can only have 1
// 	Until DateTimeUTC // nil means forever
// 	Count int64       // 0 means non

// 	// BYX
// 	// ByDay ByDayList

// 	// Optional
// 	Interval int
// }

// // Clone will deep copy current rrule reference
// func (rrule *RecurrenceRule) Clone() *RecurrenceRule {
// 	if rrule == nil {
// 		return nil
// 	}
// 	r := *rrule
// 	return &r
// }
