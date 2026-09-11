package roadPolicyService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenRoadsAreOnAndConnectionIsNotAPortal_StampsRoadTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	types := []string{"Direct", "Default", "GladiatorArena", "Proximity", "Custom", ""}

	for _, connectionType := range types {
		t.Run(connectionType, func(t *testing.T) {
			t.Parallel()
			// Arrange
			connection := template_model.Connection{From: "A", To: "B", ConnectionType: connectionType}

			// Act
			newPolicy().StampConnectionRoad(&connection, true)

			// Assert
			require.NotNil(t, connection.Road)
			assert.True(t, *connection.Road)
		})
	}
}

func TestWhenRoadsAreOffAndConnectionIsNotAPortal_StampsRoadFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	types := []string{"Direct", "Default", "GladiatorArena", "Proximity", "Custom", ""}

	for _, connectionType := range types {
		t.Run(connectionType, func(t *testing.T) {
			t.Parallel()
			// Arrange
			connection := template_model.Connection{From: "A", To: "B", ConnectionType: connectionType}

			// Act
			newPolicy().StampConnectionRoad(&connection, false)

			// Assert
			require.NotNil(t, connection.Road)
			assert.False(t, *connection.Road)
		})
	}
}

func TestWhenImportedNonPortalSaysNoRoadAndRoadsAreOn_OverwritesItWithTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{From: "A", To: "B", ConnectionType: "Direct", Road: new(false)}

	// Act
	newPolicy().StampConnectionRoad(&connection, true)

	// Assert
	assert.True(t, *connection.Road)
}

func TestWhenImportedNonPortalSaysRoadAndRoadsAreOff_OverwritesItWithFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{From: "A", To: "B", ConnectionType: "Direct", Road: new(true)}

	// Act
	newPolicy().StampConnectionRoad(&connection, false)

	// Assert
	assert.False(t, *connection.Road)
}

func TestWhenConnectionIsAnExplicitPortalWithoutARoadFlag_LeavesItOmitted(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{From: "A", To: "B", ConnectionType: "Portal"}

	// Act
	newPolicy().StampConnectionRoad(&connection, true)

	// Assert
	assert.Nil(t, connection.Road)
}

func TestWhenExplicitPortalSaysNoRoadAndRoadsAreOn_KeepsThatFlag(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{From: "A", To: "B", ConnectionType: "Portal", Road: new(false)}

	// Act
	newPolicy().StampConnectionRoad(&connection, true)

	// Assert
	assert.False(t, *connection.Road)
}

func TestWhenExplicitPortalSaysRoadAndRoadsAreOff_KeepsThatFlag(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{From: "A", To: "B", ConnectionType: "Portal", Road: new(true)}

	// Act
	newPolicy().StampConnectionRoad(&connection, false)

	// Assert
	assert.True(t, *connection.Road)
}

func TestWhenPortalTypeUsesAnotherCase_IsStillTreatedAsAPortal(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{From: "A", To: "B", ConnectionType: "portal"}

	// Act
	newPolicy().StampConnectionRoad(&connection, false)

	// Assert
	assert.Nil(t, connection.Road)
}

func TestWhenSingleNonPortalCarriesPortalPlacementRules_StillFollowsTheSetting(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{
		From:                     "A",
		To:                       "B",
		ConnectionType:           "Direct",
		PortalPlacementRulesFrom: []template_model.PlacementRule{{Type: "MainObject"}},
	}

	// Act
	newPolicy().StampConnectionRoad(&connection, false)

	// Assert
	assert.False(t, *connection.Road)
}
