package dialogs

import (
	"image"

	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
)

func (this *FileExplorerDialog) getFooterWidget(theme *material.Theme) layout.Widget {
	if this.mode == modeSaveFile && this.overwriteActive {
		return func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, widgets.NewLabelBuilder(theme).WithSizeBig().
					WithText("File already exists. Overwrite?").
					WithColor(themes.ColorsBase.WarnText).WithMaxLines(1).Build),
				layout.Rigid(widgets.NewButtonWidget(theme, "Cancel", &this.overwriteCancelBtn, false)),
				widgets.NewDefaultComponentSpacer(),
				layout.Rigid(widgets.NewBrightButtonWidget(theme, "Overwrite", &this.overwriteConfirmBtn, false)),
			)
		}
	}

	confirmLabel, showConfirm, confirmDisabled := this.confirmButtonState()
	cancelLabel := "Cancel"
	if this.mode == modeBrowse {
		cancelLabel = "Close"
	}

	children := make([]layout.FlexChild, 0, 5)
	if this.canModify() {
		children = append(
			children,
			layout.Rigid(widgets.NewButtonWidget(theme, "New Folder", &this.newFolderBtn, false)),
		)
	}
	children = append(children, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, 0)}
	}))
	children = append(children, layout.Rigid(widgets.NewButtonWidget(theme, cancelLabel, &this.cancelBtn, false)))
	if showConfirm {
		children = append(children,
			widgets.NewDefaultComponentSpacer(),
			layout.Rigid(widgets.NewBrightButtonWidget(theme, confirmLabel, &this.confirmBtn, confirmDisabled)),
		)
	}
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx, children...)
	}
}

// confirmButtonState returns the primary button's label, whether it is shown at
// all (browse mode has none) and whether it is currently disabled.
func (this *FileExplorerDialog) confirmButtonState() (label string, show bool, disabled bool) {
	switch this.mode {
	case modeOpenFile:
		return "Open", true, this.selectedPath == ""
	case modeSaveFile:
		if !this.hasResolvedSaveName() {
			return "Save", true, true
		}

		_, ok := this.resolveSaveTarget()
		return "Save", true, !ok
	case modePickFolder:
		return "Select This Folder", true, this.currentDir == ""
	case modeBrowse:
		fallthrough
	default:
		return "", false, false
	}
}

// handleConfirm processes the mode-specific confirm and overwrite buttons and
// reports whether the dialog should close.
func (this *FileExplorerDialog) handleConfirm(gtx layout.Context) bool {
	switch this.mode {
	case modeOpenFile:
		if this.confirmBtn.Clicked(gtx) && this.selectedPath != "" {
			if this.onPick != nil {
				this.onPick(this.selectedPath)
			}
			return true
		}
	case modePickFolder:
		if this.confirmBtn.Clicked(gtx) && this.currentDir != "" {
			if this.onPick != nil {
				this.onPick(this.currentDir)
			}
			return true
		}
	case modeSaveFile:
		if this.overwriteActive {
			return this.confirmOverwrite(gtx)
		}

		return this.confirmSelection(gtx)
	case modeBrowse: // noop
	}

	return false
}

// confirmOverwrite processes the overwrite prompt's buttons and reports
// whether the dialog should close.
func (this *FileExplorerDialog) confirmOverwrite(gtx layout.Context) bool {
	if this.overwriteConfirmBtn.Clicked(gtx) {
		this.overwriteActive = false

		path, ok := this.resolveSaveTarget()
		if !ok {
			// Filename was cleared while the prompt was up; abandon it.
			return false
		}

		if this.onSave != nil {
			this.onSave(path)
		}
		return true
	}

	if this.overwriteCancelBtn.Clicked(gtx) {
		this.overwriteActive = false
	}

	return false
}

// confirmSelection handles the save button: an existing file opens the
// overwrite prompt, an existing folder is refused, otherwise the file is saved
// immediately.
func (this *FileExplorerDialog) confirmSelection(gtx layout.Context) bool {
	if !this.confirmBtn.Clicked(gtx) {
		return false
	}

	this.saveErr = ""

	path, ok := this.resolveSaveTarget()
	if !ok {
		return false
	}

	if this.fileSystem.DirectoryExists(path) {
		this.saveErr = "A folder with that name already exists."
		return false
	}

	if this.fileSystem.PathExists(path) {
		this.overwriteActive = true
		return false
	}

	if this.onSave != nil {
		this.onSave(path)
	}
	return true
}
