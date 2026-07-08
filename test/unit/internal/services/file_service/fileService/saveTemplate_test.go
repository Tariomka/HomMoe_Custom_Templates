package fileService_test

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTemplateNameContainsSlash_ReturnsPathWithSanitizedName(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	rmgTemplate := &entities.RmgTemplate{Name: "My/Template"}

	// Act
	writtenPath, err := file_service.NewFileService().SaveTemplate(outputDir, rmgTemplate)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(outputDir, "My_Template.rmg.json"), writtenPath)
}

func TestWhenTemplateIsSaved_CreatesFileOnDisk(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	rmgTemplate := &entities.RmgTemplate{Name: "Plain"}

	// Act
	writtenPath, err := file_service.NewFileService().SaveTemplate(outputDir, rmgTemplate)

	// Assert
	require.NoError(t, err)
	assert.FileExists(t, writtenPath)
}

func TestWhenTemplateNameIsEmpty_FallsBackToGeneratedTemplateFileName(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	rmgTemplate := &entities.RmgTemplate{Name: ""}

	// Act
	writtenPath, err := file_service.NewFileService().SaveTemplate(outputDir, rmgTemplate)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Generated_Template.rmg.json", filepath.Base(writtenPath))
}

func TestWhenTemplateNameIsOnlyWhitespace_FallsBackToGeneratedTemplateFileName(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	rmgTemplate := &entities.RmgTemplate{Name: "   "}

	// Act
	writtenPath, err := file_service.NewFileService().SaveTemplate(outputDir, rmgTemplate)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Generated_Template.rmg.json", filepath.Base(writtenPath))
}

func TestWhenTemplateDirectoryIsMissing_CreatesIt(t *testing.T) {
	// Arrange
	outputDir := filepath.Join(t.TempDir(), "a", "b", "c")
	rmgTemplate := &entities.RmgTemplate{Name: "T"}

	// Act
	_, err := file_service.NewFileService().SaveTemplate(outputDir, rmgTemplate)

	// Assert
	require.NoError(t, err)
	assert.DirExists(t, outputDir)
}

func TestWhenTemplateIsSaved_ProducesIndentedJson(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	rmgTemplate := &entities.RmgTemplate{Name: "T", SizeX: 10}
	writtenPath, err := file_service.NewFileService().SaveTemplate(outputDir, rmgTemplate)
	require.NoError(t, err)

	// Act
	data, err := os.ReadFile(writtenPath)

	// Assert
	require.NoError(t, err)
	assert.Contains(t, string(data), "\n  ")
}

func TestWhenWrittenFileIsRead_ParsesBackIntoTemplate(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	rmgTemplate := &entities.RmgTemplate{Name: "T", SizeX: 10}
	writtenPath, err := file_service.NewFileService().SaveTemplate(outputDir, rmgTemplate)
	require.NoError(t, err)
	data, err := os.ReadFile(writtenPath)
	require.NoError(t, err)

	// Act
	var parsed entities.RmgTemplate
	parseErr := json.Unmarshal(data, &parsed)

	// Assert
	require.NoError(t, parseErr)
	assert.Equal(t, *rmgTemplate, parsed)
}

func TestWhenTemplateContainsNaNValue_ReturnsError(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	rmgTemplate := &entities.RmgTemplate{
		Name: "T",
		Variants: []entities.Variant{{
			Connections: []entities.Connection{{From: "A", To: "B", GuardWeeklyIncrement: math.NaN()}},
		}},
	}

	// Act
	_, err := file_service.NewFileService().SaveTemplate(outputDir, rmgTemplate)

	// Assert
	assert.Error(t, err)
}

func TestWhenTemplateParentPathIsAFile_ReturnsError(t *testing.T) {
	// Arrange
	blockerPath := filepath.Join(t.TempDir(), "blocker")
	require.NoError(t, os.WriteFile(blockerPath, []byte("x"), 0o644))
	outputDir := filepath.Join(blockerPath, "child")
	rmgTemplate := &entities.RmgTemplate{Name: "T"}

	// Act
	_, err := file_service.NewFileService().SaveTemplate(outputDir, rmgTemplate)

	// Assert
	assert.Error(t, err)
}

func TestWhenTargetPathIsOccupiedByDirectory_ReturnsError(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(outputDir, "T.rmg.json"), 0o755))
	rmgTemplate := &entities.RmgTemplate{Name: "T"}

	// Act
	_, err := file_service.NewFileService().SaveTemplate(outputDir, rmgTemplate)

	// Assert
	assert.Error(t, err)
}
