package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenRequestContainsTemplate_ReturnsServiceLayoutUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	rmgTemplate := entities.RmgTemplate{Variants: []entities.Variant{{
		Zones: []entities.Zone{
			{Name: "Spawn-A"},
			{Name: "Spawn-B"},
			{Name: "Spawn-C"},
		},
		Connections: []entities.Connection{
			{From: "Spawn-A", To: "Spawn-B"},
			{From: "Spawn-B", To: "Spawn-C"},
		},
		Orientation: entities.Orientation{ZeroAngleZone: "Spawn-C"},
	}}}
	request := dtos.PreviewLayoutRequestDto{
		Template:   new(mappers.NewTemplateMapper().ToModel(rmgTemplate)),
		Topology:   config.TopologyRing,
		CanvasSide: 600,
	}
	expected := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService()).BuildPreviewLayout(
		&rmgTemplate,
		request.Topology,
		request.CanvasSide,
	)
	handler := newProductionGuiHandler()

	// Act
	result, err := handler.BuildPreviewLayout(request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, result.Layout)
}

func TestWhenTemplateIsNil_ReturnsEmptyServiceLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	request := dtos.PreviewLayoutRequestDto{
		Topology:   config.TopologyRing,
		CanvasSide: 600,
	}
	expected := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService()).BuildPreviewLayout(
		nil,
		request.Topology,
		request.CanvasSide,
	)
	handler := newProductionGuiHandler()

	// Act
	result, err := handler.BuildPreviewLayout(request)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, result.Layout)
}

func TestWhenRequestContainsZones_BuildsPreviewTemplate(t *testing.T) {
	t.Parallel()
	// Arrange
	request := dtos.PreviewLayoutRequestDto{
		Zones: []entities.Zone{
			{Name: "Spawn-A"},
			{Name: "Neutral-B"},
		},
		Connections: []entities.Connection{{From: "Spawn-A", To: "Neutral-B"}},
		Topology:    config.TopologyRing,
		CanvasSide:  600,
	}
	handler := newProductionGuiHandler()

	// Act
	result, err := handler.BuildPreviewLayout(request)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result.Layout.Positions, 2)
}
