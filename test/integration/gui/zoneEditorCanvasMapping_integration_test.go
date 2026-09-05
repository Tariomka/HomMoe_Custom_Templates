//go:build integration_test && gui

package gui_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWhenEveryZonesMappedPointIsPressed_EachOneBecomesTheSelection calibrates
// the canvas-to-window mapping every other zone editor pointer test is built on.
// The mapping is one measured constant - where the space reserved for the canvas
// starts - plus the centring offset the dialog reports, so if the dialog's
// layout ever shifts, this test is what says so, rather than a dozen pointer
// tests failing for a reason none of them names.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenEveryZonesMappedPointIsPressed_EachOneBecomesTheSelection(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	zoneEditor := integration_common.NewHandler(runner).
		ClickLayoutAndZonesTab().
		SelectTopology("Square").
		OpenZoneEditor()
	require.Equal(t, config.TopologySquare, runner.CurrentState().Topology,
		"the topology dropdown did not take, so the layout is still randomised")
	positions := zoneEditor.Dialog().ZonePositions()
	require.NotEmpty(t, positions, "the canvas drew no zones to calibrate against")
	expected := make(map[string]string, len(positions))
	selected := make(map[string]string, len(positions))

	// Act
	for name := range positions {
		zoneEditor.ClickZone(name)
		expected[name] = name
		selected[name] = zoneEditor.Dialog().SelectedZone()
	}

	// Assert
	assert.Equal(t, expected, selected)
}
