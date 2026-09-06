package v1ToV2_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state/editor_state_v1"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/topology"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service/editor_state_migrator"
	"github.com/stretchr/testify/assert"
)

func TestWhenALegacyStateIsMigrated_EverySettingsGroupCrossesUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := newLegacyState()
	expected := editor_state.EditorState{
		TemplateIdentity:    editor_state.TemplateIdentity(legacy.TemplateIdentity),
		MapSettings:         editor_state.MapSettings(legacy.MapSettings),
		PlayerSettings:      editor_state.PlayerSettings(legacy.PlayerSettings),
		NeutralZoneSettings: editor_state.NeutralZoneSettings(legacy.NeutralZoneSettings),
		CastleSettings:      editor_state.CastleSettings(legacy.CastleSettings),
		GenerationSettings:  editor_state.GenerationSettings(legacy.GenerationSettings),
		GameRuleSettings:    editor_state.GameRuleSettings(legacy.GameRuleSettings),
	}

	// Act
	migrated := editor_state_migrator.MigrateToV2(legacy)
	migrated.ContentSettings = editor_state.ContentSettings{}
	migrated.ManualEditSettings = editor_state.ManualEditSettings{}
	migrated.SchemaOptions = editor_state.SchemaOptions{}

	// Assert
	assert.Equal(t, expected, migrated)
}

func TestWhenALegacyStateIsMigrated_ItLandsAtTheCurrentSchemaVersion(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := newLegacyState()

	// Act
	migrated := editor_state_migrator.MigrateToV2(legacy)

	// Assert
	assert.Equal(t, editor_state.CurrentEditorStateSchemaVersion, migrated.SchemaVersion)
}

// The version a v1 file happens to carry is not the version the migrated state
// is at, so it must not be copied across.
func TestWhenALegacyStateCarriesItsOwnVersion_TheMigratedStateOverwritesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := newLegacyState()
	legacy.SchemaVersion = editor_state_v1.SchemaVersion

	// Act
	migrated := editor_state_migrator.MigrateToV2(legacy)

	// Assert
	assert.Equal(t, editor_state.CurrentEditorStateSchemaVersion, migrated.SchemaVersion)
}

func TestWhenALegacyStateCarriesContentRows_TheirRulesCrossToo(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := newLegacyState()
	expected := []editor_state.ZoneContentRow{{
		Sid:     "row_sid",
		Count:   3,
		IsGroup: true,
		IsMine:  true,
		Rules:   []editor_state.ContentRuleRow{{Name: "Guarded", DistanceName: "Far", VariantID: new(4)}},
	}}

	// Act
	migrated := editor_state_migrator.MigrateToV2(legacy)

	// Assert
	assert.Equal(t, expected, migrated.PlayerZoneContentRows)
}

func TestWhenALegacyStateCarriesBonuses_TheirPresetOrdinalsCross(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := newLegacyState()
	expected := []editor_state.BonusEntry{{
		PresetType:     editor_state.BonusStartingGold,
		ReceiverFilter: "all_heroes",
		Param:          "9500",
		Param2:         "1",
	}}

	// Act
	migrated := editor_state_migrator.MigrateToV2(legacy)

	// Assert
	assert.Equal(t, expected, migrated.Bonuses)
}

func TestWhenALegacyZoneCarriesAnArrayPosition_ItBecomesAVector(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := newLegacyState()

	// Act
	migrated := editor_state_migrator.MigrateToV2(legacy)

	// Assert
	assert.Equal(t, data.NewVec2(0.25, 0.75), *migrated.ManualZones[0].ManualPosition)
}

// nil means the zone was never dragged, which is not the same as a zone parked
// at the top-left corner.
func TestWhenALegacyZoneCarriesNoPosition_ItStaysUnpositioned(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := newLegacyState()
	legacy.ManualZones[0].ManualPosition = nil

	// Act
	migrated := editor_state_migrator.MigrateToV2(legacy)

	// Assert
	assert.Nil(t, migrated.ManualZones[0].ManualPosition)
}

func TestWhenALegacyZoneRecordsATier_TheOrdinalCrosses(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := newLegacyState()

	// Act
	migrated := editor_state_migrator.MigrateToV2(legacy)

	// Assert
	assert.Equal(t, int8(2), *migrated.ManualZones[0].Quality)
}

// v1 had nowhere to persist the generator stamps, so a migrated zone reads as
// never stamped and the preview falls back to a computed layout.
func TestWhenALegacyZoneIsMigrated_ItCarriesNoGeneratorStamps(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := newLegacyState()

	// Act
	migrated := editor_state_migrator.MigrateToV2(legacy)

	// Assert
	assert.Equal(
		t,
		editor_state.ManualZoneSave{
			Zone:           legacy.ManualZones[0].Zone,
			ManualPosition: new(data.NewVec2(0.25, 0.75)),
			Quality:        legacy.ManualZones[0].Quality,
		},
		migrated.ManualZones[0])
}

func TestWhenALegacyStateCarriesManualConnections_TheyCrossUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := newLegacyState()
	expected := []editor_state.ManualConnectionSave{{
		Connection:  template_entity.Connection{Name: "Conn A", ConnectionType: "Portal"},
		IsUserAdded: true,
	}}

	// Act
	migrated := editor_state_migrator.MigrateToV2(legacy)

	// Assert
	assert.Equal(t, expected, migrated.ManualConnections)
}

// The change detection distinguishes a nil list from an empty one, so the
// migration must not manufacture empty slices for keys the file never had.
func TestWhenALegacyStateCarriesNoManualZones_TheListStaysNil(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := editor_state_v1.EditorState{}

	// Act
	migrated := editor_state_migrator.MigrateToV2(legacy)

	// Assert
	assert.Nil(t, migrated.ManualZones)
}

func TestWhenALegacyStateCarriesNoContentRows_TheListStaysNil(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := editor_state_v1.EditorState{}

	// Act
	migrated := editor_state_migrator.MigrateToV2(legacy)

	// Assert
	assert.Nil(t, migrated.PlayerZoneContentRows)
}

// newLegacyState populates one field of every group, so a group that stopped
// crossing shows up as a diff rather than as a silently zero value.
func newLegacyState() editor_state_v1.EditorState {
	return editor_state_v1.EditorState{
		TemplateName: "Legacy",
		GameMode:     "SingleHero",

		MapSize:              192,
		ExperimentalMapSizes: true,

		PlayerCount:  6,
		HeroCountMin: 2,
		HeroCountMax: 9,

		NeutralZoneCount:      13,
		AbandonedOutpostCount: 4,

		AdvancedMode:      true,
		PlayerZoneCastles: 4,

		PlayerZoneSize: 1.25,
		Topology:       topology.TopologyChain,
		GenerateRoads:  true,

		VictoryCondition: "win_condition_4",
		CityHoldDays:     9,

		BannedItems:  "banned_item",
		BannedMagics: "banned_magic",
		Bonuses: []editor_state_v1.BonusEntry{{
			PresetType:     editor_state_v1.BonusStartingGold,
			ReceiverFilter: "all_heroes",
			Param:          "9500",
			Param2:         "1",
		}},
		PlayerZoneContentRows: []editor_state_v1.ZoneContentRow{{
			Sid:     "row_sid",
			Count:   3,
			IsGroup: true,
			IsMine:  true,
			Rules: []editor_state_v1.ContentRuleRow{{
				Name:         "Guarded",
				DistanceName: "Far",
				VariantID:    new(4),
			}},
		}},
		ManualZones: []editor_state_v1.ManualZoneSave{{
			Zone:           template_entity.Zone{Name: "Zone A", Size: 1.35},
			ManualPosition: &[2]float64{0.25, 0.75},
			Quality:        new(int8(2)),
		}},
		ManualConnections: []editor_state_v1.ManualConnectionSave{{
			Connection:  template_entity.Connection{Name: "Conn A", ConnectionType: "Portal"},
			IsUserAdded: true,
		}},
	}
}
