//go:build integration_test && gui

package gui_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// saveSuffix mirrors the dialog's enforced save extension.
const saveSuffix = ".gen.json"

// The template name a fresh editor starts with and the file the Save To dialog
// therefore resolves it to. Every action below is compared against a golden, so
// the names the dialog renders are fixed rather than fuzzed.
const (
	defaultTemplateName = "Custom Template"
	defaultSaveFile     = defaultTemplateName + saveSuffix
)

// The names the tests seed or type. Fixed for the same reason.
const (
	roundTripTemplateName = "Round Trip"
	roundTripSaveFile     = roundTripTemplateName + saveSuffix
	discardedTemplateName = "Discarded"
	existingSaveFile      = "existing.gen.json"
	createdFolderName     = "nested"
)

// newEditor builds a snapshot-verified editor pointed at an empty fixture
// directory, which is where its Load and Save To dialogs open. The runner is
// returned alongside the handler because the editor state these tests assert on
// is read through it.
func newEditor(t *testing.T) (*integration_common.AppRunner, *integration_common.BaseHandler) {
	t.Helper()
	runner := integration_common.NewAppRunner(t)
	if !integration_common.IsHeadless() {
		runner.SetRenderDelay(500 * time.Millisecond)
	}

	return runner, integration_common.NewHandler(runner).WithFixtureDirectory().WithSnapshots()
}

// TestWhenOpenDialogConfirmsAFile_TheEditorStateIsLoaded drives the whole round
// trip through the real toolbar: save the editor to the fixture directory,
// change the state, then load the saved file back and confirm the change was
// undone by the load rather than by the dialog closing.
//
//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWhenOpenDialogConfirmsAFile_TheEditorStateIsLoaded(t *testing.T) {
	// Arrange
	runner, handler := newEditor(t)
	runner.SetTemplateName(roundTripTemplateName)
	handler.ClickSaveTo().ClickSave()
	require.FileExists(t, filepath.Join(handler.FixtureDirectory(), roundTripSaveFile))
	runner.SetTemplateName(discardedTemplateName)

	// Act
	handler.ClickLoad().ClickRow(roundTripSaveFile).ClickOpen()

	// Assert
	assert.Equal(t, roundTripTemplateName, runner.CurrentState().TemplateName)
}

// TestWhenSaveDialogIsConfirmed_TheResolvedNameBecomesTheSaveTarget pins the
// dialog's contract: the name it was opened with plus the directory it is
// showing produce the path the editor then remembers, suffix included. The name
// is never typed because the field is a read-only preview.
//
//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWhenSaveDialogIsConfirmed_TheResolvedNameBecomesTheSaveTarget(t *testing.T) {
	// Arrange
	runner, handler := newEditor(t)

	// Act
	handler.ClickSaveTo().ClickSave()

	// Assert
	assert.Equal(t, filepath.Join(handler.FixtureDirectory(), defaultSaveFile), runner.CurrentPath())
}

// TestWhenTheSaveDialogIsLaidOut_TheNameFieldIsReadOnly guards the point of the
// rename: the row previews the name the state will be written under, so the user
// must not be able to type a name that would then be discarded.
//
//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWhenTheSaveDialogIsLaidOut_TheNameFieldIsReadOnly(t *testing.T) {
	// Arrange
	_, handler := newEditor(t)

	// Act
	explorer := handler.ClickSaveTo()

	// Assert
	assert.True(t, explorer.Dialog().SaveNameReadOnly())
}

// TestWhenNoNameWasResolved_TheConfirmButtonIsDisabled: an unnamed template
// resolves to no filename, and the field can no longer be typed into, so the
// save has to be refused up front instead of writing under a fallback name.
//
//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWhenNoNameWasResolved_TheConfirmButtonIsDisabled(t *testing.T) {
	// Arrange
	runner, handler := newEditor(t)
	runner.SetTemplateName("")

	// Act
	explorer := handler.ClickSaveTo()

	// Assert
	assert.True(t, explorer.Dialog().ConfirmDisabled())
}

//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWhenANameWasResolved_TheConfirmButtonIsEnabled(t *testing.T) {
	// Arrange
	_, handler := newEditor(t)

	// Act
	explorer := handler.ClickSaveTo()

	// Assert
	assert.False(t, explorer.Dialog().ConfirmDisabled())
}

// TestWhenAFileRowIsClickedInSaveMode_TheResolvedNameIsUnchanged: a row click
// used to retarget the save at the clicked file, which would silently write the
// state under a name the template does not have.
//
//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWhenAFileRowIsClickedInSaveMode_TheResolvedNameIsUnchanged(t *testing.T) {
	// Arrange
	_, handler := newEditor(t)
	handler.WithFixtureFiles(existingSaveFile)

	// Act
	explorer := handler.ClickSaveTo().ClickRow(existingSaveFile)

	// Assert
	assert.Equal(t, defaultSaveFile, explorer.Dialog().ResolvedSaveName())
}

// TestWhenSaveDialogIsConfirmed_AFileLandsInTheChosenDirectory asserts about
// bytes on disk rather than the path the editor recorded, so the repository is
// covered as well as the driver.
//
//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWhenSaveDialogIsConfirmed_AFileLandsInTheChosenDirectory(t *testing.T) {
	// Arrange
	_, handler := newEditor(t)

	// Act
	handler.ClickSaveTo().ClickSave()

	// Assert
	written, err := filepath.Glob(filepath.Join(handler.FixtureDirectory(), "*"+saveSuffix))
	require.NoError(t, err)
	assert.Len(t, written, 1)
}

// newOverwriteProbe saves the editor once, changes the state so a second save
// would produce different bytes, then confirms that second save onto the same
// target - which must raise the overwrite prompt instead of writing. The bytes
// the first save produced come back with it, so a gated write is observable.
func newOverwriteProbe(t *testing.T) (
	explorer *integration_common.FileExplorerHandler, target string, original []byte) {
	t.Helper()
	_, handler := newEditor(t)
	handler.ClickSaveTo().ClickSave()

	target = filepath.Join(handler.FixtureDirectory(), defaultSaveFile)
	original, err := os.ReadFile(target)
	require.NoError(t, err)
	require.NotEmpty(t, original)

	explorer = handler.ClickGeneralTab().SelectGameMode(true).ClickSaveTo().ClickSave()
	require.True(t, explorer.Dialog().OverwriteActive(), "saving onto an existing file must raise the prompt")

	return explorer, target, original
}

// TestWhenSaveTargetAlreadyExists_ConfirmDoesNotWrite: raising the prompt must
// not be preceded by the write it is meant to guard.
//
//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWhenSaveTargetAlreadyExists_ConfirmDoesNotWrite(t *testing.T) {
	// Arrange
	_, target, original := newOverwriteProbe(t)

	// Act
	content, err := os.ReadFile(target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, string(original), string(content))
}

//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWhenOverwriteIsCancelled_TheExistingFileIsUntouched(t *testing.T) {
	// Arrange
	explorer, target, original := newOverwriteProbe(t)

	// Act
	explorer.ClickOverwriteCancel()

	// Assert
	content, err := os.ReadFile(target)
	require.NoError(t, err)
	assert.Equal(t, string(original), string(content))
}

//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWhenOverwriteIsConfirmed_TheFileIsRewritten(t *testing.T) {
	// Arrange
	explorer, target, original := newOverwriteProbe(t)

	// Act
	explorer.ClickOverwrite()

	// Assert
	content, err := os.ReadFile(target)
	require.NoError(t, err)
	assert.NotEqual(t, string(original), string(content))
}

// TestWhenANewFolderIsCreated_ItAppearsInTheParentListing: creation navigates
// into the new folder, so the listing is checked from the parent it was made in.
//
//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWhenANewFolderIsCreated_ItAppearsInTheParentListing(t *testing.T) {
	// Arrange
	_, handler := newEditor(t)
	explorer := handler.ClickSaveTo().ClickNewFolder().TypeFolderName(createdFolderName).ClickCreateFolder()
	require.Empty(t, explorer.Dialog().NewFolderError())

	// Act
	explorer.ClickBack()

	// Assert
	assert.Contains(t, explorer.Dialog().EntryNames(), createdFolderName)
}

// TestWhenNewFolderIsClickedTwice_TheRowIsDismissed: the button is a toggle, so
// a user who changes their mind has to be able to put the row away again.
//
//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWhenNewFolderIsClickedTwice_TheRowIsDismissed(t *testing.T) {
	// Arrange
	_, handler := newEditor(t)
	explorer := handler.ClickSaveTo().ClickNewFolder()
	require.True(t, explorer.Dialog().NewFolderActive())

	// Act
	explorer.ClickNewFolder()

	// Assert
	assert.False(t, explorer.Dialog().NewFolderActive())
}

// TestWhenSaveTargetIsAnExistingFolder_TheSaveIsRefused covers the defect where
// a directory sharing the save target's name was offered as an overwrite; the
// New Folder button makes "name.gen.json" a reachable folder name.
//
//nolint:paralleltest // Snapshots need exclusive access to the single headless GPU window.
func TestWhenSaveTargetIsAnExistingFolder_TheSaveIsRefused(t *testing.T) {
	// Arrange
	_, handler := newEditor(t)
	handler.WithFixtureFolders(defaultSaveFile)

	// Act
	explorer := handler.ClickSaveTo().ClickSave()

	// Assert
	assert.Equal(t, "A folder with that name already exists.", explorer.Dialog().SaveError())
}

// The name the recovery tests export under, and the two files that export
// produces.
const (
	recoveryTemplateName = "Recovery"
	recoveryTemplateFile = recoveryTemplateName + ".rmg.json"
	recoveryPreviewFile  = recoveryTemplateName + ".png"
)

// The recovery status the editor shows when it could not find the game's own
// templates directory. Exporting anywhere else would produce a file the game
// never reads (AGENTS.md 2.7), so the user is sent to the picker instead.
const recoveryStatus = "Game template directory not found. " +
	"Choose the game templates folder using the output folder picker before exporting."

// newRecoveryEditor builds an editor whose game directory lookup fails, so the
// recovery path is exercised even on a machine with the game installed - and so
// nothing here can browse or write into that real directory. The fixture
// directory it returns is where the picker opens and where an export lands once
// the user has confirmed it. Snapshots stay off: every frame shows that per-run
// path, which no golden can hold.
func newRecoveryEditor(t *testing.T) (*integration_common.AppRunner, *integration_common.BaseHandler, string) {
	t.Helper()
	fixtureDirectory := t.TempDir()
	runner := integration_common.NewAppRunnerWithFileSystem(
		t, test_helpers.NewFailedLookupFileSystemHandler(fixtureDirectory))
	if !integration_common.IsHeadless() {
		runner.SetRenderDelay(500 * time.Millisecond)
	}

	return runner, integration_common.NewHandler(runner), fixtureDirectory
}

//nolint:paralleltest // The GUI suite drives one shared window.
func TestWhenTheGameTemplateDirectoryIsNotFound_NoExportDestinationIsSelected(t *testing.T) {
	// Arrange & Act
	runner, _, _ := newRecoveryEditor(t)

	// Assert
	assert.Empty(t, runner.OutputPath())
}

//nolint:paralleltest // The GUI suite drives one shared window.
func TestWhenTheGameTemplateDirectoryIsNotFound_TheStatusDirectsToThePicker(t *testing.T) {
	// Arrange & Act
	runner, _, _ := newRecoveryEditor(t)

	// Assert
	message, _ := runner.Status()
	assert.Equal(t, recoveryStatus, message)
}

// The picker still has to open somewhere usable, but where it browses is not a
// destination: only confirming a folder is.
//
//nolint:paralleltest // The GUI suite drives one shared window.
func TestWhenTheOutputPickerOpensWithoutADetectedDirectory_ItStartsAtABrowsableDirectory(t *testing.T) {
	// Arrange
	_, handler, fixtureDirectory := newRecoveryEditor(t)

	// Act
	explorer := handler.ClickBrowseOutput()

	// Assert
	assert.Equal(t, fixtureDirectory, explorer.Dialog().CurrentDir())
}

//nolint:paralleltest // The GUI suite drives one shared window.
func TestWhenTheOutputPickerIsCancelled_NoExportDestinationIsSelected(t *testing.T) {
	// Arrange
	runner, handler, _ := newRecoveryEditor(t)

	// Act
	handler.ClickBrowseOutput().ClickCancel()

	// Assert
	assert.Empty(t, runner.OutputPath())
}

//nolint:paralleltest // The GUI suite drives one shared window.
func TestWhenTheOutputFolderIsConfirmed_ItBecomesTheExportDestination(t *testing.T) {
	// Arrange
	runner, handler, fixtureDirectory := newRecoveryEditor(t)

	// Act
	handler.ClickBrowseOutput().ClickSelectFolder()

	// Assert
	assert.Equal(t, fixtureDirectory, runner.OutputPath())
}

// Without a confirmed destination there is nowhere the game would read the
// result from, so neither the template nor its preview may be written.
//
//nolint:paralleltest // The GUI suite drives one shared window.
func TestWhenNoExportDestinationIsSelected_SavingTheTemplateWritesNothing(t *testing.T) {
	// Arrange
	runner, handler, fixtureDirectory := newRecoveryEditor(t)
	runner.SetTemplateName(recoveryTemplateName)
	handler.ClickGenerate()

	// Act
	handler.ClickSaveTemplate()

	// Assert
	written, err := filepath.Glob(filepath.Join(fixtureDirectory, "*"))
	require.NoError(t, err)
	assert.Empty(t, written)
}

//nolint:paralleltest // The GUI suite drives one shared window.
func TestWhenTheExportDestinationWasConfirmed_TheTemplateAndItsPreviewLandThere(t *testing.T) {
	// Arrange
	runner, handler, fixtureDirectory := newRecoveryEditor(t)
	runner.SetTemplateName(recoveryTemplateName)
	handler.ClickGenerate().ClickBrowseOutput().ClickSelectFolder()

	// Act
	handler.ClickSaveTemplate()

	// Assert
	written, err := filepath.Glob(filepath.Join(fixtureDirectory, "*"))
	require.NoError(t, err)
	assert.ElementsMatch(t,
		[]string{
			filepath.Join(fixtureDirectory, recoveryTemplateFile),
			filepath.Join(fixtureDirectory, recoveryPreviewFile),
		},
		written)
}

// The destination is a property of the machine, not of the document: a session
// that picked one must not hand it to the next one through a saved file.
//
//nolint:paralleltest // The GUI suite drives one shared window.
func TestWhenAStateSavedWithAPickedDestinationIsReloaded_TheDestinationIsNotRestored(t *testing.T) {
	// Arrange
	runner, handler, fixtureDirectory := newRecoveryEditor(t)
	runner.SetTemplateName(recoveryTemplateName)
	handler.ClickBrowseOutput().ClickSelectFolder()
	runner.SaveStateToFile(filepath.Join(fixtureDirectory, recoveryTemplateName+saveSuffix))
	savedPath := runner.CurrentPath()
	require.FileExists(t, savedPath)

	// Act
	nextRunner, _, _ := newRecoveryEditor(t)
	nextRunner.LoadStateFromFile(savedPath)

	// Assert
	assert.Empty(t, nextRunner.OutputPath())
}
