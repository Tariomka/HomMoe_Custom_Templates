package roadBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenDirtTypeIsChosen_SetsDirtTypeOnBuiltRoad(t *testing.T) {
	// Arrange
	builder := variant_content.NewRoadBuilder()

	// Act
	road := builder.WithDirtType().Build()

	// Assert
	assert.Equal(t, entities.Road{Type: "Dirt"}, road)
}
