package manualZoneEdit_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenZonesAreEqual_NoFieldIsOverridden(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := newDerivedZone()

	// Act
	edit := editor_state_model.NewManualZoneEdit(zone, zone, false)

	// Assert
	assert.Equal(t, editor_state.ZoneFieldOverrides{}, edit.FieldOverrides)
}

func TestWhenZonesAreEqual_UnsetIsEmpty(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := newDerivedZone()

	// Act
	edit := editor_state_model.NewManualZoneEdit(zone, zone, false)

	// Assert
	assert.Empty(t, edit.Unset)
}

func TestWhenOnlyAScalarDiffers_ThatScalarIsStoredAsAPointerOverride(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	edited := derived.Clone()
	edited.Size = derived.Size + 2

	// Act
	edit := editor_state_model.NewManualZoneEdit(derived, edited, false)

	// Assert
	assert.Equal(t, new(edited.Size), edit.FieldOverrides.Size)
}

func TestWhenOnlyAScalarDiffers_OtherFieldsStayUnset(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	edited := derived.Clone()
	edited.Size = derived.Size + 2

	// Act
	edit := editor_state_model.NewManualZoneEdit(derived, edited, false)

	// Assert
	assert.Nil(t, edit.FieldOverrides.Layout)
}

func TestWhenOnlyAScalarDiffers_UnsetIsEmpty(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	edited := derived.Clone()
	edited.Size = derived.Size + 2

	// Act
	edit := editor_state_model.NewManualZoneEdit(derived, edited, false)

	// Assert
	assert.Empty(t, edit.Unset)
}

func TestWhenEditedSliceIsNilAndDerivedHasValues_FieldNameIsListedInUnset(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	edited := derived.Clone()
	edited.GuardedContentPool = nil

	// Act
	edit := editor_state_model.NewManualZoneEdit(derived, edited, false)

	// Assert
	assert.Equal(t, []string{"guardedContentPool"}, edit.Unset)
}

func TestWhenEditedSliceIsNilAndDerivedHasValues_NoOverrideIsStored(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	edited := derived.Clone()
	edited.GuardedContentPool = nil

	// Act
	edit := editor_state_model.NewManualZoneEdit(derived, edited, false)

	// Assert
	assert.Nil(t, edit.FieldOverrides.GuardedContentPool)
}

func TestWhenEditedSliceIsEmptyAndDerivedHasValues_EmptyOverrideIsStored(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	edited := derived.Clone()
	edited.GuardedContentPool = []string{}

	// Act
	edit := editor_state_model.NewManualZoneEdit(derived, edited, false)

	// Assert
	assert.Equal(t, []string{}, edit.FieldOverrides.GuardedContentPool)
}

func TestWhenEditedSliceIsEmptyAndDerivedHasValues_FieldIsNotUnset(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	edited := derived.Clone()
	edited.GuardedContentPool = []string{}

	// Act
	edit := editor_state_model.NewManualZoneEdit(derived, edited, false)

	// Assert
	assert.Empty(t, edit.Unset)
}

func TestWhenEditedNumberIsZeroAndDerivedIsNot_ZeroOverrideIsStored(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	derived.GuardCutoffValue = 5
	edited := derived.Clone()
	edited.GuardCutoffValue = 0

	// Act
	edit := editor_state_model.NewManualZoneEdit(derived, edited, false)

	// Assert
	assert.Equal(t, new(0), edit.FieldOverrides.GuardCutoffValue)
}

func TestWhenEditedBoolIsFalseAndDerivedIsTrue_FalseOverrideIsStored(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	derived.RandomHireEnableWeeklyUnitIncrement = []bool{true}
	edited := derived.Clone()
	edited.RandomHireEnableWeeklyUnitIncrement = []bool{false}

	// Act
	edit := editor_state_model.NewManualZoneEdit(derived, edited, false)

	// Assert
	assert.Equal(t, []bool{false}, edit.FieldOverrides.RandomHireEnableWeeklyUnitIncrement)
}

func TestWhenEditedHasAName_NameIsCopied(t *testing.T) {
	t.Parallel()
	// Arrange
	edited := newPopulatedZone()

	// Act
	edit := editor_state_model.NewManualZoneEdit(newDerivedZone(), edited, false)

	// Assert
	assert.Equal(t, "Zone A", edit.Name)
}

func TestWhenEditedHasAQuality_QualityOrdinalIsStored(t *testing.T) {
	t.Parallel()
	// Arrange
	edited := newPopulatedZone()

	// Act
	edit := editor_state_model.NewManualZoneEdit(newDerivedZone(), edited, false)

	// Assert
	assert.Equal(t, new(int8(neutral_zone.QualityHigh)), edit.Quality)
}

func TestWhenEditedQualityIsTheLowestTier_QualityIsStoredAsZeroNotNil(t *testing.T) {
	t.Parallel()
	// Arrange
	edited := newPopulatedZone()
	edited.Quality = new(neutral_zone.QualityLowest)

	// Act
	edit := editor_state_model.NewManualZoneEdit(newDerivedZone(), edited, false)

	// Assert
	assert.Equal(t, new(int8(0)), edit.Quality)
}

func TestWhenEditedHasNoQuality_QualityIsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	edited := newPopulatedZone()
	edited.Quality = nil

	// Act
	edit := editor_state_model.NewManualZoneEdit(newDerivedZone(), edited, false)

	// Assert
	assert.Nil(t, edit.Quality)
}

func TestWhenEditedHasPositionsAndRing_GeneratorPositionIsCopied(t *testing.T) {
	t.Parallel()
	// Arrange
	edited := newPopulatedZone()

	// Act
	edit := editor_state_model.NewManualZoneEdit(newDerivedZone(), edited, false)

	// Assert
	assert.Equal(t, new(data.NewVec2(1.0, 2.0)), edit.GeneratorPosition)
}

func TestWhenEditedHasPositionsAndRing_ManualPositionIsCopied(t *testing.T) {
	t.Parallel()
	// Arrange
	edited := newPopulatedZone()

	// Act
	edit := editor_state_model.NewManualZoneEdit(newDerivedZone(), edited, false)

	// Assert
	assert.Equal(t, new(data.NewVec2(4.0, 5.0)), edit.ManualPosition)
}

func TestWhenEditedHasPositionsAndRing_GeneratorRingIsCopied(t *testing.T) {
	t.Parallel()
	// Arrange
	edited := newPopulatedZone()

	// Act
	edit := editor_state_model.NewManualZoneEdit(newDerivedZone(), edited, false)

	// Assert
	assert.Equal(t, new(3), edit.GeneratorRing)
}

func TestWhenEditedHasNoPositions_PositionsAreNil(t *testing.T) {
	t.Parallel()
	// Arrange
	edited := newPopulatedZone()
	edited.GeneratorPosition = nil
	edited.ManualPosition = nil

	// Act
	edit := editor_state_model.NewManualZoneEdit(newDerivedZone(), edited, false)

	// Assert
	assert.Nil(t, edit.ManualPosition)
}

func TestWhenEditedPositionIsMutatedAfterwards_StoredPositionIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	edited := newPopulatedZone()
	edit := editor_state_model.NewManualZoneEdit(newDerivedZone(), edited, false)

	// Act
	edited.ManualPosition.X = 99

	// Assert
	assert.InDelta(t, 4.0, edit.ManualPosition.X, 0)
}

func TestWhenEditedSliceOverrideIsMutatedAfterwards_StoredOverrideIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newDerivedZone()
	edited := derived.Clone()
	edited.GuardedContentPool = []string{"original"}
	edit := editor_state_model.NewManualZoneEdit(derived, edited, false)

	// Act
	edited.GuardedContentPool[0] = "mutated"

	// Assert
	assert.Equal(t, []string{"original"}, edit.FieldOverrides.GuardedContentPool)
}

func TestWhenEditIsBuilt_OriginIsNotRecorded(t *testing.T) {
	t.Parallel()
	// Arrange
	edited := newPopulatedZone()

	// Act
	edit := editor_state_model.NewManualZoneEdit(newDerivedZone(), edited, false)

	// Assert
	assert.False(t, edit.IsAdded())
}

func TestWhenRoadsAreEqualAndAlwaysStoreRoadsIsOff_RoadsAreNotStored(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := newPopulatedZone()

	// Act
	edit := editor_state_model.NewManualZoneEdit(zone, zone, false)

	// Assert
	assert.Nil(t, edit.FieldOverrides.Roads)
}

func TestWhenRoadsAreEqualAndAlwaysStoreRoadsIsOn_RoadsAreStored(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := newPopulatedZone()

	// Act
	edit := editor_state_model.NewManualZoneEdit(zone, zone, true)

	// Assert
	assert.Equal(t, template_model.ToZoneEntity(zone).Roads, edit.FieldOverrides.Roads)
}

func TestWhenRoadsAreEqualAndAlwaysStoreRoadsIsOn_RoadsAreNotUnset(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := newPopulatedZone()

	// Act
	edit := editor_state_model.NewManualZoneEdit(zone, zone, true)

	// Assert
	assert.Empty(t, edit.Unset)
}

func TestWhenRoadsDifferAndAlwaysStoreRoadsIsOn_EditedRoadsAreStored(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newPopulatedZone()
	edited := derived.Clone()
	edited.Roads = []template_model.Road{{Type: "other"}}

	// Act
	edit := editor_state_model.NewManualZoneEdit(derived, edited, true)

	// Assert
	assert.Equal(t, template_model.ToZoneEntity(edited).Roads, edit.FieldOverrides.Roads)
}

func TestWhenRoadsAreNilAndAlwaysStoreRoadsIsOnAndDerivedHasNone_RoadsAreListedInUnsetOnce(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := newPopulatedZone()
	zone.Roads = nil

	// Act
	edit := editor_state_model.NewManualZoneEdit(zone, zone, true)

	// Assert
	assert.Equal(t, []string{"roads"}, edit.Unset)
}

func TestWhenRoadsAreNilAndAlwaysStoreRoadsIsOnAndDerivedHasRoads_RoadsAreListedInUnsetOnce(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newPopulatedZone()
	edited := derived.Clone()
	edited.Roads = nil

	// Act
	edit := editor_state_model.NewManualZoneEdit(derived, edited, true)

	// Assert
	assert.Equal(t, []string{"roads"}, edit.Unset)
}

func TestWhenRoadsAreNilAndAlwaysStoreRoadsIsOnAndAnotherFieldIsUnset_BothAreListed(t *testing.T) {
	t.Parallel()
	// Arrange
	derived := newPopulatedZone()
	edited := derived.Clone()
	edited.Roads = nil
	edited.GuardedContentPool = nil

	// Act
	edit := editor_state_model.NewManualZoneEdit(derived, edited, true)

	// Assert
	assert.ElementsMatch(t, []string{"roads", "guardedContentPool"}, edit.Unset)
}

func TestWhenRoadsAreNilInBothZonesAndAlwaysStoreRoadsIsOff_RoadsAreNotUnset(t *testing.T) {
	t.Parallel()
	// Arrange
	zone := newPopulatedZone()
	zone.Roads = nil

	// Act
	edit := editor_state_model.NewManualZoneEdit(zone, zone, false)

	// Assert
	assert.Empty(t, edit.Unset)
}
