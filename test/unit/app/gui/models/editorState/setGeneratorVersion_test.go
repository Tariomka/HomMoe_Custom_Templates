package editorState_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenAGeneratorVersionIsSet_CurrentStateCarriesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	state.SetGeneratorVersion(7)

	// Assert
	assert.Equal(t, 7, state.GetCurrentState().GeneratorVersion)
}

func TestWhenTheGeneratorVersionIsSetAgain_TheLatestValueWins(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetGeneratorVersion(1)

	// Act
	state.SetGeneratorVersion(2)

	// Assert
	assert.Equal(t, 2, state.GetGeneratorVersion())
}
