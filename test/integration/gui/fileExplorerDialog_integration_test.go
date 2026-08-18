//go:build integration_test && gui

package gui_test

import (
	"image"
	"os"
	"path/filepath"
	"testing"

	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// saveSuffix mirrors the dialog's enforced save extension.
const saveSuffix = ".gen.json"

// frameFileExplorer lays out one dialog frame the way DialogHost does on every
// vsync - button clicks queued with widget.Clickable.Click are consumed here -
// and reports whether the dialog asked to close.
func frameFileExplorer(t *testing.T, dialog *dialogs.FileExplorerDialog, theme *material.Theme) bool {
	t.Helper()
	gtx, frameRouter := newDialogContext(image.Pt(720, 560))
	_, done := dialog.Body(gtx, theme)
	frameRouter.Frame(gtx.Ops)
	return done
}

// TestWhenOpenDialogConfirmsAFile_TheEditorStateIsLoaded drives the whole open
// flow: list a directory, click a row, confirm, and let the pick handler install
// the state - the path the review found untested because it needs a
// layout.Context.
func TestWhenOpenDialogConfirmsAFile_TheEditorStateIsLoaded(t *testing.T) {
	t.Parallel()
	// Arrange
	directory := t.TempDir()
	theme := themes.NewTheme()
	fileSystem := composition.InitializeFileSystemHandler()
	templateName := gofakeit.LetterN(8)

	source := drivers.NewUIState(
		composition.InitializeGuiHandler(), fileSystem, composition.InitializeRegenerationHandler(), false)
	source.UpdateState(func(state *dtos.EditorStateDto) { state.TemplateName = templateName })
	source.SaveStateToFile(filepath.Join(directory, gofakeit.LetterN(6)+saveSuffix))
	fixturePath := source.GetCurrentPath()
	require.FileExists(t, fixturePath)

	target := drivers.NewUIState(
		composition.InitializeGuiHandler(), fileSystem, composition.InitializeRegenerationHandler(), false)
	dialog := dialogs.NewOpenFileDialog(fileSystem, directory, []string{saveSuffix}, target.LoadStateFromFile)
	require.True(t, dialog.ClickEntry(filepath.Base(fixturePath)), "the saved fixture must be listed")
	require.False(t, frameFileExplorer(t, dialog, theme))
	require.Equal(t, fixturePath, dialog.SelectedPath())

	// Act
	dialog.ClickConfirm()
	frameFileExplorer(t, dialog, theme)

	// Assert
	assert.Equal(t, templateName, target.GetStateData().TemplateName)
}

// TestWhenSaveDialogIsConfirmed_TheResolvedNameBecomesTheSaveTarget pins the
// dialog's own contract: the name it was opened with plus the current directory
// produce the path handed to the save callback, suffix included. The name comes
// from the caller because the field is a read-only preview.
func TestWhenSaveDialogIsConfirmed_TheResolvedNameBecomesTheSaveTarget(t *testing.T) {
	t.Parallel()
	// Arrange
	directory := t.TempDir()
	theme := themes.NewTheme()
	filename := gofakeit.LetterN(8)
	var savedPath string
	dialog := dialogs.NewSaveFileDialog(
		composition.InitializeFileSystemHandler(),
		directory,
		filename,
		func(path string) { savedPath = path })

	// Act
	dialog.ClickConfirm()
	frameFileExplorer(t, dialog, theme)

	// Assert
	assert.Equal(t, filepath.Join(directory, filename+saveSuffix), savedPath)
}

// TestWhenTheSaveDialogIsLaidOut_TheNameFieldIsReadOnly guards the point of the
// rename: the row previews the name the state will be written under, so the
// user must not be able to type a name that would then be discarded.
func TestWhenTheSaveDialogIsLaidOut_TheNameFieldIsReadOnly(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := dialogs.NewSaveFileDialog(
		composition.InitializeFileSystemHandler(),
		t.TempDir(),
		gofakeit.LetterN(8),
		nil)

	// Act
	frameFileExplorer(t, dialog, themes.NewTheme())

	// Assert
	assert.True(t, dialog.SaveNameReadOnly())
}

// TestWhenNoNameWasResolved_TheConfirmButtonIsDisabled: an unnamed template
// resolves to no filename, and the field can no longer be typed into, so the
// save has to be refused up front instead of writing under a fallback name.
func TestWhenNoNameWasResolved_TheConfirmButtonIsDisabled(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := dialogs.NewSaveFileDialog(
		composition.InitializeFileSystemHandler(),
		t.TempDir(),
		"",
		nil)

	// Act
	frameFileExplorer(t, dialog, themes.NewTheme())

	// Assert
	assert.True(t, dialog.ConfirmDisabled())
}

// TestWhenAFileRowIsClickedInSaveMode_TheResolvedNameIsUnchanged: a row click
// used to retarget the save at the clicked file, which would silently write the
// state under a name the template does not have.
func TestWhenAFileRowIsClickedInSaveMode_TheResolvedNameIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	directory := t.TempDir()
	resolvedName := gofakeit.LetterN(8) + saveSuffix
	existing := gofakeit.LetterN(9) + saveSuffix
	require.NoError(t, os.WriteFile(filepath.Join(directory, existing), []byte("{}"), 0o600))
	dialog := dialogs.NewSaveFileDialog(
		composition.InitializeFileSystemHandler(), directory, resolvedName, nil)
	require.True(t, dialog.ClickEntry(existing), "the fixture must be listed")

	// Act
	frameFileExplorer(t, dialog, themes.NewTheme())

	// Assert
	assert.Equal(t, resolvedName, dialog.ResolvedSaveName())
}

// TestWhenSaveDialogIsConfirmedThroughTheDriver_AFileLandsInTheChosenDirectory
// runs the same confirm through the real state driver and repository, so the
// assertion is about bytes on disk rather than a callback argument.
func TestWhenSaveDialogIsConfirmedThroughTheDriver_AFileLandsInTheChosenDirectory(t *testing.T) {
	t.Parallel()
	// Arrange
	directory := t.TempDir()
	theme := themes.NewTheme()
	fileSystem := composition.InitializeFileSystemHandler()
	state := drivers.NewUIState(
		composition.InitializeGuiHandler(), fileSystem, composition.InitializeRegenerationHandler(), false)
	dialog := dialogs.NewSaveFileDialog(fileSystem, directory, gofakeit.LetterN(8), state.SaveStateToFile)

	// Act
	dialog.ClickConfirm()
	require.True(t, frameFileExplorer(t, dialog, theme), "a successful save must close the dialog")

	// Assert
	written, err := filepath.Glob(filepath.Join(directory, "*"+saveSuffix))
	require.NoError(t, err)
	assert.Len(t, written, 1)
}

// newOverwriteProbe opens a save dialog whose target already exists and drives
// the confirm click that must raise the overwrite prompt instead of writing. The
// save callback rewrites the file, standing in for the state driver, so a gated
// write is observable on disk.
func newOverwriteProbe(
	t *testing.T,
	theme *material.Theme) (dialog *dialogs.FileExplorerDialog, target string) {
	t.Helper()
	directory := t.TempDir()
	filename := gofakeit.LetterN(8)
	target = filepath.Join(directory, filename+saveSuffix)
	require.NoError(t, os.WriteFile(target, []byte("original"), 0o600))

	dialog = dialogs.NewSaveFileDialog(
		composition.InitializeFileSystemHandler(),
		directory,
		filename,
		func(path string) {
			require.NoError(t, os.WriteFile(path, []byte("rewritten"), 0o600))
		})

	dialog.ClickConfirm()
	require.False(t, frameFileExplorer(t, dialog, theme), "the prompt must keep the dialog open")
	require.True(t, dialog.OverwriteActive(), "confirming an existing target must raise the prompt")

	return dialog, target
}

// TestWhenSaveTargetAlreadyExists_ConfirmDoesNotWrite: raising the prompt must
// not be preceded by the write it is meant to guard.
func TestWhenSaveTargetAlreadyExists_ConfirmDoesNotWrite(t *testing.T) {
	t.Parallel()
	// Arrange
	_, target := newOverwriteProbe(t, themes.NewTheme())

	// Act
	content, err := os.ReadFile(target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "original", string(content))
}

func TestWhenOverwriteIsCancelled_TheExistingFileIsUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	dialog, target := newOverwriteProbe(t, theme)

	// Act
	dialog.ClickOverwriteCancel()
	frameFileExplorer(t, dialog, theme)

	// Assert
	content, err := os.ReadFile(target)
	require.NoError(t, err)
	assert.Equal(t, "original", string(content))
}

func TestWhenOverwriteIsConfirmed_TheFileIsRewritten(t *testing.T) {
	t.Parallel()
	// Arrange
	theme := themes.NewTheme()
	dialog, target := newOverwriteProbe(t, theme)

	// Act
	dialog.ClickOverwriteConfirm()
	frameFileExplorer(t, dialog, theme)

	// Assert
	content, err := os.ReadFile(target)
	require.NoError(t, err)
	assert.Equal(t, "rewritten", string(content))
}

// TestWhenANewFolderIsCreated_ItAppearsInTheParentListing: creation navigates
// into the new folder, so the listing is checked from the parent it was made in.
func TestWhenANewFolderIsCreated_ItAppearsInTheParentListing(t *testing.T) {
	t.Parallel()
	// Arrange
	directory := t.TempDir()
	theme := themes.NewTheme()
	folderName := gofakeit.LetterN(8)
	dialog := dialogs.NewSaveFileDialog(
		composition.InitializeFileSystemHandler(),
		directory,
		gofakeit.LetterN(6),
		nil)
	dialog.ClickNewFolder()
	require.False(t, frameFileExplorer(t, dialog, theme))

	// Act
	dialog.SetNewFolderName(folderName)
	dialog.ClickCreateFolder()
	require.False(t, frameFileExplorer(t, dialog, theme))
	require.Empty(t, dialog.NewFolderError())
	dialog.ClickUp()
	frameFileExplorer(t, dialog, theme)

	// Assert
	assert.Contains(t, dialog.EntryNames(), folderName)
}

// TestWhenSaveTargetIsAnExistingFolder_TheSaveIsRefused covers the defect where
// a directory sharing the save target's name was offered as an overwrite; the
// New Folder button makes "name.gen.json" a reachable folder name.
func TestWhenSaveTargetIsAnExistingFolder_TheSaveIsRefused(t *testing.T) {
	t.Parallel()
	// Arrange
	directory := t.TempDir()
	theme := themes.NewTheme()
	filename := gofakeit.LetterN(8)
	require.NoError(t, os.Mkdir(filepath.Join(directory, filename+saveSuffix), 0o750))
	dialog := dialogs.NewSaveFileDialog(
		composition.InitializeFileSystemHandler(),
		directory,
		filename,
		func(string) { require.Fail(t, "a folder must never be overwritten") })

	// Act
	dialog.ClickConfirm()
	frameFileExplorer(t, dialog, theme)

	// Assert
	assert.Equal(t, "A folder with that name already exists.", dialog.SaveError())
}

func TestWhenANameWasResolved_TheConfirmButtonIsEnabled(t *testing.T) {
	t.Parallel()
	// Arrange
	dialog := dialogs.NewSaveFileDialog(
		composition.InitializeFileSystemHandler(),
		t.TempDir(),
		gofakeit.LetterN(8),
		nil)

	// Act
	frameFileExplorer(t, dialog, themes.NewTheme())

	// Assert
	assert.False(t, dialog.ConfirmDisabled())
}
