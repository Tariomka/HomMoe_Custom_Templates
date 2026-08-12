package dialogs

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// saveFileSuffix is enforced on save targets so settings files keep their
// recognizable double extension.
const saveFileSuffix = ".gen.json"

// FileExplorerDialog is a self-contained, cross-platform file/folder picker
// rendered as a modal by the DialogHost. It replaces the OS-native pickers
// (gioui.org/x/explorer, PowerShell folder dialog, xdg-open) so behavior is
// identical on Windows, Linux and Steam Deck Big Picture, where no native file
// chooser is reliably available.
//
// It runs entirely inside the UI loop: results are delivered through the onPick
// / onSave callbacks (invoked synchronously when the user confirms) and the
// dialog then closes by returning done=true from Body. There is no blocking and
// there are no goroutines. Implements interfaces.IDialog structurally.
//
// The modes, the listing, the toolbar and the confirm flow live in the
// fileExplorerDialog*.go siblings.
type FileExplorerDialog struct {
	// fileSystem is the only route to the disk; the dialog itself performs no
	// path arithmetic and never imports os or path/filepath.
	fileSystem handler_interfaces.IFileSystemHandler

	mode  fileDialogMode
	title string

	// filterSuffixes restricts which files are listed in modeOpenFile
	// (case-insensitive suffix match). Empty means "all files".
	filterSuffixes []string

	currentDir string                  // "" => Windows drive-list view ("This PC")
	entries    []models.DirectoryEntry // cached listing for currentDir; rebuilt only on navigation
	listErr    string                  // inline read error; the dialog never panics
	showHidden bool                    // when false, hidden/system entries are filtered out

	selectedPath string // modeOpenFile: the highlighted file

	// Deferred navigation: a row click records the target here and it is applied
	// at the top of the next Body, so the listing slice is never mutated while
	// material.List is iterating it.
	pendingDir    string
	hasPendingDir bool

	// modeSaveFile overwrite-confirmation sub-state.
	overwriteActive bool
	saveErr         string // inline rejection of the typed filename

	// New-folder sub-state.
	newFolderActive bool
	newFolderErr    string

	list      widget.List
	rowClicks map[string]*widget.Clickable

	upBtn               widget.Clickable
	hiddenToggle        widget.Clickable
	newFolderBtn        widget.Clickable
	createFolderBtn     widget.Clickable
	cancelFolderBtn     widget.Clickable
	confirmBtn          widget.Clickable
	cancelBtn           widget.Clickable
	overwriteConfirmBtn widget.Clickable
	overwriteCancelBtn  widget.Clickable

	filenameEd  widget.Editor
	newFolderEd widget.Editor
	pathEd      widget.Editor

	onPick func(path string)
	onSave func(path string)
}

func newFileExplorerDialog(
	fileSystem handler_interfaces.IFileSystemHandler,
	mode fileDialogMode,
	title string) *FileExplorerDialog {
	dialog := &FileExplorerDialog{
		fileSystem: fileSystem,
		mode:       mode,
		title:      title,
		rowClicks:  map[string]*widget.Clickable{},
	}
	dialog.list.Axis = layout.Vertical
	dialog.filenameEd.SingleLine = true
	dialog.newFolderEd.SingleLine = true
	dialog.pathEd.SingleLine = true
	return dialog
}

func (this *FileExplorerDialog) Title() string { return this.title }

func (this *FileExplorerDialog) PreferredSize() (width, height unit.Dp) {
	return unit.Dp(720), unit.Dp(560)
}

func (this *FileExplorerDialog) Body(gtx layout.Context, theme *material.Theme) (layout.Dimensions, bool) {
	// Apply any navigation deferred from a previous frame's row click before
	// reading or laying out the listing.
	if this.hasPendingDir {
		this.hasPendingDir = false
		this.loadDir(this.pendingDir)
	}

	if this.cancelBtn.Clicked(gtx) {
		return layout.Dimensions{Size: gtx.Constraints.Min}, true
	}

	if this.upBtn.Clicked(gtx) {
		this.loadDir(this.parentDir())
	}

	if this.hiddenToggle.Clicked(gtx) {
		this.showHidden = !this.showHidden
		this.loadDir(this.currentDir)
	}

	if this.canModify() && this.newFolderBtn.Clicked(gtx) {
		this.newFolderActive = !this.newFolderActive
		this.newFolderErr = ""
		if this.newFolderActive {
			this.newFolderEd.SetText("")
		}
	}

	if this.canModify() && this.newFolderActive {
		if this.createFolderBtn.Clicked(gtx) {
			this.tryCreateFolder()
		}
		if this.cancelFolderBtn.Clicked(gtx) {
			this.newFolderActive = false
			this.newFolderErr = ""
		}
	}

	if this.handleConfirm(gtx) {
		return layout.Dimensions{Size: gtx.Constraints.Min}, true
	}

	return this.getContentWidget(theme)(gtx), false
}

// getContentWidget stacks the header, the scrollable listing, an optional error
// line, the mode-specific input rows and the footer.
func (this *FileExplorerDialog) getContentWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(this.getHeaderWidget(theme)),
			layout.Rigid(widgets.NewVerticalSpacerWidget(8)),
			layout.Flexed(1, this.getListWidget(theme)),
			layout.Rigid(this.getErrorLineWidget(theme)),
			layout.Rigid(this.getSaveRowWidget(theme)),
			layout.Rigid(this.getNewFolderRowWidget(theme)),
			layout.Rigid(widgets.NewVerticalSpacerWidget(10)),
			layout.Rigid(this.getFooterWidget(theme)),
		)
	}
}

// parentDir returns the directory to ascend to, or the current directory when
// already at the top (so the Up button can be disabled). On Windows the volume
// root ascends to the synthetic drive list ("").
func (this *FileExplorerDialog) parentDir() string {
	return this.fileSystem.ParentDirectory(this.currentDir)
}

func (this *FileExplorerDialog) requestNav(dir string) {
	this.pendingDir = dir
	this.hasPendingDir = true
}

func (this *FileExplorerDialog) onEntryClicked(entry models.DirectoryEntry) {
	if entry.IsDir {
		this.requestNav(entry.Path)
		return
	}

	switch this.mode {
	case modeOpenFile:
		this.selectedPath = entry.Path
	case modeSaveFile:
		this.filenameEd.SetText(entry.Name)
		this.overwriteActive = false
		this.saveErr = ""
	case modePickFolder, modeBrowse: // noop
	}
}

// resolveSaveTarget builds the absolute save path from the filename field,
// stripping any directory components for safety and enforcing the .gen.json
// suffix. ok is false when the field is empty.
func (this *FileExplorerDialog) resolveSaveTarget() (string, bool) {
	return this.fileSystem.ResolveSaveTarget(this.currentDir, this.filenameEd.Text(), saveFileSuffix)
}

func (this *FileExplorerDialog) canModify() bool {
	return (this.mode == modeSaveFile || this.mode == modePickFolder) && this.currentDir != ""
}

func (this *FileExplorerDialog) clickFor(path string) *widget.Clickable {
	clk := this.rowClicks[path]
	if clk == nil {
		clk = &widget.Clickable{}
		this.rowClicks[path] = clk
	}
	return clk
}

func (this *FileExplorerDialog) resetScroll() {
	this.list.Position = layout.Position{}
}
