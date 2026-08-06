package fileSystemHandler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenSaveTargetIsRequested_DelegatesToPathResolution(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	mocks.pathResolution.On("ResolveSaveTarget", "dir", "name", ".gen.json").Return("dir/name.gen.json", true)

	// Act
	handler.ResolveSaveTarget("dir", "name", ".gen.json")

	// Assert
	mocks.pathResolution.AssertCalled(t, "ResolveSaveTarget", "dir", "name", ".gen.json")
}

func TestWhenSaveTargetIsResolved_ReturnsItUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	mocks.pathResolution.On("ResolveSaveTarget", "dir", "name", ".gen.json").Return("dir/name.gen.json", true)

	// Act
	target, ok := handler.ResolveSaveTarget("dir", "name", ".gen.json")

	// Assert
	require.True(t, ok)
	assert.Equal(t, "dir/name.gen.json", target)
}

func TestWhenSaveTargetCannotBeResolved_ReportsNoTarget(t *testing.T) {
	t.Parallel()
	// Arrange
	handler, mocks := newHandlerWithMocks()
	mocks.pathResolution.On("ResolveSaveTarget", "dir", "", ".gen.json").Return("", false)

	// Act
	_, ok := handler.ResolveSaveTarget("dir", "", ".gen.json")

	// Assert
	assert.False(t, ok)
}
