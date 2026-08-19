//go:build integration_test

package integration_common

import (
	"image"

	"gioui.org/f32"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/editor"
)

// The labels the file explorer's buttons are addressed by. They are looked up
// inside the dialog panel, so a label the editor behind the scrim also uses
// (Save, Cancel) still resolves to the dialog's own button.
const (
	backButtonLabel            = "← Back"
	showHiddenButtonLabel      = "Show hidden"
	newFolderButtonLabel       = "New Folder"
	createFolderButtonLabel    = "Create Folder"
	openButtonLabel            = "Open"
	saveButtonLabel            = "Save"
	overwriteButtonLabel       = "Overwrite"
	overwriteCancelButtonLabel = "Cancel"
)

// newFolderFieldOffset is how far left of the Create Folder button the new
// folder textbox is clicked to focus it. The two are adjacent in one flex row,
// so a point just outside the button's left edge lands in the field.
const newFolderFieldOffset = 40

// FileExplorerHandler drives the load/save file dialog. It holds the originating
// BaseHandler but deliberately does not embed it: an open dialog's scrim absorbs
// every background click, so promoting the tab clicks would be a lie.
//
// It snapshots every action, but only once BaseHandler.WithFixtureDirectory has
// pointed the dialog at a fixture directory - otherwise the listing and the path
// bar show the machine's own templates directory (AGENTS.md 2.7), which no
// golden can hold.
type FileExplorerHandler struct {
	base *BaseHandler

	// pathBarMask is the mask over the dialog's path bar, tracked so it can be
	// lifted the moment the dialog closes - it covers a slice of the editor
	// underneath, which every later snapshot should still be comparing.
	pathBarMask image.Rectangle
}

// newFileExplorerHandler wraps the dialog the caller just opened, masks the path
// bar that shows a per-run directory, and snapshots the opened state.
func newFileExplorerHandler(base *BaseHandler) *FileExplorerHandler {
	base.runner.tb.Helper()
	handler := &FileExplorerHandler{base: base}
	if base.fixtureDirectory != "" {
		handler.maskPathBar()
	}
	handler.verifySnapshot()
	return handler
}

func (this *FileExplorerHandler) IsOpen() bool {
	this.base.runner.tb.Helper()
	return this.base.runner.DialogsOpen()
}

// Dialog exposes what the open dialog currently shows: its directory, listing,
// selection and inline errors.
func (this *FileExplorerHandler) Dialog() editor.IFileExplorerDialog {
	this.base.runner.tb.Helper()
	dialog, ok := this.base.runner.TopFileExplorer()
	if !ok {
		this.base.runner.tb.Fatal("the top-most dialog is not a file explorer")
	}

	return dialog
}

// FixtureDirectory returns the directory the dialog was pointed at.
func (this *FileExplorerHandler) FixtureDirectory() string {
	this.base.runner.tb.Helper()
	return this.base.FixtureDirectory()
}

// ClickShowHidden toggles whether dot-prefixed entries are listed.
func (this *FileExplorerHandler) ClickShowHidden() *FileExplorerHandler {
	return this.clickDialogButton(showHiddenButtonLabel)
}

// ClickRow taps the listed entry called name: a file becomes the selection in
// open mode, a folder is descended into.
func (this *FileExplorerHandler) ClickRow(name string) *FileExplorerHandler {
	return this.clickDialogButton(name)
}

// ClickBack ascends to the parent directory.
func (this *FileExplorerHandler) ClickBack() *FileExplorerHandler {
	return this.clickDialogButton(backButtonLabel)
}

// ClickNewFolder toggles the new folder row; clicking it again dismisses it.
func (this *FileExplorerHandler) ClickNewFolder() *FileExplorerHandler {
	return this.clickDialogButton(newFolderButtonLabel)
}

// ClickCreateFolder creates the folder named in the new folder field and
// navigates into it.
func (this *FileExplorerHandler) ClickCreateFolder() *FileExplorerHandler {
	return this.clickDialogButton(createFolderButtonLabel)
}

// ClickOpen confirms the open dialog's selection.
func (this *FileExplorerHandler) ClickOpen() *FileExplorerHandler {
	return this.clickDialogButton(openButtonLabel)
}

// ClickSave confirms the save dialog's resolved target.
func (this *FileExplorerHandler) ClickSave() *FileExplorerHandler {
	return this.clickDialogButton(saveButtonLabel)
}

// ClickOverwrite accepts the overwrite prompt.
func (this *FileExplorerHandler) ClickOverwrite() *FileExplorerHandler {
	return this.clickDialogButton(overwriteButtonLabel)
}

// ClickOverwriteCancel dismisses the overwrite prompt, leaving the file alone.
func (this *FileExplorerHandler) ClickOverwriteCancel() *FileExplorerHandler {
	return this.clickDialogButton(overwriteCancelButtonLabel)
}

// TypeFolderName focuses the new folder field and types text into it.
func (this *FileExplorerHandler) TypeFolderName(text string) *FileExplorerHandler {
	this.base.runner.tb.Helper()
	createButton := this.base.runner.ButtonBoundsIn(fileDialogRect(), createFolderButtonLabel)
	field := f32.Pt(
		float32(createButton.Min.X-newFolderFieldOffset),
		float32((createButton.Min.Y+createButton.Max.Y)/2))

	this.base.runner.ClickAt(field)
	this.base.runner.InputText(text)
	this.verifySnapshot()
	return this
}

// Scroll turns the mouse wheel over the listing by delta pixels; positive
// scrolls the content up. Gio clamps to the listing's range.
func (this *FileExplorerHandler) Scroll(delta float32) *FileExplorerHandler {
	this.base.runner.tb.Helper()
	this.base.runner.Scroll(f32.Pt(fileDialogListScrollX, fileDialogListScrollY), f32.Pt(0, delta))
	this.verifySnapshot()
	return this
}

// Close dismisses the dialog and hands control back to the editor.
func (this *FileExplorerHandler) Close() *BaseHandler {
	this.base.runner.tb.Helper()
	this.base.runner.CloseTopDialog()
	this.unmaskPathBar()
	return this.base
}

// Editor hands control back to the editor without touching the dialog, for the
// confirmations that already closed it themselves.
func (this *FileExplorerHandler) Editor() *BaseHandler {
	this.base.runner.tb.Helper()
	return this.base
}

func (this *FileExplorerHandler) clickDialogButton(label string) *FileExplorerHandler {
	this.base.runner.tb.Helper()
	this.base.runner.ClickButtonIn(fileDialogRect(), label)
	// A row click only records where to navigate; the dialog applies it at the
	// top of the next Body so the listing is not swapped mid-iteration.
	this.base.runner.NextFrame()
	this.verifySnapshot()
	return this
}

// maskPathBar blanks the path bar, which shows the per-run fixture directory.
// The bar fills the gap between the two header buttons, so its rectangle is read
// from their bounds rather than pinned to a coordinate.
func (this *FileExplorerHandler) maskPathBar() {
	this.base.runner.tb.Helper()
	back := this.base.runner.ButtonBoundsIn(fileDialogRect(), backButtonLabel)
	showHidden := this.base.runner.ButtonBoundsIn(fileDialogRect(), showHiddenButtonLabel)
	this.pathBarMask = image.Rect(
		back.Max.X, back.Min.Y-headerBarSlack,
		showHidden.Min.X, back.Max.Y+headerBarSlack)
	this.base.runner.MaskRect(this.pathBarMask)
}

// unmaskPathBar lifts the path bar mask once, if one was ever registered.
func (this *FileExplorerHandler) unmaskPathBar() {
	this.base.runner.tb.Helper()
	if this.pathBarMask.Empty() {
		return
	}

	this.base.runner.UnmaskRect(this.pathBarMask)
	this.pathBarMask = image.Rectangle{}
}

// verifySnapshot compares the golden for the action that just ran, refusing to
// snapshot a dialog that no fixture directory made deterministic. A confirm that
// closed the dialog leaves the editor on screen, so the dialog's own mask is
// lifted before that last snapshot is taken.
func (this *FileExplorerHandler) verifySnapshot() {
	this.base.runner.tb.Helper()
	if !this.base.runner.DialogsOpen() {
		this.unmaskPathBar()
	}
	if !this.base.runner.SnapshotsEnabled() {
		return
	}
	if this.base.fixtureDirectory == "" {
		this.base.runner.tb.Fatal("file explorer snapshots need BaseHandler.WithFixtureDirectory")
	}

	this.base.runner.VerifySnapshot()
}
