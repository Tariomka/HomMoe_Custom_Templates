//go:build integration_test

// Package integration_common holds shared helpers for the gated integration and
// performance suites (build tag integration_test). AppRunner drives a real
// editor.Window in one of two modes selected at runtime:
//
//   - default: no window, no display and no app.Main - safe for CI.
//   - headed: pass "headed" after -args (go test ... -args headed) so a real
//     gioui.org/app.Window opens on screen and renders the UI while the test
//     drives it.
//
// Both modes drive the UI through a single shared input.Router that lays out the
// real editor.Window and injects synthetic pointer/key events (a real app.Window
// exposes no event-injection API, so even in windowed mode the clicks are
// processed through this router while the on-screen window passively mirrors the
// resulting state). AppRunner hides the difference: write a test once and it
// behaves identically in either mode, the windowed mode merely also renders.
package integration_common

import (
	"image"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/gpu/headless"
	"gioui.org/io/input"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/editor"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common/snapshot"
)

const (
	// WindowWidth is the fixed layout width AppRunner drives the editor at.
	WindowWidth = 1600
	// WindowHeight is the fixed layout height AppRunner drives the editor at.
	WindowHeight = 900
)

// AppRunner drives a real editor.Window for a test or benchmark. In windowed
// mode it also owns an app.Window that renders the same editor.Window on a
// background goroutine; in headless mode window is nil and Start/Stop are
// no-ops.
//
// All synthetic input is applied through the aux input.Router (router/ops) by
// laying out the shared editor.Window. Because that widget tree is also laid out
// by the render goroutine in windowed mode, every access to it is guarded by mu.
type AppRunner struct {
	App   *editor.Window
	theme *material.Theme

	// aux router: the single source of truth for synthetic input in both modes.
	router input.Router
	ops    op.Ops

	window *app.Window
	mu     sync.Mutex    // guards App (laid out by both goroutines)
	done   chan struct{} // closed when the render goroutine exits
	quit   atomic.Bool   // requests a graceful window close
	winErr error         // set if the window is destroyed with an error

	// renderDelay, when > 0, pauses briefly after each on-screen invalidate so a
	// human can watch the UI change. It only applies in windowed mode and is 0
	// by default so headless timings stay clean.
	renderDelay time.Duration

	// snapshot verification (see appRunnerSnapshots.go); nil snapshotTest means
	// snapshotting is disabled (the default, e.g. for benchmarks).
	snapshotTest *testing.T
	snapshotFile string
	actionCount  int
	headlessWin  *headless.Window
	masker       snapshot.Masker
	comparer     snapshot.Comparer
	store        snapshot.Store
}

// NewAppRunner builds a runner. The window is created only in windowed mode; a
// nil window marks the runner headless.
func NewAppRunner() *AppRunner {
	runner := &AppRunner{
		App:   editor.NewWindow(),
		theme: themes.NewTheme(),
	}
	if !IsHeadless() {
		runner.window = new(app.Window)
	}
	return runner
}

func (this *AppRunner) SetRenderDelay(delay time.Duration) {
	this.renderDelay = delay
}

// Start launches the on-screen render loop. It is a no-op in headless mode.
func (this *AppRunner) Start() {
	if this.window == nil {
		return
	}

	this.window.Option(
		app.Title("Olden Era - Custom Templates (test runner)"),
		app.Size(unit.Dp(WindowWidth), unit.Dp(WindowHeight)))
	this.done = make(chan struct{})
	go this.runWindow()
}

// Stop closes the window and waits for the render loop to exit. No-op headless.
func (this *AppRunner) Stop() {
	if this.window == nil {
		return
	}

	this.quit.Store(true)
	this.window.Invalidate()
	<-this.done
}

// NextFrame lays out and commits a single frame (and mirrors it on screen). Useful
// for warm-up and for letting a freshly selected panel lay out fully.
func (this *AppRunner) NextFrame() {
	this.mu.Lock()
	this.frameLocked()
	this.mu.Unlock()
	this.invalidate()
}

// ClickAt injects a synthetic touch tap (press + release) at clickPoint. The
// leading frame registers the input areas, the trailing frame processes the tap.
// Both run under one lock so a render cannot observe a half-applied click.
func (this *AppRunner) ClickAt(point f32.Point) {
	this.mu.Lock()
	this.frameLocked()
	this.router.Queue(
		pointer.Event{Kind: pointer.Press, Source: pointer.Touch, Position: point},
		pointer.Event{Kind: pointer.Release, Source: pointer.Touch, Position: point},
		pointer.Event{Kind: pointer.Move, Source: pointer.Touch, Position: f32.Point{}})
	this.frameLocked()
	this.mu.Unlock()
	this.invalidate()
	this.verifySnapshot()
}

// MoveTo injects a synthetic touch move at point. The leading frame registers
// the input areas, the trailing frame processes the move. Both run under one
// lock so a render cannot observe a half-applied move.
func (this *AppRunner) MoveTo(point f32.Point) {
	this.mu.Lock()
	this.frameLocked()
	this.router.Queue(pointer.Event{Kind: pointer.Move, Source: pointer.Touch, Position: point})
	this.frameLocked()
	this.mu.Unlock()
	this.invalidate()
}

// DragTo injects a synthetic drag from one point to another (press, a series of
// interpolated moves, release), driving gesture.Drag widgets such as sliders.
func (this *AppRunner) DragTo(from, to image.Point) {
	const steps = 8
	this.mu.Lock()
	this.frameLocked()
	this.router.Queue(
		pointer.Event{Kind: pointer.Press, Source: pointer.Touch, Position: f32.Pt(float32(from.X), float32(from.Y))})
	this.frameLocked()
	for i := 1; i <= steps; i++ {
		x := from.X + (to.X-from.X)*i/steps
		y := from.Y + (to.Y-from.Y)*i/steps
		this.router.Queue(pointer.Event{
			Kind:     pointer.Move,
			Source:   pointer.Touch,
			Position: f32.Pt(float32(x), float32(y)),
		})
		this.frameLocked()
	}
	this.router.Queue(pointer.Event{
		Kind:     pointer.Release,
		Source:   pointer.Touch,
		Position: f32.Pt(float32(to.X), float32(to.Y)),
	})
	this.frameLocked()
	this.mu.Unlock()
	this.invalidate()
	this.verifySnapshot()
}

// InputText injects text into the currently focused widget (focus a field first,
// e.g. with ClickAt). The text replaces the focused editor's current selection.
func (this *AppRunner) InputText(text string) {
	this.mu.Lock()
	this.frameLocked()
	this.router.Queue(key.EditEvent{Text: text})
	this.frameLocked()
	this.mu.Unlock()
	this.invalidate()
	this.verifySnapshot()
}

// SelectedTabIndex returns the editor's selected tab, taken under the lock so it
// is safe to read while the render goroutine is active.
func (this *AppRunner) SelectedTabIndex() int {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.App.SelectedTabIndex()
}

// TabCount returns the number of editor tabs (lock-guarded).
func (this *AppRunner) TabCount() int {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.App.TabCount()
}

// DialogsOpen reports whether a modal dialog is open (lock-guarded).
func (this *AppRunner) DialogsOpen() bool {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.App.DialogsOpen()
}

// CloseTopDialog dismisses the top-most modal dialog (lock-guarded).
func (this *AppRunner) CloseTopDialog() {
	this.mu.Lock()
	this.App.CloseTopDialog()
	this.mu.Unlock()
}

// CurrentState returns the editor's current state snapshot (lock-guarded).
func (this *AppRunner) CurrentState() dtos.EditorStateDto {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.App.CurrentState()
}

// LoadStateFromFile loads an editor state file programmatically, mirroring the
// Load dialog picking a file (lock-guarded).
func (this *AppRunner) LoadStateFromFile(path string) {
	this.mu.Lock()
	this.App.LoadStateFromFile(path)
	this.mu.Unlock()
}

// SaveStateToFile saves the current editor state to a file programmatically,
// mirroring the Save dialog (lock-guarded).
func (this *AppRunner) SaveStateToFile(path string) {
	this.mu.Lock()
	this.App.SaveStateToFile(path)
	this.mu.Unlock()
}

// Status returns the state driver's status message and error flag (lock-guarded).
func (this *AppRunner) Status() (string, bool) {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.App.GetStateDriver().GetStatus()
}

// runWindow is the real app.Window event loop, mirroring gui.eventLoop. It only
// renders; all UI logic is driven through the aux router on the test goroutine.
// It performs a graceful close once Stop requests it.
func (this *AppRunner) runWindow() {
	defer close(this.done)
	var ops op.Ops
	for {
		switch event := this.window.Event().(type) {
		case app.DestroyEvent:
			this.winErr = event.Err
			return
		case app.FrameEvent:
			if this.quit.Load() {
				this.window.Perform(system.ActionClose)
			}
			this.mu.Lock()
			gtx := app.NewContext(&ops, event)
			this.App.Layout(gtx, this.theme)
			this.mu.Unlock()
			event.Frame(gtx.Ops)
		}
	}
}

// invalidate asks the on-screen window to redraw the current editor state. It is
// a no-op headless. window.Invalidate is safe to call from another goroutine.
func (this *AppRunner) invalidate() {
	if this.window == nil {
		return
	}

	this.window.Invalidate()
	if this.renderDelay > 0 {
		time.Sleep(this.renderDelay)
	}
}

// frameLocked lays out one editor.Window frame through the aux router. The caller
// must hold mu. This is the headless equivalent of an app.FrameEvent and, in
// windowed mode, the pass that actually processes injected input.
func (this *AppRunner) frameLocked() {
	this.ops.Reset()
	gtx := layout.Context{
		Ops:         &this.ops,
		Constraints: layout.Exact(image.Pt(WindowWidth, WindowHeight)),
		Metric:      unit.Metric{PxPerDp: 1, PxPerSp: 1},
		Source:      this.router.Source(),
		Now:         time.Now(),
	}
	this.App.Layout(gtx, this.theme)
	this.router.Frame(&this.ops)
}
