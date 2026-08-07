package pathResolutionService_test

import (
	"os"
	"path/filepath"
	"runtime"
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

func TestWhenSaveNameIsAWindowsDeviceOnWindows_ReportsNoTarget(t *testing.T) {
	t.Parallel()
	if runtime.GOOS != "windows" {
		t.Skip("DOS device names are only reserved on Windows")
	}
	names := map[string]string{
		"BareDeviceName":               "NUL",
		"LowercaseDeviceName":          "con",
		"NumberedSerialPort":           "COM1",
		"NumberedPrinterPort":          "lpt9",
		"DeviceNameCarryingSuffix":     "AUX" + saveFileSuffix,
		"DeviceNameWithOtherExtension": "PRN.txt",
	}
	for scenario, name := range names {
		t.Run(scenario+"_ReportsNoTarget", func(t *testing.T) {
			t.Parallel()
			// Arrange
			service := file_system.NewPathResolutionService()

			// Act
			_, ok := service.ResolveSaveTarget(t.TempDir(), name, saveFileSuffix)

			// Assert
			assert.False(t, ok)
		})
	}
}

func TestWhenSaveNameOnlyStartsLikeADevice_IsAccepted(t *testing.T) {
	t.Parallel()
	// Arrange
	service := file_system.NewPathResolutionService()
	directory := t.TempDir()

	// Act
	target, ok := service.ResolveSaveTarget(directory, "console", saveFileSuffix)

	// Assert
	require.True(t, ok)
	assert.Equal(t, filepath.Join(directory, "console"+saveFileSuffix), target)
}

func TestWhenSaveNameIsAWindowsDeviceOutsideWindows_IsAccepted(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "windows" {
		t.Skip("DOS device names are reserved on Windows")
	}
	// Arrange
	service := file_system.NewPathResolutionService()
	directory := t.TempDir()

	// Act
	target, ok := service.ResolveSaveTarget(directory, "NUL", saveFileSuffix)

	// Assert
	require.True(t, ok)
	assert.Equal(t, filepath.Join(directory, "NUL"+saveFileSuffix), target)
}
