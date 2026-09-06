package manualZoneSave_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
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
	saves := editor_state_model.ToManualZoneSaveEntities(zones)

	// Assert
	assert.Nil(t, saves)
}

func TestWhenZonesHaveManualPositions_PreservesEachPositionInSave(t *testing.T) {
	t.Parallel()
	// Arrange
	firstPosition := data.NewVec2(0.25, 0.75)
	secondPosition := data.NewVec2(0.5, 0.5)
	zones := []template_model.Zone{
		{Name: "Zone A", ManualPosition: &firstPosition},
		{Name: "Zone B", ManualPosition: &secondPosition},
	}
	expected := []editor_state.ManualZoneSave{
		{Zone: template_model.ToZoneEntity(zones[0]), ManualPosition: new(data.NewVec2(0.25, 0.75))},
		{Zone: template_model.ToZoneEntity(zones[1]), ManualPosition: new(data.NewVec2(0.5, 0.5))},
	}

	// Act
	saves := editor_state_model.ToManualZoneSaveEntities(zones)

	// Assert
	assert.Equal(t, expected, saves)
}

func TestWhenZoneHasNoManualPosition_SavesNilPosition(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{{Name: "Zone A"}}

	// Act
	saves := editor_state_model.ToManualZoneSaveEntities(zones)

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
	saves := editor_state_model.ToManualZoneSaveEntities(zones)

	// Assert
	assert.Equal(t, int8(neutral_zone.QualityLowest), *saves[0].Quality)
}

func TestWhenZoneRecordsNoTier_SavesNoQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{{Name: "Neutral-C"}}

	// Act
	saves := editor_state_model.ToManualZoneSaveEntities(zones)

	// Assert
	assert.Nil(t, saves[0].Quality)
}

func TestWhenZoneCarriesAGeneratorPosition_PreservesItInTheSave(t *testing.T) {
	t.Parallel()
	// Arrange
	position := data.NewVec2(0.4, 0.6)
	zones := []template_model.Zone{{Name: "Zone A", GeneratorPosition: &position}}

	// Act
	saves := editor_state_model.ToManualZoneSaveEntities(zones)

	// Assert
	assert.Equal(t, position, *saves[0].GeneratorPosition)
}

// Ring 0 is the innermost ring, not "unstamped", so it has to survive as a
// pointer all the way onto the save.
func TestWhenZoneSitsOnTheInnermostGeneratorRing_PreservesTheRingInTheSave(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{{Name: "Zone A", GeneratorRing: new(0)}}

	// Act
	saves := editor_state_model.ToManualZoneSaveEntities(zones)

	// Assert
	assert.Equal(t, 0, *saves[0].GeneratorRing)
}

func TestWhenZoneWasNeverStamped_SavesNoGeneratorPosition(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{{Name: "Zone A"}}

	// Act
	saves := editor_state_model.ToManualZoneSaveEntities(zones)

	// Assert
	assert.Nil(t, saves[0].GeneratorPosition)
}

func TestWhenZoneWasNeverStamped_SavesNoGeneratorRing(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{{Name: "Zone A"}}

	// Act
	saves := editor_state_model.ToManualZoneSaveEntities(zones)

	// Assert
	assert.Nil(t, saves[0].GeneratorRing)
}

// The save is a snapshot: dragging the zone afterwards must not rewrite it.
func TestWhenTheZoneIsDraggedAfterSaving_TheSavedPositionIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	position := data.NewVec2(0.25, 0.75)
	zones := []template_model.Zone{{Name: "Zone A", ManualPosition: &position}}
	saves := editor_state_model.ToManualZoneSaveEntities(zones)

	// Act
	position.X = 0.9

	// Assert
	assert.Equal(t, data.NewVec2(0.25, 0.75), *saves[0].ManualPosition)
}
