//go:build integration_test

package integration_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
)

// TestHandlerDialogs_LoadOpensFileExplorer proves the toolbar's Load button is
// wired to the file dialog. Nothing else tests that the button reaches it.
func TestHandlerDialogs_LoadOpensFileExplorer(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner)

	// Act
	dialog := handler.ClickLoad()

	// Assert
	assert.True(t, dialog.IsOpen())
}

func TestHandlerDialogs_ClosingFileExplorerDismissesIt(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	dialog := integration_common.NewHandler(runner).ClickLoad()

	// Act
	dialog.Close()

	// Assert
	assert.False(t, dialog.IsOpen())
}

func TestHandlerDialogs_SaveAsOpensFileExplorer(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner)

	// Act
	dialog := handler.ClickSaveAs()

	// Assert
	assert.True(t, dialog.IsOpen())
}

func TestHandlerDialogs_LayoutTabOpensZoneEditor(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner)

	// Act
	dialog := handler.ClickLayoutAndZonesTab().OpenZoneEditor()

	// Assert
	assert.True(t, dialog.IsOpen())
}

func TestHandlerDialogs_ClosingZoneEditorDismissesIt(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	dialog := integration_common.NewHandler(runner).ClickLayoutAndZonesTab().OpenZoneEditor()

	// Act
	dialog.Close()

	// Assert
	assert.False(t, dialog.IsOpen())
}

// TestHandlerToolbar_NewResetsTheEditor proves ClickNew reaches the toolbar's
// New action rather than a neighbouring button.
func TestHandlerToolbar_NewResetsTheEditor(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner)
	runner.NextFrame()

	// Act
	handler.ClickNew()

	// Assert
	status, _ := runner.Status()
	assert.Equal(t, "New settings file.", status)
}
