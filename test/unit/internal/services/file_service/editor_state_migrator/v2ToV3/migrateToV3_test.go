package v2ToV3_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state/editor_state_v2"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service/editor_state_migrator"
	"github.com/stretchr/testify/assert"
)

func TestWhenAV2StateIsMigrated_EveryFlatGroupLandsInItsSection(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := newV2State()

	// Act
	migrated, _ := editor_state_migrator.MigrateToV3(legacy)

	// Assert
	assert.Equal(t,
		[]any{"Legacy", 4, "Square", true},
		[]any{migrated.TemplateName, migrated.Players.PlayerCount, string(migrated.Generation.Topology),
			migrated.GameRules.CityHold})
}

func TestWhenAV2StateIsMigrated_TheTextListsBecomeEntries(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := newV2State()

	// Act
	migrated, _ := editor_state_migrator.MigrateToV3(legacy)

	// Assert
	assert.Equal(t,
		[]any{
			[]editor_state.BanEntry{{Sid: "item_a"}, {Sid: "item_a"}},
			[]editor_state.BanEntry{{Sid: "magic_a"}},
			[]editor_state.ValueOverrideEntry{{Sid: "mine", GuardValue: 500}},
		},
		[]any{migrated.Content.BannedItems, migrated.Content.BannedMagics, migrated.Content.ValueOverrides})
}

func TestWhenAV2OverrideLineIsMalformed_ItIsDroppedWithAWarning(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := newV2State()
	legacy.ValueOverridesText = "mine=500\nbroken"

	// Act
	_, warnings := editor_state_migrator.MigrateToV3(legacy)

	// Assert
	assert.Equal(t, []string{"value override line 2 'broken' is not sid=value and was dropped"}, warnings)
}

func TestWhenAV2StateIsMigrated_TheSnapshotIsKeptWhole(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := newV2State()

	// Act
	migrated, _ := editor_state_migrator.MigrateToV3(legacy)

	// Assert
	assert.Equal(t,
		editor_state.ManualZoneEdit{Name: "Neutral-A", FieldOverrides: editor_state.ZoneFieldOverrides{Size: new(1.5)}},
		migrated.ManualEdits.Zones[0])
}

func TestWhenAV2ZoneRecordsATier_TheTierIsKept(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := newV2State()
	legacy.ManualZones[0].Quality = new(int8(0))

	// Act
	migrated, _ := editor_state_migrator.MigrateToV3(legacy)

	// Assert
	assert.Equal(t, new(int8(0)), migrated.ManualEdits.Zones[0].Quality)
}

func TestWhenAV2ConnectionWasAddedByTheUser_TheFlagIsKept(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := newV2State()
	var connection template_entity.Connection
	connection.Name = "Manual-A-B"
	legacy.ManualConnections = []editor_state_v2.ManualConnectionSave{{Connection: connection, IsUserAdded: true}}

	// Act
	migrated, _ := editor_state_migrator.MigrateToV3(legacy)

	// Assert
	assert.True(t, migrated.ManualEdits.Connections[0].IsUserAdded)
}

func TestWhenAV2StateIsMigrated_ItAwaitsDerivationAtTheFirstDefaultsVersion(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := newV2State()

	// Act
	migrated, _ := editor_state_migrator.MigrateToV3(legacy)

	// Assert
	assert.Equal(t,
		[]int{editor_state.CurrentEditorStateSchemaVersion, 0, 1},
		[]int{migrated.SchemaVersion, migrated.GeneratorVersion, migrated.Content.ContentDefaultsVersion})
}

// A v2 file with no player rows had them seeded as defaults, so the migration
// keeps them omitted rather than clearing them.
func TestWhenAV2StateHasNoPlayerRows_TheyStayOmitted(t *testing.T) {
	t.Parallel()
	// Arrange
	legacy := newV2State()

	// Act
	migrated, _ := editor_state_migrator.MigrateToV3(legacy)

	// Assert
	assert.Nil(t, migrated.Content.PlayerZoneContentRows)
}

func newV2State() editor_state_v2.EditorState {
	var zone template_entity.Zone
	zone.Name = "Neutral-A"
	zone.Size = 1.5

	var legacy editor_state_v2.EditorState
	legacy.SchemaVersion = 2
	legacy.TemplateName = "Legacy"
	legacy.PlayerCount = 4
	legacy.Topology = "Square"
	legacy.CityHold = true
	legacy.BannedItems = "item_a\n\nitem_a"
	legacy.BannedMagics = " magic_a "
	legacy.ValueOverridesText = "mine=500"
	legacy.ManualZones = []editor_state_v2.ManualZoneSave{{Zone: zone}}
	return legacy
}
