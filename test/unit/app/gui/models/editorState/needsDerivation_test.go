package editorState_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenEditsExistWithoutAGeneratorVersion_NeedsDerivationIsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetManualEdits(newPopulatedManualEdits())

	// Act
	needsDerivation := state.NeedsDerivation()

	// Assert
	assert.True(t, needsDerivation)
}

func TestWhenEditsExistWithAGeneratorVersion_NeedsDerivationIsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetManualEdits(newPopulatedManualEdits())
	state.SetGeneratorVersion(1)

	// Act
	needsDerivation := state.NeedsDerivation()

	// Assert
	assert.False(t, needsDerivation)
}

func TestWhenThereAreNoEdits_NeedsDerivationIsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	needsDerivation := state.NeedsDerivation()

	// Assert
	assert.False(t, needsDerivation)
}
