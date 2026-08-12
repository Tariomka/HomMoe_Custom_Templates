package dtos

// ExistingBonusesDto summarises an already-composed bonus list: the hash set
// used for duplicate detection and the spell IDs to pre-exclude in the picker.
type ExistingBonusesDto struct {
	Keys     map[string]bool
	SpellIDs []string
}
