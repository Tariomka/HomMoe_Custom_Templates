package zoneEditorHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenCanvasGeometryIsRequested_ReturnsTheServiceLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	request := dtos.ZoneEditorGeometryRequestDto{
		Topology:   config.MapTopology(gofakeit.Word()),
		CanvasSide: gofakeit.Number(100, 900),
	}
	expected := models.ZoneEditorGeometry{ZoneRadius: gofakeit.Float64Range(1, 60)}
	fixture.geometry.
		On("BuildGeometry", request.Zones, request.Connections, request.Topology, request.CanvasSide).
		Return(expected)

	// Act
	geometry := fixture.handler.BuildZoneEditorGeometry(request)

	// Assert
	assert.Equal(t, expected, geometry)
}
