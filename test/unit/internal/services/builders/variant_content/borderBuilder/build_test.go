package borderBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMultipleOptionsAreChained_ReturnsBorderWithAllAccumulatedValues(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedRadius := gofakeit.Float64Range(0.01, 1)
	expectedObstaclesWidth := gofakeit.Number(1, 20)
	expectedWaterWidth := gofakeit.Number(1, 20)
	builder := variant_content.NewBorderBuilder()

	// Act
	border := builder.
		WithCornerRadius(expectedRadius).
		WithObstaclesWidth(expectedObstaclesWidth).
		WithWaterWidth(expectedWaterWidth).
		WithWaterTypeWaterGrass().
		Build()

	// Assert
	assert.Equal(t, entities.Border{
		CornerRadius:   expectedRadius,
		ObstaclesWidth: expectedObstaclesWidth,
		WaterWidth:     expectedWaterWidth,
		WaterType:      "water grass",
	}, border)
}
