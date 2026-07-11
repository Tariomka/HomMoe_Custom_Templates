package roadBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenStoneTypeIsChosen_SetsStoneTypeOnBuiltRoad(t *testing.T) {
	t.Parallel()
	// Arrange
	builder := variant_content.NewRoadBuilder()

	// Act
	road := builder.WithStoneType().Build()

	// Assert
	assert.Equal(t, entities.Road{Type: "Stone"}, road)
}
