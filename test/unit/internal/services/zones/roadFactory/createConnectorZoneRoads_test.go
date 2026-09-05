package roadFactory_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenConnectionsContainDuplicates_CreatesDistinctFanout(t *testing.T) {
	t.Parallel()
	// Arrange
	factory := zones.NewRoadFactory()

	// Act
	roads := factory.CreateConnectorZoneRoads([]string{"Gate-1", "Gate-1", "Gate-2"}, true)

	// Assert
	assert.Equal(t, []template_model.Road{{
		From: template_model.TypedRef{Type: "Connection", Args: []string{"Gate-1"}},
		To:   template_model.TypedRef{Type: "Connection", Args: []string{"Gate-2"}},
	}}, roads)
}
