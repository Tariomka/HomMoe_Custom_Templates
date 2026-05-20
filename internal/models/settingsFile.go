package models

import "github.com/Tariomka/hommoe_custom_templates/internal/models/generator"

// SettingsFile is the serialised .gen.json file produced and consumed by the
// editor
type SettingsFile struct {
	TemplateName                  string `json:"templateName"`
	MapSize                       int    `json:"mapSize"`
	PlayerCount                   int    `json:"playerCount"`
	NeutralZoneCount              int    `json:"neutralZoneCount"`
	PlayerZoneCastles             int    `json:"playerCastles"`
	NeutralZoneCastles            int    `json:"neutralCastles"`
	AdvancedMode                  bool   `json:"advancedMode"`
	NeutralLowNoCastleCount       int    `json:"neutralLowNoCastle"`
	NeutralLowCastleCount         int    `json:"neutralLowCastle"`
	NeutralMediumNoCastleCount    int    `json:"neutralMediumNoCastle"`
	NeutralMediumCastleCount      int    `json:"neutralMediumCastle"`
	NeutralHighNoCastleCount      int    `json:"neutralHighNoCastle"`
	NeutralHighCastleCount        int    `json:"neutralHighCastle"`
	MatchPlayerCastleFactions     bool   `json:"matchPlayerCastleFactions"`
	MinNeutralZonesBetweenPlayers int    `json:"minNeutralZonesBetweenPlayers"`
	// Deprecated: replaced by Topology == TopologyBalanced.
	// Retained only for one-way migration of older .gen.json files; the
	// loader copies the value into Topology and clears it again on save.
	ExperimentalBalancedZonePlacement bool                  `json:"experimentalBalancedZonePlacement,omitempty"`
	ExperimentalMapSizes              bool                  `json:"experimentalMapSizes"`
	PlayerZoneSize                    float64               `json:"playerZoneSize"`
	NeutralZoneSize                   float64               `json:"neutralZoneSize"`
	HubZoneSize                       float64               `json:"hubZoneSize"`
	HubZoneCastles                    int                   `json:"hubCastles"`
	GuardRandomization                float64               `json:"guardRandomization"`
	HeroCountMin                      int                   `json:"heroMin"`
	HeroCountMax                      int                   `json:"heroMax"`
	HeroCountIncrement                int                   `json:"heroIncrement"`
	Topology                          generator.MapTopology `json:"topology"`
	RandomPortals                     bool                  `json:"randomPortals"`
	MaxPortalConnections              int                   `json:"maxPortalConns"`
	SpawnRemoteFootholds              bool                  `json:"spawnFootholds"`
	GenerateRoads                     bool                  `json:"generateRoads"`
	NoDirectPlayerConn                bool                  `json:"isolateplayers"`
	ResourceDensityPercent            *int                  `json:"resourceDensity,omitempty"`
	StructureDensityPercent           *int                  `json:"structureDensity,omitempty"`
	NeutralStackStrengthPercent       int                   `json:"neutralStackStrength"`
	BorderGuardStrengthPercent        int                   `json:"borderGuardStrength"`
	VictoryCondition                  string                `json:"victoryCondition"`
	FactionLawsExpPercent             int                   `json:"factionLawsExp"`
	AstrologyExpPercent               int                   `json:"astrologyExp"`
	LostStartCity                     bool                  `json:"lostStartCity"`
	LostStartCityDay                  int                   `json:"lostStartCityDay"`
	LostStartHero                     bool                  `json:"lostStartHero"`
	CityHold                          bool                  `json:"cityHold"`
	CityHoldDays                      int                   `json:"cityHoldDays"`
	GladiatorArena                    bool                  `json:"gladiatorArena"`
	GladiatorArenaDaysDelayStart      int                   `json:"gladiatorArenaDaysDelayStart"`
	GladiatorArenaCountDay            int                   `json:"gladiatorArenaCountDay"`
	Tournament                        bool                  `json:"tournament"`
	TournamentFirstTournamentDay      int                   `json:"tournamentFirstTournamentDay"`
	TournamentInterval                int                   `json:"tournamentInterval"`
	TournamentPointsToWin             int                   `json:"tournamentPointsToWin"`
	TournamentSaveArmy                bool                  `json:"tournamentSaveArmy"`

	// ── Banned content / overrides / bonuses ─────────────────────────────
	BannedItems        string `json:"bannedItems"`
	BannedMagics       string `json:"bannedMagics"`
	ValueOverridesText string `json:"valueOverrides"`
	// BonusesJson stores configurable bonuses as a newline-separated list of
	// `BonusEntry.String()` lines (see ParseBonusesJson).
	BonusesJson string `json:"bonuses"`

	// ── Mandatory content rows per zone type ─────────────────────────────
	PlayerZoneContentRows    []ZoneContentRowSave `json:"playerZoneContentRows,omitempty"`
	LowNeutralContentRows    []ZoneContentRowSave `json:"lowNeutralContentRows,omitempty"`
	MediumNeutralContentRows []ZoneContentRowSave `json:"mediumNeutralContentRows,omitempty"`
	HighNeutralContentRows   []ZoneContentRowSave `json:"highNeutralContentRows,omitempty"`
	HubZoneContentRows       []ZoneContentRowSave `json:"hubZoneContentRows,omitempty"`

	// ContentDensityPercent is a legacy v0.2 setting; when present it seeds
	// both split density sliders on load.
	ContentDensityPercent *int `json:"contentDensity,omitempty"`
}

func NewSettingsFile() *SettingsFile {
	return &SettingsFile{
		TemplateName:                 "Custom Template",
		MapSize:                      160,
		PlayerCount:                  2,
		PlayerZoneCastles:            1,
		NeutralZoneCastles:           1,
		PlayerZoneSize:               1.0,
		NeutralZoneSize:              1.0,
		HubZoneSize:                  1.0,
		GuardRandomization:           0.05,
		HeroCountMin:                 4,
		HeroCountMax:                 8,
		HeroCountIncrement:           1,
		Topology:                     generator.TopologyBalanced,
		MaxPortalConnections:         32,
		SpawnRemoteFootholds:         true,
		GenerateRoads:                true,
		NeutralStackStrengthPercent:  100,
		BorderGuardStrengthPercent:   100,
		VictoryCondition:             "win_condition_1",
		FactionLawsExpPercent:        100,
		AstrologyExpPercent:          100,
		LostStartCityDay:             3,
		CityHoldDays:                 6,
		GladiatorArenaDaysDelayStart: 30,
		GladiatorArenaCountDay:       3,
		TournamentFirstTournamentDay: 14,
		TournamentInterval:           7,
		TournamentPointsToWin:        2,
		TournamentSaveArmy:           true,
	}
}

func (this *SettingsFile) EffectiveResourceDensity() int {
	if this.ResourceDensityPercent != nil {
		return *this.ResourceDensityPercent
	}
	if this.ContentDensityPercent != nil {
		return *this.ContentDensityPercent
	}
	return 100
}

func (this *SettingsFile) EffectiveStructureDensity() int {
	if this.StructureDensityPercent != nil {
		return *this.StructureDensityPercent
	}
	if this.ContentDensityPercent != nil {
		return *this.ContentDensityPercent
	}
	return 100
}
