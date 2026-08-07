package dtos

// NextStateAction names the mutation the caller must apply to its pending
// ("next") snapshot after a regeneration decision. The decision service is
// pure, so it reports the mutation instead of performing it.
type NextStateAction int

const (
	// NextStateLeave keeps the pending snapshot untouched.
	NextStateLeave NextStateAction = iota
	// NextStateClear drops the pending snapshot, cancelling any debounce.
	NextStateClear
	// NextStateSetFromCurrent re-arms the debounce against the live state.
	NextStateSetFromCurrent
)
