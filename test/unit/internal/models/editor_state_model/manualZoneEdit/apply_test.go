package manualZoneEdit_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenAppliedToTheDerivedZone_RichEditedZoneIsReproduced(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	edited := newRichEditedZone()
	edit := editor_state_model.NewManualZoneEdit(derived, edited, false)

	// Act
	result := edit.Apply(derived)

	// Assert
	assert.Equal(t, edited, result)
}

func TestWhenRoadsAreAlwaysStored_AppliedZoneReproducesTheEditedZone(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	edited := newRichEditedZone()
	edit := editor_state_model.NewManualZoneEdit(derived, edited, true)

	// Act
	result := edit.Apply(derived)

	// Assert
	assert.Equal(t, edited, result)
}

func TestWhenEditedZoneHasNilRoadsAndRoadsAreAlwaysStored_AppliedZoneHasNilRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	edited := newRichEditedZone()
	edited.Roads = nil
	edit := editor_state_model.NewManualZoneEdit(derived, edited, true)

	// Act
	result := edit.Apply(derived)

	// Assert
	assert.Nil(t, result.Roads)
}

func TestWhenEditedZoneUnsetAField_AppliedZoneHasItNil(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	edited := derived.Clone()
	edited.GuardedContentPool = nil
	edit := editor_state_model.NewManualZoneEdit(derived, edited, false)

	// Act
	result := edit.Apply(derived)

	// Assert
	assert.Nil(t, result.GuardedContentPool)
}

func TestWhenAppliedToADifferentDerivedZone_StructureComesFromTheEdit(t *testing.T) {
	t.Parallel()
	// Arrange
	edit := editor_state_model.NewManualZoneEdit(newDerivedZone(), newPopulatedZone(), false)
	otherDerived := newRichEditedZone()

	// Act
	result := edit.Apply(otherDerived)

	// Assert
	assert.Equal(t, "Zone A", result.Name)
}

func TestWhenAppliedToADifferentDerivedZone_UntouchedFieldsComeFromDerived(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	edited := derived.Clone()
	edited.Size = derived.Size + 1
	edit := editor_state_model.NewManualZoneEdit(derived, edited, false)
	otherDerived := newRichEditedZone()

	// Act
	result := edit.Apply(otherDerived)

	// Assert
	assert.Equal(t, otherDerived.Layout, result.Layout)
}

func TestWhenTheEditHasNoQuality_AppliedZoneHasNilQuality(t *testing.T) {
	t.Parallel()
	// Arrange
	edited := newPopulatedZone()
	edited.Quality = nil
	edit := editor_state_model.NewManualZoneEdit(newDerivedZone(), edited, false)

	// Act
	result := edit.Apply(newDerivedZone())

	// Assert
	assert.Nil(t, result.Quality)
}

func TestWhenTheEditQualityIsTheLowestTier_AppliedZoneKeepsTheLowestTier(t *testing.T) {
	t.Parallel()
	// Arrange
	edited := newPopulatedZone()
	edited.Quality = new(neutral_zone.QualityLowest)
	edit := editor_state_model.NewManualZoneEdit(newDerivedZone(), edited, false)

	// Act
	result := edit.Apply(newDerivedZone())

	// Assert
	assert.Equal(t, new(neutral_zone.QualityLowest), result.Quality)
}

func TestWhenTheResultSliceIsMutated_ReapplyingYieldsTheOriginalSlice(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	edited := derived.Clone()
	edited.GuardedContentPool = []string{"original"}
	edit := editor_state_model.NewManualZoneEdit(derived, edited, false)
	edit.Apply(derived).GuardedContentPool[0] = "mutated"

	// Act
	result := edit.Apply(derived)

	// Assert
	assert.Equal(t, []string{"original"}, result.GuardedContentPool)
}

func TestWhenTheResultMainObjectIsMutated_ReapplyingYieldsTheOriginalMainObject(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	edit := editor_state_model.NewManualZoneEdit(derived, newPopulatedZone(), false)
	edit.Apply(derived).MainObjects[0].Factions[0] = "mutated"

	// Act
	result := edit.Apply(derived)

	// Assert
	assert.Equal(t, []string{"faction"}, result.MainObjects[0].Factions)
}

func TestWhenTheResultEncounterHolesAreMutated_ReapplyingYieldsTheOriginalSettings(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	edit := editor_state_model.NewManualZoneEdit(derived, newPopulatedZone(), false)
	edit.Apply(derived).EncounterHolesSettings.TwoHoleEncounters = 99

	// Act
	result := edit.Apply(derived)

	// Assert
	assert.InDelta(t, 2.0, result.EncounterHolesSettings.TwoHoleEncounters, 0)
}

func TestWhenTheResultPositionIsMutated_ReapplyingYieldsTheOriginalPosition(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	edit := editor_state_model.NewManualZoneEdit(derived, newPopulatedZone(), false)
	edit.Apply(derived).ManualPosition.X = 99

	// Act
	result := edit.Apply(derived)

	// Assert
	assert.InDelta(t, 4.0, result.ManualPosition.X, 0)
}

func TestWhenTheResultRoadIsMutated_ReapplyingYieldsTheOriginalRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	edit := editor_state_model.NewManualZoneEdit(derived, newPopulatedZone(), false)
	*edit.Apply(derived).Roads[0].Road = false

	// Act
	result := edit.Apply(derived)

	// Assert
	assert.Equal(t, new(true), result.Roads[0].Road)
}

func TestWhenTheResultSliceIsMutated_DerivedZoneIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	original := derived.Clone()
	edit := editor_state_model.NewManualZoneEdit(derived, derived, false)

	// Act
	edit.Apply(derived).GuardedContentPool[0] = "mutated"

	// Assert
	assert.Equal(t, original, derived)
}

func TestWhenTheEditIsAppliedTwice_ResultsDoNotShareSlices(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	edit := editor_state_model.NewManualZoneEdit(derived, newRichEditedZone(), false)
	first := edit.Apply(derived)
	second := edit.Apply(derived)

	// Act
	first.GuardedContentPool[0] = "mutated"

	// Assert
	assert.NotEqual(t, "mutated", second.GuardedContentPool[0])
}

func TestWhenTheEditIsAppliedToAnEmptyZone_PopulatedZoneIsReproduced(t *testing.T) {
	t.Parallel()
	// Arrange
	populated := newPopulatedZone()
	edit := editor_state_model.NewManualZoneEdit(template_model.Zone{}, populated, false)

	// Act
	result := edit.Apply(template_model.Zone{})

	// Assert
	assert.Equal(t, populated, result)
}
