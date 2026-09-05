package editor_state

// ContentSettings holds the global bans, value overrides, bonuses and the
// mandatory content rows seeded into each zone type.
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
