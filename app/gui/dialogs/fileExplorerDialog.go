package dialogs

import (
	"errors"
	"fmt"
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// fileDialogMode selects which task the explorer performs.
type fileDialogMode uint8

const (
	modeOpenFile fileDialogMode = iota
	modeSaveFile
	modePickFolder
	modeBrowse
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

// NewOpenFileDialog builds a single-file picker starting at initialDir. Only
// files whose name ends with one of filterSuffixes (case-insensitive) are shown;
// pass nil to list every file. onPick receives the chosen absolute file path.
func NewOpenFileDialog(
	fileSystem handler_interfaces.IFileSystemHandler,
	initialDir string,
	filterSuffixes []string,
	onPick func(path string)) *FileExplorerDialog {
	dialog := newFileExplorerDialog(fileSystem, modeOpenFile, "Open File")
	dialog.filterSuffixes = filterSuffixes
	dialog.onPick = onPick
	dialog.loadDir(fileSystem.ResolveStartDirectory(initialDir))
	return dialog
}

// NewSaveFileDialog builds a save-location picker starting at initialDir with
// the filename field prefilled to defaultName. onSave receives the chosen
// absolute path (with a guaranteed .gen.json suffix) after any overwrite
// confirmation.
func NewSaveFileDialog(
	fileSystem handler_interfaces.IFileSystemHandler,
	initialDir, defaultName string,
	onSave func(path string)) *FileExplorerDialog {
	dialog := newFileExplorerDialog(fileSystem, modeSaveFile, "Save File")
	dialog.onSave = onSave
	dialog.filenameEd.SetText(defaultName)
	dialog.loadDir(fileSystem.ResolveStartDirectory(initialDir))
	return dialog
}

// NewPickFolderDialog builds a single-folder picker starting at initialDir.
// onPick receives the directory the user navigated into and confirmed.
func NewPickFolderDialog(
	fileSystem handler_interfaces.IFileSystemHandler,
	initialDir string,
	onPick func(dir string)) *FileExplorerDialog {
	dialog := newFileExplorerDialog(fileSystem, modePickFolder, "Select Folder")
	dialog.onPick = onPick
	dialog.loadDir(fileSystem.ResolveStartDirectory(initialDir))
	return dialog
}

// NewBrowseDialog builds a read-only viewer starting at initialDir; it has no
// confirm action and is used to inspect a directory's contents in-app.
func NewBrowseDialog(
	fileSystem handler_interfaces.IFileSystemHandler,
	initialDir string) *FileExplorerDialog {
	dialog := newFileExplorerDialog(fileSystem, modeBrowse, "Browse")
	dialog.loadDir(fileSystem.ResolveStartDirectory(initialDir))
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

func (this *FileExplorerDialog) getHeaderWidget(theme *material.Theme) layout.Widget {
	upDisabled := this.parentDir() == this.currentDir
	if this.pathEd.Text() != this.currentDir {
		if this.currentDir == "" {
			this.pathEd.SetText("This PC")
		} else {
			this.pathEd.SetText(this.currentDir)
		}
	}

	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(widgets.NewButtonWidget(theme, "Up", &this.upBtn, upDisabled)),
			widgets.NewDefaultComponentSpacer(),
			layout.Flexed(1, widgets.NewTextboxWidget(theme, &this.pathEd, "Current directory", true)),
			widgets.NewDefaultComponentSpacer(),
			layout.Rigid(widgets.NewToggleButtonWidget(theme, "Show hidden", &this.hiddenToggle, this.showHidden)),
		)
	}
}

func (this *FileExplorerDialog) getListWidget(theme *material.Theme) layout.Widget {
	if len(this.entries) == 0 {
		message := "(empty folder)"
		if this.listErr != "" {
			message = "(unable to read folder)"
		}
		return func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, widgets.NewLabelBigWidget(theme, message, themes.ColorsBase.TextDim))
		}
	}

	return func(gtx layout.Context) layout.Dimensions {
		return material.List(theme, &this.list).Layout(gtx, len(this.entries),
			func(gtx layout.Context, index int) layout.Dimensions {
				return this.getEntryRowWidget(theme, this.entries[index])(gtx)
			})
	}
}

func (this *FileExplorerDialog) getEntryRowWidget(
	theme *material.Theme,
	entry models.DirectoryEntry) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		clk := this.clickFor(entry.Path)
		if clk.Clicked(gtx) {
			this.onEntryClicked(entry)
		}
		selected := !entry.IsDir && entry.Path == this.selectedPath
		badgeText := ""
		textColor := themes.ColorsBase.Text
		if entry.IsDir {
			badgeText = "DIR"
			textColor = themes.ColorsBase.AccentBright
		}
		return material.Clickable(gtx, clk, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			dims := layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(6), Right: unit.Dp(6)}.
				Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Dp(unit.Dp(38))
							return widgets.NewStyledLabelWidget(
								theme, badgeText, themes.ColorsBase.Accent, font.Font{Weight: font.SemiBold})(gtx)
						}),
						layout.Flexed(1, widgets.NewLabelBuilder(theme).WithSizeBig().
							WithText(entry.Name).WithColor(textColor).WithMaxLines(1).Build))
				})
			call := macro.Stop()
			if selected {
				paint.FillShape(gtx.Ops, themes.ColorsBase.Selection, clip.Rect{Max: dims.Size}.Op())
			} else if clk.Hovered() {
				paint.FillShape(gtx.Ops, themes.ColorsBase.Hover, clip.Rect{Max: dims.Size}.Op())
			}
			call.Add(gtx.Ops)
			return dims
		})
	}
}

func (this *FileExplorerDialog) getErrorLineWidget(theme *material.Theme) layout.Widget {
	message := this.listErr
	if message == "" {
		message = this.saveErr
	}

	if message == "" {
		return widgets.NewEmptyWidget()
	}

	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: constants.DefaultPaddingSmall}.
			Layout(gtx, widgets.NewLabelBuilder(theme).WithSizeDefault().
				WithText(message).WithColor(themes.ColorsBase.Error).WithMaxLines(2).Build)
	}
}

func (this *FileExplorerDialog) getSaveRowWidget(theme *material.Theme) layout.Widget {
	if this.mode != modeSaveFile {
		return widgets.NewEmptyWidget()
	}

	hint := fmt.Sprintf("filename%s", saveFileSuffix)
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: constants.DefaultPadding}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(widgets.NewLabelBigWidget(theme, "Save as:", themes.ColorsBase.TextDim)),
				widgets.NewDefaultComponentSpacer(),
				layout.Flexed(1, widgets.NewTextboxWidget(theme, &this.filenameEd, hint, false)),
			)
		})
	}
}

func (this *FileExplorerDialog) getNewFolderRowWidget(theme *material.Theme) layout.Widget {
	if !this.canModify() || !this.newFolderActive {
		return widgets.NewEmptyWidget()
	}

	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: constants.DefaultPadding}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(widgets.NewLabelBigWidget(theme, "New folder:", themes.ColorsBase.TextDim)),
						widgets.NewDefaultComponentSpacer(),
						layout.Flexed(1, widgets.NewTextboxWidget(theme, &this.newFolderEd, "folder name", false)),
						widgets.NewDefaultComponentSpacer(),
						layout.Rigid(widgets.NewButtonWidget(theme, "Create", &this.createFolderBtn, false)),
						widgets.NewDefaultComponentSpacer(),
						layout.Rigid(widgets.NewButtonWidget(theme, "Cancel", &this.cancelFolderBtn, false)),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if this.newFolderErr == "" {
						return layout.Dimensions{}
					}

					return layout.Inset{Top: constants.DefaultPaddingSmall - 2}.
						Layout(gtx, widgets.NewLabelBuilder(theme).WithSizeDefault().
							WithText(this.newFolderErr).WithColor(themes.ColorsBase.Error).WithMaxLines(2).Build)
				}),
			)
		})
	}
}

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

// loadDir reads dir and replaces the cached listing. It is the single
// navigation primitive and is called only on open, on explicit navigation and
// on toggling hidden files - never per frame. On failure the previous listing
// is kept and the error is surfaced inline so navigation into an unreadable
// folder leaves the user where they were.
func (this *FileExplorerDialog) loadDir(dir string) {
	this.newFolderActive = false
	this.newFolderErr = ""
	this.overwriteActive = false
	this.listErr = ""
	this.saveErr = ""
	this.selectedPath = ""

	if dir == "" {
		this.currentDir = ""
		this.entries = this.fileSystem.ListRoots()
		this.rowClicks = map[string]*widget.Clickable{}
		this.resetScroll()
		return
	}

	entries, err := this.fileSystem.ListEntries(dir, this.filterSuffixes, this.showHidden)
	if err != nil {
		// Keep the current location; only adopt dir on the very first load so
		// the path bar is coherent when the initial directory is unreadable.
		if this.currentDir == "" {
			this.currentDir = dir
		}
		this.listErr = err.Error()
		return
	}

	this.currentDir = dir
	this.entries = entries
	this.rowClicks = map[string]*widget.Clickable{}
	this.resetScroll()
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

func (this *FileExplorerDialog) tryCreateFolder() {
	target, err := this.fileSystem.CreateDirectory(this.currentDir, this.newFolderEd.Text())
	switch {
	case errors.Is(err, common_errors.ErrDirectoryNameEmpty):
		this.newFolderErr = "Enter a folder name."
		return
	case errors.Is(err, common_errors.ErrDirectoryNameInvalid):
		this.newFolderErr = "Invalid folder name."
		return
	case err != nil:
		this.newFolderErr = err.Error()
		return
	}

	this.newFolderActive = false
	this.newFolderErr = ""
	this.loadDir(target)
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
