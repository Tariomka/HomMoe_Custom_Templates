package manualZoneSave_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneListIsEmpty_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var zones []template_model.Zone

	// Act
	saves := editor_state_model.ToManualZoneSaves(zones)

	// Assert
	assert.Nil(t, saves)
}

func TestWhenZonesHaveManualPositions_PreservesEachPositionInSave(t *testing.T) {
	t.Parallel()
	// Arrange
	firstPosition := &[2]float64{0.25, 0.75}
	secondPosition := &[2]float64{0.5, 0.5}
	zones := []template_model.Zone{
		{Name: "Zone A", ManualPosition: firstPosition},
		{Name: "Zone B", ManualPosition: secondPosition},
	}
	expected := []editor_state_model.ManualZoneSave{
		{Zone: template_model.ToZoneEntity(zones[0]), ManualPosition: firstPosition},
		{Zone: template_model.ToZoneEntity(zones[1]), ManualPosition: secondPosition},
	}

	// Act
	saves := editor_state_model.ToManualZoneSaves(zones)

	// Assert
	assert.Equal(t, expected, saves)
}

func TestWhenZoneHasNoManualPosition_SavesNilPosition(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{{Name: "Zone A"}}

	// Act
	saves := editor_state_model.ToManualZoneSaves(zones)

	// Assert
	assert.Nil(t, saves[0].ManualPosition)
}

// The Plastic tier is ordinal 0, so a value field with omitempty would drop it
// from the file entirely and the zone would load back as "never recorded".
func TestWhenZoneRecordsThePlasticTier_SavesItsOrdinal(t *testing.T) {
	t.Parallel()
	// Arrange
	quality := neutral_zone.QualityLowest
	zones := []template_model.Zone{{Name: "Neutral-C", Quality: &quality}}

	// Act
	saves := editor_state_model.ToManualZoneSaves(zones)

	// Assert
	assert.Equal(t, int8(neutral_zone.QualityLowest), *saves[0].Quality)
}

func TestWhenZoneRecordsNoTier_SavesNoQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{{Name: "Neutral-C"}}

	// Act
	saves := editor_state_model.ToManualZoneSaves(zones)

	// Assert
	assert.Nil(t, saves[0].Quality)
}
