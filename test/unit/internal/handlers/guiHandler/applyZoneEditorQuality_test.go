package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenQualityChanges_ReturnsServiceEquivalentZone(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	tuning := models.GenerationTuning{
		ContentScale:                   1,
		ResourceDensityMultiplier:      0.5,
		StructureDensityMultiplier:     1,
		NeutralStackStrengthMultiplier: 1,
		BorderGuardStrengthMultiplier:  1,
	}
	zone := connection_editor.NewDefaultZoneEditorService().
		NewDefaultNeutralZone("Z", neutral_zone.QualityLow, 0, true, tuning)
	expected := zone
	connection_editor.NewDefaultZoneEditorService().
		ApplyNeutralZoneQuality(&expected, neutral_zone.QualityHigh, 3, tuning)

	// Act
	result := handler.ApplyZoneEditorQuality(dtos.ZoneEditorQualityRequestDto{
		Zone:        zone,
		Quality:     neutral_zone.QualityHigh,
		CastleCount: 3,
		Tuning:      tuning,
	})

	// Assert
	assert.Equal(t, expected, result)
}
