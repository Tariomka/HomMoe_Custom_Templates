//go:build integration_test

package integration_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// newSaveToProbe drives State.SaveTo all the way through the save-file
// dialog's confirmation callback, which is the only path that assigns
// currentPath. The callback is captured in the dialog's unexported onSave
// field and is otherwise reachable only from a layout.Context, hence the
// integration_test accessors. writtenPath is the path the handler reports
// back, which is what the editor has to remember.
func newSaveToProbe(t *testing.T, saveResult error) (state *drivers.State, writtenPath string) {
	t.Helper()
	outputDirectory := t.TempDir()
	writtenPath = filepath.Join(outputDirectory, gofakeit.Word()+".gen.json")
	if saveResult != nil {
		writtenPath = ""
	}

	handlerMock := &test_helpers.TemplateHandlerMock{}
	handlerMock.On("SaveState", mock.Anything).Return(writtenPath, saveResult)
	state = drivers.NewUIState(
		handlerMock,
		composition.InitializeFileSystemHandler(),
		composition.InitializeRegenerationHandler(),
		mappers.NewEditorStateMapper(),
		false)

	state.SaveTo(gofakeit.ProductName())
	saveDialog, isFileExplorer := state.GetDialogHost().GetTopDialog().(*dialogs.FileExplorerDialog)
	require.True(t, isFileExplorer, "SaveTo must open the file explorer dialog in save mode")
	saveDialog.ConfirmSave(filepath.Join(outputDirectory, gofakeit.Word()+".gen.json"))

	return state, writtenPath
}

// TestWhenSaveToFails_CurrentPathIsNotRecorded: a failed Save To must not
// retarget the editor at a path that holds no file, otherwise the following
// Save silently writes to the same broken location instead of re-prompting.
func TestWhenSaveToFails_CurrentPathIsNotRecorded(t *testing.T) {
	// Arrange
	state, _ := newSaveToProbe(t, errors.New(gofakeit.Sentence(3)))

	// Act
	currentPath := state.GetCurrentPath()

	// Assert
	assert.Empty(t, currentPath)
}

// TestWhenSaveToSucceeds_CurrentPathIsRecorded: the successful path must still
// record the written file so a later Save writes to it without re-prompting.
func TestWhenSaveToSucceeds_CurrentPathIsRecorded(t *testing.T) {
	// Arrange
	state, writtenPath := newSaveToProbe(t, nil)

	// Act
	currentPath := state.GetCurrentPath()

	// Assert
	assert.Equal(t, writtenPath, currentPath)
}

// newSaveToDialog opens the save dialog for templateName and returns it, so the
// name the dialog previews can be compared with the name SaveSettings will
// actually write the state under.
func newSaveToDialog(t *testing.T, templateName string) *dialogs.FileExplorerDialog {
	t.Helper()
	state := drivers.NewUIState(
		&test_helpers.TemplateHandlerMock{},
		composition.InitializeFileSystemHandler(),
		composition.InitializeRegenerationHandler(),
		mappers.NewEditorStateMapper(),
		false)

	state.SaveTo(templateName)
	saveDialog, isFileExplorer := state.GetDialogHost().GetTopDialog().(*dialogs.FileExplorerDialog)
	require.True(t, isFileExplorer, "SaveTo must open the file explorer dialog in save mode")

	return saveDialog
}

// TestWhenSaveToOpens_TheDialogPreviewsTheResolvedFilename: the dialog picks a
// directory only, so the name it shows has to be the one the repository derives
// from the template name - sanitized and suffixed.
func TestWhenSaveToOpens_TheDialogPreviewsTheResolvedFilename(t *testing.T) {
	// Arrange
	saveDialog := newSaveToDialog(t, "  Jebus: Cross  ")

	// Act
	resolvedName := saveDialog.ResolvedSaveName()

	// Assert
	assert.Equal(t, "Jebus_ Cross.gen.json", resolvedName)
}

// TestWhenTheTemplateIsUnnamed_SaveToPreviewsNoFilename: an unnamed template
// has no name to derive from, and offering a bare ".gen.json" would be a
// promise the repository breaks by falling back to its own default name.
func TestWhenTheTemplateIsUnnamed_SaveToPreviewsNoFilename(t *testing.T) {
	// Arrange
	saveDialog := newSaveToDialog(t, "   ")

	// Act
	resolvedName := saveDialog.ResolvedSaveName()

	// Assert
	assert.Empty(t, resolvedName)
}
