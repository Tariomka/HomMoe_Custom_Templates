package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenRequestContainsTemplate_ReturnsServiceLayoutUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	template := template_model.Template{Variants: []template_model.Variant{{
		Zones: []template_model.Zone{
			{Name: "Spawn-A"},
			{Name: "Spawn-B"},
			{Name: "Spawn-C"},
		},
		Connections: []template_model.Connection{
			{From: "Spawn-A", To: "Spawn-B"},
			{From: "Spawn-B", To: "Spawn-C"},
		},
		Orientation: template_model.Orientation{ZeroAngleZone: "Spawn-C"},
	}}}
	request := dtos.PreviewLayoutRequestDto{
		Template:   &template,
		Topology:   config.TopologyRing,
		CanvasSide: 600,
	}
	expected := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService()).BuildPreviewLayout(
		&template,
		request.Topology,
		request.CanvasSide,
	)
	handler := newProductionGuiHandler()

	// Act
	result := handler.BuildPreviewLayout(request)

	// Assert
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
	result := handler.BuildPreviewLayout(request)

	// Assert
	assert.Equal(t, expected, result.Layout)
}
