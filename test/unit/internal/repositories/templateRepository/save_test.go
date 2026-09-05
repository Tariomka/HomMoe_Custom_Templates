package templateRepository_test

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTemplateIsSaved_ReturnsPathWithRmgJsonExtension(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()

	// Act
	writtenPath, err := repositories.NewTemplateRepository().Save(outputDir, "My_Template", entities.RmgTemplate{})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(outputDir, "My_Template.rmg.json"), writtenPath)
}

func TestWhenTemplateIsSaved_CreatesFileOnDisk(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()

	// Act
	writtenPath, err := repositories.NewTemplateRepository().Save(outputDir, "Plain", entities.RmgTemplate{})

	// Assert
	require.NoError(t, err)
	assert.FileExists(t, writtenPath)
}

func TestWhenTemplateNameContainsInvalidCharacters_WritesUnderASanitizedName(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()

	// Act
	writtenPath, err := repositories.NewTemplateRepository().Save(outputDir, "a/b:c", entities.RmgTemplate{})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(outputDir, "a_b_c.rmg.json"), writtenPath)
}

func TestWhenTemplateNameIsOnlyWhitespace_FallsBackToGeneratedTemplateFileName(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()

	// Act
	writtenPath, err := repositories.NewTemplateRepository().Save(outputDir, "   ", entities.RmgTemplate{})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(outputDir, "Generated_Template.rmg.json"), writtenPath)
}

func TestWhenTemplateDirectoryIsMissing_CreatesIt(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := filepath.Join(t.TempDir(), "a", "b", "c")

	// Act
	_, err := repositories.NewTemplateRepository().Save(outputDir, "T", entities.RmgTemplate{})

	// Assert
	require.NoError(t, err)
	assert.DirExists(t, outputDir)
}

func TestWhenTemplateIsSaved_ProducesIndentedJson(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()
	writtenPath, err := repositories.NewTemplateRepository().Save(
		outputDir, "T", entities.RmgTemplate{Name: "T", SizeX: 10})
	require.NoError(t, err)

	// Act
	data, err := os.ReadFile(writtenPath)

	// Assert
	require.NoError(t, err)
	assert.Contains(t, string(data), "\n  ")
}

func TestWhenWrittenFileIsRead_ParsesBackIntoTemplate(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()
	rmgTemplate := entities.RmgTemplate{Name: "T", SizeX: 10}
	writtenPath, err := repositories.NewTemplateRepository().Save(outputDir, "T", rmgTemplate)
	require.NoError(t, err)
	data, err := os.ReadFile(writtenPath)
	require.NoError(t, err)

	// Act
	var parsed entities.RmgTemplate
	parseErr := json.Unmarshal(data, &parsed)

	// Assert
	require.NoError(t, parseErr)
	assert.Equal(t, rmgTemplate, parsed)
}

func TestWhenTemplateContainsNaNValue_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()

	// Act
	_, err := repositories.NewTemplateRepository().Save(outputDir, "T", newTemplateWithNaN())

	// Assert
	assert.Error(t, err)
}

func TestWhenEncodingFailsOverAnExistingTemplate_LeavesTheDestinationUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()
	destinationPath := filepath.Join(outputDir, "T.rmg.json")
	require.NoError(t, os.WriteFile(destinationPath, []byte("previous"), 0o644))

	// Act
	_, err := repositories.NewTemplateRepository().Save(outputDir, "T", newTemplateWithNaN())

	// Assert
	require.Error(t, err)
	survivingContent, readErr := os.ReadFile(destinationPath)
	require.NoError(t, readErr)
	assert.Equal(t, "previous", string(survivingContent))
}

func TestWhenEncodingFails_RemovesTheTemporaryFile(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()

	// Act
	_, err := repositories.NewTemplateRepository().Save(outputDir, "T", newTemplateWithNaN())

	// Assert
	require.Error(t, err)
	assert.NoFileExists(t, filepath.Join(outputDir, "TEMP-T.rmg.json"))
}

func TestWhenTemplateParentPathIsAFile_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	blockerPath := filepath.Join(t.TempDir(), "blocker")
	require.NoError(t, os.WriteFile(blockerPath, []byte("x"), 0o644))
	outputDir := filepath.Join(blockerPath, "child")

	// Act
	_, err := repositories.NewTemplateRepository().Save(outputDir, "T", entities.RmgTemplate{})

	// Assert
	assert.Error(t, err)
}

func TestWhenTargetPathIsOccupiedByDirectory_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	outputDir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(outputDir, "T.rmg.json"), 0o755))

	// Act
	_, err := repositories.NewTemplateRepository().Save(outputDir, "T", entities.RmgTemplate{})

	// Assert
	assert.Error(t, err)
}

func newTemplateWithNaN() entities.RmgTemplate {
	return entities.RmgTemplate{
		Name: "T",
		Variants: []entities.Variant{{
			Connections: []entities.Connection{{From: "A", To: "B", GuardWeeklyIncrement: math.NaN()}},
		}},
	}
}
