package editor_state_v1

// ContentSettings is the frozen v1 snapshot - see editorState.go. Never edit.
type ContentSettings struct {
	BannedItems        string       `json:"bannedItems"`
	BannedMagics       string       `json:"bannedMagics"`
	ValueOverridesText string       `json:"valueOverrides"`
	Bonuses            []BonusEntry `json:"bonuses"`

	PlayerZoneContentRows    []ZoneContentRow `json:"playerZoneContentRows,omitempty"`
	LowestNeutralContentRows []ZoneContentRow `json:"lowestNeutralContentRows,omitempty"`
	LowNeutralContentRows    []ZoneContentRow `json:"lowNeutralContentRows,omitempty"`
	MediumNeutralContentRows []ZoneContentRow `json:"mediumNeutralContentRows,omitempty"`
	HighNeutralContentRows   []ZoneContentRow `json:"highNeutralContentRows,omitempty"`
	HubZoneContentRows       []ZoneContentRow `json:"hubZoneContentRows,omitempty"`
}
