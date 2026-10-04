package retiredTopologies_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_topologies"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenTopologyIsRetired_ReturnsItsFormerLabel(t *testing.T) {
	t.Parallel()
	testCases := map[config.MapTopology]string{
		"Default":     "Ring",
		"HubAndSpoke": "Hub",
		"Chain":       "Chain",
		"SharedWeb":   "Shared Web",
	}
	for topology, expected := range testCases {
		t.Run(string(topology)+"_ReturnsItsFormerLabel", func(t *testing.T) {
			t.Parallel()
			// Arrange
			var actual string

			// Act
			actual, _ = common_topologies.GetRetiredTopologyLabel(topology)

			// Assert
			assert.Equal(t, expected, actual)
		})
	}
}

func TestWhenTopologyIsRetired_ReportsFound(t *testing.T) {
	t.Parallel()
	for _, topology := range []config.MapTopology{"Default", "HubAndSpoke", "Chain", "SharedWeb"} {
		t.Run(string(topology)+"_ReportsFound", func(t *testing.T) {
			t.Parallel()
			// Arrange
			var found bool

			// Act
			_, found = common_topologies.GetRetiredTopologyLabel(topology)

			// Assert
			assert.True(t, found)
		})
	}
}

func TestWhenTopologyIsNotRetired_ReportsNotFound(t *testing.T) {
	t.Parallel()
	testCases := map[string]config.MapTopology{
		"Surviving": config.TopologyCircles,
		"Unknown":   config.MapTopology("NoSuchTopology"),
		"Empty":     config.MapTopology(""),
	}
	for name, topology := range testCases {
		t.Run(name+"_ReportsNotFound", func(t *testing.T) {
			t.Parallel()
			// Arrange
			var found bool

			// Act
			_, found = common_topologies.GetRetiredTopologyLabel(topology)

			// Assert
			assert.False(t, found)
		})
	}
}
