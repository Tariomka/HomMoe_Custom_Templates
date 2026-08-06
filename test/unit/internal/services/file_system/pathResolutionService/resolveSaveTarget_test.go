package pathResolutionService_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_system"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const saveFileSuffix = ".gen.json"

func TestWhenSaveNameIsBlank_ReportsNoTarget(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()

	// Act
	_, ok := service.ResolveSaveTarget(t.TempDir(), "   ", saveFileSuffix)

	// Assert
	assert.False(t, ok)
}

func TestWhenSaveNameIsARelativePathToken_ReportsNoTarget(t *testing.T) {
	t.Parallel()
	tokens := map[string]string{
		"CurrentDirectoryToken": ".",
		"ParentDirectoryToken":  "..",
		"PathSeparatorOnly":     string(os.PathSeparator),
	}
	for scenario, token := range tokens {
		t.Run(scenario+"_ReportsNoTarget", func(t *testing.T) {
			t.Parallel()
			// Arrange
			service := file_system.NewPathResolutionService()

			// Act
			_, ok := service.ResolveSaveTarget(t.TempDir(), token, saveFileSuffix)

			// Assert
			assert.False(t, ok)
		})
	}
}

func TestWhenSaveNameLacksTheRequiredSuffix_AppendsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()
	directory := t.TempDir()

	// Act
	target, ok := service.ResolveSaveTarget(directory, "My Template", saveFileSuffix)

	// Assert
	require.True(t, ok)
	assert.Equal(t, filepath.Join(directory, "My Template"+saveFileSuffix), target)
}

func TestWhenSaveNameAlreadyCarriesTheSuffix_LeavesItAlone(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()
	directory := t.TempDir()

	// Act
	target, ok := service.ResolveSaveTarget(directory, "My Template"+saveFileSuffix, saveFileSuffix)

	// Assert
	require.True(t, ok)
	assert.Equal(t, filepath.Join(directory, "My Template"+saveFileSuffix), target)
}

func TestWhenSaveNameCarriesTheSuffixInAnotherCase_LeavesItAlone(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()
	directory := t.TempDir()

	// Act
	target, ok := service.ResolveSaveTarget(directory, "My Template.GEN.JSON", saveFileSuffix)

	// Assert
	require.True(t, ok)
	assert.Equal(t, filepath.Join(directory, "My Template.GEN.JSON"), target)
}

func TestWhenSaveNameContainsDirectoryComponents_KeepsTheTargetInsideTheDirectory(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()
	directory := t.TempDir()

	// Act
	target, ok := service.ResolveSaveTarget(directory, "../../escape"+saveFileSuffix, saveFileSuffix)

	// Assert
	require.True(t, ok)
	assert.Equal(t, filepath.Join(directory, "escape"+saveFileSuffix), target)
}

func TestWhenNoSuffixIsRequired_UsesTheNameVerbatim(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()
	directory := t.TempDir()

	// Act
	target, ok := service.ResolveSaveTarget(directory, "plain.txt", "")

	// Assert
	require.True(t, ok)
	assert.Equal(t, filepath.Join(directory, "plain.txt"), target)
}
