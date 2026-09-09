//go:build integration_test

package integration_test

import (
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWhenTheSameLayoutIsAppliedAfterSaving_TheDocumentStaysSaved proves the
// committed snapshot is stable across the handler's rebuild: re-applying what
// the last apply produced commits nothing and must not resurrect the warning.
func TestWhenTheSameLayoutIsAppliedAfterSaving_TheDocumentStaysSaved(t *testing.T) {
	// Arrange
	state := newSavedEditedSession(t)
	template := state.GetLastTemplate()
	require.NotNil(t, template)

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{
		Zones:       append([]template_model.Zone(nil), template.Variants[0].Zones...),
		Connections: append([]template_model.Connection(nil), template.Variants[0].Connections...),
	})

	// Assert
	assert.False(t, state.IsUnsaved())
}

// A layout the user actually moved is a persisted change, so the saved
// document goes dirty again.
func TestWhenAMovedZoneIsAppliedAfterSaving_TheDocumentBecomesUnsaved(t *testing.T) {
	// Arrange
	state := newSavedEditedSession(t)
	template := state.GetLastTemplate()
	require.NotNil(t, template)
	movedZones := append([]template_model.Zone(nil), template.Variants[0].Zones...)
	movedZones[0].ManualPosition = new(data.NewVec2(0.7, 0.8))

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{
		Zones:       movedZones,
		Connections: append([]template_model.Connection(nil), template.Variants[0].Connections...),
	})

	// Assert
	assert.True(t, state.IsUnsaved())
}

// The exit confirmation is armed against the saved document; committing a new
// layout has to re-arm it the same way a scalar edit does.
func TestWhenAMovedZoneIsAppliedAfterSaving_ExitWarnsAgain(t *testing.T) {
	// Arrange
	state := newSavedEditedSession(t)
	template := state.GetLastTemplate()
	require.NotNil(t, template)
	movedZones := append([]template_model.Zone(nil), template.Variants[0].Zones...)
	movedZones[0].ManualPosition = new(data.NewVec2(0.7, 0.8))
	exitCalled := false
	state.SetOnExit(func() { exitCalled = true })

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{
		Zones:       movedZones,
		Connections: append([]template_model.Connection(nil), template.Variants[0].Connections...),
	})
	state.Exit()

	// Assert
	assert.False(t, exitCalled)
}

// newSavedEditedSession applies a manual layout through the real handler and
// saves it, leaving a clean document that already carries a manual snapshot.
func newSavedEditedSession(t *testing.T) *drivers.State {
	t.Helper()
	state := newEditedSession(t)
	state.SaveStateToFile(filepath.Join(t.TempDir(), "manualEdits.gen.json"))
	require.False(t, state.IsUnsaved(), "the session must be saved before dirty tracking can be observed")

	return state
}
