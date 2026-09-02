package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenCanvasGeometryRequested_PlacesEveryZone(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	state := editor_state_model.NewDefaultEditorStateModel()
	variant := generateDefaultTemplate(t, handler).Variants[0]
	zones := template_model.ToZoneEntities(variant.Zones)
	require.NotEmpty(t, zones)

	// Act
	geometry := handler.BuildZoneEditorGeometry(dtos.ZoneEditorGeometryRequestDto{
		Zones:       zones,
		Connections: template_model.ToConnectionEntities(variant.Connections),
		Topology:    state.Topology,
		CanvasSide:  700,
	})

	// Assert
	assert.Len(t, geometry.Positions, len(zones))
}
