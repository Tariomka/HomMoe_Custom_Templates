package manualZoneEdit_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheEditIsCloned_CloneEqualsTheSource(t *testing.T) {
	t.Parallel()
	// Arrange
	edit := newPopulatedEdit()

	// Act
	clone := edit.Clone()

	// Assert
	assert.Equal(t, edit, clone)
}

func TestWhenAMainObjectElementOfTheCloneIsMutated_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	edit := newPopulatedEdit()
	clone := edit.Clone()

	// Act
	clone.FieldOverrides.MainObjects[0].Factions[0] = "mutated"

	// Assert
	assert.Equal(t, []string{"faction"}, edit.FieldOverrides.MainObjects[0].Factions)
}

func TestWhenAnOverrideListOfTheCloneIsMutated_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	edit := newPopulatedEdit()
	clone := edit.Clone()

	// Act
	clone.FieldOverrides.GuardedContentPool[0] = "mutated"

	// Assert
	assert.Equal(t, []string{"guarded"}, edit.FieldOverrides.GuardedContentPool)
}

func TestWhenTheSizePointerOfTheCloneIsMutated_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	edit := newPopulatedEdit()
	clone := edit.Clone()

	// Act
	*clone.FieldOverrides.Size = 99

	// Assert
	assert.InDelta(t, 1.5, *edit.FieldOverrides.Size, 0)
}

func TestWhenTheEncounterHolesOfTheCloneAreMutated_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	edit := newPopulatedEdit()
	clone := edit.Clone()

	// Act
	clone.FieldOverrides.EncounterHolesSettings.TwoHoleEncounters = 99

	// Assert
	assert.InDelta(t, 2.0, edit.FieldOverrides.EncounterHolesSettings.TwoHoleEncounters, 0)
}

func TestWhenARoadOfTheCloneIsMutated_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	edit := newPopulatedEdit()
	clone := edit.Clone()

	// Act
	*clone.FieldOverrides.Roads[0].Road = false

	// Assert
	assert.Equal(t, new(true), edit.FieldOverrides.Roads[0].Road)
}

func TestWhenAnUnsetEntryOfTheCloneIsMutated_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	edit := newPopulatedEdit()
	clone := edit.Clone()

	// Act
	clone.Unset[0] = "mutated"

	// Assert
	assert.Equal(t, []string{"layout"}, edit.Unset)
}

func TestWhenThePositionsOfTheCloneAreMutated_SourcePositionsAreUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	edit := newPopulatedEdit()
	clone := edit.Clone()

	// Act
	clone.ManualPosition.X = 99
	clone.GeneratorPosition.X = 99

	// Assert
	assert.InDeltaSlice(t, []float64{1.0, 4.0}, []float64{edit.GeneratorPosition.X, edit.ManualPosition.X}, 0)
}

func TestWhenTheRingOfTheCloneIsMutated_SourceRingIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	edit := newPopulatedEdit()
	clone := edit.Clone()

	// Act
	*clone.GeneratorRing = 99

	// Assert
	assert.Equal(t, 3, *edit.GeneratorRing)
}

func TestWhenTheQualityOfTheCloneIsMutated_SourceQualityIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	edit := newPopulatedEdit()
	clone := edit.Clone()

	// Act
	*clone.Quality = int8(neutral_zone.QualityLowest)

	// Assert
	assert.Equal(t, int8(neutral_zone.QualityHigh), *edit.Quality)
}

func TestWhenTheCastlesOfTheCloneAreMutated_SourceCastlesAreUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	edit := newPopulatedEdit()
	clone := edit.Clone()

	// Act
	*clone.Castles = 99

	// Assert
	assert.Equal(t, 2, *edit.Castles)
}

func TestWhenTheEditHasNothingRecorded_CloneKeepsNilFieldsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit

	// Act
	clone := edit.Clone()

	// Assert
	assert.Equal(t, edit, clone)
}
