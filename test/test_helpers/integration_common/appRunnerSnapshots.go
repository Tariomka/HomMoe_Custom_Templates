//go:build integration_test

package integration_common

import (
	"image"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gioui.org/gpu/headless"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common/snapshot"
)

// EnableSnapshots turns on golden-snapshot verification for this runner: after
// every injected action (ClickAt, DragTo, InputText) a screenshot is rendered
// through a headless GPU window, masked, and either saved as a golden (-update)
// or validated against the stored golden (headless runs). Headed runs without
// -update leave snapshotting inert. Benchmarks simply never call this.
//
// The golden path is
// __snapshots__/{caller test file name}/{test name}_{action number}.golden.
func (this *AppRunner) EnableSnapshots() {
	this.tb.Helper()
	_, callerFile, _, ok := runtime.Caller(2)
	if !ok {
		this.tb.Fatal("EnableSnapshots: cannot resolve the calling test file")
	}

	this.snapshotFile = strings.TrimSuffix(filepath.Base(callerFile), ".go")
	this.comparer = snapshot.NewComparer()
	this.store = snapshot.NewStore()

	if !IsHeadless() && !IsUpdate() {
		return // headed without -update: no capture, no validation.
	}

	headlessWindow, err := headless.NewWindow(WindowWidth, WindowHeight)
	if err != nil {
		this.tb.Fatalf("EnableSnapshots: cannot create headless GPU window: %v", err)
	}
	this.headlessWin = headlessWindow
	this.tb.Cleanup(headlessWindow.Release)
}

// MaskRect registers a window-pixel rectangle to be blanked out of every
// snapshot before saving or comparing. Use it to hide nondeterministic regions
// such as the randomly generated map preview.
func (this *AppRunner) MaskRect(rect image.Rectangle) {
	this.tb.Helper()
	this.masker.AddRect(rect)
}

// UnmaskRect lifts a rectangle registered with MaskRect, for a region that is
// only nondeterministic while something transient is on screen.
func (this *AppRunner) UnmaskRect(rect image.Rectangle) {
	this.tb.Helper()
	if !this.masker.RemoveRect(rect) {
		this.tb.Fatalf("UnmaskRect: %v was never masked", rect)
	}
}

// SnapshotsEnabled reports whether EnableSnapshots armed this runner, including
// the headed run that captures nothing. Handlers use it to refuse an action they
// cannot produce a machine-independent golden for.
func (this *AppRunner) SnapshotsEnabled() bool {
	this.tb.Helper()
	return this.snapshotFile != ""
}

// VerifySnapshot renders, masks and then saves (-update) or validates the
// screenshot for the action that just completed. No-op unless EnableSnapshots
// armed this runner with a headless window.
func (this *AppRunner) VerifySnapshot() {
	this.tb.Helper()
	if this.headlessWin == nil {
		return
	}

	this.actionCount++
	screenshot := this.captureScreenshot()
	this.masker.Apply(screenshot)

	goldenPath := this.store.GoldenPath(this.snapshotFile, this.tb.Name(), this.actionCount)
	if IsUpdate() {
		if err := this.store.SaveGolden(goldenPath, screenshot); err != nil {
			this.tb.Fatalf("snapshot %d: cannot save golden %s: %v", this.actionCount, goldenPath, err)
		}
		return
	}

	this.validateScreenshot(goldenPath, screenshot)
}

// validateScreenshot compares the masked screenshot against the stored golden;
// a mismatch (or missing/mis-sized golden) fails the test and leaves the actual
// image beside the golden with the .failure extension.
func (this *AppRunner) validateScreenshot(goldenPath string, screenshot *image.RGBA) {
	this.tb.Helper()
	failurePath := this.store.FailurePath(this.snapshotFile, this.tb.Name(), this.actionCount)

	golden, err := this.store.LoadGolden(goldenPath)
	if err != nil {
		this.saveFailure(failurePath, screenshot)
		this.tb.Fatalf(
			"snapshot %d: cannot load golden %s (regenerate with `-args -update`): %v",
			this.actionCount, goldenPath, err)
		return
	}

	difference, err := this.comparer.Compare(golden, screenshot)
	if err != nil {
		this.saveFailure(failurePath, screenshot)
		this.tb.Fatalf("snapshot %d: %v (actual saved to %s)", this.actionCount, err, failurePath)
		return
	}
	if !this.comparer.Matches(difference) {
		this.saveFailure(failurePath, screenshot)
		this.tb.Errorf(
			"snapshot %d differs from %s: %s; actual saved to %s",
			this.actionCount, goldenPath, this.comparer.Describe(difference), failurePath)
		return
	}

	if err := this.store.DeleteFailure(failurePath); err != nil {
		this.tb.Errorf("snapshot %d: cannot remove stale failure %s: %v", this.actionCount, failurePath, err)
	}
}

// captureScreenshot lays out the current editor state into the headless GPU
// window and reads back the rendered pixels. It holds mu so the render
// goroutine (headed + -update) cannot interleave a layout.
func (this *AppRunner) captureScreenshot() *image.RGBA {
	this.tb.Helper()
	this.mu.Lock()
	defer this.mu.Unlock()

	var frameOps op.Ops
	gtx := layout.Context{
		Ops:         &frameOps,
		Constraints: layout.Exact(image.Pt(WindowWidth, WindowHeight)),
		Metric:      unit.Metric{PxPerDp: 1, PxPerSp: 1},
		Source:      this.router.Source(),
		Now:         time.Now(),
	}
	this.App.Layout(gtx, this.theme)

	if err := this.headlessWin.Frame(&frameOps); err != nil {
		this.tb.Fatalf("snapshot %d: headless frame failed: %v", this.actionCount, err)
	}
	screenshot := image.NewRGBA(image.Rectangle{Max: this.headlessWin.Size()})
	if err := this.headlessWin.Screenshot(screenshot); err != nil {
		this.tb.Fatalf("snapshot %d: screenshot failed: %v", this.actionCount, err)
	}
	return screenshot
}

// saveFailure best-effort persists the failing screenshot for inspection.
func (this *AppRunner) saveFailure(failurePath string, screenshot *image.RGBA) {
	this.tb.Helper()
	if err := this.store.SaveFailure(failurePath, screenshot); err != nil {
		this.tb.Logf("snapshot %d: cannot save failure image %s: %v", this.actionCount, failurePath, err)
	}
}
