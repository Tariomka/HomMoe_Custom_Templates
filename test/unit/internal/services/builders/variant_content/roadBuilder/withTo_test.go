package roadBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenToReferenceIsProvided_SetsToOnBuiltRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedTo := entities.TypedRef{Type: gofakeit.Word(), Args: []string{gofakeit.Word()}}
	builder := variant_content.NewRoadBuilder()

	// Act
	road := builder.WithTo(expectedTo).Build()

	// Assert
	assert.Equal(t, entities.Road{To: expectedTo}, road)
}
