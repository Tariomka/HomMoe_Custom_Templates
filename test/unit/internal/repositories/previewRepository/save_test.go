package previewRepository_test

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newPreviewImage() image.RGBA {
	return *image.NewRGBA(image.Rect(0, 0, 16, 16))
}

// A zero-sized image is rejected by [png.Encode], which is the only encode
// failure reachable without a test-only seam.
func newUnencodablePreviewImage() image.RGBA {
	return image.RGBA{}
}

func TestWhenPreviewIsSaved_ReturnsPathWithPngExtension(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()

	// Act
	writtenPath, err := repositories.NewPreviewRepository().Save(outputDir, "My_Preview", newPreviewImage())

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(outputDir, "My_Preview.png"), writtenPath)
}

func TestWhenPreviewIsWritten_CreatesNonEmptyFile(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()
	writtenPath, err := repositories.NewPreviewRepository().Save(outputDir, "Preview", newPreviewImage())
	require.NoError(t, err)

	// Act
	fileInfo, statErr := os.Stat(writtenPath)

	// Assert
	require.NoError(t, statErr)
	assert.Positive(t, fileInfo.Size())
}

func TestWhenPreviewNameContainsInvalidCharacters_WritesUnderASanitizedName(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()

	// Act
	writtenPath, err := repositories.NewPreviewRepository().Save(outputDir, "a/b:c", newPreviewImage())

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(outputDir, "a_b_c.png"), writtenPath)
}

func TestWhenPreviewNameIsOnlyWhitespace_FallsBackToGeneratedTemplateFileName(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()

	// Act
	writtenPath, err := repositories.NewPreviewRepository().Save(outputDir, "   ", newPreviewImage())

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(outputDir, "Generated_Template.png"), writtenPath)
}

func TestWhenPreviewIsSaved_WritesDecodablePng(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()
	writtenPath, err := repositories.NewPreviewRepository().Save(outputDir, "Preview", newPreviewImage())
	require.NoError(t, err)
	file, err := os.Open(writtenPath)
	require.NoError(t, err)
	defer file.Close()

	// Act
	_, decodeErr := png.Decode(file)

	// Assert
	assert.NoError(t, decodeErr)
}

func TestWhenPreviewDirectoryIsMissing_CreatesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := filepath.Join(t.TempDir(), "nested", "preview")

	// Act
	_, err := repositories.NewPreviewRepository().Save(outputDir, "T", newPreviewImage())

	// Assert
	require.NoError(t, err)
	assert.DirExists(t, outputDir)
}

func TestWhenPreviewCannotBeEncoded_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()

	// Act
	_, err := repositories.NewPreviewRepository().Save(outputDir, "T", newUnencodablePreviewImage())

	// Assert
	assert.Error(t, err)
}

func TestWhenEncodingFailsOverAnExistingPreview_LeavesTheDestinationUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()
	destinationPath := filepath.Join(outputDir, "T.png")
	require.NoError(t, os.WriteFile(destinationPath, []byte("previous"), 0o644))

	// Act
	_, err := repositories.NewPreviewRepository().Save(outputDir, "T", newUnencodablePreviewImage())

	// Assert
	require.Error(t, err)
	survivingContent, readErr := os.ReadFile(destinationPath)
	require.NoError(t, readErr)
	assert.Equal(t, "previous", string(survivingContent))
}

func TestWhenPreviewEncodingFails_RemovesTheTemporaryFile(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()

	// Act
	_, err := repositories.NewPreviewRepository().Save(outputDir, "T", newUnencodablePreviewImage())

	// Assert
	require.Error(t, err)
	assert.NoFileExists(t, filepath.Join(outputDir, "TEMP-T.png"))
}

func TestWhenPreviewParentPathIsAFile_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	blockerPath := filepath.Join(t.TempDir(), "blocker")
	require.NoError(t, os.WriteFile(blockerPath, []byte("x"), 0o644))
	outputDir := filepath.Join(blockerPath, "child")

	// Act
	_, err := repositories.NewPreviewRepository().Save(outputDir, "T", newPreviewImage())

	// Assert
	assert.Error(t, err)
}

func TestWhenPreviewTargetPathIsADirectory_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(outputDir, "T.png"), 0o755))

	// Act
	_, err := repositories.NewPreviewRepository().Save(outputDir, "T", newPreviewImage())

	// Assert
	assert.Error(t, err)
}
