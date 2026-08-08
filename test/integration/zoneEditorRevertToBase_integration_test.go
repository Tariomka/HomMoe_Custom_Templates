package integration_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newEditedSession generates a template, stamps manual positions on every zone
// and applies them, leaving a state that carries a persisted manual snapshot.
func newEditedSession(t *testing.T) *drivers.State {
	t.Helper()
	state := drivers.NewUIState(
		composition.InitializeGuiHandler(),
		composition.InitializeFileSystemHandler(),
		composition.InitializeRegenerationHandler(),
		false)
	state.Generate()
	template := state.GetLastTemplate()
	require.NotNil(t, template)
	require.NotEmpty(t, template.Variants)

	zones := append([]entities.Zone(nil), template.Variants[0].Zones...)
	for i := range zones {
		pinned := [2]float64{0.1, 0.2}
		zones[i].ManualPosition = &pinned
	}
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{
		Zones:       zones,
		Connections: template.Variants[0].Connections,
	})
	stateData := state.GetStateData()
	require.True(t, stateData.HasManualEdits(), "the edit must be applied before it can be reverted")

	return state
}

// zonesCarryManualPositions reports whether any zone of the live template is
// still pinned to a hand-placed position.
func zonesCarryManualPositions(state *drivers.State) bool {
	template := state.GetLastTemplate()
	if template == nil || len(template.Variants) == 0 {
		return false
	}
	for _, zone := range template.Variants[0].Zones {
		if zone.ManualPosition != nil {
			return true
		}
	}

	return false
}

// Previewing must not touch the live template - otherwise cancelling the
// editor leaves the base on screen with no way back to the edited layout.
func TestWhenPreviewingTheBase_TheEditedTemplateStaysOnScreen(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditedSession(t)

	// Act
	_, ok := state.PreviewBaseZones()

	// Assert
	require.True(t, ok)
	assert.True(t, zonesCarryManualPositions(state))
}

func TestWhenPreviewingTheBase_TheStoredManualEditsStay(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditedSession(t)

	// Act
	state.PreviewBaseZones()

	// Assert
	stateData := state.GetStateData()
	assert.True(t, stateData.HasManualEdits())
}

func TestWhenPreviewingTheBase_TheReturnedZonesCarryNoManualPositions(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditedSession(t)

	// Act
	base, _ := state.PreviewBaseZones()

	// Assert
	require.NotEmpty(t, base.Zones)
	for _, zone := range base.Zones {
		assert.Nilf(t, zone.ManualPosition, "zone %s came back pinned", zone.Name)
	}
}

// This is the reported defect, end to end: applying the previewed base must
// take effect on the spot, with no separate Generate() call.
func TestWhenApplyingTheBase_TheLiveTemplateLosesTheManualPositions(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditedSession(t)
	base, ok := state.PreviewBaseZones()
	require.True(t, ok)

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{
		Zones:        base.Zones,
		Connections:  base.Connections,
		RevertToBase: true,
	})

	// Assert
	assert.False(t, zonesCarryManualPositions(state))
}

func TestWhenApplyingTheBase_TheStoredManualEditsAreGone(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditedSession(t)
	base, _ := state.PreviewBaseZones()

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{
		Zones:        base.Zones,
		Connections:  base.Connections,
		RevertToBase: true,
	})

	// Assert
	stateData := state.GetStateData()
	assert.False(t, stateData.HasManualEdits())
}

// Undo is session-scoped, so applying after it re-commits the edits that were
// already applied rather than dropping them.
func TestWhenEditsAreAppliedTwice_TheManualSnapshotSurvives(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditedSession(t)
	template := state.GetLastTemplate()
	require.NotNil(t, template)

	// Act
	state.ApplyEditedZones(dtos.ZoneEditorZonesDto{
		Zones:       template.Variants[0].Zones,
		Connections: template.Variants[0].Connections,
	})

	// Assert
	stateData := state.GetStateData()
	assert.True(t, stateData.HasManualEdits())
}
