package editor_state

// BonusEntry is the editor-side record for a single configurable game-start bonus.
// Its hash lives in internal/helpers/config_helpers.
type BonusEntry struct {
	PresetType     BonusPresetType `json:"presetType"`
	ReceiverFilter string          `json:"receiverFilter"`   // "start_hero" or "all_heroes".
	Param          string          `json:"param,omitempty"`  // Spell sid / item sid / numeric value depending on type.
	Param2         string          `json:"param2,omitempty"` // For Spell: "1" = free, "0" = normal. Unused for other types.
}
