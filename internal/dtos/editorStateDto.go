package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

var (
	winConditions = registry.GetWinningConditionValues()
	gameModes     = registry.GetGameModeValues()
)

// EditorStateDto is the serialized .gen.json file produced and consumed by the
// editor
type EditorStateDto struct {
	TemplateName                  string `json:"templateName"`
	GameMode                      string `json:"gameMode"`
	MapSize                       int    `json:"mapSize"`
	PlayerCount                   int    `json:"playerCount"`
	NeutralZoneCount              int    `json:"neutralZoneCount"`
	PlayerOwnedCastles            int    `json:"playerOwnedCastles"`
	PlayerZoneCastles             int    `json:"playerCastles"`
	NeutralZoneCastles            int    `json:"neutralCastles"`
	SpawnAbandonedOutposts        bool   `json:"spawnAbandonedOutposts"`
	AbandonedOutpostCount         int    `json:"abandonedOutpostCount"`
	AdvancedMode                  bool   `json:"advancedMode"`
	NeutralLowNoCastleCount       int    `json:"neutralLowNoCastle"`
	NeutralLowCastleCount         int    `json:"neutralLowCastle"`
	NeutralMediumNoCastleCount    int    `json:"neutralMediumNoCastle"`
	NeutralMediumCastleCount      int    `json:"neutralMediumCastle"`
	NeutralHighNoCastleCount      int    `json:"neutralHighNoCastle"`
	NeutralHighCastleCount        int    `json:"neutralHighCastle"`
	NeutralLowCastlesPerZone      int    `json:"neutralLowCastlesPerZone"`
	NeutralMediumCastlesPerZone   int    `json:"neutralMedCastlesPerZone"`
	NeutralHighCastlesPerZone     int    `json:"neutralHighCastlesPerZone"`
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
	RemoteFootholdCount          int                      `json:"remoteFootholdCount"`
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
		GameMode:                     gameModes.Classic,
		MapSize:                      160,
		PlayerCount:                  2,
		NeutralZoneCastles:           1,
		NeutralLowCastlesPerZone:     1,
		NeutralMediumCastlesPerZone:  1,
		NeutralHighCastlesPerZone:    1,
		MatchPlayerCastleFactions:    true,
		PlayerZoneSize:               1.0,
		NeutralZoneSize:              1.0,
		HubZoneSize:                  1.0,
		GuardRandomization:           0.05,
		HeroCountMin:                 4,
		HeroCountMax:                 8,
		HeroCountIncrement:           1,
		Topology:                     config_inner.TopologyRandom,
		MaxPortalConnections:         32,
		SpawnRemoteFootholds:         true,
		RemoteFootholdCount:          1,
		AbandonedOutpostCount:        1,
		GenerateRoads:                true,
		ResourceDensityPercent:       100,
		StructureDensityPercent:      100,
		NeutralStackStrengthPercent:  100,
		BorderGuardStrengthPercent:   100,
		VictoryCondition:             winConditions.Standard,
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
		PlayerZoneContentRows:        DefaultPlayerZoneContentRows(),
	}
}

// DefaultPlayerZoneContentRows returns the historical default mandatory-content
// rows seeded into every player zone: the six basic mines plus an alchemy lab,
// a couple of guarded treasures, random hires and resource banks.
func DefaultPlayerZoneContentRows() []models.ZoneContentRowSave {
	interactable := registry.GetMapObjectAllInteractableValues()
	resources := registry.GetMapObjectResourceValues()
	randomItems := registry.GetMapObjectRandomItemValues()
	randomHires := registry.GetMandatoryContentRandomHiresBuildingValues()
	basicRandomHires := registry.GetMandatoryContentBasicRandomHiresBuildingValues()
	basicResourceBanks := registry.GetMandatoryContentBasicResourceBanksBuildingValues()

	trueVal := true
	return []models.ZoneContentRowSave{
		{
			Sid:    interactable.WoodMine,
			Count:  1,
			IsMine: true,
			Rules: []models.ContentRuleRowSave{
				{Name: "Guarded", IsGuarded: &trueVal},
				{Name: "Distance to town", DistanceName: "Near"},
			},
		},
		{
			Sid:    interactable.OreMine,
			Count:  1,
			IsMine: true,
			Rules: []models.ContentRuleRowSave{
				{Name: "Guarded", IsGuarded: &trueVal},
				{Name: "Distance to town", DistanceName: "Near"},
			},
		},
		{
			Sid:    interactable.GoldMine,
			Count:  1,
			IsMine: true,
			Rules: []models.ContentRuleRowSave{
				{Name: "Guarded", IsGuarded: &trueVal},
				{Name: "Distance to town", DistanceName: "Near"},
			},
		},
		{
			Sid:    interactable.CrystalMine,
			Count:  1,
			IsMine: true,
			Rules: []models.ContentRuleRowSave{
				{Name: "Guarded", IsGuarded: &trueVal},
				{Name: "Distance to road", DistanceName: "Next To"},
			},
		},
		{
			Sid:    interactable.MercuryMine,
			Count:  1,
			IsMine: true,
			Rules: []models.ContentRuleRowSave{
				{Name: "Guarded", IsGuarded: &trueVal},
				{Name: "Distance to road", DistanceName: "Next To"},
			},
		},
		{
			Sid:    interactable.GemstoneMine,
			Count:  1,
			IsMine: true,
			Rules: []models.ContentRuleRowSave{
				{Name: "Guarded", IsGuarded: &trueVal},
				{Name: "Distance to road", DistanceName: "Next To"},
			},
		},
		{
			Sid:    interactable.AlchemyLab,
			Count:  1,
			IsMine: true,
			Rules: []models.ContentRuleRowSave{
				{Name: "Guarded", IsGuarded: &trueVal},
				{Name: "Distance to road", DistanceName: "Next To"},
			},
		},
		{
			Sid:   resources.PandoraBox,
			Count: 1,
			Rules: []models.ContentRuleRowSave{{Name: "Guarded", IsGuarded: &trueVal}},
		},
		{
			Sid:   randomItems.RandomItemEpic,
			Count: 1,
			Rules: []models.ContentRuleRowSave{{Name: "Guarded", IsGuarded: &trueVal}},
		},
		{
			Sid:     randomHires.RandomHiresLowTier,
			Count:   2,
			IsGroup: true,
			Rules:   []models.ContentRuleRowSave{{Name: "Guarded", IsGuarded: &trueVal}},
		},
		{
			Sid:     randomHires.RandomHiresHighTier,
			Count:   1,
			IsGroup: true,
			Rules:   []models.ContentRuleRowSave{{Name: "Guarded", IsGuarded: &trueVal}},
		},
		{
			Sid:     basicRandomHires.BasicRandomHires,
			Count:   1,
			IsGroup: true,
			Rules:   []models.ContentRuleRowSave{{Name: "Guarded", IsGuarded: &trueVal}},
		},
		{
			Sid:     basicResourceBanks.BasicResourceBanksTier1,
			Count:   2,
			IsGroup: true,
			Rules:   []models.ContentRuleRowSave{{Name: "Guarded", IsGuarded: &trueVal}},
		},
		{
			Sid:     basicResourceBanks.BasicResourceBanksTier2,
			Count:   1,
			IsGroup: true,
			Rules:   []models.ContentRuleRowSave{{Name: "Guarded", IsGuarded: &trueVal}},
		},
	}
}

// LayoutDefiningOptionsChanged reports whether any option that changes the set
// of zones or the connection graph differs between two editor states. When
// these are unchanged, manual zone edits remain valid and can be reapplied.
func (this *EditorStateDto) LayoutDefiningOptionsChanged(incoming *EditorStateDto) bool {
	return this.PlayerCount != incoming.PlayerCount ||
		this.Topology != incoming.Topology ||
		this.GenerateRoads != incoming.GenerateRoads ||
		this.RandomPortals != incoming.RandomPortals ||
		this.NoDirectPlayerConn != incoming.NoDirectPlayerConn ||
		this.MaxPortalConnections != incoming.MaxPortalConnections ||
		this.MinNeutralZonesBetweenPlayers != incoming.MinNeutralZonesBetweenPlayers ||
		this.SpawnRemoteFootholds != incoming.SpawnRemoteFootholds ||
		this.RemoteFootholdCount != incoming.RemoteFootholdCount ||
		this.SpawnAbandonedOutposts != incoming.SpawnAbandonedOutposts ||
		this.AbandonedOutpostCount != incoming.AbandonedOutpostCount ||
		this.zoneCountOptionsChanged(incoming)
}

// zoneCountOptionsChanged reports whether the number of neutral zones differs
// between two editor states.
func (this *EditorStateDto) zoneCountOptionsChanged(incoming *EditorStateDto) bool {
	return this.AdvancedMode != incoming.AdvancedMode ||
		this.NeutralZoneCount != incoming.NeutralZoneCount ||
		this.NeutralLowNoCastleCount != incoming.NeutralLowNoCastleCount ||
		this.NeutralLowCastleCount != incoming.NeutralLowCastleCount ||
		this.NeutralMediumNoCastleCount != incoming.NeutralMediumNoCastleCount ||
		this.NeutralMediumCastleCount != incoming.NeutralMediumCastleCount ||
		this.NeutralHighNoCastleCount != incoming.NeutralHighNoCastleCount ||
		this.NeutralHighCastleCount != incoming.NeutralHighCastleCount
}
