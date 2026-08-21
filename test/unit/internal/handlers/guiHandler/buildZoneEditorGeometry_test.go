package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenCanvasGeometryRequested_PlacesEveryZone(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	state := editor_state_dto.NewDefaultEditorStateDto()
	variant := generateDefaultTemplate(t, handler).Variants[0]
	require.NotEmpty(t, variant.Zones)

	// Act
	geometry := handler.BuildZoneEditorGeometry(dtos.ZoneEditorGeometryRequestDto{
		Zones:       variant.Zones,
		Connections: variant.Connections,
		Topology:    state.Topology,
		CanvasSide:  700,
	})

	// Assert
	assert.Len(t, geometry.Positions, len(variant.Zones))
}
