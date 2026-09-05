package manualZoneSave_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenSaveListIsEmpty_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var saves []editor_state.ManualZoneSave

	// Act
	zones := editor_state_model.ToManualZoneModels(saves)

	// Assert
	assert.Nil(t, zones)
}

func TestWhenSavesCarryManualPositions_RestoresEachPositionOntoZone(t *testing.T) {
	t.Parallel()
	// Arrange
	firstPosition := &[2]float64{0.1, 0.9}
	secondPosition := &[2]float64{0.6, 0.4}
	saves := []editor_state.ManualZoneSave{
		{Zone: entities.Zone{Name: "Zone A"}, ManualPosition: firstPosition},
		{Zone: entities.Zone{Name: "Zone B"}, ManualPosition: secondPosition},
	}
	expected := []template_model.Zone{
		{Name: "Zone A", ManualPosition: new(data.NewVec2(0.1, 0.9))},
		{Name: "Zone B", ManualPosition: new(data.NewVec2(0.6, 0.4))},
	}

	// Act
	zones := editor_state_model.ToManualZoneModels(saves)

	// Assert
	assert.Equal(t, expected, zones)
}

func TestWhenSaveRecordsThePlasticTier_RestoresItOntoTheZone(t *testing.T) {
	t.Parallel()
	// Arrange
	ordinal := int8(neutral_zone.QualityLowest)
	saves := []editor_state.ManualZoneSave{
		{Zone: entities.Zone{Name: "Neutral-C"}, Quality: &ordinal},
	}

	// Act
	zones := editor_state_model.ToManualZoneModels(saves)

	// Assert
	assert.Equal(t, neutral_zone.QualityLowest, *zones[0].Quality)
}

// A .gen.json written before the tier was persisted carries no quality, and
// the zone must fall back to inference rather than claim Plastic.
func TestWhenSaveCarriesNoQuality_LeavesTheTierUnrecorded(t *testing.T) {
	t.Parallel()
	// Arrange
	saves := []editor_state.ManualZoneSave{{Zone: entities.Zone{Name: "Neutral-C"}}}

	// Act
	zones := editor_state_model.ToManualZoneModels(saves)

	// Assert
	assert.Nil(t, zones[0].Quality)
}
