package connectionLineStyle_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/stretchr/testify/assert"
)

func TestWhenAPreviewEdgeCarriesARoad_TheStyleReportsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := preview.Connection{HasRoad: true}

	// Act
	style := utils.NewPreviewConnectionLineStyle(connection)

	// Assert
	assert.Equal(t, utils.ConnectionLineStyle{HasRoad: true}, style)
}

func TestWhenAPreviewEdgeIsAnExplicitPortal_TheStyleReportsBothPortalFlags(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := preview.Connection{
		Type:           preview.ConnectionTypePortal,
		HasRoad:        true,
		ExplicitPortal: true,
	}

	// Act
	style := utils.NewPreviewConnectionLineStyle(connection)

	// Assert
	assert.Equal(t,
		utils.ConnectionLineStyle{HasRoad: true, ExplicitPortal: true, DrawsAsPortal: true},
		style)
}

// The projection classifies a placement-rule-only connection as portal-shaped
// without calling it an explicit portal, and the style has to keep the two
// apart.
func TestWhenAPreviewEdgeIsOnlyPortalShaped_TheStyleReportsNoExplicitPortal(t *testing.T) {
	t.Parallel()
	// Arrange
	connection := preview.Connection{Type: preview.ConnectionTypePortal}

	// Act
	style := utils.NewPreviewConnectionLineStyle(connection)

	// Assert
	assert.Equal(t, utils.ConnectionLineStyle{DrawsAsPortal: true}, style)
}
