package roadBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenFromReferenceIsProvided_SetsFromOnBuiltRoad(t *testing.T) {
	// Arrange
	expectedFrom := entities.TypedRef{Type: gofakeit.Word(), Args: []string{gofakeit.Word()}}
	builder := variant_content.NewRoadBuilder()

	// Act
	road := builder.WithFrom(expectedFrom).Build()

	// Assert
	assert.Equal(t, entities.Road{From: expectedFrom}, road)
}
