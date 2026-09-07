//go:build integration_test

package integration_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newExitProbe returns a State with an onExit spy installed, mirroring how
// program.go wires window.Perform(system.ActionClose) into the state. The
// initial AutoRegenerate reproduces the first frame's generation, which
// establishes the snapshot change detection (and thus the unsaved flag)
// compares against.
func newExitProbe() (state *drivers.State, exitCalled *bool) {
	state = newUIState()
	state.AutoRegenerate(time.Now())
	exitCalled = new(bool)
	state.SetOnExit(func() { *exitCalled = true })
	return state, exitCalled
}

// markUnsaved makes a real edit so the state reports unsaved changes, the
// same way panel SaveToState calls do.
func markUnsaved(state *drivers.State) {
	current := state.GetStateData()
	state.UpdateState(func(s *editor_state_model.EditorState) { s.PlayerCount = current.PlayerCount + 1 })
}

// TestExit_WithCleanState_ClosesImmediately: no unsaved changes means no
// confirmation round-trip.
func TestExit_WithCleanState_ClosesImmediately(t *testing.T) {
	state, exitCalled := newExitProbe()

	state.Exit()

	assert.True(t, *exitCalled, "exit callback should run immediately when nothing is unsaved")
}

// TestExit_WithUnsavedChanges_FirstPressOnlyWarns: the first Exit press with
// unsaved changes must warn and keep the application open.
func TestExit_WithUnsavedChanges_FirstPressOnlyWarns(t *testing.T) {
	state, exitCalled := newExitProbe()
	markUnsaved(state)

	state.Exit()

	message, isError := state.GetStatus()
	require.True(t, isError)
	require.Contains(t, message, "Unsaved changes exist")
	assert.False(t, *exitCalled, "exit callback must not run on the first press with unsaved changes")
}

// TestExit_WithUnsavedChanges_SecondPressCloses: pressing Exit again right
// after the warning confirms the exit.
func TestExit_WithUnsavedChanges_SecondPressCloses(t *testing.T) {
	state, exitCalled := newExitProbe()
	markUnsaved(state)

	state.Exit()
	state.Exit()

	assert.True(t, *exitCalled, "second Exit press should confirm and close")
}

// TestExit_NewEditsAfterWarning_RearmTheGuard is the regression test for the
// confirmExit flag never resetting: edits made after the warning must demand
// a fresh confirmation instead of exiting on the next press.
func TestExit_NewEditsAfterWarning_RearmTheGuard(t *testing.T) {
	state, exitCalled := newExitProbe()
	markUnsaved(state)
	state.Exit() // warning press arms confirmExit

	markUnsaved(state) // new edits must re-arm the guard
	state.Exit()

	assert.False(t, *exitCalled, "new unsaved edits after the warning must require a fresh confirmation")
}

// TestExit_AfterSaving_ClosesWithoutWarning: saving clears both the unsaved
// flag and any pending exit confirmation.
func TestExit_AfterSaving_ClosesWithoutWarning(t *testing.T) {
	dir := t.TempDir()
	state, exitCalled := newExitProbe()
	markUnsaved(state)
	state.Exit() // warning press
	state.SaveStateToFile(filepath.Join(dir, "exit.gen.json"))

	state.Exit()

	assert.True(t, *exitCalled, "exit after saving should close without another warning")
}

// TestExit_WithoutOnExitCallback_DoesNotPanic: an unwired state (tests,
// headless drivers) must treat Exit as a safe no-op.
func TestExit_WithoutOnExitCallback_DoesNotPanic(t *testing.T) {
	state := newUIState()

	assert.NotPanics(t, func() { state.Exit() })
}
