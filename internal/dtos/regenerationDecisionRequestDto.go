package dtos

import "time"

// RegenerationDecisionRequestDto is everything the regeneration decision
// depends on. Previous and Next are nil when absent, which is what
// distinguishes a first generation from a steady state.
//
// Now is the current frame time and DebounceDueAt is when the armed debounce
// window elapses; both are supplied by the caller so the decision stays
// deterministic and testable.
type RegenerationDecisionRequestDto struct {
	Previous      *EditorStateDto
	Current       *EditorStateDto
	Next          *EditorStateDto
	Now           time.Time
	DebounceDueAt time.Time
}
