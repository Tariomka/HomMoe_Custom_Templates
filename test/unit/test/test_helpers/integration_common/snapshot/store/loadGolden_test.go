package store_test

import (
	"image"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common/snapshot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenGoldenExists_RoundTripsImagePixels(t *testing.T) {
	t.Parallel()
	// Arrange
	store := snapshot.NewStoreWithRoot(t.TempDir())
	goldenPath := store.GoldenPath("someFile", "SomeTest", 1)
	saved := sampleScreenshot()
	require.NoError(t, store.SaveGolden(goldenPath, saved))

	// Act
	loaded, err := store.LoadGolden(goldenPath)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, saved, loaded.(*image.RGBA))
}

func TestWhenGoldenIsMissing_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	store := snapshot.NewStoreWithRoot(t.TempDir())

	// Act
	_, err := store.LoadGolden(store.GoldenPath("someFile", "MissingTest", 1))

	// Assert
	assert.Error(t, err)
}
