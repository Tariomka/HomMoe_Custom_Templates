package roadType_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/road_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheRoadStartsAtAConnection_ItIsAConnectionRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	types := registry.GetRoadConnectionTypeValues()
	road := template_model.Road{
		From: template_model.TypedRef{Type: types.Connection},
		To:   template_model.TypedRef{Type: types.MainObject},
	}

	// Act
	isConnectionRoad := road_helpers.IsRoadTypeConnection(road)

	// Assert
	assert.True(t, isConnectionRoad)
}

func TestWhenTheRoadEndsAtAConnection_ItIsAConnectionRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	types := registry.GetRoadConnectionTypeValues()
	road := template_model.Road{
		From: template_model.TypedRef{Type: types.Crossroads},
		To:   template_model.TypedRef{Type: types.Connection},
	}

	// Act
	isConnectionRoad := road_helpers.IsRoadTypeConnection(road)

	// Assert
	assert.True(t, isConnectionRoad)
}

// The castle-to-castle roads a zone owns must not be mistaken for the border
// roads a connection stamps, because the two are rebuilt by different services.
func TestWhenNeitherEndIsAConnection_ItIsNotAConnectionRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	types := registry.GetRoadConnectionTypeValues()
	road := template_model.Road{
		From: template_model.TypedRef{Type: types.MainObject},
		To:   template_model.TypedRef{Type: types.MainObject},
	}

	// Act
	isConnectionRoad := road_helpers.IsRoadTypeConnection(road)

	// Assert
	assert.False(t, isConnectionRoad)
}
