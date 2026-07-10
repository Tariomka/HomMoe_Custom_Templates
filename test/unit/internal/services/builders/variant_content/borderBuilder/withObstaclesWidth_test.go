package borderBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenObstaclesWidthIsProvided_SetsObstaclesWidthOnBuiltBorder(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedWidth := gofakeit.Number(1, 20)
	builder := variant_content.NewBorderBuilder()

	// Act
	border := builder.WithObstaclesWidth(expectedWidth).Build()

	// Assert
	assert.Equal(t, entities.Border{ObstaclesWidth: expectedWidth}, border)
}
