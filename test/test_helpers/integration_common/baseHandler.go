//go:build integration_test

package integration_common

import (
	"os"
	"path/filepath"

	"gioui.org/f32"
)

// BaseHandler drives the editor window through the AppRunner with a fluent API,
// standing in for the Playwright/Selenium-style driver that Gio has no
// equivalent of. Per-tab and per-dialog handlers, the state they need to track,
// and the scrolling support they depend on are designed in
// plans/gui-handler-framework.md; read that before adding to this file.
type BaseHandler struct {
	runner *AppRunner

	isRandomTopology bool

	// isExperimentalMapSizes tracks the checkbox that grows the map size
	// dropdown from 11 rows to 28, because the row a coordinate addresses is
	// only valid against the option count the checkbox implies.
	isExperimentalMapSizes bool

	// fixtureDirectory is the throwaway directory the file dialogs open at once
	// WithFixtureDirectory has seeded one; empty means they open wherever the
	// machine's editor state happens to live.
	fixtureDirectory string
}

func NewHandler(runner *AppRunner) *BaseHandler {
	runner.tb.Helper()
	handler := BaseHandler{runner: runner}
	handler.setRandomTopology()
	return &handler
}

func (this *BaseHandler) WithSnapshots() *BaseHandler {
	this.runner.tb.Helper()
	this.runner.EnableSnapshots()
	return this
}

// WithFixtureDirectory points the Load and Save To dialogs at an empty
// throwaway directory, so a listing golden shows what the test put there instead
// of whatever the machine's own templates directory happens to contain
// (AGENTS.md 2.7). Call it before WithSnapshots: the toolbar reports the seeded
// path, so the mask it registers has to cover every snapshot the test takes.
func (this *BaseHandler) WithFixtureDirectory() *BaseHandler {
	this.runner.tb.Helper()
	this.ensureFixtureDirectory()
	return this
}

// WithFixtureFiles seeds the fixture directory with the named files, creating
// the directory if no builder has yet. A name starting with "." is hidden on
// both Windows and Linux; only a name ending in .gen.json is listed by the open
// dialog, which filters on that suffix.
func (this *BaseHandler) WithFixtureFiles(names ...string) *BaseHandler {
	this.runner.tb.Helper()
	directory := this.ensureFixtureDirectory()
	for _, name := range names {
		if err := os.WriteFile(filepath.Join(directory, name), []byte("{}"), 0o600); err != nil {
			this.runner.tb.Fatalf("cannot create fixture file %s: %v", name, err)
		}
	}

	return this
}

// WithFixtureFolders seeds the fixture directory with the named folders. Folders
// are listed whatever they are called, so this is also how a folder named like a
// settings file is set up.
func (this *BaseHandler) WithFixtureFolders(names ...string) *BaseHandler {
	this.runner.tb.Helper()
	directory := this.ensureFixtureDirectory()
	for _, name := range names {
		if err := os.Mkdir(filepath.Join(directory, name), 0o750); err != nil {
			this.runner.tb.Fatalf("cannot create fixture folder %s: %v", name, err)
		}
	}

	return this
}

// FixtureDirectory returns the directory the fixture builders seeded, so a test
// can build the absolute paths it expects the dialog to report.
func (this *BaseHandler) FixtureDirectory() string {
	this.runner.tb.Helper()
	return this.fixtureDirectory
}

func (this *BaseHandler) ClickGeneralTab() *GeneralTabHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(generalTabX, tabStripCenterY))
	this.runner.VerifySnapshot()
	return &GeneralTabHandler{BaseHandler: this}
}

func (this *BaseHandler) ClickLayoutAndZonesTab() *LayoutAndZonesTabHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(layoutAndZonesTabX, tabStripCenterY))
	this.runner.VerifySnapshot()
	return &LayoutAndZonesTabHandler{BaseHandler: this}
}

func (this *BaseHandler) ClickBonusesAndBansTab() *BonusesAndBansTabHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(bonusesAndBansTabX, tabStripCenterY))
	this.runner.VerifySnapshot()
	return &BonusesAndBansTabHandler{BaseHandler: this}
}

// ClickNew discards the editor state and starts a fresh template. It takes no
// snapshot: without a fixture directory the toolbar reports a per-machine path.
func (this *BaseHandler) ClickNew() *BaseHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(newButtonX, toolbarCenterY))
	return this
}

func (this *BaseHandler) ClickLoad() *FileExplorerHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(loadButtonX, toolbarCenterY))
	return newFileExplorerHandler(this)
}

func (this *BaseHandler) ClickSaveTo() *FileExplorerHandler {
	this.runner.tb.Helper()
	this.runner.ClickAt(f32.Pt(saveToButtonX, toolbarCenterY))
	return newFileExplorerHandler(this)
}

// ScrollPanel turns the mouse wheel over the settings panel by delta pixels;
// positive scrolls the content up. Gio clamps to the panel's scrollable range,
// so an oversized delta scrolls to the end.
func (this *BaseHandler) ScrollPanel(delta float32) *BaseHandler {
	this.runner.tb.Helper()
	this.runner.Scroll(f32.Pt(panelScrollX, panelScrollY), f32.Pt(0, delta))
	this.runner.VerifySnapshot()
	return this
}

// commit renders one more frame so the panel's SaveToState runs. A click is
// processed on the frame it is queued against, but the panel only writes the
// resulting widget values back to the editor state on the layout after that.
func (this *BaseHandler) commit() {
	this.runner.tb.Helper()
	this.runner.NextFrame()
}

// ensureFixtureDirectory creates the fixture directory on first use and points
// the file dialogs at it.
func (this *BaseHandler) ensureFixtureDirectory() string {
	this.runner.tb.Helper()
	if this.fixtureDirectory != "" {
		return this.fixtureDirectory
	}

	this.fixtureDirectory = this.runner.tb.TempDir()
	this.runner.SetCurrentPath(this.fixtureDirectory)
	this.runner.MaskRect(fileStatusMask())
	return this.fixtureDirectory
}

// setRandomTopology records that the editor starts on the Random topology,
// whose preview is regenerated with fresh randomness on every run, and masks the
// three regions of the preview panel that cannot be compared against a golden.
func (this *BaseHandler) setRandomTopology() *BaseHandler {
	this.runner.tb.Helper()
	this.isRandomTopology = true
	this.runner.MaskRect(previewCanvasMask())
	this.runner.MaskRect(statusMessageMask())
	this.runner.MaskRect(outputDirectoryMask())
	return this
}
