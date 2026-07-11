package typedRefBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenBiomeFromListReferenceIsBuilt_SetsFromListTypeWithArguments(t *testing.T) {
	t.Parallel()
	// Arrange
	firstBiome := gofakeit.Word()
	secondBiome := gofakeit.Word()
	builder := variant_content.NewRefBuilder()

	// Act
	reference := builder.BuildBiomeFromListType(firstBiome, secondBiome)

	// Assert
	assert.Equal(t, entities.TypedRef{Type: "FromList", Args: []string{firstBiome, secondBiome}}, reference)
}
