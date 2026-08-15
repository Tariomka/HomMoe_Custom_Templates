package dtos

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

const guardedRuleName = "Guarded"

// EditorStateDto is the serialized .gen.json file produced and consumed by the editor.
type EditorStateDto struct {
	TemplateName                string `json:"templateName"`
	GameMode                    string `json:"gameMode"`
	MapSize                     int    `json:"mapSize"`
	PlayerCount                 int    `json:"playerCount"`
	NeutralZoneCount            int    `json:"neutralZoneCount"`
	PlayerOwnedCastles          int    `json:"playerOwnedCastles"`
	PlayerZoneCastles           int    `json:"playerCastles"`
	NeutralZoneCastles          int    `json:"neutralCastles"`
	SpawnAbandonedOutposts      bool   `json:"spawnAbandonedOutposts"`
	AbandonedOutpostCount       int    `json:"abandonedOutpostCount"`
	AdvancedMode                bool   `json:"advancedMode"`
	NeutralLowestNoCastleCount  int    `json:"neutralLowestNoCastle"`
	NeutralLowestCastleCount    int    `json:"neutralLowestCastle"`
	NeutralLowNoCastleCount     int    `json:"neutralLowNoCastle"`
	NeutralLowCastleCount       int    `json:"neutralLowCastle"`
	NeutralMediumNoCastleCount  int    `json:"neutralMediumNoCastle"`
	NeutralMediumCastleCount    int    `json:"neutralMediumCastle"`
	NeutralHighNoCastleCount    int    `json:"neutralHighNoCastle"`
	NeutralHighCastleCount      int    `json:"neutralHighCastle"`
	NeutralLowestCastlesPerZone int    `json:"neutralLowestCastlesPerZone"`
	NeutralLowCastlesPerZone    int    `json:"neutralLowCastlesPerZone"`
	NeutralMediumCastlesPerZone int    `json:"neutralMedCastlesPerZone"`
	NeutralHighCastlesPerZone   int    `json:"neutralHighCastlesPerZone"`
	MatchPlayerCastleFactions   bool   `json:"matchPlayerCastleFactions"`

	ExperimentalMapSizes         bool               `json:"experimentalMapSizes"`
	PlayerZoneSize               float64            `json:"playerZoneSize"`
	NeutralZoneSize              float64            `json:"neutralZoneSize"`
	HubZoneSize                  float64            `json:"hubZoneSize"`
	HubZoneCastles               int                `json:"hubCastles"`
	GuardRandomization           float64            `json:"guardRandomization"`
	HeroCountMin                 int                `json:"heroMin"`
	HeroCountMax                 int                `json:"heroMax"`
	HeroCountIncrement           int                `json:"heroIncrement"`
	Topology                     config.MapTopology `json:"topology"`
	RandomPortals                bool               `json:"randomPortals"`
	MaxPortalConnections         int                `json:"maxPortalConns"`
	SpawnRemoteFootholds         bool               `json:"spawnFootholds"`
	RemoteFootholdCount          int                `json:"remoteFootholdCount"`
	GenerateRoads                bool               `json:"generateRoads"`
	NoDirectPlayerConn           bool               `json:"isolateplayers"`
	ResourceDensityPercent       int                `json:"resourceDensity"`
	StructureDensityPercent      int                `json:"structureDensity"`
	NeutralStackStrengthPercent  int                `json:"neutralStackStrength"`
	BorderGuardStrengthPercent   int                `json:"borderGuardStrength"`
	VictoryCondition             string             `json:"victoryCondition"`
	FactionLawXpPercent          int                `json:"factionLawsExp"`
	AstrologyXpPercent           int                `json:"astrologyExp"`
	LostStartCity                bool               `json:"lostStartCity"`
	LostStartCityDay             int                `json:"lostStartCityDay"`
	LostStartHero                bool               `json:"lostStartHero"`
	CityHold                     bool               `json:"cityHold"`
	CityHoldDays                 int                `json:"cityHoldDays"`
	GladiatorArena               bool               `json:"gladiatorArena"`
	GladiatorArenaDaysDelayStart int                `json:"gladiatorArenaDaysDelayStart"`
	GladiatorArenaCountDay       int                `json:"gladiatorArenaCountDay"`
	Tournament                   bool               `json:"tournament"`
	TournamentFirstTournamentDay int                `json:"tournamentFirstTournamentDay"`
	TournamentInterval           int                `json:"tournamentInterval"`
	TournamentPointsToWin        int                `json:"tournamentPointsToWin"`
	TournamentSaveArmy           bool               `json:"tournamentSaveArmy"`

	// ── Banned content / overrides / bonuses ─────────────────────────────
	BannedItems        string              `json:"bannedItems"`
	BannedMagics       string              `json:"bannedMagics"`
	ValueOverridesText string              `json:"valueOverrides"`
	Bonuses            []config.BonusEntry `json:"bonuses"`

	// ── Mandatory content rows per zone type ─────────────────────────────
	PlayerZoneContentRows    []models.ZoneContentRowSave `json:"playerZoneContentRows,omitempty"`
	LowestNeutralContentRows []models.ZoneContentRowSave `json:"lowestNeutralContentRows,omitempty"`
	LowNeutralContentRows    []models.ZoneContentRowSave `json:"lowNeutralContentRows,omitempty"`
	MediumNeutralContentRows []models.ZoneContentRowSave `json:"mediumNeutralContentRows,omitempty"`
	HighNeutralContentRows   []models.ZoneContentRowSave `json:"highNeutralContentRows,omitempty"`
	HubZoneContentRows       []models.ZoneContentRowSave `json:"hubZoneContentRows,omitempty"`

	// ── Manual zone editor edits ─────────────────────────────────────────
	ManualZones       []editor_state_dto.ManualZoneSave       `json:"manualZones,omitempty"`
	ManualConnections []editor_state_dto.ManualConnectionSave `json:"manualConnections,omitempty"`
}

func NewDefaultEditorStateDto() EditorStateDto {
	winConditions := registry.GetWinningConditionValues()
	gameModes := registry.GetGameModeValues()
	return EditorStateDto{
		TemplateName:                 "Custom Template",
		GameMode:                     gameModes.Classic,
		MapSize:                      160,
		PlayerCount:                  2,
		NeutralZoneCastles:           1,
		NeutralLowestCastlesPerZone:  1,
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
		Topology:                     config.TopologyRandom,
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

// Clone returns a copy that shares no backing array or pointer with the
// receiver, so an in-place edit of any element on either side stays invisible
// to the other. A plain struct copy duplicates slice headers only, which would
// hide element mutations from the change detection in
// EqualsIgnoringManualEdits.
func (this *EditorStateDto) Clone() EditorStateDto {
	clone := *this

	clone.Bonuses = slices.Clone(this.Bonuses)
	clone.PlayerZoneContentRows = cloneContentRows(this.PlayerZoneContentRows)
	clone.LowestNeutralContentRows = cloneContentRows(this.LowestNeutralContentRows)
	clone.LowNeutralContentRows = cloneContentRows(this.LowNeutralContentRows)
	clone.MediumNeutralContentRows = cloneContentRows(this.MediumNeutralContentRows)
	clone.HighNeutralContentRows = cloneContentRows(this.HighNeutralContentRows)
	clone.HubZoneContentRows = cloneContentRows(this.HubZoneContentRows)

	clone.ManualZones = slices.Clone(this.ManualZones)
	for zoneIndex := range clone.ManualZones {
		clone.ManualZones[zoneIndex] = clone.ManualZones[zoneIndex].Clone()
	}

	clone.ManualConnections = slices.Clone(this.ManualConnections)
	for connectionIndex := range clone.ManualConnections {
		clone.ManualConnections[connectionIndex] = clone.ManualConnections[connectionIndex].Clone()
	}

	return clone
}

// cloneContentRows deep-copies a content-row slice, preserving a nil slice as
// nil because the change detection distinguishes nil from empty.
func cloneContentRows(rows []models.ZoneContentRowSave) []models.ZoneContentRowSave {
	clone := slices.Clone(rows)
	for rowIndex := range clone {
		clone[rowIndex] = clone[rowIndex].Clone()
	}
	return clone
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
		this.zoneCountOptionsChanged(incoming)
}

// DiffCastleSettings compares the castle-count options of this state (the one
// behind the last generation) against the incoming current state. AdvancedMode
// gates which neutral options are relevant; it cannot flip between the two
// states here because such a flip is layout-defining and discards manual edits
// before castle propagation is ever considered.
func (this *EditorStateDto) DiffCastleSettings(incoming *EditorStateDto) editor_state_dto.CastleSettingChanges {
	changes := editor_state_dto.CastleSettingChanges{
		PlayerCastles: this.PlayerZoneCastles != incoming.PlayerZoneCastles ||
			this.PlayerOwnedCastles != incoming.PlayerOwnedCastles,
		Hub: this.HubZoneCastles != incoming.HubZoneCastles,
	}
	if incoming.AdvancedMode {
		changes.NeutralLowest = this.NeutralLowestCastlesPerZone != incoming.NeutralLowestCastlesPerZone
		changes.NeutralLow = this.NeutralLowCastlesPerZone != incoming.NeutralLowCastlesPerZone
		changes.NeutralMedium = this.NeutralMediumCastlesPerZone != incoming.NeutralMediumCastlesPerZone
		changes.NeutralHigh = this.NeutralHighCastlesPerZone != incoming.NeutralHighCastlesPerZone
	} else {
		changes.NeutralSimple = this.NeutralZoneCastles != incoming.NeutralZoneCastles
	}
	return changes
}

// EqualsIgnoringManualEdits reports whether two editor states are equal when
// the manual-edit fields are disregarded. Manual zones and connections are
// reapplied to the generated template through a separate path, so they must
// not trigger an automatic regeneration on their own.
//
// The comparison is hand-rolled instead of [reflect.DeepEqual] because it runs
// on the UI hot path several times per frame. Every non-manual field must be
// covered here; the per-field mutation test on this method trips when a new
// DTO field is added without extending the comparison.
func (this *EditorStateDto) EqualsIgnoringManualEdits(other *EditorStateDto) bool {
	return this.zoneOptionScalarsEqual(other) &&
		this.generationOptionScalarsEqual(other) &&
		this.gameRuleScalarsEqual(other) &&
		(this.Bonuses == nil) == (other.Bonuses == nil) &&
		slices.Equal(this.Bonuses, other.Bonuses) &&
		contentRowSlicesEqual(this.PlayerZoneContentRows, other.PlayerZoneContentRows) &&
		contentRowSlicesEqual(this.LowestNeutralContentRows, other.LowestNeutralContentRows) &&
		contentRowSlicesEqual(this.LowNeutralContentRows, other.LowNeutralContentRows) &&
		contentRowSlicesEqual(this.MediumNeutralContentRows, other.MediumNeutralContentRows) &&
		contentRowSlicesEqual(this.HighNeutralContentRows, other.HighNeutralContentRows) &&
		contentRowSlicesEqual(this.HubZoneContentRows, other.HubZoneContentRows)
}

func (this *EditorStateDto) HasManualEdits() bool {
	return len(this.ManualZones) > 0 || len(this.ManualConnections) > 0
}

// zoneCountOptionsChanged reports whether the number of neutral zones differs
// between two editor states.
func (this *EditorStateDto) zoneCountOptionsChanged(incoming *EditorStateDto) bool {
	return this.AdvancedMode != incoming.AdvancedMode ||
		this.NeutralZoneCount != incoming.NeutralZoneCount ||
		this.NeutralLowestNoCastleCount != incoming.NeutralLowestNoCastleCount ||
		this.NeutralLowestCastleCount != incoming.NeutralLowestCastleCount ||
		this.NeutralLowNoCastleCount != incoming.NeutralLowNoCastleCount ||
		this.NeutralLowCastleCount != incoming.NeutralLowCastleCount ||
		this.NeutralMediumNoCastleCount != incoming.NeutralMediumNoCastleCount ||
		this.NeutralMediumCastleCount != incoming.NeutralMediumCastleCount ||
		this.NeutralHighNoCastleCount != incoming.NeutralHighNoCastleCount ||
		this.NeutralHighCastleCount != incoming.NeutralHighCastleCount
}

// zoneOptionScalarsEqual compares the template identity and zone/castle count
// options (first section of the struct).
//
//nolint:dupl // comparison chains over disjoint field groups are inherently similar
func (this *EditorStateDto) zoneOptionScalarsEqual(other *EditorStateDto) bool {
	return this.TemplateName == other.TemplateName &&
		this.GameMode == other.GameMode &&
		this.MapSize == other.MapSize &&
		this.PlayerCount == other.PlayerCount &&
		this.NeutralZoneCount == other.NeutralZoneCount &&
		this.PlayerOwnedCastles == other.PlayerOwnedCastles &&
		this.PlayerZoneCastles == other.PlayerZoneCastles &&
		this.NeutralZoneCastles == other.NeutralZoneCastles &&
		this.SpawnAbandonedOutposts == other.SpawnAbandonedOutposts &&
		this.AbandonedOutpostCount == other.AbandonedOutpostCount &&
		this.AdvancedMode == other.AdvancedMode &&
		this.NeutralLowestNoCastleCount == other.NeutralLowestNoCastleCount &&
		this.NeutralLowestCastleCount == other.NeutralLowestCastleCount &&
		this.NeutralLowNoCastleCount == other.NeutralLowNoCastleCount &&
		this.NeutralLowCastleCount == other.NeutralLowCastleCount &&
		this.NeutralMediumNoCastleCount == other.NeutralMediumNoCastleCount &&
		this.NeutralMediumCastleCount == other.NeutralMediumCastleCount &&
		this.NeutralHighNoCastleCount == other.NeutralHighNoCastleCount &&
		this.NeutralHighCastleCount == other.NeutralHighCastleCount &&
		this.NeutralLowestCastlesPerZone == other.NeutralLowestCastlesPerZone &&
		this.NeutralLowCastlesPerZone == other.NeutralLowCastlesPerZone &&
		this.NeutralMediumCastlesPerZone == other.NeutralMediumCastlesPerZone &&
		this.NeutralHighCastlesPerZone == other.NeutralHighCastlesPerZone &&
		this.MatchPlayerCastleFactions == other.MatchPlayerCastleFactions
}

// generationOptionScalarsEqual compares the map generation options (second
// section of the struct: sizes, heroes, topology, connections, densities).
//
//nolint:dupl // comparison chains over disjoint field groups are inherently similar
func (this *EditorStateDto) generationOptionScalarsEqual(other *EditorStateDto) bool {
	return this.ExperimentalMapSizes == other.ExperimentalMapSizes &&
		this.PlayerZoneSize == other.PlayerZoneSize &&
		this.NeutralZoneSize == other.NeutralZoneSize &&
		this.HubZoneSize == other.HubZoneSize &&
		this.HubZoneCastles == other.HubZoneCastles &&
		this.GuardRandomization == other.GuardRandomization &&
		this.HeroCountMin == other.HeroCountMin &&
		this.HeroCountMax == other.HeroCountMax &&
		this.HeroCountIncrement == other.HeroCountIncrement &&
		this.Topology == other.Topology &&
		this.RandomPortals == other.RandomPortals &&
		this.MaxPortalConnections == other.MaxPortalConnections &&
		this.SpawnRemoteFootholds == other.SpawnRemoteFootholds &&
		this.RemoteFootholdCount == other.RemoteFootholdCount &&
		this.GenerateRoads == other.GenerateRoads &&
		this.NoDirectPlayerConn == other.NoDirectPlayerConn &&
		this.ResourceDensityPercent == other.ResourceDensityPercent &&
		this.StructureDensityPercent == other.StructureDensityPercent &&
		this.NeutralStackStrengthPercent == other.NeutralStackStrengthPercent &&
		this.BorderGuardStrengthPercent == other.BorderGuardStrengthPercent
}

// gameRuleScalarsEqual compares the victory/game-rule options and the banned
// content free-text fields (third section of the struct).
//
//nolint:dupl // comparison chains over disjoint field groups are inherently similar
func (this *EditorStateDto) gameRuleScalarsEqual(other *EditorStateDto) bool {
	return this.VictoryCondition == other.VictoryCondition &&
		this.FactionLawXpPercent == other.FactionLawXpPercent &&
		this.AstrologyXpPercent == other.AstrologyXpPercent &&
		this.LostStartCity == other.LostStartCity &&
		this.LostStartCityDay == other.LostStartCityDay &&
		this.LostStartHero == other.LostStartHero &&
		this.CityHold == other.CityHold &&
		this.CityHoldDays == other.CityHoldDays &&
		this.GladiatorArena == other.GladiatorArena &&
		this.GladiatorArenaDaysDelayStart == other.GladiatorArenaDaysDelayStart &&
		this.GladiatorArenaCountDay == other.GladiatorArenaCountDay &&
		this.Tournament == other.Tournament &&
		this.TournamentFirstTournamentDay == other.TournamentFirstTournamentDay &&
		this.TournamentInterval == other.TournamentInterval &&
		this.TournamentPointsToWin == other.TournamentPointsToWin &&
		this.TournamentSaveArmy == other.TournamentSaveArmy &&
		this.BannedItems == other.BannedItems &&
		this.BannedMagics == other.BannedMagics &&
		this.ValueOverridesText == other.ValueOverridesText
}

// DefaultPlayerZoneContentRows returns the historical default mandatory-content
// rows seeded into every player zone: the six basic mines plus an alchemy lab,
// a couple of guarded treasures, random hires and resource banks.
func DefaultPlayerZoneContentRows() []models.ZoneContentRowSave {
	rows := defaultPlayerZoneMineRows()
	rows = append(rows, defaultPlayerZoneTreasureRows()...)
	rows = append(rows, defaultPlayerZoneHireAndBankRows()...)
	return rows
}

// defaultPlayerZoneMineRows returns the six basic mines plus an alchemy lab:
// the wood, ore and gold mines near the town, the rest next to a road.
func defaultPlayerZoneMineRows() []models.ZoneContentRowSave {
	interactable := registry.GetMapObjectAllInteractableValues()
	nearTown := models.ContentRuleRowSave{Name: "Distance to town", DistanceName: "Near"}
	nextToRoad := models.ContentRuleRowSave{Name: "Distance to road", DistanceName: "Next To"}
	return []models.ZoneContentRowSave{
		defaultGuardedMineRow(interactable.WoodMine, nearTown),
		defaultGuardedMineRow(interactable.OreMine, nearTown),
		defaultGuardedMineRow(interactable.GoldMine, nearTown),
		defaultGuardedMineRow(interactable.CrystalMine, nextToRoad),
		defaultGuardedMineRow(interactable.MercuryMine, nextToRoad),
		defaultGuardedMineRow(interactable.GemstoneMine, nextToRoad),
		defaultGuardedMineRow(interactable.AlchemyLab, nextToRoad),
	}
}

// defaultGuardedMineRow builds a single guarded mine row with the given
// distance rule.
func defaultGuardedMineRow(sid string, distanceRule models.ContentRuleRowSave) models.ZoneContentRowSave {
	return models.ZoneContentRowSave{
		Sid:    sid,
		Count:  1,
		IsMine: true,
		Rules: []models.ContentRuleRowSave{
			{Name: guardedRuleName, IsGuarded: new(true)},
			distanceRule,
		},
	}
}

// defaultPlayerZoneTreasureRows returns the guarded treasures: a pandora box
// and an epic random item.
func defaultPlayerZoneTreasureRows() []models.ZoneContentRowSave {
	resources := registry.GetMapObjectResourceValues()
	randomItems := registry.GetMapObjectRandomItemValues()

	return []models.ZoneContentRowSave{
		{
			Sid:   resources.PandoraBox,
			Count: 1,
			Rules: []models.ContentRuleRowSave{{Name: guardedRuleName, IsGuarded: new(true)}},
		},
		{
			Sid:   randomItems.RandomItemEpic,
			Count: 1,
			Rules: []models.ContentRuleRowSave{{Name: guardedRuleName, IsGuarded: new(true)}},
		},
	}
}

// defaultPlayerZoneHireAndBankRows returns the guarded random-hire and
// resource-bank group rows.
func defaultPlayerZoneHireAndBankRows() []models.ZoneContentRowSave {
	randomHires := registry.GetMandatoryContentRandomHiresBuildingValues()
	basicRandomHires := registry.GetMandatoryContentBasicRandomHiresBuildingValues()
	basicResourceBanks := registry.GetMandatoryContentBasicResourceBanksBuildingValues()

	trueVal := true // To not allocate multiple times
	return []models.ZoneContentRowSave{
		{
			Sid:     randomHires.RandomHiresLowTier,
			Count:   2,
			IsGroup: true,
			Rules:   []models.ContentRuleRowSave{{Name: guardedRuleName, IsGuarded: &trueVal}},
		},
		{
			Sid:     randomHires.RandomHiresHighTier,
			Count:   1,
			IsGroup: true,
			Rules:   []models.ContentRuleRowSave{{Name: guardedRuleName, IsGuarded: &trueVal}},
		},
		{
			Sid:     basicRandomHires.BasicRandomHires,
			Count:   1,
			IsGroup: true,
			Rules:   []models.ContentRuleRowSave{{Name: guardedRuleName, IsGuarded: &trueVal}},
		},
		{
			Sid:     basicResourceBanks.BasicResourceBanksTier1,
			Count:   2,
			IsGroup: true,
			Rules:   []models.ContentRuleRowSave{{Name: guardedRuleName, IsGuarded: &trueVal}},
		},
		{
			Sid:     basicResourceBanks.BasicResourceBanksTier2,
			Count:   1,
			IsGroup: true,
			Rules:   []models.ContentRuleRowSave{{Name: guardedRuleName, IsGuarded: &trueVal}},
		},
	}
}

func contentRowSlicesEqual(left, right []models.ZoneContentRowSave) bool {
	return (left == nil) == (right == nil) && slices.EqualFunc(left, right, contentRowsEqual)
}

// contentRowsEqual compares the scalar row fields and the rules element-wise.
// New ZoneContentRowSave fields must be added here; the fuzz-parity test on
// EqualsIgnoringManualEdits guards against drift.
func contentRowsEqual(left, right models.ZoneContentRowSave) bool {
	return left.Sid == right.Sid &&
		left.Count == right.Count &&
		left.IsGroup == right.IsGroup &&
		left.IsMine == right.IsMine &&
		(left.Rules == nil) == (right.Rules == nil) &&
		slices.EqualFunc(left.Rules, right.Rules, contentRulesEqual)
}

// contentRulesEqual compares two content rules; the pointer fields compare by
// pointed-to value, matching [reflect.DeepEqual].
func contentRulesEqual(left, right models.ContentRuleRowSave) bool {
	leftScalars := left
	rightScalars := right
	leftScalars.IsGuarded, rightScalars.IsGuarded = nil, nil
	leftScalars.IsSoloEncounter, rightScalars.IsSoloEncounter = nil, nil
	leftScalars.VariantID, rightScalars.VariantID = nil, nil
	return leftScalars == rightScalars &&
		pointedValuesEqual(left.IsGuarded, right.IsGuarded) &&
		pointedValuesEqual(left.IsSoloEncounter, right.IsSoloEncounter) &&
		pointedValuesEqual(left.VariantID, right.VariantID)
}

// pointedValuesEqual reports whether two pointers are both nil or point to
// equal values, matching [reflect.DeepEqual]'s pointer semantics.
func pointedValuesEqual[Value comparable](left, right *Value) bool {
	if left == nil || right == nil {
		return left == right
	}
	return *left == *right
}
