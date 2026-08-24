package editor_state_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_zone_contents"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/editor_state_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

// EditorState is the runtime editor state: the nine behaviour-free entity
// groups plus the behaviour that operates across them. The groups are embedded
// anonymously so every field stays directly selectable and encoding/json keeps
// emitting them as one flat object.
type EditorState struct {
	editor_state.EditorState
}

func NewDefaultEditorStateModel() EditorState {
	winConditions := registry.GetWinningConditionValues()
	gameModes := registry.GetGameModeValues()
	return EditorState{
		SchemaVersion:                editor_state.CurrentEditorStateSchemaVersion,
		TemplateName:                 "Custom Template",
		GameMode:                     gameModes.Classic,
		MapSize:                      160,
		PlayerCount:                  2,
		HeroCountMin:                 4,
		HeroCountMax:                 8,
		HeroCountIncrement:           1,
		AbandonedOutpostCount:        1,
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
		Topology:                     config.TopologyRandom,
		MaxPortalConnections:         32,
		SpawnRemoteFootholds:         true,
		RemoteFootholdCount:          1,
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
		PlayerZoneContentRows:        common_zone_contents.GetDefaultPlayerZoneContentRows(),
	}
}

// Clone returns a copy that shares no backing array or pointer with the
// receiver, so an in-place edit of any element on either side stays invisible
// to the other. A plain struct copy duplicates slice headers only, which would
// hide element mutations from the change detection in
// EqualsIgnoringManualEdits.
func (this *EditorState) Clone() EditorState {
	clone := *this

	clone.Bonuses = slices.Clone(this.Bonuses)
	clone.PlayerZoneContentRows = editor_state_helpers.CloneZoneContentRows(this.PlayerZoneContentRows)
	clone.LowestNeutralContentRows = editor_state_helpers.CloneZoneContentRows(this.LowestNeutralContentRows)
	clone.LowNeutralContentRows = editor_state_helpers.CloneZoneContentRows(this.LowNeutralContentRows)
	clone.MediumNeutralContentRows = editor_state_helpers.CloneZoneContentRows(this.MediumNeutralContentRows)
	clone.HighNeutralContentRows = editor_state_helpers.CloneZoneContentRows(this.HighNeutralContentRows)
	clone.HubZoneContentRows = editor_state_helpers.CloneZoneContentRows(this.HubZoneContentRows)

	clone.ManualZones = slices.Clone(this.ManualZones)
	for zoneIndex := range clone.ManualZones {
		zoneModel := ManualZoneSaveModel{ManualZoneSave: clone.ManualZones[zoneIndex]}
		clone.ManualZones[zoneIndex] = zoneModel.Clone().ManualZoneSave
	}

	clone.ManualConnections = slices.Clone(this.ManualConnections)
	for connectionIndex := range clone.ManualConnections {
		connectionModel := ManualConnectionSaveModel{
			ManualConnectionSave: clone.ManualConnections[connectionIndex],
		}
		clone.ManualConnections[connectionIndex] = connectionModel.Clone().ManualConnectionSave
	}

	return clone
}

// LayoutDefiningOptionsChanged reports whether any option that changes the set
// of zones or the connection graph differs between two editor states. When
// these are unchanged, manual zone edits remain valid and can be reapplied.
func (this *EditorState) LayoutDefiningOptionsChanged(incoming *EditorState) bool {
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
func (this *EditorState) DiffCastleSettings(incoming *EditorState) CastleSettingChanges {
	changes := CastleSettingChanges{
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
// field is added to any entity group without extending the comparison.
func (this *EditorState) EqualsIgnoringManualEdits(other *EditorState) bool {
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

func (this *EditorState) HasManualEdits() bool {
	return len(this.ManualZones) > 0 || len(this.ManualConnections) > 0
}

// zoneCountOptionsChanged reports whether the number of neutral zones differs
// between two editor states.
func (this *EditorState) zoneCountOptionsChanged(incoming *EditorState) bool {
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

// zoneOptionScalarsEqual compares the template identity and the zone and castle
// count options.
//
//nolint:dupl // comparison chains over disjoint field groups are inherently similar
func (this *EditorState) zoneOptionScalarsEqual(other *EditorState) bool {
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

// generationOptionScalarsEqual compares the map generation options: sizes,
// heroes, topology, connections and densities.
//
//nolint:dupl // comparison chains over disjoint field groups are inherently similar
func (this *EditorState) generationOptionScalarsEqual(other *EditorState) bool {
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

// gameRuleScalarsEqual compares the victory and game-rule options together with
// the banned content free-text fields.
//
//nolint:dupl // comparison chains over disjoint field groups are inherently similar
func (this *EditorState) gameRuleScalarsEqual(other *EditorState) bool {
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

func contentRowSlicesEqual(left, right []models.ZoneContentRow) bool {
	return (left == nil) == (right == nil) && slices.EqualFunc(left, right, contentRowsEqual)
}

// contentRowsEqual compares the scalar row fields and the rules element-wise.
// New ZoneContentRow fields must be added here; the fuzz-parity test on
// EqualsIgnoringManualEdits guards against drift.
func contentRowsEqual(left, right models.ZoneContentRow) bool {
	return left.Sid == right.Sid &&
		left.Count == right.Count &&
		left.IsGroup == right.IsGroup &&
		left.IsMine == right.IsMine &&
		(left.Rules == nil) == (right.Rules == nil) &&
		slices.EqualFunc(left.Rules, right.Rules, contentRulesEqual)
}

// contentRulesEqual compares two content rules; the pointer fields compare by
// pointed-to value, matching [reflect.DeepEqual].
func contentRulesEqual(left, right models.ContentRuleRow) bool {
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
