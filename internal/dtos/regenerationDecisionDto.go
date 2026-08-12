package dtos

import "time"

// RegenerationDecisionDto is the outcome of a regeneration decision: whether to
// regenerate now, what to do with the pending snapshot, and when the caller
// should wake up again.
//
// RedrawAt is only meaningful while ScheduleRedraw is true.
type RegenerationDecisionDto struct {
	Regenerate      bool
	NextStateAction NextStateAction
	RedrawAt        time.Time
	ScheduleRedraw  bool
}
