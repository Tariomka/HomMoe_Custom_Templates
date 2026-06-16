package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
)

// EditorStateDto is the serialized .gen.json file produced and consumed by the
// editor
type EditorStateDto struct {
	TemplateName                  string `json:"templateName"`
	GameMode                      string `json:"gameMode"`
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

	ExperimentalMapSizes         bool                     `json:"experimentalMapSizes"`
	PlayerZoneSize               float64                  `json:"playerZoneSize"`
	NeutralZoneSize              float64                  `json:"neutralZoneSize"`
	HubZoneSize                  float64                  `json:"hubZoneSize"`
	HubZoneCastles               int                      `json:"hubCastles"`
	GuardRandomization           float64                  `json:"guardRandomization"`
	HeroCountMin                 int                      `json:"heroMin"`
	HeroCountMax                 int                      `json:"heroMax"`
	HeroCountIncrement           int                      `json:"heroIncrement"`
	Topology                     config_inner.MapTopology `json:"topology"`
	RandomPortals                bool                     `json:"randomPortals"`
	MaxPortalConnections         int                      `json:"maxPortalConns"`
	SpawnRemoteFootholds         bool                     `json:"spawnFootholds"`
	GenerateRoads                bool                     `json:"generateRoads"`
	NoDirectPlayerConn           bool                     `json:"isolateplayers"`
	ResourceDensityPercent       int                      `json:"resourceDensity"`
	StructureDensityPercent      int                      `json:"structureDensity"`
	NeutralStackStrengthPercent  int                      `json:"neutralStackStrength"`
	BorderGuardStrengthPercent   int                      `json:"borderGuardStrength"`
	VictoryCondition             string                   `json:"victoryCondition"`
	FactionLawXpPercent          int                      `json:"factionLawsExp"`
	AstrologyXpPercent           int                      `json:"astrologyExp"`
	LostStartCity                bool                     `json:"lostStartCity"`
	LostStartCityDay             int                      `json:"lostStartCityDay"`
	LostStartHero                bool                     `json:"lostStartHero"`
	CityHold                     bool                     `json:"cityHold"`
	CityHoldDays                 int                      `json:"cityHoldDays"`
	GladiatorArena               bool                     `json:"gladiatorArena"`
	GladiatorArenaDaysDelayStart int                      `json:"gladiatorArenaDaysDelayStart"`
	GladiatorArenaCountDay       int                      `json:"gladiatorArenaCountDay"`
	Tournament                   bool                     `json:"tournament"`
	TournamentFirstTournamentDay int                      `json:"tournamentFirstTournamentDay"`
	TournamentInterval           int                      `json:"tournamentInterval"`
	TournamentPointsToWin        int                      `json:"tournamentPointsToWin"`
	TournamentSaveArmy           bool                     `json:"tournamentSaveArmy"`

	// ── Banned content / overrides / bonuses ─────────────────────────────
	BannedItems        string `json:"bannedItems"`
	BannedMagics       string `json:"bannedMagics"`
	ValueOverridesText string `json:"valueOverrides"`
	// BonusesJSON stores configurable bonuses as a newline-separated list of
	// `BonusEntry.String()` lines (see ParseBonusesJSON).
	BonusesJSON string `json:"bonuses"`

	// ── Mandatory content rows per zone type ─────────────────────────────
	PlayerZoneContentRows    []models.ZoneContentRowSave `json:"playerZoneContentRows,omitempty"`
	LowNeutralContentRows    []models.ZoneContentRowSave `json:"lowNeutralContentRows,omitempty"`
	MediumNeutralContentRows []models.ZoneContentRowSave `json:"mediumNeutralContentRows,omitempty"`
	HighNeutralContentRows   []models.ZoneContentRowSave `json:"highNeutralContentRows,omitempty"`
	HubZoneContentRows       []models.ZoneContentRowSave `json:"hubZoneContentRows,omitempty"`
}

func NewDefaultEditorStateDto() EditorStateDto {
	return EditorStateDto{
		TemplateName:                 "Custom Template",
		GameMode:                     "Classic",
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
		Topology:                     config_inner.TopologyCircles,
		MaxPortalConnections:         32,
		SpawnRemoteFootholds:         true,
		GenerateRoads:                true,
		ResourceDensityPercent:       100,
		StructureDensityPercent:      100,
		NeutralStackStrengthPercent:  100,
		BorderGuardStrengthPercent:   100,
		VictoryCondition:             "win_condition_1",
		FactionLawXpPercent:          100,
		AstrologyXpPercent:           100,
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
