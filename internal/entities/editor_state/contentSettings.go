package editor_state

const CurrentContentDefaultsVersion = 1 // Names the default player content rows this build omits from saved files.

type ContentSettings struct {
	ContentDefaultsVersion int `json:"contentDefaultsVersion"`

	BannedItems    []BanEntry           `json:"bannedItems"`
	BannedMagics   []BanEntry           `json:"bannedMagics"`
	ValueOverrides []ValueOverrideEntry `json:"valueOverrides"`
	Bonuses        []BonusEntry         `json:"bonuses"`

	PlayerZoneContentRows    []ZoneContentRow `json:"playerZoneContentRows,omitzero"`
	LowestNeutralContentRows []ZoneContentRow `json:"lowestNeutralContentRows,omitzero"`
	LowNeutralContentRows    []ZoneContentRow `json:"lowNeutralContentRows,omitzero"`
	MediumNeutralContentRows []ZoneContentRow `json:"mediumNeutralContentRows,omitzero"`
	HighNeutralContentRows   []ZoneContentRow `json:"highNeutralContentRows,omitzero"`
	HubZoneContentRows       []ZoneContentRow `json:"hubZoneContentRows,omitzero"`
}
