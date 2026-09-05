package topologyBase_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenRoadGenerationIsDisabled_NoConnectorRoadsAreCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	roads := topologyBase.CreateConnectorZoneRoads([]string{"Gate-1", "Gate-2"}, false)

	// Assert
	assert.Nil(t, roads)
}

func TestWhenNoConnectionNamesAreProvided_NoConnectorRoadsAreCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())

	// Act
	roads := topologyBase.CreateConnectorZoneRoads(nil, true)

	// Assert
	assert.Nil(t, roads)
}

func TestWhenSingleConnectionNameIsProvided_RoadLoopsBackToItself(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	expectedRoads := []template_model.Road{
		{
			From: template_model.TypedRef{Type: "Connection", Args: []string{"Gate-1"}},
			To:   template_model.TypedRef{Type: "Connection", Args: []string{"Gate-1"}},
		},
	}

	// Act
	roads := topologyBase.CreateConnectorZoneRoads([]string{"Gate-1"}, true)

	// Assert
	assert.Equal(t, expectedRoads, roads)
}

func TestWhenMultipleConnectionNamesAreProvided_RoadsFanOutFromFirstConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	expectedRoads := []template_model.Road{
		{
			From: template_model.TypedRef{Type: "Connection", Args: []string{"Gate-1"}},
			To:   template_model.TypedRef{Type: "Connection", Args: []string{"Gate-2"}},
		},
		{
			From: template_model.TypedRef{Type: "Connection", Args: []string{"Gate-1"}},
			To:   template_model.TypedRef{Type: "Connection", Args: []string{"Gate-3"}},
		},
	}

	// Act
	roads := topologyBase.CreateConnectorZoneRoads([]string{"Gate-1", "Gate-2", "Gate-3"}, true)

	// Assert
	assert.Equal(t, expectedRoads, roads)
}

func TestWhenDuplicateConnectionNamesAreProvided_DuplicatesAreIgnored(t *testing.T) {
	t.Parallel()
	// Arrange
	topologyBase := base.NewTopologyBase(test_helpers.NewZoneFactories())
	expectedRoads := []template_model.Road{
		{
			From: template_model.TypedRef{Type: "Connection", Args: []string{"Gate-1"}},
			To:   template_model.TypedRef{Type: "Connection", Args: []string{"Gate-2"}},
		},
	}

	// Act
	roads := topologyBase.CreateConnectorZoneRoads([]string{"Gate-1", "Gate-1", "Gate-2"}, true)

	// Assert
	assert.Equal(t, expectedRoads, roads)
}
