package editor_state_model

// CastleSettingChanges reports which castle-count options differ between two
// editor states. These are the only generator options that override manual
// zone edits: when one changes after manual editing, the new count is pushed
// into the matching manually edited zones on the next regeneration.
type CastleSettingChanges struct {
	PlayerCastles bool // PlayerZoneCastles or PlayerOwnedCastles
	NeutralSimple bool // NeutralZoneCastles (simple mode: every neutral zone)
	NeutralLowest bool // NeutralLowestCastlesPerZone (advanced mode)
	NeutralLow    bool // NeutralLowCastlesPerZone (advanced mode)
	NeutralMedium bool // NeutralMediumCastlesPerZone (advanced mode)
	NeutralHigh   bool // NeutralHighCastlesPerZone (advanced mode)
	Hub           bool // HubZoneCastles
}

func (this CastleSettingChanges) Any() bool {
	return this.PlayerCastles || this.NeutralSimple || this.NeutralLowest ||
		this.NeutralLow || this.NeutralMedium || this.NeutralHigh || this.Hub
}
