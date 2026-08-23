package guiHandler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTemplateToSaveIsNil_ReturnsNothingToSaveError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	templateDto := dtos.TemplateSaveDto{Template: nil, OutputPath: t.TempDir()}

	// Act
	_, err := handler.SaveTemplate(templateDto)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNothingToSave)
}

func TestWhenTemplateOutputPathIsEmpty_ReturnsNoOutputPathError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	templateDto := dtos.TemplateSaveDto{
		Template:   &entities.RmgTemplate{Name: "Empty Path Template"},
		OutputPath: "",
	}

	// Act
	_, err := handler.SaveTemplate(templateDto)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoOutputPath)
}

func TestWhenTemplateOutputPathIsWhitespaceOnly_ReturnsNoOutputPathError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	templateDto := dtos.TemplateSaveDto{
		Template:   &entities.RmgTemplate{Name: "Whitespace Path Template"},
		OutputPath: "   \t ",
	}

	// Act
	_, err := handler.SaveTemplate(templateDto)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoOutputPath)
}

func TestWhenTemplateAndOutputPathAreValid_ReturnsTemplateFilePath(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	outputDirectory := t.TempDir()
	template := generateDefaultTemplate(t, handler)
	template.Name = "Valid Save Template"
	templateDto := dtos.TemplateSaveDto{
		Template:   template,
		Topology:   config.TopologyRing,
		OutputPath: outputDirectory,
	}

	// Act
	savedPath, err := handler.SaveTemplate(templateDto)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(outputDirectory, "Valid Save Template.rmg.json"), savedPath)
}

func TestWhenTemplateAndOutputPathAreValid_WritesTemplateFile(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	outputDirectory := t.TempDir()
	template := generateDefaultTemplate(t, handler)
	template.Name = "Written Save Template"
	templateDto := dtos.TemplateSaveDto{
		Template:   template,
		Topology:   config.TopologyRing,
		OutputPath: outputDirectory,
	}

	// Act
	savedPath, err := handler.SaveTemplate(templateDto)

	// Assert
	require.NoError(t, err)
	assert.FileExists(t, savedPath)
}

func TestWhenTemplateAndOutputPathAreValid_WritesPreviewImage(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	outputDirectory := t.TempDir()
	template := generateDefaultTemplate(t, handler)
	template.Name = "Preview Save Template"
	templateDto := dtos.TemplateSaveDto{
		Template:   template,
		Topology:   config.TopologyRing,
		OutputPath: outputDirectory,
	}

	// Act
	_, err := handler.SaveTemplate(templateDto)

	// Assert
	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(outputDirectory, "Preview Save Template.png"))
}

func TestWhenTemplateOutputPathPointsToExistingFile_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	blockingFile := filepath.Join(t.TempDir(), "not-a-directory.txt")
	require.NoError(t, os.WriteFile(blockingFile, []byte("occupied"), 0o644))
	template := generateDefaultTemplate(t, handler)
	templateDto := dtos.TemplateSaveDto{
		Template:   template,
		Topology:   config.TopologyRing,
		OutputPath: blockingFile,
	}

	// Act
	_, err := handler.SaveTemplate(templateDto)

	// Assert
	assert.Error(t, err)
}

func TestWhenPreviewImageCannotBeWritten_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	outputDirectory := t.TempDir()
	template := generateDefaultTemplate(t, handler)
	template.Name = "Blocked Preview Template"
	require.NoError(t, os.Mkdir(filepath.Join(outputDirectory, "Blocked Preview Template.png"), 0o755))
	templateDto := dtos.TemplateSaveDto{
		Template:   template,
		Topology:   config.TopologyRing,
		OutputPath: outputDirectory,
	}

	// Act
	_, err := handler.SaveTemplate(templateDto)

	// Assert
	assert.Error(t, err)
}

func TestWhenPreviewImageCannotBeWritten_StillReturnsTemplateFilePath(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	outputDirectory := t.TempDir()
	template := generateDefaultTemplate(t, handler)
	template.Name = "Blocked Preview Path Template"
	require.NoError(t, os.Mkdir(filepath.Join(outputDirectory, "Blocked Preview Path Template.png"), 0o755))
	templateDto := dtos.TemplateSaveDto{
		Template:   template,
		Topology:   config.TopologyRing,
		OutputPath: outputDirectory,
	}

	// Act
	savedPath, err := handler.SaveTemplate(templateDto)

	// Assert
	require.Error(t, err)
	assert.Equal(t, filepath.Join(outputDirectory, "Blocked Preview Path Template.rmg.json"), savedPath)
}
