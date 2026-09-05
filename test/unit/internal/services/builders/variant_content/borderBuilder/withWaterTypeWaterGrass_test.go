package borderBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenWaterGrassTypeIsChosen_SetsWaterGrassTypeOnBuiltBorder(t *testing.T) {
	t.Parallel()
	// Arrange
	builder := variant_content.NewBorderBuilder()

	// Act
	border := builder.WithWaterTypeWaterGrass().Build()

	// Assert
	assert.Equal(t, template_model.Border{WaterType: "water grass"}, border)
}
