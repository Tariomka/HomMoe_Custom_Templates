package roadType_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/road_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/stretchr/testify/assert"
)

func TestWhenBothEndsAreMainObjects_ItIsACastleRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	types := registry.GetRoadConnectionTypeValues()
	road := template_model.Road{
		From: template_model.TypedRef{Type: types.MainObject},
		To:   template_model.TypedRef{Type: types.MainObject},
	}

	// Act
	isCastleRoad := road_helpers.IsRoadTypeCastle(road)

	// Assert
	assert.True(t, isCastleRoad)
}

func TestWhenOnlyOneEndIsAMainObject_ItIsNotACastleRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	types := registry.GetRoadConnectionTypeValues()
	road := template_model.Road{
		From: template_model.TypedRef{Type: types.MainObject},
		To:   template_model.TypedRef{Type: types.Connection},
	}

	// Act
	isCastleRoad := road_helpers.IsRoadTypeCastle(road)

	// Assert
	assert.False(t, isCastleRoad)
}
