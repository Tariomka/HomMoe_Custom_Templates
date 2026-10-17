package editorState_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenAGeneratorVersionIsSet_GetGeneratorVersionReturnsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.SetGeneratorVersion(4)

	// Act
	version := state.GetGeneratorVersion()

	// Assert
	assert.Equal(t, 4, version)
}

func TestWhenNoGeneratorVersionWasSet_GetGeneratorVersionReturnsZero(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()

	// Act
	version := state.GetGeneratorVersion()

	// Assert
	assert.Equal(t, 0, version)
}
