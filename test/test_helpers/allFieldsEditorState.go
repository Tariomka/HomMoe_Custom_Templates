package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// NewAllFieldsEditorStateDto wraps the all-fields model in the persistence
// shell, which is what the frozen .gen.json fixture is written from.
func NewAllFieldsEditorStateDto() editor_state_dto.EditorStateDto {
	return editor_state_dto.NewEditorStateDto(NewAllFieldsEditorStateModel())
}

// NewAllFieldsEditorStateModel builds an editor state in which every persisted
// field carries a distinctive value that differs from both the Go zero value
// and the seeded default. It backs the frozen .gen.json fixture guarding the
// on-disk wire format: unmarshalling the fixture into a zero-valued state must
// reproduce exactly this value, so a field that stops being written or read is
// impossible to miss.
//
//nolint:funlen // one assignment per persisted field; splitting it would hide the 1:1 field coverage this fixture exists to guarantee.
func NewAllFieldsEditorStateModel() editor_state_model.EditorState {
	return editor_state_model.EditorState{
		TemplateName: "All Fields Fixture",
		GameMode:     "SingleHero",

		MapSize:              208,
		ExperimentalMapSizes: true,

		PlayerCount:        7,
		HeroCountMin:       2,
		HeroCountMax:       9,
		HeroCountIncrement: 3,

		NeutralZoneCount:           13,
		SpawnAbandonedOutposts:     true,
		AbandonedOutpostCount:      4,
		NeutralLowestNoCastleCount: 21,
		NeutralLowestCastleCount:   22,
		NeutralLowNoCastleCount:    23,
		NeutralLowCastleCount:      24,
		NeutralMediumNoCastleCount: 25,
		NeutralMediumCastleCount:   26,
		NeutralHighNoCastleCount:   27,
		NeutralHighCastleCount:     28,

		AdvancedMode:                true,
		PlayerOwnedCastles:          3,
		PlayerZoneCastles:           4,
		NeutralZoneCastles:          2,
		HubZoneCastles:              5,
		NeutralLowestCastlesPerZone: 31,
		NeutralLowCastlesPerZone:    32,
		NeutralMediumCastlesPerZone: 33,
		NeutralHighCastlesPerZone:   34,
		MatchPlayerCastleFactions:   true,

		PlayerZoneSize:              1.25,
		NeutralZoneSize:             1.5,
		HubZoneSize:                 1.75,
		GuardRandomization:          0.35,
		Topology:                    config.TopologyChain,
		RandomPortals:               true,
		MaxPortalConnections:        17,
		SpawnRemoteFootholds:        true,
		RemoteFootholdCount:         6,
		GenerateRoads:               true,
		NoDirectPlayerConn:          true,
		ResourceDensityPercent:      111,
		StructureDensityPercent:     122,
		NeutralStackStrengthPercent: 133,
		BorderGuardStrengthPercent:  144,

		VictoryCondition:             "win_condition_4",
		FactionLawXpPercent:          155,
		AstrologyXpPercent:           166,
		LostStartCity:                true,
		LostStartCityDay:             8,
		LostStartHero:                true,
		CityHold:                     true,
		CityHoldDays:                 9,
		GladiatorArena:               true,
		GladiatorArenaDaysDelayStart: 41,
		GladiatorArenaCountDay:       5,
		Tournament:                   true,
		TournamentFirstTournamentDay: 18,
		TournamentInterval:           11,
		TournamentPointsToWin:        7,
		TournamentSaveArmy:           true,

		BannedItems:        "fixture_banned_item_one\nfixture_banned_item_two",
		BannedMagics:       "fixture_banned_magic",
		ValueOverridesText: "fixture_override_sid=12345",
		Bonuses: []config.BonusEntry{{
			PresetType:     config.BonusStartingGold,
			ReceiverFilter: "all_heroes",
			Param:          "9500",
			Param2:         "1",
		}},

		PlayerZoneContentRows:    allFieldsContentRows("fixture_player_zone_row", 41),
		LowestNeutralContentRows: allFieldsContentRows("fixture_lowest_neutral_row", 42),
		LowNeutralContentRows:    allFieldsContentRows("fixture_low_neutral_row", 43),
		MediumNeutralContentRows: allFieldsContentRows("fixture_medium_neutral_row", 44),
		HighNeutralContentRows:   allFieldsContentRows("fixture_high_neutral_row", 45),
		HubZoneContentRows:       allFieldsContentRows("fixture_hub_zone_row", 46),

		ManualZones:       allFieldsManualZones(),
		ManualConnections: allFieldsManualConnections(),
	}
}

// allFieldsContentRows builds a one-row content list whose every field, and
// every field of its single rule, is populated.
func allFieldsContentRows(sid string, variantID int) []models.ZoneContentRowSave {
	return []models.ZoneContentRowSave{{
		Sid:     sid,
		Count:   variantID,
		IsGroup: true,
		IsMine:  true,
		Rules: []models.ContentRuleRowSave{{
			Name:            "Guarded",
			DistanceName:    "Very Far",
			IsGuarded:       new(true),
			IsSoloEncounter: new(true),
			VariantID:       new(variantID),
		}},
	}}
}

// allFieldsManualZones builds a manual zone carrying the normalized position
// that entities.Zone itself omits from JSON.
func allFieldsManualZones() []editor_state.ManualZoneSave {
	return []editor_state.ManualZoneSave{{
		Zone: entities.Zone{
			Name:                 "Fixture-Spawn-A",
			Size:                 1.35,
			Layout:               "Sides",
			GuardCutoffValue:     17500,
			GuardRandomization:   0.15,
			GuardWeeklyIncrement: 0.2,
		},
		ManualPosition: &[2]float64{0.25, 0.75},
	}}
}

// allFieldsManualConnections builds a manual connection carrying the
// user-added flag that entities.Connection itself omits from JSON.
func allFieldsManualConnections() []editor_state.ManualConnectionSave {
	return []editor_state.ManualConnectionSave{{
		Connection: entities.Connection{
			Name:                 "Fixture-Conn-A-B",
			From:                 "Fixture-Spawn-A",
			To:                   "Fixture-Spawn-B",
			ConnectionType:       "Portal",
			GuardValue:           26500,
			GuardWeeklyIncrement: 0.15,
		},
		IsUserAdded: true,
	}}
}
