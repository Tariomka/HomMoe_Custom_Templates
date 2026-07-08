package templateWriter_test

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenNameContainsSlash_ReturnsPathWithSanitizedName(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	rmgTemplate := &template.RmgTemplate{Name: "My/Template"}

	// Act
	writtenPath, err := services.WriteTemplate(outputDir, rmgTemplate)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(outputDir, "My_Template.rmg.json"), writtenPath)
}

func TestWhenTemplateIsWritten_CreatesFileOnDisk(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	rmgTemplate := &template.RmgTemplate{Name: "Plain"}

	// Act
	writtenPath, err := services.WriteTemplate(outputDir, rmgTemplate)

	// Assert
	require.NoError(t, err)
	assert.FileExists(t, writtenPath)
}

func TestWhenNameIsEmpty_FallsBackToGeneratedTemplateFileName(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	rmgTemplate := &template.RmgTemplate{Name: ""}

	// Act
	writtenPath, err := services.WriteTemplate(outputDir, rmgTemplate)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Generated_Template.rmg.json", filepath.Base(writtenPath))
}

func TestWhenNameIsOnlyWhitespace_FallsBackToGeneratedTemplateFileName(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	rmgTemplate := &template.RmgTemplate{Name: "   "}

	// Act
	writtenPath, err := services.WriteTemplate(outputDir, rmgTemplate)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Generated_Template.rmg.json", filepath.Base(writtenPath))
}

func TestWhenTargetDirectoryIsMissing_CreatesIt(t *testing.T) {
	// Arrange
	outputDir := filepath.Join(t.TempDir(), "a", "b", "c")
	rmgTemplate := &template.RmgTemplate{Name: "T"}

	// Act
	_, err := services.WriteTemplate(outputDir, rmgTemplate)

	// Assert
	require.NoError(t, err)
	assert.DirExists(t, outputDir)
}

func TestWhenTemplateIsWritten_ProducesIndentedJson(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	rmgTemplate := &template.RmgTemplate{Name: "T", SizeX: 10}
	writtenPath, err := services.WriteTemplate(outputDir, rmgTemplate)
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
	rmgTemplate := &template.RmgTemplate{Name: "T", SizeX: 10}
	writtenPath, err := services.WriteTemplate(outputDir, rmgTemplate)
	require.NoError(t, err)
	data, err := os.ReadFile(writtenPath)
	require.NoError(t, err)

	// Act
	var parsed template.RmgTemplate
	parseErr := json.Unmarshal(data, &parsed)

	// Assert
	assert.NoError(t, parseErr)
}

func TestWhenParentPathIsAFile_ReturnsError(t *testing.T) {
	// Arrange
	blockerPath := filepath.Join(t.TempDir(), "blocker")
	require.NoError(t, os.WriteFile(blockerPath, []byte("x"), 0o644))
	outputDir := filepath.Join(blockerPath, "child")
	rmgTemplate := &template.RmgTemplate{Name: "T"}

	// Act
	_, err := services.WriteTemplate(outputDir, rmgTemplate)

	// Assert
	assert.Error(t, err)
}

func TestWhenTemplateContainsNaNValue_ReturnsError(t *testing.T) {
	// Arrange
	outputDir := t.TempDir()
	rmgTemplate := &template.RmgTemplate{
		Name: "T",
		Variants: []entities.Variant{{
			Connections: []entities.Connection{{From: "A", To: "B", GuardWeeklyIncrement: math.NaN()}},
		}},
	}

	// Act
	_, err := services.WriteTemplate(outputDir, rmgTemplate)

	// Assert
	assert.Error(t, err)
}
