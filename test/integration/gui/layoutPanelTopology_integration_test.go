//go:build integration_test && gui

package gui_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
)

// The retired Ring, Hub, Chain and Shared Web choices must not be offered.
//
//nolint:paralleltest // Driving the window needs exclusive access to the single headless GPU window.
func TestWhenTheTopologyDropdownIsOpened_ItOffersExactlyTheSupportedTopologiesInOrder(t *testing.T) {
	// Arrange
	runner := integration_common.NewAppRunner(t)
	tab := integration_common.NewHandler(runner).WithFixtureDirectory().ClickLayoutAndZonesTab()

	// Act
	labels := tab.TopologyOptionLabels()

	// Assert
	assert.Equal(t,
		[]string{"Random", "Circles", "Geometric Hub", "Square", "Geometric", "Cross", "Fractal"},
		labels)
}
