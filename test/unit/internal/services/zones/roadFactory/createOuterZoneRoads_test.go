package roadFactory_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
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
	assert.Equal(t, []template_model.Road{
		{
			Type: "Stone",
			From: template_model.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   template_model.TypedRef{Type: "MainObject", Args: []string{"1"}},
		},
		{
			From: template_model.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   template_model.TypedRef{Type: "MandatoryContent", Args: []string{"name_remote_foothold_1"}},
		},
		{
			From: template_model.TypedRef{Type: "MainObject", Args: []string{"0"}},
			To:   template_model.TypedRef{Type: "Connection", Args: []string{"Gate-1"}},
		},
	}, roads)
}
