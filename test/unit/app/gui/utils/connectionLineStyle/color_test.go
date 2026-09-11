package connectionLineStyle_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/stretchr/testify/assert"
)

func TestWhenAnEdgeCarriesARoad_ItKeepsTheDirectColour(t *testing.T) {
	t.Parallel()
	// Arrange
	style := utils.ConnectionLineStyle{HasRoad: true}

	// Act
	lineColor := style.Color()

	// Assert
	assert.Equal(t, themes.ColorsPreview.DirectLine, lineColor)
}

func TestWhenARoadedEdgeIsDrawnAsAPortal_ItKeepsThePortalColour(t *testing.T) {
	t.Parallel()
	// Arrange
	style := utils.ConnectionLineStyle{HasRoad: true, ExplicitPortal: true, DrawsAsPortal: true}

	// Act
	lineColor := style.Color()

	// Assert
	assert.Equal(t, themes.ColorsPreview.PortalLine, lineColor)
}

// A connection that only carries portal placement rules is drawn with the
// portal shape but is not an explicit Portal, so its road state is the
// non-portal one.
func TestWhenAPortalShapedEdgeIsNotAnExplicitPortal_ARoadlessOneGoesGrey(t *testing.T) {
	t.Parallel()
	// Arrange
	style := utils.ConnectionLineStyle{DrawsAsPortal: true}

	// Act
	lineColor := style.Color()

	// Assert
	assert.Equal(t, themes.ColorsPreview.NoRoadLine, lineColor)
}

func TestWhenAnEdgeCarriesNoRoad_ItGoesGrey(t *testing.T) {
	t.Parallel()
	// Arrange
	style := utils.ConnectionLineStyle{}

	// Act
	lineColor := style.Color()

	// Assert
	assert.Equal(t, themes.ColorsPreview.NoRoadLine, lineColor)
}

func TestWhenAnExplicitPortalCarriesNoRoad_ItGoesGreen(t *testing.T) {
	t.Parallel()
	// Arrange
	style := utils.ConnectionLineStyle{ExplicitPortal: true, DrawsAsPortal: true}

	// Act
	lineColor := style.Color()

	// Assert
	assert.Equal(t, themes.ColorsPreview.PortalNoRoadLine, lineColor)
}

// The roadless colours are the one place the palette must not blend with the
// canvas behind it.
func TestWhenAnEdgeCarriesNoRoad_ItsColourIsOpaque(t *testing.T) {
	t.Parallel()
	// Arrange
	roadless := utils.ConnectionLineStyle{}
	roadlessPortal := utils.ConnectionLineStyle{ExplicitPortal: true}

	// Act
	alphas := []uint8{roadless.Color().A, roadlessPortal.Color().A}

	// Assert
	assert.Equal(t, []uint8{255, 255}, alphas)
}
