package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenHoldCityWinConIsChosen_MarksBuiltObjectAsHoldCityWinCondition(t *testing.T) {
	t.Parallel()
	// Arrange
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithHoldCityWinCon().Build()

	// Assert
	assert.Equal(t, template_model.MainObject{HoldCityWinCon: true}, mainObject)
}
