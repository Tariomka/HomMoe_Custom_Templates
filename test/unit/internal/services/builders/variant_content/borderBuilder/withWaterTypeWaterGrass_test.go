package borderBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenWaterGrassTypeIsChosen_SetsWaterGrassTypeOnBuiltBorder(t *testing.T) {
	// Arrange
	builder := variant_content.NewBorderBuilder()

	// Act
	border := builder.WithWaterTypeWaterGrass().Build()

	// Assert
	assert.Equal(t, entities.Border{WaterType: "water grass"}, border)
}
