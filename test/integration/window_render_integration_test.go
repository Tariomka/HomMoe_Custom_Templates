//go:build integration_test

package integration_test

import (
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// renderFrames lays out and commits the given number of frames on the runner.
func renderFrames(runner *integration_common.AppRunner, frameCount int) {
	for range frameCount {
		runner.NextFrame()
	}
}

// TestWindow_RendersFramesWithoutPanic ensures the entire editor UI (toolbar,
// tabs, every panel and the live preview) lays out cleanly across many frames.
func TestWindow_RendersFramesWithoutPanic(t *testing.T) {
	runner := integration_common.NewAppRunner(t)
	require.NotPanics(t, func() { renderFrames(runner, 10) })

	// The first frames auto-generate a preview from the default state.
	state := runner.CurrentState()
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
	author := newUIState()
	author.UpdateState(func(s *dtos.EditorStateDto) {
		s.TemplateName = "Window Loaded"
		s.PlayerCount = 5
		s.Topology = config.TopologyCross
	})
	author.SaveStateToFile(savedPath)
	message, irError := author.GetStatus()
	require.False(t, irError)
	assert.Equal(t, "Saved "+savedPath, message)

	runner := integration_common.NewAppRunner(t)

	// Render baseline frames at the defaults.
	renderFrames(runner, 3)
	require.Equal(t, 2, runner.CurrentState().PlayerCount)

	// Load the saved state, mirroring the Load dialog picking a file.
	runner.LoadStateFromFile(savedPath)
	message, irError = runner.Status()
	require.False(t, irError)
	assert.Equal(t, "Loaded "+savedPath, message)

	// Render several more frames. Each frame runs the window's save() (panels →
	// state); the loaded values must survive instead of being overwritten.
	renderFrames(runner, 5)

	got := runner.CurrentState()
	assert.Equal(t, "Window Loaded", got.TemplateName)
	assert.Equal(t, 5, got.PlayerCount)
	assert.Equal(t, config.TopologyCross, got.Topology)
}
