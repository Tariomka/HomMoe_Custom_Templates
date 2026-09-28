package zoneEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTypeIsPicked_RecordsItOnTheConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{From: "Spawn-A", To: "Spawn-B", ConnectionType: "Direct"}

	// Act
	test_helpers.NewZoneEditorService().ChangeConnectionType(&connection, "Portal", true)

	// Assert
	assert.Equal(t, "Portal", connection.ConnectionType)
}

func TestWhenTypeStaysNonPortalAndRoadsAreOn_StampsRoadTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{From: "Spawn-A", To: "Spawn-B", ConnectionType: "Direct"}

	// Act
	test_helpers.NewZoneEditorService().ChangeConnectionType(&connection, "Direct", true)

	// Assert
	require.NotNil(t, connection.Road)
	assert.True(t, *connection.Road)
}

func TestWhenTypeStaysNonPortalAndRoadsAreOff_StampsRoadFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{From: "Spawn-A", To: "Spawn-B", ConnectionType: "Direct"}

	// Act
	test_helpers.NewZoneEditorService().ChangeConnectionType(&connection, "Direct", false)

	// Assert
	require.NotNil(t, connection.Road)
	assert.False(t, *connection.Road)
}

func TestWhenTypeBecomesPortal_KeepsTheRoadFlagItCarried(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{
		From:           "Spawn-A",
		To:             "Spawn-B",
		ConnectionType: "Direct",
		Road:           new(false),
	}

	// Act
	test_helpers.NewZoneEditorService().ChangeConnectionType(&connection, "Portal", true)

	// Assert
	assert.False(t, *connection.Road)
}

func TestWhenTypeBecomesPortalWithoutARoadFlag_LeavesItOmitted(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{From: "Spawn-A", To: "Spawn-B", ConnectionType: "Direct"}

	// Act
	test_helpers.NewZoneEditorService().ChangeConnectionType(&connection, "Portal", false)

	// Assert
	assert.Nil(t, connection.Road)
}

func TestWhenPortalBecomesDirectAndRoadsAreOn_StampsRoadTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{
		From:           "Spawn-A",
		To:             "Spawn-B",
		ConnectionType: "Portal",
		Road:           new(false),
	}

	// Act
	test_helpers.NewZoneEditorService().ChangeConnectionType(&connection, "Direct", true)

	// Assert
	assert.True(t, *connection.Road)
}

func TestWhenPortalBecomesDirectAndRoadsAreOff_StampsRoadFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{
		From:           "Spawn-A",
		To:             "Spawn-B",
		ConnectionType: "Portal",
		Road:           new(true),
	}

	// Act
	test_helpers.NewZoneEditorService().ChangeConnectionType(&connection, "Direct", false)

	// Assert
	assert.False(t, *connection.Road)
}

func TestWhenTypeChanges_LeavesTheConnectionsOtherFieldsAlone(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{
		Name:            "Rnd-A-B",
		From:            "Spawn-A",
		To:              "Spawn-B",
		ConnectionType:  "Direct",
		GuardValue:      7500,
		GuardZone:       "Spawn-A",
		GuardMatchGroup: "rnd_guard_A_B",
		IsUserAdded:     true,
	}
	expected := connection
	expected.ConnectionType = "Portal"

	// Act
	test_helpers.NewZoneEditorService().ChangeConnectionType(&connection, "Portal", true)

	// Assert
	assert.Equal(t, expected, connection)
}

func TestWhenAPortalBecomesDirect_ItsSourcePlacementRulesAreCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := portalWithPlacementRules()

	// Act
	test_helpers.NewZoneEditorService().ChangeConnectionType(&connection, "Direct", true)

	// Assert
	assert.Nil(t, connection.PortalPlacementRulesFrom)
}

func TestWhenAPortalBecomesDirect_ItsTargetPlacementRulesAreCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := portalWithPlacementRules()

	// Act
	test_helpers.NewZoneEditorService().ChangeConnectionType(&connection, "Direct", true)

	// Assert
	assert.Nil(t, connection.PortalPlacementRulesTo)
}

func TestWhenAPortalBecomesProximity_ItStopsBeingAnEffectivePortal(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := portalWithPlacementRules()

	// Act
	test_helpers.NewZoneEditorService().ChangeConnectionType(&connection, "Proximity", true)

	// Assert
	assert.False(t, connection.IsEffectivePortal())
}

func TestWhenTypeBecomesPortal_ThePlacementRulesAreKept(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := portalWithPlacementRules()
	connection.ConnectionType = "Proximity"
	expected := portalWithPlacementRules().PortalPlacementRulesFrom

	// Act
	test_helpers.NewZoneEditorService().ChangeConnectionType(&connection, "Portal", true)

	// Assert
	assert.Equal(t, expected, connection.PortalPlacementRulesFrom)
}

func TestWhenTypeBecomesALowerCasedPortal_ThePlacementRulesAreKept(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := portalWithPlacementRules()
	expected := portalWithPlacementRules().PortalPlacementRulesTo

	// Act
	test_helpers.NewZoneEditorService().ChangeConnectionType(&connection, "portal", true)

	// Assert
	assert.Equal(t, expected, connection.PortalPlacementRulesTo)
}

// portalWithPlacementRules is a portal carrying placement rules on both ends,
// which is what the generator's topology portals look like.
func portalWithPlacementRules() template_model.Connection {
	return template_model.Connection{
		From:                     "Spawn-A",
		To:                       "Spawn-B",
		ConnectionType:           "Portal",
		PortalPlacementRulesFrom: []template_model.PlacementRule{{Type: "MainObject"}},
		PortalPlacementRulesTo:   []template_model.PlacementRule{{Type: "MainObject"}},
	}
}
