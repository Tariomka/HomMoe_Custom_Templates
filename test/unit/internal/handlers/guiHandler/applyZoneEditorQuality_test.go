package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenQualityChanges_ReturnsServiceEquivalentMutation(t *testing.T) {
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
	zone := test_helpers.NewZoneEditorService().
		NewDefaultNeutralZone("Z", neutral_zone.QualityLow, 0, true, tuning)
	zones := []template_model.Zone{zone, {Name: "Player-A"}}
	connections := []template_model.Connection{
		{Name: "c1", From: "Player-A", To: zone.Name, GuardValue: 15000},
	}
	expectedZones, expectedConnections := test_helpers.NewZoneEditorService().
		ApplyNeutralZoneQualityEdit(models.NeutralZoneQualityEditRequest{
			Zones:           zones,
			Connections:     connections,
			PlayerZoneNames: []string{"Player-A"},
			ZoneName:        zone.Name,
			Quality:         neutral_zone.QualityHigh,
			CastleCount:     3,
			Tuning:          tuning,
		})

	// Act
	result := handler.ApplyZoneEditorQuality(dtos.ZoneEditorQualityRequestDto{
		Zone:            zone,
		Quality:         neutral_zone.QualityHigh,
		CastleCount:     3,
		Tuning:          tuning,
		Zones:           zones,
		Connections:     connections,
		PlayerZoneNames: []string{"Player-A"},
	})

	// Assert
	assert.Equal(
		t,
		dtos.ZoneEditorMutationDto{Zones: expectedZones, Connections: expectedConnections},
		result)
}
