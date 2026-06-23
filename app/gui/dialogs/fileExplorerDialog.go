package dialogs

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
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

// fileEntry is one selectable row in the listing.
type fileEntry struct {
	name  string
	path  string
	isDir bool
}

// FileExplorerDialog is a self-contained, cross-platform file/folder picker
// rendered as a modal by the DialogHost. It replaces the OS-native pickers
// (gioui.org/x/explorer, PowerShell folder dialog, xdg-open) so behavior is
// identical on Windows, Linux and Steam Deck Big Picture, where no native file
// chooser is reliably available.
//
// It runs entirely inside the UI loop: results are delivered through the onPick
// / onSave callbacks (invoked synchronously when the user confirms) and the
// dialog then closes by returning done=true from Body. There is no blocking and
// there are no goroutines. Implements interfaces.Dialog structurally.
type FileExplorerDialog struct {
	mode  fileDialogMode
	title string

	// filterSuffixes restricts which files are listed in modeOpenFile
	// (case-insensitive suffix match). Empty means "all files".
	filterSuffixes []string

	currentDir string      // "" => Windows drive-list view ("This PC")
	entries    []fileEntry // cached listing for currentDir; rebuilt only on navigation
	listErr    string      // inline read error; the dialog never panics
	showHidden bool        // when false, hidden/system entries are filtered out

	selectedPath string // modeOpenFile: the highlighted file

	// Deferred navigation: a row click records the target here and it is applied
	// at the top of the next Body, so the listing slice is never mutated while
	// material.List is iterating it.
	pendingDir    string
	hasPendingDir bool

	// modeSaveFile overwrite-confirmation sub-state.
	overwriteActive bool

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

	onPick func(path string)
	onSave func(path string)
}

func newFileExplorerDialog(mode fileDialogMode, title string) *FileExplorerDialog {
	dialog := &FileExplorerDialog{
		mode:      mode,
		title:     title,
		rowClicks: map[string]*widget.Clickable{},
	}
	dialog.list.Axis = layout.Vertical
	dialog.filenameEd.SingleLine = true
	dialog.newFolderEd.SingleLine = true
	return dialog
}

// NewOpenFileDialog builds a single-file picker starting at initialDir. Only
// files whose name ends with one of filterSuffixes (case-insensitive) are shown;
// pass nil to list every file. onPick receives the chosen absolute file path.
func NewOpenFileDialog(initialDir string, filterSuffixes []string, onPick func(path string)) *FileExplorerDialog {
	dialog := newFileExplorerDialog(modeOpenFile, "Open File")
	dialog.filterSuffixes = filterSuffixes
	dialog.onPick = onPick
	dialog.loadDir(resolveInitialDir(initialDir))
	return dialog
}

// NewSaveFileDialog builds a save-location picker starting at initialDir with
// the filename field prefilled to defaultName. onSave receives the chosen
// absolute path (with a guaranteed .gen.json suffix) after any overwrite
// confirmation.
func NewSaveFileDialog(initialDir, defaultName string, onSave func(path string)) *FileExplorerDialog {
	dialog := newFileExplorerDialog(modeSaveFile, "Save File")
	dialog.onSave = onSave
	dialog.filenameEd.SetText(defaultName)
	dialog.loadDir(resolveInitialDir(initialDir))
	return dialog
}

// NewPickFolderDialog builds a single-folder picker starting at initialDir.
// onPick receives the directory the user navigated into and confirmed.
func NewPickFolderDialog(initialDir string, onPick func(dir string)) *FileExplorerDialog {
	dialog := newFileExplorerDialog(modePickFolder, "Select Folder")
	dialog.onPick = onPick
	dialog.loadDir(resolveInitialDir(initialDir))
	return dialog
}

// NewBrowseDialog builds a read-only viewer starting at initialDir; it has no
// confirm action and is used to inspect a directory's contents in-app.
func NewBrowseDialog(initialDir string) *FileExplorerDialog {
	dialog := newFileExplorerDialog(modeBrowse, "Browse")
	dialog.loadDir(resolveInitialDir(initialDir))
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

	if this.canModify() {
		if this.newFolderBtn.Clicked(gtx) {
			this.newFolderActive = !this.newFolderActive
			this.newFolderErr = ""
			if this.newFolderActive {
				this.newFolderEd.SetText("")
			}
		}
		if this.newFolderActive {
			if this.createFolderBtn.Clicked(gtx) {
				this.tryCreateFolder()
			}
			if this.cancelFolderBtn.Clicked(gtx) {
				this.newFolderActive = false
				this.newFolderErr = ""
			}
		}
	}

	if this.handleConfirm(gtx) {
		return layout.Dimensions{Size: gtx.Constraints.Min}, true
	}

	return this.layoutContent(gtx, theme), false
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
			if this.overwriteConfirmBtn.Clicked(gtx) {
				if path, ok := this.resolveSaveTarget(); ok {
					this.overwriteActive = false
					if this.onSave != nil {
						this.onSave(path)
					}
					return true
				}
				// Filename was cleared while the prompt was up; abandon it.
				this.overwriteActive = false
			}
			if this.overwriteCancelBtn.Clicked(gtx) {
				this.overwriteActive = false
			}
		} else if this.confirmBtn.Clicked(gtx) {
			if path, ok := this.resolveSaveTarget(); ok {
				if pathExists(path) {
					this.overwriteActive = true
				} else {
					if this.onSave != nil {
						this.onSave(path)
					}
					return true
				}
			}
		}
	}
	return false
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
	this.selectedPath = ""

	if dir == "" {
		this.currentDir = ""
		this.entries = listWindowsDrives()
		this.rowClicks = map[string]*widget.Clickable{}
		this.resetScroll()
		return
	}

	osEntries, err := os.ReadDir(dir)
	if err != nil {
		// Keep the current location; only adopt dir on the very first load so
		// the path bar is coherent when the initial directory is unreadable.
		if this.currentDir == "" {
			this.currentDir = dir
		}
		this.listErr = err.Error()
		return
	}

	dirs := make([]fileEntry, 0, len(osEntries))
	files := make([]fileEntry, 0, len(osEntries))
	for _, entry := range osEntries {
		name := entry.Name()
		if !this.showHidden && this.isHidden(name, entry) {
			continue
		}
		full := filepath.Join(dir, name)
		if entry.IsDir() {
			dirs = append(dirs, fileEntry{name: name, path: full, isDir: true})
			continue
		}
		if this.mode == modeOpenFile && !this.matchesFilter(name) {
			continue
		}
		files = append(files, fileEntry{name: name, path: full, isDir: false})
	}
	sort.SliceStable(dirs, func(i, j int) bool {
		return strings.ToLower(dirs[i].name) < strings.ToLower(dirs[j].name)
	})
	sort.SliceStable(files, func(i, j int) bool {
		return strings.ToLower(files[i].name) < strings.ToLower(files[j].name)
	})

	this.currentDir = dir
	this.entries = append(dirs, files...)
	this.rowClicks = map[string]*widget.Clickable{}
	this.resetScroll()
}

// isHidden reports whether an entry should be filtered out while showHidden is
// off: dotfiles everywhere, plus hidden/system entries on Windows.
func (this *FileExplorerDialog) isHidden(name string, entry os.DirEntry) bool {
	if strings.HasPrefix(name, ".") {
		return true
	}
	info, err := entry.Info()
	if err != nil {
		return false
	}
	return hasHiddenAttr(info)
}

func (this *FileExplorerDialog) matchesFilter(name string) bool {
	if len(this.filterSuffixes) == 0 {
		return true
	}
	lower := strings.ToLower(name)
	for _, suffix := range this.filterSuffixes {
		if strings.HasSuffix(lower, strings.ToLower(suffix)) {
			return true
		}
	}
	return false
}

// parentDir returns the directory to ascend to, or the current directory when
// already at the top (so the Up button can be disabled). On Windows the volume
// root ascends to the synthetic drive list ("").
func (this *FileExplorerDialog) parentDir() string {
	if this.currentDir == "" {
		return ""
	}
	parent := filepath.Dir(this.currentDir)
	if parent == this.currentDir {
		if runtime.GOOS == "windows" {
			return ""
		}
		return this.currentDir
	}
	return parent
}

func (this *FileExplorerDialog) requestNav(dir string) {
	this.pendingDir = dir
	this.hasPendingDir = true
}

func (this *FileExplorerDialog) onEntryClicked(entry fileEntry) {
	if entry.isDir {
		this.requestNav(entry.path)
		return
	}
	switch this.mode {
	case modeOpenFile:
		this.selectedPath = entry.path
	case modeSaveFile:
		this.filenameEd.SetText(entry.name)
		this.overwriteActive = false
	}
}

// resolveSaveTarget builds the absolute save path from the filename field,
// stripping any directory components for safety and enforcing the .gen.json
// suffix. ok is false when the field is empty.
func (this *FileExplorerDialog) resolveSaveTarget() (string, bool) {
	name := strings.TrimSpace(this.filenameEd.Text())
	if name == "" {
		return "", false
	}
	name = filepath.Base(name)
	if name == "." || name == ".." || name == string(os.PathSeparator) {
		return "", false
	}
	if !strings.HasSuffix(strings.ToLower(name), saveFileSuffix) {
		name += saveFileSuffix
	}
	return filepath.Join(this.currentDir, name), true
}

func (this *FileExplorerDialog) tryCreateFolder() {
	name := strings.TrimSpace(this.newFolderEd.Text())
	if name == "" {
		this.newFolderErr = "Enter a folder name."
		return
	}
	if name == "." || name == ".." ||
		strings.ContainsRune(name, '/') ||
		strings.ContainsRune(name, os.PathSeparator) ||
		strings.ContainsRune(name, ':') {
		this.newFolderErr = "Invalid folder name."
		return
	}
	target := filepath.Join(this.currentDir, name)
	if err := os.Mkdir(target, 0o755); err != nil {
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

// listWindowsDrives enumerates the existing fixed/removable drive roots. It is
// only ever shown after ascending past a Windows volume root; on other systems
// it returns an empty slice (and is never reached).
func listWindowsDrives() []fileEntry {
	drives := make([]fileEntry, 0, 26)
	for letter := 'A'; letter <= 'Z'; letter++ {
		root := string(letter) + `:\`
		if _, err := os.Stat(root); err == nil {
			drives = append(drives, fileEntry{name: root, path: root, isDir: true})
		}
	}
	return drives
}

// resolveInitialDir normalizes the requested start directory and guarantees an
// existing directory is returned: it climbs to the nearest existing ancestor
// (also handling being handed a file path), then falls back to the user's home
// directory and finally the working directory.
func resolveInitialDir(dir string) string {
	dir = strings.TrimSpace(dir)
	if dir != "" {
		if abs, err := filepath.Abs(dir); err == nil {
			dir = abs
		}
		dir = filepath.Clean(dir)
		for candidate := dir; ; {
			if info, err := os.Stat(candidate); err == nil && info.IsDir() {
				return candidate
			}
			parent := filepath.Dir(candidate)
			if parent == candidate {
				break
			}
			candidate = parent
		}
	}
	if home, err := os.UserHomeDir(); err == nil {
		if info, err := os.Stat(home); err == nil && info.IsDir() {
			return home
		}
	}
	if workingDir, err := os.Getwd(); err == nil {
		return workingDir
	}
	return "."
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
