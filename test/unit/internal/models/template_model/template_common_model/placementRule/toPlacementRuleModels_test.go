package placementRule_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_common_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenRulesAreLifted_EveryRuleIsCarriedAcross(t *testing.T) {
	t.Parallel()
	// Arrange
	entities := []template_entity.PlacementRule{
		{Type: "Road", Args: []any{"fromArg"}, Weight: 1},
		{Type: "Crossroads", Args: []any{"toArg"}, Weight: 2},
	}

	// Act
	models := template_common_model.ToPlacementRuleModels(entities)

	// Assert
	expected := []template_common_model.PlacementRule{
		{Type: "Road", Args: []any{"fromArg"}, Weight: 1},
		{Type: "Crossroads", Args: []any{"toArg"}, Weight: 2},
	}
	assert.Equal(t, expected, models)
}

// nil and empty are different on the wire - null versus [] - so lifting must
// not normalise one into the other.
func TestWhenThereAreNoRules_TheNilSliceStaysNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var entities []template_entity.PlacementRule

	// Act
	models := template_common_model.ToPlacementRuleModels(entities)

	// Assert
	assert.Nil(t, models)
}
