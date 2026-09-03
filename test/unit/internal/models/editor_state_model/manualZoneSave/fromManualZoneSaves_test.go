package manualZoneSave_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenSaveListIsEmpty_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var saves []editor_state_model.ManualZoneSave

	// Act
	zones := editor_state_model.FromManualZoneSaves(saves)

	// Assert
	assert.Nil(t, zones)
}

func TestWhenSavesCarryManualPositions_RestoresEachPositionOntoZone(t *testing.T) {
	t.Parallel()
	// Arrange
	firstPosition := &[2]float64{0.1, 0.9}
	secondPosition := &[2]float64{0.6, 0.4}
	saves := []editor_state_model.ManualZoneSave{
		{Zone: entities.Zone{Name: "Zone A"}, ManualPosition: firstPosition},
		{Zone: entities.Zone{Name: "Zone B"}, ManualPosition: secondPosition},
	}
	expected := []template_model.Zone{
		{Name: "Zone A", ManualPosition: firstPosition},
		{Name: "Zone B", ManualPosition: secondPosition},
	}

	// Act
	zones := editor_state_model.FromManualZoneSaves(saves)

	// Assert
	assert.Equal(t, expected, zones)
}

func TestWhenSavePositionDiffersFromEmbeddedZonePosition_SavePositionWins(t *testing.T) {
	t.Parallel()
	// Arrange
	savedPosition := &[2]float64{0.2, 0.3}
	staleEmbeddedPosition := &[2]float64{0.8, 0.8}
	saves := []editor_state_model.ManualZoneSave{
		{
			Zone:           entities.Zone{Name: "Zone A", ManualPosition: staleEmbeddedPosition},
			ManualPosition: savedPosition,
		},
	}

	// Act
	zones := editor_state_model.FromManualZoneSaves(saves)

	// Assert
	assert.Same(t, savedPosition, zones[0].ManualPosition)
}

func TestWhenSaveRecordsThePlasticTier_RestoresItOntoTheZone(t *testing.T) {
	t.Parallel()
	// Arrange
	ordinal := int8(neutral_zone.QualityLowest)
	saves := []editor_state_model.ManualZoneSave{
		{Zone: entities.Zone{Name: "Neutral-C"}, Quality: &ordinal},
	}

	// Act
	zones := editor_state_model.FromManualZoneSaves(saves)

	// Assert
	assert.Equal(t, neutral_zone.QualityLowest, *zones[0].Quality)
}

// A .gen.json written before the tier was persisted carries no quality, and
// the zone must fall back to inference rather than claim Plastic.
func TestWhenSaveCarriesNoQuality_LeavesTheTierUnrecorded(t *testing.T) {
	t.Parallel()
	// Arrange
	saves := []editor_state_model.ManualZoneSave{{Zone: entities.Zone{Name: "Neutral-C"}}}

	// Act
	zones := editor_state_model.FromManualZoneSaves(saves)

	// Assert
	assert.Nil(t, zones[0].Quality)
}
