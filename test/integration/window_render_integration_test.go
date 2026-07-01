//go:build integration_test

package integration_test

import (
	"image"
	"path/filepath"
	"testing"
	"time"

	"gioui.org/io/input"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/editor"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// headlessRenderer drives a real editor.Window through the genuine Gio layout
// pipeline without a GPU or OS window, following the pattern used by Gio's own
// widget tests (build a layout.Context backed by an input.Router, lay out, then
// commit the frame).
type headlessRenderer struct {
	window *editor.Window
	theme  *material.Theme
	router input.Router
	ops    op.Ops
}

func newHeadlessRenderer() *headlessRenderer {
	return &headlessRenderer{
		window: editor.NewWindow(),
		theme:  themes.NewTheme(),
	}
}

func (r *headlessRenderer) frame() {
	r.ops.Reset()
	gtx := layout.Context{
		Ops:         &r.ops,
		Constraints: layout.Exact(image.Pt(1600, 900)),
		Metric:      unit.Metric{PxPerDp: 1, PxPerSp: 1},
		Source:      r.router.Source(),
		Now:         time.Now(),
	}
	r.window.Layout(gtx, r.theme)
	r.router.Frame(&r.ops)
}

func (r *headlessRenderer) frames(n int) {
	for range n {
		r.frame()
	}
}

// TestWindow_RendersFramesWithoutPanic ensures the entire editor UI (toolbar,
// tabs, every panel and the live preview) lays out cleanly across many frames.
func TestWindow_RendersFramesWithoutPanic(t *testing.T) {
	renderer := newHeadlessRenderer()
	require.NotPanics(t, func() { renderer.frames(10) })

	// The first frames auto-generate a preview from the default state.
	state := renderer.window.CurrentState()
	assert.Equal(t, 2, state.PlayerCount)
}

// TestWindow_LoadReflectsInRenderedUI is the end-to-end UI test for bug #1: a
// programmatic load (equivalent to picking a file in the Load dialog) must be
// reflected by the rendered panels and must persist across subsequent rendered
// frames, where the per-frame save() previously clobbered the loaded state.
func TestWindow_LoadReflectsInRenderedUI(t *testing.T) {
	dir := t.TempDir()
	savedPath := filepath.Join(dir, "windowload.gen.json")

	// Author a distinctive saved state through the real save path.
	author := drivers.NewUIState()
	author.UpdateState(func(s *dtos.EditorStateDto) {
		s.TemplateName = "Window Loaded"
		s.PlayerCount = 5
		s.Topology = config.TopologyCross
	})
	require.NoError(t, author.SaveStateToFile(savedPath))

	renderer := newHeadlessRenderer()

	// Render baseline frames at the defaults.
	renderer.frames(3)
	require.Equal(t, 2, renderer.window.CurrentState().PlayerCount)

	// Load the saved state, mirroring the Load dialog picking a file.
	require.NoError(t, renderer.window.LoadStateFromFile(savedPath))

	// Render several more frames. Each frame runs the window's save() (panels →
	// state); the loaded values must survive instead of being overwritten.
	renderer.frames(5)

	got := renderer.window.CurrentState()
	assert.Equal(t, "Window Loaded", got.TemplateName)
	assert.Equal(t, 5, got.PlayerCount)
	assert.Equal(t, config.TopologyCross, got.Topology)
}
