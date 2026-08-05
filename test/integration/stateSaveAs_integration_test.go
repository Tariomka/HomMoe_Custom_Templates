//go:build integration_test

package integration_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// newSaveAsProbe drives State.SaveAs all the way through the save-file
// dialog's confirmation callback, which is the only path that assigns
// currentPath. The callback is captured in the dialog's unexported onSave
// field and is otherwise reachable only from a layout.Context, hence the
// integration_test accessors.
func newSaveAsProbe(t *testing.T, saveResult error) (state *drivers.State, chosenPath string) {
	t.Helper()
	handlerMock := &test_helpers.TemplateHandlerMock{}
	handlerMock.On("SaveState", mock.Anything).Return("", saveResult)
	state = drivers.NewUIState(handlerMock, false)
	chosenPath = filepath.Join(t.TempDir(), gofakeit.Word()+".gen.json")

	state.SaveAs(gofakeit.ProductName())
	saveDialog, isFileExplorer := state.GetDialogHost().GetTopDialog().(*dialogs.FileExplorerDialog)
	require.True(t, isFileExplorer, "SaveAs must open the file explorer dialog in save mode")
	saveDialog.ConfirmSave(chosenPath)

	return state, chosenPath
}

// TestWhenSaveAsFails_CurrentPathIsNotRecorded: a failed Save As must not
// retarget the editor at a path that holds no file, otherwise the following
// Save silently writes to the same broken location instead of re-prompting.
func TestWhenSaveAsFails_CurrentPathIsNotRecorded(t *testing.T) {
	// Arrange
	state, _ := newSaveAsProbe(t, errors.New(gofakeit.Sentence(3)))

	// Act
	currentPath := state.GetCurrentPath()

	// Assert
	assert.Empty(t, currentPath)
}

// TestWhenSaveAsSucceeds_CurrentPathIsRecorded: the successful path must still
// record the chosen file so a later Save writes to it without re-prompting.
func TestWhenSaveAsSucceeds_CurrentPathIsRecorded(t *testing.T) {
	// Arrange
	state, chosenPath := newSaveAsProbe(t, nil)

	// Act
	currentPath := state.GetCurrentPath()

	// Assert
	assert.Equal(t, chosenPath, currentPath)
}
