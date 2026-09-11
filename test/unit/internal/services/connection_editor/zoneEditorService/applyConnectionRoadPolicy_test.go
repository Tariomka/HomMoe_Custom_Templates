package zoneEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenRoadsAreOn_StampsThePendingConnectionTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{From: "Spawn-A", To: "Spawn-B", ConnectionType: "Direct"}

	// Act
	test_helpers.NewZoneEditorService().ApplyConnectionRoadPolicy(&connection, true)

	// Assert
	require.NotNil(t, connection.Road)
	assert.True(t, *connection.Road)
}

func TestWhenRoadsAreOff_StampsThePendingConnectionFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{From: "Spawn-A", To: "Spawn-B", ConnectionType: "Direct"}

	// Act
	test_helpers.NewZoneEditorService().ApplyConnectionRoadPolicy(&connection, false)

	// Assert
	require.NotNil(t, connection.Road)
	assert.False(t, *connection.Road)
}

func TestWhenPendingConnectionIsAnExplicitPortal_LeavesItsOmittedFlagAlone(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{From: "Spawn-A", To: "Spawn-B", ConnectionType: "Portal"}

	// Act
	test_helpers.NewZoneEditorService().ApplyConnectionRoadPolicy(&connection, false)

	// Assert
	assert.Nil(t, connection.Road)
}

func TestWhenPendingPortalSaysNoRoadAndRoadsAreOn_KeepsThatFlag(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{
		From:           "Spawn-A",
		To:             "Spawn-B",
		ConnectionType: "Portal",
		Road:           new(false),
	}

	// Act
	test_helpers.NewZoneEditorService().ApplyConnectionRoadPolicy(&connection, true)

	// Assert
	assert.False(t, *connection.Road)
}

func TestWhenPendingConnectionIsStamped_LeavesItsOtherFieldsAlone(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{
		Name:           "Rnd-A-B",
		From:           "Spawn-A",
		To:             "Spawn-B",
		ConnectionType: "Direct",
		GuardValue:     7500,
		IsUserAdded:    true,
	}
	expected := connection

	// Act
	test_helpers.NewZoneEditorService().ApplyConnectionRoadPolicy(&connection, true)

	// Assert
	connection.Road = nil
	assert.Equal(t, expected, connection)
}
