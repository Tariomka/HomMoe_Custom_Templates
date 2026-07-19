package snapshotStore_test

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sampleScreenshot builds a small deterministic image for snapshot IO tests.
func sampleScreenshot() *image.RGBA {
	screenshot := image.NewRGBA(image.Rect(0, 0, 3, 2))
	for row := range 2 {
		for column := range 3 {
			screenshot.SetRGBA(column, row, color.RGBA{R: uint8(40 * column), G: uint8(90 * row), B: 200, A: 255})
		}
	}
	return screenshot
}

func TestWhenDirectoriesMissing_CreatesThemAndWritesFile(t *testing.T) {
	t.Parallel()
	// Arrange
	store := integration_common.NewSnapshotStoreWithRoot(t.TempDir())
	goldenPath := store.GoldenPath("someFile", "SomeTest", 1)

	// Act
	err := store.SaveGolden(goldenPath, sampleScreenshot())

	// Assert
	require.NoError(t, err)
	assert.FileExists(t, goldenPath)
}

func TestWhenStaleFailureExists_RemovesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	store := integration_common.NewSnapshotStoreWithRoot(t.TempDir())
	goldenPath := store.GoldenPath("someFile", "SomeTest", 1)
	failurePath := store.FailurePath("someFile", "SomeTest", 1)
	require.NoError(t, store.SaveFailure(failurePath, sampleScreenshot()))

	// Act
	err := store.SaveGolden(goldenPath, sampleScreenshot())

	// Assert
	require.NoError(t, err)
	assert.NoFileExists(t, failurePath)
}

func TestWhenRootIsExistingFile_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	blockingFile := filepath.Join(t.TempDir(), "blocking")
	require.NoError(t, os.WriteFile(blockingFile, []byte("not a directory"), 0o644))
	store := integration_common.NewSnapshotStoreWithRoot(blockingFile)

	// Act
	err := store.SaveGolden(store.GoldenPath("someFile", "SomeTest", 1), sampleScreenshot())

	// Assert
	assert.Error(t, err)
}
