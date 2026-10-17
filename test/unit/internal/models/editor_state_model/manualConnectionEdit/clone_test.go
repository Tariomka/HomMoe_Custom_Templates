package manualConnectionEdit_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
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

func TestWhenARuleOfTheCloneIsMutated_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	edit := newPopulatedEdit()
	clone := edit.Clone()

	// Act
	clone.FieldOverrides.PortalPlacementRulesFrom[0].Type = "mutated"

	// Assert
	assert.Equal(t, "edited", edit.FieldOverrides.PortalPlacementRulesFrom[0].Type)
}

func TestWhenARuleArgumentOfTheCloneIsMutated_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	edit := newPopulatedEdit()
	clone := edit.Clone()

	// Act
	clone.FieldOverrides.PortalPlacementRulesTo[0].Args[0] = "mutated"

	// Assert
	assert.Equal(t, []any{"arg"}, edit.FieldOverrides.PortalPlacementRulesTo[0].Args)
}

func TestWhenTheRoadFlagOfTheCloneIsMutated_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	edit := newPopulatedEdit()
	clone := edit.Clone()

	// Act
	*clone.FieldOverrides.Road = true

	// Assert
	assert.Equal(t, new(false), edit.FieldOverrides.Road)
}

func TestWhenTheLengthPointerOfTheCloneIsMutated_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	edit := newPopulatedEdit()
	clone := edit.Clone()

	// Act
	*clone.FieldOverrides.Length = 99

	// Assert
	assert.Equal(t, new(2.5), edit.FieldOverrides.Length)
}

func TestWhenAnUnsetEntryOfTheCloneIsMutated_SourceIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	edit := newPopulatedEdit()
	clone := edit.Clone()

	// Act
	clone.Unset[0] = "mutated"

	// Assert
	assert.Equal(t, []string{"guardZone"}, edit.Unset)
}

func TestWhenTheEditHasNothingRecorded_CloneKeepsNilFieldsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualConnectionEdit

	// Act
	clone := edit.Clone()

	// Assert
	assert.Equal(t, edit, clone)
}
