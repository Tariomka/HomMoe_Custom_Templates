package roadFactory_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneHasObjectsFootholdAndConnection_CreatesAllRoadKinds(t *testing.T) {
	t.Parallel()
	// Arrange
	factory := zones.NewRoadFactory()

	// Act
	roads := factory.CreateOuterZoneRoads([]string{"Gate-1"}, 2, 1, true)

	// Assert
	assert.Equal(t, []entities.Road{
		{
			Type: "Stone",
			From: entities.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   entities.TypedRef{Type: "MainObject", Args: []string{"1"}},
		},
		{
			From: entities.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   entities.TypedRef{Type: "MandatoryContent", Args: []string{"name_remote_foothold_1"}},
		},
		{
			From: entities.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   entities.TypedRef{Type: "Connection", Args: []string{"Gate-1"}},
		},
	}, roads)
}
