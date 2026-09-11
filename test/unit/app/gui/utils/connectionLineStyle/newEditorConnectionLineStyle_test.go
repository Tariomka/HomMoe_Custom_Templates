package connectionLineStyle_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenAnEditedConnectionCarriesAnExplicitRoadFlag_TheStyleReportsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	road := true
	connection := template_model.Connection{ConnectionType: "Direct", Road: &road}

	// Act
	style := utils.NewEditorConnectionLineStyle(connection)

	// Assert
	assert.Equal(t, utils.ConnectionLineStyle{HasRoad: true}, style)
}

func TestWhenAnEditedConnectionCarriesNoRoadFlag_TheStyleReportsNoRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{ConnectionType: "Direct"}

	// Act
	style := utils.NewEditorConnectionLineStyle(connection)

	// Assert
	assert.Equal(t, utils.ConnectionLineStyle{}, style)
}

// A portal without a flag of its own still carries a road, and the editor's
// type comparison is case-insensitive.
func TestWhenAnEditedPortalCarriesNoRoadFlag_TheStyleReportsARoadedPortal(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := template_model.Connection{ConnectionType: "portal"}

	// Act
	style := utils.NewEditorConnectionLineStyle(connection)

	// Assert
	assert.Equal(t,
		utils.ConnectionLineStyle{HasRoad: true, ExplicitPortal: true, DrawsAsPortal: true},
		style)
}

func TestWhenAnEditedPortalIsFlaggedRoadless_TheStyleReportsARoadlessPortal(t *testing.T) {
	t.Parallel()
	// Arrange
	road := false
	connection := template_model.Connection{ConnectionType: "Portal", Road: &road}

	// Act
	style := utils.NewEditorConnectionLineStyle(connection)

	// Assert
	assert.Equal(t,
		utils.ConnectionLineStyle{ExplicitPortal: true, DrawsAsPortal: true},
		style)
}
