package placementRule_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_common_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenRulesAreFlattened_EveryRuleIsCarriedAcross(t *testing.T) {
	t.Parallel()
	// Arrange
	models := []template_common_model.PlacementRule{
		{Type: "Road", Args: []any{"fromArg"}, Weight: 1},
		{Type: "Crossroads", Args: []any{"toArg"}, Weight: 2},
	}

	// Act
	entities := template_common_model.ToPlacementRuleEntities(models)

	// Assert
	expected := []template_entity.PlacementRule{
		{Type: "Road", Args: []any{"fromArg"}, Weight: 1},
		{Type: "Crossroads", Args: []any{"toArg"}, Weight: 2},
	}
	assert.Equal(t, expected, entities)
}

func TestWhenThereAreNoRules_TheNilSliceStaysNilOnTheWay(t *testing.T) {
	t.Parallel()
	// Arrange
	var models []template_common_model.PlacementRule

	// Act
	entities := template_common_model.ToPlacementRuleEntities(models)

	// Assert
	assert.Nil(t, entities)
}
