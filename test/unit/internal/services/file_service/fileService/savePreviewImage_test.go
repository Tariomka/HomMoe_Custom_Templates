package fileService_test

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newPreviewImage() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, 16, 16))
}

func TestWhenPreviewNameContainsSlash_ReturnsPathWithSanitizedName(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()

	// Act
	writtenPath, err := file_service.NewFileService().SavePreviewImage(outputDir, newPreviewImage(), "My/Preview")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(outputDir, "My_Preview.png"), writtenPath)
}

func TestWhenPreviewIsWritten_CreatesNonEmptyFile(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	writtenPath, err := file_service.NewFileService().SavePreviewImage(outputDir, newPreviewImage(), "Preview")
	require.NoError(t, err)

	// Act
	fileInfo, statErr := os.Stat(writtenPath)

	// Assert
	require.NoError(t, statErr)
	assert.Positive(t, fileInfo.Size())
}

func TestWhenPreviewIsSaved_WritesDecodablePng(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	writtenPath, err := file_service.NewFileService().SavePreviewImage(outputDir, newPreviewImage(), "Preview")
	require.NoError(t, err)
	file, err := os.Open(writtenPath)
	require.NoError(t, err)
	defer file.Close()

	// Act
	_, decodeErr := png.Decode(file)

	// Assert
	assert.NoError(t, decodeErr)
}

func TestWhenPreviewNameIsEmpty_FallsBackToGeneratedTemplateFileName(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()

	// Act
	writtenPath, err := file_service.NewFileService().SavePreviewImage(outputDir, newPreviewImage(), "")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Generated_Template.png", filepath.Base(writtenPath))
}

func TestWhenPreviewDirectoryIsMissing_CreatesIt(t *testing.T) {
	// Arrange
	outputDir := filepath.Join(t.TempDir(), "nested", "preview")

	// Act
	_, err := file_service.NewFileService().SavePreviewImage(outputDir, newPreviewImage(), "T")

	// Assert
	require.NoError(t, err)
	assert.DirExists(t, outputDir)
}

func TestWhenPreviewParentPathIsAFile_ReturnsError(t *testing.T) {
	// Arrange
	blockerPath := filepath.Join(t.TempDir(), "blocker")
	require.NoError(t, os.WriteFile(blockerPath, []byte("x"), 0o644))
	outputDir := filepath.Join(blockerPath, "child")

	// Act
	_, err := file_service.NewFileService().SavePreviewImage(outputDir, newPreviewImage(), "T")

	// Assert
	assert.Error(t, err)
}

func TestWhenPreviewTargetPathIsADirectory_ReturnsError(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(outputDir, "T.png"), 0o755))

	// Act
	_, err := file_service.NewFileService().SavePreviewImage(outputDir, newPreviewImage(), "T")

	// Assert
	assert.Error(t, err)
}
