//go:build integration_test

package integration_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHandlerScroll_MovesTheLayoutPanel guards the snapshot in the gated GUI
// suite: without it, a Scroll that silently did nothing would still match a
// golden captured from an equally unscrolled frame.
func TestHandlerScroll_MovesTheLayoutPanel(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner).
		ClickLayoutAndZonesTab().
		ToggleAdvancedZoneControl()
	firstBefore, offsetBefore, ok := runner.SelectedPanelScrollPosition()
	require.True(t, ok, "the Layout & Zones panel must expose its list position")

	// Act
	handler.ScrollPanel(400)

	// Assert
	firstAfter, offsetAfter, _ := runner.SelectedPanelScrollPosition()
	assert.NotEqual(t,
		[]int{firstBefore, offsetBefore},
		[]int{firstAfter, offsetAfter})
}

// TestHandlerScroll_ScrollingUpAtTheTopStaysPut shows why the position is read
// from the real widget.List rather than accumulated from the injected deltas:
// layout.List clamps, so the two would diverge permanently after this event.
func TestHandlerScroll_ScrollingUpAtTheTopStaysPut(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	handler := integration_common.NewHandler(runner).
		ClickLayoutAndZonesTab().
		ToggleAdvancedZoneControl()

	// Act
	handler.ScrollPanel(-400)

	// Assert
	first, offset, _ := runner.SelectedPanelScrollPosition()
	assert.Equal(t, []int{0, 0}, []int{first, offset})
}
