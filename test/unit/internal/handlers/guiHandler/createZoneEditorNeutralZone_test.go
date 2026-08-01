package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenNeutralZoneIsCreated_ReturnsServiceEquivalentZone(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	request := dtos.ZoneEditorNeutralZoneRequestDto{
		Label:         "Q",
		Quality:       neutral_zone.QualityMedium,
		CastleCount:   1,
		GenerateRoads: true,
		Tuning: models.GenerationTuning{
			ContentScale:                   1,
			ResourceDensityMultiplier:      0.5,
			StructureDensityMultiplier:     1,
			NeutralStackStrengthMultiplier: 1,
			BorderGuardStrengthMultiplier:  1,
		},
	}
	expected := connection_editor.NewZoneEditorService().NewDefaultNeutralZone(
		request.Label, request.Quality, request.CastleCount, request.GenerateRoads, request.Tuning)

	// Act
	result := handler.CreateZoneEditorNeutralZone(request)

	// Assert
	assert.Equal(t, expected, result)
}
