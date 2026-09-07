package manualZoneSave_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
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
	saves := []editor_state.ManualZoneSave{
		{Zone: template_entity.Zone{Name: "Zone A"}, ManualPosition: new(data.NewVec2(0.1, 0.9))},
		{Zone: template_entity.Zone{Name: "Zone B"}, ManualPosition: new(data.NewVec2(0.6, 0.4))},
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
		{Zone: template_entity.Zone{Name: "Neutral-C"}, Quality: &ordinal},
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
	saves := []editor_state.ManualZoneSave{{Zone: template_entity.Zone{Name: "Neutral-C"}}}

	// Act
	zones := editor_state_model.ToManualZoneModels(saves)

	// Assert
	assert.Nil(t, zones[0].Quality)
}

func TestWhenSaveCarriesAGeneratorPosition_RestoresItOntoTheZone(t *testing.T) {
	t.Parallel()
	// Arrange
	position := data.NewVec2(0.4, 0.6)
	saves := []editor_state.ManualZoneSave{
		{Zone: template_entity.Zone{Name: "Zone A"}, GeneratorPosition: &position},
	}

	// Act
	zones := editor_state_model.ToManualZoneModels(saves)

	// Assert
	assert.Equal(t, position, *zones[0].GeneratorPosition)
}

func TestWhenSaveCarriesTheInnermostGeneratorRing_RestoresItOntoTheZone(t *testing.T) {
	t.Parallel()
	// Arrange
	saves := []editor_state.ManualZoneSave{
		{Zone: template_entity.Zone{Name: "Zone A"}, GeneratorRing: new(0)},
	}

	// Act
	zones := editor_state_model.ToManualZoneModels(saves)

	// Assert
	assert.Equal(t, 0, *zones[0].GeneratorRing)
}

// A .gen.json written before schema v2 carries no stamps, and the preview has
// to fall back to a computed layout rather than believe the zone sits at 0,0.
func TestWhenSaveCarriesNoGeneratorPosition_LeavesTheZoneUnstamped(t *testing.T) {
	t.Parallel()
	// Arrange
	saves := []editor_state.ManualZoneSave{{Zone: template_entity.Zone{Name: "Zone A"}}}

	// Act
	zones := editor_state_model.ToManualZoneModels(saves)

	// Assert
	assert.Nil(t, zones[0].GeneratorPosition)
}

func TestWhenSaveCarriesNoGeneratorRing_LeavesTheZoneUnstamped(t *testing.T) {
	t.Parallel()
	// Arrange
	saves := []editor_state.ManualZoneSave{{Zone: template_entity.Zone{Name: "Zone A"}}}

	// Act
	zones := editor_state_model.ToManualZoneModels(saves)

	// Assert
	assert.Nil(t, zones[0].GeneratorRing)
}
