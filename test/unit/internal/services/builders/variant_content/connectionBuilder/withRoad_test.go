package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenRoadFlagIsProvided_SetsRoadPointerOnBuiltConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedRoad := gofakeit.Bool()
	builder := variant_content.NewConnectionBuilder()

	// Act
	connection := builder.WithRoad(expectedRoad).Build()

	// Assert
	assert.Equal(t, entities.Connection{Road: &expectedRoad}, connection)
}
