package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/models/regeneration"

// NextStateAction names the mutation the caller must apply to its pending
// ("next") snapshot after a regeneration decision. The decision service is
// pure, so it reports the mutation instead of performing it.
//
// The enum itself lives on the model side, because the service that produces it
// sits below the DTO line; these aliases name the identical type, so callers
// may keep spelling it dtos.NextStateLeave.
type NextStateAction = regeneration.NextStateAction

const (
	// NextStateLeave keeps the pending snapshot untouched.
	NextStateLeave = regeneration.NextStateLeave
	// NextStateClear drops the pending snapshot, cancelling any debounce.
	NextStateClear = regeneration.NextStateClear
	// NextStateSetFromCurrent re-arms the debounce against the live state.
	NextStateSetFromCurrent = regeneration.NextStateSetFromCurrent
)
