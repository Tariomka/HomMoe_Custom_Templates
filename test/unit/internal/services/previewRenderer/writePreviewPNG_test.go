package previewRenderer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenNameContainsSlash_ReturnsPathWithSanitizedName(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	rmgTemplate := simpleTemplate("My/Preview")

	// Act
	writtenPath, err := services.WritePreviewPNG(outputDir, rmgTemplate, config.TopologyRing)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(outputDir, "My_Preview.png"), writtenPath)
}

func TestWhenPreviewIsWritten_CreatesNonEmptyFile(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	rmgTemplate := simpleTemplate("Preview")
	writtenPath, err := services.WritePreviewPNG(outputDir, rmgTemplate, config.TopologyRing)
	require.NoError(t, err)

	// Act
	fileInfo, statErr := os.Stat(writtenPath)

	// Assert
	require.NoError(t, statErr)
	assert.Positive(t, fileInfo.Size())
}

func TestWhenNameIsEmpty_FallsBackToGeneratedTemplateFileName(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	rmgTemplate := simpleTemplate("")

	// Act
	writtenPath, err := services.WritePreviewPNG(outputDir, rmgTemplate, config.TopologyRing)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Generated_Template.png", filepath.Base(writtenPath))
}

func TestWhenTargetDirectoryIsMissing_CreatesIt(t *testing.T) {
	// Arrange
	outputDir := filepath.Join(t.TempDir(), "nested", "preview")
	rmgTemplate := simpleTemplate("T")

	// Act
	_, err := services.WritePreviewPNG(outputDir, rmgTemplate, config.TopologyRing)

	// Assert
	require.NoError(t, err)
	assert.DirExists(t, outputDir)
}

func TestWhenParentPathIsAFile_ReturnsError(t *testing.T) {
	// Arrange
	blockerPath := filepath.Join(t.TempDir(), "blocker")
	require.NoError(t, os.WriteFile(blockerPath, []byte("x"), 0o644))
	outputDir := filepath.Join(blockerPath, "child")
	rmgTemplate := simpleTemplate("T")

	// Act
	_, err := services.WritePreviewPNG(outputDir, rmgTemplate, config.TopologyRing)

	// Assert
	assert.Error(t, err)
}

func TestWhenTargetPathIsOccupiedByDirectory_ReturnsError(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(outputDir, "T.png"), 0o755))
	rmgTemplate := simpleTemplate("T")

	// Act
	_, err := services.WritePreviewPNG(outputDir, rmgTemplate, config.TopologyRing)

	// Assert
	assert.Error(t, err)
}
